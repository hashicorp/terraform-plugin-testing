// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

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
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectIdentityValueMatchesStateAtPath_CheckState_ResourceNotFound(t *testing.T) {
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
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"examplecloud_thing.two",
						tfjsonpath.New("id"),
						tfjsonpath.New("id"),
					),
				},
				ExpectError: regexp.MustCompile("examplecloud_thing.two - Resource not found in state"),
			},
		},
	})
}

func TestExpectIdentityValueMatchesStateAtPath_CheckState_No_Terraform_Identity_Support(t *testing.T) {
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
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						tfjsonpath.New("id"),
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Identity not found in state. Either the resource ` +
					`does not support identity or the Terraform version running the test does not support identity. \(must be v1.12\+\)`,
				),
			},
		},
	})
}

func TestExpectIdentityValueMatchesStateAtPath_CheckState_No_Identity(t *testing.T) {
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
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						tfjsonpath.New("id"),
					),
				},
				ExpectError: regexp.MustCompile(`examplecloud_thing.one - Identity not found in state. Either the resource ` +
					`does not support identity or the Terraform version running the test does not support identity. \(must be v1.12\+\)`,
				),
			},
		},
	})
}

func TestExpectIdentityValueMatchesStateAtPath_CheckState_String_Matches(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentityDifferentPaths(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"examplecloud_thing.one",
						tfjsonpath.New("identity_id"),
						tfjsonpath.New("state_id"),
					),
				},
			},
		},
	})
}

func TestExpectIdentityValueMatchesStateAtPath_CheckState_String_DoesntMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithMismatchedResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"examplecloud_thing.one",
						tfjsonpath.New("id"),
						tfjsonpath.New("id"),
					),
				},
				ExpectError: regexp.MustCompile(`expected identity \(examplecloud_thing.one.id\) and state value \(examplecloud_thing.one.id\) to match, but they differ: identity value: id-123, state value: 321-di`),
			},
		},
	})
}

func TestExpectIdentityValueMatchesStateAtPath_CheckState_List(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithResourceIdentityDifferentPaths(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"examplecloud_thing.one",
						tfjsonpath.New("identity_list_of_numbers"),
						tfjsonpath.New("state_list_of_numbers"),
					),
				},
			},
		},
	})
}

func TestExpectIdentityValueMatchesStateAtPath_CheckState_List_DoesntMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProviderWithMismatchedResourceIdentity(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"examplecloud_thing.one",
						tfjsonpath.New("list_of_numbers"),
						tfjsonpath.New("list_of_numbers"),
					),
				},
				ExpectError: regexp.MustCompile(`expected identity \(examplecloud_thing.one.list_of_numbers\) and state value \(examplecloud_thing.one.list_of_numbers\) to match, but they differ: identity value: \[1 2 3 4\], state value: \[4 3 2 1\]`),
			},
		},
	})
}

func examplecloudProviderWithResourceIdentityDifferentPaths() func() (tfprotov6.ProviderServer, error) {
	return providerserver.NewProviderServer(testprovider.Provider{
		Resources: map[string]testprovider.Resource{
			"examplecloud_thing": {
				CreateResponse: &resource.CreateResponse{
					NewState: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name":                  tftypes.String,
								"state_id":              tftypes.String,
								"state_list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"name":     tftypes.NewValue(tftypes.String, "test value"),
							"state_id": tftypes.NewValue(tftypes.String, "id-123"),
							"state_list_of_numbers": tftypes.NewValue(
								tftypes.List{ElementType: tftypes.Number},
								[]tftypes.Value{
									tftypes.NewValue(tftypes.Number, 1),
									tftypes.NewValue(tftypes.Number, 2),
									tftypes.NewValue(tftypes.Number, 3),
									tftypes.NewValue(tftypes.Number, 4),
								},
							),
						},
					),
					NewIdentity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"identity_id":              tftypes.String,
								"identity_list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"identity_id": tftypes.NewValue(tftypes.String, "id-123"),
							"identity_list_of_numbers": tftypes.NewValue(
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
								"name":                  tftypes.String,
								"state_id":              tftypes.String,
								"state_list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"name":     tftypes.NewValue(tftypes.String, "test value"),
							"state_id": tftypes.NewValue(tftypes.String, "id-123"),
							"state_list_of_numbers": tftypes.NewValue(
								tftypes.List{ElementType: tftypes.Number},
								[]tftypes.Value{
									tftypes.NewValue(tftypes.Number, 1),
									tftypes.NewValue(tftypes.Number, 2),
									tftypes.NewValue(tftypes.Number, 3),
									tftypes.NewValue(tftypes.Number, 4),
								},
							),
						},
					),
					NewIdentity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"identity_id":              tftypes.String,
								"identity_list_of_numbers": tftypes.List{ElementType: tftypes.Number},
							},
						},
						map[string]tftypes.Value{
							"identity_id": tftypes.NewValue(tftypes.String, "id-123"),
							"identity_list_of_numbers": tftypes.NewValue(
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
								Name:              "identity_id",
								Type:              tftypes.String,
								RequiredForImport: true,
							},
							{
								Name:              "identity_list_of_numbers",
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
									Name:     "state_id",
									Type:     tftypes.String,
									Computed: true,
								},
								{
									Name:     "state_list_of_numbers",
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
