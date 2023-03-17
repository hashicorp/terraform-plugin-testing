package plancheck

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_ExpectedResourceAction_NoOp(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionNoop),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_NoOp_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionNoop),
					},
				},
				ExpectError: regexp.MustCompile(`expected NoOp, got action\(s\): \[create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_Create(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionCreate),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_Create_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 15
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionCreate),
					},
				},
				ExpectError: regexp.MustCompile(`expected Create, got action\(s\): \[delete create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_Read(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
			"null": {
				Source: "registry.terraform.io/hashicorp/null",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 15
				}

				data "null_data_source" "two" {
					inputs = {
						unknown_val = random_string.one.result
					}
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("data.null_data_source.two", ResourceActionRead),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_Read_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
			"null": {
				Source: "registry.terraform.io/hashicorp/null",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 15
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionRead),
					},
				},
				ExpectError: regexp.MustCompile(`expected Read, got action\(s\): \[create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_Update(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_offset" "one" {
					offset_days = 1
				  }`,
			},
			{
				Config: `resource "time_offset" "one" {
					offset_days = 2
				  }`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("time_offset.one", ResourceActionUpdate),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_Update_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_offset" "one" {
					offset_days = 1
				  }`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("time_offset.one", ResourceActionUpdate),
					},
				},
				ExpectError: regexp.MustCompile(`expected Update, got action\(s\): \[create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_Destroy(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
			},
			{
				Config: ` `,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionDestroy),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_Destroy_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionDestroy),
					},
				},
				ExpectError: regexp.MustCompile(`expected Destroy, got action\(s\): \[create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_DestroyBeforeCreate(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 15
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_DestroyBeforeCreate_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionDestroyBeforeCreate),
					},
				},
				ExpectError: regexp.MustCompile(`expected DestroyBeforeCreate, got action\(s\): \[create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_CreateBeforeDestroy(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
					lifecycle {
						create_before_destroy = true
					}
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 15
					lifecycle {
						create_before_destroy = true
					}
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionCreateBeforeDestroy),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_CreateBeforeDestroy_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionCreateBeforeDestroy),
					},
				},
				ExpectError: regexp.MustCompile(`expected CreateBeforeDestroy, got action\(s\): \[create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_Replace(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}

				resource "random_string" "two" {
					length = 16
					lifecycle {
						create_before_destroy = true
					}
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 15
				}

				resource "random_string" "two" {
					length = 15
					lifecycle {
						create_before_destroy = true
					}
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionReplace),
						ExpectResourceAction("random_string.two", ResourceActionReplace),
					},
				},
			},
		},
	})
}

func Test_ExpectedResourceAction_Replace_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", ResourceActionReplace),
					},
				},
				ExpectError: regexp.MustCompile(`expected Replace, got action\(s\): \[create\]`),
			},
			{
				Config: `resource "random_string" "two" {
					length = 16
					lifecycle {
						create_before_destroy = true
					}
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.two", ResourceActionReplace),
					},
				},
				ExpectError: regexp.MustCompile(`expected Replace, got action\(s\): \[create\]`),
			},
		},
	})
}

func Test_ExpectedResourceAction_NoResourceFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.doesntexist", ResourceActionCreate),
					},
				},
				ExpectError: regexp.MustCompile(`random_string.doesntexist - Resource not found in plan ResourceChanges`),
			},
		},
	})
}

func Test_ExpectedResourceAction_InvalidResourceActionType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", 0),
					},
				},
				ExpectError: regexp.MustCompile(`random_string.one - unexpected ResourceActionType byte: 0`),
			},
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []r.PlanCheck{
						ExpectResourceAction("random_string.one", 9),
					},
				},
				ExpectError: regexp.MustCompile(`random_string.one - unexpected ResourceActionType byte: 9`),
			},
		},
	})
}
