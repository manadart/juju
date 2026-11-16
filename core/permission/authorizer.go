// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package permission

import "context"

// Authorizer evaluates a subject's access using the persisted permission model.
// It is shared by access services. Its worker owns startup and shutdown;
// consumers must not close it. Checks fail when the worker is shutting down.
type Authorizer interface {
	Check(context.Context, string, ID, Access) (bool, error)
}
