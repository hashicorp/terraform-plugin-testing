package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

var _ provider.Provider = Provider{}

// Provider is a declarative provider implementation for unit testing in this
// Go module.
type Provider struct {
	ConfigureResponse      *provider.ConfigureResponse
	DataSources            map[string]DataSource
	Resources              map[string]Resource
	SchemaResponse         *provider.SchemaResponse
	StopResponse           *provider.StopResponse
	ValidateConfigResponse *provider.ValidateConfigResponse
}

func (p Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if p.ConfigureResponse != nil {
		resp.Diagnostics = p.ConfigureResponse.Diagnostics
	}
}

func (p Provider) DataSourcesMap() map[string]datasource.DataSource {
	datasources := make(map[string]datasource.DataSource, len(p.DataSources))

	for typeName, d := range p.DataSources {
		datasources[typeName] = d
	}

	return datasources
}

func (p Provider) ResourcesMap() map[string]resource.Resource {
	resources := make(map[string]resource.Resource, len(p.Resources))

	for typeName, d := range p.Resources {
		resources[typeName] = d
	}

	return resources
}

func (p Provider) Stop(ctx context.Context, req provider.StopRequest, resp *provider.StopResponse) {
	if p.StopResponse != nil {
		resp.Error = p.StopResponse.Error
	}
}

func (p Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	if p.SchemaResponse != nil {
		resp.Diagnostics = p.SchemaResponse.Diagnostics
		resp.Schema = p.SchemaResponse.Schema
	}

	resp.Schema = &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{},
	}
}

func (p Provider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	if p.ValidateConfigResponse != nil {
		resp.Diagnostics = p.ValidateConfigResponse.Diagnostics
	}
}
