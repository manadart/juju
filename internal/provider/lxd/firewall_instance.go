// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"context"
	"maps"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/canonical/lxd/shared/api"
	"github.com/juju/collections/set"
	"github.com/juju/errors"

	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/network/firewall"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/internal/container/lxd"
)

// prepareContainerFirewall materialises only the OVN device overrides, even
// for the single eth0 profile case where cloud-init needs no custom network
// configuration. The container must start with its ACLs already attached.
func (env *environ) prepareContainerFirewall(ctx context.Context, spec *lxd.ContainerSpec, controllerUUID string, devices map[string]map[string]string) error {
	overrides, err := env.prepareOVNFirewall(ctx, env.server(), spec.Name, controllerUUID, devices)
	if err != nil {
		return errors.Annotate(err, "preparing OVN firewall")
	}
	if len(overrides) > 0 {
		if spec.Devices == nil {
			spec.Devices = make(map[string]map[string]string)
		}
		maps.Copy(spec.Devices, overrides)
	}
	return nil
}

// ensureInstanceFirewall attaches ACLs to an existing instance, using ETags
// to preserve concurrent instance configuration changes. Returning nil means
// that the instance has no OVN NICs and needs no provider firewall.
func (env *environ) ensureInstanceFirewall(ctx context.Context, instanceName string) (networkACLServer, error) {
	if env.Config().FirewallMode() == config.FwNone {
		return nil, ctx.Err()
	}
	srv := env.server()
	for range aclUpdateAttempts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		inst, etag, err := srv.GetInstance(instanceName)
		if err != nil {
			return nil, errors.Annotatef(err, "reading firewall instance %q", instanceName)
		}
		nics, err := ovnFirewallNICs(ctx, srv, inst.ExpandedDevices)
		if err != nil || len(nics) == 0 {
			return nil, err
		}
		if inst.Config[jujuACLModelKey] != env.uuid {
			return nil, errors.Errorf("instance %q does not belong to model %q", instanceName, env.uuid)
		}
		overrides, err := env.prepareOVNNICFirewall(ctx, srv, instanceName,
			inst.Config[jujuACLControllerKey], nics)
		if err != nil || len(overrides) == 0 {
			return nil, err
		}
		client, err := aclServer(srv)
		if err != nil {
			return nil, err
		}
		changed := false
		put := inst.Writable()
		put.Devices = maps.Clone(inst.Devices)
		if put.Devices == nil {
			put.Devices = make(map[string]map[string]string)
		}
		for _, name := range slices.Sorted(maps.Keys(overrides)) {
			device := overrides[name]
			if reflect.DeepEqual(inst.ExpandedDevices[name], device) {
				continue
			}
			put.Devices[name] = device
			changed = true
		}
		if !changed {
			return client, nil
		}
		op, err := client.UpdateInstance(instanceName, put, etag)
		if err == nil {
			err = op.WaitContext(ctx)
		}
		if api.StatusErrorCheck(err, http.StatusPreconditionFailed) {
			continue
		}
		return client, errors.Annotatef(err, "attaching OVN firewall to instance %q", instanceName)
	}
	return nil, errors.Errorf("instance %q kept changing during firewall attachment", instanceName)
}

// prepareOVNFirewall creates the required ACLs and returns device overrides
// to install before starting an instance. Only OVN NICs are changed, including
// NICs inherited from profiles. Network-wide ACL assignments are not changed.
func (env *environ) prepareOVNFirewall(ctx context.Context, srv Server, instanceName, controllerUUID string, devices map[string]map[string]string) (map[string]map[string]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if env.Config().FirewallMode() == config.FwNone {
		return nil, nil
	}
	nics, err := ovnFirewallNICs(ctx, srv, devices)
	if err != nil || len(nics) == 0 {
		return nil, err
	}
	return env.prepareOVNNICFirewall(ctx, srv, instanceName, controllerUUID, nics)
}

