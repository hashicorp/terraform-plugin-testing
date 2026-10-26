// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Command semantic-coverage-report compares semantic coverage findings written
// by acceptance tests (TF_ACC_SEMANTIC_COVERAGE_DIR) with a
// `terraform providers schema -json` document.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
)

type evidence struct {
	Package string `json:"package,omitempty"`
	Test    string `json:"test"`
	Step    int    `json:"step"`
	Address string `json:"address,omitempty"`
	Kind    string `json:"kind"`
}
type attributeResult struct {
	Path     string     `json:"path"`
	Mode     string     `json:"mode"`
	Covered  bool       `json:"covered"`
	Evidence []evidence `json:"evidence,omitempty"`
}
type target struct {
	Provider   string                `json:"provider"`
	Kind       semanticcoverage.Kind `json:"kind"`
	Name       string                `json:"name"`
	Exercised  bool                  `json:"exercised"`
	Evidence   []evidence            `json:"evidence,omitempty"`
	Attributes []attributeResult     `json:"attributes,omitempty"`
	Planned    map[string][]evidence `json:"planned,omitempty"`
	Completed  map[string][]evidence `json:"completed,omitempty"`
}
type metric struct {
	Covered int `json:"covered"`
	Total   int `json:"total"`
}
type output struct {
	Version        int       `json:"version"`
	AttributeScope string    `json:"attribute_scope"`
	Findings       int       `json:"findings"`
	Entities       metric    `json:"entities"`
	Inputs         metric    `json:"inputs"`
	Computed       metric    `json:"computed"`
	Targets        []*target `json:"targets"`
	Warnings       []string  `json:"warnings,omitempty"`
}

