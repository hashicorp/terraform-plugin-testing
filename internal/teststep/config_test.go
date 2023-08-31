// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestPrepareConfigurationRequest_Exec(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		prepareConfigRequest PrepareConfigurationRequest
		expected             ConfigurationRequest
	}{
		"directory": {
			prepareConfigRequest: PrepareConfigurationRequest{
				Directory: func(request config.TestStepConfigRequest) string { return "directory" },
			},
			expected: ConfigurationRequest{
				Directory: Pointer("directory"),
				File:      Pointer(""),
				Raw:       Pointer(""),
			},
		},
		"file": {
			prepareConfigRequest: PrepareConfigurationRequest{
				File: func(request config.TestStepConfigRequest) string { return "file" },
			},
			expected: ConfigurationRequest{
				Directory: Pointer(""),
				File:      Pointer("file"),
				Raw:       Pointer(""),
			},
		},
		"raw": {
			prepareConfigRequest: PrepareConfigurationRequest{
				Raw: "str",
			},
			expected: ConfigurationRequest{
				Directory: Pointer(""),
				File:      Pointer(""),
				Raw:       Pointer("str"),
			},
		},
		"directory-file": {
			prepareConfigRequest: PrepareConfigurationRequest{
				Directory: func(request config.TestStepConfigRequest) string { return "directory" },
				File:      func(request config.TestStepConfigRequest) string { return "file" },
			},
			expected: ConfigurationRequest{
				Directory: Pointer("directory"),
				File:      Pointer("file"),
				Raw:       Pointer(""),
			},
		},
		"directory-raw": {
			prepareConfigRequest: PrepareConfigurationRequest{
				Directory: func(request config.TestStepConfigRequest) string { return "directory" },
				Raw:       "str",
			},
			expected: ConfigurationRequest{
				Directory: Pointer("directory"),
				File:      Pointer(""),
				Raw:       Pointer("str"),
			},
		},
		"file-raw": {
			prepareConfigRequest: PrepareConfigurationRequest{
				File: func(request config.TestStepConfigRequest) string { return "file" },
				Raw:  "str",
			},
			expected: ConfigurationRequest{
				Directory: Pointer(""),
				File:      Pointer("file"),
				Raw:       Pointer("str"),
			},
		},
		"directory-file-raw": {
			prepareConfigRequest: PrepareConfigurationRequest{
				Directory: func(request config.TestStepConfigRequest) string { return "directory" },
				File:      func(request config.TestStepConfigRequest) string { return "file" },
				Raw:       "str",
			},
			expected: ConfigurationRequest{
				Directory: Pointer("directory"),
				File:      Pointer("file"),
				Raw:       Pointer("str"),
			},
		},
	}

	comparer := cmp.Comparer(func(x, y ConfigurationRequest) bool {
		if x.Directory != nil && y.Directory == nil {
			return false
		}

		if x.Directory == nil && y.Directory != nil {
			return false
		}

		if *x.Directory != *y.Directory {
			return false
		}

		if x.File != nil && y.File == nil {
			return false
		}

		if x.File == nil && y.File != nil {
			return false
		}

		if *x.File != *y.File {
			return false
		}

		if x.Raw != nil && y.Raw == nil {
			return false
		}

		if x.Raw == nil && y.Raw != nil {
			return false
		}

		if *x.Raw != *y.Raw {
			return false
		}

		return true
	})

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.prepareConfigRequest.Exec()

			if diff := cmp.Diff(testCase.expected, got, comparer); diff != "" {
				t.Errorf("expected %+v, got %+v", testCase.expected, got)
			}

		})
	}
}

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

func TestConfiguration(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configRequest ConfigurationRequest
		expected      Config
	}{
		"directory": {
			configRequest: ConfigurationRequest{
				Directory: Pointer("directory"),
			},
			expected: configurationDirectory{
				directory: "directory",
			},
		},
		"file": {
			configRequest: ConfigurationRequest{
				File: Pointer("file"),
			},
			expected: configurationFile{
				file: "file",
			},
		},
		"raw": {
			configRequest: ConfigurationRequest{
				Raw: Pointer("str"),
			},
			expected: configurationString{
				raw: "str",
			},
		},
	}

	allowUnexported := cmp.AllowUnexported(
		configurationDirectory{},
		configurationFile{},
		configurationString{},
	)

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := Configuration(testCase.configRequest)

			if diff := cmp.Diff(testCase.expected, got, allowUnexported); diff != "" {
				t.Errorf("expected %+v, got %+v", testCase.expected, got)
			}
		})
	}
}
