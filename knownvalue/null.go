// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = null{}

type null struct{}

// CheckValue determines whether the passed value is nil.
func (v null) CheckValue(other any) error {
	if other != nil {
		return fmt.Errorf("expected value nil for Null check, got: %T", other)
	}

	return nil
}

// String returns the string representation of null.
func (v null) String() string {
	return "null"
}

// Null returns a Check for asserting equality nil
// and the value passed to the CheckValue method.
func Null() null {
	return null{}
}
