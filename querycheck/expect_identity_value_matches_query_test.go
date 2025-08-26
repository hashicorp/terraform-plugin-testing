// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectIdentityValueMatchesQuery_CheckQuery_ResourceNotFound(t *testing.T) {
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
					querycheck.ExpectIdentityValueMatchesQuery(
						"examplecloud_thing.two",
						tfjsonpath.New("id"),
					),
				},
				ExpectError: regexp.MustCompile("examplecloud_thing.two - Resource not found in query"),
			},
		},
	})
}

func TestExpectIdentityValueMatchesQuery_CheckQuery_No_Identity(t *testing.T) {
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
					querycheck.ExpectIdentityValueMatchesQuery(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Identity not found in query. Either the resource ` +
					`does not support identity or the Terraform version running the test does not support identity. \(must be v1.14\+\)`,
				),
			},
		},
	})
}

func TestExpectIdentityValueMatchesQuery_CheckQuery_String_Matches(t *testing.T) {
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
					querycheck.ExpectIdentityValueMatchesQuery(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestExpectIdentityValueMatchesQuery_CheckQuery_String_DoesntMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithMismatchedResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentityValueMatchesQuery(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
					),
				},
				ExpectError: regexp.MustCompile("expected identity and query value at path to match, but they differ: examplecloud_thing.one.id, identity value: id-123, query value: 321-di"),
			},
		},
	})
}

func TestExpectIdentityValueMatchesQuery_CheckQuery_List(t *testing.T) {
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
					querycheck.ExpectIdentityValueMatchesQuery(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers"),
					),
				},
			},
		},
	})
}

func TestExpectIdentityValueMatchesQuery_CheckQuery_List_DoesntMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithMismatchedResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentityValueMatchesQuery(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers"),
					),
				},
				ExpectError: regexp.MustCompile(`expected identity and query value at path to match, but they differ: examplecloud_thing.one.list_of_numbers, identity value: \[1 2 3 4\], query value: \[4 3 2 1\]`),
			},
		},
	})
}

func examplecloudProviderWithMismatchedResourceIdentity() func() (tfprotov6.ProviderServer, error) {
	return providerserver.NewProviderServer(testprovider.Provider{
		Resources: map[string]testprovider.Resource{
			"examplecloud_thing": {
				CreateResponse: &resource.CreateResponse{
					NewQuery: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name":            tftypes.String,
								"id":              tftypes.String,
								"list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "test value"),
							"id":   tftypes.NewValue(tftypes.String, "321-di"), // doesn't match identity -> id
							"list_of_numbers": tftypes.NewValue(
								tftypes.List{ElementType: tftypes.Number},
								[]tftypes.Value{
									tftypes.NewValue(tftypes.Number, 4), // doesn't match identity -> list_of_numbers[0]
									tftypes.NewValue(tftypes.Number, 3), // doesn't match identity -> list_of_numbers[1]
									tftypes.NewValue(tftypes.Number, 2), // doesn't match identity -> list_of_numbers[2]
									tftypes.NewValue(tftypes.Number, 1), // doesn't match identity -> list_of_numbers[3]
								},
							),
						},
					),
					NewIdentity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":              tftypes.String,
								"list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"id": tftypes.NewValue(tftypes.String, "id-123"),
							"list_of_numbers": tftypes.NewValue(
								tftypes.List{ElementType: tftypes.Number},
								[]tftypes.Value{
									tftypes.NewValue(tftypes.Number, 1),
									tftypes.NewValue(tftypes.Number, 2),
									tftypes.NewValue(tftypes.Number, 3),
									tftypes.NewValue(tftypes.Number, 4),
								},
							),
						},
					)),
				},
				ReadResponse: &resource.ReadResponse{
					NewQuery: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name":            tftypes.String,
								"id":              tftypes.String,
								"list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "test value"),
							"id":   tftypes.NewValue(tftypes.String, "321-di"), // doesn't match identity -> id
							"list_of_numbers": tftypes.NewValue(
								tftypes.List{ElementType: tftypes.Number},
								[]tftypes.Value{
									tftypes.NewValue(tftypes.Number, 4), // doesn't match identity -> list_of_numbers[0]
									tftypes.NewValue(tftypes.Number, 3), // doesn't match identity -> list_of_numbers[1]
									tftypes.NewValue(tftypes.Number, 2), // doesn't match identity -> list_of_numbers[2]
									tftypes.NewValue(tftypes.Number, 1), // doesn't match identity -> list_of_numbers[3]
								},
							),
						},
					),
					NewIdentity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":              tftypes.String,
								"list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"id": tftypes.NewValue(tftypes.String, "id-123"),
							"list_of_numbers": tftypes.NewValue(
								tftypes.List{ElementType: tftypes.Number},
								[]tftypes.Value{
									tftypes.NewValue(tftypes.Number, 1),
									tftypes.NewValue(tftypes.Number, 2),
									tftypes.NewValue(tftypes.Number, 3),
									tftypes.NewValue(tftypes.Number, 4),
								},
							),
						},
					)),
				},
				IdentitySchemaResponse: &resource.IdentitySchemaResponse{
					Schema: &tfprotov6.ResourceIdentitySchema{
						IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
							{
								Name:              "id",
								Type:              tftypes.String,
								RequiredForImport: true,
							},
							{
								Name:              "list_of_numbers",
								Type:              tftypes.List{ElementType: tftypes.Number},
								OptionalForImport: true,
							},
						},
					},
				},
				SchemaResponse: &resource.SchemaResponse{
					Schema: &tfprotov6.Schema{
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:     "name",
									Type:     tftypes.String,
									Computed: true,
								},
								{
									Name:     "id",
									Type:     tftypes.String,
									Computed: true,
								},
								{
									Name:     "list_of_numbers",
									Type:     tftypes.List{ElementType: tftypes.Number},
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	})
}
