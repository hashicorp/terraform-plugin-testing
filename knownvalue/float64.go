// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"strconv"
)

var _ KnownValue = Float64Value{}

// Float64Value is a KnownValue for asserting equality between the value
// supplied to Float64ValueExact and the value passed to the Equal method.
type Float64Value struct {
	value float64
}

// Equal determines whether the passed value is of type float64, and
// contains a matching float64 value.
func (v Float64Value) Equal(other any) bool {
	otherVal, ok := other.(float64)

	if !ok {
		return false
	}

	if otherVal != v.value {
		return false
	}

	return true
}

// String returns the string representation of the float64 value.
func (v Float64Value) String() string {
	return strconv.FormatFloat(v.value, 'f', -1, 64)
}

// Float64ValueExact returns a KnownValue for asserting equality between the
// supplied float64 and the value passed to the Equal method.
func Float64ValueExact(value float64) Float64Value {
	return Float64Value{
		value: value,
	}
}
