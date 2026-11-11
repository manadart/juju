// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package openfga

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/tuple"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	internaldatabase "github.com/juju/juju/internal/database"
)

type tupleRow struct {
	id    string
	tuple *openfgav1.Tuple
}

// Read materialises its results before the transaction can be retried.
func (d *Datastore) Read(ctx context.Context, store string, filter storage.ReadFilter, _ storage.ReadOptions) (storage.TupleIterator, error) {
	rows, err := d.read(ctx, store, filter)
	if err != nil {
		return nil, err
	}
	return storage.NewStaticTupleIterator(tuples(rows)), nil
}

func (d *Datastore) ReadPage(ctx context.Context, store string, filter storage.ReadFilter, options storage.ReadPageOptions) ([]*openfgav1.Tuple, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	rows, err := d.read(ctx, store, filter)
	if err != nil {
		return nil, "", err
	}
	rows = slices.DeleteFunc(rows, func(r tupleRow) bool { return r.id <= options.Pagination.From })
	var next string
	if n := pageSize(options.Pagination); len(rows) > n {
		rows = rows[:n]
		next = rows[n-1].id
	}
	return tuples(rows), next, nil
}

func (d *Datastore) ReadUserTuple(ctx context.Context, store string, filter storage.ReadUserTupleFilter, _ storage.ReadUserTupleOptions) (*openfgav1.Tuple, error) {
	rows, err := d.read(ctx, store, filter)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		key := row.tuple.Key
		if key.Object == filter.Object && key.Relation == filter.Relation && key.User == filter.User {
			return row.tuple, nil
		}
	}
	return nil, storage.ErrNotFound
}

func (d *Datastore) ReadUsersetTuples(ctx context.Context, store string, filter storage.ReadUsersetTuplesFilter, _ storage.ReadUsersetTuplesOptions) (storage.TupleIterator, error) {
	rows, err := d.read(ctx, store, storage.ReadFilter{Object: filter.Object, Relation: filter.Relation, Conditions: filter.Conditions})
	if err != nil {
		return nil, err
	}
	var result []*openfgav1.Tuple
	for _, row := range rows {
		user := row.tuple.Key.User
		userObject, relation, userset := strings.Cut(user, "#")
		typ, id := tuple.SplitObject(userObject)
		if !userset && id != "*" {
			continue
		}
		allowed := len(filter.AllowedUserTypeRestrictions) == 0
		for _, restriction := range filter.AllowedUserTypeRestrictions {
			if restriction.GetType() == typ && ((userset && restriction.GetRelation() == relation) || (id == "*" && restriction.GetWildcard() != nil)) {
				allowed = true
			}
		}
		if allowed {
			result = append(result, row.tuple)
		}
	}
	return storage.NewStaticTupleIterator(result), nil
}

func (d *Datastore) ReadStartingWithUser(ctx context.Context, store string, filter storage.ReadStartingWithUserFilter, options storage.ReadStartingWithUserOptions) (storage.TupleIterator, error) {
	rows, err := d.read(ctx, store, storage.ReadFilter{Object: filter.ObjectType + ":", Relation: filter.Relation, Conditions: filter.Conditions})
	if err != nil {
		return nil, err
	}
	var result []*openfgav1.Tuple
	for _, row := range rows {
		_, id := tuple.SplitObject(row.tuple.Key.Object)
		if filter.ObjectIDs != nil && !filter.ObjectIDs.Exists(id) {
			continue
		}
		for _, user := range filter.UserFilter {
			want := user.Object
			if user.Relation != "" {
				want += "#" + user.Relation
			}
			if row.tuple.Key.User == want {
				result = append(result, row.tuple)
				break
			}
		}
	}
	if options.WithResultsSortedAscending {
		slices.SortFunc(result, func(a, b *openfgav1.Tuple) int { return strings.Compare(a.Key.Object, b.Key.Object) })
	}
	return storage.NewStaticTupleIterator(result), nil
}

func (d *Datastore) read(ctx context.Context, store string, filter storage.ReadFilter) ([]tupleRow, error) {
	query := `SELECT ulid, object_type, object_id, relation, user, condition, inserted_at FROM openfga_tuple WHERE store_id = ?`
	args := []any{store}
	if filter.Object != "" {
		typ, id := tuple.SplitObject(filter.Object)
		query += " AND object_type = ?"
		args = append(args, typ)
		if id != "" {
			query += " AND object_id = ?"
			args = append(args, id)
		}
	}
	if filter.Relation != "" {
		query += " AND relation = ?"
		args = append(args, filter.Relation)
	}
	if filter.User != "" {
		query += " AND user = ?"
		args = append(args, filter.User)
	}
	query += " ORDER BY ulid"
	var result []tupleRow
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		result = nil
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, typ, object, relation, user, condition, created string
			if err := rows.Scan(&id, &typ, &object, &relation, &user, &condition, &created); err != nil {
				return err
			}
			key := &openfgav1.TupleKey{Object: typ + ":" + object, Relation: relation, User: user}
			if condition != "" {
				key.Condition = &openfgav1.RelationshipCondition{}
				if err := protojson.Unmarshal([]byte(condition), key.Condition); err != nil {
					return err
				}
			}
			if len(filter.Conditions) > 0 && !slices.Contains(filter.Conditions, key.GetCondition().GetName()) {
				continue
			}
			t, err := time.Parse(time.RFC3339Nano, created)
			if err != nil {
				return err
			}
			result = append(result, tupleRow{id: id, tuple: &openfgav1.Tuple{Key: key, Timestamp: timestamppb.New(t)}})
		}
		return rows.Err()
	})
	return result, err
}

