// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/list"
)

var _ list.ListResource = ListResource{}

type ListResource struct {
	SchemaResponse             *list.SchemaResponse
	ListResultsStream          *list.ListResultsStream
	ValidateListConfigResponse *list.ValidateListConfigResponse
}

func (r ListResource) Schema(ctx context.Context, req list.SchemaRequest, resp *list.SchemaResponse) {
	if r.SchemaResponse != nil {
		resp.Diagnostics = r.SchemaResponse.Diagnostics
		resp.Schema = r.SchemaResponse.Schema
	}
}
func (r ListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	stream.Results = r.ListResultsStream.Results
}
