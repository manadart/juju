// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package openfga

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (d *Datastore) CreateStore(ctx context.Context, store *openfgav1.Store) (*openfgav1.Store, error) {
	if store.GetId() == "" || store.GetName() == "" {
		return nil, fmt.Errorf("store ID and name are required")
	}
	now := time.Now().UTC()
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO openfga_store (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`, store.Id, store.Name, timestamp(now), timestamp(now))
		return err
	})
	if err != nil {
		return nil, collision(err)
	}
	return &openfgav1.Store{Id: store.Id, Name: store.Name, CreatedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}, nil
}

func (d *Datastore) DeleteStore(ctx context.Context, id string) error {
	return d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "UPDATE openfga_store SET deleted_at = ? WHERE id = ?", timestamp(time.Now()), id)
		return err
	})
}

func (d *Datastore) GetStore(ctx context.Context, id string) (*openfgav1.Store, error) {
	var name, created, updated string
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT name, created_at, updated_at FROM openfga_store WHERE id = ? AND deleted_at IS NULL`, id).Scan(&name, &created, &updated)
	})
	if err != nil {
		return nil, err
	}
	return decodeStore(id, name, created, updated)
}

func (d *Datastore) ListStores(ctx context.Context, options storage.ListStoresOptions) ([]*openfgav1.Store, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	var result []*openfgav1.Store
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		result = nil
		rows, err := tx.QueryContext(ctx, `SELECT id, name, created_at, updated_at FROM openfga_store WHERE deleted_at IS NULL AND id > ? ORDER BY id`, options.Pagination.From)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, name, created, updated string
			if err := rows.Scan(&id, &name, &created, &updated); err != nil {
				return err
			}
			if options.Name != "" && name != options.Name {
				continue
			}
			if len(options.IDs) > 0 && !slices.Contains(options.IDs, id) {
				continue
			}
			store, err := decodeStore(id, name, created, updated)
			if err != nil {
				return err
			}
			result = append(result, store)
			if len(result) > pageSize(options.Pagination) {
				break
			}
		}
		return rows.Err()
	})
	var next string
	if n := pageSize(options.Pagination); len(result) > n {
		result = result[:n]
		next = result[n-1].Id
	}
	return result, next, err
}

func decodeStore(id, name, created, updated string) (*openfgav1.Store, error) {
	c, err := time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return nil, err
	}
	u, err := time.Parse(time.RFC3339Nano, updated)
	if err != nil {
		return nil, err
	}
	return &openfgav1.Store{Id: id, Name: name, CreatedAt: timestamppb.New(c), UpdatedAt: timestamppb.New(u)}, nil
}
