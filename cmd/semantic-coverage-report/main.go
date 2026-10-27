// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Command semantic-coverage-report compares semantic coverage findings written
// by acceptance tests (TF_ACC_SEMANTIC_COVERAGE_DIR) with a
// `terraform providers schema -json` document.
package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
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
	Path string `json:"path"`
	// Configurable attributes are covered by configuration ("configured"
	// evidence): required or optional attributes, and blocks.
	Configurable bool `json:"configurable"`
	// Computed attributes are covered by a known value in state that the
	// step didn't configure ("state_known" evidence). Optional-computed
	// attributes are both, and need both: one test that sets them and one
	// that leaves them to the provider.
	Computed bool `json:"computed"`
	// Covered reports whether every side of the attribute has evidence.
	Covered  bool       `json:"covered"`
	Evidence []evidence `json:"evidence,omitempty"`
}

// has reports whether a has evidence of the given kind.
func (a attributeResult) has(kind string) bool {
	return slices.ContainsFunc(a.Evidence, func(e evidence) bool { return e.Kind == kind })
}

// missing returns the evidence kinds that a still needs to be covered.
func (a attributeResult) missing() []string {
	var kinds []string
	if a.Configurable && !a.has("configured") {
		kinds = append(kinds, "configured")
	}
	if a.Computed && !a.has("state_known") {
		kinds = append(kinds, "state_known")
	}
	return kinds
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
	Version        int    `json:"version"`
	AttributeScope string `json:"attribute_scope"`
	Findings       int    `json:"findings"`
	Entities       metric `json:"entities"`
	// Configurable counts configurable attributes covered by configuration.
	Configurable metric `json:"configurable"`
	// Computed counts computed attributes covered by a provider-computed
	// state value. Optional-computed attributes count in both totals.
	Computed metric    `json:"computed"`
	Targets  []*target `json:"targets"`
	Warnings []string  `json:"warnings,omitempty"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	findingsDir := flags.String("findings-dir", "findings", "directory of test-derived JSON findings")
	schemaFile := flags.String("schema-file", "", "terraform providers schema -json file")
	jsonOut := flags.String("json-out", "coverage.json", "JSON report path")
	textOut := flags.String("text-out", "coverage.txt", "text report path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *schemaFile == "" {
		return errors.New("-schema-file is required")
	}
	docBytes, err := os.ReadFile(*schemaFile)
	if err != nil {
		return err
	}
	// Unmarshalling also validates the schema format version.
	var doc tfjson.ProviderSchemas
	if err = json.Unmarshal(docBytes, &doc); err != nil {
		return fmt.Errorf("%s: %w", *schemaFile, err)
	}
	if len(doc.Schemas) == 0 {
		return errors.New("schema has no providers")
	}
	items, err := readFindings(*findingsDir)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("no step findings in %s", *findingsDir)
	}
	result := aggregate(&doc, items)
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(*jsonOut, append(data, '\n'), 0600); err != nil {
		return err
	}
	var text strings.Builder
	fmt.Fprintf(&text, "Test-derived semantic coverage: %d steps\nAttribute scope: top-level schema attributes and blocks only; nested descendants are excluded.\nEntities %d/%d; configurable %d/%d; computed %d/%d\n", result.Findings, result.Entities.Covered, result.Entities.Total, result.Configurable.Covered, result.Configurable.Total, result.Computed.Covered, result.Computed.Total)
	for _, t := range result.Targets {
		covered := 0
		for _, a := range t.Attributes {
			if a.Covered {
				covered++
			}
		}
		fmt.Fprintf(&text, "%s %s %s: exercised=%t attributes=%d/%d planned=%v completed=%v\n", t.Provider, t.Kind, t.Name, t.Exercised, covered, len(t.Attributes), slices.Sorted(maps.Keys(t.Planned)), slices.Sorted(maps.Keys(t.Completed)))
		for _, a := range t.Attributes {
			if !a.Covered {
				fmt.Fprintf(&text, "  uncovered %s (needs %s)\n", a.Path, strings.Join(a.missing(), ", "))
			}
		}
	}
	for _, w := range result.Warnings {
		fmt.Fprintf(&text, "warning: %s\n", w)
	}
	if err = os.WriteFile(*textOut, []byte(text.String()), 0600); err != nil {
		return err
	}
	fmt.Printf("read %d step findings; wrote %s and %s\n", len(items), *jsonOut, *textOut)
	return nil
}

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

func aggregate(doc *tfjson.ProviderSchemas, items []semanticcoverage.Finding) output {
	result := output{Version: 1, AttributeScope: "top_level", Findings: len(items)}
	for _, address := range slices.Sorted(maps.Keys(doc.Schemas)) {
		p := doc.Schemas[address]
		if p.ConfigSchema != nil {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindProvider, "", p.ConfigSchema.Block))
		}
		for _, name := range slices.Sorted(maps.Keys(p.ResourceSchemas)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindResource, name, p.ResourceSchemas[name].Block))
		}
		for _, name := range slices.Sorted(maps.Keys(p.DataSourceSchemas)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindDataSource, name, p.DataSourceSchemas[name].Block))
		}
		for _, name := range slices.Sorted(maps.Keys(p.EphemeralResourceSchemas)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindEphemeralResource, name, p.EphemeralResourceSchemas[name].Block))
		}
		for _, name := range slices.Sorted(maps.Keys(p.ListResourceSchemas)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindListResource, name, p.ListResourceSchemas[name].Block))
		}
		for _, name := range slices.Sorted(maps.Keys(p.ActionSchemas)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindAction, name, p.ActionSchemas[name].Block))
		}
		for _, name := range slices.Sorted(maps.Keys(p.StateStoreSchemas)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindStateStore, name, p.StateStoreSchemas[name].Block))
		}
		// Identity attributes are only observed in state and have no
		// configurable/computed distinction (only required/optional for
		// import), so they don't fit either attribute total and are left out.
		for _, name := range slices.Sorted(maps.Keys(p.ResourceIdentitySchemas)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindIdentity, name, nil))
		}
		for _, name := range slices.Sorted(maps.Keys(p.Functions)) {
			result.Targets = append(result.Targets, newTarget(address, semanticcoverage.KindFunction, name, nil))
		}
	}
	for _, item := range items {
		for _, warning := range item.Warnings {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s step %d: %s", item.Test, item.Step, warning))
		}
		// A failed step can stop anywhere, and an unexpected error means
		// the provider didn't do what the test expected, so none of its
		// evidence shows that a feature works. Other findings for the same
		// step still count.
		if item.Outcome == semanticcoverage.OutcomeFailed {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s step %d failed; evidence excluded from coverage", item.Test, item.Step))
			continue
		}
		// Match every feature first, so that the step's configured attributes
		// are known per target before any state evidence is recorded.
		type match struct {
			f semanticcoverage.Feature
			t *target
		}
		var matched []match
		configured := make(map[*target][]string)
		for _, f := range item.Features {
			var matches []*target
			for _, t := range result.Targets {
				if t.matches(f) {
					matches = append(matches, t)
				}
			}
			// Features of providers missing from the schema match nothing.
			// Several matches mean config-only evidence (no provider address)
			// for a type name that more than one provider in the schema
			// defines.
			if len(matches) != 1 {
				if len(matches) > 1 {
					result.Warnings = append(result.Warnings, fmt.Sprintf("%s step %d: %s %s has %d schema matches", item.Test, item.Step, f.Kind, f.Name, len(matches)))
				}
				continue
			}
			matched = append(matched, match{f, matches[0]})
			configured[matches[0]] = append(configured[matches[0]], f.Inputs...)
		}
		for _, m := range matched {
			f, t := m.f, m.t
			ref := evidence{Package: item.Package, Test: item.Test, Step: item.Step, Address: f.Address, Kind: "configured"}
			if len(f.Inputs) > 0 || len(f.Known) > 0 || len(f.Planned) > 0 || len(f.Completed) > 0 || f.Kind == semanticcoverage.KindFunction || f.Kind == semanticcoverage.KindIdentity {
				t.Evidence = append(t.Evidence, ref)
				t.Exercised = true
			}
			for _, action := range f.Planned {
				e := ref
				e.Kind = "planned"
				t.Planned[action] = append(t.Planned[action], e)
			}
			// An expected error leaves no completed operation and no reliable
			// state.
			success := item.Outcome == semanticcoverage.OutcomeSuccess
			if success {
				for _, operation := range f.Completed {
					e := ref
					e.Kind = "completed"
					t.Completed[operation] = append(t.Completed[operation], e)
				}
			}
			for i := range t.Attributes {
				a := &t.Attributes[i]
				if a.Configurable && slices.Contains(f.Inputs, a.Path) {
					a.Evidence = append(a.Evidence, ref)
				}
				// A configured value is echoed back into state, so a known
				// value only shows that the provider computed it if the step
				// didn't configure the attribute. Instances in a step can't be
				// paired reliably between configuration and state (count,
				// for_each, modules), so any configuration of the type in the
				// step rules it out.
				if a.Computed && success && slices.Contains(f.Known, a.Path) && !slices.Contains(configured[t], a.Path) {
					e := ref
					e.Kind = "state_known"
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
		for i := range t.Attributes {
			a := &t.Attributes[i]
			if a.Configurable {
				result.Configurable.Total++
				if a.has("configured") {
					result.Configurable.Covered++
				}
			}
			if a.Computed {
				result.Computed.Total++
				if a.has("state_known") {
					result.Computed.Covered++
				}
			}
			a.Covered = len(a.missing()) == 0
		}
	}
	slices.SortFunc(result.Targets, func(a, b *target) int {
		return cmp.Or(
			cmp.Compare(a.Provider, b.Provider),
			cmp.Compare(a.Kind, b.Kind),
			cmp.Compare(a.Name, b.Name),
		)
	})
	slices.Sort(result.Warnings)
	return result
}

// newTarget returns a report target measured against the top-level attributes
// of b. A nil b gives a target without attributes.
func newTarget(provider string, kind semanticcoverage.Kind, name string, b *tfjson.SchemaBlock) *target {
	return &target{
		Provider:   provider,
		Kind:       kind,
		Name:       name,
		Attributes: topLevelAttributes(b),
		Planned:    make(map[string][]evidence),
		Completed:  make(map[string][]evidence),
	}
}

// matches reports whether f is evidence for t.
func (t *target) matches(f semanticcoverage.Feature) bool {
	if t.Kind != f.Kind {
		return false
	}
	if t.Kind == semanticcoverage.KindProvider {
		// Provider blocks are labelled with the provider's local name.
		return f.Name != "" && providerMatches(t.Provider, f.Name)
	}
	return t.Name == f.Name && providerMatches(t.Provider, f.Provider)
}

// providerMatches reports whether ref refers to the provider with the given
// source address. ref is a full source address (from plans and state), a
// local name (from configuration), matched against the address's type
// segment, or empty, which matches any provider.
func providerMatches(address, ref string) bool {
	switch {
	case ref == "":
		return true
	case strings.Contains(ref, "/"):
		return address == ref
	default:
		return strings.HasSuffix(address, "/"+ref)
	}
}

// topLevelAttributes returns the top-level attributes and nested blocks of b,
// sorted by path, which are what attribute coverage is measured against.
// Nested attributes are left out: state can't tell their names apart from
// user-provided map keys. A nil b has none.
func topLevelAttributes(b *tfjson.SchemaBlock) []attributeResult {
	if b == nil {
		return nil
	}
	var result []attributeResult
	for name, a := range b.Attributes {
		result = append(result, attributeResult{Path: name, Configurable: a.Optional || a.Required, Computed: a.Computed})
	}
	for name := range b.NestedBlocks {
		result = append(result, attributeResult{Path: name, Configurable: true})
	}
	slices.SortFunc(result, func(a, b attributeResult) int { return strings.Compare(a.Path, b.Path) })
	return result
}
