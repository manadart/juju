// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

// ResizeTarget is the resolved resource and attachment used by resize workers.
// Identity binds all mutable execution inputs; capacities are observations and
// are deliberately excluded from that identity.
type ResizeTarget struct {
	Identity       string
	SizeMiB        uint64
	ProviderType   string
	PoolName       string
	PoolAttributes map[string]string
	Kind           string
	MachineScoped  bool
	MachineUUID    string
	InstanceID     string
	AttachmentID   string
	VolumeID       string
	ProviderID     string
	DeviceName     string
	DeviceLink     string
	FilesystemID   string
	FilesystemUUID string
	MountPoint     string
	Alive          bool
	ReadOnly       bool
	Shared         bool
}

// ResizeStage is durable progress for either backend or attachment work.
type ResizeStage struct {
	State             string
	Message           string
	Prerequisites     []string
	FollowUp          []string
	ObservedSizeMiB   uint64
	FilesystemSizeMiB uint64
}

// ResizePlan holds the approved target and separate executor observations.
// Empty Identity means no plan has been committed. Empty BackendMachineUUID
// identifies model execution; empty NodeMachineUUID means no node stage.
type ResizePlan struct {
	UUID               string
	Identity           string
	BackendMachineUUID string
	NodeMachineUUID    string
	ApprovedSizeMiB    uint64
	Target             ResizeTarget
	Backend            ResizeStage
	Attachment         ResizeStage
}

// Next returns the eligible stage and phase for the given executor. The empty
// machine UUID identifies the model worker. State also uses this decision when
// atomically accepting a report, to guard concurrent advancement.
func (p ResizePlan) Next(machineUUID string) (stage, phase string) {
	if p.BackendMachineUUID == machineUUID {
		if p.Backend.State == "checking" {
			return "backend", "check"
		}
		if p.ApprovedSizeMiB != 0 && (p.Backend.State == "checked" || p.Backend.State == "running") {
			return "backend", "ensure"
		}
	}
	if machineUUID != "" && p.NodeMachineUUID == machineUUID {
		if p.Backend.State != "checking" && p.Attachment.State == "checking" {
			return "attachment", "check"
		}
		if p.ApprovedSizeMiB != 0 && p.Backend.State == "completed" && (p.Attachment.State == "checked" || p.Attachment.State == "running") {
			return "attachment", "ensure"
		}
	}
	return "", ""
}

// ResizeStageWork is one bounded reconciliation, not permission to skip checks.
type ResizeStageWork struct {
	UUID     string
	Identity string
	Stage    string
	Phase    string
	Target   ResizeTarget
}

// ResizeStageReport records a check or execution step. Status is blocked,
// checked, pending, completed or failed. Identity and phase fence stale reports.
type ResizeStageReport struct {
	UUID              string
	Identity          string
	Stage             string
	Phase             string
	Status            string
	Message           string
	Prerequisites     []string
	FollowUp          []string
	ObservedSizeMiB   uint64
	FilesystemSizeMiB uint64
}
