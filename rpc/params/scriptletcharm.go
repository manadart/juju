// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

// RegisterScriptletCharmArgs contains raw scriptlet charm deployment data.
type RegisterScriptletCharmArgs struct {
	ApplicationName string `json:"application-name"`
	Scriptlet       string `json:"scriptlet"`
}
