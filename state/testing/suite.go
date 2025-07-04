// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"context"
	"sync"
	"time"

	"github.com/juju/clock/testclock"
	mgotesting "github.com/juju/mgo/v3/testing"
	"github.com/juju/names/v6"
	"github.com/juju/tc"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/testhelpers"
	"github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/testing/factory"
	"github.com/juju/juju/state"
	statewatcher "github.com/juju/juju/state/watcher"
)

// StateSuite provides setup and teardown for tests that require a
// state.State.
type StateSuite struct {
	mgotesting.MgoSuite
	testing.BaseSuite
	NewPolicy                 state.NewPolicyFunc
	Controller                *state.Controller
	StatePool                 *state.StatePool
	State                     *state.State
	Model                     *state.Model
	Owner                     names.UserTag
	Factory                   *factory.Factory
	InitialConfig             *config.Config
	ControllerConfig          map[string]interface{}
	ControllerInheritedConfig map[string]interface{}
	ControllerModelType       state.ModelType
	RegionConfig              cloud.RegionConfig
	Clock                     testclock.AdvanceableClock
	modelWatcherIdle          chan string
	modelWatcherMutex         *sync.Mutex
	InstancePrechecker        func(*tc.C, *state.State) environs.InstancePrechecker
}

func (s *StateSuite) SetUpSuite(c *tc.C) {
	s.MgoSuite.SetUpSuite(c)
	s.BaseSuite.SetUpSuite(c)
}

func (s *StateSuite) TearDownSuite(c *tc.C) {
	s.BaseSuite.TearDownSuite(c)
	s.MgoSuite.TearDownSuite(c)
}

func (s *StateSuite) SetUpTest(c *tc.C) {
	s.MgoSuite.SetUpTest(c)
	s.BaseSuite.SetUpTest(c)

	s.modelWatcherIdle = nil
	s.modelWatcherMutex = &sync.Mutex{}
	s.PatchValue(&statewatcher.HubWatcherIdleFunc, s.hubWatcherIdleFunc)

	s.Owner = names.NewLocalUserTag("test-admin")

	if s.Clock == nil {
		s.Clock = testclock.NewDilatedWallClock(100 * time.Millisecond)
	}

	s.Controller = InitializeWithArgs(c, InitializeArgs{
		Owner:                     s.Owner,
		InitialConfig:             s.InitialConfig,
		ControllerConfig:          s.ControllerConfig,
		ControllerInheritedConfig: s.ControllerInheritedConfig,
		ControllerModelType:       s.ControllerModelType,
		RegionConfig:              s.RegionConfig,
		NewPolicy:                 s.NewPolicy,
		Clock:                     s.Clock,
	})
	s.AddCleanup(func(*tc.C) {
		_ = s.Controller.Close()
	})
	s.StatePool = s.Controller.StatePool()
	var err error
	s.State, err = s.StatePool.SystemState()
	c.Assert(err, tc.ErrorIsNil)
	model, err := s.State.Model()
	c.Assert(err, tc.ErrorIsNil)
	s.Model = model

	s.Factory = factory.NewFactory(s.State, s.StatePool, s.ControllerConfig)
}

func (s *StateSuite) TearDownTest(c *tc.C) {
	s.BaseSuite.TearDownTest(c)
	s.MgoSuite.TearDownTest(c)
}

func (s *StateSuite) hubWatcherIdleFunc(modelUUID string) {
	s.modelWatcherMutex.Lock()
	idleChan := s.modelWatcherIdle
	s.modelWatcherMutex.Unlock()
	if idleChan == nil {
		return
	}
	// There is a very small race condition between when the
	// idle channel is cleared and when the function exits.
	// Under normal circumstances, there is a goroutine in a tight loop
	// reading off the idle channel. If the channel isn't read
	// within a short wait, we don't send the message.
	select {
	case idleChan <- modelUUID:
	case <-time.After(testing.ShortWait):
	}
}

// WaitForModelWatchersIdle firstly waits for the txn poller to process
// all pending changes, then waits for the hub watcher on the state object
// to have finished processing all those events.
func (s *StateSuite) WaitForModelWatchersIdle(c *tc.C, modelUUID string) {
	// Use a logger rather than c.Log so we get timestamps.
	logger := loggertesting.WrapCheckLog(c)
	logger.Infof(context.TODO(), "waiting for model %s to be idle", modelUUID)
	// Create idle channel after the sync so as to be sure that at least
	// one sync is complete before signalling the idle timer.
	s.modelWatcherMutex.Lock()
	idleChan := make(chan string)
	s.modelWatcherIdle = idleChan
	s.modelWatcherMutex.Unlock()

	defer func() {
		s.modelWatcherMutex.Lock()
		s.modelWatcherIdle = nil
		s.modelWatcherMutex.Unlock()
		// Clear out any pending events.
		for {
			select {
			case <-idleChan:
			default:
				return
			}
		}
	}()

	timeout := time.After(testhelpers.LongWait)
	for {
		loop := time.After(10 * time.Millisecond)
		select {
		case <-loop:
		case uuid := <-idleChan:
			if uuid == modelUUID {
				return
			} else {
				logger.Infof(context.TODO(), "model %s is idle", uuid)
			}
		case <-timeout:
			c.Fatal("no sync event sent, is the watcher dead?")
		}
	}
}
