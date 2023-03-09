package planassert

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ resource.PlanAssert = expectEmptyPlan{}

type expectEmptyPlan struct{}

func (e expectEmptyPlan) RunAssert(plan *tfjson.Plan) error {
	var result *multierror.Error

	for _, rc := range plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			result = multierror.Append(result, fmt.Errorf("expected empty plan, but %s has planned action(s): %v", rc.Address, rc.Change.Actions))
		}
	}

	return result.ErrorOrNil()
}

// TODO: document
func ExpectEmptyPlan() resource.PlanAssert {
	return expectEmptyPlan{}
}
