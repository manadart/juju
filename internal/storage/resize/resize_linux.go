// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package resize

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/juju/names/v6"

	"github.com/juju/juju/internal/errors"
)

const bytesPerMiB = 1024 * 1024

// Resize grows the backing file, refreshes the attached device, then grows a
// managed ext4 filesystem. Revalidate runs before each mutation. Cancellation
// leaves the request retryable; the same stages are safe to repeat.
func (e *Executor) Resize(ctx context.Context, request Request, revalidate func(context.Context) error) (Observation, error) {
	if _, err := e.CheckBackend(ctx, request); err != nil {
		return Observation{}, err
	}
	if _, err := e.CheckAttachment(ctx, request); err != nil {
		return Observation{}, err
	}
	if _, err := e.EnsureBackend(ctx, request, revalidate); err != nil {
		return Observation{}, err
	}
	observed, err := e.EnsureAttachment(ctx, request, revalidate)
	if err == nil && (observed.VolumeSizeMiB < request.SizeMiB || (request.FilesystemID != "" && observed.FilesystemSizeMiB < request.SizeMiB)) {
		err = errors.New("local capacity has not reached the requested target")
	}
	return observed, err
}

// CheckBackend inspects the existing loop backing file without modifying it.
func (e *Executor) CheckBackend(ctx context.Context, request Request) (Observation, error) {
	backing, device, info, err := e.loopTarget(ctx, request)
	if err != nil {
		return Observation{}, err
	}
	if err := e.checkLoop(ctx, backing, device, info); err != nil {
		return Observation{}, err
	}
	return Observation{VolumeSizeMiB: uint64(info.Size()) / bytesPerMiB}, nil
}

// EnsureBackend grows only the local backing file. Device and filesystem
// observations belong to the attachment stage, even for local storage.
func (e *Executor) EnsureBackend(ctx context.Context, request Request, revalidate func(context.Context) error) (Observation, error) {
	if revalidate == nil {
		return Observation{}, errors.New("missing resize revalidation")
	}
	if err := revalidate(ctx); err != nil {
		return Observation{}, err
	}
	backing, device, info, err := e.loopTarget(ctx, request)
	if err != nil {
		return Observation{}, err
	}
	if err := e.checkLoop(ctx, backing, device, info); err != nil {
		return Observation{}, err
	}
	err = e.ensureBackingSize(ctx, backing, info, request.SizeMiB)
	current, statErr := os.Stat(backing)
	if statErr != nil {
		return Observation{}, statErr
	}
	if !os.SameFile(info, current) {
		return Observation{}, errors.New("loop backing file was replaced")
	}
	return Observation{VolumeSizeMiB: uint64(current.Size()) / bytesPerMiB}, err
}

// CheckAttachment inspects a whole device and optional ext4 mount. Cloud
// devices currently require a stable by-id link and a SCSI rescan interface.
func (e *Executor) CheckAttachment(ctx context.Context, request Request) (Observation, error) {
	var observed Observation
	device, err := e.attachmentDevice(ctx, request)
	if err != nil {
		return observed, err
	}
	if request.ProviderType != "loop" {
		readOnly, err := e.Run(ctx, "blockdev", "--getro", device)
		if err != nil {
			return observed, err
		}
		if strings.TrimSpace(readOnly) != "0" {
			return observed, errors.New("block device is read-only")
		}
	}
	if err := e.checkLayout(ctx, device, request.FilesystemID != ""); err != nil {
		return observed, err
	}
	size, err := e.deviceSize(ctx, device)
	if err != nil {
		return observed, err
	}
	observed.VolumeSizeMiB = size / bytesPerMiB
	if request.FilesystemID != "" {
		for _, tool := range []string{"findmnt", "dumpe2fs", "resize2fs"} {
			if _, err := commandPath(tool); err != nil {
				return observed, err
			}
		}
		if err := e.checkMount(ctx, device, request.MountPoint); err != nil {
			return observed, err
		}
		size, uuid, err := e.filesystemSize(ctx, device)
		if err != nil {
			return observed, err
		}
		if request.FilesystemUUID != "" && uuid != request.FilesystemUUID {
			return observed, errors.New("filesystem identity no longer matches the checked target")
		}
		observed.FilesystemSizeMiB = size / bytesPerMiB
	}
	return observed, nil
}

