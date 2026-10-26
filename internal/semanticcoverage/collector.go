// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Package semanticcoverage records structural evidence from acceptance-test steps.
// It never serializes configuration or Terraform values.
package semanticcoverage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfjson "github.com/hashicorp/terraform-json"
)

const EnvDirectory = "TF_ACC_SEMANTIC_COVERAGE_DIR"

// FindingVersion is the schema version of the JSON findings written per step.
//
// Findings live in a user-owned directory (EnvDirectory) that outlives the
// test run. They are written by the test binary and read by
// cmd/semantic-coverage-report, and the two can be built from different
// versions of this module. Bump FindingVersion whenever the meaning or shape
// of Finding changes incompatibly, so that the report rejects stale findings
// instead of misreading them.
const FindingVersion = 1

// Terraform instance keys can contain arbitrary strings, including ']'.
// Never persist their contents or the suffix of a malformed address.
func safeAddress(address string) string {
	var result strings.Builder
	for i := 0; i < len(address); i++ {
		if address[i] != '[' {
			result.WriteByte(address[i])
			continue
		}
		result.WriteString("[*]")
		quoted, escaped := false, false
		closed := false
		for i++; i < len(address); i++ {
			switch {
			case escaped:
				escaped = false
			case quoted && address[i] == '\\':
				escaped = true
			case address[i] == '"':
				quoted = !quoted
			case address[i] == ']' && !quoted:
				closed = true
			}
			if closed {
				break
			}
		}
		if !closed {
			break
		}
	}
	return result.String()
}

type Feature struct {
	Kind      string   `json:"kind"`
	Name      string   `json:"name"`
	Provider  string   `json:"provider,omitempty"`
	Address   string   `json:"address,omitempty"`
	Inputs    []string `json:"inputs,omitempty"`
	Known     []string `json:"known,omitempty"`
	Null      []string `json:"null,omitempty"`
	Planned   []string `json:"planned,omitempty"`
	Completed []string `json:"completed,omitempty"`
}

// Outcome is the final status of a test step, as written to a Finding.
type Outcome string

const (
	// OutcomeSuccess means the step completed and passed its checks.
	OutcomeSuccess Outcome = "success"
	// OutcomeExpectedError means the step failed with an error matching
	// TestStep.ExpectError.
	OutcomeExpectedError Outcome = "expected_error"
	// OutcomeFailed means the test failed during the step.
	OutcomeFailed Outcome = "failed"
)

// Valid reports whether o is one of the known outcomes.
func (o Outcome) Valid() bool {
	switch o {
	case OutcomeSuccess, OutcomeExpectedError, OutcomeFailed:
		return true
	}
	return false
}

// Completed reports whether the step ran to completion, either passing or
// failing with the expected error.
func (o Outcome) Completed() bool {
	return o == OutcomeSuccess || o == OutcomeExpectedError
}

// Warning means that one type of evidence for a step may be incomplete. The
// codes don't say why; use the logs or the collector tests for that.
type Warning string

const (
	// WarningConfigUnavailable means some configuration could not be read or
	// parsed, so input and feature evidence may be missing.
	WarningConfigUnavailable Warning = "config_unavailable"
	// WarningPlanUnavailable means a plan could not be read, so planned and
	// completed operation evidence may be missing.
	WarningPlanUnavailable Warning = "plan_unavailable"
	// WarningStateUnavailable means the post-operation state could not be
	// read, so known and null attribute evidence may be missing.
	WarningStateUnavailable Warning = "state_unavailable"
)

// Valid reports whether w is one of the known warnings.
func (w Warning) Valid() bool {
	switch w {
	case WarningConfigUnavailable, WarningPlanUnavailable, WarningStateUnavailable:
		return true
	}
	return false
}

type Finding struct {
	Version  int       `json:"version"`
	Package  string    `json:"package,omitempty"`
	Test     string    `json:"test"`
	Step     int       `json:"step"`
	Outcome  Outcome   `json:"outcome"`
	Features []Feature `json:"features,omitempty"`
	Warnings []Warning `json:"warnings,omitempty"`
}

