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

func TestTest_TestStep_ImportBlockVerify(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// import blocks are only available in v1.5.0 and later
			tfversion.SkipBelow(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_bucket": exampleCloudBucketResource(t),
				},
			}),
		},
		Steps: []TestStep{
			{
				Config: `
				resource "examplecloud_bucket" "storage" {
					bucket = "test-bucket"
					description = "A bucket for testing."
				}`,
			},
			{
				ImportState:                          true,
				ImportStateKind:                      ImportBlockWithResourceIdentity,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "bucket", // upgrade to resource identity
				ResourceName:                         "examplecloud_bucket.storage",
			},
		},
	})
}

func exampleCloudBucketResource(t *testing.T) testprovider.Resource {
	t.Helper()

	return testprovider.Resource{
		CreateResponse: &resource.CreateResponse{
			NewResourceIdentityData: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"region": tftypes.String,
						"bucket": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"region": tftypes.NewValue(tftypes.String, "test-region"),
					"bucket": tftypes.NewValue(tftypes.String, "test-bucket"),
				},
			),
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"bucket":      tftypes.String,
						"description": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"bucket":      tftypes.NewValue(tftypes.String, "test-bucket"),
					"description": tftypes.NewValue(tftypes.String, "A bucket for testing."),
				},
			),
		},
		ImportStateResponse: &resource.ImportStateResponse{
			State: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"bucket":      tftypes.String,
						"description": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"bucket":      tftypes.NewValue(tftypes.String, "test-bucket"),
					"description": tftypes.NewValue(tftypes.String, "A bucket for testing."),
				},
			),
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "bucket",
							Type:     tftypes.String,
							Required: true,
						},
						{
							Name:     "description",
							Type:     tftypes.String,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
