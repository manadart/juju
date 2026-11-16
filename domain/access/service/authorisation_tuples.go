// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/tuple"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/juju/juju/domain/access/internal"
	"github.com/juju/juju/internal/errors"
)

type tupleRow struct {
	id    string
	tuple *openfgav1.Tuple
}

// Read materialises its results before the transaction can be retried.
func (s *AuthorisationService) Read(ctx context.Context, store string, filter storage.ReadFilter, _ storage.ReadOptions) (storage.TupleIterator, error) {
	rows, err := s.read(ctx, store, filter)
	if err != nil {
		return nil, err
	}
	return storage.NewStaticTupleIterator(tuples(rows)), nil
}

func (s *AuthorisationService) ReadPage(ctx context.Context, store string, filter storage.ReadFilter, options storage.ReadPageOptions) ([]*openfgav1.Tuple, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	rows, err := s.read(ctx, store, filter)
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

func (s *AuthorisationService) ReadUserTuple(ctx context.Context, store string, filter storage.ReadUserTupleFilter, _ storage.ReadUserTupleOptions) (*openfgav1.Tuple, error) {
	rows, err := s.read(ctx, store, filter)
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

func (s *AuthorisationService) ReadUsersetTuples(ctx context.Context, store string, filter storage.ReadUsersetTuplesFilter, _ storage.ReadUsersetTuplesOptions) (storage.TupleIterator, error) {
	rows, err := s.read(ctx, store, storage.ReadFilter{Object: filter.Object, Relation: filter.Relation, Conditions: filter.Conditions})
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

func (s *AuthorisationService) ReadStartingWithUser(ctx context.Context, store string, filter storage.ReadStartingWithUserFilter, options storage.ReadStartingWithUserOptions) (storage.TupleIterator, error) {
	rows, err := s.read(ctx, store, storage.ReadFilter{Object: filter.ObjectType + ":", Relation: filter.Relation, Conditions: filter.Conditions})
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

func (s *AuthorisationService) read(ctx context.Context, store string, filter storage.ReadFilter) ([]tupleRow, error) {
	typ, id := tuple.SplitObject(filter.Object)
	rows, err := s.st.GetAuthorisationTuples(ctx, store, internal.TupleFilter{
		ObjectType: typ, ObjectID: id, Relation: filter.Relation, User: filter.User,
	})
	if err != nil {
		return nil, authorisationStorageError(err)
	}
	result := make([]tupleRow, 0, len(rows))
	for _, row := range rows {
		t, err := decodeTuple(row)
		if err != nil {
			return nil, err
		}
		if len(filter.Conditions) > 0 && !slices.Contains(filter.Conditions, t.Key.GetCondition().GetName()) {
			continue
		}
		result = append(result, tupleRow{id: row.ID, tuple: t})
	}
	return result, nil
}

func decodeTuple(row internal.Tuple) (*openfgav1.Tuple, error) {
	key := &openfgav1.TupleKey{Object: row.ObjectType + ":" + row.ObjectID, Relation: row.Relation, User: row.User}
	if row.Condition != "" {
		key.Condition = &openfgav1.RelationshipCondition{}
		if err := protojson.Unmarshal([]byte(row.Condition), key.Condition); err != nil {
			return nil, err
		}
	}
	created, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &openfgav1.Tuple{Key: key, Timestamp: timestamppb.New(created)}, nil
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

func (s *AuthorisationService) MaxTuplesPerWrite() int { return storage.DefaultMaxTuplesPerWrite }

// Write records a replay marker in the same transaction as tuples and changes.
// An ambiguous commit can therefore be retried without duplicating the write.
func (s *AuthorisationService) Write(ctx context.Context, store string, deletes storage.Deletes, writes storage.Writes, options ...storage.TupleWriteOption) error {
	opts := storage.NewTupleWriteOptions(options...)
	write := internal.TupleWrite{
		OperationID:           ulid.Make().String(),
		IgnoreMissingDeletes:  opts.OnMissingDelete == storage.OnMissingDeleteIgnore,
		IgnoreDuplicateWrites: opts.OnDuplicateInsert == storage.OnDuplicateInsertIgnore,
	}
	for _, key := range deletes {
		typ, id := tuple.SplitObject(key.Object)
		write.Deletes = append(write.Deletes, internal.Tuple{ObjectType: typ, ObjectID: id, Relation: key.Relation, User: key.User})
	}
	now := timestamp(s.clock.Now())
	for _, key := range writes {
		typ, id := tuple.SplitObject(key.Object)
		row := internal.Tuple{ID: ulid.Make().String(), ObjectType: typ, ObjectID: id, Relation: key.Relation, User: key.User, CreatedAt: now}
		if key.Condition != nil {
			data, err := protojson.Marshal(key.Condition)
			if err != nil {
				return err
			}
			row.Condition = string(data)
		}
		write.Writes = append(write.Writes, row)
	}
	err := s.st.WriteAuthorisationTuples(ctx, store, write)
	var invalid *internal.TupleWriteError
	if errors.As(err, &invalid) {
		switch invalid.Reason {
		case internal.TupleDeleteMissing:
			return storage.InvalidWriteInputError(deletes[invalid.Index], openfgav1.TupleOperation_TUPLE_OPERATION_DELETE)
		case internal.TupleInsertExists:
			return storage.InvalidWriteInputError(writes[invalid.Index], openfgav1.TupleOperation_TUPLE_OPERATION_WRITE)
		case internal.TupleConditionConflict:
			return storage.TupleConditionConflictError(writes[invalid.Index])
		}
	}
	if errors.Is(err, internal.ErrAuthorisationCollision) {
		return errors.Errorf("%w: %w", storage.ErrTransactionalWriteFailed, err)
	}
	return err
}

func (s *AuthorisationService) ReadChanges(ctx context.Context, store string, filter storage.ReadChangesFilter, options storage.ReadChangesOptions) ([]*openfgav1.TupleChange, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	rows, err := s.st.GetAuthorisationChanges(ctx, store, internal.ChangesFilter{
		ObjectType: filter.ObjectType, From: options.Pagination.From,
		Until:      timestamp(s.clock.Now().Add(-filter.HorizonOffset)),
		Descending: options.SortDesc, Limit: pageSize(options.Pagination),
	})
	if err != nil {
		return nil, "", authorisationStorageError(err)
	}
	if len(rows) == 0 {
		return nil, "", storage.ErrNotFound
	}
	result := make([]*openfgav1.TupleChange, 0, len(rows))
	for _, row := range rows {
		t, err := decodeTuple(row.Tuple)
		if err != nil {
			return nil, "", err
		}
		result = append(result, &openfgav1.TupleChange{TupleKey: t.Key, Operation: openfgav1.TupleOperation(row.Operation), Timestamp: t.Timestamp})
	}
	return result, rows[len(rows)-1].ID, nil
}