// Collector belongs to one TestCase, and is shared by its import working directory.
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
	actions []Feature
	// latestPlan is the plan that Plan observed most recently.
	latestPlan    *tfjson.Plan
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
	c.latestPlan = nil
	c.expectedError = false
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

// Config reads the files actually copied into the test working directory.
func (c *Collector) Config(dir string) {
	if c == nil {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		c.Warning(WarningConfigUnavailable)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".tf") && !strings.HasSuffix(name, ".tf.json") && !strings.HasSuffix(name, ".tfquery.hcl") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			c.Warning(WarningConfigUnavailable)
			continue
		}
		c.Parse(name, data)
	}
}

func (c *Collector) QueryConfig(dir string) {
	if c == nil {
		return
	}
	data, err := os.ReadFile(filepath.Join(dir, "terraform_plugin_test.tfquery.hcl"))
	if err != nil {
		c.Warning(WarningConfigUnavailable)
		return
	}
	c.Parse("terraform_plugin_test.tfquery.hcl", data)
}

func (c *Collector) Parse(name string, data []byte) {
	if c == nil {
		return
	}
	if strings.HasSuffix(name, ".json") {
		var object map[string]any
		if json.Unmarshal(data, &object) != nil {
			c.Warning(WarningConfigUnavailable)
			return
		}
		if c.active {
			c.jsonBlocks(object)
		}
		return
	}
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(data, name)
	if diags.HasErrors() {
		c.Warning(WarningConfigUnavailable)
		return
	}
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		c.Warning(WarningConfigUnavailable)
		return
	}
	if c.active {
		c.blocks(body.Blocks)
		c.functions(body)
	}
}

func (c *Collector) blocks(blocks []*hclsyntax.Block) {
	for _, b := range blocks {
		kind, name, address := "", "", ""
		switch b.Type {
		case "provider":
			if len(b.Labels) > 0 && !strings.HasSuffix(b.DefRange().Filename, ".tfquery.hcl") {
				kind, name = "provider", b.Labels[0]
			}
		case "resource", "data", "ephemeral", "action", "list":
			if len(b.Labels) > 1 {
				kind, name, address = b.Type, b.Labels[0], b.Labels[0]+"."+b.Labels[1]
				switch kind {
				case "data":
					kind, address = "data_source", "data."+address
				case "list":
					kind, address = "list_resource", "list."+address
				case "ephemeral":
					kind, address = "ephemeral_resource", "ephemeral."+address
				case "action":
					address = "action." + address
				}
			}
		case "state_store":
			if len(b.Labels) > 0 {
				kind, name = "state_store", b.Labels[0]
			}
		}
		if kind != "" {
			f := c.feature(kind, name, address)
			switch kind {
			case "action", "list_resource", "state_store":
				// Schema attributes live in a nested config block.
				for _, nested := range b.Body.Blocks {
					if nested.Type == "config" {
						walkInputs(f, nested.Body, "")
					}
				}
			default:
				walkInputs(f, b.Body, "")
			}
		}
		c.functions(b.Body)
		c.blocks(b.Body.Blocks)
	}
}

func walkInputs(f *Feature, body *hclsyntax.Body, prefix string) {
	for name := range body.Attributes {
		path := prefix + name
		add(&f.Inputs, path)
	}
	for _, b := range body.Blocks {
		path := prefix + b.Type
		add(&f.Inputs, path)
		walkInputs(f, b.Body, path+".")
	}
}

func (c *Collector) functions(body *hclsyntax.Body) {
	for _, attr := range body.Attributes {
		diags := hclsyntax.VisitAll(attr.Expr, func(node hclsyntax.Node) hcl.Diagnostics {
			if call, ok := node.(*hclsyntax.FunctionCallExpr); ok && strings.HasPrefix(call.Name, "provider::") {
				parts := strings.Split(call.Name, "::")
				if len(parts) == 3 {
					c.feature("function", parts[2], "").Provider = parts[1]
				}
			}
			return nil
		})
		if diags.HasErrors() {
			add(&c.warnings, WarningConfigUnavailable)
		}
	}
}

