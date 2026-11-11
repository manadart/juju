// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

//go:build !linux

package resize

import (
	"context"

	"github.com/juju/juju/internal/errors"
)

// Resize is only implemented for Linux loop devices.
func (*Executor) Resize(context.Context, Request, func(context.Context) error) (Observation, error) {
	return Observation{}, errors.New("local storage resize requires Linux")
}

func runCommand(context.Context, string, ...string) (string, error) {
	return "", errors.New("local storage resize requires Linux")
}

func (*Executor) CheckBackend(context.Context, Request) (Observation, error) {
	return Observation{}, errors.New("local storage resize requires Linux")
}
func (*Executor) EnsureBackend(context.Context, Request, func(context.Context) error) (Observation, error) {
	return Observation{}, errors.New("local storage resize requires Linux")
}
func (*Executor) CheckAttachment(context.Context, Request) (Observation, error) {
	return Observation{}, errors.New("local storage resize requires Linux")
}
func (*Executor) EnsureAttachment(context.Context, Request, func(context.Context) error) (Observation, error) {
	return Observation{}, errors.New("local storage resize requires Linux")
}
