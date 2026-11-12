// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import "strings"

// DirectiveUpdate preserves which fields an allocation-default edit supplies.
// Nil size and count, and an empty pool, retain the existing directive values.
type DirectiveUpdate struct {
	Pool    string
	SizeMiB *uint64
	Count   *uint64
}

// ParseDirectiveUpdate parses the usual directive syntax without defaulting
// omitted fields. In particular, a size-only edit does not reset count to one.
func ParseDirectiveUpdate(s string) (DirectiveUpdate, error) {
	directive, err := ParseDirective(s)
	if err != nil {
		return DirectiveUpdate{}, err
	}
	update := DirectiveUpdate{Pool: directive.Pool}
	if directive.Size != 0 {
		update.SizeMiB = &directive.Size
	}
	for field := range strings.SplitSeq(s, ",") {
		if _, ok, _ := parseCount(field); ok {
			update.Count = &directive.Count
			break
		}
	}
	return update, nil
}
