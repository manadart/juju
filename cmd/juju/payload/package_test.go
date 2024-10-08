// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package payload

import (
	"context"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/jujuclient/jujuclienttesting"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

func NewListCommandForTest(newClient func(ctx context.Context) (ListAPI, error)) *ListCommand {
	cmd := &ListCommand{newAPIClient: newClient}
	cmd.SetClientStore(jujuclienttesting.MinimalStore())
	return cmd
}
