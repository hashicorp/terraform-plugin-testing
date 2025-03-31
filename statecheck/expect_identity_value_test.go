// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectIdentityValue_CheckState_ResourceNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.two",
						tfjsonpath.New("id"),
						knownvalue.StringExact("id-123"),
					),
				},
				ExpectError: regexp.MustCompile("examplecloud_thing.two - Resource not found in state"),
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_No_Terraform_Identity_Support(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						knownvalue.StringExact("id-123"),
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Identity not found in state. Either the resource ` +
					`does not support identity or the Terraform version running the test does not support identity. \(must be v1.12\+\)`,
				),
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_No_Identity(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			// Resource does not support identity
			"examplecloud": examplecloudProviderNoIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						knownvalue.StringExact("id-123"),
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Identity not found in state. Either the resource ` +
					`does not support identity or the Terraform version running the test does not support identity. \(must be v1.12\+\)`,
				),
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_String(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
			// TODO: There is currently a bug in Terraform v1.12.0-alpha20250319 that causes a panic
			// when refreshing a resource that has an identity stored via protocol v6.
			//
			// We can remove this skip once the bug fix is merged/released:
			// - https://github.com/hashicorp/terraform/pull/36756
			tfversion.SkipIf(version.Must(version.NewVersion("1.12.0-alpha20250319"))),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						knownvalue.StringExact("id-123")),
				},
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_String_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						knownvalue.Bool(true)),
				},
				ExpectError: regexp.MustCompile("expected bool value for Bool check, got: string"),
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_String_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						knownvalue.StringExact("321-id")),
				},
				ExpectError: regexp.MustCompile("expected value 321-id for StringExact check, got: id-123"),
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_List(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
			// TODO: There is currently a bug in Terraform v1.12.0-alpha20250319 that causes a panic
			// when refreshing a resource that has an identity stored via protocol v6.
			//
			// We can remove this skip once the bug fix is merged/released:
			// - https://github.com/hashicorp/terraform/pull/36756
			tfversion.SkipIf(version.Must(version.NewVersion("1.12.0-alpha20250319"))),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers").AtSliceIndex(0),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers").AtSliceIndex(1),
						knownvalue.Int64Exact(2),
					),
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers").AtSliceIndex(2),
						knownvalue.Int64Exact(3),
					),
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers").AtSliceIndex(3),
						knownvalue.Int64Exact(4),
					),
				},
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_List_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers"),
						knownvalue.MapExact(map[string]knownvalue.Check{}),
					),
				},
				ExpectError: regexp.MustCompile(`expected map\[string\]any value for MapExact check, got: \[\]interface {}`),
			},
		},
	})
}

func TestExpectIdentityValue_CheckState_List_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.Int64Exact(4),
							knownvalue.Int64Exact(3),
							knownvalue.Int64Exact(2),
							knownvalue.Int64Exact(1),
						}),
					),
				},
				ExpectError: regexp.MustCompile(`list element index 0: expected value 4 for Int64Exact check, got: 1`),
			},
		},
	})
}

func examplecloudProviderWithResourceIdentity() func() (tfprotov6.ProviderServer, error) {
	return providerserver.NewProviderServer(testprovider.Provider{
		Resources: map[string]testprovider.Resource{
			"examplecloud_thing": {
				CreateResponse: &resource.CreateResponse{
					NewState: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "test value"),
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
					NewState: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "test value"),
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
							},
						},
					},
				},
			},
		},
	})
}

func examplecloudProviderNoIdentity() func() (tfprotov6.ProviderServer, error) {
	return providerserver.NewProviderServer(testprovider.Provider{
		Resources: map[string]testprovider.Resource{
			"examplecloud_thing": {
				CreateResponse: &resource.CreateResponse{
					NewState: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "test value"),
						},
					),
				},
				ReadResponse: &resource.ReadResponse{
					NewState: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "test value"),
						},
					),
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
							},
						},
					},
				},
			},
		},
	})
}
