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

// reportVersion is the schema version of the JSON report. It is independent
// of semanticcoverage.FindingVersion, which versions the findings that the
// report reads.
const reportVersion = 1

type output struct {
	Version        int    `json:"version"`
	AttributeScope string `json:"attribute_scope"`
	Findings       int    `json:"findings"`
	Totals         totals `json:"totals"`
	// Features lists every feature in the schema, sorted by provider, kind,
	// and name.
	Features []*featureResult `json:"features"`
	Warnings []string         `json:"warnings,omitempty"`
}

type totals struct {
	Features metric `json:"features"`
	// Configurable counts configurable attributes set by a test.
	Configurable metric `json:"configurable"`
	// Computed counts computed attributes that a test left to the provider.
	// Optional-computed attributes count in both totals.
	Computed metric `json:"computed"`
}

type metric struct {
	Covered int `json:"covered"`
	Total   int `json:"total"`
}

func (m *metric) add(covered bool) {
	m.Total++
	if covered {
		m.Covered++
	}
}

// featureResult is the coverage of one feature in the schema.
type featureResult struct {
	// Provider is the provider's source address.
	Provider string                `json:"provider"`
	Kind     semanticcoverage.Kind `json:"kind"`
	Name     string                `json:"name"`
	// Covered reports whether the feature appears in a finding that didn't
	// fail.
	Covered    bool                                     `json:"covered"`
	CoveredBy  []stepRef                                `json:"covered_by,omitempty"`
	Attributes []attributeResult                        `json:"attributes,omitempty"`
	Planned    map[semanticcoverage.Operation][]stepRef `json:"planned,omitempty"`
	Completed  map[semanticcoverage.Operation][]stepRef `json:"completed,omitempty"`
}

// attributeResult is the coverage of one top-level attribute or block.
type attributeResult struct {
	Path string `json:"path"`
	// Configurable attributes (required, optional, and blocks) need a test
	// that sets them.
	Configurable bool `json:"configurable"`
	// Computed attributes need a test that leaves them to the provider and
	// finds them populated. Optional-computed attributes are both, and need
	// both tests.
	Computed bool `json:"computed"`
	// Covered reports whether every side of the attribute has evidence.
	Covered      bool      `json:"covered"`
	ConfiguredBy []stepRef `json:"configured_by,omitempty"`
	ComputedBy   []stepRef `json:"computed_by,omitempty"`
}

// missing describes the tests that a still needs before it is covered.
func (a attributeResult) missing() []string {
	var tests []string
	if a.Configurable && len(a.ConfiguredBy) == 0 {
		tests = append(tests, "a test that sets it")
	}
	if a.Computed && len(a.ComputedBy) == 0 {
		tests = append(tests, "a test that leaves it to the provider")
	}
	return tests
}

// stepRef identifies the test step of a finding.
type stepRef struct {
	Package string `json:"package,omitempty"`
	Test    string `json:"test"`
	Step    int    `json:"step"`
}

func (r stepRef) String() string {
	if r.Package == "" {
		return fmt.Sprintf("%s step %d", r.Test, r.Step)
	}
	return fmt.Sprintf("%s %s step %d", r.Package, r.Test, r.Step)
}

// appendRef appends ref to refs unless it is already there. A step can have
// several findings for the same feature, such as a provider named both by its
// local name and by its source address. aggregate adds findings in step order,
// so a duplicate can only be the last ref.
func appendRef(refs []stepRef, ref stepRef) []stepRef {
	if n := len(refs); n > 0 && refs[n-1] == ref {
		return refs
	}
	return append(refs, ref)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	findingsDir := flags.String("findings-dir", "findings", "directory of test-derived JSON findings")
	schemaFile := flags.String("schema-file", "", "terraform providers schema -json file")
	jsonOut := flags.String("json-out", "coverage.json", "JSON report path")
	textOut := flags.String("text-out", "coverage.txt", "text report path")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *schemaFile == "" {
		return errors.New("-schema-file is required")
	}
	doc, err := readSchema(*schemaFile)
	if err != nil {
		return err
	}
	findings, err := readFindings(*findingsDir)
	if err != nil {
		return err
	}
	if len(findings) == 0 {
		return fmt.Errorf("no findings in %s", *findingsDir)
	}
	result := aggregate(doc, findings)
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(*jsonOut, append(data, '\n'), 0600); err != nil {
		return err
	}
	if err = os.WriteFile(*textOut, []byte(textReport(result)), 0600); err != nil {
		return err
	}
	fmt.Printf("read %d findings; wrote %s and %s\n", len(findings), *jsonOut, *textOut)
	return nil
}

