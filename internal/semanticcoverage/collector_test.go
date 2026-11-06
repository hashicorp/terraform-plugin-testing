// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestFindingIdentifiesTheStep(t *testing.T) {
	t.Setenv(EnvDirectory, t.TempDir())

	c, err := New("TestExample")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.BeginStep(2).Finish(); err != nil {
		t.Fatal(err)
	}
	want := Finding{
		Version: FindingVersion,
		// The full import path keeps same-named tests in packages with the
		// same directory name apart.
		Package: "github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage",
		Test:    "TestExample",
		Step:    2,
	}
	if diff := cmp.Diff([]Finding{want}, writtenFindings(t, c)); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

func TestFindingRecordsTheServedProviders(t *testing.T) {
	t.Parallel()

	c := testCollector(t)
	step := c.BeginStep(1)
	step.ObserveServedProviders([]string{"registry.terraform.io/stale/example"})
	// TF_ACC_PROVIDER_NAMESPACE can be mixed case, but Terraform reports
	// lower case source addresses.
	step.ObserveServedProviders([]string{"registry.terraform.io/hashicorp/echo", "registry.terraform.io/Example-Org/example"})
	if err := step.Finish(); err != nil {
		t.Fatal(err)
	}
	want := []string{"registry.terraform.io/example-org/example", "registry.terraform.io/hashicorp/echo"}
	if diff := cmp.Diff(want, writtenFindings(t, c)[0].ServedProviders); diff != "" {
		t.Errorf("served providers mismatch (-want +got):\n%s", diff)
	}
}

func TestOnlyFinishedStepsWriteFindings(t *testing.T) {
	t.Parallel()

	c := testCollector(t)
	skipped := c.BeginStep(1)
	skipped.ObserveState(state(resource("example_item", map[string]any{"id": "x"})))
	if err := c.BeginStep(2).Finish(); err != nil {
		t.Fatal(err)
	}
	findings := writtenFindings(t, c)
	if len(findings) != 1 || findings[0].Step != 2 || len(findings[0].Features) != 0 {
		t.Fatalf("findings = %+v, want only an empty step 2", findings)
	}
}

func TestFindingKeepsNoValues(t *testing.T) {
	t.Parallel()

	c := testCollector(t)
	s := c.BeginStep(1)
	s.ObserveConfigDir(configDir(t, `
resource "example_item" "secret_name" {
  token = "secret-configuration"
}
`))
	p := plan(itemChange(tfjson.ActionCreate))
	p.ResourceChanges[0].Address = `example_item.secret_name["secret-instance-key"]`
	p.ResourceChanges[0].Change.After = map[string]any{"token": "secret-planned-value"}
	s.ObservePlan(p)
	r := resource("example_item", map[string]any{"token": "secret-value", "map": map[string]any{"secret-map-key": "secret-value"}})
	r.Address = `example_item.secret_name["secret-instance-key"]`
	s.ObserveState(state(r))
	if err := s.Finish(); err != nil {
		t.Fatal(err)
	}
	data := writtenFiles(t, c)[0]
	for _, secret := range []string{"secret-configuration", "secret_name", "secret-instance-key", "secret-planned-value", "secret-map-key", "secret-value"} {
		if strings.Contains(string(data), secret) {
			t.Errorf("finding persisted %q: %s", secret, data)
		}
	}
}

func TestDisabledCollectionDoesNothing(t *testing.T) {
	t.Parallel()

	var c *Collector
	s := c.BeginStep(1)
	if s.Enabled() {
		t.Fatal("step of a nil collector is enabled")
	}
	s.ObserveConfigDir(t.TempDir())
	s.ObserveQueryConfig("missing.tfquery.hcl")
	s.ObservePlan(nil)
	s.ApplyCompleted(nil)
	s.ObserveState(nil)
	s.ImportCompleted(nil, "")
	s.RefreshCompleted(nil)
	s.QueryCompleted()
	s.StateStoreCompleted(true)
	s.StateStoreLockVerified()
	if err := s.Finish(); err != nil {
		t.Fatal(err)
	}
	d := c.BeginPostTestDestroy(2)
	d.DestroyCompleted(nil)
	if err := d.Finish(); err != nil {
		t.Fatal(err)
	}
}

