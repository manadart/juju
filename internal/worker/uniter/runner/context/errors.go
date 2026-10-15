// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package context

import (
	"github.com/juju/errors"
)

var ErrRequeueAndReboot = errors.New("reboot now")
var ErrReboot = errors.New("reboot after hook")

// ErrHookDeferred indicates that the hook committed successfully and must run
// again before its lifecycle can advance.
var ErrHookDeferred = errors.New("hook deferred after commit")
