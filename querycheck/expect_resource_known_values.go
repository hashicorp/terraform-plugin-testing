// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ QueryResultCheck = expectResourceKnownValues{}
var _ QueryResultCheckWithFilters = expectResourceKnownValues{}

type expectResourceKnownValues struct {
	listResourceAddress string
	filter              queryfilter.QueryFilter
	knownValues         []KnownValueCheck
}

func (e expectResourceKnownValues) QueryFilters(ctx context.Context) []queryfilter.QueryFilter {
	if e.filter == nil {
		return []queryfilter.QueryFilter{}
	}

	return []queryfilter.QueryFilter{
		e.filter,
	}
}

func (e expectResourceKnownValues) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	for _, res := range req.Query {
		var diags []error

		if e.listResourceAddress == strings.TrimPrefix(res.Address, "list.") {
			if res.ResourceObject == nil {
				resp.Error = fmt.Errorf("%s - no resource object was returned, ensure `include_resource` has been set to `true` in the list resource config`", e.listResourceAddress)
				return
			}

			for _, c := range e.knownValues {
				resource, err := tfjsonpath.Traverse(res.ResourceObject, c.Path)
				if err != nil {
					resp.Error = err
					return
				}

				if err := c.KnownValue.CheckValue(resource); err != nil {
					diags = append(diags, fmt.Errorf("error checking value for attribute at path: %s for resource with identity %s, err: %s", c.Path.String(), e.filter, err))
				}

				if diags == nil {
					return
				}
			}
		}

		if diags != nil {
			var diagsStr string
			for _, diag := range diags {
				diagsStr += diag.Error() + "; "
			}
			resp.Error = fmt.Errorf("the following errors were found while checking values: %s", diagsStr)
			return
		}
	}

	resp.Error = fmt.Errorf("%s - the resource %s was not found", e.listResourceAddress, e.filter)
}

// ExpectResourceKnownValues returns a query check that asserts the specified attribute values are present for a given resource object
// returned by a list query. The resource object can only be identified by providing the list resource address as well as
// a query filter.
//
// This query check can only be used with managed resources that support resource identity and query. Query is only supported in Terraform v1.14+
func ExpectResourceKnownValues(listResourceAddress string, filter queryfilter.QueryFilter, knownValues []KnownValueCheck) QueryResultCheck {
	return expectResourceKnownValues{
		listResourceAddress: listResourceAddress,
		filter:              filter,
		knownValues:         knownValues,
	}
}
