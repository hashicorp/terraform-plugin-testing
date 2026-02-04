// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

func examplecloudDataSource() testprovider.DataSource {
	return testprovider.DataSource{
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
						ComputedStringAttribute("id"),
					},
				},
			},
		},
	}
}

func examplecloudResource() testprovider.Resource {
	return testprovider.Resource{
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
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
				},
			)),
		},
		ReadResponse: &resource.ReadResponse{
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
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
				},
			)),
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
			Identity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
				},
			)),
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						ComputedStringAttribute("id"),
						RequiredStringAttribute("location"),
						RequiredStringAttribute("name"),
					},
				},
			},
		},
		IdentitySchemaResponse: &resource.IdentitySchemaResponse{
			Schema: &tfprotov6.ResourceIdentitySchema{
				Version: 1,
				IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
					{
						Name:              "id",
						Type:              tftypes.String,
						RequiredForImport: true,
					},
				},
			},
		},
	}
}

// examplecloudZone is a test resource that mimics a DNS zone resource.
func examplecloudZone() testprovider.Resource {
	value := tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":   tftypes.String,
				"name": tftypes.String,
			},
		},
		map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "5381dd14-6d75-4f32-9096-47f5500b1507"),
			"name": tftypes.NewValue(tftypes.String, "example.net"),
		},
	)

	return testprovider.Resource{
		CreateResponse: &resource.CreateResponse{
			NewState: value,
		},
		ReadResponse: &resource.ReadResponse{
			NewState: value,
		},
		ImportStateResponse: &resource.ImportStateResponse{
			State: value,
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						ComputedStringAttribute("id"),
						RequiredStringAttribute("name"),
					},
				},
			},
		},
	}
}

// examplecloudZoneRecord is a test resource that mimics a DNS zone record resource.
// It models a resource dependency; specifically, it depends on a DNS zone ID and will
// plan a replacement if the zone ID changes.
func examplecloudZoneRecord() testprovider.Resource {
	value := tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":      tftypes.String,
				"zone_id": tftypes.String,
				"name":    tftypes.String,
			},
		},
		map[string]tftypes.Value{
			"id":      tftypes.NewValue(tftypes.String, "f00911be-e188-433d-9ccd-d0393a9f5d05"),
			"zone_id": tftypes.NewValue(tftypes.String, "5381dd14-6d75-4f32-9096-47f5500b1507"),
			"name":    tftypes.NewValue(tftypes.String, "www"),
		},
	)

	return testprovider.Resource{
		CreateResponse: &resource.CreateResponse{
			NewState: value,
		},
		PlanChangeFunc: func(ctx context.Context, req resource.PlanChangeRequest, resp *resource.PlanChangeResponse) {
			resp.RequiresReplace = []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("zone_id"),
			}
		},
		ReadResponse: &resource.ReadResponse{
			NewState: value,
		},
		ImportStateResponse: &resource.ImportStateResponse{
			State: value,
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						ComputedStringAttribute("id"),
						RequiredStringAttribute("zone_id"),
						RequiredStringAttribute("name"),
					},
				},
			},
		},
	}
}

