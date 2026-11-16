// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/juju/juju/domain/access/internal"
	"github.com/juju/juju/internal/errors"
)

// GetAuthorisationTuples materialises matching rows inside the transaction.
func (s *AuthorisationState) GetAuthorisationTuples(ctx context.Context, store string, filter internal.TupleFilter) ([]internal.Tuple, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	in := authTuple{Store: store, ObjectType: filter.ObjectType, ObjectID: filter.ObjectID, Relation: filter.Relation, User: filter.User}
	query := `
SELECT (t.ulid, t.object_type, t.object_id, t.relation, t.user, t.condition, t.inserted_at) AS (&authTuple.*)
FROM openfga_tuple AS t WHERE t.store_id = $authTuple.store_id
`
	if filter.ObjectType != "" {
		query += " AND t.object_type = $authTuple.object_type"
	}
	if filter.ObjectID != "" {
		query += " AND t.object_id = $authTuple.object_id"
	}
	if filter.Relation != "" {
		query += " AND t.relation = $authTuple.relation"
	}
	if filter.User != "" {
		query += " AND t.user = $authTuple.user"
	}
	query += " ORDER BY t.ulid"
	stmt, err := s.Prepare(query, in)
	if err != nil {
		return nil, errors.Capture(err)
	}
	var rows []authTuple
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, in).GetAll(&rows)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return err
	})
	if err != nil {
		return nil, errors.Capture(err)
	}
	result := make([]internal.Tuple, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.toTuple())
	}
	return result, nil
}

// WriteAuthorisationTuples validates and applies all mutations atomically.
// Conditions on duplicate tuples must be compared within the same transaction.
func (s *AuthorisationState) WriteAuthorisationTuples(ctx context.Context, store string, write internal.TupleWrite) error {
	deleteStmt, err := s.Prepare(`
DELETE FROM openfga_tuple AS t
WHERE t.store_id = $authTuple.store_id AND t.object_type = $authTuple.object_type
AND t.object_id = $authTuple.object_id AND t.relation = $authTuple.relation
AND t.user = $authTuple.user
`, authTuple{})
	if err != nil {
		return errors.Capture(err)
	}
	existingStmt, err := s.Prepare(`
SELECT t.condition AS &authTuple.condition FROM openfga_tuple AS t
WHERE t.store_id = $authTuple.store_id AND t.object_type = $authTuple.object_type
AND t.object_id = $authTuple.object_id AND t.relation = $authTuple.relation
AND t.user = $authTuple.user
`, authTuple{})
	if err != nil {
		return errors.Capture(err)
	}
	insertStmt, err := s.Prepare(`
INSERT INTO openfga_tuple (store_id, ulid, object_type, object_id, relation, user, condition, inserted_at)
VALUES ($authTuple.store_id, $authTuple.ulid, $authTuple.object_type, $authTuple.object_id,
        $authTuple.relation, $authTuple.user, $authTuple.condition, $authTuple.inserted_at)
`, authTuple{})
	if err != nil {
		return errors.Capture(err)
	}
	return s.write(ctx, write.OperationID, func(ctx context.Context, tx *sqlair.TX) error {
		for i, tuple := range write.Deletes {
			var outcome sqlair.Outcome
			if err := tx.Query(ctx, deleteStmt, newAuthTuple(store, tuple)).Get(&outcome); err != nil {
				return err
			}
			n, err := outcome.Result().RowsAffected()
			if err != nil {
				return err
			}
			if n == 0 && !write.IgnoreMissingDeletes {
				return &internal.TupleWriteError{Index: i, Reason: internal.TupleDeleteMissing}
			}
		}
		for i, tuple := range write.Writes {
			in := newAuthTuple(store, tuple)
			var existing authTuple
			err := tx.Query(ctx, existingStmt, in).Get(&existing)
			if err == nil {
				if !write.IgnoreDuplicateWrites {
					return &internal.TupleWriteError{Index: i, Reason: internal.TupleInsertExists}
				}
				equal, err := equalTupleConditions(existing.Condition, in.Condition)
				if err != nil {
					return err
				}
				if !equal {
					return &internal.TupleWriteError{Index: i, Reason: internal.TupleConditionConflict}
				}
				continue
			}
			if !errors.Is(err, sqlair.ErrNoRows) {
				return err
			}
			if err := tx.Query(ctx, insertStmt, in).Run(); err != nil {
				return err
			}
		}
		return nil
	})
}

