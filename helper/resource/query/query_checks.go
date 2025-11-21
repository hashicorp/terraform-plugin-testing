// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package query

import (
	"context"
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/hack"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
)

func RunQueryChecks(ctx context.Context, t hack.BaseT, query []tfjson.LogMsg, queryChecks []querycheck.QueryResultCheck) error {
	t.Helper()

	var result []error

	if query == nil {
		result = append(result, fmt.Errorf("no query results found"))
	}

	found := make([]tfjson.ListResourceFoundData, 0)
	summary := tfjson.ListCompleteData{}

	for _, msg := range query {
		switch v := msg.(type) {
		case tfjson.ListResourceFoundMessage:
			found = append(found, v.ListResourceFound)
		case tfjson.ListCompleteMessage:
			summary = v.ListComplete
			// TODO diagnostics and errors?
		default:
			continue
		}
	}

	for _, queryCheck := range queryChecks {
		resp := querycheck.CheckQueryResponse{}
		queryCheck.CheckQuery(ctx, querycheck.CheckQueryRequest{
			Query:        found,
			QuerySummary: &summary,
		}, &resp)

		result = append(result, resp.Error)
	}

	return errors.Join(result...)
}
