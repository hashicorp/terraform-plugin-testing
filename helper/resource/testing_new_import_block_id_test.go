// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestTest_TestStep_ImportBlockId_SkipDataSourceState(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ProtoV6ProviderFactories
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
		Steps: []TestStep{
			{
				Config: `
					data "examplecloud_thing" "test" {}
					resource "examplecloud_thing" "test" {}
				`,
			},
			{
				ResourceName:    "examplecloud_thing.test",
				ImportState:     true,
				ImportStateKind: ImportBlockWithId,
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

func TestTest_TestStep_ImportBlockId_ImportStateVerifyIgnore_Real_Example(t *testing.T) {
	/*
		This test tries to imitate a real world example of behaviour we often see in the AzureRM provider which requires
		the use of `ImportStateVerifyIgnore` when testing the import of a resource using the import command.

		A sensitive field e.g. a password can be supplied on create but isn't returned in the API response on a subsequent
		read, resulting in a different value for password in the two states.

		In the AzureRM provider this is usually handled one of two ways, both requiring `ImportStateVerifyIgnore` to make
		the test pass:

		1. Property doesn't get set in the read
			* in pluginSDK at create the config gets written to state because that's what we're expecting
			* the subsequent read updates the values to create a post-apply diff and update computed values
		 	* since we don't do anything to the property in the read the imported resource's state has the password missing
		      compared to the created resource's state

		2. We retrieve the value from config and set that into state
			* the config isn't available at import time using only the import command (I think?) so there is nothing to
		      retrieve and set into state when importing

		For this test to pass I needed to add a `PlanChangeFunc` to the resource to set the id to a known value in the plan - see comment in the `PlanChangeFunc`

		I also need to omit the `password` in the import config, otherwise the value in the config is used when importing the resource and the test
		ends up passing regardless of whether `ImportStateVerifyIgnore` has been specified or not

		Ultimately it looks like:
			* Terraform is saying there's a bug in the provider? (see comment in `PlanChangeFunc`)
			* The import behaviour using a block vs. the command appears to differ
	*/
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
										"name":     tftypes.String,
										"password": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":       tftypes.NewValue(tftypes.String, "sometestid"),
									"name":     tftypes.NewValue(tftypes.String, "somename"),
									"password": tftypes.NewValue(tftypes.String, "somevalue"),
								},
							),
						},
						ImportStateResponse: &resource.ImportStateResponse{
							State: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":       tftypes.String,
										"name":     tftypes.String,
										"password": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":       tftypes.NewValue(tftypes.String, "sometestid"),
									"name":     tftypes.NewValue(tftypes.String, "somename"),
									"password": tftypes.NewValue(tftypes.String, nil), // this simulates an absent property when importing
								},
							),
						},
						PlanChangeFunc: func(ctx context.Context, request resource.PlanChangeRequest, response *resource.PlanChangeResponse) {
							/*
								Returning a nil for another attribute to simulate a situation where `ImportStateVerifyIgnore`
								should be used results in the error below from Terraform

								Error: Provider returned invalid result object after apply

								        After the apply operation, the provider still indicated an unknown value for
								        examplecloud_container.test.id. All values must be known after apply, so this
								        is always a bug in the provider and should be reported in the provider's own
								        repository. Terraform will still save the other known object values in the
								        state.

								Modifying the plan to set the id to a known value appears to be the only way to
								circumvent this behaviour, the cause of which I don't fully understand.

								This doesn't seem great, because this gets applied to all Plans that happen in this
								test - so we're modifying plans in steps that we might not want to.
							*/

							objVal := map[string]tftypes.Value{}

							if !response.PlannedState.IsNull() {
								_ = response.PlannedState.As(&objVal)
								objVal["id"] = tftypes.NewValue(tftypes.String, "sometestid")
							}
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
											Name:     "name",
											Type:     tftypes.String,
											Required: true,
										},
										{
											Name:     "password",
											Type:     tftypes.String,
											Optional: true,
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
					name     = "somename"
  					password = "somevalue"
				}`,
			},
			{
				Config: `
				terraform {
  					required_providers {
    					examplecloud = {
      						source = "registry.terraform.io/hashicorp/examplecloud"
						}
					}
				}

				resource "examplecloud_container" "test" {
					name = "somename"
				}`,
				ResourceName:            "examplecloud_container.test",
				ImportState:             true,
				ImportStateKind:         ImportBlockWithId,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestTest_TestStep_ImportBlockId_ImportStateVerifyIgnore(t *testing.T) {
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
										"name":     tftypes.String,
										"password": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":       tftypes.NewValue(tftypes.String, "sometestid"),
									"name":     tftypes.NewValue(tftypes.String, "somename"),
									"password": tftypes.NewValue(tftypes.String, "somevalue"),
								},
							),
						},
						ImportStateResponse: &resource.ImportStateResponse{
							State: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id":       tftypes.String,
										"name":     tftypes.String,
										"password": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id":       tftypes.NewValue(tftypes.String, "sometestid"),
									"name":     tftypes.NewValue(tftypes.String, "somename"),
									"password": tftypes.NewValue(tftypes.String, nil), // this simulates an absent property when importing
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
											Name:     "name",
											Type:     tftypes.String,
											Computed: true,
										},
										{
											Name:     "password",
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
		Steps: []TestStep{
			{
				Config: `resource "examplecloud_container" "test" {}`,
			},
			{
				ResourceName:            "examplecloud_container.test",
				ImportState:             true,
				ImportStateKind:         ImportBlockWithId,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
