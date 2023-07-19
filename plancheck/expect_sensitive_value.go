// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ PlanCheck = expectSensitiveValue{}

type expectSensitiveValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

// CheckPlan implements the plan check logic.
func (e expectSensitiveValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var result error

	for _, rc := range req.Plan.ResourceChanges {
		if e.resourceAddress != rc.Address {
			continue
		}

		result, err := tfjsonpath.Traverse(rc.Change.AfterSensitive, e.attributePath)
		if err != nil {
			resp.Error = err
			return
		}

		isSensitive, ok := result.(bool)
		if !ok {
			resp.Error = fmt.Errorf("path not found: cannot convert final value to bool")
			return
		}

		if !isSensitive {
			resp.Error = fmt.Errorf("attribute at path is not sensitive")
			return
		}
	}

	resp.Error = result
}

// ExpectSensitiveValue returns a plan check that asserts that the specified attribute at the given resource has a sensitive value.
//
// Due to implementation differences between the terraform-plugin-sdk and the terraform-plugin-framework, representation of sensitive
// values may differ. For example, terraform-plugin-sdk based providers may have less precise representations of sensitive values, such
// as marking whole maps as sensitive rather than individual element values.
func ExpectSensitiveValue(resourceAddress string, attributePath tfjsonpath.Path) PlanCheck {
	return expectSensitiveValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
