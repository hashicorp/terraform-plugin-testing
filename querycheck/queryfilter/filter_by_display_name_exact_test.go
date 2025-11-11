package queryfilter_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
)

func TestByDisplayNameExact(t *testing.T) {
	testCases := map[string]struct {
		displayName   string
		queryItem     tfjson.ListResourceFoundData
		expectInclude bool
		expectedError error
	}{
		"nil-query-result": {
			displayName:   "test",
			expectInclude: false,
		},
		"empty-display-name": {
			displayName:   "",
			expectInclude: true,
		},
		"included": {
			displayName: "test",
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
			},
			expectInclude: true,
		},
		"not-included": {
			displayName: "test",
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

			queryfilter.ByDisplayNameExact(testCase.displayName).Filter(context.TODO(), req, resp)

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
