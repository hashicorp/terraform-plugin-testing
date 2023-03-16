package plancheck

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ resource.PlanCheck = expectNonEmptyPlan{}

type expectNonEmptyPlan struct{}

func (e expectNonEmptyPlan) RunCheck(req resource.PlanCheckRequest, resp *resource.PlanCheckResponse) {
	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			return
		}
	}

	resp.Error = errors.New("expected a non-empty plan, but got an empty plan")
}

// TODO: document
func ExpectNonEmptyPlan() resource.PlanCheck {
	return expectNonEmptyPlan{}
}
