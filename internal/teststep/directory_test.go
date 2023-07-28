// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigurationDirectory_HasProviderBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
	}{
		"no-config": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_dir",
			},
			expected: false,
		},
		"provider-meta-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_meta_attribute",
			},
			expected: false,
		},
		"provider-object-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_object_attribute",
			},
			expected: false,
		},
		"provider-string-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_string_attribute",
			},
			expected: false,
		},
		"provider-block-quoted-with-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_quoted_with_attributes",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_with_attributes",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_quoted_without_attributes",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_without_attributes",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configDirectory.HasProviderBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}

func TestConfigurationDirectory_HasProviderBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
	}{
		"no-config": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_dir",
			},
			expected: false,
		},
		"provider-meta-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_meta_attribute",
			},
			expected: false,
		},
		"provider-object-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_object_attribute",
			},
			expected: false,
		},
		"provider-string-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_string_attribute",
			},
			expected: false,
		},
		"provider-block-quoted-with-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_quoted_with_attributes",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_with_attributes",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_quoted_without_attributes",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_without_attributes",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pwd, err := os.Getwd()

			if err != nil {
				t.Errorf("error getting wd: %s", err)
			}

			testCase.configDirectory.directory = filepath.Join(pwd, testCase.configDirectory.directory)

			got, err := testCase.configDirectory.HasProviderBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}

func TestConfigurationDirectory_HasTerraformBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
	}{
		"no-config": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_dir",
			},
			expected: false,
		},
		"terraform-meta-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_meta_attribute",
			},
			expected: false,
		},
		"terraform-object-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_object_attribute",
			},
			expected: false,
		},
		"terraform-string-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_string_attribute",
			},
			expected: false,
		},
		"terraform-block": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_block",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configDirectory.HasTerraformBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}

func TestConfigurationDirectory_HasTerraformBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
	}{
		"no-config": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_dir",
			},
			expected: false,
		},
		"terraform-meta-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_meta_attribute",
			},
			expected: false,
		},
		"terraform-object-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_object_attribute",
			},
			expected: false,
		},
		"terraform-string-attribute": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_string_attribute",
			},
			expected: false,
		},
		"terraform-block": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_block",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pwd, err := os.Getwd()

			if err != nil {
				t.Errorf("error getting wd: %s", err)
			}

			testCase.configDirectory.directory = filepath.Join(pwd, testCase.configDirectory.directory)

			got, err := testCase.configDirectory.HasTerraformBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}
