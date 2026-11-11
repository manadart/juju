// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package openfga

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/storage"
	"google.golang.org/protobuf/encoding/protojson"
)

func (d *Datastore) MaxTypesPerAuthorizationModel() int {
	return storage.DefaultMaxTypesPerAuthorizationModel
}

func (d *Datastore) WriteAuthorizationModel(ctx context.Context, store string, model *openfgav1.AuthorizationModel) error {
	data, err := protojson.Marshal(model)
	if err != nil {
		return err
	}
	return d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO openfga_authorization_model (store_id, model_id, model) VALUES (?, ?, ?) ON CONFLICT (store_id, model_id) DO UPDATE SET model = excluded.model`, store, model.Id, string(data))
		return err
	})
}

func (d *Datastore) ReadAuthorizationModel(ctx context.Context, store, id string) (*openfgav1.AuthorizationModel, error) {
	var data string
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT model FROM openfga_authorization_model WHERE store_id = ? AND model_id = ?`, store, id).Scan(&data)
	})
	if err != nil {
		return nil, err
	}
	model := &openfgav1.AuthorizationModel{}
	if err := protojson.Unmarshal([]byte(data), model); err != nil {
		return nil, err
	}
	if len(model.TypeDefinitions) == 0 {
		return nil, storage.ErrNotFound
	}
	return model, nil
}

func (d *Datastore) ReadAuthorizationModels(ctx context.Context, store string, options storage.ReadAuthorizationModelsOptions) ([]*openfgav1.AuthorizationModel, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	query := "SELECT model FROM openfga_authorization_model WHERE store_id = ?"
	args := []any{store}
	if options.Pagination.From != "" {
		query += " AND model_id < ?"
		args = append(args, options.Pagination.From)
	}
	query += " ORDER BY model_id DESC LIMIT ?"
	args = append(args, pageSize(options.Pagination)+1)
	var result []*openfgav1.AuthorizationModel
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		result = nil
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var data string
			if err := rows.Scan(&data); err != nil {
				return err
			}
			model := &openfgav1.AuthorizationModel{}
			if err := protojson.Unmarshal([]byte(data), model); err != nil {
				return err
			}
			result = append(result, model)
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

func (d *Datastore) FindLatestAuthorizationModel(ctx context.Context, store string) (*openfgav1.AuthorizationModel, error) {
	models, _, err := d.ReadAuthorizationModels(ctx, store, storage.ReadAuthorizationModelsOptions{Pagination: storage.NewPaginationOptions(1, "")})
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, storage.ErrNotFound
	}
	return models[0], nil
}

func (d *Datastore) WriteAssertions(ctx context.Context, store, model string, assertions []*openfgav1.Assertion) error {
	data := make([]json.RawMessage, 0, len(assertions))
	for _, assertion := range assertions {
		b, err := protojson.Marshal(assertion)
		if err != nil {
			return err
		}
		data = append(data, b)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO openfga_assertion (store_id, model_id, assertions) VALUES (?, ?, ?) ON CONFLICT (store_id, model_id) DO UPDATE SET assertions = excluded.assertions`, store, model, string(b))
		return err
	})
}

func (d *Datastore) ReadAssertions(ctx context.Context, store, model string) ([]*openfgav1.Assertion, error) {
	var data string
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT assertions FROM openfga_assertion WHERE store_id = ? AND model_id = ?`, store, model).Scan(&data)
	})
	if errors.Is(err, storage.ErrNotFound) {
		return []*openfgav1.Assertion{}, nil
	}
	if err != nil {
		return nil, err
	}
	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}
	result := make([]*openfgav1.Assertion, 0, len(raw))
	for _, data := range raw {
		assertion := &openfgav1.Assertion{}
		if err := protojson.Unmarshal(data, assertion); err != nil {
			return nil, err
		}
		result = append(result, assertion)
	}
	return result, nil
}
