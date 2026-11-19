// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
)

// Collector belongs to one TestCase and starts a Step for each test step.
//
// A nil *Collector means collection is disabled: New returns nil for an
// empty directory, and BeginStep then returns a nil *Step.
type Collector struct {
	directory   string
	packageName string
	testName    string
}

// New returns a Collector that writes the findings of the test with the given
// name to dir, or nil if dir is empty. It creates dir, so that a bad
// directory fails the test before its first step.
func New(dir, testName string) (*Collector, error) {
	if dir == "" {
		return nil, nil
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("coverage directory: %w", err)
	}
	return &Collector{directory: dir, packageName: testPackage(), testName: testName}, nil
}

// testPackage returns the import path of the package under test, so that
// tests with the same name in different packages get distinct findings. A test
// binary's main package is the package's import path with a ".test" suffix.
// It returns "" for binaries without one, such as go test with file arguments.
func testPackage() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	path, ok := strings.CutSuffix(info.Path, ".test")
	if !ok || path == "command-line-arguments" {
		return ""
	}
	return path
}

// BeginStep starts collecting evidence for the test step with the given
// number.
func (c *Collector) BeginStep(number int) *Step {
	if c == nil {
		return nil
	}
	return &Step{collector: c, number: number, features: make(map[featureKey]*Feature)}
}

// BeginPostTestDestroy starts collecting evidence for the destroy that runs
// after the last test step. number is the number after the last step's.
func (c *Collector) BeginPostTestDestroy(number int) *Step {
	s := c.BeginStep(number)
	if s != nil {
		s.postTestDestroy = true
	}
	return s
}

// Step collects the evidence of one test step. End writes it as a Finding
// unless the test has failed; a step that is never ended writes nothing.
//
// A nil *Step is valid and means collection is disabled: every method is a
// no-op on a nil receiver (Enabled returns false). Callers can therefore
// invoke methods unconditionally and only need Enabled to skip work done
// solely to feed the step, such as extra Terraform commands. New methods must
// preserve this by returning early when s == nil.
//
// A Step is not safe for concurrent use. Only the goroutine that runs the
// test steps uses it.
type Step struct {
	collector       *Collector
	number          int
	postTestDestroy bool
	ended           bool
	served          []string
	features        map[featureKey]*Feature
	warnings        []Warning
}

// Enabled reports whether evidence is being collected.
func (s *Step) Enabled() bool { return s != nil }

// ObserveServedProviders records the source addresses under which the test
// harness serves provider factories from the test binary during the step.
// Terraform normalizes source addresses to lower case, so they are recorded
// in lower case.
//
// The report needs these addresses to match the step's evidence to the
// provider schema. The harness serves a provider under the hashicorp
// namespace by default, while the schema usually has the address the provider
// is published under, which can be in another namespace. The served addresses
// tell the report which providers are the code under test, and which schema
// provider each one is.
func (s *Step) ObserveServedProviders(addresses []string) {
	if s == nil {
		return
	}
	s.served = nil
	for _, address := range addresses {
		s.served = append(s.served, strings.ToLower(address))
	}
}

// End stops recording the step and writes its Finding, unless the test t has
// failed: a test that reports an error without stopping fails all of its
// remaining steps too. Only the first call has an effect. Errors are reported
// to t and logged.
func (s *Step) End(ctx context.Context, t testing.T) {
	t.Helper()

	if s == nil || s.ended {
		return
	}
	s.ended = true
	if t.Failed() {
		return
	}
	finding, err := s.finish()
	if err != nil {
		logging.HelperResourceError(ctx, "Error writing semantic coverage finding", map[string]interface{}{logging.KeyError: err})
		t.Errorf("Error writing semantic coverage finding: %s", err)
		return
	}
	for _, w := range finding.Warnings {
		logging.HelperResourceWarn(ctx, "Semantic coverage evidence incomplete", map[string]interface{}{logging.KeySemanticCoverageWarning: string(w)})
	}
	logging.HelperResourceDebug(ctx, "Wrote semantic coverage finding")
}

// finish writes the step's Finding and returns it.
func (s *Step) finish() (Finding, error) {
	finding := Finding{
		Version:         FindingVersion,
		Package:         s.collector.packageName,
		Test:            s.collector.testName,
		Step:            s.number,
		PostTestDestroy: s.postTestDestroy,
		ServedProviders: sortedSet(s.served),
		Warnings:        sortedSet(s.warnings),
	}
	for _, f := range s.features {
		f.Configured = sortedSet(f.Configured)
		f.Populated = sortedSet(f.Populated)
		f.Planned = sortedSet(f.Planned)
		f.Completed = sortedSet(f.Completed)
		finding.Features = append(finding.Features, *f)
	}
	slices.SortFunc(finding.Features, func(a, b Feature) int {
		return cmp.Or(cmp.Compare(a.Kind, b.Kind), cmp.Compare(a.Name, b.Name), cmp.Compare(a.Provider, b.Provider))
	})

	// Parallel tests share a process and can finish equally numbered steps
	// at once, so CreateTemp makes the name unique.
	file, err := os.CreateTemp(s.collector.directory, fmt.Sprintf("%d-*.json", s.number))
	if err != nil {
		return finding, err
	}
	if err := json.NewEncoder(file).Encode(finding); err != nil {
		_ = file.Close()
		return finding, err
	}
	return finding, file.Close()
}

// featureKey identifies a feature within a step. Configuration names no
// provider address, so the features of one type are merged across
// configuration, plans, and state. Functions are the exception: a call names
// its provider, so equally named functions of two providers stay apart.
type featureKey struct {
	kind     Kind
	name     string
	provider string
}

// feature returns the step's feature of the given kind and name, creating it
// if needed. A non-empty provider is recorded.
func (s *Step) feature(kind Kind, name, provider string) *Feature {
	key := featureKey{kind: kind, name: name}
	if kind == KindFunction {
		key.provider = provider
	}
	f := s.features[key]
	if f == nil {
		f = &Feature{Kind: kind, Name: name}
		s.features[key] = f
	}
	if provider != "" {
		f.Provider = provider
	}
	return f
}

// warn records that evidence of the given type may be incomplete.
func (s *Step) warn(w Warning) {
	s.warnings = append(s.warnings, w)
}

// appendNonEmpty appends value to the list, unless it is empty. finish
// removes duplicates.
func appendNonEmpty[T ~string](list *[]T, value T) {
	if value != "" {
		*list = append(*list, value)
	}
}

// sortedSet sorts values and removes duplicates.
func sortedSet[T ~string](values []T) []T {
	slices.Sort(values)
	return slices.Compact(values)
}
