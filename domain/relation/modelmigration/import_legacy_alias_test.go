// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmigration

import (
	"github.com/juju/description/v12"
	"github.com/juju/tc"

	corerelation "github.com/juju/juju/core/relation"
	domainmodelmigration "github.com/juju/juju/domain/modelmigration/modelmigration"
)

// TestImportLegacyAliasPreservesRelationIdentityAndUnitSettings checks an
// established relation to the second alias of an offer. Deduplication must not
// replace its remote token, or leave settings addressed to a unit name that no
// longer matches the imported endpoint. Check the returned name so preserving
// the alias or consistently remapping it can both satisfy the contract.
func (s *importSuite) TestImportLegacyAliasPreservesRelationIdentityAndUnitSettings(c *tc.C) {
	m := description.NewModel(description.ModelArgs{})
	primary := m.AddRemoteApplication(description.RemoteApplicationArgs{Name: "first"})
	duplicate := m.AddRemoteApplication(description.RemoteApplicationArgs{Name: "second"})
	r := m.AddRelation(description.RelationArgs{Id: 42, Key: "client:db second:db"})
	r.AddEndpoint(description.EndpointArgs{ApplicationName: "client", Name: "db", Role: "requirer"})
	ep := r.AddEndpoint(description.EndpointArgs{ApplicationName: "second", Name: "db", Role: "provider"})
	ep.SetUnitSettings("second/0", map[string]any{"token": "keep-me"})
	key, err := corerelation.NewKeyFromString(r.Key())
	c.Assert(err, tc.ErrorIsNil)
	const token = "6049aa01-76c9-462d-8440-964a6e26aac2"
	arg, err := (&importOperation{}).createRemoteImportArg(r,
		domainmodelmigration.RemoteApplicationOfferer{Primary: primary, Duplicates: []description.RemoteApplication{duplicate}},
		[]relationRemoteEntity{{RelationKey: key, RelationUUID: token}},
	)
	c.Assert(err, tc.ErrorIsNil)
	c.Check(arg.ID, tc.Equals, 42)
	c.Check(arg.UUID.String(), tc.Equals, token)
	c.Assert(arg.Endpoints, tc.HasLen, 2)
	c.Check(arg.Endpoints[1].UnitSettings[arg.Endpoints[1].ApplicationName+"/0"], tc.DeepEquals, map[string]any{"token": "keep-me"})
}
