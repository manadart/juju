// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package osenv_test

import (
	stdtesting "testing"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	coretesting "github.com/juju/juju/internal/testing"
)

func Test(t *stdtesting.T) {
	gc.TestingT(t)
}

type importSuite struct {
}

var _ = gc.Suite(&importSuite{})

func (*importSuite) TestDependencies(c *gc.C) {
	c.Assert(coretesting.FindJujuCoreImports(c, "github.com/juju/juju/juju/osenv"), jc.SameContents, []string{
		"internal/featureflag",
	})
}
