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
	val, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("wrong type: %T, expected map[string]any", other)
	}

	if len(val) != v.num {
		return fmt.Errorf("wrong length: %d, expected %d", len(val), v.num)
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
