// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import "github.com/juju/errors"

// DirectiveUpdate changes the size of an allocation default, retaining its pool
// and count. SizeMiB must be non-nil and positive.
type DirectiveUpdate struct {
	SizeMiB *uint64
}

// ParseDirectiveUpdate accepts only a size, such as 20GiB.
func ParseDirectiveUpdate(s string) (DirectiveUpdate, error) {
	size, ok, err := parseSize(s)
	if err != nil {
		return DirectiveUpdate{}, err
	}
	if !ok || size == 0 {
		return DirectiveUpdate{}, errors.NotValidf("storage size %q (expected a positive size such as 20GiB)", s)
	}
	return DirectiveUpdate{SizeMiB: &size}, nil
}
