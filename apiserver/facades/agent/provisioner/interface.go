// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provisioner

import (
	"github.com/juju/names/v6"

	corenetwork "github.com/juju/juju/core/network"
	"github.com/juju/juju/internal/network/containerizer"
)

// Machine is an indirection for use in container provisioning.
// It is an indirection for both containerizer.Machine and
// containerizer.Container as well as state.Machine locally.
type Machine interface {
	containerizer.Container

	MachineTag() names.MachineTag
}

// BridgePolicy is an indirection for containerizer.BridgePolicy.
type BridgePolicy interface {
	// PopulateContainerLinkLayerDevices sets the link-layer devices of the input
	// guest, setting each device to be a child of the corresponding bridge on the
	// host machine.
	PopulateContainerLinkLayerDevices(
		containerizer.Machine, containerizer.Container, bool,
	) (corenetwork.InterfaceInfos, error)
}
