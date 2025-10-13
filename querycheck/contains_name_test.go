// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/list"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func dessertsThatStartWithPResource() testprovider.Resource {
	return testprovider.Resource{
		CreateResponse: &resource.CreateResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"name": tftypes.NewValue(tftypes.String, "pie"),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"name": tftypes.NewValue(tftypes.String, "pie"),
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
					"name": tftypes.NewValue(tftypes.String, "pie"),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"name": tftypes.NewValue(tftypes.String, "pie"),
				},
			)),
		},
		ImportStateResponse: &resource.ImportStateResponse{
			State: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"name": tftypes.NewValue(tftypes.String, "pie"),
				},
			),
			Identity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"name": tftypes.NewValue(tftypes.String, "pie"),
				},
			)),
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "name",
							Type:     tftypes.String,
							Required: true,
						},
					},
				},
			},
		},
		IdentitySchemaResponse: &resource.IdentitySchemaResponse{
			Schema: &tfprotov6.ResourceIdentitySchema{
				Version: 1,
				IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
					{
						Name:              "name",
						Type:              tftypes.String,
						RequiredForImport: true,
					},
				},
			},
		},
	}
}

func dessertsThatStartWithPListResource() testprovider.ListResource {
	return testprovider.ListResource{
		IncludeResource: true,
		SchemaResponse: &list.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "group",
							Type:     tftypes.String,
							Computed: true,
						},
					},
				},
			},
		},
		ListResultsStream: &list.ListResultsStream{
			Results: func(push func(list.ListResult) bool) {
				push(list.ListResult{
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "pie"),
						},
					)),
					DisplayName: "pie",
				})
				push(list.ListResult{
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, "pudding"),
						},
					)),
					DisplayName: "pudding",
				})
			},
		},
	}
}

func TestContainsResourceWithName(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"dessertcloud": providerserver.NewProviderServer(testprovider.Provider{
				ListResources: map[string]testprovider.ListResource{
					"dessert_letter_p": dessertsThatStartWithPListResource(),
				},
				Resources: map[string]testprovider.Resource{
					"dessert_letter_p": dessertsThatStartWithPResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Query: true,
				Config: `
				provider "dessertcloud" {} 

				list "dessert_letter_p" "test" {
					provider = dessertcloud

					config {
						group = "foo"
 					}
				}

				list "dessert_letter_p" "test2" {
					provider = dessertcloud

					config {
						group = "bar"
					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ContainsResourceWithName("dessert_letter_p.test", "pie"),
					querycheck.ContainsResourceWithName("dessert_letter_p.test", "pudding"),
				},
			},
		},
	})
}

func TestContainsResourceWithName_NotFound(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"dessertcloud": providerserver.NewProviderServer(testprovider.Provider{
				ListResources: map[string]testprovider.ListResource{
					"dessert_letter_p": dessertsThatStartWithPListResource(),
				},
				Resources: map[string]testprovider.Resource{
					"dessert_letter_p": dessertsThatStartWithPResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Query: true,
				Config: `
				provider "dessertcloud" {}

				list "dessert_letter_p" "test" {
					provider = dessertcloud

					config {
						group = "foo"
					}
				}

				list "dessert_letter_p" "test2" {
					provider = dessertcloud

					config {
						group = "bar"
					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ContainsResourceWithName("dessert_letter_p.test", "pavlova"),
				},
				ExpectError: regexp.MustCompile("expected to find resource with display name \"pavlova\" in results but resource was not found"),
			},
		},
	})
}
