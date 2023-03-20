package resource

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

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

func Test_ConfigPlanChecks_PreApply_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{}
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
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigPlanChecks.PreApply spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigPlanChecks.PreApply spy2 to be called at least once")
	}
}

func Test_ConfigPlanChecks_PreApply_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &planCheckSpy{
		err: errors.New("spy3 check failed"),
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
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						spy1,
						spy2,
						spy3,
					},
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 check failed)\n.*?(spy3 check failed)`),
			},
		},
	})
}

func Test_ConfigPlanChecks_PostApplyPreRefresh_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{}
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
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigPlanChecks.PostApplyPreRefresh spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigPlanChecks.PostApplyPreRefresh spy2 to be called at least once")
	}
}

func Test_ConfigPlanChecks_PostApplyPreRefresh_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &planCheckSpy{
		err: errors.New("spy3 check failed"),
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
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
						spy3,
					},
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 check failed)\n.*?(spy3 check failed)`),
			},
		},
	})
}

func Test_ConfigPlanChecks_PostApplyPostRefresh_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{}
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
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigPlanChecks.PostApplyPostRefresh spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigPlanChecks.PostApplyPostRefresh spy2 to be called at least once")
	}
}

func Test_ConfigPlanChecks_PostApplyPostRefresh_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &planCheckSpy{
		err: errors.New("spy3 check failed"),
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
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
						spy3,
					},
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 check failed)\n.*?(spy3 check failed)`),
			},
		},
	})
}
