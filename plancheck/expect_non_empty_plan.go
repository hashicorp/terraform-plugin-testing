package plancheck

import (
	"context"
	"errors"
)

var _ PlanCheck = expectNonEmptyPlan{}

type expectNonEmptyPlan struct{}

func (e expectNonEmptyPlan) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			return
		}
	}

	resp.Error = errors.New("expected a non-empty plan, but got an empty plan")
}

// TODO: document
func ExpectNonEmptyPlan() PlanCheck {
	return expectNonEmptyPlan{}
}
