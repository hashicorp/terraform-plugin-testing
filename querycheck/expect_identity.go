// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ QueryCheck = expectIdentity{}

type expectIdentity struct {
	resourceAddress string
	check           map[string]knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectIdentity) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	//var resource map[string]json.RawMessage
	var foundIdentities []map[string]json.RawMessage
	var found bool

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")
		return
	}

	// iterate through query messages and find the identity message
	// if we find the identity remove message from slice and if we are at the end then it wasn't found :(

	for _, v := range *req.Query {
		switch idk := v.(type) {
		case tfjson.ListResourceFoundMessage:
			prefix := "list."
			if strings.TrimPrefix(idk.ListResourceFound.Address, prefix) == e.resourceAddress {
				foundIdentities = append(foundIdentities, idk.ListResourceFound.Identity)
			}
		default:
			continue
		}
	}

	//if resource == nil {
	//	resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddress)
	//
	//	return
	//}

	if len(foundIdentities) == 0 {
		resp.Error = fmt.Errorf("%s - Identity not found in query. Either the resource does not support identity or the Terraform version running the test does not support identity. (must be v1.14+)", e.resourceAddress)

		return
	}

	//if len(resource.IdentityValues) != len(e.identity) {
	//deltaMsg := ""
	//if len(resource.IdentityValues) > len(e.identity) {
	//	deltaMsg = createDeltaString(resource.IdentityValues, e.identity, "actual identity has extra attribute(s): ")
	//} else {
	//	deltaMsg = createDeltaString(e.identity, resource.IdentityValues, "actual identity is missing attribute(s): ")
	//}

	//resp.Error = fmt.Errorf("%s - Expected %d attribute(s) in the actual identity object, got %d attribute(s): %s", e.resourceAddress, len(e.identity), len(resource.IdentityValues), deltaMsg)
	//}

	var err error

	for _, resultIdentity := range foundIdentities {
		for attribute := range e.check {
			var val any
			var ok bool
			if val, ok = resultIdentity[attribute]; !ok {
				resp.Error = fmt.Errorf("%s - expected attribute %q not in actual identity object", e.resourceAddress, attribute)
				return
			}

			var unmarshed string

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

			if err = e.check[attribute].CheckValue(unmarshed); err == nil {
				found = true
				return
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
func ExpectIdentity(resourceAddress string, identity map[string]knownvalue.Check) QueryCheck {
	return expectIdentity{
		resourceAddress: resourceAddress,
		check:           identity,
	}
}

// createDeltaString prints the map keys that are present in mapA and not present in mapB
func createDeltaString[T any, V any](mapA map[string]T, mapB map[string]V, msgPrefix string) string {
	deltaMsg := ""

	deltaMap := make(map[string]T, len(mapA))
	maps.Copy(deltaMap, mapA)
	for key := range mapB {
		delete(deltaMap, key)
	}

	deltaKeys := slices.Sorted(maps.Keys(deltaMap))

	for i, k := range deltaKeys {
		if i == 0 {
			deltaMsg += msgPrefix
		} else {
			deltaMsg += ", "
		}
		deltaMsg += fmt.Sprintf("%q", k)
	}

	return deltaMsg
}

const (
	MessageListResourceFound tfjson.LogMessageType = "list_resource_found"
)

type ListResourceFoundMessage struct {
	baseLogMessage
	Address        string                     `json:"address"`
	DisplayName    string                     `json:"display_name"`
	Identity       map[string]json.RawMessage `json:"identity"`
	ResourceType   string                     `json:"resource_type"`
	ResourceObject map[string]json.RawMessage `json:"resource_object,omitempty"`
	Config         string                     `json:"config,omitempty"`
	ImportConfig   string                     `json:"import_config,omitempty"`
}

type baseLogMessage struct {
	Lvl  tfjson.LogMessageLevel `json:"@level"`
	Msg  string                 `json:"@message"`
	Time time.Time              `json:"@timestamp"`
}

type msgType struct {
	Type tfjson.LogMessageType `json:"type"`
}

type Result struct {
	ListResourceFoundMessage `json:"list_resource_found"`
}

func unmarshalResult(t tfjson.LogMessageType, b []byte) (*tfjson.ListResourceFoundData, error) {
	v := tfjson.ListResourceFoundData{}
	switch t {
	case MessageListResourceFound:
		return &v, json.Unmarshal(b, &v)
	}

	return nil, nil
}
