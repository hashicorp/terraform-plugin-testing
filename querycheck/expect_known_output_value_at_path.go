// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource Query Check
var _ QueryCheck = expectKnownOutputValueAtPath{}

type expectKnownOutputValueAtPath struct {
	outputAddress string
	outputPath    tfjsonpath.Path
	knownValue    knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectKnownOutputValueAtPath) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	var output *tfjson.QueryOutput

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")

		return
	}

	if req.Query.Values == nil {
		resp.Error = fmt.Errorf("query does not contain any query values")

		return
	}

	for address, oc := range req.Query.Values.Outputs {
		if e.outputAddress == address {
			output = oc

			break
		}
	}

	if output == nil {
		resp.Error = fmt.Errorf("%s - Output not found in query", e.outputAddress)

		return
	}

	result, err := tfjsonpath.Traverse(output.Value, e.outputPath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for output at path: %s.%s, err: %s", e.outputAddress, e.outputPath.String(), err)

		return
	}
}

// ExpectKnownOutputValueAtPath returns a query check that asserts that the specified output at the given path
// has a known type and value.
func ExpectKnownOutputValueAtPath(outputAddress string, outputPath tfjsonpath.Path, knownValue knownvalue.Check) QueryCheck {
	return expectKnownOutputValueAtPath{
		outputAddress: outputAddress,
		outputPath:    outputPath,
		knownValue:    knownValue,
	}
}
