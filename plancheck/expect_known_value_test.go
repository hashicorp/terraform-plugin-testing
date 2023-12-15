// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

//func TestExpectKnownValue_CheckPlan_Float64(t *testing.T) {
//	t.Parallel()
//
//	testCases := map[string]struct {
//		knownValue knownvalue.KnownValue
//		req        plancheck.CheckPlanRequest
//		expected   plancheck.CheckPlanResponse
//	}{
//		"equal": {
//			knownValue: knownvalue.NewInt64Value(123),
//			req: plancheck.CheckPlanRequest{
//				Plan: &tfjson.Plan{
//					ResourceChanges: []*tfjson.ResourceChange{
//						{
//							Address: "example_resource.test",
//							Change: &tfjson.Change{
//								After: map[string]any{
//									"attribute": float64(123), // tfjson.Plan handles all numerical values as float64.
//								},
//							},
//						},
//					},
//				},
//			},
//			expected: plancheck.CheckPlanResponse{},
//		},
//	}
//
//	for name, testCase := range testCases {
//		name, testCase := name, testCase
//
//		t.Run(name, func(t *testing.T) {
//			t.Parallel()
//
//			e := plancheck.ExpectKnownValue("example_resource.test", tfjsonpath.New("attribute"), testCase.knownValue)
//
//			resp := plancheck.CheckPlanResponse{}
//
//			e.CheckPlan(context.Background(), testCase.req, &resp)
//
//			if diff := cmp.Diff(resp, testCase.expected); diff != "" {
//				t.Errorf("unexpected difference: %s", diff)
//			}
//		})
//	}
//}

func TestExpectKnownValue_CheckPlan_Bool(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("bool_attribute"),
							knownvalue.NewBoolValue(true),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Float64(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("float_attribute"),
							knownvalue.NewFloat64Value(1.23),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Int64(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("int_attribute"),
							knownvalue.NewInt64Value(123),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_List(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value1"),
								knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								0: knownvalue.NewStringValue("value1"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNestedBlock(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.NewStringValue("str"),
								}),
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNestedBlockPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								1: knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNestedBlockNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Map(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
								"key1": knownvalue.NewStringValue("value1"),
								"key2": knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.NewMapValuePartial(map[string]knownvalue.KnownValue{
								"key1": knownvalue.NewStringValue("value1"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Set(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.NewSetValue([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value1"),
								knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.NewSetValuePartial([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNestedBlock(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.NewSetValue([]knownvalue.KnownValue{
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.NewStringValue("str"),
								}),
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNestedBlockPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.NewSetValuePartial([]knownvalue.KnownValue{
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNestedBlockNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_String(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					string_attribute = "str"
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("test_resource.one", tfjsonpath.New("string_attribute"), knownvalue.NewStringValue("str")),
					},
				},
			},
		},
	})
}
