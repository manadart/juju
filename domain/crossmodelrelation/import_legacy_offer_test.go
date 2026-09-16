// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package crossmodelrelation_test

import (
	"context"
	"database/sql"

	"github.com/juju/clock"
	"github.com/juju/description/v12"
	"github.com/juju/tc"

	"github.com/juju/juju/core/model"
	coremodelmigration "github.com/juju/juju/core/modelmigration"
	applicationmigration "github.com/juju/juju/domain/application/modelmigration"
	cmrmigration "github.com/juju/juju/domain/crossmodelrelation/modelmigration"
	migrationtesting "github.com/juju/juju/domain/modelmigration/testing"
	relationmigration "github.com/juju/juju/domain/relation/modelmigration"
	sequencemigration "github.com/juju/juju/domain/sequence/modelmigration"
	loggertesting "github.com/juju/juju/internal/logger/testing"
)

// TestImportLegacyOfferRelationPreservesIdentityAndData checks the offering side
// of an established 3.6 CMR before workers can republish any data. Reallocating
// the numeric relation ID breaks agent checkpoints; losing settings or scope
// membership changes what existing hooks see after migration.
func (s *importSuite) TestImportLegacyOfferRelationPreservesIdentityAndData(c *tc.C) {
	m := description.NewModel(description.ModelArgs{Type: model.IAAS.String()})
	a := m.AddApplication(description.ApplicationArgs{Name: "mysql", CharmURL: "ch:mysql-1"})
	a.SetCharmOrigin(description.CharmOriginArgs{
		Source: "charm-hub", ID: "deadbeef", Hash: "deadbeef2", Revision: 1,
		Channel: "latest/stable", Platform: "amd64/ubuntu/20.04",
	})
	a.SetCharmMetadata(description.CharmMetadataArgs{
		Name: "mysql", Provides: map[string]description.CharmMetadataRelation{
			"db": migrationtesting.Relation{Name_: "db", Role_: "provider", InterfaceName_: "db", Scope_: "global"},
		},
	})
	a.SetCharmManifest(description.CharmManifestArgs{Bases: []description.CharmManifestBase{
		migrationtesting.ManifestBase{Name_: "ubuntu", Channel_: "stable", Architectures_: []string{"amd64"}},
	}})
	const offerUUID = "cfa46843-ebf2-4fff-8519-c1fb5a9816f3"
	const consumerUUID = "13ea2791-5e78-40d8-88c5-e9451444b45d"
	const relationUUID = "6049aa01-76c9-462d-8440-964a6e26aac2"
	const remote = "remote-13ea27915e7840d888c5e9451444b45d"
	a.AddOffer(description.ApplicationOfferArgs{
		OfferUUID: offerUUID, OfferName: "mysql", ApplicationName: "mysql", Endpoints: map[string]string{"db": "db"},
	})
	rapp := m.AddRemoteApplication(description.RemoteApplicationArgs{
		Name: remote, IsConsumerProxy: true, SourceModelUUID: "4ddd6454-931d-4278-8779-b0b7208994d9",
	})
	rapp.AddEndpoint(description.RemoteEndpointArgs{Name: "db", Role: "requirer", Interface: "db"})
	r := m.AddRelation(description.RelationArgs{Id: 42, Key: remote + ":db mysql:db"})
	ep := r.AddEndpoint(description.EndpointArgs{ApplicationName: "mysql", Name: "db", Role: "provider", Interface: "db"})
	ep.SetApplicationSettings(map[string]any{"password": "keep-me"})
	ep = r.AddEndpoint(description.EndpointArgs{ApplicationName: remote, Name: "db", Role: "requirer", Interface: "db"})
	ep.SetApplicationSettings(map[string]any{"database": "keep-me-too"})
	ep.SetUnitSettings(remote+"/0", map[string]any{"request": "keep-unit-data"})
	m.AddRemoteEntity(description.RemoteEntityArgs{ID: "application-" + remote, Token: consumerUUID})
	m.AddRemoteEntity(description.RemoteEntityArgs{ID: "relation-" + remote + ".db#mysql.db", Token: relationUUID})
	m.AddOfferConnection(description.OfferConnectionArgs{
		OfferUUID: offerUUID, RelationID: 42, RelationKey: r.Key(),
		SourceModelUUID: "4ddd6454-931d-4278-8779-b0b7208994d9", UserName: "admin",
	})
	_, scope, _ := s.setupCoordinatorScopeAndService(c)
	logger := loggertesting.WrapCheckLog(c)
	coordinator := coremodelmigration.NewCoordinator(logger)
	// Keep the production order, including conversion of the 3.6 sequence
	// (next ID) to the 4.0 sequence (last allocated ID).
	m.SetSequence("relation", 43)
	sequencemigration.RegisterImport(coordinator)
	applicationmigration.RegisterImport(coordinator, clock.WallClock, logger)
	cmrmigration.RegisterImport(coordinator, clock.WallClock, logger)
	relationmigration.RegisterImport(coordinator, clock.WallClock, logger)
	c.Assert(coordinator.Perform(c.Context(), scope, m), tc.ErrorIsNil)

	runner, err := scope.ModelDB()(c.Context())
	c.Assert(err, tc.ErrorIsNil)
	var id int
	appSettings := make(map[string]string)
	unitSettings := make(map[string]string)
	var units []string
	err = runner.StdTxn(c.Context(), func(ctx context.Context, tx *sql.Tx) error {
		appSettings = make(map[string]string)
		unitSettings = make(map[string]string)
		units = nil
		if err := tx.QueryRowContext(ctx,
			"SELECT relation_id FROM relation WHERE uuid = ?", relationUUID).Scan(&id); err != nil {
			return err
		}
		rows, err := tx.QueryContext(ctx, `
SELECT a.name, s.key, s.value
FROM relation_application_setting AS s
JOIN relation_endpoint AS re ON re.uuid = s.relation_endpoint_uuid
JOIN application_endpoint AS ae ON ae.uuid = re.endpoint_uuid
JOIN application AS a ON a.uuid = ae.application_uuid
WHERE re.relation_uuid = ?`, relationUUID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var app, key, value string
			if err := rows.Scan(&app, &key, &value); err != nil {
				return err
			}
			appSettings[app+":"+key] = value
		}
		if err := rows.Err(); err != nil {
			return err
		}
		rows, err = tx.QueryContext(ctx, `
SELECT u.name, s.key, s.value
FROM relation_unit_setting AS s
JOIN relation_unit AS ru ON ru.uuid = s.relation_unit_uuid
JOIN relation_endpoint AS re ON re.uuid = ru.relation_endpoint_uuid
JOIN unit AS u ON u.uuid = ru.unit_uuid
WHERE re.relation_uuid = ?`, relationUUID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var unit, key, value string
			if err := rows.Scan(&unit, &key, &value); err != nil {
				return err
			}
			unitSettings[unit+":"+key] = value
		}
		if err := rows.Err(); err != nil {
			return err
		}
		rows, err = tx.QueryContext(ctx, `
SELECT u.name
FROM relation_unit AS ru
JOIN relation_endpoint AS re ON re.uuid = ru.relation_endpoint_uuid
JOIN unit AS u ON u.uuid = ru.unit_uuid
WHERE re.relation_uuid = ?`, relationUUID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return err
			}
			units = append(units, name)
		}
		return rows.Err()
	})
	c.Assert(err, tc.ErrorIsNil)
	c.Check(id, tc.Equals, 42)
	c.Check(appSettings, tc.DeepEquals, map[string]string{
		"mysql:password": "keep-me", remote + ":database": "keep-me-too",
	})
	c.Check(unitSettings, tc.DeepEquals, map[string]string{
		remote + "/0:request": "keep-unit-data",
	})
	c.Check(units, tc.SameContents, []string{remote + "/0"})
}
