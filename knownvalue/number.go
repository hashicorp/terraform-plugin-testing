// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"math/big"
)

var _ Check = NumberValue{}

// NumberValue is a KnownValue for asserting equality between the value supplied
// to NumberValueExact and the value passed to the CheckValue method.
type NumberValue struct {
	value *big.Float
}

// CheckValue determines whether the passed value is of type *big.Float, and
// contains a matching *big.Float value.
func (v NumberValue) CheckValue(other any) error {
	if v.value == nil {
		return fmt.Errorf("known value type is nil")
	}

	otherVal, ok := other.(*big.Float)

	if !ok {
		return fmt.Errorf("wrong type: %T, known value type is *big.Float", other)
	}

	if v.value.Cmp(otherVal) != 0 {
		return fmt.Errorf("value: %s does not equal expected value: %s", otherVal.Text('f', -1), v.String())
	}

	return nil
}

// String returns the string representation of the *big.Float value.
func (v NumberValue) String() string {
	return v.value.Text('f', -1)
}

// NumberValueExact returns a Check for asserting equality between the
// supplied *big.Float and the value passed to the CheckValue method.
func NumberValueExact(value *big.Float) NumberValue {
	return NumberValue{
		value: value,
	}
}
