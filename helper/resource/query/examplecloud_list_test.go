package query_test

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/list"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

func examplecloudListResource() testprovider.ListResource {
	return testprovider.ListResource{
		IncludeResource: true,
		SchemaResponse: &list.SchemaResponse{
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
		ListResultsStream: &list.ListResultsStream{
			Results: func(push func(list.ListResult) bool) {
				push(list.ListResult{
					Resource: teststep.Pointer(tftypes.NewValue(
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
					)),
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":       tftypes.String,
								"location": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue1"),
							"location": tftypes.NewValue(tftypes.String, "westeurope"),
						},
					)),
				})
				push(list.ListResult{
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":       tftypes.String,
								"location": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue2"),
							"location": tftypes.NewValue(tftypes.String, "westeurope2"),
						},
					)),
				})
				push(list.ListResult{
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":       tftypes.String,
								"location": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue3"),
							"location": tftypes.NewValue(tftypes.String, "westeurope3"),
						},
					)),
				})
			},
		},
	}
}
