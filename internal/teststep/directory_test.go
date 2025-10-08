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
		"terraform-block-no-space": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_block_no_space",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
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
		"terraform-block-no-space": {
			configDirectory: configurationDirectory{
				directory: "testdata/terraform_block_no_space",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
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
				directory: "testdata/empty_dir",
			},
		},
		"dir-single-file": {
			configDirectory: configurationDirectory{
				directory: "testdata/random",
			},
		},
		"dir-multiple-files": {
			configDirectory: configurationDirectory{
				directory: "testdata/random_multiple_files",
			},
		},
	}

	for name, testCase := range testCases {
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
					t.Fatalf("error reading directory: %s", err)
				}

				tempDirEntries, err := os.ReadDir(tempDir)

				if err != nil {
					t.Fatalf("error reading temp directory: %s", err)
				}

				files := filesOnly(dirEntries)
				tempDirFiles := filesOnly(tempDirEntries)

				if len(tempDirFiles)-len(files) != 0 {
					t.Errorf("expected %d files, got %d files", len(files), tempDirFiles)
				}

				for i, file := range files {
					dirEntryInfo, err := file.Info()

					if err != nil {
						t.Errorf("error getting dir entry info: %s", err)
					}

					tempDirEntry := tempDirFiles[i]
					tempDirEntryInfo, err := tempDirEntry.Info()

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
				directory: "testdata/empty_dir",
			},
		},
		"dir-single-file": {
			configDirectory: configurationDirectory{
				directory: "testdata/random",
			},
		},
		"dir-multiple-files": {
			configDirectory: configurationDirectory{
				directory: "testdata/random_multiple_files",
			},
		},
	}

	for name, testCase := range testCases {
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
					t.Fatalf("error reading directory: %s", err)
				}

				tempDirEntries, err := os.ReadDir(tempDir)

				if err != nil {
					t.Fatalf("error reading temp directory: %s", err)
				}

				files := filesOnly(dirEntries)
				tempDirFiles := filesOnly(tempDirEntries)

				if len(tempDirFiles)-len(files) != 0 {
					t.Errorf("expected %d files, got %d files", len(files), tempDirFiles)
				}

				for i, file := range files {
					dirEntryInfo, err := file.Info()

					if err != nil {
						t.Errorf("error getting dir entry info: %s", err)
					}

					tempDirEntryInfo, err := tempDirFiles[i].Info()

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

func TestConfigurationDirectory_Write_WithAppendedConfig(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configDirectory configurationDirectory
		expectedError   *regexp.Regexp
	}{
		"dir-single-file": {
			configDirectory: configurationDirectory{
				directory: "testdata/random",
				appendedConfig: map[string]string{
					"import.tf": `terraform {\nimport\n{\nto = satellite.the_moon\nid = "moon"\n}\n}\n`,
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			err := testCase.configDirectory.Write(context.Background(), tempDir)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}

			dirEntries, err := os.ReadDir(testCase.configDirectory.directory)
			if err != nil {
				t.Errorf("error reading directory: %s", err)
			}

			tempDirEntries, err := os.ReadDir(tempDir)
			if err != nil {
				t.Errorf("error reading temp directory: %s", err)
			}

			files := filesOnly(dirEntries)
			tempDirFiles := filesOnly(tempDirEntries)

			if len(tempDirFiles)-len(files) != 1 {
				t.Errorf("expected %d files, got %d files", len(files)+1, tempDirFiles)
			}

			for _, file := range files {
				filename := file.Name()
				expectedContent, err := os.ReadFile(filepath.Join(testCase.configDirectory.directory, filename))
				if err != nil {
					t.Errorf("error reading file from config directory %s: %s", filename, err)
				}

				content, err := os.ReadFile(filepath.Join(tempDir, filename))
				if err != nil {
					t.Errorf("error reading generated file %s: %s", filename, err)
				}

				if diff := cmp.Diff(expectedContent, content); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			}

			appendedConfigFiles := testCase.configDirectory.appendedConfig
			for filename, expectedContent := range appendedConfigFiles {
				content, err := os.ReadFile(filepath.Join(tempDir, filename))
				if err != nil {
					t.Errorf("error reading appended config file %s: %s", filename, err)
				}

				if diff := cmp.Diff([]byte(expectedContent), content); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
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

func filesOnly(entries []os.DirEntry) []os.DirEntry {
	files := []os.DirEntry{}
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e)
		}
	}
	return files
}