func readSchema(path string) (*tfjson.ProviderSchemas, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Unmarshalling also validates the schema format version.
	var doc tfjson.ProviderSchemas
	if err = json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(doc.Schemas) == 0 {
		return nil, fmt.Errorf("%s: schema has no providers", path)
	}
	return &doc, nil
}

func readFindings(dir string) ([]semanticcoverage.Finding, error) {
	var findings []semanticcoverage.Finding
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
		var finding semanticcoverage.Finding
		if err = json.Unmarshal(data, &finding); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if err = finding.Validate(); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		findings = append(findings, finding)
		return nil
	})
	return findings, err
}

func textReport(r output) string {
	var text strings.Builder
	fmt.Fprintf(&text, "Test-derived semantic coverage: %d findings\n", r.Findings)
	text.WriteString("Attribute scope: top-level schema attributes and blocks only; nested descendants are excluded.\n")
	fmt.Fprintf(&text, "Features %d/%d; configurable %d/%d; computed %d/%d\n",
		r.Totals.Features.Covered, r.Totals.Features.Total,
		r.Totals.Configurable.Covered, r.Totals.Configurable.Total,
		r.Totals.Computed.Covered, r.Totals.Computed.Total)
	for _, f := range r.Features {
		covered := 0
		for _, a := range f.Attributes {
			if a.Covered {
				covered++
			}
		}
		fmt.Fprintf(&text, "%s %s %s: covered=%t attributes=%d/%d planned=%v completed=%v\n",
			f.Provider, f.Kind, f.Name, f.Covered, covered, len(f.Attributes),
			slices.Sorted(maps.Keys(f.Planned)), slices.Sorted(maps.Keys(f.Completed)))
		for _, a := range f.Attributes {
			if !a.Covered {
				fmt.Fprintf(&text, "  uncovered %s (needs %s)\n", a.Path, strings.Join(a.missing(), " and "))
			}
		}
	}
	for _, w := range r.Warnings {
		fmt.Fprintf(&text, "warning: %s\n", w)
	}
	return text.String()
}

func aggregate(doc *tfjson.ProviderSchemas, findings []semanticcoverage.Finding) output {
	features := schemaFeatures(doc)
	// Sorting the findings puts the warnings in step order.
	findings = slices.Clone(findings)
	slices.SortStableFunc(findings, func(a, b semanticcoverage.Finding) int {
		return cmp.Or(cmp.Compare(a.Package, b.Package), cmp.Compare(a.Test, b.Test), cmp.Compare(a.Step, b.Step))
	})
	providers := newProviderMap(doc)
	var warnings []string
	for _, f := range findings {
		warnings = append(warnings, addFinding(features, providers, f)...)
	}
	warnings = append(warnings, providers.warnings()...)
	return summarize(features, warnings, len(findings))
}

