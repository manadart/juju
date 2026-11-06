// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"context"
	"fmt"
	"maps"
	"net"
	"slices"
	"time"

	"github.com/canonical/lxd/shared/api"
	"github.com/juju/clock"
	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"github.com/juju/retry"

	"github.com/juju/juju/internal/container/lxd"
)

const (
	jujuInstanceForwardKey = lxd.UserNamespacePrefix + "juju-instance"
	jujuDeviceForwardKey   = lxd.UserNamespacePrefix + "juju-device"
)

// ensureNetworkForwards creates ingress addresses before provisioning returns.
func (env *environ) ensureNetworkForwards(ctx context.Context, srv Server, container *lxd.Container) error {
	pools, err := newConfig(env.Config()).bridgeForwardPools()
	if err != nil {
		return errors.Trace(err)
	}
	if err := ensureOVNNetworkForwards(ctx, srv, container, env.clock); err != nil {
		return errors.Trace(err)
	}
	return ensureBridgeNetworkForwards(ctx, srv, container, pools, env.clock)
}

// ensureOVNNetworkForwards allocates external addresses for each OVN
// interface. Ownership tags allow repeated starts to reuse or replace forwards.
// Only families enabled on the guest network and permitted for external
// allocation are requested. All selected guest addresses must become ready.
func ensureOVNNetworkForwards(ctx context.Context, srv Server, container *lxd.Container, clk clock.Clock) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// ExpandedDevices includes NICs inherited from profiles. Avoid querying
	// networks when the instance has no NICs attached to a named network.
	nics := make(map[string]map[string]string)
	for name, device := range container.ExpandedDevices {
		if device["type"] == "nic" && lxd.NetworkName(device) != "" {
			nics[name] = device
		}
	}
	if len(nics) == 0 {
		return nil
	}
	networks, err := srv.GetNetworks()
	if err != nil {
		return errors.Annotate(err, "retrieving networks")
	}
	ovnNetworks := make(map[string]api.Network)
	for _, network := range networks {
		if network.Type == networkTypeOVN {
			ovnNetworks[network.Name] = network
		}
	}
	attachedNetworks := make(map[string]api.Network)
	for name, device := range nics {
		networkName := lxd.NetworkName(device)
		if network, ok := ovnNetworks[networkName]; ok {
			attachedNetworks[networkName] = network
		} else {
			delete(nics, name)
		}
	}
	if len(nics) == 0 {
		return nil
	}

	// Check capabilities only for OVN, so bridge instances remain compatible
	// with older servers. Automatic listen-address allocation is also required.
	for _, extension := range []string{"network_forward", "network_allocate_external_ips"} {
		if !srv.HasExtension(extension) {
			return errors.NotSupportedf("OVN network forwards require the %q LXD API extension; upgrade LXD", extension)
		}
	}

	families, err := ovnForwardFamilies(ctx, srv, attachedNetworks)
	if err != nil {
		return errors.Trace(err)
	}

	// Wait for all guest addresses before allocating any external addresses.
	logger.Debugf(ctx, "Waiting for OVN network addresses for instance %q", container.Name)
	state, err := waitForForwardAddresses(ctx, srv, container.Name, nics, families, "OVN", clk)
	if err != nil {
		return errors.Trace(err)
	}
	// Each interface/family reconciles a disjoint set of owned forwards, so
	// one snapshot per network is sufficient for the entire instance.
	forwards := make(map[string][]api.NetworkForward)
	for _, name := range slices.Sorted(maps.Keys(nics)) {
		device := nics[name]
		iface := ovnInterfaceName(name, device)
		networkName := lxd.NetworkName(device)
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, ok := forwards[networkName]; !ok {
			forwards[networkName], err = srv.GetNetworkForwards(networkName)
			if err != nil {
				return errors.Annotatef(err, "retrieving forwards for OVN network %q", networkName)
			}
		}
		for _, family := range families[networkName] {
			address := interfaceForwardAddress(state.Network[iface], family)
			if err := ensureOVNNetworkForward(ctx, srv, networkName, container.Name, iface, address, forwards[networkName]); err != nil {
				return errors.Trace(err)
			}
		}
	}
	return nil
}

func waitForForwardAddresses(ctx context.Context, srv Server, instanceName string, nics map[string]map[string]string, families map[string][]ovnForwardFamily, networkType string, clk clock.Clock) (*api.InstanceState, error) {
	var state *api.InstanceState
	err := retry.Call(retry.CallArgs{
		Clock:    clk,
		Stop:     ctx.Done(),
		Delay:    time.Second,
		Attempts: 60,
		IsFatalError: func(err error) bool {
			return !errors.Is(err, errors.NotAssigned)
		},
		Func: func() error {
			// Stop only interrupts retry's wait, not the next invocation.
			if err := ctx.Err(); err != nil {
				return err
			}
			var err error
			state, _, err = srv.GetInstanceState(instanceName)
			if err != nil {
				return errors.Trace(err)
			}
			for name, device := range nics {
				iface := ovnInterfaceName(name, device)
				for _, family := range families[lxd.NetworkName(device)] {
					if interfaceForwardAddress(state.Network[iface], family) == "" {
						return errors.NotAssignedf("IPv%d address for %s interface %q", family.version, networkType, iface)
					}
				}
			}
			return nil
		},
	})
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return state, errors.Trace(err)
}

