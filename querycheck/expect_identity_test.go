// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectIdentity_CheckQuery_ResourceNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentity(
						"examplecloud_thing.two",
						map[string]knownvalue.Check{
							"id": knownvalue.StringExact("id-123"),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
				ExpectError: regexp.MustCompile("examplecloud_thing.two - Resource not found in query"),
			},
		},
	})
}

func TestExpectIdentity_CheckQuery_No_Terraform_Identity_Support(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
			tfversion.SkipAbove(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			// Resource support identity, but the Terraform versions running will not.
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"id": knownvalue.StringExact("id-123"),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Identity not found in query. Either the resource ` +
					`does not support identity or the Terraform version running the test does not support identity. \(must be v1.12\+\)`,
				),
			},
		},
	})
}

func TestExpectIdentity_CheckQuery_No_Identity(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			// Resource does not support identity
			"examplecloud": examplecloudProviderNoIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"id": knownvalue.StringExact("id-123"),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Identity not found in query. Either the resource ` +
					`does not support identity or the Terraform version running the test does not support identity. \(must be v1.12\+\)`,
				),
			},
		},
	})
}

func TestExpectIdentity_CheckQuery(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
			},
			{
				Query: true,
				Config: `
				provider "examplecloud" {}
				list "examplecloud_containerette" "test" {
					provider = examplecloud
			
					config {
						id = "bat"
					}
				}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"id": knownvalue.StringExact("id-123"),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
			},
		},
	})
}

func TestExpectIdentity_CheckQuery_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{

					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"id": knownvalue.Bool(true),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - "id" identity attribute: expected bool value for Bool check, got: string`),
			},
		},
	})
}

func TestExpectIdentity_CheckQuery_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{

					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"id": knownvalue.StringExact("321-id"),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - "id" identity attribute: expected value 321-id for StringExact check, got: id-123`),
			},
		},
	})
}

func TestExpectIdentity_CheckQuery_ExtraAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{

					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"id": knownvalue.StringExact("321-id"),
						},
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Expected 1 attribute\(s\) in the actual identity object, got 2 attribute\(s\): actual identity has extra attribute\(s\): "list_of_numbers"`),
			},
		},
	})
}

func TestExpectIdentity_CheckQuery_MissingAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{

					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"id":               knownvalue.StringExact("id-123"),
							"nonexistent_attr": knownvalue.StringExact("hello"),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Expected 3 attribute\(s\) in the actual identity object, got 2 attribute\(s\): actual identity is missing attribute\(s\): "nonexistent_attr"`),
			},
		},
	})
}

func TestExpectIdentity_CheckQuery_MismatchedAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentity(
						"examplecloud_thing.one",
						map[string]knownvalue.Check{
							"not_id": knownvalue.StringExact("id-123"),
							"list_of_numbers": knownvalue.ListExact(
								[]knownvalue.Check{
									knownvalue.Int64Exact(1),
									knownvalue.Int64Exact(2),
									knownvalue.Int64Exact(3),
									knownvalue.Int64Exact(4),
								},
							),
						},
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - missing attribute "not_id" in actual identity object`),
			},
		},
	})
}