func (env *environ) prepareOVNNICFirewall(ctx context.Context, srv Server, instanceName, controllerUUID string, nics map[string]map[string]string) (map[string]map[string]string, error) {
	client, err := aclServer(srv)
	if err != nil {
		return nil, err
	}
	modelName := modelACLName(env.uuid)
	owner := aclMetadata(env.uuid, controllerUUID, "", "model")
	var sshRules firewall.IngressRules
	if sources := env.Config().SSHAllow(); len(sources) > 0 {
		sshRules = append(sshRules, firewall.NewIngressRule(network.MustParsePortRange("22/tcp"), sources...))
	}
	ingress, err := ingressACLRules(sshRules)
	if err != nil {
		return nil, err
	}
	if err := ensureNetworkACL(ctx, client, modelName, owner, ingress); err != nil {
		return nil, err
	}
	if err := updateNetworkACL(ctx, client, modelName, owner, func(acl *api.NetworkACLPut) error {
		setModelACLBaseline(acl, modelName)
		return nil
	}); err != nil {
		return nil, err
	}

	instanceACL := machineACLName(env.uuid, instanceName)
	instanceOwner := aclMetadata(env.uuid, controllerUUID, instanceName, "machine")
	if env.Config().FirewallMode() == config.FwGlobal {
		instanceACL = globalACLName(env.uuid)
		instanceOwner = aclMetadata(env.uuid, controllerUUID, "", "global")
	}
	if err := ensureNetworkACL(ctx, client, instanceACL, instanceOwner, nil); err != nil {
		return nil, err
	}
	for name, nic := range nics {
		device := maps.Clone(nic)
		assigned := set.NewStrings(modelName, instanceACL)
		for _, acl := range strings.Split(device["security.acls"], ",") {
			if acl = strings.TrimSpace(acl); acl != "" {
				assigned.Add(acl)
			}
		}
		device["security.acls"] = strings.Join(assigned.SortedValues(), ",")
		device["security.acls.default.ingress.action"] = "drop"
		device["security.acls.default.egress.action"] = "drop"
		nics[name] = device
	}
	return nics, nil
}

// ovnFirewallNICs resolves the actual managed network type; neither the NIC
// name nor the presence of a network property identifies an OVN interface.
func ovnFirewallNICs(ctx context.Context, srv Server, devices map[string]map[string]string) (map[string]map[string]string, error) {
	hasNetworks := false
	for _, device := range devices {
		if device["type"] == "nic" && device["network"] != "" {
			hasNetworks = true
			break
		}
	}
	if !hasNetworks {
		return nil, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	networks, err := srv.GetNetworks()
	if err != nil {
		return nil, errors.Annotate(err, "retrieving networks for OVN firewall")
	}
	networkTypes := make(map[string]string)
	for _, net := range networks {
		networkTypes[net.Name] = net.Type
	}
	result := make(map[string]map[string]string)
	for name, device := range devices {
		if device["type"] != "nic" || device["network"] == "" {
			continue
		}
		networkType, found := networkTypes[device["network"]]
		if !found {
			return nil, errors.NotFoundf("network %q for firewall NIC %q", device["network"], name)
		}
		if networkType == networkTypeOVN {
			result[name] = device
		}
	}
	return result, nil
}

func aclMetadata(modelUUID, controllerUUID, instanceName, kind string) map[string]string {
	metadata := map[string]string{
		jujuACLModelKey: modelUUID,
		jujuACLKindKey:  kind,
	}
	if controllerUUID != "" {
		metadata[jujuACLControllerKey] = controllerUUID
	}
	if instanceName != "" {
		metadata[jujuACLInstanceKey] = instanceName
	}
	return metadata
}

func setModelACLBaseline(acl *api.NetworkACLPut, modelName string) {
	var ingress, egress []api.NetworkACLRule
	for _, protocol := range []string{"tcp", "udp", "icmp4", "icmp6"} {
		ingress = append(ingress, api.NetworkACLRule{
			Action: "allow", State: "enabled", Protocol: protocol,
			Source: modelName, Description: jujuBaselineDescription,
		})
	}
	// LXD evaluates ingress and egress in the same OVN ACL stage. An
	// unrestricted egress allow would also permit traffic to unexposed
	// machines in other models on the switch. Restrict this permission to
	// the router; same-model switch traffic matches the ingress rules above.
	// Explicit allow rules are stateful, including their return traffic.
	egress = append(egress, api.NetworkACLRule{
		Action: "allow", State: "enabled", Destination: "@external",
		Description: jujuBaselineDescription,
	})
	for _, rule := range acl.Ingress {
		if rule.Description != jujuBaselineDescription {
			ingress = append(ingress, rule)
		}
	}
	for _, rule := range acl.Egress {
		if rule.Description != jujuBaselineDescription {
			egress = append(egress, rule)
		}
	}
	acl.Ingress, acl.Egress = ingress, egress
}
