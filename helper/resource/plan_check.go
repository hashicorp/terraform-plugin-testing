package resource

import (
	"context"

	"github.com/hashicorp/go-multierror"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"
)

// TODO: document
type PlanCheck interface {
	CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse)
}

// TODO: document
type CheckPlanRequest struct {
	Plan *tfjson.Plan
}

// TODO: document
type CheckPlanResponse struct {
	Error    error
	SkipTest bool
}

func runPlanChecks(ctx context.Context, t testing.T, plan *tfjson.Plan, planChecks []PlanCheck) error {
	t.Helper()

	var result *multierror.Error

	for _, planCheck := range planChecks {
		resp := CheckPlanResponse{}
		planCheck.CheckPlan(ctx, CheckPlanRequest{Plan: plan}, &resp)

		if resp.SkipTest {
			// TODO: better msg
			t.Skip("skipping test caused by plan check")
		}
		if resp.Error != nil {
			result = multierror.Append(result, resp.Error)
		}
	}

	return result.ErrorOrNil()
}
