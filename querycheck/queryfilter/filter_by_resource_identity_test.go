package queryfilter_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
)

func TestByResourceIdentity(t *testing.T) {
	testCases := map[string]struct {
		identity      map[string]knownvalue.Check
		queryItem     tfjson.ListResourceFoundData
		expectInclude bool
		expectedError error
	}{
		"nil-query-result": {
			identity: map[string]knownvalue.Check{
				"id": knownvalue.StringExact("id-123"),
			},
			expectInclude: false,
		},
		"nil-identity": {
			expectInclude: true,
		},
		"included": {
			identity: map[string]knownvalue.Check{
				"id": knownvalue.StringExact("id-123"),
				"list_of_numbers": knownvalue.ListExact(
					[]knownvalue.Check{
						knownvalue.Int64Exact(1),
						knownvalue.Int64Exact(2),
						knownvalue.Int64Exact(3),
						knownvalue.Int64Exact(4),
					},
				),
			},
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
				Identity: map[string]any{
					"id": "id-123",
					"list_of_numbers": []any{
						json.Number("1"),
						json.Number("2"),
						json.Number("3"),
						json.Number("4"),
					},
				},
			},
			expectInclude: true,
		},
		"not-included-nonexistent-attribute": {
			identity: map[string]knownvalue.Check{
				"id":               knownvalue.StringExact("id-123"),
				"nonexistent_attr": knownvalue.StringExact("hello"),
				"list_of_numbers": knownvalue.ListExact(
					[]knownvalue.Check{
						knownvalue.Int64Exact(1),
						knownvalue.Int64Exact(2),
						knownvalue.Int64Exact(3),
						knownvalue.Int64Exact(4),
					},
				),
			},
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
				Identity: map[string]any{
					"id": "id-123",
					"list_of_numbers": []any{
						json.Number("1"),
						json.Number("2"),
						json.Number("3"),
						json.Number("4"),
					},
				},
			},
			expectInclude: false,
		},
		"not-included-incorrect-string": {
			identity: map[string]knownvalue.Check{
				"id": knownvalue.StringExact("id-123"),
				"list_of_numbers": knownvalue.ListExact(
					[]knownvalue.Check{
						knownvalue.Int64Exact(1),
						knownvalue.Int64Exact(2),
						knownvalue.Int64Exact(3),
						knownvalue.Int64Exact(4),
					},
				),
			},
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
				Identity: map[string]any{
					"id": "incorrect",
					"list_of_numbers": []any{
						json.Number("1"),
						json.Number("2"),
						json.Number("3"),
						json.Number("4"),
					},
				},
			},
			expectInclude: false,
		},
		"not-included-incorrect-list-item": {
			identity: map[string]knownvalue.Check{
				"id": knownvalue.StringExact("id-123"),
				"list_of_numbers": knownvalue.ListExact(
					[]knownvalue.Check{
						knownvalue.Int64Exact(1),
						knownvalue.Int64Exact(2),
						knownvalue.Int64Exact(3),
						knownvalue.Int64Exact(4),
					},
				),
			},
			queryItem: tfjson.ListResourceFoundData{
				DisplayName: "test",
				Identity: map[string]any{
					"id": "id-123",
					"list_of_numbers": []any{
						json.Number("1"),
						json.Number("2"),
						json.Number("333"),
						json.Number("4"),
					},
				},
			},
			expectInclude: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := queryfilter.FilterQueryRequest{QueryItem: testCase.queryItem}

			resp := &queryfilter.FilterQueryResponse{}

			queryfilter.ByResourceIdentity(testCase.identity).Filter(context.TODO(), req, resp)

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
