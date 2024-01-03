// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = BoolValue{}

// BoolValue is a KnownValue for asserting equality between the value supplied
// to BoolValueExact and the value passed to the CheckValue method.
type BoolValue struct {
	value bool
}

// CheckValue determines whether the passed value is of type bool, and
// contains a matching bool value.
func (v BoolValue) CheckValue(other any) error {
	otherVal, ok := other.(bool)

	if !ok {
		return fmt.Errorf("wrong type: %T, known value type is bool", other)
	}

	if otherVal != v.value {
		return fmt.Errorf("value: %t does not equal expected value: %t", otherVal, v.value)
	}

	return nil
}

// String returns the string representation of the bool value.
func (v BoolValue) String() string {
	return strconv.FormatBool(v.value)
}

// BoolValueExact returns a Check for asserting equality between the
// supplied bool and the value passed to the CheckValue method.
func BoolValueExact(value bool) BoolValue {
	return BoolValue{
		value: value,
	}
}
