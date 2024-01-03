// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "fmt"

var _ Check = StringValue{}

// StringValue is a KnownValue for asserting equality between the value
// supplied to StringValueExact and the value passed to the CheckValue method.
type StringValue struct {
	value string
}

// CheckValue determines whether the passed value is of type string, and
// contains a matching sequence of bytes.
func (v StringValue) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringValue check, got: %T", other)
	}

	if otherVal != v.value {
		return fmt.Errorf("expected value %s for StringValue check, got: %s", v.value, otherVal)
	}

	return nil
}

// String returns the string representation of the value.
func (v StringValue) String() string {
	return v.value
}

// StringValueExact returns a Check for asserting equality between the
// supplied string and a value passed to the CheckValue method.
func StringValueExact(value string) StringValue {
	return StringValue{
		value: value,
	}
}
