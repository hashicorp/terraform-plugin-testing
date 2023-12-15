// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"strconv"
)

var _ KnownValue = Int64Value{}

type Int64Value struct {
	value int64
}

// Equal determines whether the passed value is of type int64, and
// contains a matching int64 value.
func (v Int64Value) Equal(other any) bool {
	otherVal, ok := other.(int64)

	if !ok {
		return false
	}

	if otherVal != v.value {
		return false
	}

	return true
}

// String returns the string representation of the int64 value.
func (v Int64Value) String() string {
	return strconv.FormatInt(v.value, 10)
}

// NewInt64Value returns a KnownValue for asserting equality between the
// supplied int64 and the value passed to the Equal method.
func NewInt64Value(value int64) Int64Value {
	return Int64Value{
		value: value,
	}
}
