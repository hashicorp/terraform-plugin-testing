// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = ObjectAttributes{}

// ObjectAttributes is a Check for asserting equality between the value supplied
// to ObjectAttributesExact and the value passed to the CheckValue method.
type ObjectAttributes struct {
	num int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v ObjectAttributes) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for ObjectAttributes check, got: %T", other)
	}

	if len(otherVal) != v.num {
		expectedAttributes := "attributes"
		actualAttributes := "attributes"

		if v.num == 1 {
			expectedAttributes = "attribute"
		}

		if len(otherVal) == 1 {
			actualAttributes = "attribute"
		}

		return fmt.Errorf("expected %d %s for ObjectAttributes check, got %d %s", v.num, expectedAttributes, len(otherVal), actualAttributes)
	}

	return nil
}

// String returns the string representation of the value.
func (v ObjectAttributes) String() string {
	return strconv.FormatInt(int64(v.num), 10)
}

// ObjectAttributesExact returns a Check for asserting that
// a list num elements.
func ObjectAttributesExact(num int) ObjectAttributes {
	return ObjectAttributes{
		num: num,
	}
}