func main() {
	findingsDir := flag.String("findings-dir", "findings", "directory of test-derived JSON findings")
	schemaFile := flag.String("schema-file", "", "terraform providers schema -json file")
	jsonOut := flag.String("json-out", "coverage.json", "JSON report path")
	textOut := flag.String("text-out", "coverage.txt", "text report path")
	flag.Parse()
	if *schemaFile == "" {
		die(fmt.Errorf("-schema-file is required"))
	}
	docBytes, err := os.ReadFile(*schemaFile)
	if err != nil {
		die(err)
	}
	// Unmarshalling also validates the schema format version.
	var doc tfjson.ProviderSchemas
	if err = json.Unmarshal(docBytes, &doc); err != nil {
		die(fmt.Errorf("%s: %w", *schemaFile, err))
	}
	if len(doc.Schemas) == 0 {
		die(fmt.Errorf("schema has no providers"))
	}
	items, err := readFindings(*findingsDir)
	if err != nil {
		die(err)
	}
	if len(items) == 0 {
		die(fmt.Errorf("no step findings in %s", *findingsDir))
	}
	result := aggregate(&doc, items)
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		die(err)
	}
	if err = os.WriteFile(*jsonOut, append(data, '\n'), 0600); err != nil {
		die(err)
	}
	var text strings.Builder
	fmt.Fprintf(&text, "Test-derived semantic coverage: %d steps\nAttribute scope: top-level schema attributes and blocks only; nested descendants are excluded.\nEntities %d/%d; inputs %d/%d; computed %d/%d\n", result.Findings, result.Entities.Covered, result.Entities.Total, result.Inputs.Covered, result.Inputs.Total, result.Computed.Covered, result.Computed.Total)
	for _, t := range result.Targets {
		covered := 0
		for _, a := range t.Attributes {
			if a.Covered {
				covered++
			}
		}
		fmt.Fprintf(&text, "%s %s %s: exercised=%t attributes=%d/%d planned=%v completed=%v\n", t.Provider, t.Kind, t.Name, t.Exercised, covered, len(t.Attributes), sortedKeys(t.Planned), sortedKeys(t.Completed))
		for _, a := range t.Attributes {
			if !a.Covered {
				fmt.Fprintf(&text, "  uncovered %s (%s)\n", a.Path, a.Mode)
			}
		}
	}
	for _, w := range result.Warnings {
		fmt.Fprintf(&text, "warning: %s\n", w)
	}
	if err = os.WriteFile(*textOut, []byte(text.String()), 0600); err != nil {
		die(err)
	}
	fmt.Printf("read %d step findings; wrote %s and %s\n", len(items), *jsonOut, *textOut)
}
func die(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
func readFindings(dir string) ([]semanticcoverage.Finding, error) {
	var items []semanticcoverage.Finding
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var item semanticcoverage.Finding
		if err = json.Unmarshal(data, &item); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if item.Version != semanticcoverage.FindingVersion {
			return fmt.Errorf("%s: finding version %d, this report reads version %d; re-run the tests with the same terraform-plugin-testing version", path, item.Version, semanticcoverage.FindingVersion)
		}
		if item.Test == "" || item.Step < 1 {
			return fmt.Errorf("%s: missing test name or step", path)
		}
		if !item.Outcome.Valid() {
			return fmt.Errorf("%s: invalid outcome %q", path, item.Outcome)
		}
		for _, warning := range item.Warnings {
			if !warning.Valid() {
				return fmt.Errorf("%s: invalid warning %q", path, warning)
			}
		}
		for _, feature := range item.Features {
			if !feature.Kind.Valid() {
				return fmt.Errorf("%s: invalid feature kind %q", path, feature.Kind)
			}
		}
		items = append(items, item)
		return nil
	})
	return items, err
}
func sortedKeys[V any](m map[string]V) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func aggregate(doc *tfjson.ProviderSchemas, items []semanticcoverage.Finding) output {
	result := output{Version: 1, AttributeScope: "top_level", Findings: len(items)}
	seen := make(map[string]bool)
	for _, address := range sortedKeys(doc.Schemas) {
		p := doc.Schemas[address]
		if p.ConfigSchema != nil {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindProvider, "", p.ConfigSchema.Block))
		}
		for _, name := range sortedKeys(p.ResourceSchemas) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindResource, name, p.ResourceSchemas[name].Block))
		}
		for _, name := range sortedKeys(p.DataSourceSchemas) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindDataSource, name, p.DataSourceSchemas[name].Block))
		}
		for _, name := range sortedKeys(p.EphemeralResourceSchemas) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindEphemeralResource, name, p.EphemeralResourceSchemas[name].Block))
		}
		for _, name := range sortedKeys(p.ListResourceSchemas) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindListResource, name, p.ListResourceSchemas[name].Block))
		}
		for _, name := range sortedKeys(p.ActionSchemas) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindAction, name, p.ActionSchemas[name].Block))
		}
		for _, name := range sortedKeys(p.StateStoreSchemas) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindStateStore, name, p.StateStoreSchemas[name].Block))
		}
		// Identity attributes are only observed in state, and the schema
		// doesn't say which of them are computed, so they are left out of
		// the attribute totals.
		for _, name := range sortedKeys(p.ResourceIdentitySchemas) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindIdentity, name, nil))
		}
		for _, name := range sortedKeys(p.Functions) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindFunction, name, nil))
		}
	}
	for _, item := range items {
		id := fmt.Sprintf("%s\x00%s\x00%d", item.Package, item.Test, item.Step)
		if seen[id] {
			result.Warnings = append(result.Warnings, fmt.Sprintf("duplicate semanticcoverage.Finding for %s step %d", item.Test, item.Step))
			continue
		}
		seen[id] = true
		if item.Outcome == semanticcoverage.OutcomeFailed {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s step %d failed; evidence excluded from coverage", item.Test, item.Step))
		}
		for _, warning := range item.Warnings {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s step %d: %s", item.Test, item.Step, warning))
		}
		for _, f := range item.Features {
			var matches []*target
			for _, t := range result.Targets {
				providerMatch := t.Kind == semanticcoverage.KindProvider && f.Name != "" && (t.Provider == f.Name || strings.HasSuffix(t.Provider, "/"+f.Name))
				providerSourceMatch := f.Provider == "" || t.Provider == f.Provider || !strings.Contains(f.Provider, "/") && strings.HasSuffix(t.Provider, "/"+f.Provider)
				if t.Kind == f.Kind && (t.Name == f.Name || providerMatch) && providerSourceMatch {
					matches = append(matches, t)
				}
			}
			if len(matches) != 1 {
				if f.Name != "terraform_data" && len(matches) > 1 {
					result.Warnings = append(result.Warnings, fmt.Sprintf("%s step %d: %s %s has %d schema matches", item.Test, item.Step, f.Kind, f.Name, len(matches)))
				}
				continue
			}
			t := matches[0]
			if item.Outcome == semanticcoverage.OutcomeFailed {
				continue
			}
			ref := evidence{Package: item.Package, Test: item.Test, Step: item.Step, Address: f.Address, Kind: "configured"}
			if len(f.Inputs) > 0 || len(f.Known) > 0 || len(f.Planned) > 0 || len(f.Completed) > 0 || f.Kind == semanticcoverage.KindFunction || f.Kind == semanticcoverage.KindIdentity {
				t.Evidence = append(t.Evidence, ref)
				if item.Outcome.Completed() {
					t.Exercised = true
				}
			}
			for _, action := range f.Planned {
				e := ref
				e.Kind = "planned"
				t.Planned[action] = append(t.Planned[action], e)
			}
			if item.Outcome == semanticcoverage.OutcomeSuccess {
				for _, operation := range f.Completed {
					e := ref
					e.Kind = "completed"
					t.Completed[operation] = append(t.Completed[operation], e)
				}
			}
			if !item.Outcome.Completed() {
				continue
			}
			for i := range t.Attributes {
				a := &t.Attributes[i]
				seen := false
				if a.Mode == "input" {
					seen = contains(f.Inputs, a.Path)
				}
				if a.Mode == "computed" {
					seen = item.Outcome == semanticcoverage.OutcomeSuccess && contains(f.Known, a.Path)
				}
				if a.Mode == "either" {
					seen = contains(f.Inputs, a.Path) || (item.Outcome == semanticcoverage.OutcomeSuccess && contains(f.Known, a.Path))
				}
				if seen {
					a.Covered = true
					e := ref
					if a.Mode == "computed" {
						e.Kind = "state_known"
					}
					a.Evidence = append(a.Evidence, e)
				}
			}
		}
	}
	for _, t := range result.Targets {
		if t.Kind != semanticcoverage.KindIdentity {
			result.Entities.Total++
			if t.Exercised {
				result.Entities.Covered++
			}
		}
		for _, a := range t.Attributes {
			if a.Mode == "computed" {
				result.Computed.Total++
				if a.Covered {
					result.Computed.Covered++
				}
			} else {
				result.Inputs.Total++
				if a.Covered {
					result.Inputs.Covered++
				}
			}
		}
	}
	sort.Slice(result.Targets, func(i, j int) bool {
		a, b := result.Targets[i], result.Targets[j]
		return a.Provider+string(a.Kind)+a.Name < b.Provider+string(b.Kind)+b.Name
	})
	sort.Strings(result.Warnings)
	return result
}
func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// newTarget returns a report target with the top-level attributes of b, sorted
// by path. A nil b gives a target without attributes.
func newTarget(provider string, kind semanticcoverage.Kind, name string, b *tfjson.SchemaBlock) *target {
	t := &target{Provider: provider, Kind: kind, Name: name, Planned: make(map[string][]evidence), Completed: make(map[string][]evidence)}
	flatten(&t.Attributes, "", b)
	sort.Slice(t.Attributes, func(i, j int) bool { return t.Attributes[i].Path < t.Attributes[j].Path })
	return t
}
func addAttribute(to *[]attributeResult, path string, a *tfjson.SchemaAttribute) {
	if a.Deprecated {
		return
	}
	mode := "input"
	if a.Computed && (a.Optional || a.Required) {
		mode = "either"
	} else if a.Computed {
		mode = "computed"
	}
	*to = append(*to, attributeResult{Path: path, Mode: mode})
}
func flatten(to *[]attributeResult, prefix string, b *tfjson.SchemaBlock) {
	if b == nil || b.Deprecated {
		return
	}
	for name, a := range b.Attributes {
		addAttribute(to, prefix+name, a)
	}
	for name := range b.NestedBlocks {
		*to = append(*to, attributeResult{Path: prefix + name, Mode: "input"})
	}
}
