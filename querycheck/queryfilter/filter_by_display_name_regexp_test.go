package queryfilter_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
)

func TestByDisplayNameRegexp(t *testing.T) {
	testCases := map[string]struct {
		regexp        *regexp.Regexp
		queryItem     tfjson.ListResourceFoundData
		expectInclude bool
		expectedError error
	}{
		"nil-query-result": {
			regexp:        regexp.MustCompile("display"),
			expectInclude: false,
		},
		"empty-regexp": {
			regexp:        regexp.MustCompile(""),
			expectInclude: true,
		},
		"included": {
			regexp: regexp.MustCompile("test"),
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
			},
			expectInclude: true,
		},
		"not-included": {
			regexp: regexp.MustCompile("invalid"),
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

			queryfilter.ByDisplayNameRegexp(testCase.regexp).Filter(context.TODO(), req, resp)

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
