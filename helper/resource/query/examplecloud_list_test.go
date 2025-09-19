// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
							Name:     "resource_group_name",
							Type:     tftypes.String,
							Required: true,
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
								"id":                  tftypes.String,
								"location":            tftypes.String,
								"name":                tftypes.String,
								"resource_group_name": tftypes.String,
								"instances":           tftypes.Number,
							},
						},
						map[string]tftypes.Value{
							"id":                  tftypes.NewValue(tftypes.String, "foo/banane"),
							"location":            tftypes.NewValue(tftypes.String, "westeurope"),
							"name":                tftypes.NewValue(tftypes.String, "banane"),
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"instances":           tftypes.NewValue(tftypes.Number, 5),
						},
					)),
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
								"name":                tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"name":                tftypes.NewValue(tftypes.String, "banane"),
						},
					)),
					DisplayName: "banane",
				})
				push(list.ListResult{
					Resource: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":                  tftypes.String,
								"location":            tftypes.String,
								"name":                tftypes.String,
								"resource_group_name": tftypes.String,
								"instances":           tftypes.Number,
							},
						},
						map[string]tftypes.Value{
							"id":                  tftypes.NewValue(tftypes.String, "foo/ananas"),
							"location":            tftypes.NewValue(tftypes.String, "westeurope"),
							"name":                tftypes.NewValue(tftypes.String, "ananas"),
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"instances":           tftypes.NewValue(tftypes.Number, 9000),
						},
					)),
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
								"name":                tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"name":                tftypes.NewValue(tftypes.String, "ananas"),
						},
					)),
					DisplayName: "ananas",
				})
				push(list.ListResult{
					Resource: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":                  tftypes.String,
								"location":            tftypes.String,
								"name":                tftypes.String,
								"resource_group_name": tftypes.String,
								"instances":           tftypes.Number,
							},
						},
						map[string]tftypes.Value{
							"id":                  tftypes.NewValue(tftypes.String, "foo/kiwi"),
							"location":            tftypes.NewValue(tftypes.String, "westeurope"),
							"name":                tftypes.NewValue(tftypes.String, "kiwi"),
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"instances":           tftypes.NewValue(tftypes.Number, 88),
						},
					)),
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
								"name":                tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"name":                tftypes.NewValue(tftypes.String, "kiwi"),
						},
					)),
					DisplayName: "kiwi",
				})
				push(list.ListResult{
					Resource: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":                  tftypes.String,
								"location":            tftypes.String,
								"name":                tftypes.String,
								"resource_group_name": tftypes.String,
								"instances":           tftypes.Number,
							},
						},
						map[string]tftypes.Value{
							"id":                  tftypes.NewValue(tftypes.String, "bar/papaya"),
							"location":            tftypes.NewValue(tftypes.String, "westeurope"),
							"name":                tftypes.NewValue(tftypes.String, "banane"),
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"instances":           tftypes.NewValue(tftypes.Number, 3),
						},
					)),
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
								"name":                tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "bar"),
							"name":                tftypes.NewValue(tftypes.String, "papaya"),
						},
					)),
					DisplayName: "papaya",
				})
				push(list.ListResult{
					Resource: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":                  tftypes.String,
								"location":            tftypes.String,
								"name":                tftypes.String,
								"resource_group_name": tftypes.String,
								"instances":           tftypes.Number,
							},
						},
						map[string]tftypes.Value{
							"id":                  tftypes.NewValue(tftypes.String, "bar/birne"),
							"location":            tftypes.NewValue(tftypes.String, "westeurope"),
							"name":                tftypes.NewValue(tftypes.String, "birne"),
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"instances":           tftypes.NewValue(tftypes.Number, 8564),
						},
					)),
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
								"name":                tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "bar"),
							"name":                tftypes.NewValue(tftypes.String, "birne"),
						},
					)),
					DisplayName: "birne",
				})
				push(list.ListResult{
					Resource: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":                  tftypes.String,
								"location":            tftypes.String,
								"name":                tftypes.String,
								"resource_group_name": tftypes.String,
								"instances":           tftypes.Number,
							},
						},
						map[string]tftypes.Value{
							"id":                  tftypes.NewValue(tftypes.String, "bar/kirsche"),
							"location":            tftypes.NewValue(tftypes.String, "westeurope"),
							"name":                tftypes.NewValue(tftypes.String, "kirsche"),
							"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
							"instances":           tftypes.NewValue(tftypes.Number, 500),
						},
					)),
					Identity: teststep.Pointer(tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"resource_group_name": tftypes.String,
								"name":                tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, "bar"),
							"name":                tftypes.NewValue(tftypes.String, "kirsche"),
						},
					)),
					DisplayName: "kirsche",
				})
			},
		},
	}
}
