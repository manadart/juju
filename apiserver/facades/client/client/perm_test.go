// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client_test

import (
	"context"
	"strings"
	"time"

	"github.com/juju/names/v5"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/annotations"
	"github.com/juju/juju/api/client/application"
	apiclient "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/api/client/modelconfig"
	apicharm "github.com/juju/juju/api/common/charm"
	"github.com/juju/juju/apiserver/facades/client/client"
	"github.com/juju/juju/core/constraints"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/juju/testing"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
)

type permSuite struct {
	baseSuite
}

var _ = gc.Suite(&permSuite{})

// Most (if not all) of the permission tests below aim to test
// end-to-end operations execution through the API, but do not care
// about the results. They only test that a call is succeeds or fails
// (usually due to "permission denied"). There are separate test cases
// testing each individual API call data flow later on.

func allowed(allow []names.Tag) map[names.Tag]bool {
	p := make(map[names.Tag]bool)
	if allow != nil {
		for _, e := range allow {
			p[e] = true
		}
		return p
	}
	return p
}

func (s *permSuite) SetUpTest(c *gc.C) {
	s.baseSuite.SetUpTest(c)
	client.SkipReplicaCheck(s)
}

type watcherPermSuite struct {
	permSuite
}

var _ = gc.Suite(&watcherPermSuite{})

func (s *permSuite) TestOperationPermClientSetApplicationConstraints(c *gc.C) {
	s.testOperationPerm(c, opClientSetApplicationConstraints)
}

func (s *permSuite) TestOperationPermClientSetModelConstraints(c *gc.C) {
	s.testOperationPerm(c, opClientSetModelConstraints)
}

func (s *permSuite) TestOperationPermClientModelGet(c *gc.C) {
	s.testOperationPerm(c, opClientModelGet)
}

func (s *permSuite) TestOperationPermClientModelSet(c *gc.C) {
	s.testOperationPerm(c, opClientModelSet)
}

func (s *permSuite) TestOperationPermApplicationAddRelation(c *gc.C) {
	s.testOperationPerm(c, opClientAddRelation)
}

func (s *permSuite) TestOperationPermApplicationDestroyRelation(c *gc.C) {
	s.testOperationPerm(c, opClientDestroyRelation)
}

func (s *permSuite) TestOperationPermApplicationGetConstraints(c *gc.C) {
	s.testOperationPerm(c, opClientGetApplicationConstraints)
}

func (s *permSuite) TestOperationPermDestroyUnits(c *gc.C) {
	s.testOperationPerm(c, opClientDestroyUnit)
}

func (s *permSuite) TestOperationPermApplicationAddUnits(c *gc.C) {
	s.testOperationPerm(c, opClientAddApplicationUnits)
}

func (s *permSuite) TestOperationPermApplicationGet(c *gc.C) {
	s.testOperationPerm(c, opClientApplicationGet)
}

func (s *permSuite) TestOperationPermAnnotationsGetAnnotations(c *gc.C) {
	s.testOperationPerm(c, opClientGetAnnotations)
}

func (s *permSuite) TestOperationPermClientStatus(c *gc.C) {
	s.testOperationPerm(c, opClientStatus)
}

func (s *permSuite) TestOperationPermApplicationResolveUnitErrors(c *gc.C) {
	s.testOperationPerm(c, opClientResolved)
}

func (s *permSuite) TestOperationPermApplicationExpose(c *gc.C) {
	s.testOperationPerm(c, opClientApplicationExpose)
}

func (s *permSuite) TestOperationPermApplicationUnexpose(c *gc.C) {
	s.testOperationPerm(c, opClientApplicationUnexpose)
}

func (s *permSuite) TestOperationPermAnnotationsSetAnnotations(c *gc.C) {
	s.testOperationPerm(c, opClientSetAnnotations)
}

func (s *permSuite) TestOperationPermApplicationDestroyUnits(c *gc.C) {
	s.testOperationPerm(c, opClientDestroyApplicationUnits)
}

