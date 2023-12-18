// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

var _ KnownValue = StringValue{}

// StringValue is a KnownValue for asserting equality between the value
// supplied to NewStringValue and the value passed to the Equal method.
type StringValue struct {
	value string
}

// Equal determines whether the passed value is of type string, and
// contains a matching sequence of bytes.
func (v StringValue) Equal(other any) bool {
	otherVal, ok := other.(string)

	if !ok {
		return false
	}

	if otherVal != v.value {
		return false
	}

	return true
}

// String returns the string representation of the value.
func (v StringValue) String() string {
	return v.value
}

// NewStringValue returns a KnownValue for asserting equality between the
// supplied string and a value passed to the Equal method.
func NewStringValue(value string) StringValue {
	return StringValue{
		value: value,
	}
}
