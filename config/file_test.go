// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestTestStepConfigFunc_Exec_File(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testStepConfigFunc    config.TestStepConfigFunc
		testStepConfigRequest config.TestStepConfigRequest
		expected              string
	}{
		"static_file": {
			testStepConfigFunc: config.StaticFile("name_of_file"),
			expected:           "name_of_file",
		},
		"test_name_file": {
			testStepConfigFunc: config.TestNameFile("test.tf"),
			testStepConfigRequest: config.TestStepConfigRequest{
				TestName: "TestTestStepConfigFunc_Exec",
			},
			expected: "testdata/TestTestStepConfigFunc_Exec/test.tf",
		},
		"test_step_file": {
			testStepConfigFunc: config.TestStepFile("test.tf"),
			testStepConfigRequest: config.TestStepConfigRequest{
				StepNumber: 1,
				TestName:   "TestTestStepConfigFunc_Exec",
			},
			expected: "testdata/TestTestStepConfigFunc_Exec/1/test.tf",
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
