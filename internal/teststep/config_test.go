// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
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
		"raw": {
			configRequest: ConfigurationRequest{
				Raw: Pointer("raw"),
			},
		},
		"directory-raw": {
			configRequest: ConfigurationRequest{
				Directory: Pointer("directory"),
				Raw:       Pointer("raw"),
			},
			expectedError: `both "directory" and "raw" are populated, only one configuration option is allowed`,
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

func TestConfiguration_HasProviderBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config   configuration
		expected bool
	}{
		"no-config": {
			config:   configuration{},
			expected: false,
		},
		"provider-meta-attribute": {
			config: configuration{
				raw: `
resource "test_test" "test" {
  provider = test.test
}
`,
			},
			expected: false,
		},
		"provider-object-attribute": {
			config: configuration{
				raw: `
resource "test_test" "test" {
  test = {
	provider = {
	  test = true
	}
  }
}
`,
			},
			expected: false,
		},
		"provider-string-attribute": {
			config: configuration{
				raw: `
resource "test_test" "test" {
  test = {
	provider = "test"
  }
}
`,
			},
			expected: false,
		},
		"provider-block-quoted-with-attributes": {
			config: configuration{
				raw: `
provider "test" {
  test = true
}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			config: configuration{
				raw: `
provider test {
  test = true
}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			config: configuration{
				raw: `
provider "test" {}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			config: configuration{
				raw: `
provider test {}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.config.HasProviderBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}
