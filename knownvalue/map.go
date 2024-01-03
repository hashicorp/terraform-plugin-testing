// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"sort"
)

var _ Check = MapValue{}

// MapValue is a KnownValue for asserting equality between the value supplied
// to MapValueExact and the value passed to the CheckValue method.
type MapValue struct {
	value map[string]Check
}

// CheckValue determines whether the passed value is of type map[string]any, and
// contains matching map entries.
func (v MapValue) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("wrong type: %T, known value type is map[string]Check", other)
	}

	if len(otherVal) != len(v.value) {
		return fmt.Errorf("wrong length: %d, known value length is %d", len(otherVal), len(v.value))
	}

	var keys []string

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		otherValItem, ok := otherVal[k]

		if !ok {
			return fmt.Errorf("missing key: %s", k)
		}

		if err := v.value[k].CheckValue(otherValItem); err != nil {
			return err
		}
	}

	return nil
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
		mapVals[k] = v.value[k].String()
	}

	return fmt.Sprintf("%v", mapVals)
}

// MapValueExact returns a Check for asserting equality between the
// supplied map[string]KnownValue and the value passed to the CheckValue method.
func MapValueExact(value map[string]Check) MapValue {
	return MapValue{
		value: value,
	}
}
