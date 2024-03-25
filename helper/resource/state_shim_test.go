// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

func TestStateShimOutput_String(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		Steps: []TestStep{
			{
				Config: `output "test" {
					value = "hello world"
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact("hello world"),
					),
				},
			},
		},
	})
}

func TestStateShimOutput_List(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		Steps: []TestStep{
			{
				Config: `
				output "list_of_strings" {
					value = tolist(["hello", "world"])
				}
				output "list_of_bools" {
					value = tolist([true, false])
				}
				output "list_of_numbers" {
					value = tolist([1.23, 4.56, 500])
				}
				output "list_of_objects" {
					value = tolist([{a = "hey", b = "there"}, {a = "and", b = "another"}])
				}
				output "list_of_lists" {
					value = tolist([tolist(["hey", "there"]), tolist(["and", "another"])])
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"list_of_strings",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.StringExact("hello"),
								knownvalue.StringExact("world"),
							},
						),
					),
					statecheck.ExpectKnownOutputValue(
						"list_of_bools",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.Bool(true),
								knownvalue.Bool(false),
							},
						),
					),
					statecheck.ExpectKnownOutputValue(
						"list_of_numbers",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.Float64Exact(1.23),
								knownvalue.Float64Exact(4.56),
								knownvalue.Float64Exact(500),
							},
						),
					),
					statecheck.ExpectKnownOutputValue(
						"list_of_objects",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"a": knownvalue.StringExact("hey"),
									"b": knownvalue.StringExact("there"),
								}),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"a": knownvalue.StringExact("and"),
									"b": knownvalue.StringExact("another"),
								}),
							},
						),
					),
					statecheck.ExpectKnownOutputValue(
						"list_of_lists",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.ListExact(
									[]knownvalue.Check{
										knownvalue.StringExact("hey"),
										knownvalue.StringExact("there"),
									},
								),
								knownvalue.ListExact(
									[]knownvalue.Check{
										knownvalue.StringExact("and"),
										knownvalue.StringExact("another"),
									},
								),
							},
						),
					),
				},
			},
		},
	})
}

// Ref: https://github.com/hashicorp/terraform-plugin-testing/issues/310
func TestStateShimOutput_Tuple(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		Steps: []TestStep{
			{
				Config: `output "test" {
					value = [true, "hello", 1.23, ["hello", "world"], {a = false}]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.Bool(true),
								knownvalue.StringExact("hello"),
								knownvalue.Float64Exact(1.23),
								knownvalue.ListExact(
									[]knownvalue.Check{
										knownvalue.StringExact("hello"),
										knownvalue.StringExact("world"),
									},
								),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"a": knownvalue.Bool(false),
								}),
							},
						),
					),
				},
			},
			{
				Config: `output "test" {
					value = [{a = false}, true, "hello", 1.23, ["hello", "world"]]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"a": knownvalue.Bool(false),
								}),
								knownvalue.Bool(true),
								knownvalue.StringExact("hello"),
								knownvalue.Float64Exact(1.23),
								knownvalue.ListExact(
									[]knownvalue.Check{
										knownvalue.StringExact("hello"),
										knownvalue.StringExact("world"),
									},
								),
							},
						),
					),
				},
			},
			{
				Config: `output "test" {
					value = [["hello", "world"], {a = false}, true, "hello", 1.23]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.ListExact(
									[]knownvalue.Check{
										knownvalue.StringExact("hello"),
										knownvalue.StringExact("world"),
									},
								),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"a": knownvalue.Bool(false),
								}),
								knownvalue.Bool(true),
								knownvalue.StringExact("hello"),
								knownvalue.Float64Exact(1.23),
							},
						),
					),
				},
			},
			{
				Config: `output "test" {
					value = [1.23, ["hello", "world"], {a = false}, true, "hello"]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.Float64Exact(1.23),
								knownvalue.ListExact(
									[]knownvalue.Check{
										knownvalue.StringExact("hello"),
										knownvalue.StringExact("world"),
									},
								),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"a": knownvalue.Bool(false),
								}),
								knownvalue.Bool(true),
								knownvalue.StringExact("hello"),
							},
						),
					),
				},
			},
			{
				Config: `output "test" {
					value = ["hello", 1.23, ["hello", "world"], {a = false}, true]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.StringExact("hello"),
								knownvalue.Float64Exact(1.23),
								knownvalue.ListExact(
									[]knownvalue.Check{
										knownvalue.StringExact("hello"),
										knownvalue.StringExact("world"),
									},
								),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"a": knownvalue.Bool(false),
								}),
								knownvalue.Bool(true),
							},
						),
					),
				},
			},
		},
	})
}

func TestStateShimOutput_Object(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		Steps: []TestStep{
			{
				Config: `output "test" {
					value = {
						a = "hey",
						b = "there",
						c = true,
						d = 1.23,
					}
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"a": knownvalue.StringExact("hey"),
							"b": knownvalue.StringExact("there"),
							"c": knownvalue.Bool(true),
							"d": knownvalue.Float64Exact(1.23),
						}),
					),
				},
			},
		},
	})
}
