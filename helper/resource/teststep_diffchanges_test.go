package resource

import (
	"regexp"
	"testing"
)

func TestTest_TestStep_ExpectedResourceChanges_NoOp(t *testing.T) {
	t.Parallel()

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
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffNoop,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_NoOp_NoMatch(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffNoop,
				},
				ExpectError: regexp.MustCompile(`expected NoOp, got action\(s\): \[create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Create(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffCreate,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Create_NoMatch(t *testing.T) {
	t.Parallel()

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
				Config: `resource "random_string" "one" {
					length = 15
				}`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffCreate,
				},
				ExpectError: regexp.MustCompile(`expected Create, got action\(s\): \[delete create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Read(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
			"null": {
				Source: "registry.terraform.io/hashicorp/null",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 15
				}

				data "null_data_source" "two" {
					inputs = {
						unknown_val = random_string.one.result
					}
				}`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"data.null_data_source.two": DiffRead,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Read_NoMatch(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
			"null": {
				Source: "registry.terraform.io/hashicorp/null",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 15
				}`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffRead,
				},
				ExpectError: regexp.MustCompile(`expected Read, got action\(s\): \[create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Update(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "time_offset" "one" {
					offset_days = 1
				  }`,
			},
			{
				Config: `resource "time_offset" "one" {
					offset_days = 2
				  }`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"time_offset.one": DiffUpdate,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Update_NoMatch(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "time_offset" "one" {
					offset_days = 1
				  }`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"time_offset.one": DiffUpdate,
				},
				ExpectError: regexp.MustCompile(`expected Update, got action\(s\): \[create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Destroy(t *testing.T) {
	t.Parallel()

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
				Config: ` `,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffDestroy,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Destroy_NoMatch(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffDestroy,
				},
				ExpectError: regexp.MustCompile(`expected Destroy, got action\(s\): \[create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_DestroyBeforeCreate(t *testing.T) {
	t.Parallel()

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
				Config: `resource "random_string" "one" {
					length = 15
				}`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffDestroyBeforeCreate,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_DestroyBeforeCreate_NoMatch(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffDestroyBeforeCreate,
				},
				ExpectError: regexp.MustCompile(`expected DestroyBeforeCreate, got action\(s\): \[create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_CreateBeforeDestroy(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffCreateBeforeDestroy,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_CreateBeforeDestroy_NoMatch(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffCreateBeforeDestroy,
				},
				ExpectError: regexp.MustCompile(`expected CreateBeforeDestroy, got action\(s\): \[create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Replace(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffReplace,
					"random_string.two": DiffReplace,
				},
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_Replace_NoMatch(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": DiffReplace,
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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.two": DiffReplace,
				},
				ExpectError: regexp.MustCompile(`expected Replace, got action\(s\): \[create\]`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_NoResourceFound(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.doesntexist": DiffCreate,
				},
				ExpectError: regexp.MustCompile(`random_string.doesntexist - Resource not found in planned ResourceChanges`),
			},
		},
	})
}

func TestTest_TestStep_ExpectedResourceChanges_InvalidDiffChangeType(t *testing.T) {
	t.Parallel()

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
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": 0,
				},
				ExpectError: regexp.MustCompile(`random_string.one - unexpected DiffChangeType byte: 0`),
			},
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ExpectedResourceChanges: map[string]DiffChangeType{
					"random_string.one": 9,
				},
				ExpectError: regexp.MustCompile(`random_string.one - unexpected DiffChangeType byte: 9`),
			},
		},
	})
}
