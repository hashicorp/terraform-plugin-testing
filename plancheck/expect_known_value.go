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
	knownValue      knownvalue.Check
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
		resp.Error = fmt.Errorf("value is null")

		return
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.Bool,
		reflect.Map,
		reflect.Slice,
		reflect.String:
		if err := e.knownValue.CheckValue(result); err != nil {
			resp.Error = err

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
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) PlanCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
