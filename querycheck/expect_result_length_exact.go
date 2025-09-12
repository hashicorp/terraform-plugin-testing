// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ QueryResultCheck = expectLength{}

type expectLength struct {
	resourceAddress string
	check           knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectLength) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.CompletedQuery == nil {
		resp.Error = fmt.Errorf("no completed query information available")
		return
	}

	lengthCheck := json.Number(strconv.Itoa(req.CompletedQuery.Total))

	if err := e.check.CheckValue(lengthCheck); err != nil {
		resp.Error = fmt.Errorf("Query result of length %v - expected but got %v.", e.check, req.CompletedQuery.Total)
		return
	}

	return
}

// ExpectLength returns a query check that asserts that the length of the query result is exactly the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectLength(resourceAddress string, length knownvalue.Check) QueryResultCheck {
	return expectLength{
		resourceAddress: resourceAddress,
		check:           length,
	}
}
