// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = boolExact{}

type boolExact struct {
	value bool
}

// CheckValue determines whether the passed value is of type bool, and
// contains a matching bool value.
func (v boolExact) CheckValue(other any) error {
	otherVal, ok := other.(bool)

	if !ok {
		return fmt.Errorf("expected bool value for BoolExact check, got: %T", other)
	}

	if otherVal != v.value {
		return fmt.Errorf("expected value %t for BoolExact check, got: %t", v.value, otherVal)
	}

	return nil
}

// String returns the string representation of the bool value.
func (v boolExact) String() string {
	return strconv.FormatBool(v.value)
}

// BoolExact returns a Check for asserting equality between the
// supplied bool and the value passed to the CheckValue method.
func BoolExact(value bool) boolExact {
	return boolExact{
		value: value,
	}
}
