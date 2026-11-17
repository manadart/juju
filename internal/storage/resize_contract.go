// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import "context"

// Resizer is an optional capability of a VolumeSource or FilesystemSource.
// The selected source determines the kind of resource being resized.
// Checks must have no side effects. Unsupported requests return an error
// matching core/errors.NotSupported; inspection failures must remain errors.
// EnsureSize must preserve identity and data, honour cancellation, and
// safely resume using OperationID after an ambiguous response or worker restart.
// EnsureSize must revalidate changed conditions and refuse new, unmet
// prerequisites before mutation. Acceptance is progress, not completion.
// Each call performs a bounded step.
type Resizer interface {
	CheckResize(context.Context, ResizeRequest) (ResizeRequirements, error)
	EnsureSize(context.Context, ResizeRequest) (ResizeProgress, error)
}

// ResizeAttachment identifies the attachment relevant to a provider operation.
type ResizeAttachment struct {
	AttachmentID string
	InstanceID   string
}

// ResizeRequest describes growth of an existing provider resource.
// ProviderID identifies the resource within the selected volume or filesystem
// source. OperationID and target remain stable across retries. The provider
// must discover existing remote progress when a previous response was lost.
type ResizeRequest struct {
	OperationID string
	ProviderID  string
	SizeMiB     uint64
	Attachments []ResizeAttachment
}

// ResizeRequirement identifies a prerequisite or required follow-up step.
type ResizeRequirement string

const (
	ResizeUnmount        ResizeRequirement = "unmount"
	ResizeDetach         ResizeRequirement = "detach"
	ResizeDeviceRefresh  ResizeRequirement = "device-refresh"
	ResizePartitionGrow  ResizeRequirement = "partition-grow"
	ResizeFilesystemGrow ResizeRequirement = "filesystem-grow"
)

// ResizeRequirements describes eligibility at the time of inspection. Blocked
// is temporary. AlreadySatisfied applies only to this resource layer, and does
// not waive other stages. Unknown or unhandled requirements prevent approval.
type ResizeRequirements struct {
	Blocked          bool
	AlreadySatisfied bool
	Prerequisites    []ResizeRequirement
	FollowUp         []ResizeRequirement
	ObservedSizeMiB  uint64
	Message          string
}

// ResizeProgress contains provider observations from one reconciliation.
// Complete requires observed capacity at least the requested target. A false
// value includes accepted operations and operations still converging.
type ResizeProgress struct {
	Complete        bool
	ObservedSizeMiB uint64
	Message         string
}
