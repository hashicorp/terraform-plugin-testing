// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import (
	"fmt"
	"reflect"
)

var _ ValueComparer = valuesDifferAll{}

type valuesDifferAll struct{}

// CompareValues determines whether each value in the supplied values
// is unique.
func (v valuesDifferAll) CompareValues(values ...any) error {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if reflect.DeepEqual(values[i], values[j]) {
				return fmt.Errorf("expected values to differ, but value is duplicated: %v", values[i])
			}
		}
	}

	return nil
}

// ValuesDifferAll returns a ValueComparer for asserting that each value in the
// values supplied to the CompareValues method is unique.
func ValuesDifferAll() valuesDifferAll {
	return valuesDifferAll{}
}
