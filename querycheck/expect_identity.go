// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ QueryResultCheck = expectIdentity{}

type expectIdentity struct {
	resourceAddress string
	check           map[string]knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectIdentity) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	var foundIdentities []map[string]json.RawMessage

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")
		return
	}

	for _, v := range *req.Query {
		switch i := v.(type) {
		case tfjson.ListResourceFoundMessage:
			prefix := "list."
			if strings.TrimPrefix(i.ListResourceFound.Address, prefix) == e.resourceAddress {
				foundIdentities = append(foundIdentities, i.ListResourceFound.Identity)
			}
		default:
			continue
		}
	}

	if len(foundIdentities) == 0 {
		resp.Error = fmt.Errorf("%s - Identity not found in query. Either the resource does not support query or the Terraform version running the test does not support query. (must be v1.14+)", e.resourceAddress)

		return
	}

	var err error

	for _, resultIdentity := range foundIdentities {
		for attribute := range e.check {
			var val any
			var ok bool
			var unmarshalledVal any

			if val, ok = resultIdentity[attribute]; !ok {
				resp.Error = fmt.Errorf("%s - expected attribute %q not in actual identity object", e.resourceAddress, attribute)
				return
			}

			rawMessage, ok := val.(json.RawMessage)
			if !ok {
				resp.Error = fmt.Errorf("%s - expected json.RawMessage but got %T", e.resourceAddress, val)
				return
			}
			err = json.Unmarshal(rawMessage, &unmarshed)

			if err != nil {
				resp.Error = fmt.Errorf("%s - Error decoding message type: %s", e.resourceAddress, err)
				return
			}

			if err = e.check[attribute].CheckValue(unmarshalledVal); err != nil {
				errCollection = append(errCollection, fmt.Errorf("%s - %q identity attribute: %s\n", e.resourceAddress, e.check, err))
			}
		}
	}

	if !found {
		resp.Error = fmt.Errorf("%s - %q identity attribute: %s", e.resourceAddress, e.check, err)
	}

	return
}

// ExpectIdentity returns a query check that asserts that the identity at the given resource matches a known object, where each
// map key represents an identity attribute name. The identity in query must exactly match the given object and any missing/extra
// attributes will raise a diagnostic.
//
// This query check can only be used with managed resources that support resource identity. Resource identity is only supported in Terraform v1.12+
func ExpectIdentity(resourceAddress string, identity map[string]knownvalue.Check) QueryResultCheck {
	return expectIdentity{
		resourceAddress: resourceAddress,
		check:           identity,
	}
}
