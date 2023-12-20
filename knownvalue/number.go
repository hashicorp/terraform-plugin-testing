// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"math/big"
)

var _ KnownValue = NumberValue{}

// NumberValue is a KnownValue for asserting equality between the value
// supplied to NewNumberValue and the value passed to the Equal method.
type NumberValue struct {
	value *big.Float
}

// Equal determines whether the passed value is of type *big.Float, and
// contains a matching *big.Float value.
func (v NumberValue) Equal(other any) bool {
	otherVal, ok := other.(*big.Float)

	if !ok {
		return false
	}

	if v.value.Cmp(otherVal) != 0 {
		return false
	}

	return true
}

// String returns the string representation of the *big.Float value.
func (v NumberValue) String() string {
	return v.value.Text('f', -1)
}

// NewNumberValue returns a KnownValue for asserting equality between the
// supplied *big.Float and the value passed to the Equal method.
func NewNumberValue(value *big.Float) NumberValue {
	return NumberValue{
		value: value,
	}
}
