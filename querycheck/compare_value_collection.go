// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"sort"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource Query Check
var _ QueryCheck = &compareValueCollection{}

type compareValueCollection struct {
	resourceAddressOne string
	collectionPath     []tfjsonpath.Path
	resourceAddressTwo string
	attributePath      tfjsonpath.Path
	comparer           compare.ValueComparer
}

func walkCollectionPath(obj any, paths []tfjsonpath.Path, results []any) ([]any, error) {
	switch t := obj.(type) {
	case []any:
		for _, v := range t {
			if len(paths) == 0 {
				results = append(results, v)
				continue
			}

			x, err := tfjsonpath.Traverse(v, paths[0])

			if err != nil {
				return results, err
			}

			results, err = walkCollectionPath(x, paths[1:], results)

			if err != nil {
				return results, err
			}
		}
	case map[string]any:
		keys := make([]string, 0, len(t))

		for k := range t {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, key := range keys {
			if len(paths) == 0 {
				results = append(results, t[key])
				continue
			}

			x, err := tfjsonpath.Traverse(t, paths[0])

			if err != nil {
				return results, err
			}

			results, err = walkCollectionPath(x, paths[1:], results)

			if err != nil {
				return results, err
			}
		}
	default:
		results = append(results, obj)
	}

	return results, nil
}

// CheckQuery implements the query check logic.
func (e *compareValueCollection) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	var resourceOne *plugintest.QueryResult
	var resourceTwo *plugintest.QueryResult

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")

		return
	}

	if e.resourceAddressOne == req.Query.Address {
		resourceOne = req.Query
	}

	if resourceOne == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddressOne)

		return
	}

	if len(e.collectionPath) == 0 {
		resp.Error = fmt.Errorf("%s - No collection path was provided", e.resourceAddressOne)

		return
	}

	resultOne, err := tfjsonpath.Traverse(resourceOne.ResourceObject, e.collectionPath[0])

	if err != nil {
		resp.Error = err

		return
	}

	// Verify resultOne is a collection.
	switch t := resultOne.(type) {
	case []any, map[string]any:
		// Collection found.
	default:
		var pathStr string

		for _, v := range e.collectionPath {
			pathStr += fmt.Sprintf(".%s", v.String())
		}

		resp.Error = fmt.Errorf("%s%s is not a collection type: %T", e.resourceAddressOne, pathStr, t)

		return
	}

	var results []any

	results, err = walkCollectionPath(resultOne, e.collectionPath[1:], results)

	if err != nil {
		resp.Error = err

		return
	}

	if e.resourceAddressTwo == req.Query.Address {
		resourceTwo = req.Query
	}

	if resourceTwo == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddressTwo)

		return
	}

	resultTwo, err := tfjsonpath.Traverse(resourceTwo.ResourceObject, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	var errs []error

	for _, v := range results {
		switch resultTwo.(type) {
		case []any:
			errs = append(errs, e.comparer.CompareValues([]any{v}, resultTwo))
		default:
			errs = append(errs, e.comparer.CompareValues(v, resultTwo))
		}
	}

	for _, err = range errs {
		if err == nil {
			return
		}
	}

	errMsgs := map[string]struct{}{}

	for _, err = range errs {
		if _, ok := errMsgs[err.Error()]; ok {
			continue
		}

		resp.Error = errors.Join(resp.Error, err)

		errMsgs[err.Error()] = struct{}{}
	}
}

// CompareValueCollection returns a query check that iterates over each element in a collection and compares the value of each element
// with the value of an attribute using the given value comparer.
func CompareValueCollection(resourceAddressOne string, collectionPath []tfjsonpath.Path, resourceAddressTwo string, attributePath tfjsonpath.Path, comparer compare.ValueComparer) QueryCheck {
	return &compareValueCollection{
		resourceAddressOne: resourceAddressOne,
		collectionPath:     collectionPath,
		resourceAddressTwo: resourceAddressTwo,
		attributePath:      attributePath,
		comparer:           comparer,
	}
}
