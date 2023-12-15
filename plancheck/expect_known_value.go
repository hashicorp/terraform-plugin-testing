// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource Plan Check
var _ PlanCheck = expectKnownValue{}

type expectKnownValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	knownValue      knownvalue.KnownValue
}

// CheckPlan implements the plan check logic.
func (e expectKnownValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var rc *tfjson.ResourceChange

	for _, resourceChange := range req.Plan.ResourceChanges {
		if e.resourceAddress == resourceChange.Address {
			rc = resourceChange

			break
		}
	}

	if rc == nil {
		resp.Error = fmt.Errorf("%s - resource not found in plan", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(rc.Change.After, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if result == nil {
		resp.Error = fmt.Errorf("attribute value is null")

		return
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.Bool:
		v, ok := e.knownValue.(knownvalue.BoolValue)

		if !ok {
			resp.Error = fmt.Errorf("wrong type: attribute value is bool, known value type is %T", e.knownValue)

			return
		}

		if !v.Equal(reflect.ValueOf(result).Interface()) {
			resp.Error = fmt.Errorf("attribute value does not equal expected value: %v != %v", result, v)
		}
	// Float64 is the default type for all numerical values.
	case reflect.Float64:
		switch t := e.knownValue.(type) {
		case
			knownvalue.Float64Value:
			if !t.Equal(result) {
				resp.Error = fmt.Errorf("attribute value %v is not equal to %s", result, t)

				return
			}
		case knownvalue.Int64Value:
			// nolint:forcetypeassert // result is reflect.Float64 Kind
			f := result.(float64)

			if !t.Equal(int64(f)) {
				resp.Error = fmt.Errorf("attribute value %v is not equal to %s", result, t)

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: attribute type is %T, known value type is %T", result, t)
		}
	case reflect.Map:
		switch t := e.knownValue.(type) {
		case knownvalue.MapValue,
			knownvalue.MapValuePartial,
			knownvalue.NumElements,
			knownvalue.ObjectValue,
			knownvalue.ObjectValuePartial:

			elems := make(map[string]any)

			val := reflect.ValueOf(result)

			for _, key := range val.MapKeys() {
				elems[key.String()] = val.MapIndex(key).Interface()
			}

			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute %v is not equal to expected value %v", elems, t)

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: attribute type is list, or set, known value type is %T", t)

			return
		}
	case reflect.Slice:
		switch t := e.knownValue.(type) {
		case knownvalue.ListValue,
			knownvalue.ListValuePartial,
			knownvalue.NumElements,
			knownvalue.SetValue,
			knownvalue.SetValuePartial:

			var elems []any

			var elemsWithIndex []string

			val := reflect.ValueOf(result)

			for i := 0; i < val.Len(); i++ {
				elems = append(elems, val.Index(i).Interface())
				elemsWithIndex = append(elemsWithIndex, fmt.Sprintf("%d:%v", i, val.Index(i).Interface()))
			}

			if !t.Equal(elems) {
				switch e.knownValue.(type) {
				case knownvalue.ListValuePartial:
					resp.Error = fmt.Errorf("attribute %v does not contain elements at the specified indices %v", elemsWithIndex, t)

					return
				case knownvalue.NumElements:
					resp.Error = fmt.Errorf("attribute contains %d elements, expected %v", len(elems), t)

					return
				case knownvalue.SetValuePartial:
					resp.Error = fmt.Errorf("attribute %v does not contain %v", elems, t)

					return
				}

				resp.Error = fmt.Errorf("attribute %v is not equal to expected value %v", elems, t)

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: attribute type is list, or set, known value type is %T", t)

			return
		}
	// String will need to handle json.Number if tfjson.Plan is modified to use json.Number for numerical values.
	case reflect.String:
		v, ok := e.knownValue.(knownvalue.StringValue)

		if !ok {
			resp.Error = fmt.Errorf("wrong type: attribute value is string, known value type is %T", e.knownValue)

			return
		}

		if !v.Equal(reflect.ValueOf(result).Interface()) {
			resp.Error = fmt.Errorf("attribute value does not equal expected value: %v != %v", result, v)

			return
		}
	default:
		resp.Error = fmt.Errorf("wrong type: attribute type is %T, known value type is %T", result, e.knownValue)

		return
	}
}

// ExpectKnownValue returns a plan check that asserts that the specified attribute at the given resource has a known type, and value.
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.KnownValue) PlanCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