// EnsureAttachment refreshes the device and grows its managed filesystem.
// Insufficient device capacity is a pending observation so propagation can be
// retried. No partition creation, filesystem formatting or device replacement
// is permitted on this path.
func (e *Executor) EnsureAttachment(ctx context.Context, request Request, revalidate func(context.Context) error) (Observation, error) {
	var observed Observation
	if revalidate == nil {
		return observed, errors.New("missing resize revalidation")
	}
	if err := revalidate(ctx); err != nil {
		return observed, err
	}
	observed, err := e.CheckAttachment(ctx, request)
	if err != nil {
		return observed, err
	}
	// An already satisfied attachment needs observation only.
	if observed.VolumeSizeMiB >= request.SizeMiB && (request.FilesystemID == "" || observed.FilesystemSizeMiB >= request.SizeMiB) {
		return observed, nil
	}
	device, err := e.attachmentDevice(ctx, request)
	if err != nil {
		return observed, err
	}
	if request.ProviderType == "loop" {
		if _, err := e.Run(ctx, "losetup", "--set-capacity", device); err != nil {
			return observed, err
		}
	} else {
		if err := ctx.Err(); err != nil {
			return observed, err
		}
		// The kernel exposes this only for transports supporting this rescan.
		if err := os.WriteFile(filepath.Join("/sys/class/block", filepath.Base(device), "device/rescan"), []byte("1\n"), 0200); err != nil {
			return observed, errors.Errorf("rescanning block device: %w", err)
		}
	}
	size, err := e.deviceSize(ctx, device)
	if err != nil {
		return observed, err
	}
	observed.VolumeSizeMiB = size / bytesPerMiB
	if observed.VolumeSizeMiB < request.SizeMiB || request.FilesystemID == "" {
		return observed, nil
	}
	if err := revalidate(ctx); err != nil {
		return observed, err
	}
	current, err := e.CheckAttachment(ctx, request)
	if err != nil {
		return observed, err
	}
	observed.FilesystemSizeMiB = current.FilesystemSizeMiB
	if observed.FilesystemSizeMiB < request.SizeMiB {
		_, growErr := e.Run(ctx, "resize2fs", device, fmt.Sprintf("%dM", request.SizeMiB))
		size, uuid, readErr := e.filesystemSize(ctx, device)
		if readErr == nil && (request.FilesystemUUID == "" || uuid == request.FilesystemUUID) {
			observed.FilesystemSizeMiB = size / bytesPerMiB
		}
		if growErr != nil {
			return observed, errors.Errorf("growing ext4 filesystem: %w", growErr)
		}
		if readErr != nil {
			return observed, readErr
		}
		if request.FilesystemUUID != "" && uuid != request.FilesystemUUID {
			return observed, errors.New("filesystem identity changed during growth")
		}
	}
	return observed, nil
}

func (e *Executor) loopTarget(ctx context.Context, request Request) (string, string, os.FileInfo, error) {
	if e.Run == nil || e.StorageDir == "" {
		return "", "", nil, errors.New("incomplete resize executor configuration")
	}
	if request.ProviderType != "loop" || !names.IsValidVolume(request.VolumeID) {
		return "", "", nil, errors.New("only existing loop backends can be grown locally")
	}
	if request.SizeMiB == 0 || request.SizeMiB > math.MaxInt64/bytesPerMiB {
		return "", "", nil, errors.New("invalid resize target")
	}
	tag := names.NewVolumeTag(request.VolumeID).String()
	if tag != request.VolumeProviderID {
		return "", "", nil, errors.New("loop provider identity does not match")
	}
	name := strings.TrimPrefix(request.DeviceName, "/dev/")
	if !strings.HasPrefix(name, "loop") {
		return "", "", nil, errors.New("expected a loop device")
	}
	if _, err := strconv.ParseUint(strings.TrimPrefix(name, "loop"), 10, 32); err != nil {
		return "", "", nil, errors.New("invalid loop device name")
	}
	path := filepath.Join(e.StorageDir, "loop", tag)
	info, err := os.Lstat(path)
	if err != nil {
		return "", "", nil, err
	}
	if !info.Mode().IsRegular() {
		return "", "", nil, errors.New("loop backing path is not a regular file")
	}
	for _, tool := range []string{"losetup", "lsblk", "blockdev", "fallocate"} {
		if _, err := commandPath(tool); err != nil {
			return "", "", nil, err
		}
	}
	return path, filepath.Join("/dev", name), info, nil
}

