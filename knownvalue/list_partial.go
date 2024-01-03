// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

var _ Check = ListValuePartial{}

// ListValuePartial is a KnownValue for asserting equality between the value supplied
// to ListValuePartialMatch and the value passed to the CheckValue method.
type ListValuePartial struct {
	value map[int]Check
}

// CheckValue determines whether the passed value is of type []any, and
// contains matching slice entries in the same sequence.
func (v ListValuePartial) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("wrong type: %T, known value type is map[int]Check", other)
	}

	var keys []int

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		if len(otherVal) <= k {
			return fmt.Errorf("index out of bounds: %d", k)
		}

		if err := v.value[k].CheckValue(otherVal[k]); err != nil {
			return err
		}
	}

	return nil
}

// String returns the string representation of the value.
func (v ListValuePartial) String() string {
	var b bytes.Buffer

	b.WriteString("[")

	var keys []int

	var listVals []string

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		listVals = append(listVals, fmt.Sprintf("%d:%s", k, v.value[k]))
	}

	b.WriteString(strings.Join(listVals, " "))

	b.WriteString("]")

	return b.String()
}

// ListValuePartialMatch returns a Check for asserting equality between the
// supplied map[int]KnownValue and the value passed to the CheckValue method. The
// map keys correspond to the position of the zero-ordered element within the
// list that is being checked.
func ListValuePartialMatch(value map[int]Check) ListValuePartial {
	return ListValuePartial{
		value: value,
	}
}
