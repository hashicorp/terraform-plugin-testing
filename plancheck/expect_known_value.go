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
		resp.Error = fmt.Errorf("%s - Resource not found in plan ResourceChanges", e.resourceAddress)

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
			resp.Error = fmt.Errorf("attribute value: %v does not equal expected value: %s", result, v)
		}
	// Float64 is the default type for all numerical values in tfjson.Plan.
	case reflect.Float64:
		switch t := e.knownValue.(type) {
		case
			knownvalue.Float64Value:
			if !t.Equal(result) {
				resp.Error = fmt.Errorf("attribute value: %v does not equal expected value: %s", result, t)

				return
			}
		case knownvalue.Int64Value:
			// nolint:forcetypeassert // result is reflect.Float64 Kind
			f := result.(float64)

			if !t.Equal(int64(f)) {
				resp.Error = fmt.Errorf("attribute value: %v does not equal expected value: %s", result, t)

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: attribute value is float64 or int64, known value type is %T", t)
		}
	case reflect.Map:
		elems := make(map[string]any)

		val := reflect.ValueOf(result)

		for _, key := range val.MapKeys() {
			elems[key.String()] = val.MapIndex(key).Interface()
		}

		switch t := e.knownValue.(type) {
		case knownvalue.MapValue,
			knownvalue.ObjectValue:
			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute value: %v does not equal expected value: %s", elems, t)

				return
			}
		case knownvalue.MapValuePartial,
			knownvalue.ObjectValuePartial:
			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute value: %v does not contain: %v", elems, t)

				return
			}
		case knownvalue.NumElementsValue:
			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute contains %d elements, expected %v", len(elems), t)

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: attribute value is map, or object, known value type is %T", t)

			return
		}
	case reflect.Slice:
		var elems []any

		var elemsWithIndex []string

		val := reflect.ValueOf(result)

		for i := 0; i < val.Len(); i++ {
			elems = append(elems, val.Index(i).Interface())
			elemsWithIndex = append(elemsWithIndex, fmt.Sprintf("%d:%v", i, val.Index(i).Interface()))
		}

		switch t := e.knownValue.(type) {
		case knownvalue.ListValue,
			knownvalue.SetValue:
			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute value: %v does not equal expected value: %s", elems, t)

				return
			}
		case knownvalue.ListValuePartial:
			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute value: %v does not contain elements at the specified indices: %v", elemsWithIndex, t)

				return
			}
		case knownvalue.NumElementsValue:
			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute contains %d elements, expected %v", len(elems), t)

				return
			}
		case knownvalue.SetValuePartial:
			if !t.Equal(elems) {
				resp.Error = fmt.Errorf("attribute value: %v does not contain: %v", elems, t)

				return
			}
		default:
			resp.Error = fmt.Errorf("wrong type: attribute value is list, or set, known value type is %T", t)

			return
		}
	case reflect.String:
		v, ok := e.knownValue.(knownvalue.StringValue)

		if !ok {
			resp.Error = fmt.Errorf("wrong type: attribute value is string, known value type is %T", e.knownValue)

			return
		}

		if !v.Equal(reflect.ValueOf(result).Interface()) {
			resp.Error = fmt.Errorf("attribute value: %v does not equal expected value: %s", result, v)

			return
		}
	default:
		errorStr := fmt.Sprintf("unrecognised attribute type: %T, known value type is %T", result, e.knownValue)
		errorStr += "\n\nThis is an error in plancheck.ExpectKnownValue.\nPlease report this to the maintainers."

		resp.Error = fmt.Errorf(errorStr)

		return
	}
}

// ExpectKnownValue returns a plan check that asserts that the specified attribute at the given resource
// has a known type and value.
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.KnownValue) PlanCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
