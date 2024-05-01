// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import (
	"fmt"
	"reflect"
)

var _ ValueComparer = valuesDifferAny{}

type valuesDifferAny struct{}

// CompareValues determines whether any value in the supplied values
// is unique.
func (v valuesDifferAny) CompareValues(values ...any) error {
	if len(values) < 2 {
		return nil
	}

	var val any

	for i := 0; i < len(values); i++ {
		val = values[i]
		for j := 1; j < len(values); j++ {
			if !reflect.DeepEqual(values[i], values[j]) {
				return nil
			}
		}
	}

	return fmt.Errorf("expected values to differ, but value is duplicated: %v", val)
}

// ValuesDifferAny returns a ValueComparer for asserting that any value in the
// values supplied to the CompareValues method is unique.
func ValuesDifferAny() valuesDifferAny {
	return valuesDifferAny{}
}
