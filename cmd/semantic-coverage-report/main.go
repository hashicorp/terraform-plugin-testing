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

	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
)

type schemaDoc struct {
	FormatVersion string                    `json:"format_version"`
	Providers     map[string]providerSchema `json:"provider_schemas"`
}
type providerSchema struct {
	Provider    *schema                    `json:"provider"`
	Resources   map[string]schema          `json:"resource_schemas"`
	DataSources map[string]schema          `json:"data_source_schemas"`
	Ephemeral   map[string]schema          `json:"ephemeral_resource_schemas"`
	Lists       map[string]schema          `json:"list_resource_schemas"`
	Actions     map[string]schema          `json:"action_schemas"`
	Stores      map[string]schema          `json:"state_store_schemas"`
	Identities  map[string]schema          `json:"resource_identity_schemas"`
	Functions   map[string]json.RawMessage `json:"functions"`
}
type schema struct {
	Block      *block               `json:"block"`
	Attributes map[string]attribute `json:"attributes"`
}
type block struct {
	Attributes map[string]attribute   `json:"attributes"`
	Blocks     map[string]nestedBlock `json:"block_types"`
	Deprecated bool                   `json:"deprecated"`
}
type nestedBlock struct {
	Nesting string `json:"nesting_mode"`
	Block   *block `json:"block"`
}
type attribute struct {
	Required   bool             `json:"required"`
	Optional   bool             `json:"optional"`
	Computed   bool             `json:"computed"`
	WriteOnly  bool             `json:"write_only"`
	Deprecated bool             `json:"deprecated"`
	Nested     *nestedAttribute `json:"nested_type"`
}
type nestedAttribute struct {
	Nesting    string               `json:"nesting_mode"`
	Attributes map[string]attribute `json:"attributes"`
}

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
	Kind       string                `json:"kind"`
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
	findings := flag.String("findings", "findings", "directory of test-derived JSON findings")
	schemaFile := flag.String("schema", "", "terraform providers schema -json file")
	jsonFile := flag.String("json", "coverage.json", "JSON report path")
	textFile := flag.String("text", "coverage.txt", "text report path")
	flag.Parse()
	if *schemaFile == "" {
		die(fmt.Errorf("-schema is required"))
	}
	docBytes, err := os.ReadFile(*schemaFile)
	if err != nil {
		die(err)
	}
	var doc schemaDoc
	if err = json.Unmarshal(docBytes, &doc); err != nil {
		die(err)
	}
	if len(doc.Providers) == 0 {
		die(fmt.Errorf("schema has no providers"))
	}
	if !strings.HasPrefix(doc.FormatVersion, "1.") {
		die(fmt.Errorf("unsupported provider schema format version %q", doc.FormatVersion))
	}
	items, err := readFindings(*findings)
	if err != nil {
		die(err)
	}
	if len(items) == 0 {
		die(fmt.Errorf("no step findings in %s", *findings))
	}
	result := aggregate(doc, items)
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		die(err)
	}
	if err = os.WriteFile(*jsonFile, append(data, '\n'), 0600); err != nil {
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
	if err = os.WriteFile(*textFile, []byte(text.String()), 0600); err != nil {
		die(err)
	}
	fmt.Printf("read %d step findings; wrote %s and %s\n", len(items), *jsonFile, *textFile)
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
			return fmt.Errorf("%s: unsupported semanticcoverage.Finding version %d", path, item.Version)
		}
		if item.Test == "" || item.Step < 1 {
			return fmt.Errorf("%s: missing test name or step", path)
		}
		if !item.Outcome.Valid() {
			return fmt.Errorf("%s: invalid outcome %q", path, item.Outcome)
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

func aggregate(doc schemaDoc, items []semanticcoverage.Finding) output {
	result := output{Version: 1, AttributeScope: "top_level", Findings: len(items)}
	seen := make(map[string]bool)
	for _, address := range sortedKeys(doc.Providers) {
		p := doc.Providers[address]
		register := func(kind, name string, s schema) {
			t := &target{Provider: address, Kind: kind, Name: name, Planned: make(map[string][]evidence), Completed: make(map[string][]evidence)}
			if s.Block != nil {
				flatten(&t.Attributes, "", s.Block)
			} else {
				for name, a := range s.Attributes {
					addAttribute(&t.Attributes, name, a)
				}
			}
			sort.Slice(t.Attributes, func(i, j int) bool { return t.Attributes[i].Path < t.Attributes[j].Path })
			result.Targets = append(result.Targets, t)
		}
		if p.Provider != nil {
			register("provider", "", *p.Provider)
		}
		for _, group := range []struct {
			kind    string
			members map[string]schema
		}{{"resource", p.Resources}, {"data_source", p.DataSources}, {"ephemeral_resource", p.Ephemeral}, {"list_resource", p.Lists}, {"action", p.Actions}, {"state_store", p.Stores}, {"identity", p.Identities}} {
			for _, name := range sortedKeys(group.members) {
				register(group.kind, name, group.members[name])
			}
		}
		for _, name := range sortedKeys(p.Functions) {
			register("function", name, schema{})
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
				providerMatch := t.Kind == "provider" && f.Name != "" && (t.Provider == f.Name || strings.HasSuffix(t.Provider, "/"+f.Name))
				providerSourceMatch := f.Provider == "" || t.Provider == f.Provider || !strings.Contains(f.Provider, "/") && strings.HasSuffix(t.Provider, "/"+f.Provider)
				if t.Kind == f.Kind && (t.Name == f.Name || providerMatch) && providerSourceMatch {
					matches = append(matches, t)
				}
			}
			if len(matches) != 1 {
				if f.Kind != "test_step" && f.Name != "terraform_data" && len(matches) > 1 {
					result.Warnings = append(result.Warnings, fmt.Sprintf("%s step %d: %s %s has %d schema matches", item.Test, item.Step, f.Kind, f.Name, len(matches)))
				}
				continue
			}
			t := matches[0]
			if item.Outcome == semanticcoverage.OutcomeSkipped {
				continue
			}
			if item.Outcome == semanticcoverage.OutcomeFailed {
				continue
			}
			ref := evidence{Package: item.Package, Test: item.Test, Step: item.Step, Address: f.Address, Kind: "configured"}
			if len(f.Inputs) > 0 || len(f.Known) > 0 || len(f.Planned) > 0 || len(f.Completed) > 0 || f.Kind == "function" || f.Kind == "identity" {
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
		if t.Kind != "identity" {
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
		return a.Provider+a.Kind+a.Name < b.Provider+b.Kind+b.Name
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
func addAttribute(to *[]attributeResult, path string, a attribute) {
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
func flatten(to *[]attributeResult, prefix string, b *block) {
	if b == nil || b.Deprecated {
		return
	}
	for name, a := range b.Attributes {
		addAttribute(to, prefix+name, a)
	}
	for name := range b.Blocks {
		*to = append(*to, attributeResult{Path: prefix + name, Mode: "input"})
	}
}