func (s *permSuite) TestOperationPermApplicationDestroy(c *gc.C) {
	s.testOperationPerm(c, opClientApplicationDestroy)
}

func (s *permSuite) TestOperationPermApplicationDestroyApplication(c *gc.C) {
	s.testOperationPerm(c, opClientDestroyApplication)
}

func (s *permSuite) TestOperationPermApplicationSetCharm(c *gc.C) {
	s.testOperationPerm(c, opClientApplicationSetCharm)
}

func (s *permSuite) testOperationPerm(
	c *gc.C,
	op func(c *gc.C, st api.Connection, mst *state.State) (reset func(), err error),
) {
	allow := allowed([]names.Tag{testing.AdminUser, names.NewLocalUserTag("other")})
	for j, e := range s.setUpScenario(c) {
		c.Logf("\n------\ntest %d; entity %q", j, e)
		st := s.openAs(c, e)
		reset, err := op(c, st, s.ControllerModel(c).State())
		if allow[e] {
			c.Check(err, jc.ErrorIsNil)
		} else {
			c.Check(err, jc.Satisfies, params.IsCodeUnauthorized)
		}
		reset()
		_ = st.Close()
	}
}

func opClientAddRelation(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).AddRelation(context.Background(), []string{"nosuch1", "nosuch2"}, nil)
	if params.IsCodeNotFound(err) {
		err = nil
	}
	return func() {}, err
}

func opClientDestroyRelation(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	err := application.NewClient(st).DestroyRelation(context.Background(), (*bool)(nil), (*time.Duration)(nil), "nosuch1", "nosuch2")
	if params.IsCodeNotFound(err) {
		err = nil
	}
	return func() {}, err
}

func opClientStatus(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	status, err := apiclient.NewClient(st, loggertesting.WrapCheckLog(c)).Status(context.Background(), nil)
	if err != nil {
		c.Check(status, gc.IsNil)
		return func() {}, err
	}
	clearSinceTimes(status)
	clearSinceTimes(scenarioStatus)
	clearContollerTimestamp(status)
	c.Assert(status, jc.DeepEquals, scenarioStatus)
	return func() {}, nil
}

func opClientApplicationGet(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).Get(context.Background(), "wordpress")
	if err != nil {
		return func() {}, err
	}
	return func() {}, nil
}

func opClientApplicationExpose(c *gc.C, st api.Connection, mst *state.State) (func(), error) {
	err := application.NewClient(st).Expose(context.Background(), "wordpress", nil)
	if err != nil {
		return func() {}, err
	}
	return func() {
		svc, err := mst.Application("wordpress")
		c.Assert(err, jc.ErrorIsNil)
		err = svc.ClearExposed()
		c.Assert(err, jc.ErrorIsNil)
	}, nil
}

func opClientApplicationUnexpose(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	err := application.NewClient(st).Unexpose(context.Background(), "wordpress", nil)
	if err != nil {
		return func() {}, err
	}
	return func() {}, nil
}

func opClientResolved(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	err := application.NewClient(st).ResolveUnitErrors(context.Background(), []string{"wordpress/1"}, false, false)
	// There are several scenarios in which this test is called, one is
	// that the user is not authorized.  In that case we want to exit now,
	// letting the error percolate out so the caller knows that the
	// permission error was correctly generated.
	if err != nil && params.IsCodeUnauthorized(err) {
		return func() {}, err
	}
	// Otherwise, the user was authorized, but we expect an error anyway
	// because the unit is not in an error state when we tried to resolve
	// the error.  Therefore, since it is complaining it means that the
	// call to Resolved worked, so we're happy.
	c.Assert(err, gc.NotNil)
	c.Assert(err.Error(), gc.Equals, `unit "wordpress/1" is not in an error state`)
	return func() {}, nil
}

