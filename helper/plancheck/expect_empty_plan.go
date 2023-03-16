package plancheck

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ resource.PlanCheck = expectEmptyPlan{}

type expectEmptyPlan struct{}

func (e expectEmptyPlan) RunCheck(req resource.PlanCheckRequest, resp *resource.PlanCheckResponse) {
	var result *multierror.Error

	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			result = multierror.Append(result, fmt.Errorf("expected empty plan, but %s has planned action(s): %v", rc.Address, rc.Change.Actions))
		}
	}

	resp.Error = result.ErrorOrNil()
}

// TODO: document
func ExpectEmptyPlan() resource.PlanCheck {
	return expectEmptyPlan{}
}
