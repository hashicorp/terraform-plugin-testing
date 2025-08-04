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
var _ QueryCheck = &compareValue{}

type compareValue struct {
	resourceAddresses []string
	attributePaths    []tfjsonpath.Path
	queryValues       []any
	comparer          compare.ValueComparer
}

func (e *compareValue) AddQueryValue(resourceAddress string, attributePath tfjsonpath.Path) QueryCheck {
	e.resourceAddresses = append(e.resourceAddresses, resourceAddress)
	e.attributePaths = append(e.attributePaths, attributePath)

	return e
}

// CheckQuery implements the query check logic.
func (e *compareValue) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	var resource *tfjson.QueryResource

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

	// All calls to AddQueryValue occur before any TestStep is run, populating the resourceAddresses
	// and attributePaths slices. The queryValues slice is populated during execution of each TestStep.
	// Each call to CheckQuery happens sequentially during each TestStep.
	// The currentIndex is reflective of the current query value being checked.
	currentIndex := len(e.queryValues)

	if len(e.resourceAddresses) <= currentIndex {
		resp.Error = fmt.Errorf("resource addresses index out of bounds: %d", currentIndex)

		return
	}

	resourceAddress := e.resourceAddresses[currentIndex]

	for _, r := range req.Query.Values.RootModule.Resources {
		if resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in query", resourceAddress)

		return
	}

	if len(e.attributePaths) <= currentIndex {
		resp.Error = fmt.Errorf("attribute paths index out of bounds: %d", currentIndex)

		return
	}

	attributePath := e.attributePaths[currentIndex]

	result, err := tfjsonpath.Traverse(resource.AttributeValues, attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	e.queryValues = append(e.queryValues, result)

	err = e.comparer.CompareValues(e.queryValues...)

	if err != nil {
		resp.Error = err
	}
}

// CompareValue returns a query check that compares values retrieved from query using the
// supplied value comparer.
func CompareValue(comparer compare.ValueComparer) *compareValue {
	return &compareValue{
		comparer: comparer,
	}
}
