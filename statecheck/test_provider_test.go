// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

type aTestProvider struct {
}

func (t aTestProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestProvider) DataSourcesMap() map[string]datasource.DataSource {
	//TODO implement me
	panic("implement me")
}

func (t aTestProvider) ResourcesMap() map[string]resource.Resource {
	//TODO implement me
	panic("implement me")
}

func (t aTestProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestProvider) Stop(ctx context.Context, request provider.StopRequest, response *provider.StopResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestProvider) ValidateConfig(ctx context.Context, request provider.ValidateConfigRequest, response *provider.ValidateConfigResponse) {
	//TODO implement me
	panic("implement me")
}
