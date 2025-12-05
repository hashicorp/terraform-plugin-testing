// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	testinginterface "github.com/mitchellh/go-testing-interface"
)

func TestRunTFVersionChecks(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		versionChecks []tfversion.TerraformVersionCheck
		tfVersion     *version.Version
		expectError   bool
	}{
		"run-test": {
			versionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipIf(version.Must(version.NewVersion("1.1.0"))),
				tfversion.RequireAbove(version.Must(version.NewVersion("1.2.0"))),
			},
			tfVersion:   version.Must(version.NewVersion("1.3.0")),
			expectError: false,
		},
		"skip-test": {
			versionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipIf(version.Must(version.NewVersion("1.1.0"))),
			},
			tfVersion:   version.Must(version.NewVersion("1.1.0")),
			expectError: false,
		},
		"fail-test": {
			versionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireNot(version.Must(version.NewVersion("1.1.0"))),
			},
			tfVersion:   version.Must(version.NewVersion("1.1.0")),
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if test.expectError {
				plugintest.TestExpectTFatal(t, func() {
					runTFVersionChecks(context.Background(), &testinginterface.RuntimeT{}, test.tfVersion, test.versionChecks)
				})
			} else {
				runTFVersionChecks(context.Background(), t, test.tfVersion, test.versionChecks)
			}
		})
	}
}
