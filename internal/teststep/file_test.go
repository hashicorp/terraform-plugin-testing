// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestConfigurationFile_HasProviderBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile    configurationFile
		expected      bool
		expectedError *regexp.Regexp
	}{
		"not-file": {
			configFile: configurationFile{
				file: "testdata/empty_file/not_a_real_file.tf",
			},
			expectedError: regexp.MustCompile(`.*no such file or directory`),
		},
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
		"provider-block-quoted-with-attributes-no-spaces": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_with_attributes_no_spaces/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_with_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes-no-trailing-space": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_with_attributes_no_trailing_space/main.tf",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_without_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes-no-spaces": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_without_attributes_no_spaces/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_without_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes-no-trailing-space": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_without_attributes_no_trailing_space/main.tf",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configFile.HasProviderBlock(context.Background())

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

func TestConfigurationFile_HasProviderBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile    configurationFile
		expected      bool
		expectedError *regexp.Regexp
	}{
		"not-file": {
			configFile: configurationFile{
				file: "testdata/empty_file/not_a_real_file.tf",
			},
			expectedError: regexp.MustCompile(`.*no such file or directory`),
		},
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
		"provider-block-quoted-with-attributes-no-spaces": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_with_attributes_no_spaces/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_with_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes-no-trailing-space": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_with_attributes_no_trailing_space/main.tf",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_without_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes-no-spaces": {
			configFile: configurationFile{
				file: "testdata/provider_block_quoted_without_attributes_no_spaces/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_without_attributes/main.tf",
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes-no-trailing-space": {
			configFile: configurationFile{
				file: "testdata/provider_block_unquoted_without_attributes_no_trailing_space/main.tf",
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

			testCase.configFile.file = filepath.Join(pwd, testCase.configFile.file)

			got, err := testCase.configFile.HasProviderBlock(context.Background())

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

func TestConfigurationFile_HasTerraformBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile    configurationFile
		expected      bool
		expectedError *regexp.Regexp
	}{
		"not-file": {
			configFile: configurationFile{
				file: "testdata/empty_file/not_a_real_file.tf",
			},
			expectedError: regexp.MustCompile(`.*no such file or directory`),
		},
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
		"terraform-block-no-space": {
			configFile: configurationFile{
				file: "testdata/terraform_block_no_space/main.tf",
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.configFile.HasTerraformBlock(context.Background())

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

func TestConfigurationFile_HasTerraformBlock_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile    configurationFile
		expected      bool
		expectedError *regexp.Regexp
	}{
		"not-file": {
			configFile: configurationFile{
				file: "testdata/empty_file/not_a_real_file.tf",
			},
			expectedError: regexp.MustCompile(`.*no such file or directory`),
		},
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
		"terraform-block-no-space": {
			configFile: configurationFile{
				file: "testdata/terraform_block_no_space/main.tf",
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

			testCase.configFile.file = filepath.Join(pwd, testCase.configFile.file)

			got, err := testCase.configFile.HasTerraformBlock(context.Background())

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

func TestConfigurationFile_Write(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile    configurationFile
		expectedError *regexp.Regexp
	}{
		"not-file": {
			configFile: configurationFile{
				file: "testdata/empty_file/not_a_real_file.tf",
			},
			expectedError: regexp.MustCompile(`.*no such file or directory`),
		},
		"file": {
			configFile: configurationFile{
				file: "testdata/random/random.tf",
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			err := testCase.configFile.Write(context.Background(), tempDir)

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
				fileInfo, err := os.Lstat(testCase.configFile.file)

				if err != nil {
					t.Errorf("error getting dir entry info: %s", err)
				}

				tempFileInfo, err := os.Lstat(filepath.Join(tempDir, filepath.Base(testCase.configFile.file)))

				if err != nil {
					t.Errorf("error getting temp dir entry info: %s", err)
				}

				if diff := cmp.Diff(tempFileInfo, fileInfo, fileInfoComparer); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			}
		})
	}
}

func TestConfigurationFile_Write_AbsolutePath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configFile    configurationFile
		expectedError *regexp.Regexp
	}{
		"not-file": {
			configFile: configurationFile{
				file: "testdata/empty_file/not_a_real_file.tf",
			},
			expectedError: regexp.MustCompile(`.*no such file or directory`),
		},
		"file": {
			configFile: configurationFile{
				file: "testdata/random/random.tf",
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

			testCase.configFile.file = filepath.Join(pwd, testCase.configFile.file)

			err = testCase.configFile.Write(context.Background(), tempDir)

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
				fileInfo, err := os.Lstat(testCase.configFile.file)

				if err != nil {
					t.Errorf("error getting dir entry info: %s", err)
				}

				tempFileInfo, err := os.Lstat(filepath.Join(tempDir, filepath.Base(testCase.configFile.file)))

				if err != nil {
					t.Errorf("error getting temp dir entry info: %s", err)
				}

				if diff := cmp.Diff(tempFileInfo, fileInfo, fileInfoComparer); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			}
		})
	}
}

func TestConfigFile_Append(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		filename        string
		appendContent   string
		expectedContent string
	}{
		"append content to a ConfigFile": {
			filename:        `testdata/main.tf`, // Contains `// Hello world`
			appendContent:   `terraform {}`,
			expectedContent: "# Copyright IBM Corp. 2014, 2026\n# SPDX-License-Identifier: MPL-2.0\n\n// Hello world" + "\n" + "terraform {}",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			prepareConfigRequest := PrepareConfigurationRequest{
				File: func(config.TestStepConfigRequest) string {
					return testCase.filename
				},
			}

			teststepConfig := Configuration(prepareConfigRequest.Exec())
			teststepConfig = teststepConfig.Append(testCase.appendContent)

			tempdir := t.TempDir()
			if err := teststepConfig.Write(context.Background(), tempdir); err != nil {
				t.Fatalf("failed to write file: %s", err)
			}

			got, err := os.ReadFile(filepath.Join(tempdir, filepath.Base(testCase.filename)))
			if err != nil {
				t.Fatalf("failed to read file: %s", err)
			}

			gotS := string(got[:])
			if diff := cmp.Diff(testCase.expectedContent, gotS); diff != "" {
				t.Errorf("expected %+v, got %+v", testCase.expectedContent, gotS)
			}
		})
	}
}
