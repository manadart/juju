// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import "context"

// VolumeResizer is an optional capability of a VolumeSource. Checks must have
// no side effects. Unsupported requests return an error matching
// core/errors.NotSupported; inspection failures must remain errors.
// EnsureVolumeSize must preserve identity and data, honour cancellation, and
// safely resume using OperationID after an ambiguous response or worker restart.
// Ensure must revalidate changed conditions and refuse new, unmet
// prerequisites before mutation. Acceptance is progress, not completion.
// Each call performs a bounded step.
type VolumeResizer interface {
	CheckVolumeResize(context.Context, VolumeResizeRequest) (ResizeRequirements, error)
	EnsureVolumeSize(context.Context, VolumeResizeRequest) (ResizeProgress, error)
}

// FilesystemResizer is the equivalent capability for a provider-native
// filesystem. The same check, retry and completion contracts apply.
type FilesystemResizer interface {
	CheckFilesystemResize(context.Context, FilesystemResizeRequest) (ResizeRequirements, error)
	EnsureFilesystemSize(context.Context, FilesystemResizeRequest) (ResizeProgress, error)
}

// ResizeAttachment identifies the attachment relevant to a provider operation.
type ResizeAttachment struct {
	AttachmentID string
	InstanceID   string
}

// VolumeResizeRequest describes growth of an existing provider volume.
// OperationID and target remain stable across retries. The provider must
// discover existing remote progress when a previous response was lost.
type VolumeResizeRequest struct {
	OperationID string
	VolumeID    string
	SizeMiB     uint64
	Attachments []ResizeAttachment
}

// FilesystemResizeRequest describes growth of a provider-native filesystem.
type FilesystemResizeRequest struct {
	OperationID  string
	FilesystemID string
	SizeMiB      uint64
	Attachments  []ResizeAttachment
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
