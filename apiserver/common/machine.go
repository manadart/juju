// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"context"
	"time"

	"github.com/juju/errors"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/machine"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/core/objectstore"
	"github.com/juju/juju/core/status"
	machineerrors "github.com/juju/juju/domain/machine/errors"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
)

// StateJobs translates a slice of multiwatcher jobs to their equivalents in state.
func StateJobs(jobs []model.MachineJob) ([]state.MachineJob, error) {
	newJobs := make([]state.MachineJob, len(jobs))
	for i, job := range jobs {
		newJob, err := machineJobFromParams(job)
		if err != nil {
			return nil, err
		}
		newJobs[i] = newJob
	}
	return newJobs, nil
}

// machineJobFromParams returns the job corresponding to model.MachineJob.
func machineJobFromParams(job model.MachineJob) (state.MachineJob, error) {
	switch job {
	case model.JobHostUnits:
		return state.JobHostUnits, nil
	case model.JobManageModel:
		return state.JobManageModel, nil
	default:
		return -1, errors.Errorf("invalid machine job %q", job)
	}
}

type origStateInterface interface {
	Machine(string) (*state.Machine, error)
}

type stateInterface interface {
	Machine(string) (Machine, error)
}

type stateShim struct {
	origStateInterface
}

func (st *stateShim) Machine(id string) (Machine, error) {
	return st.origStateInterface.Machine(id)
}

type ControllerNode interface {
	Id() string
	HasVote() bool
	WantsVote() bool
}

type Machine interface {
	Id() string
	Status() (status.StatusInfo, error)
	ContainerType() instance.ContainerType
	Life() state.Life
	ForceDestroy(time.Duration) error
	Destroy(objectstore.ObjectStore) error
	IsManager() bool
}

func DestroyMachines(st origStateInterface, store objectstore.ObjectStore, force bool, maxWait time.Duration, ids ...string) error {
	return destroyMachines(&stateShim{st}, store, force, maxWait, ids...)
}

func destroyMachines(st stateInterface, store objectstore.ObjectStore, force bool, maxWait time.Duration, ids ...string) error {
	var errs []error
	for _, id := range ids {
		machine, err := st.Machine(id)
		switch {
		case errors.Is(err, errors.NotFound):
			err = errors.Errorf("machine %s does not exist", id)
		case err != nil:
		case force:
			err = machine.ForceDestroy(maxWait)
		case machine.Life() != state.Alive:
			continue
		default:
			err = machine.Destroy(store)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return apiservererrors.DestroyErr("machines", ids, errs)
}

// ModelMachineInfo returns information about machine hardware for
// alive top level machines (not containers).
func ModelMachineInfo(ctx context.Context, st ModelManagerBackend, machineService MachineService) (machineInfo []params.ModelMachineInfo, _ error) {
	machines, err := st.AllMachines()
	if err != nil {
		return nil, errors.Trace(err)
	}
	controllerNodes, err := st.ControllerNodes()
	if err != nil {
		return nil, errors.Trace(err)
	}
	hasVote := make(map[string]bool)
	wantsVote := make(map[string]bool)
	for _, n := range controllerNodes {
		hasVote[n.Id()] = n.HasVote()
		wantsVote[n.Id()] = n.WantsVote()
	}
	var primaryID string
	primaryHA, err := st.HAPrimaryMachine()
	if err != nil {
		// We do not want to return any errors here as they are all
		// non-fatal for this call since we can still
		// get machine info even if we could not get HA Primary determined.
		// Also on some non-HA setups, i.e. where mongo was not run with --replSet,
		// this call will return an error.
		logger.Warningf(context.TODO(), "could not determine if there is a primary HA machine: %v", err)
	}
	if len(controllerNodes) > 1 {
		primaryID = primaryHA.Id()
	}

	for _, m := range machines {
		if m.Life() != state.Alive {
			continue
		}
		var aStatus string
		// This is suboptimal as if there are many machines,
		// we are making many calls into the DB for each machine.
		var statusMessage string
		statusInfo, err := m.Status()
		if err == nil {
			aStatus = string(statusInfo.Status)
			statusMessage = statusInfo.Message
		} else {
			aStatus = err.Error()
		}
		mInfo := params.ModelMachineInfo{
			Id:        m.Id(),
			HasVote:   hasVote[m.Id()],
			WantsVote: wantsVote[m.Id()],
			Status:    aStatus,
			Message:   statusMessage,
		}
		if primaryID != "" {
			if isPrimary := primaryID == m.Id(); isPrimary {
				mInfo.HAPrimary = &isPrimary
			}
		}
		machineUUID, err := machineService.GetMachineUUID(ctx, machine.Name(m.Id()))
		if errors.Is(err, machineerrors.MachineNotFound) {
			return nil, errors.NotFoundf("machine %q", m.Id())
		} else if err != nil {
			return nil, errors.Trace(err)
		}
		instanceID, displayName, err := machineService.InstanceIDAndName(ctx, machineUUID)
		switch {
		case err == nil:
			mInfo.InstanceId = instanceID.String()
			mInfo.DisplayName = displayName
		case errors.Is(err, machineerrors.MachineNotFound):
			return nil, errors.NotFoundf("machine %q", m.Id())
		case errors.Is(err, machineerrors.NotProvisioned):
			// ok, but no instance ID to get.
		default:
			return nil, errors.Trace(err)
		}
		if m.ContainerType() != "" && m.ContainerType() != instance.NONE {
			machineInfo = append(machineInfo, mInfo)
			continue
		}
		// Only include cores for physical machines.
		hw, err := machineService.HardwareCharacteristics(ctx, machineUUID)
		if errors.Is(err, machineerrors.MachineNotFound) {
			return nil, errors.NotFoundf("machine %q", m.Id())
		} else if err != nil && !errors.Is(err, machineerrors.NotProvisioned) {
			return nil, errors.Trace(err)
		}
		if hw != nil && hw.String() != "" {
			hwParams := &params.MachineHardware{
				Cores:            hw.CpuCores,
				Arch:             hw.Arch,
				Mem:              hw.Mem,
				RootDisk:         hw.RootDisk,
				CpuPower:         hw.CpuPower,
				Tags:             hw.Tags,
				AvailabilityZone: hw.AvailabilityZone,
				VirtType:         hw.VirtType,
			}
			mInfo.Hardware = hwParams
		}
		machineInfo = append(machineInfo, mInfo)
	}
	return machineInfo, nil
}
