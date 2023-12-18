// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ KnownValue = ListValue{}

// ListValue is a KnownValue for asserting equality between the value
// supplied to NewListValue and the value passed to the Equal method.
type ListValue struct {
	value []KnownValue
}

// Equal determines whether the passed value is of type []any, and
// contains matching slice entries in the same sequence.
func (v ListValue) Equal(other any) bool {
	otherVal, ok := other.([]any)

	if !ok {
		return false
	}

	if len(otherVal) != len(v.value) {
		return false
	}

	for i := 0; i < len(v.value); i++ {
		if !v.value[i].Equal(otherVal[i]) {
			return false
		}
	}

	return true
}

// String returns the string representation of the value.
func (v ListValue) String() string {
	var listVals []string

	for _, val := range v.value {
		listVals = append(listVals, fmt.Sprintf("%s", val))
	}

	return fmt.Sprintf("%s", listVals)
}

// NewListValue returns a KnownValue for asserting equality between the
// supplied []KnownValue and the value passed to the Equal method.
func NewListValue(value []KnownValue) ListValue {
	return ListValue{
		value: value,
	}
}
