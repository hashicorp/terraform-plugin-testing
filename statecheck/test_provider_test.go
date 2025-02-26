// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

type aTestProvider struct {
}

func (t aTestProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
}

func (t aTestProvider) DataSourcesMap() map[string]datasource.DataSource {
	return nil
}

func (t aTestProvider) ResourcesMap() map[string]resource.Resource {
	return map[string]resource.Resource{
		"test_resource": &aTestResource{},
	}
}

func (t aTestProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	schema := tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name: "planet",
					Type: tftypes.String,
				},
			},
		},
	}
	response.Schema = &schema
}

func (t aTestProvider) Stop(ctx context.Context, request provider.StopRequest, response *provider.StopResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestProvider) ValidateConfig(ctx context.Context, request provider.ValidateConfigRequest, response *provider.ValidateConfigResponse) {
}
