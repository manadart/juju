// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package database

import (
	"strings"

	"github.com/juju/collections/transform"
)

// SliceToPlaceholder returns a string that can be used in a SQL/DML
// statement as a parameter list for a [NOT] IN clause.
// For example, passing []int{1, 2, 3} would return "?,?,?".
// It also returns a suitable transformed slice of the input values to type any.
func SliceToPlaceholder[T any](in []T) (string, []any) {
	vals := make([]any, 0, len(in))
	return strings.Join(transform.Slice(in, func(item T) string {
		vals = append(vals, item)
		return "?"
	}), ","), vals
}

// SliceToPlaceholderTransform returns a string that can be used in SQL/DML
// statement as a parameter list for a [NOT] IN clause.
// For example, passing []int{1, 2, 3} would return "?,?,?".
// Also takes a transform function to alter the type and meaning of the in slice
// into a new slice that can be used with the parameters.
func SliceToPlaceholderTransform[T any](in []T, trans func(T) any) (string, []any) {
	vals := make([]any, 0, len(in))
	return strings.Join(transform.Slice(in, func(item T) string {
		vals = append(vals, trans(item))
		return "?"
	}), ","), vals
}