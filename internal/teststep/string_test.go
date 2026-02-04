// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		"provider-block-quoted-with-attributes-no-spaces": {
			configRaw: configurationString{
				raw: `
provider"test"{
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
		"provider-block-unquoted-with-attributes-no-trailing-space": {
			configRaw: configurationString{
				raw: `
provider test{
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
		"provider-block-quoted-without-attributes-no-spaces": {
			configRaw: configurationString{
				raw: `
provider"test"{}

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
		"provider-block-unquoted-without-attributes-no-trailing-space": {
			configRaw: configurationString{
				raw: `
provider test{}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configRaw.HasProviderBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(testCase.expected, got); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfiguration_HasTerraformBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configRaw configurationString
		expected  bool
	}{
		"no-config": {
			expected: false,
		},
		"terraform-meta-attribute": {
			configRaw: configurationString{
				raw: `
resource "test_test" "test" {
  terraform = test.test
}
`,
			},
			expected: false,
		},
		"terraform-object-attribute": {
			configRaw: configurationString{
				raw: `
resource "test_test" "test" {
  test = {
    terraform = {
      test = true
    }
  }
}
`,
			},
			expected: false,
		},
		"terraform-string-attribute": {
			configRaw: configurationString{
				raw: `
resource "test_test" "test" {
  test = {
    terraform = "test"
  }
}
`,
			},
			expected: false,
		},
		"terraform-block": {
			configRaw: configurationString{
				raw: `
terraform {
  test = true
}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"terraform-block-no-space": {
			configRaw: configurationString{
				raw: `
terraform{
  test = true
}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configRaw.HasTerraformBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(testCase.expected, got); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigurationString_Write(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configRaw configurationString
	}{
		"raw": {
			configRaw: configurationString{
				`
provider "test" {
 test = true
}

resource "test_test" "test" {}
`,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			err := testCase.configRaw.Write(context.Background(), tempDir)

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			expectedBytes := []byte(testCase.configRaw.raw)

			gotBytes, err := os.ReadFile(filepath.Join(tempDir, rawConfigFileName))

			if err != nil {
				t.Errorf("error reading file: %s", err)
			}

			if diff := cmp.Diff(gotBytes, expectedBytes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}

		})
	}
}
