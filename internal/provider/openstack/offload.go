// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package openstack

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-goose/goose/v5/client"
	goosehttp "github.com/go-goose/goose/v5/http"
	"github.com/go-goose/goose/v5/neutron"
	"github.com/go-goose/goose/v5/nova"
	"github.com/juju/errors"

	"github.com/juju/juju/core/network"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/internal/cloudconfig/cloudinit"
)

const offloadPortDescription = "Port created by Juju for network-offload"

// offloadedNetworksForInstance creates all ports before asking Nova to boot.
// Their names retain model and machine ownership after a controller restart.
func (e *Environ) offloadedNetworksForInstance(
	ctx context.Context,
	args environs.StartInstanceParams,
	cloudCfg cloudinit.NetworkingConfig,
) (_ []nova.ServerNetworks, err error) {
	if !e.supportsNeutron() || e.ecfg().useOpenstackGBP() {
		return nil, errors.New("network-offload=true requires Neutron networking without OpenStack GBP")
	}
	if args.Constraints.HasVirtType() && *args.Constraints.VirtType == "lxd" {
		return nil, errors.New("network-offload=true is not supported for LXD instances")
	}

	networks, err := e.networksForModel(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(networks) == 0 {
		return nil, errors.New("network-offload=true requires an OpenStack network")
	}

	// Preserve the normal network/space selection. Without space bindings,
	// create one port per selected network and let Neutron allocate its IPs.
	var ports []neutron.PortV2
	if len(args.SubnetsToZones) == 0 {
		for _, net := range networks {
			ports = append(ports, neutron.PortV2{NetworkId: net.Id})
		}
	} else {
		subnets, err := subnetInZone(ctx, args.AvailabilityZone, args.SubnetsToZones)
		if err != nil {
			return nil, errors.Trace(err)
		}
		for _, subnetID := range subnets {
			net, err := networkForSubnet(networks, subnetID)
			if err != nil {
				return nil, errors.Trace(err)
			}
			ports = append(ports, neutron.PortV2{
				NetworkId: net.Id,
				FixedIPs:  []neutron.PortFixedIPsV2{{SubnetID: subnetID.String()}},
			})
		}
	}
	if len(ports) == 0 {
		return nil, errors.New("no networks satisfy network-offload=true")
	}

	for _, port := range ports {
		net, err := e.neutron().GetNetworkV2(port.NetworkId)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if net.PortSecurityEnabled != nil && !*net.PortSecurityEnabled {
			return nil, errors.Errorf("network %q has port security disabled; network-offload requires port security", net.Id)
		}
	}

	// A failed previous attempt may have left unbound ports. Only reap ports
	// with our exact ownership prefix and description, never attached ports.
	if err := e.deleteUnusedOffloadPorts(ctx, args.InstanceConfig.MachineId); err != nil {
		return nil, errors.Trace(err)
	}

	var created []nova.ServerNetworks
	defer func() {
		if err == nil {
			return
		}
		e.cleanupOffloadAttempt(ctx, created, args.InstanceConfig.MachineId)
	}()

	groupNames, err := e.firewaller.SetUpGroups(ctx, args.ControllerUUID, args.InstanceConfig.MachineId)
	if err != nil {
		return nil, errors.Annotate(err, "cannot set up groups for offloaded ports")
	}
	groupIDs := make([]string, 0, len(groupNames))
	for _, name := range groupNames {
		groups, err := e.neutron().SecurityGroupByNameV2(name)
		if err != nil {
			return nil, errors.Trace(err)
		}
		var matches []string
		for _, group := range groups {
			if group.TenantId == e.client().TenantId() && group.Name == name {
				matches = append(matches, group.Id)
			}
		}
		if len(matches) != 1 {
			return nil, errors.Errorf("expected one security group %q in the current project, found %d", name, len(matches))
		}
		groupIDs = append(groupIDs, matches[0])
	}

	netInfo := make(network.InterfaceInfos, 0, len(ports))
	for i, port := range ports {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		port.Name = generateUniquePortName(e.modelUUID + "-offload-" + args.InstanceConfig.MachineId)
		port.Description = offloadPortDescription
		port.PortSecurityEnabled = true
		port.SecurityGroups = groupIDs
		createdPort, err := e.createOffloadedPort(ctx, port)
		if err != nil {
			return nil, errors.Annotatef(err, "creating offloaded port on network %q", port.NetworkId)
		}
		// Record ownership before validating the returned binding attributes,
		// so a cloud that ignores them does not leave an orphaned port.
		created = append(created, nova.ServerNetworks{NetworkId: port.NetworkId, PortId: createdPort.Id})
		if createdPort.VNICType != "direct" || !slices.Contains(createdPort.BindingProfile.Capabilities, "switchdev") {
			return nil, errors.Errorf("port %q did not retain the requested direct/switchdev binding", createdPort.Id)
		}
		if !createdPort.PortSecurityEnabled {
			return nil, errors.Errorf("port %q did not retain port security", createdPort.Id)
		}
		actualGroups := slices.Clone(createdPort.SecurityGroups)
		expectedGroups := slices.Clone(groupIDs)
		slices.Sort(actualGroups)
		slices.Sort(expectedGroups)
		if !slices.Equal(actualGroups, expectedGroups) {
			return nil, errors.Errorf("port %q did not retain Juju security groups", createdPort.Id)
		}

		ips := make([]string, len(createdPort.FixedIPs))
		for j, ip := range createdPort.FixedIPs {
			ips[j] = ip.IPAddress
		}
		netInfo = append(netInfo, network.InterfaceInfo{
			InterfaceName: fmt.Sprintf("eth%d", i),
			MACAddress:    createdPort.MACAddress,
			Addresses:     network.NewMachineAddresses(ips).AsProviderAddresses(),
			ConfigType:    network.ConfigDHCP,
			Origin:        network.OriginProvider,
		})
	}
	if err := cloudCfg.AddNetworkConfig(netInfo); err != nil {
		return nil, errors.Annotate(err, "configuring offloaded interfaces")
	}
	return created, nil
}

// offloadedPort adds the port-binding extension fields missing from goose's
// PortV2. Keep the extension at the OpenStack API boundary.
type offloadedPort struct {
	neutron.PortV2
	VNICType       string `json:"binding:vnic_type"`
	BindingProfile struct {
		Capabilities []string `json:"capabilities"`
	} `json:"binding:profile"`
	// An explicit empty list prevents Neutron attaching its default group.
	SecurityGroups []string `json:"security_groups"`
}

func (e *Environ) createOffloadedPort(ctx context.Context, port neutron.PortV2) (*offloadedPort, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	request := struct {
		Port offloadedPort `json:"port"`
	}{Port: offloadedPort{PortV2: port, VNICType: "direct", SecurityGroups: port.SecurityGroups}}
	request.Port.BindingProfile.Capabilities = []string{"switchdev"}
	var response struct {
		Port offloadedPort `json:"port"`
	}
	data := goosehttp.RequestData{
		ReqValue:       request,
		RespValue:      &response,
		ExpectedStatus: []int{http.StatusCreated},
	}
	if err := e.client().SendRequest(client.POST, "network", "v2.0", neutron.ApiPortsV2, &data); err != nil {
		return nil, errors.Trace(err)
	}
	if response.Port.Id == "" {
		return nil, errors.New("Neutron returned an offloaded port without an ID")
	}
	return &response.Port, nil
}

func (e *Environ) cleanupOffloadAttempt(ctx context.Context, ports []nova.ServerNetworks, machineID string) {
	for _, net := range ports {
		port, err := e.neutron().PortByIdV2(net.PortId)
		if isNotFoundError(err) {
			continue
		}
		if err != nil {
			logger.Errorf(ctx, "cannot inspect offloaded port %q for cleanup: %v", net.PortId, err)
			continue
		}
		// Failed server deletion must not cause us to remove a live NIC.
		// Normal termination captures and deletes its ports after the server.
		if port.DeviceId != "" {
			logger.Warningf(ctx, "leaving offloaded port %q attached to instance %q", port.Id, port.DeviceId)
			continue
		}
		if err := e.neutron().DeletePortV2(port.Id); err != nil && !isNotFoundError(err) {
			logger.Errorf(ctx, "cannot clean up offloaded port %q for machine %q: %v", port.Id, machineID, err)
		}
	}
	if err := e.firewaller.DeleteMachineGroup(ctx, machineID); err != nil {
		logger.Errorf(ctx, "cannot clean up security groups for machine %q: %v", machineID, err)
	}
}

// deleteUnusedOffloadPorts recovers unbound ports left by interrupted attempts.
// An empty machine ID selects all such ports owned by this model.
func (e *Environ) deleteUnusedOffloadPorts(ctx context.Context, machineID string) error {
	prefix := "juju-" + e.modelUUID + "-offload-"
	if machineID != "" {
		prefix += machineID + "-"
	}
	ports, err := e.neutron().ListPortsV2(projectIDFilter(e.client().TenantId()))
	if err != nil {
		return errors.Trace(err)
	}
	for _, port := range ports {
		if err := ctx.Err(); err != nil {
			return err
		}
		if port.DeviceId != "" || port.Description != offloadPortDescription || !strings.HasPrefix(port.Name, prefix) {
			continue
		}
		if err := e.neutron().DeletePortV2(port.Id); err != nil && !isNotFoundError(err) {
			return errors.Annotatef(err, "deleting unused offloaded port %q", port.Id)
		}
	}
	return nil
}
