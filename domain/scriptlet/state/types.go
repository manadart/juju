// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

type scriptletCharm struct {
	ApplicationName string `db:"application_name"`
	Scriptlet       string `db:"scriptlet"`
}
