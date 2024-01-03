// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = MapElements{}

// MapElements is a Check for asserting equality between the value supplied
// to MapElementsExact and the value passed to the CheckValue method.
type MapElements struct {
	num int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v MapElements) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for MapElements check, got: %T", other)
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

		return fmt.Errorf("expected %d %s for MapElements check, got %d %s", v.num, expectedElements, len(otherVal), actualElements)
	}

	return nil
}

// String returns the string representation of the value.
func (v MapElements) String() string {
	return strconv.FormatInt(int64(v.num), 10)
}

// MapElementsExact returns a Check for asserting that
// a list num elements.
func MapElementsExact(num int) MapElements {
	return MapElements{
		num: num,
	}
}
