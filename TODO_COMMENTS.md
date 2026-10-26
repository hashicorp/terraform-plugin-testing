# TODO: fix and add code comments in the semantic coverage code

Source: review comments.

## What

| Location | Change |
|---|---|
| `helper/resource/testing_new_config.go`, before the post-apply `ObserveState` block | Replace "Inspect the actual post-apply state even when the test has no Check." with the reason. Known and null attribute evidence must come from the state that an apply actually produced. Plan values are not used for this (see the comment at the end of `Collector.Plan`). Without this read, a step that has no `Check` would record no state evidence at all. |
| same place | Delete "This does not persist any of the state's values." The privacy rules belong in one place: the `redact.go` package comment from `PLAN_FIX_COLLECTOR_SPLIT.md`. |
| `Collector.Plan` doc comment | Delete "not inferred provider RPCs". Say what the method does, not what it doesn't do. For example: "Plan records the planned operation of each resource change and action invocation in plan. A nil plan means the saved plan could not be read." |
| `Collector.MatchedExpectedError` | Add a comment explaining why it locks. Every exported method takes `c.mu`, because a provider server or a `Check` function can call into the collector from another goroutine than the step runner. Verify that this is actually true before you write it down. If all calls come from the step goroutine, the mutex may not be needed at all. In that case, document that instead, or remove the mutex everywhere. |
| `Collector.lastPlan` | Add a field comment. It holds the resource changes of the plan that `PreApplyPlan` selected, and `ApplyCompleted` marks them completed. Also see the rename in `TODO_COLLECTOR_NAMING.md`. |
| `helper/resource/testing_new_query.go`, `ResetConfig` call | Explain the reason, or delete the call (see below). |
| `FindingVersion` doc comment (`collector.go:26`) | Keep the constant, and explain why it exists. Findings live in a user-owned directory (`TF_ACC_SEMANTIC_COVERAGE_DIR`) that outlives the test run. The writer is the test binary and the reader is `cmd/semantic-coverage-report`, and the two can be built from different module versions. Bump the constant whenever the meaning or shape of `Finding` changes incompatibly. Add one sentence to `SEMANTIC_COVERAGE.md` as well: re-run the tests after upgrading, because the report rejects old findings. |
| `cmd/semantic-coverage-report`, version check | Make the error message actionable, e.g. `"%s: finding version %d, this report reads version %d; re-run the tests with the same terraform-plugin-testing version"`. |

## The query step's `ResetConfig` may be dead code

The current comment says the call drops input evidence from "any earlier config in this step". I could not find such a config:

- `Begin` already replaces `c.features` with an empty map for every step.
- `runNewTest` skips `wd.SetConfig` when `step.Query` is true (`testing_new.go:284`). Before `testStepNewQuery`, nothing in a query step calls `Config`.

The only user of the behaviour seems to be `TestQueryConfigDoesNotIncludePreviousStep` in `collector_test.go`, and that test calls the collector directly.

Steps:

1. Remove the `ResetConfig` call. Run `go test ./helper/resource/...` and the query acceptance tests with `TF_ACC_SEMANTIC_COVERAGE_DIR` set. Compare the findings with and without the call.
2. If they match, delete `Collector.ResetConfig` and the test that covers it.
3. If they differ, the diff shows which earlier hook records config. Write that down as the reason in the comment. Better still, stop that hook from firing for query steps.
