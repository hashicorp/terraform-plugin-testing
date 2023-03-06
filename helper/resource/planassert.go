package resource

import (
	"github.com/hashicorp/go-multierror"
	tfjson "github.com/hashicorp/terraform-json"
)

// TODO: document
type PlanAssert interface {
	RunAssert(*tfjson.Plan) error
}

func runPlanAssertions(plan *tfjson.Plan, planAsserts []PlanAssert) error {
	var result *multierror.Error

	for _, planAssert := range planAsserts {
		err := planAssert.RunAssert(plan)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
