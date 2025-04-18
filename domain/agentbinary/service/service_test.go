// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/domain/agentbinary"
)

type serviceSuite struct {
	testing.IsolationSuite

	mockModelState, mockControllerState *MockAgentBinaryState
}

var _ = gc.Suite(&serviceSuite{})

func (s *serviceSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.mockModelState = NewMockAgentBinaryState(ctrl)
	s.mockControllerState = NewMockAgentBinaryState(ctrl)
	return ctrl
}

// TestListAgentBinaries tests the ListAgentBinaries method of the
// AgentBinaryService. It verifies that the method correctly merges
// agent binaries from the controller and model stores, with the model
// binaries taking precedence over the controller binaries.
func (s *serviceSuite) TestListAgentBinaries(c *gc.C) {
	defer s.setupMocks(c).Finish()

	controllerBinaries := []agentbinary.Metadata{
		{
			Version: "4.0.0",
			Size:    1,
			SHA256:  "sha256hash-1",
		},
		{
			Version: "4.0.1",
			Size:    2,
			SHA256:  "sha256hash-2",
		},
	}
	modelBinaries := []agentbinary.Metadata{
		{
			Version: "4.0.1",
			Size:    222,
			// A same SHA should never have a different size, but this is just for testing the merge logic.
			SHA256: "sha256hash-2",
		},
		{
			Version: "4.0.2",
			Size:    3,
			SHA256:  "sha256hash-3",
		},
	}
	expected := []agentbinary.Metadata{
		{
			Version: "4.0.0",
			Size:    1,
			SHA256:  "sha256hash-1",
		},
		{
			Version: "4.0.1",
			Size:    222,
			SHA256:  "sha256hash-2",
		},
		{
			Version: "4.0.2",
			Size:    3,
			SHA256:  "sha256hash-3",
		},
	}
	s.mockControllerState.EXPECT().ListAgentBinaries(gomock.Any()).Return(controllerBinaries, nil)
	s.mockModelState.EXPECT().ListAgentBinaries(gomock.Any()).Return(modelBinaries, nil)

	svc := NewAgentBinaryService(s.mockControllerState, s.mockModelState)
	result, err := svc.ListAgentBinaries(context.Background())
	c.Assert(err, gc.IsNil)
	c.Assert(result, jc.SameContents, expected)
}
