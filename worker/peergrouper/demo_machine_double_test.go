package peergrouper

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	coretesting "github.com/juju/juju/testing"
	"github.com/golang/mock/gomock"
)

type newSuite struct {
	coretesting.BaseSuite
	//hub   Hub
}

var _ = gc.Suite(&newSuite{})

func newMachine(ctrl *gomock.Controller) Machine {
	m := NewMockMachine(ctrl)

	hasVoteCall := m.EXPECT().HasVote().AnyTimes()
	hasVoteCall.Return(false)

	setHasVote := func(hasVote bool) {
		hasVoteCall.Return(hasVote)
	}
	m.EXPECT().SetHasVote(gomock.Any()).Do(setHasVote).AnyTimes()

	return m
}

func (s *workerSuite) TestMachineDouble(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	m := newMachine(ctrl)
	c.Check(m.HasVote(), jc.IsFalse)

	m.SetHasVote(true)
	c.Check(m.HasVote(), jc.IsTrue)
	c.Check(m.HasVote(), jc.IsTrue)

	m.SetHasVote(false)
	c.Check(m.HasVote(), jc.IsFalse)
}

