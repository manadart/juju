// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"
	"testing"

	"github.com/juju/tc"

	schematesting "github.com/juju/juju/domain/schema/testing"
	"github.com/juju/juju/internal/uuid"
)

type stateSuite struct {
	schematesting.ModelSuite
}

func TestStateSuite(t *testing.T) {
	tc.Run(t, &stateSuite{})
}

func (s *stateSuite) TestGetEnvironment(c *tc.C) {
	controllerUUID := uuid.MustNewUUID().String()

	err := s.TxnRunner().StdTxn(c.Context(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
INSERT INTO model (uuid, controller_uuid, name, qualifier, type, cloud, cloud_type)
VALUES (?, ?, ?, 'prod', 'iaas', 'aws', 'ec2')
`, uuid.MustNewUUID().String(), controllerUUID, "my-model")
		return err
	})
	c.Assert(err, tc.ErrorIsNil)

	st := NewState(s.TxnRunnerFactory())
	environment, err := st.GetEnvironment(c.Context())
	c.Assert(err, tc.ErrorIsNil)
	c.Check(environment, tc.DeepEquals, map[string]string{
		"model-name":      "my-model",
		"controller-uuid": controllerUUID,
	})
}
