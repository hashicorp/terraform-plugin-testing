// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = SetElements{}

// SetElements is a Check for asserting equality between the value supplied
// to SetElementsExact and the value passed to the CheckValue method.
type SetElements struct {
	num int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v SetElements) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for SetElements check, got: %T", other)
	}

	if len(otherVal) != v.num {
		expectedElements := "elements"
		actualElements := "elements"

		if v.num == 1 {
			expectedElements = "element"
		}

		if len(otherVal) == 1 {
			actualElements = "element"
		}

		return fmt.Errorf("expected %d %s for SetElements check, got %d %s", v.num, expectedElements, len(otherVal), actualElements)
	}

	return nil
}

// String returns the string representation of the value.
func (v SetElements) String() string {
	return strconv.FormatInt(int64(v.num), 10)
}

// SetElementsExact returns a Check for asserting that
// a list num elements.
func SetElementsExact(num int) SetElements {
	return SetElements{
		num: num,
	}
}
