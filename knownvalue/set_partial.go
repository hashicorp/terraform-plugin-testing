// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ KnownValue = SetValuePartial{}

// SetValuePartial is a KnownValue for asserting equality between the value
// supplied to NewSetValuePartial and the value passed to the Equal method.
type SetValuePartial struct {
	value []KnownValue
}

// Equal determines whether the passed value is of type []any, and
// contains matching slice entries in the same sequence.
func (v SetValuePartial) Equal(other any) bool {
	otherVal, ok := other.([]any)

	if !ok {
		return false
	}

	otherValCopy := make([]any, len(otherVal))

	copy(otherValCopy, otherVal)

	for i := 0; i < len(v.value); i++ {
		var equal bool

		for j := 0; j < len(otherValCopy); j++ {
			equal = v.value[i].Equal(otherValCopy[j])

			if equal {
				otherValCopy[j] = otherValCopy[len(otherValCopy)-1]
				otherValCopy = otherValCopy[:len(otherValCopy)-1]

				break
			}
		}

		if !equal {
			return false
		}
	}

	return true
}

// String returns the string representation of the value.
func (v SetValuePartial) String() string {
	var setVals []string

	for _, val := range v.value {
		setVals = append(setVals, fmt.Sprintf("%s", val))
	}

	return fmt.Sprintf("%s", setVals)
}

// NewSetValuePartial returns a KnownValue for asserting equality of the elements
// supplied in []KnownValue and the elements in the value passed to the Equal method.
func NewSetValuePartial(value []KnownValue) SetValuePartial {
	return SetValuePartial{
		value: value,
	}
}
