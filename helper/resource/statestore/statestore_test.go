package statestore_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestStateStore_validation_error(t *testing.T) {
	// Setting this environment variable ensures TF core uses pluggable state storage during init.
	// This is only temporary while PSS is experimental.
	t.Setenv("TF_ENABLE_PLUGGABLE_STATE_STORAGE", "1")

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
					"examplecloud_inmem": {
						SchemaResponse: &statestore.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "path",
											Type:     tftypes.String,
											Required: true,
										},
									},
								},
							},
						},
						ValidateConfigResponse: &statestore.ValidateConfigResponse{
							Diagnostics: []*tfprotov6.Diagnostic{
								{
									Severity:  tfprotov6.DiagnosticSeverityError,
									Summary:   "WHOOPS",
									Detail:    "Something isn't right about that request :D, error it is!",
									Attribute: tftypes.NewAttributePath().WithAttributeName("path"),
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				StateStore: true,
				Config: `
					terraform {
					  required_providers {
					  	examplecloud = {
						  source = "registry.terraform.io/hashicorp/examplecloud"
						}
					  }
					  state_store "examplecloud_inmem" {
					  	provider "examplecloud" {}

					    path = "test_state_file.tfstate"
					  }
					}
				`,
				ExpectError: regexp.MustCompile(`Something isn't right about that request :D, error it is!`),
			},
		},
	})
}

func TestStateStore_configure_error(t *testing.T) {
	// Setting this environment variable ensures TF core uses pluggable state storage during init.
	// This is only temporary while PSS is experimental.
	t.Setenv("TF_ENABLE_PLUGGABLE_STATE_STORAGE", "1")

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
					"examplecloud_inmem": {
						SchemaResponse: &statestore.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{Attributes: []*tfprotov6.SchemaAttribute{}},
							},
						},
						ConfigureResponse: &statestore.ConfigureResponse{
							Diagnostics: []*tfprotov6.Diagnostic{
								{
									Severity: tfprotov6.DiagnosticSeverityError,
									Summary:  "WHOOPS",
									Detail:   "The configure has failed us! :P",
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				StateStore: true,
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
				ExpectError: regexp.MustCompile(`The configure has failed us! :P`),
			},
		},
	})
}

func TestStateStore_inmem(t *testing.T) {
	// Setting this environment variable ensures TF core uses pluggable state storage during init.
	// This is only temporary while PSS is experimental.
	t.Setenv("TF_ENABLE_PLUGGABLE_STATE_STORAGE", "1")

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
					"examplecloud_inmem": exampleCloudStateStoreFS(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				StateStore: true,
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

// What we can test so far w/ smoke test:
// - State store validation
// - State store configuration
// - State store can write state
// - State store can read state
// - State store handles multiple workspaces

// What is kind of awkward to test w/ smoke test
// - State size / chunking. You can do this by just changing the config you pass in w/ StateStore (we always inject a single resource to verify state apply will work)
// - Multiple applies to the same state (currently it just does one apply to the bar workspace)

// What we can't test w/ smoke test
// - State store locks
// - State store unlocks
// - State store locks with a ton of clients (i.e. load test)
