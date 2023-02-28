package resource

import (
	"errors"
	"regexp"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

var _ PlanAssert = &planAssertSpy{}

type planAssertSpy struct {
	err    error
	called bool
}

func (f *planAssertSpy) RunAssert(_ *tfjson.Plan) error {
	f.called = true
	return f.err
}

func Test_PreApplyPlanAssert_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planAssertSpy{}
	spy2 := &planAssertSpy{}
	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				PreApplyPlanAsserts: []PlanAssert{
					spy1,
					spy2,
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected PreApplyPlanAssert spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected PreApplyPlanAssert spy2 to be called at least once")
	}
}

func Test_PreApplyPlanAssert_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planAssertSpy{}
	spy2 := &planAssertSpy{
		err: errors.New("spy2 assert failed"),
	}
	spy3 := &planAssertSpy{
		err: errors.New("spy3 assert failed"),
	}
	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				PreApplyPlanAsserts: []PlanAssert{
					spy1,
					spy2,
					spy3,
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 assert failed)\n.*?(spy3 assert failed)`),
			},
		},
	})
}
