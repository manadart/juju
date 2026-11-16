// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package internal

import "github.com/juju/juju/internal/errors"

const (
	ErrAuthorisationNotFound  = errors.ConstError("authorisation record not found")
	ErrAuthorisationCollision = errors.ConstError("authorisation record already exists")
)

// Tuple is a persisted relationship, including its ordering and condition.
type Tuple struct {
	ID, ObjectType, ObjectID, Relation, User, Condition, CreatedAt string
}

// TupleFilter restricts tuple retrieval. Empty fields impose no restriction.
type TupleFilter struct {
	ObjectType, ObjectID, Relation, User string
}

// TupleWrite describes one atomic, replayable set of tuple mutations.
type TupleWrite struct {
	OperationID                                 string
	Deletes, Writes                             []Tuple
	IgnoreMissingDeletes, IgnoreDuplicateWrites bool
}

// TupleWriteError identifies the rejected mutation without exposing SQL errors.
type TupleWriteError struct {
	Index  int
	Reason TupleWriteFailure
}

type TupleWriteFailure int

const (
	TupleDeleteMissing TupleWriteFailure = iota
	TupleInsertExists
	TupleConditionConflict
)

func (e *TupleWriteError) Error() string {
	return "invalid authorisation tuple mutation"
}

// TupleChange describes an actual insertion or deletion in the changelog.
type TupleChange struct {
	Tuple
	Operation int32
}

// ChangesFilter bounds a changelog page in either direction.
type ChangesFilter struct {
	ObjectType, From, Until string
	Descending              bool
	Limit                   int
}

// AuthorisationStore contains the metadata of an active store.
type AuthorisationStore struct {
	ID, Name, CreatedAt, UpdatedAt string
}
