// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = SetValue{}

// SetValue is a KnownValue for asserting equality between the value supplied
// to SetValueExact and the value passed to the CheckValue method.
type SetValue struct {
	value []Check
}

// CheckValue determines whether the passed value is of type []any, and
// contains matching slice entries independent of the sequence.
func (v SetValue) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("wrong type: %T, known value type is []Check", other)
	}

	if len(otherVal) != len(v.value) {
		return fmt.Errorf("wrong length: %d, known value length is %d", len(otherVal), len(v.value))
	}

	otherValCopy := make([]any, len(otherVal))

	copy(otherValCopy, otherVal)

	for i := 0; i < len(v.value); i++ {
		err := fmt.Errorf("expected value not found: %s", v.value[i].String())

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
func (v SetValue) String() string {
	var setVals []string

	for _, val := range v.value {
		setVals = append(setVals, fmt.Sprintf("%s", val))
	}

	return fmt.Sprintf("%s", setVals)
}

// SetValueExact returns a Check for asserting equality between the
// supplied []Check and the value passed to the CheckValue method.
func SetValueExact(value []Check) SetValue {
	return SetValue{
		value: value,
	}
}
