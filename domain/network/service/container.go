// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"github.com/juju/utils/set"

	"github.com/juju/juju/core/machine"
	"github.com/juju/juju/core/trace"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/domain/network"
)

type ContainerState interface {
	// GetMachinePositiveSpaceConstraints retrieves the positive
	// space constraints for the machine with the input UUID.
	GetMachinePositiveSpaceConstraints(ctx context.Context, machineUUID machine.UUID) ([]string, error)
	// GetMachineAppBindings retrieves the bound spaces for applications
	// with units assigned to the machine with the input UUID.
	GetMachineAppBindings(ctx context.Context, machineUUID machine.UUID) ([]string, error)

	// NICsInSpaces retrieves the link-layer devices on the machine with the
	// input net node UUID that are connected the input spaces.
	NICsInSpaces(ctx context.Context, netNode string, spaces []string) ([]network.NetInterface, error)
}

// DevicesToBridge accepts the UUID of a host machine and a guest container/VM.
// It returns the information needed for creating network bridges that will be
// parents of the guest's virtual network devices.
// This determination is made based on the guest's space constraints, bindings
// of applications to run on the guest, and any host bridges that already exist.
func (s *Service) DevicesToBridge(
	ctx context.Context, hostUUID, guestUUID machine.UUID,
) ([]network.DeviceToBridge, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := hostUUID.Validate(); err != nil {
		return nil, errors.Errorf("invalid machine UUID: %w", err)
	}
	if err := guestUUID.Validate(); err != nil {
		return nil, errors.Errorf("invalid machine UUID: %w", err)
	}

	spaces, err := s.spaceReqsForMachine(ctx, guestUUID)
	if err != nil {
		return nil, errors.Capture(err)
	}
	s.logger.Debugf(ctx, "machine %q needs spaces %v spaces", guestUUID, spaces)

	nics, err := s.NICsInSpaces(ctx, hostUUID, spaces)
	if err != nil {
		return nil, errors.Capture(err)
	}

	return nil, nil
}

// spacesForMachine returns UUIDs for the *positive* space
// requirements of the machine with the input UUID.
func (s *Service) spaceReqsForMachine(ctx context.Context, machineUUID machine.UUID) ([]string, error) {
	spaces, err := s.st.GetMachinePositiveSpaceConstraints(ctx, machineUUID)
	if err != nil {
		return nil, errors.Errorf("retrieving positive space constraints for machine %q: %w", machineUUID, err)
	}

	bound, err := s.st.GetMachineAppBindings(ctx, machineUUID)
	if err != nil {
		return nil, errors.Errorf("retrieving app bindings for machine %q: %w", machineUUID, err)
	}

	return set.NewStrings(append(spaces, bound...)...).Values(), nil
}

// TODO: This needs to map space names to NICs so we can ensure satisfaction.
func (s *Service) NICsInSpaces(
	ctx context.Context, mUUID machine.UUID, spaces []string,
) ([]network.NetInterface, error) {
	nodeUUID, err := s.st.GetMachineNetNodeUUID(ctx, mUUID.String())
	if err != nil {
		return nil, errors.Errorf("retrieving net node for machine %q: %w", mUUID., err)
	}

	nics, err := s.st.NICsInSpaces(ctx, nodeUUID, spaces)
	if err != nil {
		return nil, errors.Errorf("retrieving NICs for machine %q in spaces %v: %w", mUUID, spaces, err)
	}

	return nics, nil
}