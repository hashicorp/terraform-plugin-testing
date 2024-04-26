// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"
)

var _ PlanCheck = expectDeferredReason{}

type expectDeferredReason struct {
	resourceAddress string
	reason          DeferredReason
}

// CheckPlan implements the plan check logic.
func (e expectDeferredReason) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	foundResource := false

	for _, dc := range req.Plan.DeferredChanges {
		if dc.ResourceChange == nil || e.resourceAddress != dc.ResourceChange.Address {
			continue
		}

		if e.reason != DeferredReason(dc.Reason) {
			resp.Error = fmt.Errorf("'%s' - expected %s, got deferred reason: %s", dc.ResourceChange.Address, e.reason, dc.Reason)
			return
		}

		foundResource = true
		break
	}

	if !foundResource {
		resp.Error = fmt.Errorf("%s - Resource not found in plan DeferredChanges", e.resourceAddress)
		return
	}
}

// TODO: doc
func ExpectDeferredReason(resourceAddress string, reason DeferredReason) PlanCheck {
	return expectDeferredReason{
		resourceAddress: resourceAddress,
		reason:          reason,
	}
}
