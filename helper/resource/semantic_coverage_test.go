// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const coverageTestConfig = `resource "test_resource" "secret_name" {
  name = "secret-value"
}
`

func TestSemanticCoverageWritesAFindingPerStep(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(semanticcoverage.EnvDirectory, dir)

	UnitTest(t, coverageTestCase(
		TestStep{Config: coverageTestConfig},
		// A step that expects an error writes no finding.
		TestStep{
			Config:      coverageTestConfig + `resource "test_invalid" "b" {}`,
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		TestStep{
			Config:   coverageTestConfig,
			SkipFunc: func() (bool, error) { return true, nil },
		},
		// Restores a valid configuration for the post-test destroy.
		TestStep{Config: coverageTestConfig},
		// Refresh uses the configuration in the working directory, so the
		// configured name counts as configured, not as computed.
		TestStep{RefreshState: true},
	))

	const provider = "registry.terraform.io/hashicorp/test"
	item := semanticcoverage.Feature{
		Kind:       semanticcoverage.KindResource,
		Name:       "test_resource",
		Provider:   provider,
		Configured: []string{"name"},
		Populated:  []string{"name"},
		// Every passing config step ends with a plan that refreshes.
		Completed: []semanticcoverage.Operation{semanticcoverage.OperationRead},
	}
	created := item
	created.Planned = []semanticcoverage.Operation{semanticcoverage.OperationCreate}
	created.Completed = []semanticcoverage.Operation{semanticcoverage.OperationCreate, semanticcoverage.OperationRead}
	configured := semanticcoverage.Feature{
		Kind:      semanticcoverage.KindProvider,
		Name:      provider,
		Completed: []semanticcoverage.Operation{semanticcoverage.OperationConfigure},
	}
	served := []string{provider}
	want := []semanticcoverage.Finding{
		{
			Step:            1,
			ServedProviders: served,
			Features:        []semanticcoverage.Feature{configured, created},
		},
		{
			// A no-op apply has no operation.
			Step:            4,
			ServedProviders: served,
			Features:        []semanticcoverage.Feature{configured, item},
		},
		{
			Step:            5,
			ServedProviders: served,
			Features:        []semanticcoverage.Feature{configured, item},
		},
		{
			// The destroy after the last step deletes what is in state.
			Step:            6,
			PostTestDestroy: true,
			ServedProviders: served,
			Features: []semanticcoverage.Feature{{
				Kind:      semanticcoverage.KindResource,
				Name:      "test_resource",
				Provider:  provider,
				Completed: []semanticcoverage.Operation{semanticcoverage.OperationDelete},
			}},
		},
	}
	got := writtenCoverage(t, dir)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

func TestSemanticCoverageOfStepsThatStopTheTest(t *testing.T) {
	tests := map[string]struct {
		step         TestStep
		checkDestroy TestCheckFunc
		wantFindings int
	}{
		"a failed step writes no finding": {
			step: TestStep{
				Config: coverageTestConfig,
				Check:  func(*terraform.State) error { return errors.New("check failed") },
			},
			wantFindings: 0,
		},
		// The failed CheckDestroy also keeps the post-test destroy from
		// writing a finding.
		"a failing cleanup doesn't drop the last step": {
			step:         TestStep{Config: coverageTestConfig},
			checkDestroy: func(*terraform.State) error { return errors.New("destroy check failed") },
			wantFindings: 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv(semanticcoverage.EnvDirectory, dir)

			c := coverageTestCase(test.step)
			c.CheckDestroy = test.checkDestroy
			func() {
				defer func() {
					if r := recover(); r != errStopTest {
						panic(r)
					}
				}()
				UnitTest(&stoppingT{}, c)
			}()

			got := writtenCoverage(t, dir)
			if len(got) != test.wantFindings {
				t.Fatalf("findings = %+v, want %d", got, test.wantFindings)
			}
			// Cleanup destroys the resource, but isn't part of the step.
			for _, finding := range got {
				for _, f := range finding.Features {
					if slices.Contains(f.Planned, semanticcoverage.OperationDelete) {
						t.Errorf("cleanup recorded as evidence: %+v", f)
					}
				}
			}
		})
	}
}

func TestServedAddresses(t *testing.T) {
	factories := &providerFactories{
		legacy:  sdkProviderFactories{"terraform-provider-legacy": nil},
		protov5: protov5ProviderFactories{"five": nil},
		protov6: protov6ProviderFactories{"six": nil},
	}
	tests := map[string]struct {
		namespace, host string
		want            []string
	}{
		// The namespaceless address only exists for Terraform 0.12.
		"default": {want: []string{
			"registry.terraform.io/hashicorp/five",
			"registry.terraform.io/hashicorp/legacy",
			"registry.terraform.io/hashicorp/six",
		}},
		"overridden": {namespace: "example-org", host: "example.com", want: []string{
			"example.com/example-org/five",
			"example.com/example-org/legacy",
			"example.com/example-org/six",
		}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv(EnvTfAccProviderNamespace, test.namespace)
			t.Setenv(EnvTfAccProviderHost, test.host)

			got := factories.servedAddresses()
			slices.Sort(got)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("served addresses mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

var errStopTest = errors.New("test stopped")

// stoppingT stops the test like testing.T does: Fatal marks it failed and
// stops it.
type stoppingT struct {
	mockT
	failed bool
}

func (t *stoppingT) Fatal(...any) { t.FailNow() }

func (t *stoppingT) Fatalf(string, ...any) { t.FailNow() }

func (t *stoppingT) FailNow() {
	t.failed = true
	panic(errStopTest)
}

func (t *stoppingT) Failed() bool { return t.failed || t.mockT.Failed() }

// coverageTestCase returns a test case of the given steps with a provider
// "test". Its resource test_resource has a name, and its resource
// test_invalid fails validation.
func coverageTestCase(steps ...TestStep) TestCase {
	return TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{Name: "name", Type: tftypes.String, Required: true},
									},
								},
							},
						},
					},
					"test_invalid": {
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{Block: &tfprotov6.SchemaBlock{}},
						},
						ValidateConfigResponse: &resource.ValidateConfigResponse{
							Diagnostics: []*tfprotov6.Diagnostic{
								{Severity: tfprotov6.DiagnosticSeverityError, Summary: "Invalid Attribute Value"},
							},
						},
					},
				},
			}),
		},
		Steps: steps,
	}
}

// writtenCoverage returns the findings in dir, sorted by step, without the
// fields that identify the test. It fails the test if a finding contains a
// value from coverageTestConfig.
func writtenCoverage(t *testing.T, dir string) []semanticcoverage.Finding {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var findings []semanticcoverage.Finding
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		for _, secret := range []string{"secret_name", "secret-value"} {
			if strings.Contains(string(data), secret) {
				t.Errorf("finding persisted %q: %s", secret, data)
			}
		}
		var f semanticcoverage.Finding
		if err := json.Unmarshal(data, &f); err != nil {
			t.Fatal(err)
		}
		if err := f.Validate(); err != nil {
			t.Fatal(err)
		}
		f.Version, f.Package, f.Test = 0, "", ""
		findings = append(findings, f)
	}
	slices.SortFunc(findings, func(a, b semanticcoverage.Finding) int { return a.Step - b.Step })
	return findings
}
