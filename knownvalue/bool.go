// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "strconv"

var _ KnownValue = BoolValue{}

// BoolValue is a KnownValue for asserting equality between the value
// supplied to NewBoolValue and the value passed to the Equal method.
type BoolValue struct {
	value bool
}

// Equal determines whether the passed value is of type bool, and
// contains a matching bool value.
func (v BoolValue) Equal(other any) bool {
	otherVal, ok := other.(bool)

	if !ok {
		return false
	}

	if otherVal != v.value {
		return false
	}

	return true
}

// String returns the string representation of the bool value.
func (v BoolValue) String() string {
	return strconv.FormatBool(v.value)
}

// NewBoolValue returns a KnownValue for asserting equality between the
// supplied bool and the value passed to the Equal method.
func NewBoolValue(value bool) BoolValue {
	return BoolValue{
		value: value,
	}
}