func (e *Executor) attachmentDevice(ctx context.Context, request Request) (string, error) {
	if e.Run == nil || request.SizeMiB == 0 || request.SizeMiB > math.MaxInt64/bytesPerMiB {
		return "", errors.New("invalid attachment resize configuration")
	}
	if request.ProviderType == "loop" {
		backing, device, info, err := e.loopTarget(ctx, request)
		if err != nil {
			return "", err
		}
		return device, e.checkLoop(ctx, backing, device, info)
	}
	name := strings.TrimPrefix(request.DeviceName, "/dev/")
	if name == "" || name == "." || name == ".." || filepath.Base(name) != name {
		return "", errors.New("invalid block device name")
	}
	device := filepath.Join("/dev", name)
	if !strings.HasPrefix(request.DeviceLink, "/dev/disk/by-id/") {
		return "", errors.New("cloud attachment resize requires a stable device link")
	}
	actual, err := os.Stat(device)
	if err != nil {
		return "", err
	}
	link, err := os.Stat(request.DeviceLink)
	if err != nil {
		return "", err
	}
	if !os.SameFile(actual, link) || actual.Mode()&os.ModeDevice == 0 || actual.Mode()&os.ModeCharDevice != 0 {
		return "", errors.New("block device identity does not match its stable link")
	}
	if _, err := os.Stat(filepath.Join("/sys/class/block", name, "device/rescan")); err != nil {
		return "", errors.New("attachment transport does not support the implemented SCSI rescan")
	}
	for _, tool := range []string{"lsblk", "blockdev"} {
		if _, err := commandPath(tool); err != nil {
			return "", err
		}
	}
	return device, nil
}

func (e *Executor) checkLoop(ctx context.Context, backing, device string, original os.FileInfo) error {
	current, err := os.Lstat(backing)
	if err != nil {
		return err
	}
	if !current.Mode().IsRegular() || !os.SameFile(original, current) {
		return errors.New("loop backing file was replaced")
	}
	out, err := e.Run(ctx, "losetup", "--list", "--noheadings", "--raw", "--output", "NAME,BACK-INO,OFFSET,SIZELIMIT,RO", "--associated", backing)
	if err != nil {
		return errors.Errorf("checking loop association: %w", err)
	}
	fields := strings.Fields(out)
	if len(fields) != 5 || fields[0] != device {
		return errors.New("expected exactly one matching loop attachment")
	}
	inode, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return err
	}
	if inode != original.Sys().(*syscall.Stat_t).Ino || fields[2] != "0" || fields[3] != "0" || fields[4] != "0" {
		return errors.New("loop attachment is read-only, offset, limited, or has changed identity")
	}
	return nil
}

func (e *Executor) checkLayout(ctx context.Context, device string, managed bool) error {
	out, err := e.Run(ctx, "lsblk", "--json", "--output", "NAME,TYPE,FSTYPE", device)
	if err != nil {
		return errors.Errorf("checking device layout: %w", err)
	}
	type node struct {
		Name       string            `json:"name"`
		Type       string            `json:"type"`
		Filesystem string            `json:"fstype"`
		Children   []json.RawMessage `json:"children"`
	}
	var layout struct {
		Devices []node `json:"blockdevices"`
	}
	if err := json.Unmarshal([]byte(out), &layout); err != nil {
		return err
	}
	if len(layout.Devices) != 1 || layout.Devices[0].Name != filepath.Base(device) || (layout.Devices[0].Type != "loop" && layout.Devices[0].Type != "disk") || len(layout.Devices[0].Children) != 0 {
		return errors.New("partitioned or stacked devices are not supported")
	}
	if managed && layout.Devices[0].Filesystem != "ext4" {
		return errors.New("only whole-device ext4 managed filesystems are supported")
	}
	return nil
}

