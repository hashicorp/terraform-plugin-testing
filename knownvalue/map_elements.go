// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = mapElementsExact{}

type mapElementsExact struct {
	num int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v mapElementsExact) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for MapElementsExact check, got: %T", other)
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

		return fmt.Errorf("expected %d %s for MapElementsExact check, got %d %s", v.num, expectedElements, len(otherVal), actualElements)
	}

	return nil
}

// String returns the string representation of the value.
func (v mapElementsExact) String() string {
	return strconv.Itoa(v.num)
}

// MapElementsExact returns a Check for asserting that
// a map has num elements.
func MapElementsExact(num int) mapElementsExact {
	return mapElementsExact{
		num: num,
	}
}
