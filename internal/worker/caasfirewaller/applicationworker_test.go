// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasfirewaller_test

import (
	"context"
	"time"

	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v4"
	"github.com/juju/worker/v4/workertest"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/caas"
	caasmocks "github.com/juju/juju/caas/mocks"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/watchertest"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/worker/caasfirewaller"
	"github.com/juju/juju/internal/worker/caasfirewaller/mocks"
)

type appWorkerSuite struct {
	testing.BaseSuite

	appName string

	firewallerAPI *mocks.MockCAASFirewallerAPI
	lifeGetter    *mocks.MockLifeGetter
	broker        *mocks.MockCAASBroker
	brokerApp     *caasmocks.MockApplication

	applicationChanges chan struct{}
	portsChanges       chan []string

	appsWatcher  watcher.NotifyWatcher
	portsWatcher watcher.StringsWatcher
}

var _ = gc.Suite(&appWorkerSuite{})

func (s *appWorkerSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	s.appName = "app1"
	s.applicationChanges = make(chan struct{})
	s.portsChanges = make(chan []string)
}

func (s *appWorkerSuite) getController(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.appsWatcher = watchertest.NewMockNotifyWatcher(s.applicationChanges)
	s.portsWatcher = watchertest.NewMockStringsWatcher(s.portsChanges)

	s.firewallerAPI = mocks.NewMockCAASFirewallerAPI(ctrl)

	s.lifeGetter = mocks.NewMockLifeGetter(ctrl)
	s.broker = mocks.NewMockCAASBroker(ctrl)
	s.brokerApp = caasmocks.NewMockApplication(ctrl)

	return ctrl
}

func (s *appWorkerSuite) getWorker(c *gc.C) worker.Worker {
	w, err := caasfirewaller.NewApplicationWorker(
		testing.ControllerTag.Id(),
		testing.ModelTag.Id(),
		s.appName,
		s.firewallerAPI,
		s.broker,
		s.lifeGetter,
		loggertesting.WrapCheckLog(c),
	)
	c.Assert(err, jc.ErrorIsNil)
	return w
}

func (s *appWorkerSuite) TestWorker(c *gc.C) {
	ctrl := s.getController(c)
	defer ctrl.Finish()

	done := make(chan struct{})

	go func() {
		// 1st port change event.
		s.portsChanges <- []string{s.appName}
		// 2nd port change event.
		s.portsChanges <- []string{s.appName}
		// 3rd port change event, including another application.
		s.portsChanges <- []string{s.appName, "another-app"}

		// port change event for another application.
		s.portsChanges <- []string{"another-app"}

		s.applicationChanges <- struct{}{}
	}()

	gpr1 := network.GroupedPortRanges{
		"": []network.PortRange{
			network.MustParsePortRange("1000/tcp"),
		},
	}

	gpr2 := network.GroupedPortRanges{
		"": []network.PortRange{
			network.MustParsePortRange("1000/tcp"),
		},
		"monitoring-port": []network.PortRange{
			network.MustParsePortRange("2000/udp"),
		},
	}

	gomock.InOrder(
		s.firewallerAPI.EXPECT().WatchApplication(gomock.Any(), s.appName).Return(s.appsWatcher, nil),
		s.firewallerAPI.EXPECT().WatchOpenedPorts(gomock.Any()).Return(s.portsWatcher, nil),
		s.broker.EXPECT().Application(s.appName, caas.DeploymentStateful).Return(s.brokerApp),

		// initial fetch.
		s.firewallerAPI.EXPECT().GetOpenedPorts(gomock.Any(), s.appName).Return(network.GroupedPortRanges{}, nil),

		// 1st triggerred by port change event.
		s.firewallerAPI.EXPECT().GetOpenedPorts(gomock.Any(), s.appName).Return(gpr1, nil),
		s.brokerApp.EXPECT().UpdatePorts([]caas.ServicePort{
			{
				Name:       "1000-tcp",
				Port:       1000,
				TargetPort: 1000,
				Protocol:   "tcp",
			},
		}, false).Return(nil),

		// 2nd triggerred by port change event, no UpdatePorts because no diff on the portchanges.
		s.firewallerAPI.EXPECT().GetOpenedPorts(gomock.Any(), s.appName).Return(gpr1, nil),

		// 1rd triggerred by port change event.
		s.firewallerAPI.EXPECT().GetOpenedPorts(gomock.Any(), s.appName).Return(gpr2, nil),
		s.brokerApp.EXPECT().UpdatePorts([]caas.ServicePort{
			{
				Name:       "1000-tcp",
				Port:       1000,
				TargetPort: 1000,
				Protocol:   "tcp",
			},
			{
				Name:       "2000-udp",
				Port:       2000,
				TargetPort: 2000,
				Protocol:   "udp",
			},
		}, false).Return(nil),

		s.firewallerAPI.EXPECT().IsExposed(gomock.Any(), s.appName).DoAndReturn(func(_ context.Context, _ string) (bool, error) {
			close(done)
			return false, nil
		}),
	)

	w := s.getWorker(c)

	select {
	case <-done:
	case <-time.After(testing.ShortWait):
		c.Errorf("timed out waiting for worker")
	}
	workertest.CleanKill(c, w)
}
