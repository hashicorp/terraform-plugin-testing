// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ PlanCheck = expectUnknownValue{}

type expectUnknownValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

// CheckPlan implements the plan check logic.
func (e expectUnknownValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {

	for _, rc := range req.Plan.ResourceChanges {
		if e.resourceAddress != rc.Address {
			continue
		}

		result, err := tfjsonpath.Traverse(rc.Change.AfterUnknown, e.attributePath)
		if err != nil {
			// If we find the attribute in the known values, return a more explicit message
			knownVal, knownErr := tfjsonpath.Traverse(rc.Change.After, e.attributePath)
			if knownErr == nil {
				resp.Error = fmt.Errorf("Expected unknown value at %q, but found known value: \"%v\"", e.attributePath.String(), knownVal)
				return
			}

			resp.Error = err
			return
		}

		isUnknown, ok := result.(bool)
		if !ok {
			resp.Error = fmt.Errorf("invalid path: the path value cannot be asserted as bool")
			return
		}

		if !isUnknown {
			// The attribute should have a known value, look first to return a more explicit message
			knownVal, knownErr := tfjsonpath.Traverse(rc.Change.After, e.attributePath)
			if knownErr == nil {
				resp.Error = fmt.Errorf("Expected unknown value at %q, but found known value: \"%v\"", e.attributePath.String(), knownVal)
				return
			}

			resp.Error = fmt.Errorf("Expected unknown value at %q, but found known value", e.attributePath.String())
			return
		}

		return
	}

	resp.Error = fmt.Errorf("%s - Resource not found in plan ResourceChanges", e.resourceAddress)
}

// ExpectUnknownValue returns a plan check that asserts that the specified attribute at the given resource has an unknown value.
//
// Due to implementation differences between the terraform-plugin-sdk and the terraform-plugin-framework, representation of unknown
// values may differ. For example, terraform-plugin-sdk based providers may have less precise representations of unknown values, such
// as marking whole maps as unknown rather than individual element values.
func ExpectUnknownValue(resourceAddress string, attributePath tfjsonpath.Path) PlanCheck {
	return expectUnknownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