func TestPostTestDestroyFinding(t *testing.T) {
	t.Parallel()

	c := testCollector(t)
	s := c.BeginPostTestDestroy(3)
	data := resource("example_lookup", nil)
	data.Mode = tfjson.DataResourceMode
	s.DestroyCompleted(state(resource("example_item", map[string]any{"id": "x"}), data))
	if err := s.Finish(); err != nil {
		t.Fatal(err)
	}
	want := []Finding{{
		Version:         FindingVersion,
		Test:            "TestExample",
		Step:            3,
		PostTestDestroy: true,
		// Only managed resources are deleted, and the destroy records no
		// attributes.
		Features: []Feature{{Kind: KindResource, Name: "example_item", Provider: exampleProvider, Completed: []Operation{OperationDelete}}},
	}}
	if diff := cmp.Diff(want, writtenFindings(t, c)); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

func TestNewFailsEarlyForUnusableDirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvDirectory, filepath.Join(file, "findings"))
	if c, err := New("TestExample"); err == nil || c != nil {
		t.Fatalf("collector=%v err=%v", c, err)
	}
	t.Setenv(EnvDirectory, "")
	if c, err := New("TestExample"); err != nil || c != nil {
		t.Fatalf("disabled: collector=%v err=%v", c, err)
	}
}

// testCollector returns a collector that writes to a new directory.
func testCollector(t *testing.T) *Collector {
	t.Helper()

	return &Collector{directory: t.TempDir(), testName: "TestExample"}
}

// stepFinding records evidence in step 1 with record and returns the finding
// that the step writes.
func stepFinding(t *testing.T, record func(s *Step)) Finding {
	t.Helper()

	c := testCollector(t)
	s := c.BeginStep(1)
	record(s)
	if err := s.Finish(); err != nil {
		t.Fatal(err)
	}
	findings := writtenFindings(t, c)
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1", len(findings))
	}
	return findings[0]
}

func writtenFindings(t *testing.T, c *Collector) []Finding {
	t.Helper()

	var findings []Finding
	for _, data := range writtenFiles(t, c) {
		var f Finding
		if err := json.Unmarshal(data, &f); err != nil {
			t.Fatal(err)
		}
		findings = append(findings, f)
	}
	return findings
}

func writtenFiles(t *testing.T, c *Collector) [][]byte {
	t.Helper()

	entries, err := os.ReadDir(c.directory)
	if err != nil {
		t.Fatal(err)
	}
	var files [][]byte
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(c.directory, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, data)
	}
	return files
}

// configDir returns a directory with main.tf containing config.
func configDir(t *testing.T, config string) string {
	t.Helper()

	return configFiles(t, map[string]string{"main.tf": config})
}

// configFiles returns a directory with the given files.
func configFiles(t *testing.T, files map[string]string) string {
	t.Helper()

	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

const exampleProvider = "registry.terraform.io/hashicorp/example"

func plan(changes ...*tfjson.ResourceChange) *tfjson.Plan {
	return &tfjson.Plan{ResourceChanges: changes}
}

// itemChange returns a change of the managed resource example_item.test.
func itemChange(actions ...tfjson.Action) *tfjson.ResourceChange {
	return &tfjson.ResourceChange{
		Address:      "example_item.test",
		Type:         "example_item",
		Mode:         tfjson.ManagedResourceMode,
		ProviderName: exampleProvider,
		Change:       &tfjson.Change{Actions: actions},
	}
}

func state(resources ...*tfjson.StateResource) *tfjson.State {
	return &tfjson.State{Values: &tfjson.StateValues{RootModule: &tfjson.StateModule{Resources: resources}}}
}

// resource returns the managed resource typ.test with the given values.
func resource(typ string, values map[string]any) *tfjson.StateResource {
	return &tfjson.StateResource{
		Address:         typ + ".test",
		Type:            typ,
		Mode:            tfjson.ManagedResourceMode,
		ProviderName:    exampleProvider,
		AttributeValues: values,
	}
}