func tuples(rows []tupleRow) []*openfgav1.Tuple {
	result := make([]*openfgav1.Tuple, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.tuple)
	}
	return result
}

func validToken(token string) error {
	if token != "" {
		if _, err := ulid.ParseStrict(token); err != nil {
			return storage.ErrInvalidContinuationToken
		}
	}
	return nil
}

func (d *Datastore) MaxTuplesPerWrite() int { return storage.DefaultMaxTuplesPerWrite }

// Write records a replay marker in the same transaction as tuples and changes.
// An ambiguous commit can therefore be retried without duplicating the write.
func (d *Datastore) Write(ctx context.Context, store string, deletes storage.Deletes, writes storage.Writes, options ...storage.TupleWriteOption) error {
	opts := storage.NewTupleWriteOptions(options...)
	operation := ulid.Make().String()
	now := timestamp(time.Now())
	ids := make([]string, len(deletes)+len(writes))
	for i := range ids {
		ids[i] = ulid.Make().String()
	}
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var found string
		err := tx.QueryRowContext(ctx, "SELECT ulid FROM openfga_operation WHERE ulid = ?", operation).Scan(&found)
		if err == nil {
			return nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		for _, key := range deletes {
			typ, id := tuple.SplitObject(key.Object)
			res, err := tx.ExecContext(ctx, `DELETE FROM openfga_tuple WHERE store_id = ? AND object_type = ? AND object_id = ? AND relation = ? AND user = ?`, store, typ, id, key.Relation, key.User)
			if err != nil {
				return err
			}
			n, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if n == 0 {
				if opts.OnMissingDelete == storage.OnMissingDeleteIgnore {
					continue
				}
				return storage.InvalidWriteInputError(key, openfgav1.TupleOperation_TUPLE_OPERATION_DELETE)
			}

		}
		for i, key := range writes {
			typ, id := tuple.SplitObject(key.Object)
			var existing string
			err := tx.QueryRowContext(ctx, `SELECT condition FROM openfga_tuple WHERE store_id = ? AND object_type = ? AND object_id = ? AND relation = ? AND user = ?`, store, typ, id, key.Relation, key.User).Scan(&existing)
			if err == nil {
				if opts.OnDuplicateInsert != storage.OnDuplicateInsertIgnore {
					return storage.InvalidWriteInputError(key, openfgav1.TupleOperation_TUPLE_OPERATION_WRITE)
				}
				var old *openfgav1.RelationshipCondition
				if existing != "" {
					old = &openfgav1.RelationshipCondition{}
					if err := protojson.Unmarshal([]byte(existing), old); err != nil {
						return err
					}
				}
				if !proto.Equal(old, key.Condition) {
					return storage.TupleConditionConflictError(key)
				}
				continue
			}
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			condition := ""
			if key.Condition != nil {
				b, err := protojson.Marshal(key.Condition)
				if err != nil {
					return err
				}
				condition = string(b)
			}
			uid := ids[len(deletes)+i]
			if _, err := tx.ExecContext(ctx, `INSERT INTO openfga_tuple (store_id, ulid, object_type, object_id, relation, user, condition, inserted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, store, uid, typ, id, key.Relation, key.User, condition, now); err != nil {
				return err
			}

		}
		_, err = tx.ExecContext(ctx, "INSERT INTO openfga_operation (ulid) VALUES (?)", operation)
		return err
	})
	if internaldatabase.IsErrConstraintUnique(err) || internaldatabase.IsErrConstraintPrimaryKey(err) {
		return fmt.Errorf("%w: %w", storage.ErrTransactionalWriteFailed, err)
	}
	return err
}

func (d *Datastore) ReadChanges(ctx context.Context, store string, filter storage.ReadChangesFilter, options storage.ReadChangesOptions) ([]*openfgav1.TupleChange, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	query := `SELECT ulid, object_type, object_id, relation, user, condition, operation, inserted_at FROM openfga_changelog WHERE store_id = ? AND inserted_at <= ?`
	args := []any{store, timestamp(time.Now().Add(-filter.HorizonOffset))}
	if filter.ObjectType != "" {
		query += " AND object_type = ?"
		args = append(args, filter.ObjectType)
	}
	if options.Pagination.From != "" {
		if options.SortDesc {
			query += " AND ulid < ?"
		} else {
			query += " AND ulid > ?"
		}
		args = append(args, options.Pagination.From)
	}
	if options.SortDesc {
		query += " ORDER BY ulid DESC"
	} else {
		query += " ORDER BY ulid ASC"
	}
	query += " LIMIT ?"
	args = append(args, pageSize(options.Pagination))
	var result []*openfgav1.TupleChange
	var next string
	err := d.transact(ctx, func(ctx context.Context, tx *sql.Tx) error {
		result, next = nil, ""
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var typ, id, relation, user, condition, created string
			var operation int32
			if err := rows.Scan(&next, &typ, &id, &relation, &user, &condition, &operation, &created); err != nil {
				return err
			}
			key := &openfgav1.TupleKey{Object: typ + ":" + id, Relation: relation, User: user}
			if condition != "" {
				key.Condition = &openfgav1.RelationshipCondition{}
				if err := protojson.Unmarshal([]byte(condition), key.Condition); err != nil {
					return err
				}
			}
			t, err := time.Parse(time.RFC3339Nano, created)
			if err != nil {
				return err
			}
			result = append(result, &openfgav1.TupleChange{TupleKey: key, Operation: openfgav1.TupleOperation(operation), Timestamp: timestamppb.New(t)})
		}
		return rows.Err()
	})
	if err == nil && len(result) == 0 {
		err = storage.ErrNotFound
	}
	return result, next, err
}
