// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmigration_test

import (
	"github.com/juju/description/v12"
	"github.com/juju/tc"

	"github.com/juju/juju/core/annotations"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/core/version"
	annotationservice "github.com/juju/juju/domain/annotation/service"
	annotationstate "github.com/juju/juju/domain/annotation/state"
)

// TestAnnotations checks every entity kind whose annotations 3.6 exports.
// Compare exact maps through the annotation service after the complete import,
// including values shared by key but different by entity. This catches both
// omission and accidental attachment to the wrong newly allocated UUID.
func (s *legacyImportSuite) TestAnnotations(c *tc.C) {
	desc, coordinator, scope := s.setupImport(c)
	modelAnnotations := map[string]string{"owner": "model-owner", "purpose": "production"}
	machineAnnotations := map[string]string{"owner": "machine-owner", "rack": "rack-3"}
	applicationAnnotations := map[string]string{"owner": "application-owner", "team": "database"}
	unitAnnotations := map[string]string{"owner": "unit-owner", "ticket": "42"}
	desc.SetAnnotations(modelAnnotations)
	machine := desc.AddMachine(description.MachineArgs{Id: "0", Base: "ubuntu@22.04"})
	machine.SetAnnotations(machineAnnotations)
	machine.SetConstraints(description.ConstraintsArgs{Architecture: "amd64"})
	machine.SetTools(description.AgentToolsArgs{Version: version.Current.String() + "-ubuntu-amd64"})
	machine.SetStatus(description.StatusArgs{Value: "started"})
	machine.SetInstance(description.CloudInstanceArgs{InstanceId: "instance-0", Architecture: "amd64"})
	machine.Instance().SetStatus(description.StatusArgs{Value: "running"})
	app := addLegacyApplication(desc, "client")
	app.SetAnnotations(applicationAnnotations)
	unit := app.AddUnit(description.UnitArgs{Name: "client/0", Type: model.IAAS.String(), Machine: "0"})
	unit.SetAnnotations(unitAnnotations)
	unit.SetTools(description.AgentToolsArgs{Version: version.Current.String() + "-ubuntu-amd64"})
	unit.SetAgentStatus(description.StatusArgs{Value: "idle"})
	unit.SetWorkloadStatus(description.StatusArgs{Value: "active"})
	c.Assert(coordinator.Perform(c.Context(), scope, desc), tc.ErrorIsNil)
	svc := annotationservice.NewService(annotationstate.NewState(scope.ModelDB()))
	for _, test := range []struct {
		id   annotations.ID
		want map[string]string
	}{
		{annotations.ID{Kind: annotations.KindModel, Name: desc.UUID()}, modelAnnotations},
		{annotations.ID{Kind: annotations.KindMachine, Name: "0"}, machineAnnotations},
		{annotations.ID{Kind: annotations.KindApplication, Name: "client"}, applicationAnnotations},
		{annotations.ID{Kind: annotations.KindUnit, Name: "client/0"}, unitAnnotations},
	} {
		got, err := svc.GetAnnotations(c.Context(), test.id)
		c.Assert(err, tc.ErrorIsNil)
		c.Check(got, tc.DeepEquals, test.want, tc.Commentf("annotations for %s", test.id))
	}
}
