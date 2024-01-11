// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = setValueExact{}

type setValueExact struct {
	value []Check
}

// CheckValue determines whether the passed value is of type []any, and
// contains matching slice entries independent of the sequence.
func (v setValueExact) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for SetValueExact check, got: %T", other)
	}

	if len(otherVal) != len(v.value) {
		expectedElements := "elements"
		actualElements := "elements"

		if len(v.value) == 1 {
			expectedElements = "element"
		}

		if len(otherVal) == 1 {
			actualElements = "element"
		}

		return fmt.Errorf("expected %d %s for SetValueExact check, got %d %s", len(v.value), expectedElements, len(otherVal), actualElements)
	}

	otherValCopy := make([]any, len(otherVal))

	copy(otherValCopy, otherVal)

	for i := 0; i < len(v.value); i++ {
		err := fmt.Errorf("missing value %s for SetValueExact check", v.value[i].String())

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
func (v setValueExact) String() string {
	var setVals []string

	for _, val := range v.value {
		setVals = append(setVals, val.String())
	}

	return fmt.Sprintf("%s", setVals)
}

// SetValueExact returns a Check for asserting equality between the
// supplied []Check and the value passed to the CheckValue method.
// This is an order-independent check.
func SetValueExact(value []Check) setValueExact {
	return setValueExact{
		value: value,
	}
}
