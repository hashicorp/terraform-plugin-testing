package plancheck

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ PlanCheck = expectNoValue{}

type expectNoValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

// CheckPlan implements validation that the plan does not contain any values for the specified resource attribute.
func (e expectNoValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	for _, rc := range req.Plan.ResourceChanges {
		if e.resourceAddress != rc.Address {
			continue
		}
		v, err := tfjsonpath.Traverse(rc.Change.After, e.attributePath)
		if v == nil || err != nil {
			v, err := tfjsonpath.Traverse(rc.Change.AfterUnknown, e.attributePath)
			if v == nil || err != nil {
				return
			}
			resp.Error = fmt.Errorf(
				"expected %s to have no value, got %v",
				e.attributePath,
				v,
			)
			return
		}

		resp.Error = fmt.Errorf(
			"expected %s to have no value, got %v",
			e.attributePath,
			v,
		)
		return
	}
}

func ExpectNoValue(resourceAddress string, attributePath tfjsonpath.Path) PlanCheck {
	return expectNoValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
