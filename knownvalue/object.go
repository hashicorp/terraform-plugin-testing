// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"sort"
)

var _ KnownValue = ObjectValue{}

// ObjectValue is a KnownValue for asserting equality between the value
// supplied to ObjectValueExact and the value passed to the Equal method.
type ObjectValue struct {
	value map[string]KnownValue
}

// Equal determines whether the passed value is of type map[string]any, and
// contains matching object entries.
func (v ObjectValue) Equal(other any) bool {
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
func (v ObjectValue) String() string {
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

// ObjectValueExact returns a KnownValue for asserting equality between the
// supplied map[string]KnownValue and the value passed to the Equal method.
func ObjectValueExact(value map[string]KnownValue) ObjectValue {
	return ObjectValue{
		value: value,
	}
}
