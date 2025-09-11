// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package query

import (
	"context"
	"errors"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/querycheck"
)

func RunQueryChecks(ctx context.Context, t testing.T, query *[]tfjson.LogMsg, queryChecks []querycheck.QueryResultCheck) error {
	t.Helper()

	var result []error

	for _, queryCheck := range queryChecks {
		resp := querycheck.CheckQueryResponse{}
		queryCheck.CheckQuery(ctx, querycheck.CheckQueryRequest{Query: query}, &resp)

		result = append(result, resp.Error)
	}

	return errors.Join(result...)
}