func (c *Collector) jsonBlocks(root map[string]any) {
	for rawKind, value := range root {
		kind := rawKind
		if kind == "data" {
			kind = "data_source"
		}
		if kind == "list" {
			kind = "list_resource"
		}
		if kind == "ephemeral" {
			kind = "ephemeral_resource"
		}
		if kind == "resource" || kind == "data_source" || kind == "ephemeral_resource" || kind == "action" || kind == "list_resource" || kind == "provider" || kind == "state_store" {
			for name, raw := range asMap(value) {
				if kind == "provider" {
					f := c.feature(kind, name, "")
					jsonInputs(f, asMap(raw))
					continue
				}
				if kind == "state_store" {
					jsonInputs(c.feature(kind, name, ""), asMap(asMap(raw)["config"]))
					continue
				}
				for label, cfg := range asMap(raw) {
					address := name + "." + label
					if kind == "data_source" {
						address = "data." + address
					}
					if kind == "ephemeral_resource" {
						address = "ephemeral." + address
					}
					if kind == "action" || kind == "list_resource" {
						cfg = asMap(cfg)["config"]
					}
					if kind == "action" {
						address = "action." + address
					}
					if kind == "list_resource" {
						address = "list." + name + "." + label
					}
					jsonInputs(c.feature(kind, name, address), asMap(cfg))
				}
			}
		}
	}
}

func asMap(value any) map[string]any { m, _ := value.(map[string]any); return m }
func jsonInputs(f *Feature, object map[string]any) {
	for name := range object {
		add(&f.Inputs, name)
	}
}

// Plan records the planned operation of each resource change and action
// invocation in plan. A nil plan means the saved plan could not be read.
func (c *Collector) Plan(plan *tfjson.Plan) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	c.latestPlan = plan
	if plan == nil {
		add(&c.warnings, WarningPlanUnavailable)
		return
	}
	if plan.Config != nil {
		for _, provider := range plan.Config.ProviderConfigs {
			if provider != nil && provider.FullName != "" {
				c.feature("provider", provider.FullName, "")
			}
		}
	}
	// Several plans run during one step. Only the pre-apply plan can justify
	// completed actions; follow-up no-op plans must not replace that evidence.
	for _, rc := range plan.ResourceChanges {
		if rc == nil || rc.Change == nil {
			continue
		}
		kind := "resource"
		if rc.Mode == tfjson.DataResourceMode {
			kind = "data_source"
		} else if rc.Mode != tfjson.ManagedResourceMode {
			continue
		}
		f := c.feature(kind, rc.Type, rc.Address)
		c.setProvider(kind, rc.Type, rc.Address, rc.ProviderName)
		add(&f.Planned, action(rc.Change.Actions))
	}
	for _, invocation := range plan.ActionInvocations {
		if invocation == nil {
			continue
		}
		f := c.feature("action", invocation.Type, invocation.Address)
		c.setProvider("action", invocation.Type, invocation.Address, invocation.ProviderName)
		add(&f.Planned, "invoke")
	}
	// PlannedValues is a plan, not the state produced by a completed apply.
	// Computed coverage is taken from post-apply or refresh state instead.
}

// PreApplyPlan marks the most recently observed plan as the one which will
// actually be applied.
func (c *Collector) PreApplyPlan() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	c.lastPlan = nil
	c.actions = nil
	plan := c.latestPlan
	if plan == nil {
		add(&c.warnings, WarningPlanUnavailable)
		return
	}
	for _, rc := range plan.ResourceChanges {
		if rc != nil && rc.Change != nil && rc.Mode == tfjson.ManagedResourceMode {
			c.lastPlan = append(c.lastPlan, Feature{Kind: "resource", Name: rc.Type, Provider: rc.ProviderName, Address: safeAddress(rc.Address), Planned: []string{action(rc.Change.Actions)}})
		}
	}
	for _, invocation := range plan.ActionInvocations {
		if invocation != nil {
			c.actions = append(c.actions, Feature{Kind: "action", Name: invocation.Type, Provider: invocation.ProviderName, Address: safeAddress(invocation.Address)})
		}
	}
}

func (c *Collector) ImportPlan(plan *tfjson.Plan) {
	if c == nil || plan == nil {
		return
	}
	if !c.active {
		return
	}
	for _, rc := range plan.ResourceChanges {
		if rc != nil && rc.Change != nil && rc.Change.Importing != nil {
			f := c.feature("resource", rc.Type, rc.Address)
			if f.Provider == "" {
				f.Provider = rc.ProviderName
			}
			add(&f.Planned, "import")
		}
	}
}

