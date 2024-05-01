// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import (
	"fmt"
	"reflect"
)

var _ ValueComparer = valuesSameAny{}

type valuesSameAny struct{}

// CompareValues determines whether any value in the supplied values
// matches any other value.
func (v valuesSameAny) CompareValues(values ...any) error {
	if len(values) < 2 {
		return nil
	}

	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if reflect.DeepEqual(values[i], values[j]) {
				return nil
			}
		}
	}

	return fmt.Errorf("expected at least two values to be the same, but all values differ")
}

// ValuesSameAny returns a ValueComparer for asserting whether any value in the
// values supplied to the CompareValues method is the same as any other value.
func ValuesSameAny() valuesSameAny {
	return valuesSameAny{}
}
