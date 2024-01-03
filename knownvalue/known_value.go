// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

// Check defines an interface that is implemented to determine whether type and value match. Individual
// implementations determine how the match is performed (e.g., exact match, partial match).
type Check interface {
	// CheckValue should perform match testing.
	CheckValue(value any) error
	// String should return a string representation of the type and value.
	String() string
}
