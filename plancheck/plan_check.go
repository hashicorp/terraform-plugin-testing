package plancheck

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
)

// TODO: document
type PlanCheck interface {
	CheckPlan(context.Context, CheckPlanRequest, *CheckPlanResponse)
}

// TODO: document
type CheckPlanRequest struct {
	Plan *tfjson.Plan
}

// TODO: document
type CheckPlanResponse struct {
	Error error
}
