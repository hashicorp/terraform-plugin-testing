package planassert

import (
	"errors"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ resource.PlanAssert = expectNonEmptyPlan{}

type expectNonEmptyPlan struct{}

func (e expectNonEmptyPlan) RunAssert(plan *tfjson.Plan) error {
	for _, rc := range plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			return nil
		}
	}

	return errors.New("expected a non-empty plan, but got an empty plan")
}

// TODO: document
func ExpectNonEmptyPlan() resource.PlanAssert {
	return expectNonEmptyPlan{}
}