func examplecloudResourceWithEveryIdentitySchemaType() testprovider.Resource {
	return testprovider.Resource{
		CreateResponse: &resource.CreateResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"hostname":      tftypes.String,
						"cabinet":       tftypes.String,
						"unit":          tftypes.Number,
						"active":        tftypes.Bool,
						"tags":          tftypes.List{ElementType: tftypes.String},
						"magic_numbers": tftypes.List{ElementType: tftypes.Number},
						"beep_boop":     tftypes.List{ElementType: tftypes.Bool},
					},
				},
				map[string]tftypes.Value{
					"hostname": tftypes.NewValue(tftypes.String, "mail.example.net"),
					"cabinet":  tftypes.NewValue(tftypes.String, "A1"),
					"unit":     tftypes.NewValue(tftypes.Number, 14),
					"active":   tftypes.NewValue(tftypes.Bool, true),
					"tags": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
							tftypes.NewValue(tftypes.String, "storage"),
							tftypes.NewValue(tftypes.String, "fast")}),
					"magic_numbers": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{
							tftypes.NewValue(tftypes.Number, 5),
							tftypes.NewValue(tftypes.Number, 2)}),
					"beep_boop": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Bool}, []tftypes.Value{
							tftypes.NewValue(tftypes.Bool, false),
							tftypes.NewValue(tftypes.Bool, true),
						}),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"cabinet":       tftypes.String,
						"unit":          tftypes.Number,
						"active":        tftypes.Bool,
						"tags":          tftypes.List{ElementType: tftypes.String},
						"magic_numbers": tftypes.List{ElementType: tftypes.Number},
						"beep_boop":     tftypes.List{ElementType: tftypes.Bool},
					},
				},
				map[string]tftypes.Value{
					"cabinet": tftypes.NewValue(tftypes.String, "A1"),
					"unit":    tftypes.NewValue(tftypes.Number, 14),
					"active":  tftypes.NewValue(tftypes.Bool, true),
					"tags": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
							tftypes.NewValue(tftypes.String, "storage"),
							tftypes.NewValue(tftypes.String, "fast"),
						}),
					"magic_numbers": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{
							tftypes.NewValue(tftypes.Number, 5),
							tftypes.NewValue(tftypes.Number, 2)}),
					"beep_boop": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Bool}, []tftypes.Value{
							tftypes.NewValue(tftypes.Bool, false),
							tftypes.NewValue(tftypes.Bool, true),
						}),
				},
			)),
		},
		ReadResponse: &resource.ReadResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"hostname":      tftypes.String,
						"cabinet":       tftypes.String,
						"unit":          tftypes.Number,
						"active":        tftypes.Bool,
						"tags":          tftypes.List{ElementType: tftypes.String},
						"magic_numbers": tftypes.List{ElementType: tftypes.Number},
						"beep_boop":     tftypes.List{ElementType: tftypes.Bool},
					},
				},
				map[string]tftypes.Value{
					"hostname": tftypes.NewValue(tftypes.String, "mail.example.net"),
					"cabinet":  tftypes.NewValue(tftypes.String, "A1"),
					"unit":     tftypes.NewValue(tftypes.Number, 14),
					"active":   tftypes.NewValue(tftypes.Bool, true),
					"tags": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
							tftypes.NewValue(tftypes.String, "storage"),
							tftypes.NewValue(tftypes.String, "fast")}),
					"magic_numbers": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{
							tftypes.NewValue(tftypes.Number, 5),
							tftypes.NewValue(tftypes.Number, 2)}),
					"beep_boop": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Bool}, []tftypes.Value{
							tftypes.NewValue(tftypes.Bool, false),
							tftypes.NewValue(tftypes.Bool, true),
						}),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"cabinet":       tftypes.String,
						"unit":          tftypes.Number,
						"active":        tftypes.Bool,
						"tags":          tftypes.List{ElementType: tftypes.String},
						"magic_numbers": tftypes.List{ElementType: tftypes.Number},
						"beep_boop":     tftypes.List{ElementType: tftypes.Bool},
					},
				},
				map[string]tftypes.Value{
					"cabinet": tftypes.NewValue(tftypes.String, "A1"),
					"unit":    tftypes.NewValue(tftypes.Number, 14),
					"active":  tftypes.NewValue(tftypes.Bool, true),
					"tags": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
							tftypes.NewValue(tftypes.String, "storage"),
							tftypes.NewValue(tftypes.String, "fast")}),
					"magic_numbers": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{
							tftypes.NewValue(tftypes.Number, 5),
							tftypes.NewValue(tftypes.Number, 2)}),
					"beep_boop": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Bool}, []tftypes.Value{
							tftypes.NewValue(tftypes.Bool, false),
							tftypes.NewValue(tftypes.Bool, true),
						}),
				},
			)),
		},
		ImportStateResponse: &resource.ImportStateResponse{
			State: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"hostname":      tftypes.String,
						"cabinet":       tftypes.String,
						"unit":          tftypes.Number,
						"active":        tftypes.Bool,
						"tags":          tftypes.List{ElementType: tftypes.String},
						"magic_numbers": tftypes.List{ElementType: tftypes.Number},
						"beep_boop":     tftypes.List{ElementType: tftypes.Bool},
					},
				},
				map[string]tftypes.Value{
					"hostname": tftypes.NewValue(tftypes.String, "mail.example.net"),
					"cabinet":  tftypes.NewValue(tftypes.String, "A1"),
					"unit":     tftypes.NewValue(tftypes.Number, 14),
					"active":   tftypes.NewValue(tftypes.Bool, true),
					"tags": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
							tftypes.NewValue(tftypes.String, "storage"),
							tftypes.NewValue(tftypes.String, "fast")}),
					"magic_numbers": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{
							tftypes.NewValue(tftypes.Number, 5),
							tftypes.NewValue(tftypes.Number, 2)}),
					"beep_boop": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Bool}, []tftypes.Value{
							tftypes.NewValue(tftypes.Bool, false),
							tftypes.NewValue(tftypes.Bool, true),
						}),
				},
			),
			Identity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"cabinet":       tftypes.String,
						"unit":          tftypes.Number,
						"active":        tftypes.Bool,
						"tags":          tftypes.List{ElementType: tftypes.String},
						"magic_numbers": tftypes.List{ElementType: tftypes.Number},
						"beep_boop":     tftypes.List{ElementType: tftypes.Bool},
					},
				},
				map[string]tftypes.Value{
					"cabinet": tftypes.NewValue(tftypes.String, "A1"),
					"unit":    tftypes.NewValue(tftypes.Number, 14),
					"active":  tftypes.NewValue(tftypes.Bool, true),
					"tags": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
							tftypes.NewValue(tftypes.String, "storage"),
							tftypes.NewValue(tftypes.String, "fast")}),
					"magic_numbers": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{
							tftypes.NewValue(tftypes.Number, 5),
							tftypes.NewValue(tftypes.Number, 2)}),
					"beep_boop": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Bool}, []tftypes.Value{
							tftypes.NewValue(tftypes.Bool, false),
							tftypes.NewValue(tftypes.Bool, true),
						}),
				},
			)),
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						ComputedStringAttribute("hostname"),
						RequiredStringAttribute("cabinet"),
						RequiredNumberAttribute("unit"),
						RequiredBoolAttribute("active"),
						RequiredListAttribute("tags", tftypes.String),
						OptionalComputedListAttribute("magic_numbers", tftypes.Number),
					},
				},
			},
		},
		IdentitySchemaResponse: &resource.IdentitySchemaResponse{
			Schema: &tfprotov6.ResourceIdentitySchema{
				Version: 1,
				IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
					{
						Name:              "cabinet",
						Type:              tftypes.String,
						RequiredForImport: true,
					},
					{
						Name:              "unit",
						Type:              tftypes.Number,
						OptionalForImport: true,
					},
					{
						Name:              "active",
						Type:              tftypes.Bool,
						OptionalForImport: true,
					},
					{
						Name: "tags",
						Type: tftypes.List{
							ElementType: tftypes.String,
						},
						OptionalForImport: true,
					},
					{
						Name: "magic_numbers",
						Type: tftypes.List{
							ElementType: tftypes.Number,
						},
						OptionalForImport: true,
					},
				},
			},
		},
	}
}

