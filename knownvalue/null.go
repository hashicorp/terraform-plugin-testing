// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = nullExact{}

type nullExact struct{}

// CheckValue determines whether the passed value is nil.
func (v nullExact) CheckValue(other any) error {
	if other != nil {
		return fmt.Errorf("expected value nil for NullExact check, got: %T", other)
	}

	return nil
}

// String returns the string representation of null.
func (v nullExact) String() string {
	return "null"
}

// NullExact returns a Check for asserting equality nil
// and the value passed to the CheckValue method.
func NullExact() nullExact {
	return nullExact{}
}
