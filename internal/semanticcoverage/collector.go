// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"time"
)

// Collector belongs to one TestCase. It records evidence for one step at a
// time.
//
// A nil *Collector is valid and means collection is disabled: New returns nil
// when EnvDirectory is unset, and every exported method is a no-op on a nil
// receiver (FinishStep returns nil, Enabled returns false). Callers can
// therefore invoke methods unconditionally and only need Enabled to skip work
// done solely to feed the collector, such as extra Terraform commands. New
// exported methods must preserve this by returning early when c == nil.
//
// A Collector is not safe for concurrent use. All calls come from the
// goroutine that runs the test steps: Terraform commands and checks run
// synchronously there, and provider servers never see the collector.
type Collector struct {
	directory   string
	packageName string
	testName    string
	stepNumber  int
	features    map[string]*featureEvidence
	warnings    set[Warning]
	active      bool
	// pendingApplyChanges holds the resource changes of the plan that
	// PreApplyPlan selected. ApplyCompleted marks them completed.
	pendingApplyChanges []Feature
	// pendingApplyActions holds the action invocations of the plan that
	// PreApplyPlan selected. ApplyCompleted marks them invoked.
	pendingApplyActions []Feature
	expectedError       bool
}

// Enabled reports whether findings are being collected. Callers use it to
// skip Terraform commands that only exist to feed the collector.
func (c *Collector) Enabled() bool { return c != nil }

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
	packageName := ""
	if wd, err := os.Getwd(); err == nil {
		packageName = filepath.Base(wd)
	}
	return &Collector{directory: dir, packageName: packageName, testName: testName}, nil
}

// BeginStep starts collecting evidence for a step. Any unfinished step is
// discarded, so callers must FinishStep or DiscardStep the previous step first.
func (c *Collector) BeginStep(stepNumber int) {
	if c == nil {
		return
	}
	c.stepNumber, c.active = stepNumber, true
	c.features = make(map[string]*featureEvidence)
	c.warnings = make(set[Warning])
	c.pendingApplyChanges = nil
	c.pendingApplyActions = nil
	c.expectedError = false
}

// FinishStep writes one private finding for the active test step, if any. A
// successful step that matched TestStep.ExpectError is written as
// OutcomeExpectedError.
func (c *Collector) FinishStep(outcome Outcome) error {
	if c == nil {
		return nil
	}
	if !c.active {
		return nil
	}
	c.active = false
	if c.expectedError && outcome == OutcomeSuccess {
		outcome = OutcomeExpectedError
	}
	finding := Finding{Version: FindingVersion, Package: c.packageName, Test: c.testName, Step: c.stepNumber, Outcome: outcome, Warnings: c.warnings.sorted()}
	for _, f := range c.features {
		finding.Features = append(finding.Features, f.finding())
	}
	sort.Slice(finding.Features, func(i, j int) bool {
		a, b := finding.Features[i], finding.Features[j]
		return string(a.Kind)+a.Name+a.Address < string(b.Kind)+b.Name+b.Address
	})
	if err := os.MkdirAll(c.directory, 0700); err != nil {
		return fmt.Errorf("coverage directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(c.directory, fmt.Sprintf("%d-%d-%d.json", time.Now().UnixNano(), os.Getpid(), c.stepNumber)), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(file).Encode(finding); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

// DiscardStep ends the active step, if any, without writing a finding. It is
// used for skipped steps, which contribute no coverage.
func (c *Collector) DiscardStep() {
	if c == nil {
		return
	}
	c.active = false
}

// featureEvidence is the evidence for one feature in the active step. It is
// converted to a Feature when the step is written.
type featureEvidence struct {
	kind      Kind
	name      string
	provider  string
	address   string
	inputs    set[string]
	known     set[string]
	null      set[string]
	planned   set[string]
	completed set[string]
}

func (f *featureEvidence) finding() Feature {
	return Feature{
		Kind:      f.kind,
		Name:      f.name,
		Provider:  f.provider,
		Address:   f.address,
		Inputs:    f.inputs.sorted(),
		Known:     f.known.sorted(),
		Null:      f.null.sorted(),
		Planned:   f.planned.sorted(),
		Completed: f.completed.sorted(),
	}
}

func (c *Collector) feature(kind Kind, name, address string) *featureEvidence {
	address = safeAddress(address)
	key := string(kind) + "\x00" + name + "\x00" + address
	f := c.features[key]
	if f == nil {
		f = &featureEvidence{
			kind:      kind,
			name:      name,
			address:   address,
			inputs:    make(set[string]),
			known:     make(set[string]),
			null:      make(set[string]),
			planned:   make(set[string]),
			completed: make(set[string]),
		}
		c.features[key] = f
	}
	return f
}

func (c *Collector) setProvider(kind Kind, name, address, provider string) {
	if provider != "" {
		c.feature(kind, name, address).provider = provider
	}
}

// set is a set of strings. The empty string is never a member.
type set[T ~string] map[T]struct{}

func (s set[T]) add(value T) {
	if value != "" {
		s[value] = struct{}{}
	}
}

// sorted returns the members of s in order, or nil if s is empty.
func (s set[T]) sorted() []T {
	return slices.Sorted(maps.Keys(s))
}

// Warning records that evidence of the given type may be incomplete for the
// active step.
func (c *Collector) Warning(code Warning) {
	if c == nil {
		return
	}
	if c.active {
		c.warnings.add(code)
	}
}

// MatchedExpectedError records that the active step failed with an error
// matching TestStep.ExpectError. FinishStep then writes OutcomeExpectedError
// instead of OutcomeSuccess.
func (c *Collector) MatchedExpectedError() {
	if c == nil {
		return
	}
	if c.active {
		c.expectedError = true
	}
}
