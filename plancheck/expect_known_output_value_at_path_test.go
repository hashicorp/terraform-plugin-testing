// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectKnownOutputValueAtPath_CheckPlan_ResourceNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_two_output",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(true),
						),
					},
				},
				ExpectError: regexp.MustCompile("test_resource_two_output - Output not found in plan OutputChanges"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_AttributeValueNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(true),
						),
					},
				},
				ExpectError: regexp.MustCompile("value is null"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Bool(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(true),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Bool_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("bool_attribute"),
							knownvalue.Float64ValueExact(1.23),
						),
					},
				},
				ExpectError: regexp.MustCompile(`expected json\.Number value for Float64Value check, got: bool`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Bool_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(false),
						),
					},
				},
				ExpectError: regexp.MustCompile("expected value false for BoolValue check, got: true"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Float64(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("float_attribute"),
							knownvalue.Float64ValueExact(1.23),
						),
					},
				},
			},
		},
	})
}

// We do not need equivalent tests for Int64 and Number as they all test the same logic.
func TestExpectKnownOutputValueAtPath_CheckPlan_Float64_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("float_attribute"),
							knownvalue.StringValueExact("str"),
						),
					},
				},
				ExpectError: regexp.MustCompile(`expected string value for StringValue check, got: json\.Number`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Float64_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("float_attribute"),
							knownvalue.Float64ValueExact(3.21),
						),
					},
				},
				ExpectError: regexp.MustCompile("expected value 3.21 for Float64Value check, got: 1.23"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Int64(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("int_attribute"),
							knownvalue.Int64ValueExact(123),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Int64_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("int_attribute"),
							knownvalue.Int64ValueExact(321),
						),
					},
				},
				ExpectError: regexp.MustCompile("expected value 321 for Int64Value check, got: 123"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_List(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_attribute"),
							knownvalue.ListValueExact([]knownvalue.Check{
								knownvalue.StringValueExact("value1"),
								knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_List_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_attribute"),
							knownvalue.MapValueExact(map[string]knownvalue.Check{}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`expected map\[string\]any value for MapValue check, got: \[\]interface {}`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_List_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_attribute"),
							knownvalue.ListValueExact([]knownvalue.Check{
								knownvalue.StringValueExact("value3"),
								knownvalue.StringValueExact("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`list element index 0: expected value value3 for StringValue check, got: value1`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_ListPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_attribute"),
							knownvalue.ListValuePartialMatch(map[int]knownvalue.Check{
								0: knownvalue.StringValueExact("value1"),
							}),
						),
					},
				},
			},
		},
	})
}

