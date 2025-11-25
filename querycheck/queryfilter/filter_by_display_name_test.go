// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package queryfilter_test

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
)

func TestByDisplayName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		displayName   knownvalue.StringCheck
		queryItem     tfjson.ListResourceFoundData
		expectInclude bool
		expectedError error
	}{
		"nil-query-result-exact": {
			displayName:   knownvalue.StringExact("test"),
			expectInclude: false,
		},
		"nil-query-result-regex": {
			displayName:   knownvalue.StringRegexp(regexp.MustCompile("display")),
			expectInclude: false,
		},
		"empty-display-name-exact": {
			displayName:   knownvalue.StringExact(""),
			expectInclude: true,
		},
		"empty-display-name-regexp": {
			displayName:   knownvalue.StringRegexp(regexp.MustCompile("")),
			expectInclude: true,
		},
		"included-exact": {
			displayName: knownvalue.StringExact("test"),
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
			},
			expectInclude: true,
		},
		"included-regex": {
			displayName: knownvalue.StringRegexp(regexp.MustCompile("test")),
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
			},
			expectInclude: true,
		},
		"not-included-exact": {
			displayName: knownvalue.StringExact("test"),
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "testsss",
			},
			expectInclude: false,
		},
		"not-included-regex": {
			displayName: knownvalue.StringRegexp(regexp.MustCompile("invalid")),
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "testsss",
			},
			expectInclude: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := queryfilter.FilterQueryRequest{QueryItem: testCase.queryItem}

			resp := &queryfilter.FilterQueryResponse{}

			queryfilter.ByDisplayName(testCase.displayName).Filter(t.Context(), req, resp)

			if testCase.expectInclude != resp.Include {
				t.Fatalf("expected included: %t, but got %t", testCase.expectInclude, resp.Include)
			}

			if testCase.expectedError == nil && resp.Error != nil {
				t.Errorf("unexpected error %s", resp.Error)
			}

			if testCase.expectedError != nil && resp.Error == nil {
				t.Errorf("expected error but got none")
			}

			if diff := cmp.Diff(resp.Error, testCase.expectedError); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
