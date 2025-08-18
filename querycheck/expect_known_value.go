// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource Query Check
var _ QueryCheck = expectKnownValue{}

type expectKnownValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	knownValue      knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectKnownValue) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	var resource *plugintest.QueryResult

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")

		return
	}

	if len(req.Query.Address) == 0 {
		resp.Error = fmt.Errorf("query does not contain any address values")

		return
	}

	if e.resourceAddress == req.Query.Address {
		resource = req.Query
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(resource.ResourceObject, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for attribute at path: %s.%s, err: %s", e.resourceAddress, e.attributePath.String(), err)

		return
	}
}

// ExpectKnownValue returns a query check that asserts that the specified attribute at the given resource
// has a known type and value.
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) QueryCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
