// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"math/big"
)

var _ Check = NumberValue{}

// NumberValue is a Check for asserting equality between the value supplied
// to NumberValueExact and the value passed to the CheckValue method.
type NumberValue struct {
	value *big.Float
}

// CheckValue determines whether the passed value is of type *big.Float, and
// contains a matching *big.Float value.
func (v NumberValue) CheckValue(other any) error {
	if v.value == nil {
		return fmt.Errorf("value in NumberValue check is nil")
	}

	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for NumberValue check, got: %T", other)
	}

	otherVal, _, err := big.ParseFloat(jsonNum.String(), 10, 512, big.ToNearestEven)

	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as big.Float value for NumberValue check: %s", err)
	}

	if v.value.Cmp(otherVal) != 0 {
		return fmt.Errorf("expected value %s for NumberValue check, got: %s", v.String(), otherVal.Text('f', -1))
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
