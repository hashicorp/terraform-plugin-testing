package terraform

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/internal/configs/hcl2shim"
)

func TestUnknownValueWalk(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    any
		expected bool
	}{
		"primitive": {
			value:    42,
			expected: false,
		},
		"primitive computed": {
			value:    hcl2shim.UnknownVariableValue,
			expected: true,
		},
		"list": {
			value: []any{
				"foo",
				hcl2shim.UnknownVariableValue,
			},
			expected: true,
		},
		"nested list": {
			value: []any{
				"foo",
				[]any{hcl2shim.UnknownVariableValue},
			},
			expected: true,
		},
		"map": {
			value: map[string]any{
				"testkey1": "foo",
				"testkey2": hcl2shim.UnknownVariableValue,
			},
			expected: true,
		},
		"nested map": {
			value: map[string]any{
				"testkey1": "foo",
				"testkey2": map[string]any{"testkey": hcl2shim.UnknownVariableValue},
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := unknownValueWalk(reflect.ValueOf(testCase.value))

			if got != testCase.expected {
				t.Errorf("expected: %t, got: %t", testCase.expected, got)
			}
		})
	}
}
