// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "strconv"

var _ KnownValue = NumElementsValue{}

// NumElementsValue is a KnownValue for asserting equality between the value
// supplied to NumElementsExact and the value passed to the Equal method.
type NumElementsValue struct {
	num int
}

// Equal verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v NumElementsValue) Equal(other any) bool {
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
func (v NumElementsValue) String() string {
	return strconv.FormatInt(int64(v.num), 10)
}

// NumElementsExact returns a KnownValue for asserting that
// a list, map, object, or set contains num elements.
func NumElementsExact(num int) NumElementsValue {
	return NumElementsValue{
		num: num,
	}
}
