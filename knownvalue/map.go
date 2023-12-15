// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"sort"
)

var _ KnownValue = MapValue{}

type MapValue struct {
	value map[string]KnownValue
}

// Equal determines whether the passed value is of type map[string]any, and
// contains matching map entries.
func (v MapValue) Equal(other any) bool {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return false
	}

	if len(otherVal) != len(v.value) {
		return false
	}

	for k, v := range v.value {
		otherValItem, ok := otherVal[k]

		if !ok {
			return false
		}

		if !v.Equal(otherValItem) {
			return false
		}
	}

	return true
}

// String returns the string representation of the value.
func (v MapValue) String() string {
	var keys []string

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	mapVals := make(map[string]string, len(keys))

	for _, k := range keys {
		mapVals[k] = fmt.Sprintf("%s", v.value[k])
	}

	return fmt.Sprintf("%v", mapVals)
}

// NewMapValue returns a KnownValue for asserting equality between the
// supplied map[string]KnownValue and the value passed to the Equal method.
func NewMapValue(value map[string]KnownValue) MapValue {
	return MapValue{
		value: value,
	}
}
