# TODO: one `FinishStep(outcome)`, and don't record skipped steps

Source: review discussion of `Finish` and `FinishStep`.

## Decisions

- Rename `Begin(step)` to `BeginStep(stepNumber)`. It is called once per step, at the top of the step loop in `runNewTest`.
- Keep a single `FinishStep(outcome Outcome) error`, which is today's `Finish` renamed. Delete the `FinishStep(failed, skipped bool)` variant.
- **Skipped steps write no finding.** The report already ignores `skipped` findings exactly like `failed` ones (`cmd/semantic-coverage-report/main.go`, the `OutcomeSkipped` `continue`), and it doesn't count them anywhere. So nothing is lost.
- Failed steps are still written as `failed`.
- The `expected_error` override (a matched `ExpectError` turns `success` into `expected_error`) stays inside the collector.

## Changes

1. **Remove `OutcomeSkipped`** from `Outcome`, `Valid()`, the report's `continue`, and `SEMANTIC_COVERAGE.md` (the outcome list and "Failed and skipped steps contribute no coverage").
2. **Add `(c *Collector) DiscardStep()`.** It sets `active = false` without writing anything, and it replaces the `skipped` field.
   - SkipFunc path (`testing_new.go`): `coverage.DiscardStep()` instead of `coverage.Warning("step_skipped")`. This also removes the side effect from `Warning` (see `TODO_WARNING_CODES.md`).
   - Deferred call in `runNewTest`:
     ```go
     defer func() {
         if t.Skipped() {
             coverage.DiscardStep()
             return
         }
         outcome := semanticcoverage.OutcomeSuccess
         if t.Failed() {
             outcome = semanticcoverage.OutcomeFailed
         }
         if err := coverage.FinishStep(outcome); err != nil { ... }
     }()
     ```
     If `PLAN_FIX_SHOTGUN_SURGERY.md` lands first, put this into `semantic_coverage.go` as a `finishFromTest(t)` helper.
3. **Rename** `Begin` → `BeginStep` and `Finish` → `FinishStep` at all call sites and in the tests. Update the `Begin` doc comment ("callers must FinishStep or DiscardStep the previous step first").
4. Remove the `Begin`/`Finish`/`FinishStep` row from `TODO_COLLECTOR_NAMING.md`. This file replaces it.

## Check

- `go test ./internal/semanticcoverage/... ./helper/resource/... ./cmd/...`
- Run an acceptance test with a `SkipFunc` step. Check that no finding is written for that step, and that the steps around it are still written.
- Bump `FindingVersion` (removing an outcome value is an incompatible change). Do it together with `TODO_WARNING_CODES.md` so that there is only one bump.
