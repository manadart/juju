// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package crossmodelrelation_test

import (
	"github.com/juju/description/v12"
	"github.com/juju/tc"

	"github.com/juju/juju/core/model"
)

// TestImportLegacyConsumedOfferWithoutRelation checks the valid 3.6 state after
// consume but before integrate. Registration has not assigned a remote entity
// token yet. Migration must preserve the consumed offer so it can be integrated
// later, including the macaroon needed to authenticate to the offering model.
func (s *importSuite) TestImportLegacyConsumedOfferWithoutRelation(c *tc.C) {
	m := description.NewModel(description.ModelArgs{Type: model.IAAS.String()})
	mac := newMacaroon(c, "unrelated-offer")
	macJSON, err := mac.MarshalJSON()
	c.Assert(err, tc.ErrorIsNil)
	const offerUUID = "cfa46843-ebf2-4fff-8519-c1fb5a9816f3"
	const modelUUID = "4ddd6454-931d-4278-8779-b0b7208994d9"
	const offerURL = "source:admin/model.database"
	app := m.AddRemoteApplication(description.RemoteApplicationArgs{
		Name: "database", OfferUUID: offerUUID, SourceModelUUID: modelUUID,
		URL: offerURL, Macaroon: string(macJSON),
	})
	app.AddEndpoint(description.RemoteEndpointArgs{Name: "db", Role: "provider", Interface: "db"})
	coordinator, scope, svc := s.setupCoordinatorScopeAndService(c)
	c.Assert(coordinator.Perform(c.Context(), scope, m), tc.ErrorIsNil)
	offers, err := svc.GetRemoteApplicationOfferers(c.Context())
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(offers, tc.HasLen, 1)
	c.Check(offers[0].ApplicationName, tc.Equals, "database")
	c.Check(offers[0].OfferUUID, tc.Equals, offerUUID)
	c.Check(offers[0].OffererModelUUID, tc.Equals, modelUUID)
	c.Check(offers[0].OfferURL, tc.Equals, offerURL)
	c.Check(offers[0].Macaroon, tc.DeepEquals, mac)
}
