// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"context"
	"iter"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type ListResource interface {
	Schema(context.Context, SchemaRequest, *SchemaResponse)
	List(context.Context, ListRequest, *ListResultsStream)
}

type ListRequest struct {
}

type ListResultsStream struct {
	Results iter.Seq[ListResult]
}

type ListResult struct {
}

type ValidateListConfigResponse struct {
}

type SchemaRequest struct{}

type SchemaResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	Schema      *tfprotov6.Schema
}
