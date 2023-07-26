// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"testing"
)

func TestConfiguration_HasProviderBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configRaw configurationString
		expected  bool
	}{
		"no-config": {
			expected: false,
		},
		"provider-meta-attribute": {
			configRaw: configurationString{
				raw: `
resource "test_test" "test" {
 provider = test.test
}
`,
			},
			expected: false,
		},
		"provider-object-attribute": {
			configRaw: configurationString{
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
			configRaw: configurationString{
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
			configRaw: configurationString{
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
			configRaw: configurationString{
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
			configRaw: configurationString{
				raw: `
provider "test" {}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configRaw: configurationString{
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

			got, err := testCase.configRaw.HasProviderBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}
