// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = Int64Value{}

// Int64Value is a KnownValue for asserting equality between the value supplied
// to Int64ValueExact and the value passed to the CheckValue method.
type Int64Value struct {
	value int64
}

// CheckValue determines whether the passed value is of type int64, and
// contains a matching int64 value.
func (v Int64Value) CheckValue(other any) error {
	otherVal, ok := other.(int64)

	if !ok {
		return fmt.Errorf("expected int64 value for Int64Value check, got: %T", other)
	}

	if otherVal != v.value {
		return fmt.Errorf("expected value %d for Int64Value check, got: %d", v.value, otherVal)
	}

	return nil
}

// String returns the string representation of the int64 value.
func (v Int64Value) String() string {
	return strconv.FormatInt(v.value, 10)
}

// Int64ValueExact returns a Check for asserting equality between the
// supplied int64 and the value passed to the CheckValue method.
func Int64ValueExact(value int64) Int64Value {
	return Int64Value{
		value: value,
	}
}
