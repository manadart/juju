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

// ScriptletConfigOption describes a single config option for a scriptlet charm.
type ScriptletConfigOption struct {
	Key          string `json:"key"`
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	DefaultValue string `json:"default-value,omitempty"`
}

// DeployScriptletCharmArgs contains all data needed to register a scriptlet
// charm and create the corresponding application in one shot.
type DeployScriptletCharmArgs struct {
	ApplicationName string                  `json:"application-name"`
	Scriptlet       string                  `json:"scriptlet"`
	Relations       []ScriptletRelation     `json:"relations,omitempty"`
	Config          []ScriptletConfigOption `json:"config,omitempty"`
	Runtime         string                  `json:"runtime,omitempty"`
	App             string                  `json:"app,omitempty"`
	Events          []string                `json:"events,omitempty"`
}
