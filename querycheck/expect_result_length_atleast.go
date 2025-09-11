// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

var _ QueryResultCheck = expectLengthAtLeast{}

type expectLengthAtLeast struct {
	resourceAddress string
	check           int
}

// CheckQuery implements the query check logic.
func (e expectLengthAtLeast) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")
		return
	}

	for _, v := range *req.Query {
		switch i := v.(type) {
		case tfjson.ListCompleteMessage:
			prefix := "list."

			if strings.TrimPrefix(i.ListComplete.Address, prefix) == e.resourceAddress {
				if i.ListComplete.Total < e.check {
					resp.Error = fmt.Errorf("Query result of at least length %v - expected but got %v.", e.check, i.ListComplete.Total)
					return
				} else {
					return
				}
			}
		default:
			continue
		}
	}

	resp.Error = fmt.Errorf("%s - Address not found in query result.", e.resourceAddress)

	return
}

// ExpectLengthAtLeast returns a query check that asserts that the length of the query result is at least the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectLengthAtLeast(resourceAddress string, length int) QueryResultCheck {
	return expectLengthAtLeast{
		resourceAddress: resourceAddress,
		check:           length,
	}
}
