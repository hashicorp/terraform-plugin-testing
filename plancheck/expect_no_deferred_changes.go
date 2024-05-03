// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"errors"
	"fmt"
)

var _ PlanCheck = expectNoDeferredChanges{}

type expectNoDeferredChanges struct{}

// CheckPlan implements the plan check logic.
func (e expectNoDeferredChanges) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	if len(req.Plan.DeferredChanges) == 0 {
		return
	}

	var result []error
	for _, deferred := range req.Plan.DeferredChanges {
		resourceAddress := "unknown"
		if deferred.ResourceChange != nil {
			resourceAddress = deferred.ResourceChange.Address
		}

		result = append(result, fmt.Errorf("expected no deferred changes, but resource %q is deferred with reason: %q", resourceAddress, deferred.Reason))
	}

	resp.Error = errors.Join(result...)
	if resp.Error != nil {
		return
	}

	if !req.Plan.Complete {
		resp.Error = errors.New("expected plan to be marked as complete, but complete was \"false\", indicating that at least one more plan/apply round is needed to converge.")
		return
	}
}

// ExpectNoDeferredChanges returns a plan check that asserts that there are no deffered changes
// for any resources in the plan.
func ExpectNoDeferredChanges() PlanCheck {
	return expectNoDeferredChanges{}
}
