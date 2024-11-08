// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package echoprovider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEchoProviderServer_primitive(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				provider "echo" {
					data = "hello world"
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("echo.test_one", tfjsonpath.New("data")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test_one", tfjsonpath.New("data"), knownvalue.StringExact("hello world")),
				},
			},
			{
				Config: `
				provider "echo" {
					data = 200
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("echo.test_one", tfjsonpath.New("data")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test_one", tfjsonpath.New("data"), knownvalue.Int64Exact(200)),
				},
			},
			{
				Config: `
				provider "echo" {
					data = true
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("echo.test_one", tfjsonpath.New("data")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test_one", tfjsonpath.New("data"), knownvalue.Bool(true)),
				},
			},
			{
				Config: `
				provider "echo" {
					data = true
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestEchoProviderServer_complex(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				provider "echo" {
					data = tolist(["hello", "world"])
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("echo.test_one", tfjsonpath.New("data")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test_one", tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("hello"),
							knownvalue.StringExact("world"),
						}),
					),
				},
			},
			{
				Config: `
				provider "echo" {
					data = tomap({"key1": "hello", "key2": "world"})
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("echo.test_one", tfjsonpath.New("data")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test_one", tfjsonpath.New("data"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"key1": knownvalue.StringExact("hello"),
							"key2": knownvalue.StringExact("world"),
						}),
					),
				},
			},
			{
				Config: `
				provider "echo" {
					data = tomap({"key1": "hello", "key2": "world"})
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestEchoProviderServer_null(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				provider "echo" {}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("echo.test_one", tfjsonpath.New("data")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test_one", tfjsonpath.New("data"), knownvalue.Null()),
				},
			},
			{
				Config: `
				provider "echo" {
					data = null
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestEchoProviderServer_unknown(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
			},
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				resource "random_string" "str" {
					length = 12
				}
				provider "echo" {
					data = random_string.str.result
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("echo.test_one", tfjsonpath.New("data")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test_one", tfjsonpath.New("data"), knownvalue.StringRegexp(regexp.MustCompile(`\S{12}`))),
				},
			},
			{
				Config: `
				resource "random_string" "str" {
					length = 12
				}
				provider "echo" {
					data = random_string.str.result
				}
				resource "echo" "test_one" {}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
