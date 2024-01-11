// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"sort"
)

var _ Check = mapValuePartial{}

type mapValuePartial struct {
	value map[string]Check
}

// CheckValue determines whether the passed value is of type map[string]any, and
// contains matching map entries.
func (v mapValuePartial) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for MapValuePartial check, got: %T", other)
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
			return fmt.Errorf("missing element %s for MapValuePartial check", k)
		}

		if err := v.value[k].CheckValue(otherValItem); err != nil {
			return fmt.Errorf("%s map element: %s", k, err)
		}
	}

	return nil
}

// String returns the string representation of the value.
func (v mapValuePartial) String() string {
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

// MapValuePartial returns a Check for asserting partial equality between the
// supplied map[string]Check and the value passed to the CheckValue method. Only
// the elements at the map keys defined within the supplied map[string]Check are
// checked.
func MapValuePartial(value map[string]Check) mapValuePartial {
	return mapValuePartial{
		value: value,
	}
}
