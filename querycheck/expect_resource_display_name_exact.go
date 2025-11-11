// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
)

var _ QueryResultCheck = expectResourceDisplayNameExact{}
var _ QueryResultCheckWithFilters = expectResourceDisplayNameExact{}

type expectResourceDisplayNameExact struct {
	listResourceAddress string
	filter              queryfilter.QueryFilter
	displayName         string
}

func (e expectResourceDisplayNameExact) QueryFilters(ctx context.Context) []queryfilter.QueryFilter {
	if e.filter == nil {
		return []queryfilter.QueryFilter{}
	}

	return []queryfilter.QueryFilter{
		e.filter,
	}
}

func (e expectResourceDisplayNameExact) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	listRes := make([]tfjson.ListResourceFoundData, 0)
	for _, result := range req.Query {
		if strings.TrimPrefix(result.Address, "list.") == e.listResourceAddress {
			listRes = append(listRes, result)
		}
	}

	if len(listRes) == 0 {
		resp.Error = fmt.Errorf("%s - no query results found after filtering", e.listResourceAddress)
		return
	}

	if len(listRes) > 1 {
		resp.Error = fmt.Errorf("%s - more than 1 query result found after filtering", e.listResourceAddress)
		return
	}
	res := listRes[0]
	if strings.EqualFold(e.displayName, res.DisplayName) {
		return
	}

	resp.Error = fmt.Errorf("expected to find resource with display name %q in results but resource was not found", e.displayName)
}

// ExpectResourceDisplayNameExact returns a query check that asserts that a resource with a given display name exists within the returned results of the query.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectResourceDisplayNameExact(listResourceAddress string, filter queryfilter.QueryFilter, displayName string) QueryResultCheck {
	return expectResourceDisplayNameExact{
		listResourceAddress: listResourceAddress,
		filter:              filter,
		displayName:         displayName,
	}
}
