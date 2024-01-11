// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = objectAttributesExact{}

type objectAttributesExact struct {
	num int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v objectAttributesExact) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for ObjectAttributesExact check, got: %T", other)
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

		return fmt.Errorf("expected %d %s for ObjectAttributesExact check, got %d %s", v.num, expectedAttributes, len(otherVal), actualAttributes)
	}

	return nil
}

// String returns the string representation of the value.
func (v objectAttributesExact) String() string {
	return strconv.FormatInt(int64(v.num), 10)
}

// ObjectAttributesExact returns a Check for asserting that
// an object has num attributes.
func ObjectAttributesExact(num int) objectAttributesExact {
	return objectAttributesExact{
		num: num,
	}
}
