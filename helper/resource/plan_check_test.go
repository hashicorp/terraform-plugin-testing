package resource

var _ PlanCheck = &planCheckSpy{}

type planCheckSpy struct {
	err    error
	called bool
}

func (s *planCheckSpy) RunCheck(req PlanCheckRequest, resp *PlanCheckResponse) {
	s.called = true
	resp.Error = s.err
}
