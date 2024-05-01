// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import "fmt"

var _ ValueComparer = valuesDifferAll{}

type valuesDifferAll struct{}

// CompareValues determines whether each value in the supplied values
// is unique.
func (v valuesDifferAll) CompareValues(values ...any) error {
	vals := map[any]struct{}{}

	for i := 0; i < len(values); i++ {
		if _, ok := vals[values[i]]; ok {
			return fmt.Errorf("expected values to differ, but value is duplicated: %v", values[i])
		}

		vals[values[i]] = struct{}{}
	}

	return nil
}

// ValuesDifferAll returns a ValueComparer for asserting that each value in the
// values supplied to the CompareValues method is unique.
func ValuesDifferAll() valuesDifferAll {
	return valuesDifferAll{}
}
