// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "fmt"

var _ Check = stringValueExact{}

type stringValueExact struct {
	value string
}

// CheckValue determines whether the passed value is of type string, and
// contains a matching sequence of bytes.
func (v stringValueExact) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringValueExact check, got: %T", other)
	}

	if otherVal != v.value {
		return fmt.Errorf("expected value %s for StringValueExact check, got: %s", v.value, otherVal)
	}

	return nil
}

// String returns the string representation of the value.
func (v stringValueExact) String() string {
	return v.value
}

// StringValueExact returns a Check for asserting equality between the
// supplied string and a value passed to the CheckValue method.
func StringValueExact(value string) stringValueExact {
	return stringValueExact{
		value: value,
	}
}
