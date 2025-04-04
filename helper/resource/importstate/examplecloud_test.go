// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
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
						ComputedStringAttribute("id"),
						RequiredStringAttribute("location"),
						RequiredStringAttribute("name"),
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
