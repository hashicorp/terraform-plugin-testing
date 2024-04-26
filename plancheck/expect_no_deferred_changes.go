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
}

// TODO: doc
func ExpectNoDeferredChanges() PlanCheck {
	return expectNoDeferredChanges{}
}
