package resource

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/errorshim"
	"github.com/mitchellh/go-testing-interface"
)

// TODO: document
type PlanCheck interface {
	CheckPlan(context.Context, CheckPlanRequest, *CheckPlanResponse)
}

// TODO: document
type CheckPlanRequest struct {
	Plan *tfjson.Plan
}

// TODO: document
type CheckPlanResponse struct {
	Error error
}

func runPlanChecks(ctx context.Context, t testing.T, plan *tfjson.Plan, planChecks []PlanCheck) error {
	t.Helper()

	var result error

	for _, planCheck := range planChecks {
		resp := CheckPlanResponse{}
		planCheck.CheckPlan(ctx, CheckPlanRequest{Plan: plan}, &resp)

		if resp.Error != nil {
			// TODO: Once Go 1.20 is the minimum supported version for this module, replace with `errors.Join` function
			// - https://github.com/hashicorp/terraform-plugin-testing/issues/99
			result = errorshim.Join(result, resp.Error)
		}
	}

	return result
}
