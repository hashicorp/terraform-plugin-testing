// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package query

import (
	"context"
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/querycheck"
)

func RunQueryChecks(ctx context.Context, t testing.T, query *[]tfjson.LogMsg, queryChecks []querycheck.QueryResultCheck) error {
	t.Helper()

	var result []error

	if query == nil || len(*query) == 0 {
		result = append(result, fmt.Errorf("No query results found"))
	}

	found := make([]tfjson.ListResourceFoundData, 0)
	complete := tfjson.ListCompleteData{}

	for _, msg := range *query {
		switch v := msg.(type) {
		case tfjson.ListResourceFoundMessage:
			found = append(found, v.ListResourceFound)
		case tfjson.ListCompleteMessage:
			complete = v.ListComplete
			// TODO diagnostics and errors?
		default:
			// ignore other message types
		}
	}

	// TODO check diagnostics in LogMsg to see if there are any errors we can return here?
	var err error
	if len(found) == 0 {
		return fmt.Errorf("no resources found by query: %+v", err)
	}

	for _, queryCheck := range queryChecks {
		resp := querycheck.CheckQueryResponse{}
		queryCheck.CheckQuery(ctx, querycheck.CheckQueryRequest{
			Query:          &found,
			CompletedQuery: &complete,
		}, &resp)

		result = append(result, resp.Error)
	}

	return errors.Join(result...)
}
