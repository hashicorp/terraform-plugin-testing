package resource

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/errorshim"
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
	Error error

	// Skip, if non-empty, immediately skips further TestStep and defines a message
	// to be included with a call to (*testing.T).Skip(), which is visible in the test output.
	//
	// If a state has been applied via this TestStep or a prior TestStep, the testing will still
	// invoke Terraform to destroy that state before finalizing the skipped test result.
	Skip string
}

func runPlanChecks(ctx context.Context, t testing.T, plan *tfjson.Plan, planChecks []PlanCheck) error {
	t.Helper()

	var result error

	for _, planCheck := range planChecks {
		resp := CheckPlanResponse{}
		planCheck.CheckPlan(ctx, CheckPlanRequest{Plan: plan}, &resp)

		if resp.Skip != "" {
			t.Skip(fmt.Sprintf("plan check forced test skip: %s", resp.Skip))
		}
		if resp.Error != nil {
			// TODO: Once Go 1.20 is the minimum supported version for this module, replace with `errors.Join` function
			// - https://github.com/hashicorp/terraform-plugin-testing/issues/99
			result = errorshim.Join(result, resp.Error)
		}
	}

	return result
}
