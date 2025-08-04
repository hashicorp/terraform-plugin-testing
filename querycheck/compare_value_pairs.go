// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource Query Check
var _ QueryCheck = &compareValuePairs{}

type compareValuePairs struct {
	resourceAddressOne string
	attributePathOne   tfjsonpath.Path
	resourceAddressTwo string
	attributePathTwo   tfjsonpath.Path
	comparer           compare.ValueComparer
}

// CheckQuery implements the query check logic.
func (e *compareValuePairs) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	var resourceOne *tfjson.QueryResource
	var resourceTwo *tfjson.QueryResource

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")

		return
	}

	if req.Query.Values == nil {
		resp.Error = fmt.Errorf("query does not contain any query values")

		return
	}

	if req.Query.Values.RootModule == nil {
		resp.Error = fmt.Errorf("query does not contain a root module")

		return
	}

	for _, r := range req.Query.Values.RootModule.Resources {
		if e.resourceAddressOne == r.Address {
			resourceOne = r

			break
		}
	}

	if resourceOne == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddressOne)

		return
	}

	resultOne, err := tfjsonpath.Traverse(resourceOne.AttributeValues, e.attributePathOne)

	if err != nil {
		resp.Error = err

		return
	}

	for _, r := range req.Query.Values.RootModule.Resources {
		if e.resourceAddressTwo == r.Address {
			resourceTwo = r

			break
		}
	}

	if resourceTwo == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddressTwo)

		return
	}

	resultTwo, err := tfjsonpath.Traverse(resourceTwo.AttributeValues, e.attributePathTwo)

	if err != nil {
		resp.Error = err

		return
	}

	err = e.comparer.CompareValues(resultOne, resultTwo)

	if err != nil {
		resp.Error = err
	}
}

// CompareValuePairs returns a query check that compares the value in query for the first given resource address and
// path with the value in query for the second given resource address and path using the supplied value comparer.
func CompareValuePairs(resourceAddressOne string, attributePathOne tfjsonpath.Path, resourceAddressTwo string, attributePathTwo tfjsonpath.Path, comparer compare.ValueComparer) QueryCheck {
	return &compareValuePairs{
		resourceAddressOne: resourceAddressOne,
		attributePathOne:   attributePathOne,
		resourceAddressTwo: resourceAddressTwo,
		attributePathTwo:   attributePathTwo,
		comparer:           comparer,
	}
}
