package resource

var _ PlanCheck = &planCheckSpy{}

type planCheckSpy struct {
	err    error
	skip   bool
	called bool
}

func (s *planCheckSpy) RunCheck(req PlanCheckRequest, resp *PlanCheckResponse) {
	s.called = true
	resp.SkipTest = s.skip
	resp.Error = s.err
}
