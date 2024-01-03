// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"sort"
)

var _ Check = ObjectValuePartial{}

// ObjectValuePartial is a KnownValue for asserting equality between the value supplied
// to ObjectValuePartialMatch and the value passed to the CheckValue method.
type ObjectValuePartial struct {
	value map[string]Check
}

// CheckValue determines whether the passed value is of type map[string]any, and
// contains matching map entries.
func (v ObjectValuePartial) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("wrong type: %T, known value type is map[string]Check", other)
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
func (v ObjectValuePartial) String() string {
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

// ObjectValuePartialMatch returns a Check for asserting partial equality between the
// supplied map[string]KnownValue and the value passed to the CheckValue method.
func ObjectValuePartialMatch(value map[string]Check) ObjectValuePartial {
	return ObjectValuePartial{
		value: value,
	}
}
