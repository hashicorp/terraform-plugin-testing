// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestTestStepConfigFunc_Exec(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testStepConfigFunc    config.TestStepConfigFunc
		testStepConfigRequest config.TestStepConfigRequest
		expected              string
	}{
		"static_directory": {
			testStepConfigFunc: config.StaticDirectory("name_of_directory"),
			expected:           "name_of_directory",
		},
		"test_name_directory": {
			testStepConfigFunc: config.TestNameDirectory(t),
			expected:           "TestTestStepConfigFunc_Exec",
		},
		"test_step_directory": {
			testStepConfigFunc: config.TestStepDirectory(t),
			testStepConfigRequest: config.TestStepConfigRequest{
				StepNumber: 1,
			},
			expected: "TestTestStepConfigFunc_Exec/1",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.testStepConfigFunc.Exec(testCase.testStepConfigRequest)

			if testCase.expected != got {
				t.Errorf("expected %s, got %s", testCase.expected, got)
			}
		})
	}
}
