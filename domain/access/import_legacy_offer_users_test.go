// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package access_test

import (
	"github.com/juju/clock"
	"github.com/juju/description/v12"
	"github.com/juju/tc"

	"github.com/juju/juju/core/user"
	"github.com/juju/juju/domain/access/modelmigration"
	loggertesting "github.com/juju/juju/internal/logger/testing"
)

// TestImportLegacyOfferOnlyExternalUsers checks grantees absent from model
// membership and unknown to the target controller. Offer ACLs must import with
// their original permissions, without accidentally granting model access.
func (s *importSuite) TestImportLegacyOfferOnlyExternalUsers(c *tc.C) {
	s.seedEveryoneExternalUser(c)
	modelmigration.RegisterExternalUsersImport(s.coordinator, clock.WallClock, loggertesting.WrapCheckLog(c))
	modelmigration.RegisterOfferAccessImport(s.coordinator, clock.WallClock, loggertesting.WrapCheckLog(c))
	m := description.NewModel(description.ModelArgs{})
	a := m.AddApplication(description.ApplicationArgs{Name: "mysql"})
	a.AddOffer(description.ApplicationOfferArgs{
		OfferUUID: "cfa46843-ebf2-4fff-8519-c1fb5a9816f3", OfferName: "mysql",
		ApplicationName: "mysql", Endpoints: map[string]string{"db": "db"},
		ACL: map[string]string{"alice@external": "consume", "bob@external": "read", "carol@external": "admin"},
	})
	err := s.coordinator.Perform(c.Context(), s.scope, m)
	c.Assert(err, tc.ErrorIsNil)
	var expected []userAccess
	for name, access := range map[string]string{
		"alice@external": "consume", "bob@external": "read", "carol@external": "admin", "everyone@external": "read",
	} {
		username, err := user.NewName(name)
		c.Assert(err, tc.ErrorIsNil)
		imported, err := s.svc.GetUserByName(c.Context(), username)
		c.Assert(err, tc.ErrorIsNil)
		expected = append(expected, userAccess{GrantTo: imported.UUID.String(), GrantOn: "cfa46843-ebf2-4fff-8519-c1fb5a9816f3", AccessType: access})
	}
	c.Check(s.getPermissions(c, "v_permission_offer"), tc.SameContents, expected)
	c.Check(s.getPermissions(c, "v_permission_model"), tc.HasLen, 0)
}
