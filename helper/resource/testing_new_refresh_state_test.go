// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func Test_RefreshPlanChecks_PostRefresh_Called(t *testing.T) {
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
			},
			{
				RefreshState: true,
				RefreshPlanChecks: RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected RefreshPlanChecks.PostRefresh spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected RefreshPlanChecks.PostRefresh spy2 to be called at least once")
	}
}

func Test_RefreshPlanChecks_PostRefresh_Errors(t *testing.T) {
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
			},
			{
				RefreshState: true,
				RefreshPlanChecks: RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
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
