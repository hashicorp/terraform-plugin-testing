// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ Check = float64ValueExact{}

type float64ValueExact struct {
	value float64
}

// CheckValue determines whether the passed value is of type float64, and
// contains a matching float64 value.
func (v float64ValueExact) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Float64ValueExact check, got: %T", other)
	}

	otherVal, err := jsonNum.Float64()

	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as float64 value for Float64ValueExact check: %s", err)
	}

	if otherVal != v.value {
		return fmt.Errorf("expected value %s for Float64ValueExact check, got: %s", v.String(), strconv.FormatFloat(otherVal, 'f', -1, 64))
	}

	return nil
}

// String returns the string representation of the float64 value.
func (v float64ValueExact) String() string {
	return strconv.FormatFloat(v.value, 'f', -1, 64)
}

// Float64ValueExact returns a Check for asserting equality between the
// supplied float64 and the value passed to the CheckValue method.
func Float64ValueExact(value float64) float64ValueExact {
	return float64ValueExact{
		value: value,
	}
}
