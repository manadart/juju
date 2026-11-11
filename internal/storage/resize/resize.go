// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package resize

import "context"

// Request identifies the existing volume and optional managed filesystem.
type Request struct {
	ProviderType     string
	VolumeID         string
	VolumeProviderID string
	DeviceName       string
	FilesystemID     string
	MountPoint       string
	SizeMiB          uint64
}

// Observation contains capacities measured from the device and filesystem.
type Observation struct {
	VolumeSizeMiB     uint64
	FilesystemSizeMiB uint64
}

// RunCommandFunc executes a command with cancellation and returns its output.
type RunCommandFunc func(context.Context, string, ...string) (string, error)

// Executor performs machine-local resize stages on existing resources.
type Executor struct {
	StorageDir string
	Run        RunCommandFunc
}

// NewExecutor creates an executor for the machine agent's storage directory.
func NewExecutor(storageDir string) *Executor {
	return &Executor{StorageDir: storageDir, Run: runCommand}
}
