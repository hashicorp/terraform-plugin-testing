// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

// KnownValue defines an interface that is implemented to determine equality on the basis of type and value.
type KnownValue interface {
	// Equal should perform equality testing.
	Equal(other any) bool
}
