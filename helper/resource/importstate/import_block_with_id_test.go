// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_TestStep_ImportBlockId(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
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
				ImportStateKind:   r.ImportBlockWithID,
				ImportStateVerify: true,
			},
		},
	})
}

func TestTest_TestStep_ImportBlockId_ExpectError(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "examplecloud_container" "test" {
					location = "westeurope"
					name     = "somevalue"
				}`,
			},
			{
				Config: `
				resource "examplecloud_container" "test" {
					location = "eastus"
					name     = "somevalue"
				}`,
				ResourceName:      "examplecloud_container.test",
				ImportState:       true,
				ImportStateKind:   r.ImportBlockWithID,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`importing resource examplecloud_container.test: expected a no-op resource action, got "update" action with plan(.?)`),
			},
		},
	})
}

func TestTest_TestStep_ImportBlockId_FailWhenPlannableImportIsNotSupported(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
			tfversion.SkipAbove(tfversion.Version1_4_0), // ImportBlockWithId requires Terraform 1.5.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "examplecloud_container" "test" {
					location = "westeurope"
					name     = "somevalue"
				}`,
			},
			{
				Config: `
				resource "examplecloud_container" "test" {
					location = "eastus"
					name     = "somevalue"
				}`,
				ResourceName:      "examplecloud_container.test",
				ImportState:       true,
				ImportStateKind:   r.ImportBlockWithID,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`Terraform 1.5.0`),
			},
		},
	})
}

func TestTest_TestStep_ImportBlockId_SkipDataSourceState(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later

		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				DataSources: map[string]testprovider.DataSource{
					"examplecloud_thing": examplecloudDataSource(),
				},
				Resources: map[string]testprovider.Resource{
					"examplecloud_thing": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `
					data "examplecloud_thing" "test" {}
					resource "examplecloud_thing" "test" {
						name = "somevalue"
						location = "westeurope"
					}
				`,
			},
			{
				ResourceName:    "examplecloud_thing.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithID,
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

// These tests currently pass but only because the `getState` function which is used on the imported resource
// to do the state comparison doesn't return an error if there is no state in the working directory
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

			I also need to omit the `password` in the import config, otherwise the value in the config is used when importing the
		    with an import block and the test ends up passing regardless of whether `ImportStateVerifyIgnore` has been specified or not
	*/
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
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
		Steps: []r.TestStep{
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
				ImportStateKind:         r.ImportBlockWithID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestTest_TestStep_ImportBlockId_ImportStateVerifyIgnore(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
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
		Steps: []r.TestStep{
			{
				Config: `resource "examplecloud_container" "test" {}`,
			},
			{
				ResourceName:            "examplecloud_container.test",
				ImportState:             true,
				ImportStateKind:         r.ImportBlockWithID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
