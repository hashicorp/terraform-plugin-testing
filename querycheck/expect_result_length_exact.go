// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ QueryResultCheck = expectLength{}

type expectLength struct {
	resourceAddress string
	check           knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectLength) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.Query == nil {
		resp.Error = fmt.Errorf("Query is nil")
		return
	}

	for _, v := range *req.Query {
		switch i := v.(type) {
		case tfjson.ListCompleteMessage:
			prefix := "list."
			lengthCheck := json.Number((strconv.Itoa(i.ListComplete.Total)))

			if strings.TrimPrefix(i.ListComplete.Address, prefix) == e.resourceAddress {
				if err := e.check.CheckValue(lengthCheck); err != nil {
					resp.Error = fmt.Errorf("Query result of length %v - expected but got %v.", e.check, i.ListComplete.Total)
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

// ExpectLength returns a query check that asserts that the length of the query result is exactly the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectLength(resourceAddress string, length knownvalue.Check) QueryResultCheck {
	return expectLength{
		resourceAddress: resourceAddress,
		check:           length,
	}
}
