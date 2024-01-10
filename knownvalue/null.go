// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = NullValue{}

// NullValue is a Check for asserting equality between the value supplied
// to NullValueExact and the value passed to the CheckValue method.
type NullValue struct{}

// CheckValue determines whether the passed value is of nil.
func (v NullValue) CheckValue(other any) error {
	if other != nil {
		return fmt.Errorf("expected value nil for NullValue check, got: %T", other)
	}

	return nil
}

// String returns the string representation of nil.
func (v NullValue) String() string {
	return "nil"
}

// NullValueExact returns a Check for asserting equality nil
// and the value passed to the CheckValue method.
func NullValueExact() NullValue {
	return NullValue{}
}
