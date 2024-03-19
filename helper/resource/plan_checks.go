// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/testingiface"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func runPlanChecks(ctx context.Context, t testingiface.T, plan *tfjson.Plan, planChecks []plancheck.PlanCheck) error {
	t.Helper()

	var result []error

	for _, planCheck := range planChecks {
		resp := plancheck.CheckPlanResponse{}
		planCheck.CheckPlan(ctx, plancheck.CheckPlanRequest{Plan: plan}, &resp)

		result = append(result, resp.Error)
	}

	return errors.Join(result...)
}