// schemaFeatures returns a featureResult for every feature in doc.
func schemaFeatures(doc *tfjson.ProviderSchemas) []*featureResult {
	var features []*featureResult
	for provider, p := range doc.Schemas {
		add := func(kind semanticcoverage.Kind, name string, b *tfjson.SchemaBlock) {
			features = append(features, &featureResult{
				Provider:   provider,
				Kind:       kind,
				Name:       name,
				Attributes: topLevelAttributes(b),
				Planned:    make(map[semanticcoverage.Operation][]stepRef),
				Completed:  make(map[semanticcoverage.Operation][]stepRef),
			})
		}
		if p.ConfigSchema != nil {
			add(semanticcoverage.KindProvider, "", p.ConfigSchema.Block)
		}
		for kind, schemas := range map[semanticcoverage.Kind]map[string]*tfjson.Schema{
			semanticcoverage.KindResource:          p.ResourceSchemas,
			semanticcoverage.KindDataSource:        p.DataSourceSchemas,
			semanticcoverage.KindEphemeralResource: p.EphemeralResourceSchemas,
			semanticcoverage.KindListResource:      p.ListResourceSchemas,
			semanticcoverage.KindStateStore:        p.StateStoreSchemas,
		} {
			for name, schema := range schemas {
				add(kind, name, schema.Block)
			}
		}
		for name, schema := range p.ActionSchemas {
			add(semanticcoverage.KindAction, name, schema.Block)
		}
		for name := range p.Functions {
			add(semanticcoverage.KindFunction, name, nil)
		}
	}
	slices.SortFunc(features, func(a, b *featureResult) int {
		return cmp.Or(
			cmp.Compare(a.Provider, b.Provider),
			cmp.Compare(a.Kind, b.Kind),
			cmp.Compare(a.Name, b.Name),
		)
	})
	return features
}

// addFinding records the evidence of finding f in features and returns the
// finding's warnings.
func addFinding(features []*featureResult, providers *providerMap, f semanticcoverage.Finding) []string {
	ref := stepRef{Package: f.Package, Test: f.Test, Step: f.Step}
	var warnings []string
	for _, w := range f.Warnings {
		warnings = append(warnings, fmt.Sprintf("%s: %s", ref, w))
	}
	served := providers.resolve(f.ServedProviders)
	for _, feature := range f.Features {
		feature, ok := inProcess(feature, served)
		if !ok {
			// The feature belongs to a provider that Terraform installed,
			// such as an external provider or the released provider in an
			// upgrade step, not to the code under test.
			continue
		}
		var matches []*featureResult
		for _, r := range features {
			if r.matches(feature) {
				matches = append(matches, r)
			}
		}
		switch len(matches) {
		case 0:
			// The feature belongs to a provider that isn't in the schema,
			// such as another provider that the test uses.
		case 1:
			matches[0].record(ref, feature)
		default:
			// Configuration names no provider address, so a type name
			// that several providers in the schema define is ambiguous.
			warnings = append(warnings, fmt.Sprintf("%s: %s %s has %d schema matches", ref, feature.Kind, feature.Name, len(matches)))
		}
	}
	return warnings
}

// record records the evidence f of the step ref.
func (r *featureResult) record(ref stepRef, f semanticcoverage.Feature) {
	r.CoveredBy = appendRef(r.CoveredBy, ref)
	for _, operation := range f.Planned {
		r.Planned[operation] = appendRef(r.Planned[operation], ref)
	}
	for _, operation := range f.Completed {
		r.Completed[operation] = appendRef(r.Completed[operation], ref)
	}
	for i := range r.Attributes {
		a := &r.Attributes[i]
		configured := slices.Contains(f.Configured, a.Path)
		if a.Configurable && configured {
			a.ConfiguredBy = appendRef(a.ConfiguredBy, ref)
		}
		// A configured value is copied into state, so a populated value
		// shows that the provider computed it only if the step didn't
		// configure the attribute. Configuration and state can't be paired
		// per instance (count, for_each, modules), so the collector merges
		// all uses of a type in a step. Configuring the attribute on any of
		// them therefore rules out the step.
		if a.Computed && !configured && slices.Contains(f.Populated, a.Path) {
			a.ComputedBy = appendRef(a.ComputedBy, ref)
		}
	}
}

