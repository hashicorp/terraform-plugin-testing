// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestTest_TestStep_ImportBlockId(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ProtoV6ProviderFactories
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": {
						CreateResponse: &resource.CreateResponse{
							NewState: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":       tftypes.String,
										"location": tftypes.String,
										"name":     tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
									"location": tftypes.NewValue(tftypes.String, "westeurope"),
									"name":     tftypes.NewValue(tftypes.String, "somevalue"),
								},
							),
						},
						ImportStateResponse: &resource.ImportStateResponse{
							State: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":       tftypes.String,
										"location": tftypes.String,
										"name":     tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
									"location": tftypes.NewValue(tftypes.String, "westeurope"),
									"name":     tftypes.NewValue(tftypes.String, "somevalue"),
								},
							),
						},
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "id",
											Type:     tftypes.String,
											Computed: true,
										},
										{
											Name:     "location",
											Type:     tftypes.String,
											Required: true,
										},
										{
											Name:     "name",
											Type:     tftypes.String,
											Required: true,
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []TestStep{
			{
				Config: `
				resource "examplecloud_container" "test" {
					location = "westeurope"
					name     = "somevalue"
				}`,
			},
			{
				ResourceName:      "examplecloud_container.test",
				ImportState:       true,
				ImportStateKind:   ImportBlockWithId,
				ImportStateVerify: true,
			},
		},
	})
}