func opClientGetAnnotations(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	ann, err := annotations.NewClient(st).Get(context.Background(), []string{"application-wordpress"})
	if err != nil {
		return func() {}, err
	}
	c.Assert(ann, gc.HasLen, 1)
	c.Assert(ann[0].EntityTag, gc.Equals, "application-wordpress")
	c.Assert(ann[0].Annotations, gc.HasLen, 0)
	return func() {}, nil
}

func opClientSetAnnotations(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	pairs := map[string]string{"key1": "value1", "key2": "value2"}
	setParams := map[string]map[string]string{
		"application-wordpress": pairs,
	}
	_, err := annotations.NewClient(st).Set(context.Background(), setParams)
	if err != nil {
		return func() {}, err
	}
	return func() {
		pairs := map[string]string{"key1": "", "key2": ""}
		setParams := map[string]map[string]string{
			"application-wordpress": pairs,
		}
		_, err := annotations.NewClient(st).Set(context.Background(), setParams)
		c.Assert(err, jc.ErrorIsNil)
	}, nil
}

func opClientApplicationSetCharm(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	cfg := application.SetCharmConfig{
		ApplicationName: "nosuch",
		CharmID: application.CharmID{
			URL:    "local:wordpress",
			Origin: apicharm.Origin{Source: "local"},
		},
	}
	err := application.NewClient(st).SetCharm(context.Background(), cfg)
	if params.IsCodeNotFound(err) {
		err = nil
	}
	return func() {}, err
}

func opClientAddApplicationUnits(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).AddUnits(context.Background(), application.AddUnitsParams{
		ApplicationName: "nosuch",
		NumUnits:        1,
	})
	if params.IsCodeNotFound(err) {
		err = nil
	}
	return func() {}, err
}

func opClientDestroyApplicationUnits(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).DestroyUnits(
		context.Background(),
		application.DestroyUnitsParams{Units: []string{"wordpress/99"}})
	if err != nil && strings.HasPrefix(err.Error(), "no units were destroyed") {
		err = nil
	}
	return func() {}, err
}

func opClientDestroyUnit(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).DestroyUnits(context.Background(),
		application.DestroyUnitsParams{
			Units: []string{"wordpress/99"},
		})
	return func() {}, err
}

func opClientApplicationDestroy(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).DestroyApplications(
		context.Background(),
		application.DestroyApplicationsParams{Applications: []string{"non-existent"}})
	if params.IsCodeNotFound(err) {
		err = nil
	}
	return func() {}, err
}

func opClientDestroyApplication(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).DestroyApplications(
		context.Background(),
		application.DestroyApplicationsParams{
			Applications: []string{"non-existent"},
		})
	return func() {}, err
}

func opClientGetApplicationConstraints(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := application.NewClient(st).GetConstraints(context.Background(), "wordpress")
	return func() {}, err
}

func opClientSetApplicationConstraints(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	nullConstraints := constraints.Value{}
	err := application.NewClient(st).SetConstraints(context.Background(), "wordpress", nullConstraints)
	if err != nil {
		return func() {}, err
	}
	return func() {}, nil
}

func opClientSetModelConstraints(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	nullConstraints := constraints.Value{}
	err := modelconfig.NewClient(st).SetModelConstraints(context.Background(), nullConstraints)
	if err != nil {
		return func() {}, err
	}
	return func() {}, nil
}

func opClientModelGet(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	_, err := modelconfig.NewClient(st).ModelGet(context.Background())
	if err != nil {
		return func() {}, err
	}
	return func() {}, nil
}

func opClientModelSet(c *gc.C, st api.Connection, _ *state.State) (func(), error) {
	args := map[string]interface{}{"some-key": "some-value"}
	err := modelconfig.NewClient(st).ModelSet(context.Background(), args)
	if err != nil {
		return func() {}, err
	}
	return func() {
		args["some-key"] = nil
		modelconfig.NewClient(st).ModelSet(context.Background(), args)
	}, nil
}
