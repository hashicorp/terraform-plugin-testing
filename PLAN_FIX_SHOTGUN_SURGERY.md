# Plan: fix Shotgun Surgery in semantic coverage hooks

## Problem

Coverage events are recorded in two layers and seven files:

| Layer | File | Hook |
|---|---|---|
| `plugintest.WorkingDir` | `internal/plugintest/working_dir.go` | `Config` (in `SetConfig`), `QueryConfig` (in `SetQuery`), `Plan` (in `CreatePlan`) |
| step runner | `helper/resource/testing_new.go` | `BeginStep`, `FinishStep` x2, `DiscardStep` x2, `MatchedExpectedError` x5, `RefreshCompleted`, `QueryCompleted`, `StateStoreCompleted` |
| step runner | `testing_new_config.go` | `PreApplyPlan` x2, `ApplyCompleted`, `ProviderStepCompleted`, post-apply `ObserveState` |
| step runner | `testing_new_import_state.go` | `ImportPlan`, `ObserveState`, `OperationForAddress` |
| step runner | `testing_new_refresh_state.go` | `ObserveState` |
| step runner | `testing_new_query.go` | `ResetConfig`, `QueryResults` |

Adding a new evidence type, or changing what counts as "completed",
means editing several of these files. It is also hard to see the whole
collection lifecycle in one place. `PreApplyPlan` only works because
`WorkingDir.CreatePlan` called `Plan` earlier, and that implicit ordering
crosses the package boundary.

## Proposal: record everything in the step layer (`helper/resource`)

`WorkingDir` is a thin wrapper around Terraform commands. It should not
know about coverage. Every hook moves to `helper/resource`, and all
coverage wiring sits behind one file-local adapter.

1. **Remove `WorkingDir.Coverage`.** Delete the field and its three hooks
   from `working_dir.go`. `plugintest` no longer imports
   `semanticcoverage`. This also fixes the Feature Envy item in
   `TODO_WORKINGDIR_COVERAGE_COUPLING.md`.
2. **Add `helper/resource/semantic_coverage.go`** with a small type that
   owns the collector and wraps the `WorkingDir` operations that produce
   evidence:

   ```go
   type stepCoverage struct{ c *semanticcoverage.Collector }

   func (s stepCoverage) setConfig(ctx, wd, cfg, vars) error // wd.SetConfig + c.Config(wd.WorkingDirectory())
   func (s stepCoverage) setQuery(ctx, wd, cfg, vars) error  // c.ResetConfig + wd.SetQuery + c.QueryConfig
   func (s stepCoverage) createPlan(ctx, wd, opts...) error  // wd.CreatePlan + SavedPlan + c.Plan
   func (s stepCoverage) observeState(state *tfjson.State)
   func (s stepCoverage) stepResult(mode stepMode, err error, step TestStep)
   ```

   `stepResult` replaces the four
   `if err == nil && step.ExpectError == nil { coverage.XCompleted() }`
   blocks and the five `MatchedExpectedError()` calls with one switch on
   the step mode.
3. **Pass `stepCoverage` explicitly** to `testStepNewConfig`,
   `testStepNewImportState`, `testStepNewRefreshState` and
   `testStepNewQuery`, as an extra parameter. Do not hang it off
   `WorkingDir`. The import step's separate working directory then needs
   no `workingDir.Coverage = ...` copy.
4. **Replace the call sites.** Change `wd.SetConfig`, `wd.SetQuery` and
   `wd.CreatePlan` to the `stepCoverage` wrappers only where evidence is
   wanted. Do not use them in setup or cleanup paths, which matches the
   current "not setup or cleanup reads" rule.
5. **Fold `PreApplyPlan`** into a `createPlan(..., preApply bool)`
   variant. The "latest plan" is then never implicit state shared between
   two calls.

## Result

- Coverage-aware code lives in `semantic_coverage.go`. The `testing_new_*.go`
  files gain one-line calls instead of inline logic.
- `plugintest` is unchanged from `main`.
- A new evidence type touches `collector.go` and `semantic_coverage.go`
  only.

## Risks / checks

- Before step 4, check whether the query step's `ResetConfig` is
  dead code (`TODO_COMMENTS.md`). If it is, don't carry it into `setQuery`.
- Combine with `TODO_STATE_READ_GYMNASTICS.md`: call `observeState` after
  the `runProviderCommand` closure, not inside it.
- `SetConfig` is called from several places, including
  `testing_new_config.go` and the import step. Audit every call so the
  set of config reads stays the same. Today all of them feed coverage.
- `CreatePlan` is called inside `runProviderCommand` closures, and the
  wrapper must run inside the same closure.
- Verify with `go test ./helper/resource/... ./internal/semanticcoverage/...`
  and an acceptance run with `TF_ACC_SEMANTIC_COVERAGE_DIR` set. Diff the
  findings before and after; they should be identical.
