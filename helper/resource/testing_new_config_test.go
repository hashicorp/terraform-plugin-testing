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

func TestTest_TestStep_ExpectError_NewConfig(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source:            "registry.terraform.io/hashicorp/random",
				VersionConstraint: "3.4.3",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 2
					min_upper = 4
				}`,
				ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value`),
			},
		},
	})
}

func Test_ConfigPlanAsserts_PreApply_Called(t *testing.T) {
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
				ConfigPlanAsserts: ConfigPlanAsserts{
					PreApply: []PlanAssert{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigPlanAsserts.PreApply spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigPlanAsserts.PreApply spy2 to be called at least once")
	}
}

func Test_ConfigPlanAsserts_PreApply_Errors(t *testing.T) {
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
				ConfigPlanAsserts: ConfigPlanAsserts{
					PreApply: []PlanAssert{
						spy1,
						spy2,
						spy3,
					},
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 assert failed)\n.*?(spy3 assert failed)`),
			},
		},
	})
}

// TODO: add post apply, pre refresh plan assert tests
// TODO: add post apply, post refresh plan assert tests
