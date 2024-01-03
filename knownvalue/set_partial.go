// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = SetValuePartial{}

// SetValuePartial is a KnownValue for asserting equality between the value supplied
// to SetValuePartialMatch and the value passed to the CheckValue method.
type SetValuePartial struct {
	value []Check
}

// CheckValue determines whether the passed value is of type []any, and
// contains matching slice entries in any sequence.
func (v SetValuePartial) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for SetValuePartial check, got: %T", other)
	}

	otherValCopy := make([]any, len(otherVal))

	copy(otherValCopy, otherVal)

	for i := 0; i < len(v.value); i++ {
		err := fmt.Errorf("missing value %s for SetValuePartial check", v.value[i].String())

		for j := 0; j < len(otherValCopy); j++ {
			checkValueErr := v.value[i].CheckValue(otherValCopy[j])

			if checkValueErr == nil {
				otherValCopy[j] = otherValCopy[len(otherValCopy)-1]
				otherValCopy = otherValCopy[:len(otherValCopy)-1]

				err = nil

				break
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// String returns the string representation of the value.
func (v SetValuePartial) String() string {
	var setVals []string

	for _, val := range v.value {
		setVals = append(setVals, val.String())
	}

	return fmt.Sprintf("%s", setVals)
}

// SetValuePartialMatch returns a Check for asserting equality of the elements
// supplied in []Check and the elements in the value passed to the CheckValue method.
func SetValuePartialMatch(value []Check) SetValuePartial {
	return SetValuePartial{
		value: value,
	}
}
