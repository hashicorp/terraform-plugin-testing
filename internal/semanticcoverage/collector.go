// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"
)

// Collector belongs to one TestCase and starts a Step for each test step.
//
// A nil *Collector means collection is disabled: New returns nil when
// EnvDirectory is unset, and BeginStep then returns a nil *Step.
type Collector struct {
	directory   string
	packageName string
	testName    string
}

// New returns a Collector for the test with the given name, or nil if
// EnvDirectory is unset. It creates the findings directory, so that a bad
// directory fails the test before its first step.
func New(testName string) (*Collector, error) {
	dir := os.Getenv(EnvDirectory)
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
// number. expectError reports whether the step sets TestStep.ExpectError.
func (c *Collector) BeginStep(number int, expectError bool) *Step {
	if c == nil {
		return nil
	}
	return &Step{collector: c, number: number, expectError: expectError, features: make(map[featureKey]*Feature)}
}

// Step collects the evidence of one test step. Finish writes it as a Finding;
// a step that is never finished, such as a skipped one, writes nothing.
//
// A nil *Step is valid and means collection is disabled: every method is a
// no-op on a nil receiver (Finish returns nil, Enabled returns false). Callers
// can therefore invoke methods unconditionally and only need Enabled to skip
// work done solely to feed the step, such as extra Terraform commands. New
// methods must preserve this by returning early when s == nil.
//
// A Step is not safe for concurrent use. All calls come from the goroutine
// that runs the test steps: Terraform commands and checks run synchronously
// there, and provider servers never see the step.
type Step struct {
	collector   *Collector
	number      int
	expectError bool
	served      []string
	features    map[featureKey]*Feature
	warnings    []Warning
}

// Enabled reports whether evidence is being collected.
func (s *Step) Enabled() bool { return s != nil }

// ObserveServedProviders records the source addresses under which the test
// harness serves provider factories from the test binary during the step.
// Terraform normalizes source addresses to lower case, so they are recorded
// in lower case.
func (s *Step) ObserveServedProviders(addresses []string) {
	if s == nil {
		return
	}
	s.served = nil
	for _, address := range addresses {
		s.served = append(s.served, strings.ToLower(address))
	}
}

// Finish writes the step's Finding. failed reports whether the test failed
// during the step. A step that didn't fail and sets TestStep.ExpectError
// passed only because the error matched, so it is written as
// OutcomeExpectedError.
func (s *Step) Finish(failed bool) error {
	if s == nil {
		return nil
	}
	outcome := OutcomeSuccess
	switch {
	case failed:
		outcome = OutcomeFailed
	case s.expectError:
		outcome = OutcomeExpectedError
	}
	finding := Finding{
		Version:         FindingVersion,
		Package:         s.collector.packageName,
		Test:            s.collector.testName,
		Step:            s.number,
		Outcome:         outcome,
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
		return err
	}
	if err := json.NewEncoder(file).Encode(finding); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
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

// add appends value to the list, unless it is empty. Finish removes
// duplicates.
func add[T ~string](list *[]T, value T) {
	if value != "" {
		*list = append(*list, value)
	}
}

// sortedSet sorts values and removes duplicates.
func sortedSet[T ~string](values []T) []T {
	slices.Sort(values)
	return slices.Compact(values)
}
