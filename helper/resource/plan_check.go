package resource

import (
	"github.com/hashicorp/go-multierror"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"
)

// TODO: document
type PlanCheck interface {
	RunCheck(req PlanCheckRequest, resp *PlanCheckResponse)
}

// TODO: document
type PlanCheckRequest struct {
	Plan *tfjson.Plan
}

// TODO: document
type PlanCheckResponse struct {
	Error    error
	SkipTest bool
}

func runPlanChecks(t testing.T, plan *tfjson.Plan, planChecks []PlanCheck) error {
	t.Helper()

	var result *multierror.Error

	for _, planCheck := range planChecks {
		resp := PlanCheckResponse{}
		planCheck.RunCheck(PlanCheckRequest{Plan: plan}, &resp)

		if resp.SkipTest {
			t.Skip("skipping test caused by plan check")
		}
		if resp.Error != nil {
			result = multierror.Append(result, resp.Error)
		}
	}

	return result.ErrorOrNil()
}
