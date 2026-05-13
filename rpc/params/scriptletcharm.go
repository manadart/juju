// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

// ScriptletRelation carries one relation endpoint parsed from a scriptlet
// charm's metadata.yaml.
type ScriptletRelation struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Interface string `json:"interface"`
	Scope     string `json:"scope,omitempty"`
	Optional  bool   `json:"optional,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// RegisterScriptletCharmArgs contains raw scriptlet charm deployment data.
type RegisterScriptletCharmArgs struct {
	ApplicationName string              `json:"application-name"`
	Scriptlet       string              `json:"scriptlet"`
	Relations       []ScriptletRelation `json:"relations,omitempty"`
}
