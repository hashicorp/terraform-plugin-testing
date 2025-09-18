// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
)

// QueryResultCheck defines an interface for implementing test logic to apply an assertion against a collection of found
// resources that were returned by a query. It returns an error if the query results do not match what is expected.
type QueryResultCheck interface {
	// CheckQuery should perform the query check.
	CheckQuery(context.Context, CheckQueryRequest, *CheckQueryResponse)
}

// CheckQueryRequest is a request for an invoke of the CheckQuery function.
type CheckQueryRequest struct {
	// Query represents the parsed log messages relating to found resources returned by the `terraform query -json` command.
	Query []tfjson.ListResourceFoundData

	// QuerySummary contains a summary of the completed query operation
	QuerySummary *tfjson.ListCompleteData
}

// CheckQueryResponse is a response to an invoke of the CheckQuery function.
type CheckQueryResponse struct {
	// Error is used to report the failure of a query check assertion and is combined with other QueryResultCheck errors
	// to be reported as a test failure.
	Error error
}
