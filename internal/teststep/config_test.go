// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfigurationRequest_Validate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configRequest ConfigurationRequest
		expectedError string
	}{
		"directory": {
			configRequest: ConfigurationRequest{
				Directory: Pointer("directory"),
			},
		},
		"file": {
			configRequest: ConfigurationRequest{
				Raw: Pointer("file"),
			},
		},
		"raw": {
			configRequest: ConfigurationRequest{
				Raw: Pointer("raw"),
			},
		},
		"directory-file": {
			configRequest: ConfigurationRequest{
				Directory: Pointer("directory"),
				File:      Pointer("file"),
			},
			expectedError: `directory and file are populated, only one of "directory", "file", or "raw"  is allowed`,
		},
		"directory-raw": {
			configRequest: ConfigurationRequest{
				Directory: Pointer("directory"),
				Raw:       Pointer("raw"),
			},
			expectedError: `directory and raw are populated, only one of "directory", "file", or "raw"  is allowed`,
		},
		"directory-file-raw": {
			configRequest: ConfigurationRequest{
				Directory: Pointer("directory"),
				File:      Pointer("file"),
				Raw:       Pointer("raw"),
			},
			expectedError: `directory, file and raw are populated, only one of "directory", "file", or "raw"  is allowed`,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.configRequest.Validate()

			if testCase.expectedError == "" && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != "" && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != "" && err != nil {
				if diff := cmp.Diff(err.Error(), testCase.expectedError); diff != "" {
					t.Errorf("expected error %s, got error %s", testCase.expectedError, err)
				}
			}
		})
	}
}
