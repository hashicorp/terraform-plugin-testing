// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_ImportCommmandWithId_SkipDataSourceState(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				DataSources: map[string]testprovider.DataSource{
					"examplecloud_thing": {
						ReadResponse: &datasource.ReadResponse{
							State: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id": tftypes.NewValue(tftypes.String, "datasource-test"),
								},
							),
						},
						SchemaResponse: &datasource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "id",
											Type:     tftypes.String,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
				Resources: map[string]testprovider.Resource{
					"examplecloud_thing": {
						CreateResponse: &resource.CreateResponse{
							NewState: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id": tftypes.NewValue(tftypes.String, "resource-test"),
								},
							),
						},
						ImportStateResponse: &resource.ImportStateResponse{
							State: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id": tftypes.NewValue(tftypes.String, "resource-test"),
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
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `
					data "examplecloud_thing" "test" {}
					resource "examplecloud_thing" "test" {}
				`,
			},
			{
				ResourceName: "examplecloud_thing.test",
				ImportState:  true,
				ImportStateCheck: func(is []*terraform.InstanceState) error {
					if len(is) > 1 {
						return fmt.Errorf("expected 1 state, got: %d", len(is))
					}

					return nil
				},
			},
		},
	})
}

func Test_ImportCommandWithId_ImportStateVerify(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_thing": {
						CreateResponse: &resource.CreateResponse{
							NewState: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":    tftypes.String,
										"other": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":    tftypes.NewValue(tftypes.String, "resource-test"),
									"other": tftypes.NewValue(tftypes.String, "testvalue"),
								},
							),
						},
						ImportStateResponse: &resource.ImportStateResponse{
							State: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":    tftypes.String,
										"other": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":    tftypes.NewValue(tftypes.String, "resource-test"),
									"other": tftypes.NewValue(tftypes.String, "testvalue"),
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
											Name:     "other",
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
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "test" {}`,
			},
			{
				ResourceName:      "examplecloud_thing.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func Test_TestStep_ImportStateVerifyIgnore(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_thing": {
						CreateResponse: &resource.CreateResponse{
							NewState: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":          tftypes.String,
										"create_only": tftypes.String,
										"read_only":   tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":          tftypes.NewValue(tftypes.String, "resource-test"),
									"create_only": tftypes.NewValue(tftypes.String, "testvalue"),
									"read_only":   tftypes.NewValue(tftypes.String, "testvalue"),
								},
							),
						},
						ImportStateResponse: &resource.ImportStateResponse{
							State: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":          tftypes.String,
										"create_only": tftypes.String,
										"read_only":   tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":          tftypes.NewValue(tftypes.String, "resource-test"),
									"create_only": tftypes.NewValue(tftypes.String, nil), // intentional
									"read_only":   tftypes.NewValue(tftypes.String, "testvalue"),
								},
							),
						},
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "create_only",
											Type:     tftypes.String,
											Computed: true,
										},
										{
											Name:     "id",
											Type:     tftypes.String,
											Computed: true,
										},
										{
											Name:     "read_only",
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
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_thing" "test" {}`,
			},
			{
				ResourceName:            "examplecloud_thing.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_only"},
			},
		},
	})
}

func Test_TestStep_ExpectError_ImportState(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						CreateResponse: &resource.CreateResponse{
							NewState: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id": tftypes.NewValue(tftypes.String, "resource-test"),
								},
							),
						},
						ImportStateResponse: &resource.ImportStateResponse{
							Diagnostics: []*tfprotov6.Diagnostic{
								{
									Severity: tfprotov6.DiagnosticSeverityError,
									Summary:  "Invalid Import ID",
									Detail:   "Diagnostic details",
								},
							},
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
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config:        `resource "test_resource" "test" {}`,
				ImportStateId: "invalid time string",
				ResourceName:  "test_resource.test",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Error: Invalid Import ID`),
			},
		},
	})
}
