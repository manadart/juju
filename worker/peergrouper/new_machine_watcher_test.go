package peergrouper

import (
	jc "github.com/juju/testing/checkers"
	"github.com/golang/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/state"
)

var noOp = func(){}

func newNotifyWatcher(ctrl *gomock.Controller) (state.NotifyWatcher, *func()) {
	w := NewMockNotifyWatcher(ctrl)

	ch := make(chan struct{})
	change := func(){ch <- struct{}{}}
	op := &change

	w.EXPECT().Wait().Do(func(){*op = noOp}).Return(nil).MaxTimes(1)
	w.EXPECT().Changes().Return(ch).AnyTimes()

	return w, op
}

func newWatchableMachine(ctrl *gomock.Controller) Machine {
	m := NewMockMachine(ctrl)

	var notifyOps []*func()

	watch := func() state.NotifyWatcher {
		w, op := newNotifyWatcher(ctrl)
		notifyOps = append(notifyOps, op)
		return w
	}
	m.EXPECT().Watch().DoAndReturn(watch).AnyTimes()

	hasVoteCall := m.EXPECT().HasVote().AnyTimes()
	hasVoteCall.Return(false)

	setHasVote := func(hasVote bool) {
		hasVoteCall.Return(hasVote)
		for _, op := range notifyOps {
			f := *op
			f()
		}
	}
	m.EXPECT().SetHasVote(gomock.Any()).Do(setHasVote).AnyTimes()

	return m
}

func (s *newSuite) TestMachineWatcher(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	m := newWatchableMachine(ctrl)
	w := m.Watch()

	stop := make(chan struct{})
	var changes int
	go func() {
		for {
			select {
			case <-w.Changes():
				changes++
			case <-stop:
				return
			}
		}
	}()

	m.SetHasVote(true)
	c.Check(m.HasVote(), jc.IsTrue)

	m.SetHasVote(false)
	c.Check(m.HasVote(), jc.IsFalse)

	w.Wait()

	// After the watcher is done, this should not cause a notification.
	m.SetHasVote(true)
	c.Check(m.HasVote(), jc.IsTrue)

	stop <- struct{}{}

	c.Check(changes, gc.Equals, 2)
}

func (s *newSuite) TestMachineMultiWatcher(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	m := newWatchableMachine(ctrl)
	w1 := m.Watch()
	w2 := m.Watch()

	stop := make(chan struct{})
	var changes1 int
	var changes2 int
	go func() {
		for {
			select {
			case <-w1.Changes():
				changes1++
			case <-w2.Changes():
				changes2++
			case <-stop:
				return
			}
		}
	}()

	m.SetHasVote(true)
	c.Check(m.HasVote(), jc.IsTrue)

	m.SetHasVote(false)
	c.Check(m.HasVote(), jc.IsFalse)

	w1.Wait()

	// Only watcher 2 should see a notification for this.
	m.SetHasVote(true)
	c.Check(m.HasVote(), jc.IsTrue)

	w2.Wait()

	stop <- struct{}{}

	c.Check(changes1, gc.Equals, 2)
	c.Check(changes2, gc.Equals, 3)
}