func examplecloudResourceWithNullIdentityAttr() testprovider.Resource {
	return testprovider.Resource{
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
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":                        tftypes.String,
						"value_we_dont_always_need": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id":                        tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
					"value_we_dont_always_need": tftypes.NewValue(tftypes.String, nil),
				},
			)),
		},
		ReadResponse: &resource.ReadResponse{
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
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":                        tftypes.String,
						"value_we_dont_always_need": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id":                        tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
					"value_we_dont_always_need": tftypes.NewValue(tftypes.String, nil),
				},
			)),
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
			Identity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":                        tftypes.String,
						"value_we_dont_always_need": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id":                        tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
					"value_we_dont_always_need": tftypes.NewValue(tftypes.String, nil),
				},
			)),
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						ComputedStringAttribute("id"),
						RequiredStringAttribute("location"),
						RequiredStringAttribute("name"),
					},
				},
			},
		},
		IdentitySchemaResponse: &resource.IdentitySchemaResponse{
			Schema: &tfprotov6.ResourceIdentitySchema{
				Version: 1,
				IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
					{
						Name:              "id",
						Type:              tftypes.String,
						RequiredForImport: true,
					},
					{
						Name:              "value_we_dont_always_need",
						Type:              tftypes.String,
						OptionalForImport: true,
					},
				},
			},
		},
	}
}

// This example resource, on update plans, will plan a different identity to test that
// our testing framework assertions catch an identity that differs after import/refresh.
func examplecloudResourceWithChangingIdentity() testprovider.Resource {
	exampleCloudResource := examplecloudResource()

	exampleCloudResource.PlanChangeFunc = func(ctx context.Context, req resource.PlanChangeRequest, resp *resource.PlanChangeResponse) {
		// Only on update
		if !req.PriorState.IsNull() && !req.ProposedNewState.IsNull() {
			resp.PlannedIdentity = teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "easteurope/someothervalue"),
				},
			))
		}
	}

	return exampleCloudResource
}