// summarize sets the covered flags and returns the report.
func summarize(features []*featureResult, warnings []string, findings int) output {
	result := output{Version: reportVersion, AttributeScope: "top_level", Findings: findings, Features: features, Warnings: warnings}
	for _, f := range features {
		f.Covered = len(f.CoveredBy) > 0
		result.Totals.Features.add(f.Covered)
		for i := range f.Attributes {
			a := &f.Attributes[i]
			if a.Configurable {
				result.Totals.Configurable.add(len(a.ConfiguredBy) > 0)
			}
			if a.Computed {
				result.Totals.Computed.add(len(a.ComputedBy) > 0)
			}
			a.Covered = len(a.missing()) == 0
		}
	}
	return result
}

// matches reports whether f is evidence for r.
func (r *featureResult) matches(f semanticcoverage.Feature) bool {
	if r.Kind != f.Kind {
		return false
	}
	if r.Kind == semanticcoverage.KindProvider {
		// Provider blocks are labelled with the provider's local name.
		return f.Name != "" && providerMatches(r.Provider, f.Name)
	}
	return r.Name == f.Name && providerMatches(r.Provider, f.Provider)
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

// providerMap maps the provider addresses under which the test harness served
// provider factories from the test binary to the schema's provider addresses.
//
// The harness serves a factory as registry.terraform.io/hashicorp/<name>
// unless TF_ACC_PROVIDER_HOST or TF_ACC_PROVIDER_NAMESPACE say otherwise,
// while the schema is usually generated under the address the provider is
// published under. A served address maps to the schema provider with the same
// address or, failing that, to the only schema provider of the same type.
type providerMap struct {
	// byType lists the schema's provider addresses by provider type.
	byType map[string][]string
	// served records every served address of all findings.
	served map[string]bool
	// mapped reports whether some served address maps to the schema.
	mapped bool
}

func newProviderMap(doc *tfjson.ProviderSchemas) *providerMap {
	m := &providerMap{byType: make(map[string][]string), served: make(map[string]bool)}
	for address := range doc.Schemas {
		typ := providerType(address)
		m.byType[typ] = append(m.byType[typ], address)
	}
	return m
}

// resolve returns the schema address of each served address of a step. A
// served provider that no single schema provider matches, such as a provider
// that only the tests use, maps to "".
func (m *providerMap) resolve(served []string) map[string]string {
	result := make(map[string]string, len(served))
	for _, address := range served {
		m.served[address] = true
		candidates := m.byType[providerType(address)]
		var target string
		switch {
		case slices.Contains(candidates, address):
			target = address
		case len(candidates) == 1:
			target = candidates[0]
		}
		result[address] = target
		m.mapped = m.mapped || target != ""
	}
	return result
}

// warnings returns a warning if no served provider maps to the schema.
func (m *providerMap) warnings() []string {
	if m.mapped {
		return nil
	}
	served := "none"
	if len(m.served) > 0 {
		served = strings.Join(slices.Sorted(maps.Keys(m.served)), ", ")
	}
	return []string{fmt.Sprintf("no provider served from the test binary matches a single schema provider (served: %s); generate the schema from the provider under test", served)}
}

// inProcess reports whether feature f belongs to a provider that the step
// served from the test binary, and returns f with that provider's source
// address replaced by its schema address. served is the result of
// providerMap.resolve for the step.
func inProcess(f semanticcoverage.Feature, served map[string]string) (semanticcoverage.Feature, bool) {
	ref := &f.Provider
	if f.Kind == semanticcoverage.KindProvider {
		ref = &f.Name
	}
	if strings.Contains(*ref, "/") {
		address, ok := served[*ref]
		if address != "" {
			*ref = address
		}
		return f, ok
	}
	// A local name, or nothing for configuration, where Terraform implies
	// the provider from the type name's prefix.
	typ := *ref
	if typ == "" {
		typ, _, _ = strings.Cut(f.Name, "_")
	}
	for address := range served {
		if providerType(address) == typ {
			return f, true
		}
	}
	return f, false
}

// providerType returns the type of a provider source address, such as
// "aws" for registry.terraform.io/hashicorp/aws.
func providerType(address string) string {
	return address[strings.LastIndex(address, "/")+1:]
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