func ovnInterfaceName(name string, device map[string]string) string {
	if iface := device["name"]; iface != "" {
		return iface
	}
	return name
}

func interfaceForwardAddress(state api.InstanceStateNetwork, family ovnForwardFamily) string {
	for _, address := range state.Addresses {
		ip := net.ParseIP(address.Address)
		if ip.IsGlobalUnicast() && family.subnet.Contains(ip) {
			return ip.String()
		}
	}
	return ""
}

func ensureOVNNetworkForward(ctx context.Context, srv Server, networkName, instanceName, iface, address string, forwards []api.NetworkForward) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	isIPv4 := net.ParseIP(address).To4() != nil
	var found bool
	for _, forward := range forwards {
		if forward.Config[jujuInstanceForwardKey] != instanceName ||
			forward.Config[jujuDeviceForwardKey] != iface {
			continue
		}
		// IPv4 and IPv6 forwards for the same interface are independent.
		listenIP := net.ParseIP(forward.ListenAddress)
		if listenIP == nil || (listenIP.To4() != nil) != isIPv4 {
			continue
		}
		if !found && forward.Config["target_address"] == address {
			found = true
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		op, err := srv.DeleteNetworkForward(networkName, forward.ListenAddress)
		if err == nil {
			err = op.WaitContext(ctx)
		}
		if err != nil && !lxd.IsLXDNotFound(errors.Cause(err)) {
			return errors.Annotatef(err, "deleting stale forward %q on network %q", forward.ListenAddress, networkName)
		}
	}
	if found {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	listenAddress := "::"
	if isIPv4 {
		listenAddress = "0.0.0.0"
	}
	err := createInstanceNetworkForward(ctx, srv, networkName, instanceName, iface, listenAddress, address)
	return errors.Annotatef(err, "creating forward for OVN interface %q on network %q", iface, networkName)
}

func createInstanceNetworkForward(ctx context.Context, srv Server, networkName, instanceName, iface, listenAddress, address string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	forward := api.NetworkForwardsPost{
		ListenAddress: listenAddress,
		NetworkForwardPut: api.NetworkForwardPut{
			Description: fmt.Sprintf("Juju instance %s interface %s", instanceName, iface),
			Config: map[string]string{
				"target_address":       address,
				jujuInstanceForwardKey: instanceName,
				jujuDeviceForwardKey:   iface,
			},
		},
	}
	op, err := srv.CreateNetworkForward(networkName, forward)
	if err == nil {
		err = op.WaitContext(ctx)
	}
	return errors.Trace(err)
}

// removeInstances cleans up owned forwards before deleting provider instances.
// A cleanup failure leaves the instances available for a subsequent retry.
func removeInstances(ctx context.Context, srv Server, names []string) error {
	if err := removeInstanceNetworkForwards(ctx, srv, names); err != nil {
		return errors.Annotate(err, "removing instance network forwards")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.Trace(srv.RemoveContainers(names))
}

func removeInstanceNetworkForwards(ctx context.Context, srv Server, names []string) error {
	if len(names) == 0 {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if !srv.HasExtension("network_forward") {
		return nil
	}
	// Scan by ownership, so cleanup also works after an instance or one of
	// its NICs has disappeared. Only forwards for the requested names qualify.
	owners := set.NewStrings(names...)
	networks, err := srv.GetNetworks()
	if err != nil {
		return errors.Annotate(err, "retrieving networks")
	}
	for _, network := range networks {
		if network.Type != networkTypeOVN && !(network.Type == "bridge" && network.Managed) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		forwards, err := srv.GetNetworkForwards(network.Name)
		if err != nil {
			if lxd.IsLXDNotFound(errors.Cause(err)) {
				continue
			}
			return errors.Annotatef(err, "retrieving forwards for network %q", network.Name)
		}
		for _, forward := range forwards {
			owner := forward.Config[jujuInstanceForwardKey]
			if !owners.Contains(owner) {
				continue
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := deleteInstanceNetworkForward(ctx, srv, network.Name, forward); err != nil {
				return errors.Annotatef(err, "deleting forward %q for instance %q on network %q", forward.ListenAddress, owner, network.Name)
			}
		}
	}
	return nil
}

func deleteInstanceNetworkForward(ctx context.Context, srv Server, networkName string, forward api.NetworkForward) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// OVN forwards have no location. Bridge forwards must be removed from
	// their owning member even if the instance has moved or disappeared.
	if forward.Location != "" && srv.IsClustered() {
		target, err := srv.UseTargetServer(ctx, forward.Location)
		if err != nil {
			return errors.Trace(err)
		}
		srv = &providerServer{target}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	op, err := srv.DeleteNetworkForward(networkName, forward.ListenAddress)
	if err == nil {
		err = op.WaitContext(ctx)
	}
	if lxd.IsLXDNotFound(errors.Cause(err)) {
		return nil
	}
	return errors.Trace(err)
}
