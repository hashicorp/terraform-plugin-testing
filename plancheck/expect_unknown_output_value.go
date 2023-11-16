// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ PlanCheck = expectUnknownOutputValue{}

type expectUnknownOutputValue struct {
	outputAddress string
	attributePath tfjsonpath.Path
}

// CheckPlan implements the plan check logic.
func (e expectUnknownOutputValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
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

	result, err := tfjsonpath.Traverse(change.AfterUnknown, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	isUnknown, ok := result.(bool)

	if !ok {
		resp.Error = fmt.Errorf("invalid path: the path value cannot be asserted as bool")

		return
	}

	if !isUnknown {
		resp.Error = fmt.Errorf("attribute at path is known")

		return
	}

	return

}

// ExpectUnknownOutputValue returns a plan check that asserts that the specified attribute at the given resource has an unknown value.
//
// Due to implementation differences between the terraform-plugin-sdk and the terraform-plugin-framework, representation of unknown
// values may differ. For example, terraform-plugin-sdk based providers may have less precise representations of unknown values, such
// as marking whole maps as unknown rather than individual element values.
func ExpectUnknownOutputValue(params OutputValueParams) PlanCheck {
	return expectUnknownOutputValue{
		outputAddress: params.OutputAddress,
		attributePath: params.AttributePath,
	}
}

// OutputValueParams is used during the creation of a plan check for output values, and specifies
// the address and optional attribute path for an output value.
//
// For example, if an output has been specified to point at a specific value:
//
//	resource "time_static" "one" {}
//
//	output "string_attribute" {
//	    value = time_static.one.rfc3339
//	}
//
// Then the value can be addressed directly and does not require an attributePath:
//
//	plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
//	    OutputAddress: "string_attribute",
//	}),
//
// However, if an output has been specified to point at an object or a collection.
// For example:
//
//	resource "time_static" "one" {}
//
//	output "string_attribute" {
//	    value = time_static.one
//	}
//
// Then the value cannot be addressed directly and requires an attributePath:
//
//	plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
//	    OutputAddress: "string_attribute",
//	    AttributePath: tfjsonpath.New("rfc3339"),
//	}),
type OutputValueParams struct {
	OutputAddress string
	AttributePath tfjsonpath.Path
}