func equalTupleConditions(a, b string) (bool, error) {
	var left, right *openfgav1.RelationshipCondition
	if a != "" {
		left = &openfgav1.RelationshipCondition{}
		if err := protojson.Unmarshal([]byte(a), left); err != nil {
			return false, err
		}
	}
	if b != "" {
		right = &openfgav1.RelationshipCondition{}
		if err := protojson.Unmarshal([]byte(b), right); err != nil {
			return false, err
		}
	}
	return proto.Equal(left, right), nil
}

// GetAuthorisationChanges retrieves a bounded page of committed tuple changes.
func (s *AuthorisationState) GetAuthorisationChanges(ctx context.Context, store string, filter internal.ChangesFilter) ([]internal.TupleChange, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	in := authChangesFilter{Store: store, ObjectType: filter.ObjectType, From: filter.From, Until: filter.Until, Limit: filter.Limit}
	query := `
SELECT (c.ulid, c.object_type, c.object_id, c.relation, c.user, c.condition, c.inserted_at) AS (&authTuple.*),
       c.operation AS &authChangeOperation.operation
FROM openfga_changelog AS c
WHERE c.store_id = $authChangesFilter.store_id AND c.inserted_at <= $authChangesFilter.until
`
	if filter.ObjectType != "" {
		query += " AND c.object_type = $authChangesFilter.object_type"
	}
	if filter.From != "" {
		if filter.Descending {
			query += " AND c.ulid < $authChangesFilter.from"
		} else {
			query += " AND c.ulid > $authChangesFilter.from"
		}
	}
	if filter.Descending {
		query += " ORDER BY c.ulid DESC"
	} else {
		query += " ORDER BY c.ulid ASC"
	}
	query += " LIMIT $authChangesFilter.limit"
	stmt, err := s.Prepare(query, in, authTuple{}, authChangeOperation{})
	if err != nil {
		return nil, errors.Capture(err)
	}
	var rows []authTuple
	var operations []authChangeOperation
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, in).GetAll(&rows, &operations)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return err
	})
	if err != nil {
		return nil, errors.Capture(err)
	}
	result := make([]internal.TupleChange, 0, len(rows))
	for i, row := range rows {
		result = append(result, internal.TupleChange{Tuple: row.toTuple(), Operation: operations[i].Operation})
	}
	return result, nil
}

type authTuple struct {
	Store      string `db:"store_id"`
	ID         string `db:"ulid"`
	ObjectType string `db:"object_type"`
	ObjectID   string `db:"object_id"`
	Relation   string `db:"relation"`
	User       string `db:"user"`
	Condition  string `db:"condition"`
	CreatedAt  string `db:"inserted_at"`
}

func newAuthTuple(store string, t internal.Tuple) authTuple {
	return authTuple{Store: store, ID: t.ID, ObjectType: t.ObjectType, ObjectID: t.ObjectID, Relation: t.Relation, User: t.User, Condition: t.Condition, CreatedAt: t.CreatedAt}
}
func (t authTuple) toTuple() internal.Tuple {
	return internal.Tuple{ID: t.ID, ObjectType: t.ObjectType, ObjectID: t.ObjectID, Relation: t.Relation, User: t.User, Condition: t.Condition, CreatedAt: t.CreatedAt}
}

type authChangeOperation struct {
	Operation int32 `db:"operation"`
}
type authChangesFilter struct {
	Store      string `db:"store_id"`
	ObjectType string `db:"object_type"`
	From       string `db:"from"`
	Until      string `db:"until"`
	Limit      int    `db:"limit"`
}