func action(actions tfjson.Actions) string {
	if actions.Create() {
		return "create"
	}
	if actions.Delete() {
		return "delete"
	}
	if actions.Update() {
		return "update"
	}
	if actions.Replace() {
		return "replace"
	}
	if actions.NoOp() {
		return "no_op"
	}
	if actions.Read() {
		return "read"
	}
	return "unknown"
}

// ObserveState records only schema attribute names and null/known status from
// state explicitly obtained after a step operation, not setup or cleanup reads.
func (c *Collector) ObserveState(state *tfjson.State) {
	if c == nil || state == nil || state.Values == nil || state.Values.RootModule == nil {
		return
	}
	if c.active {
		c.module(state.Values.RootModule)
	}
}

func (c *Collector) module(m *tfjson.StateModule) {
	for _, r := range m.Resources {
		if r == nil {
			continue
		}
		kind := "resource"
		if r.Mode == tfjson.DataResourceMode {
			kind = "data_source"
		}
		f := c.feature(kind, r.Type, r.Address)
		c.setProvider(kind, r.Type, r.Address, r.ProviderName)
		states(f, r.AttributeValues)
		if r.IdentitySchemaVersion != nil && len(r.IdentityValues) > 0 {
			identity := c.feature("identity", r.Type, r.Address)
			identity.Provider = r.ProviderName
			states(identity, r.IdentityValues)
		}
	}
	for _, child := range m.ChildModules {
		c.module(child)
	}
}

func states(f *Feature, values map[string]any) {
	// Only top-level keys are unambiguously schema names. Nested map keys
	// could be sensitive user-provided values, so do not persist them.
	for name, value := range values {
		if value == nil {
			add(&f.Null, name)
		} else {
			add(&f.Known, name)
		}
	}
}

func (c *Collector) OperationForAddress(address, operation string) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	address = safeAddress(address)
	for _, f := range c.features {
		if f.Address == address && f.Kind == "resource" {
			add(&f.Completed, operation)
			return
		}
	}
	// Import steps without a config still identify the resource by address.
	parts := strings.Split(address, ".")
	if len(parts) > 0 {
		f := c.feature("resource", parts[0], address)
		if operation == "import" {
			add(&f.Completed, operation)
		}
	}
}

func (c *Collector) ApplyCompleted(destroy bool) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, p := range c.lastPlan {
		f := c.feature(p.Kind, p.Name, p.Address)
		if f.Provider == "" {
			f.Provider = p.Provider
		}
		if !destroy {
			add(&f.Completed, "apply")
		}
		for _, operation := range p.Planned {
			if operation == "unknown" {
				continue
			}
			if destroy {
				if operation != "delete" {
					continue
				}
				add(&f.Completed, "destroy")
			}
			add(&f.Completed, operation)
		}
	}
	for _, p := range c.actions {
		f := c.feature(p.Kind, p.Name, p.Address)
		if f.Provider == "" {
			f.Provider = p.Provider
		}
		add(&f.Completed, "invoke")
	}
}

// QueryCompleted marks list resources in a query step that passed its checks.
func (c *Collector) QueryCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.Kind == "list_resource" {
			add(&f.Completed, "query")
		}
	}
}

// QueryResults reads query message types without persisting result values,
// display names, instance keys, or provider-generated configuration.
func (c *Collector) QueryResults(messages []tfjson.LogMsg) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, message := range messages {
		switch m := message.(type) {
		case tfjson.ListResourceFoundMessage:
			for _, f := range c.features {
				if f.Kind == "list_resource" && f.Name == m.ListResourceFound.ResourceType {
					add(&f.Completed, "result")
				}
			}
		}
	}
}

func (c *Collector) RefreshCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.Kind == "resource" {
			add(&f.Completed, "refresh")
		}
	}
}

func (c *Collector) StateStoreCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.Kind == "state_store" {
			add(&f.Completed, "state_store")
		}
	}
}

func (c *Collector) ProviderStepCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.Kind == "provider" {
			add(&f.Completed, "step")
		}
	}
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
