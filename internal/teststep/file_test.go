// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigurationFile_HasProviderBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile configurationFile
		expected   bool
	}{
		"no-config": {
			configFile: configurationFile{
				file: "testdata/empty_file/main.tf",
			},
			expected: false,
		},
		"provider-meta-attribute": {
			configFile: configurationFile{
				file: "testdata/provider_meta_attribute/main.tf",
			},
			expected: false,
		},
		"provider-object-attribute": {
			configFile: configurationFile{
				file: "testdata/provider_object_attribute/main.tf",
			},
			expected: false,
		},
		"provider-string-attribute": {
			configFile: configurationFile{
				file: "testdata/provider_string_attribute/main.tf",
			},
			expected: false,
		},
		"provider-block-quoted-with-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_with_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_with_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_without_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_without_attributes/main.tf",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configFile.HasProviderBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}

func TestConfigurationFile_HasProviderBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile configurationFile
		expected   bool
	}{
		"no-config": {
			configFile: configurationFile{
				file: "testdata/empty_file/main.tf",
			},
			expected: false,
		},
		"provider-meta-attribute": {
			configFile: configurationFile{
				file: "testdata/provider_meta_attribute/main.tf",
			},
			expected: false,
		},
		"provider-object-attribute": {
			configFile: configurationFile{
				file: "testdata/provider_object_attribute/main.tf",
			},
			expected: false,
		},
		"provider-string-attribute": {
			configFile: configurationFile{
				file: "testdata/provider_string_attribute/main.tf",
			},
			expected: false,
		},
		"provider-block-quoted-with-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_with_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_with_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_without_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_without_attributes/main.tf",
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

			testCase.configFile.file = filepath.Join(pwd, testCase.configFile.file)

			got, err := testCase.configFile.HasProviderBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}

func TestConfigurationFile_HasTerraformBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile configurationFile
		expected   bool
	}{
		"no-config": {
			configFile: configurationFile{
				file: "testdata/empty_file/main.tf",
			},
			expected: false,
		},
		"terraform-meta-attribute": {
			configFile: configurationFile{
				file: "testdata/terraform_meta_attribute/main.tf",
			},
			expected: false,
		},
		"terraform-object-attribute": {
			configFile: configurationFile{
				file: "testdata/terraform_object_attribute/main.tf",
			},
			expected: false,
		},
		"terraform-string-attribute": {
			configFile: configurationFile{
				file: "testdata/terraform_string_attribute/main.tf",
			},
			expected: false,
		},
		"terraform-block": {
			configFile: configurationFile{
				file: "testdata/terraform_block/main.tf",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configFile.HasTerraformBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}

func TestConfigurationFile_HasTerraformBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile configurationFile
		expected   bool
	}{
		"no-config": {
			configFile: configurationFile{
				file: "testdata/empty_file/main.tf",
			},
			expected: false,
		},
		"terraform-meta-attribute": {
			configFile: configurationFile{
				file: "testdata/terraform_meta_attribute/main.tf",
			},
			expected: false,
		},
		"terraform-object-attribute": {
			configFile: configurationFile{
				file: "testdata/terraform_object_attribute/main.tf",
			},
			expected: false,
		},
		"terraform-string-attribute": {
			configFile: configurationFile{
				file: "testdata/terraform_string_attribute/main.tf",
			},
			expected: false,
		},
		"terraform-block": {
			configFile: configurationFile{
				file: "testdata/terraform_block/main.tf",
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

			testCase.configFile.file = filepath.Join(pwd, testCase.configFile.file)

			got, err := testCase.configFile.HasTerraformBlock(context.Background())

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}
