package resource

import tfjson "github.com/hashicorp/terraform-json"

var _ PlanAssert = &planAssertSpy{}

type planAssertSpy struct {
	err    error
	called bool
}

func (f *planAssertSpy) RunAssert(_ *tfjson.Plan) error {
	f.called = true
	return f.err
}
