// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

var _ KnownValue = ListValuePartial{}

type ListValuePartial struct {
	value map[int]KnownValue
}

// Equal determines whether the passed value is of type []any, and
// contains matching slice entries in the same sequence.
func (v ListValuePartial) Equal(other any) bool {
	otherVal, ok := other.([]any)

	if !ok {
		return false
	}

	for k, val := range v.value {
		if len(otherVal) <= k {
			return false
		}

		if !val.Equal(otherVal[k]) {
			return false
		}
	}

	return true
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

// NewListValuePartial returns a KnownValue for asserting equality between the
// supplied map[int]KnownValue and the value passed to the Equal method.
func NewListValuePartial(value map[int]KnownValue) ListValuePartial {
	return ListValuePartial{
		value: value,
	}
}
