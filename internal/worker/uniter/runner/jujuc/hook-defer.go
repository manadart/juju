// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package jujuc

import (
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/cmd"
)

// HookDeferCommand requests another invocation of the current hook.
type HookDeferCommand struct {
	cmd.CommandBase
	ctx Context
}

// NewHookDeferCommand creates a hook-defer command.
func NewHookDeferCommand(ctx Context) (cmd.Command, error) {
	return &HookDeferCommand{ctx: ctx}, nil
}

// Info implements cmd.Command.
func (c *HookDeferCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "hook-defer",
		Purpose: "Commits a successful hook's changes and schedules it to run again.",
		Doc: `Request another invocation of the current hook.

The command returns normally. The hook must still exit successfully before
Juju commits its changes and runs the same hook again with a fresh context.
No other hook runs on the unit between these invocations. The logical hook
remains pending until an invocation succeeds without requesting deferral.

Repeated calls request only one additional invocation. Reissued invocations
receive JUJU_HOOK_REISSUED=1. A failed hook or commit follows normal error
handling instead of successful deferral.

This command is not supported in actions or juju exec, and cannot be combined
with a reboot request in the same invocation.`,
	})
}

// Init implements cmd.Command.
func (c *HookDeferCommand) Init(args []string) error {
	return cmd.CheckEmpty(args)
}

// Run implements cmd.Command.
func (c *HookDeferCommand) Run(*cmd.Context) error {
	return c.ctx.RequestHookDeferral()
}
