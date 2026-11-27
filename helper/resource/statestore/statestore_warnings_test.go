// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestStateStore_warnings ensures that warning diagnostics returned by a state
// store implementation's GetStates method are not mistaken for workspace names
// by the `terraform workspace list` assertion, which would cause the state
// store test to fail. See issue #692.
func TestStateStore_warnings(t *testing.T) {
	// Setting this environment variable ensures TF core uses pluggable state storage during init.
	// This is only temporary while PSS is experimental.
	t.Setenv("TF_ENABLE_PLUGGABLE_STATE_STORAGE", "1")

	stateStore := exampleCloudValidStateStore()

	getStatesFunc := stateStore.GetStatesFunc
	stateStore.GetStatesFunc = func(ctx context.Context, req statestore.GetStatesRequest, resp *statestore.GetStatesResponse) {
		getStatesFunc(ctx, req, resp)

		// Report the "default" workspace as existing while no states have been
		// written yet, so the state store test assertions that expect the
		// "default" workspace after init continue to hold. See issue #691 for
		// the ongoing discussion about that expectation.
		if len(resp.StateIDs) == 0 {
			resp.StateIDs = []string{"default"}
		}

		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityWarning,
			Summary:  "test warning from the GetStates method",
			Detail:   "This warning is intentionally returned by the test state store.",
		})
	}

	r.UnitTest(t, r.TestCase{
		// State stores currently are only available in alpha releases or built from source
		// with experiments enabled: `go install -ldflags="-X main.experimentsAllowed=yes" .`
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_15_0),
			tfversion.SkipIfNotPrerelease(),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				StateStores: map[string]*testprovider.StateStore{
					"examplecloud_inmem": stateStore,
				},
			}),
		},
		Steps: []r.TestStep{
			{
				StateStore:           true,
				DefaultWorkspaceOnly: true,
				Config: `
					terraform {
					  required_providers {
					  	examplecloud = {
						  source = "registry.terraform.io/hashicorp/examplecloud"
						}
					  }
					  state_store "examplecloud_inmem" {
					  	provider "examplecloud" {}
					  }
					}
				`,
			},
		},
	})
}
