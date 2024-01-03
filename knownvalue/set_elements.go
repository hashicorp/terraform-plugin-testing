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
	val, ok := other.([]any)

	if !ok {
		return fmt.Errorf("wrong type: %T, expected []any", other)
	}

	if len(val) != v.num {
		return fmt.Errorf("wrong length: %d, expected %d", len(val), v.num)
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
