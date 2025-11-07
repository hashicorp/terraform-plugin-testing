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
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
)

func RunQueryChecks(ctx context.Context, t testing.T, query []tfjson.LogMsg, queryChecks []querycheck.QueryResultCheck) error {
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
		if filterCheck, ok := queryCheck.(querycheck.QueryResultCheckWithFilters); ok {
			var err error
			found, err = runQueryFilters(ctx, filterCheck, found)
			if err != nil {
				return err
			}
		}
		resp := querycheck.CheckQueryResponse{}
		queryCheck.CheckQuery(ctx, querycheck.CheckQueryRequest{
			Query:        found,
			QuerySummary: &summary,
		}, &resp)

		result = append(result, resp.Error)
	}

	return errors.Join(result...)
}

func runQueryFilters(ctx context.Context, filterCheck querycheck.QueryResultCheckWithFilters, queryResults []tfjson.ListResourceFoundData) ([]tfjson.ListResourceFoundData, error) {
	filters := filterCheck.QueryFilters(ctx)
	filteredResults := make([]tfjson.ListResourceFoundData, 0)

	for _, result := range queryResults {
		keepResult := false

		for _, filter := range filters {

			resp := queryfilter.FilterQueryResponse{}
			filter.Filter(ctx, queryfilter.FilterQueryRequest{QueryItem: result}, &resp)

			if resp.Include {
				keepResult = true
			}

			if resp.Error != nil {
				return nil, resp.Error
			}
		}

		if keepResult {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults, nil
}
