// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package openfga embeds OpenFGA using Juju's database transaction boundary.
// The storage contract is pinned to github.com/openfga/openfga v1.22.1.
package openfga

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/openfga/openfga/pkg/storage"

	"github.com/juju/juju/core/database"
	internaldatabase "github.com/juju/juju/internal/database"
)

// StoreID and ModelID identify Juju's persisted authorisation store and model.
// Changing the model requires a new model ID and a schema migration.
const (
	StoreID = "01J00000000000000000000000"
	ModelID = "01J00000000000000000000001"
)

// Datastore owns no connections or background work. The database accessor
// supplies the current connection and retries complete transaction closures.
type Datastore struct {
	factory database.TxnRunnerFactory
}

var _ storage.OpenFGADatastore = (*Datastore)(nil)

// NewDatastore constructs a datastore backed by the controller database.
func NewDatastore(factory database.TxnRunnerFactory) *Datastore {
	return &Datastore{factory: factory}
}

func (d *Datastore) transact(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {
	db, err := d.factory(ctx)
	if err != nil {
		return err
	}
	// Preserve driver errors until the accessor has finished retrying.
	err = db.StdTxn(ctx, fn)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNotFound
	}
	return err
}

// Close leaves the accessor-owned database open.
func (d *Datastore) Close() {}

// IsReady checks the accessor lifecycle and the installed schema.
func (d *Datastore) IsReady(ctx context.Context) (storage.ReadinessStatus, error) {
	db, err := d.factory(ctx)
	if err != nil {
		return storage.ReadinessStatus{}, err
	}
	select {
	case <-ctx.Done():
		return storage.ReadinessStatus{}, ctx.Err()
	case <-db.Dying():
		return storage.ReadinessStatus{Message: "controller database is shutting down"}, nil
	default:
	}
	err = db.StdTxn(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var count int
		return tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM openfga_store").Scan(&count)
	})
	return storage.ReadinessStatus{IsReady: err == nil}, err
}

func collision(err error) error {
	if internaldatabase.IsErrConstraintUnique(err) || internaldatabase.IsErrConstraintPrimaryKey(err) {
		return fmt.Errorf("%w: %w", storage.ErrCollision, err)
	}
	return err
}

func pageSize(p storage.PaginationOptions) int {
	if p.PageSize <= 0 {
		return storage.DefaultPageSize
	}
	return p.PageSize
}

func timestamp(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000000000Z")
}
