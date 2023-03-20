package resource

import "context"

var _ PlanCheck = &planCheckSpy{}

type planCheckSpy struct {
	err    error
	skip   string
	called bool
}

func (s *planCheckSpy) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	s.called = true
	resp.Skip = s.skip
	resp.Error = s.err
}
