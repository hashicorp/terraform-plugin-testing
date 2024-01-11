// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = setElementsExact{}

type setElementsExact struct {
	num int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v setElementsExact) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for SetElementExact check, got: %T", other)
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

		return fmt.Errorf("expected %d %s for SetElementExact check, got %d %s", v.num, expectedElements, len(otherVal), actualElements)
	}

	return nil
}

// String returns the string representation of the value.
func (v setElementsExact) String() string {
	return strconv.FormatInt(int64(v.num), 10)
}

// SetElementsExact returns a Check for asserting that
// a set has num elements.
func SetElementsExact(num int) setElementsExact {
	return setElementsExact{
		num: num,
	}
}
