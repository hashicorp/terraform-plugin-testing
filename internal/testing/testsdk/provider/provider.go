// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

type Provider interface {
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)
	DataSourcesMap() map[string]datasource.DataSource
	ResourcesMap() map[string]resource.Resource
	Schema(context.Context, SchemaRequest, *SchemaResponse)
	Stop(context.Context, StopRequest, *StopResponse)
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}

type ConfigureRequest struct {
	Config tftypes.Value
}

type ConfigureResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}

type SchemaRequest struct{}

type SchemaResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	Schema      *tfprotov6.Schema
}

type StopRequest struct{}

type StopResponse struct {
	Error error
}

type ValidateConfigRequest struct {
	Config tftypes.Value
}

type ValidateConfigResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}
