// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "strconv"

var _ KnownValue = NumElements{}

type NumElements struct {
	num int
}

// Equal verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v NumElements) Equal(other any) bool {
	mapVal, mapOk := other.(map[string]any)

	sliceVal, sliceOk := other.([]any)

	if !mapOk && !sliceOk {
		return false
	}

	if mapOk && len(mapVal) != v.num {
		return false
	}

	if sliceOk && len(sliceVal) != v.num {
		return false
	}

	return true
}

// String returns the string representation of the value.
func (v NumElements) String() string {
	return strconv.FormatInt(int64(v.num), 10)
}

// NewNumElements returns a KnownValue for asserting that
// a list, map, object, or set contains num elements.
func NewNumElements(num int) NumElements {
	return NumElements{
		num: num,
	}
}
