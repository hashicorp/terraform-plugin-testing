// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
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