// No need to check KnownValueWrongType for ListPartial as all lists, and sets are []any in
// tfjson.Plan.
func TestExpectKnownOutputValueAtPath_CheckPlan_ListPartial_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_attribute"),
							knownvalue.ListValuePartialMatch(map[int]knownvalue.Check{
								0: knownvalue.StringValueExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`list element 0: expected value value3 for StringValue check, got: value1`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_ListElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_attribute"),
							knownvalue.ListElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_ListElements_WrongNum(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_attribute"),
							knownvalue.ListElementsExact(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("expected 3 elements for ListElements check, got 2 elements"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_ListNestedBlock(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
					list_nested_block {
						list_nested_block_attribute = "rts"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_nested_block"),
							knownvalue.ListValueExact([]knownvalue.Check{
								knownvalue.MapValueExact(map[string]knownvalue.Check{
									"list_nested_block_attribute": knownvalue.StringValueExact("str"),
								}),
								knownvalue.MapValueExact(map[string]knownvalue.Check{
									"list_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_ListNestedBlockPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
					list_nested_block {
						list_nested_block_attribute = "rts"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_nested_block"),
							knownvalue.ListValuePartialMatch(map[int]knownvalue.Check{
								1: knownvalue.MapValueExact(map[string]knownvalue.Check{
									"list_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_ListNestedBlockElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
					list_nested_block {
						list_nested_block_attribute = "rts"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("list_nested_block"),
							knownvalue.ListElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Map(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValueExact(map[string]knownvalue.Check{
								"key1": knownvalue.StringValueExact("value1"),
								"key2": knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Map_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("map_attribute"),
							knownvalue.ListValueExact([]knownvalue.Check{}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`expected \[\]any value for ListValue check, got: map\[string\]interface {}`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Map_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValueExact(map[string]knownvalue.Check{
								"key3": knownvalue.StringValueExact("value3"),
								"key4": knownvalue.StringValueExact("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`missing element key3 for MapValue check`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_MapPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
								"key1": knownvalue.StringValueExact("value1"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_MapPartial_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
								"key3": knownvalue.StringValueExact("value1"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`missing element key3 for MapValuePartial check`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_MapElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_MapElements_WrongNum(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapElementsExact(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("expected 3 elements for MapElements check, got 2 elements"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Number(t *testing.T) {
	t.Parallel()

	f, _, err := big.ParseFloat("123", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("int_attribute"),
							knownvalue.NumberValueExact(f),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Number_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	f, _, err := big.ParseFloat("321", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("int_attribute"),
							knownvalue.NumberValueExact(f),
						),
					},
				},
				ExpectError: regexp.MustCompile("expected value 321 for NumberValue check, got: 123"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Set(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValueExact([]knownvalue.Check{
								knownvalue.StringValueExact("value1"),
								knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_Set_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValueExact([]knownvalue.Check{
								knownvalue.StringValueExact("value1"),
								knownvalue.StringValueExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`missing value value3 for SetValue check`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_SetPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValuePartialMatch([]knownvalue.Check{
								knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_SetPartial_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValuePartialMatch([]knownvalue.Check{
								knownvalue.StringValueExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`missing value value3 for SetValuePartial check`),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_SetElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_SetNestedBlock(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_nested_block {
						set_nested_block_attribute = "str"
					}
					set_nested_block {
						set_nested_block_attribute = "rts"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_nested_block"),
							knownvalue.SetValueExact([]knownvalue.Check{
								knownvalue.MapValueExact(map[string]knownvalue.Check{
									"set_nested_block_attribute": knownvalue.StringValueExact("str"),
								}),
								knownvalue.MapValueExact(map[string]knownvalue.Check{
									"set_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_SetNestedBlockPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_nested_block {
						set_nested_block_attribute = "str"
					}
					set_nested_block {
						set_nested_block_attribute = "rts"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_nested_block"),
							knownvalue.SetValuePartialMatch([]knownvalue.Check{
								knownvalue.MapValueExact(map[string]knownvalue.Check{
									"set_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_SetNestedBlockElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_nested_block {
						set_nested_block_attribute = "str"
					}
					set_nested_block {
						set_nested_block_attribute = "rts"
					}
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("set_nested_block"),
							knownvalue.SetElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_String(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					string_attribute = "str"
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("string_attribute"),
							knownvalue.StringValueExact("str")),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_String_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					string_attribute = "str"
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("string_attribute"),
							knownvalue.BoolValueExact(true)),
					},
				},
				ExpectError: regexp.MustCompile("expected bool value for BoolValue check, got: string"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_String_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					string_attribute = "str"
				}

				output test_resource_one_output {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValueAtPath(
							"test_resource_one_output",
							tfjsonpath.New("string_attribute"),
							knownvalue.StringValueExact("rts")),
					},
				},
				ExpectError: regexp.MustCompile("expected value rts for StringValue check, got: str"),
			},
		},
	})
}

func TestExpectKnownOutputValueAtPath_CheckPlan_UnknownAttributeType(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		knownValue  knownvalue.Check
		req         plancheck.CheckPlanRequest
		expectedErr error
	}{
		"unrecognised-type": {
			knownValue: knownvalue.Int64ValueExact(123),
			req: plancheck.CheckPlanRequest{
				Plan: &tfjson.Plan{
					OutputChanges: map[string]*tfjson.Change{
						"float32_output": {
							After: float32(123),
						},
					},
				},
			},
			expectedErr: fmt.Errorf("unrecognised output type: float32, known value type is knownvalue.Int64Value\n\nThis is an error in plancheck.ExpectKnownOutputValueAtPath.\nPlease report this to the maintainers."),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := plancheck.ExpectKnownOutputValueAtPath("float32_output", tfjsonpath.Path{}, testCase.knownValue)

			resp := plancheck.CheckPlanResponse{}

			e.CheckPlan(context.Background(), testCase.req, &resp)

			if diff := cmp.Diff(resp.Error, testCase.expectedErr, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
