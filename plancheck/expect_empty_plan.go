package plancheck

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/internal/errorshim"
)

var _ PlanCheck = expectEmptyPlan{}

type expectEmptyPlan struct{}

func (e expectEmptyPlan) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var result error

	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			// TODO: Once Go 1.20 is the minimum supported version for this module, replace with `errors.Join` function
			// - https://github.com/hashicorp/terraform-plugin-testing/issues/99
			result = errorshim.Join(result, fmt.Errorf("expected empty plan, but %s has planned action(s): %v", rc.Address, rc.Change.Actions))
		}
	}

	resp.Error = result
}

// TODO: document
func ExpectEmptyPlan() PlanCheck {
	return expectEmptyPlan{}
}
