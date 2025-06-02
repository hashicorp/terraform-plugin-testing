// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
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

// Example Usage
// -------------
func TestExpectStateContainsIdentityValueAtPath_CheckState_String_Contains(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": examplecloudProvider(),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectStateContainsIdentityValueAtPath(
						"examplecloud_thing.one",
						tfjsonpath.New("resource_group_name"), // identity path
						tfjsonpath.New("parent_resource_id"),  // state path
					),
				},
			},
		},
	})
}

// -------------

func examplecloudProvider() func() (tfprotov6.ProviderServer, error) {
	return providerserver.NewProviderServer(testprovider.Provider{
		Resources: map[string]testprovider.Resource{
			"examplecloud_thing": {
				CreateResponse: &resource.CreateResponse{
					NewState: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name":               tftypes.String,
								"parent_resource_id": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name":               tftypes.NewValue(tftypes.String, "test value"),
							"parent_resource_id": tftypes.NewValue(tftypes.String, "/subscriptionId/exampleSub/resourceGroupName/exampleRG"), // doesn't match identity -> id
						},
					),
					NewIdentity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "exampleRG"),
						},
					)),
				},
				ReadResponse: &resource.ReadResponse{
					NewState: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name":               tftypes.String,
								"parent_resource_id": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name":               tftypes.NewValue(tftypes.String, "test value"),
							"parent_resource_id": tftypes.NewValue(tftypes.String, "/subscriptionId/exampleSub/resourceGroupName/exampleRG"),
						},
					),
					NewIdentity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "exampleRG"),
						},
					)),
				},
				IdentitySchemaResponse: &resource.IdentitySchemaResponse{
					Schema: &tfprotov6.ResourceIdentitySchema{
						IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
							{
								Name:              "resource_group_name",
								Type:              tftypes.String,
								RequiredForImport: true,
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
									Name:     "parent_resource_id",
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
