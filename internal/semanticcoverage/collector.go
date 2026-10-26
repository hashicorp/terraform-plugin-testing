// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"encoding/json"
	"fmt"
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
// receiver (Ready and FinishStep return nil, Enabled returns false). Callers can
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
	test        string
	step        int
	features    map[string]*Feature
	warnings    []Warning
	active      bool
	// lastPlan holds the resource changes of the plan that PreApplyPlan
	// selected. ApplyCompleted marks them completed.
	lastPlan []Feature
	// actions holds the action invocations of the plan that PreApplyPlan
	// selected. ApplyCompleted marks them invoked.
	actions       []Feature
	expectedError bool
}

// Enabled reports whether findings are being collected. Callers use it to
// skip Terraform commands that only exist to feed the collector.
func (c *Collector) Enabled() bool { return c != nil }

func New(test string) *Collector {
	dir := os.Getenv(EnvDirectory)
	if dir == "" {
		return nil
	}
	packageName := ""
	if wd, err := os.Getwd(); err == nil {
		packageName = filepath.Base(wd)
	}
	return &Collector{directory: dir, packageName: packageName, test: test}
}

// BeginStep starts collecting evidence for a step. Any unfinished step is
// discarded, so callers must FinishStep or DiscardStep the previous step first.
func (c *Collector) BeginStep(stepNumber int) {
	if c == nil {
		return
	}
	c.step, c.active = stepNumber, true
	c.features = make(map[string]*Feature)
	c.warnings = nil
	c.lastPlan = nil
	c.actions = nil
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
	finding := Finding{Version: FindingVersion, Package: c.packageName, Test: c.test, Step: c.step, Outcome: outcome, Warnings: c.warnings}
	for _, f := range c.features {
		for _, xs := range []*[]string{&f.Inputs, &f.Known, &f.Null, &f.Planned, &f.Completed} {
			sort.Strings(*xs)
		}
		finding.Features = append(finding.Features, *f)
	}
	sort.Slice(finding.Features, func(i, j int) bool {
		a, b := finding.Features[i], finding.Features[j]
		return a.Kind+a.Name+a.Address < b.Kind+b.Name+b.Address
	})
	slices.Sort(finding.Warnings)
	if err := os.MkdirAll(c.directory, 0700); err != nil {
		return fmt.Errorf("coverage directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(c.directory, fmt.Sprintf("%d-%d-%d.json", time.Now().UnixNano(), os.Getpid(), c.step)), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
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

// Ready creates the findings directory before running a test step.
func (c *Collector) Ready() error {
	if c == nil {
		return nil
	}
	if err := os.MkdirAll(c.directory, 0700); err != nil {
		return fmt.Errorf("coverage directory: %w", err)
	}
	return nil
}

func (c *Collector) feature(kind, name, address string) *Feature {
	address = safeAddress(address)
	key := kind + "\x00" + name + "\x00" + address
	f := c.features[key]
	if f == nil {
		f = &Feature{Kind: kind, Name: name, Address: address}
		c.features[key] = f
	}
	return f
}

func (c *Collector) setProvider(kind, name, address, provider string) {
	if provider != "" {
		c.feature(kind, name, address).Provider = provider
	}
}

func add[T ~string](xs *[]T, value T) {
	if value == "" {
		return
	}
	for _, x := range *xs {
		if x == value {
			return
		}
	}
	*xs = append(*xs, value)
}

// Warning records that evidence of the given type may be incomplete for the
// active step.
func (c *Collector) Warning(code Warning) {
	if c == nil {
		return
	}
	if c.active {
		add(&c.warnings, code)
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
