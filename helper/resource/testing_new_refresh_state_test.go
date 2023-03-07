package resource

import (
	"errors"
	"regexp"
	"testing"
)

func Test_RefreshPlanAsserts_PostRefresh_Called(t *testing.T) {
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
			},
			{
				RefreshState: true,
				RefreshPlanAsserts: RefreshPlanAsserts{
					PostRefresh: []PlanAssert{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected RefreshPlanAsserts.PostRefresh spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected RefreshPlanAsserts.PostRefresh spy2 to be called at least once")
	}
}

func Test_RefreshPlanAsserts_PostRefresh_Errors(t *testing.T) {
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
			},
			{
				RefreshState: true,
				RefreshPlanAsserts: RefreshPlanAsserts{
					PostRefresh: []PlanAssert{
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
