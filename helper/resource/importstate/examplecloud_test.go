package importstate_test

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

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
	}
}
