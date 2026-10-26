# Plan: split `internal/semanticcoverage/collector.go`

## Problem (Divergent Change)

`collector.go` is about 830 lines and changes for unrelated reasons:

1. The finding wire format (`Feature`, `Finding`, `Outcome`,
   `FindingVersion`). This is shared with `cmd/semantic-coverage-report`.
2. Step lifecycle and output (`New`, `Ready`, `Begin`, `Finish`,
   `FinishStep`, `Warning`, `MatchedExpectedError`, file writing).
3. Config parsing for HCL and Terraform JSON (`Config`, `QueryConfig`,
   `Parse`, `blocks`, `walkInputs`, `functions`, `jsonBlocks`,
   `jsonInputs`).
4. Plan observation (`Plan`, `PreApplyPlan`, `ImportPlan`, `action`).
5. State observation (`ObserveState`, `module`, `states`).
6. Completion marking (`ApplyCompleted`, `QueryCompleted`,
   `QueryResults`, `RefreshCompleted`, `StateStoreCompleted`,
   `ProviderStepCompleted`, `OperationForAddress`).
7. Address redaction (`safeAddress`). This is a privacy rule.

A change to HCL parsing and a change to the JSON format both touch the
same file, which makes review and blame noisy. The privacy-critical code
(`safeAddress`, and the "never persist values" rules in `states`) is also
buried among unrelated code.

## Proposal

Split the file within the same package. This is a pure move: no API
change and no behaviour change.

| File | Contents |
|---|---|
| `finding.go` | Package doc, `EnvDirectory`, `FindingVersion`, `Outcome` and its constants and methods, `Feature`, `Finding` |
| `collector.go` | `Collector` struct and its nil-receiver contract doc, `New`, `Enabled`, `Ready`, `Begin`, `Finish`, `FinishStep`, `Warning`, `MatchedExpectedError`, `feature`, `setProvider`, `add` |
| `config.go` | `Config`, `QueryConfig`, `Parse`, `blocks`, `walkInputs`, `functions`, `jsonBlocks`, `jsonInputs`, `asMap` |
| `plan.go` | `Plan`, `PreApplyPlan`, `ImportPlan`, `action` |
| `state.go` | `ObserveState`, `module`, `states` |
| `completion.go` | `ApplyCompleted`, `ProviderStepCompleted`, `RefreshCompleted`, `QueryCompleted`, `QueryResults`, `StateStoreCompleted`, `OperationForAddress` |
| `redact.go` | `safeAddress`, plus a package comment that states the "no values persisted" rules in one place |

Split `collector_test.go` the same way: `config_test.go`, `plan_test.go`,
`redact_test.go`, and so on. Each test moves next to the code it covers.

## Follow-ups (separate commits, optional)

- **Kind constants.** The reviewer decided to make `Kind` an enum. See
  `TODO_COLLECTOR_SIMPLIFY.md` item 1. Kind strings repeat across `blocks`,
  `jsonBlocks`, `Plan` and `module`: `"resource"`, `"data_source"`,
  `"list_resource"` and others. One `blockKinds` map from the HCL/JSON
  block type to `(kind, address prefix)` would replace the nested `switch`
  in `blocks` and the `if` chain in `jsonBlocks` (Repeated Switches).
  The decided approach is a single Kind table that gives each kind its
  address prefix, see `TODO_COLLECTOR_SIMPLIFY.md` item 1.
- Other small cleanups (a real set instead of `add`, inlining
  `jsonInputs`, the query file name constant) are in
  `TODO_COLLECTOR_SIMPLIFY.md`.
- **Guard helper.** Most exported methods start with the same lines
  (`nil` check, `active` check). A
  `func (c *Collector) withActive(fn func())` helper would remove some
  repetition. It would also make the nil-receiver contract hold by
  construction.
- Apply the renames from `TODO_COLLECTOR_NAMING.md` after the move, not
  during it. That keeps `git log --follow` and review simple.

## Verification

- Commit the move on its own. `git diff -M --stat` should show mostly
  moved lines.
- `go build ./... && go test ./internal/semanticcoverage/... ./cmd/... ./helper/resource/...`
- `golangci-lint run ./internal/semanticcoverage/...`
