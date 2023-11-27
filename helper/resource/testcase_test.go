// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/84
func TestTestCase_NoDataSourceIdRequirement(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		Steps: []TestStep{
			{
				Check: ComposeAggregateTestCheckFunc(
					TestCheckNoResourceAttr("data.test_datasource.test", "id"),
					TestCheckResourceAttr("data.test_datasource.test", "not_id", "test"),
				),
				Config: `data "test_datasource" "test" {}`,
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": providerserver.NewProviderServer(testprovider.Provider{
						DataSources: map[string]testprovider.DataSource{
							"test_datasource": {
								ReadResponse: &datasource.ReadResponse{
									State: tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"not_id": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"not_id": tftypes.NewValue(tftypes.String, "test"),
										},
									),
								},
								SchemaResponse: &datasource.SchemaResponse{
									Schema: &tfprotov6.Schema{
										Block: &tfprotov6.SchemaBlock{
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:     "not_id",
													Type:     tftypes.String,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					}),
				},
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/84
func TestTestCase_NoResourceIdRequirement(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		Steps: []TestStep{
			{
				Check: ComposeAggregateTestCheckFunc(
					TestCheckNoResourceAttr("test_resource.test", "id"),
					TestCheckResourceAttr("test_resource.test", "not_id", "test"),
				),
				Config: `resource "test_resource" "test" {}`,
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": providerserver.NewProviderServer(testprovider.Provider{
						Resources: map[string]testprovider.Resource{
							"test_resource": {
								CreateResponse: &resource.CreateResponse{
									NewState: tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"not_id": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"not_id": tftypes.NewValue(tftypes.String, "test"),
										},
									),
								},
								SchemaResponse: &resource.SchemaResponse{
									Schema: &tfprotov6.Schema{
										Block: &tfprotov6.SchemaBlock{
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:     "not_id",
													Type:     tftypes.String,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					}),
				},
			},
		},
	})
}
