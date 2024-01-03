// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource Plan Check
var _ PlanCheck = expectKnownOutputValue{}

type expectKnownOutputValue struct {
	outputAddress string
	knownValue    knownvalue.Check
}

// CheckPlan implements the plan check logic.
func (e expectKnownOutputValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var change *tfjson.Change

	for address, oc := range req.Plan.OutputChanges {
		if e.outputAddress == address {
			change = oc

			break
		}
	}

	if change == nil {
		resp.Error = fmt.Errorf("%s - Output not found in plan OutputChanges", e.outputAddress)

		return
	}

	result, err := tfjsonpath.Traverse(change.After, tfjsonpath.Path{})

	if err != nil {
		resp.Error = err

		return
	}

	if result == nil {
		resp.Error = fmt.Errorf("value is null")

		return
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.Bool:
		v, ok := e.knownValue.(knownvalue.BoolValue)

		if !ok {
			resp.Error = fmt.Errorf("wrong type: value is bool, known value type is %T", e.knownValue)

			return
		}

		if err := v.CheckValue(result); err != nil {
			resp.Error = err

			return
		}
	case reflect.Map:
		elems := make(map[string]any)

		val := reflect.ValueOf(result)

		for _, key := range val.MapKeys() {
			elems[key.String()] = val.MapIndex(key).Interface()
		}

		switch t := e.knownValue.(type) {
		case knownvalue.MapElements,
			knownvalue.MapValue,
			knownvalue.MapValuePartial,
			knownvalue.ObjectAttributes,
			knownvalue.ObjectValue,
			knownvalue.ObjectValuePartial:
			if err := t.CheckValue(elems); err != nil {
				resp.Error = err

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: value is map, or object, known value type is %T", t)

			return
		}
	case reflect.Slice:
		var elems []any

		val := reflect.ValueOf(result)

		for i := 0; i < val.Len(); i++ {
			elems = append(elems, val.Index(i).Interface())
		}

		switch t := e.knownValue.(type) {
		case knownvalue.ListElements,
			knownvalue.ListValue,
			knownvalue.ListValuePartial,
			knownvalue.SetElements,
			knownvalue.SetValue,
			knownvalue.SetValuePartial:
			if err := t.CheckValue(elems); err != nil {
				resp.Error = err

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: value is list, or set, known value type is %T", t)

			return
		}
	case reflect.String:
		jsonNum, jsonNumOk := result.(json.Number)

		if jsonNumOk {
			float64Val, float64ValOk := e.knownValue.(knownvalue.Float64Value)

			int64Val, int64ValOk := e.knownValue.(knownvalue.Int64Value)

			numberValue, numberValOk := e.knownValue.(knownvalue.NumberValue)

			if !float64ValOk && !int64ValOk && !numberValOk {
				resp.Error = fmt.Errorf("wrong type: value is number, known value type is %T", e.knownValue)
			}

			switch {
			case float64ValOk:
				f, err := jsonNum.Float64()

				if err != nil {
					resp.Error = fmt.Errorf("%q could not be parsed as float64", jsonNum.String())

					return
				}

				if err := float64Val.CheckValue(f); err != nil {
					resp.Error = err

					return
				}
			case int64ValOk:
				i, err := jsonNum.Int64()

				if err != nil {
					resp.Error = fmt.Errorf("%q could not be parsed as int64", jsonNum.String())

					return
				}

				if err := int64Val.CheckValue(i); err != nil {
					resp.Error = err

					return
				}
			case numberValOk:
				f, _, err := big.ParseFloat(jsonNum.String(), 10, 512, big.ToNearestEven)

				if err != nil {
					resp.Error = fmt.Errorf("%q could not be parsed as big.Float", jsonNum.String())

					return
				}

				if err := numberValue.CheckValue(f); err != nil {
					resp.Error = err

					return
				}
			}
		} else {
			v, ok := e.knownValue.(knownvalue.StringValue)

			if !ok {
				resp.Error = fmt.Errorf("wrong type: value is string, known value type is %T", e.knownValue)

				return
			}

			if err := v.CheckValue(result); err != nil {
				resp.Error = err

				return
			}
		}
	default:
		errorStr := fmt.Sprintf("unrecognised output type: %T, known value type is %T", result, e.knownValue)
		errorStr += "\n\nThis is an error in plancheck.ExpectKnownOutputValue.\nPlease report this to the maintainers."

		resp.Error = fmt.Errorf(errorStr)

		return
	}
}

// ExpectKnownOutputValue returns a plan check that asserts that the specified value
// has a known type, and value.
func ExpectKnownOutputValue(outputAddress string, knownValue knownvalue.Check) PlanCheck {
	return expectKnownOutputValue{
		outputAddress: outputAddress,
		knownValue:    knownValue,
	}
}
