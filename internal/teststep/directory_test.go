// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfigurationDirectory_HasProviderBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
		expectedError   *regexp.Regexp
	}{
		"not-directory": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_file/main.tf",
			},
			expectedError: regexp.MustCompile(`.*not a directory`),
		},
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
		"provider-block-quoted-with-attributes-no-spaces": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_quoted_with_attributes_no_spaces",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_with_attributes",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes_no-trailing-space": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_with_attributes_no_trailing_space",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_quoted_without_attributes",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes-no-spaces": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_quoted_without_attributes_no_spaces",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_without_attributes",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes-no-trailing-space": {
			configDirectory: configurationDirectory{
				directory: "testdata/provider_block_unquoted_without_attributes_no_trailing_space",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configDirectory.HasProviderBlock(context.Background())

			if testCase.expectedError == nil && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != nil && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != nil && err != nil {
				if !testCase.expectedError.MatchString(err.Error()) {
					t.Errorf("expected error %s, got error %s", testCase.expectedError.String(), err)
				}
			}

			if diff := cmp.Diff(testCase.expected, got); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigurationDirectory_HasProviderBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
		expectedError   *regexp.Regexp
	}{
		"not-directory": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_file/main.tf",
			},
			expectedError: regexp.MustCompile(`.*not a directory`),
		},
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

			if testCase.expectedError == nil && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != nil && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != nil && err != nil {
				if !testCase.expectedError.MatchString(err.Error()) {
					t.Errorf("expected error %s, got error %s", testCase.expectedError.String(), err)
				}
			}

			if diff := cmp.Diff(testCase.expected, got); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigurationDirectory_HasTerraformBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
		expectedError   *regexp.Regexp
	}{
		"not-directory": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_file/main.tf",
			},
			expectedError: regexp.MustCompile(`.*not a directory`),
		},
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

			if testCase.expectedError == nil && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != nil && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != nil && err != nil {
				if !testCase.expectedError.MatchString(err.Error()) {
					t.Errorf("expected error %s, got error %s", testCase.expectedError.String(), err)
				}
			}

			if diff := cmp.Diff(testCase.expected, got); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigurationDirectory_HasTerraformBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expected        bool
		expectedError   *regexp.Regexp
	}{
		"not-directory": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_file/main.tf",
			},
			expectedError: regexp.MustCompile(`.*not a directory`),
		},
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

			if testCase.expectedError == nil && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != nil && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != nil && err != nil {
				if !testCase.expectedError.MatchString(err.Error()) {
					t.Errorf("expected error %s, got error %s", testCase.expectedError.String(), err)
				}
			}

			if diff := cmp.Diff(testCase.expected, got); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigurationDirectory_Write(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expectedError   *regexp.Regexp
	}{
		"not-directory": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_file/main.tf",
			},
			expectedError: regexp.MustCompile(`.*not a directory`),
		},
		"no-config": {
			configDirectory: configurationDirectory{
				"testdata/empty_dir",
			},
		},
		"dir-single-file": {
			configDirectory: configurationDirectory{
				"testdata/random",
			},
		},
		"dir-multiple-files": {
			configDirectory: configurationDirectory{
				"testdata/random_multiple_files",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			err := testCase.configDirectory.Write(context.Background(), tempDir)

			if testCase.expectedError == nil && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != nil && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != nil && err != nil {
				if !testCase.expectedError.MatchString(err.Error()) {
					t.Errorf("expected error %s, got error %s", testCase.expectedError.String(), err)
				}
			}

			if err == nil {
				dirEntries, err := os.ReadDir(testCase.configDirectory.directory)

				if err != nil {
					t.Errorf("error reading directory: %s", err)
				}

				tempDirEntries, err := os.ReadDir(tempDir)

				if err != nil {
					t.Errorf("error reading temp directory: %s", err)
				}

				if len(dirEntries) != len(tempDirEntries) {
					t.Errorf("expected %d dir entries, got %d dir entries", dirEntries, tempDirEntries)
				}

				for k, v := range dirEntries {
					dirEntryInfo, err := v.Info()

					if err != nil {
						t.Errorf("error getting dir entry info: %s", err)
					}

					tempDirEntryInfo, err := tempDirEntries[k].Info()

					if err != nil {
						t.Errorf("error getting temp dir entry info: %s", err)
					}

					if diff := cmp.Diff(tempDirEntryInfo, dirEntryInfo, fileInfoComparer); diff != "" {
						t.Errorf("unexpected difference: %s", diff)
					}
				}
			}
		})
	}
}

func TestConfigurationDirectory_Write_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expectedError   *regexp.Regexp
	}{
		"not-directory": {
			configDirectory: configurationDirectory{
				directory: "testdata/empty_file/main.tf",
			},
			expectedError: regexp.MustCompile(`.*not a directory`),
		},
		"no-config": {
			configDirectory: configurationDirectory{
				"testdata/empty_dir",
			},
		},
		"dir-single-file": {
			configDirectory: configurationDirectory{
				"testdata/random",
			},
		},
		"dir-multiple-files": {
			configDirectory: configurationDirectory{
				"testdata/random_multiple_files",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			pwd, err := os.Getwd()

			if err != nil {
				t.Errorf("error getting wd: %s", err)
			}

			testCase.configDirectory.directory = filepath.Join(pwd, testCase.configDirectory.directory)

			err = testCase.configDirectory.Write(context.Background(), tempDir)

			if testCase.expectedError == nil && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != nil && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != nil && err != nil {
				if !testCase.expectedError.MatchString(err.Error()) {
					t.Errorf("expected error %s, got error %s", testCase.expectedError.String(), err)
				}
			}

			if err == nil {
				dirEntries, err := os.ReadDir(testCase.configDirectory.directory)

				if err != nil {
					t.Errorf("error reading directory: %s", err)
				}

				tempDirEntries, err := os.ReadDir(tempDir)

				if err != nil {
					t.Errorf("error reading temp directory: %s", err)
				}

				if len(dirEntries) != len(tempDirEntries) {
					t.Errorf("expected %d dir entries, got %d dir entries", dirEntries, tempDirEntries)
				}

				for k, v := range dirEntries {
					dirEntryInfo, err := v.Info()

					if err != nil {
						t.Errorf("error getting dir entry info: %s", err)
					}

					tempDirEntryInfo, err := tempDirEntries[k].Info()

					if err != nil {
						t.Errorf("error getting temp dir entry info: %s", err)
					}

					if diff := cmp.Diff(tempDirEntryInfo, dirEntryInfo, fileInfoComparer); diff != "" {
						t.Errorf("unexpected difference: %s", diff)
					}
				}
			}
		})
	}
}

var fileInfoComparer = cmp.Comparer(func(x, y os.FileInfo) bool {
	if x.Name() != y.Name() {
		return false
	}

	if x.Mode() != y.Mode() {
		return false
	}

	if x.Size() != y.Size() {
		return false
	}

	return true
})