func (e *Executor) checkMount(ctx context.Context, device, mountPoint string) error {
	if !filepath.IsAbs(mountPoint) {
		return errors.New("invalid filesystem mount point")
	}
	out, err := e.Run(ctx, "findmnt", "--json", "--mountpoint", mountPoint, "--output", "SOURCE,FSTYPE,OPTIONS,FSROOT")
	if err != nil {
		return errors.Errorf("checking filesystem mount: %w", err)
	}
	var mounts struct {
		Filesystems []struct {
			Source  string `json:"source"`
			Type    string `json:"fstype"`
			Options string `json:"options"`
			Root    string `json:"fsroot"`
		} `json:"filesystems"`
	}
	if err := json.Unmarshal([]byte(out), &mounts); err != nil {
		return err
	}
	if len(mounts.Filesystems) != 1 {
		return errors.New("expected exactly one mounted filesystem")
	}
	mount := mounts.Filesystems[0]
	if mount.Type != "ext4" || mount.Root != "/" || !strings.Contains(","+mount.Options+",", ",rw,") {
		return errors.New("filesystem is not a writable whole ext4 mount")
	}
	source, err := os.Stat(mount.Source)
	if err != nil {
		return err
	}
	target, err := os.Stat(device)
	if err != nil {
		return err
	}
	if source.Mode()&os.ModeDevice == 0 || source.Mode()&os.ModeCharDevice != 0 || source.Sys().(*syscall.Stat_t).Rdev != target.Sys().(*syscall.Stat_t).Rdev {
		return errors.New("mounted filesystem is not on the expected block device")
	}
	return nil
}

func (e *Executor) ensureBackingSize(ctx context.Context, path string, original os.FileInfo, sizeMiB uint64) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !os.SameFile(original, info) {
		return errors.New("loop backing file was replaced")
	}
	if info.Size() >= int64(sizeMiB*bytesPerMiB) {
		return nil
	}
	if _, err := e.Run(ctx, "fallocate", "--length", strconv.FormatUint(sizeMiB*bytesPerMiB, 10), path); err != nil {
		return errors.Errorf("growing loop backing file: %w", err)
	}
	return f.Sync()
}

func (e *Executor) deviceSize(ctx context.Context, device string) (uint64, error) {
	out, err := e.Run(ctx, "blockdev", "--getsize64", device)
	if err != nil {
		return 0, errors.Errorf("reading device capacity: %w", err)
	}
	return strconv.ParseUint(strings.TrimSpace(out), 10, 64)
}

func (e *Executor) filesystemSize(ctx context.Context, device string) (uint64, string, error) {
	out, err := e.Run(ctx, "dumpe2fs", "-h", device)
	if err != nil {
		return 0, "", errors.Errorf("reading ext4 geometry: %w", err)
	}
	var blocks, blockSize uint64
	var uuid string
	for _, line := range strings.Split(out, "\n") {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Block count":
			blocks, err = strconv.ParseUint(strings.TrimSpace(value), 10, 64)
		case "Block size":
			blockSize, err = strconv.ParseUint(strings.TrimSpace(value), 10, 64)
		case "Filesystem UUID":
			uuid = strings.TrimSpace(value)
		}
		if err != nil {
			return 0, "", err
		}
	}
	if blocks == 0 || blockSize == 0 || uuid == "" || blocks > math.MaxUint64/blockSize {
		return 0, "", errors.New("invalid ext4 geometry")
	}
	return blocks * blockSize, uuid, nil
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	path, err := commandPath(name)
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), errors.Errorf("%s: %s: %w", name, strings.TrimSpace(string(out)), err)
	}
	return string(out), nil
}

func commandPath(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	for _, dir := range []string{"/usr/sbin", "/sbin"} {
		if path, err := exec.LookPath(filepath.Join(dir, name)); err == nil {
			return path, nil
		}
	}
	return "", errors.Errorf("required resize command %q is unavailable", name)
}
