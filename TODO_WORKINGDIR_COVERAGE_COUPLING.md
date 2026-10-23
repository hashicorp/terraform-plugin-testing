# TODO: decouple `plugintest.WorkingDir` from the semantic coverage collector

Smell: **Feature Envy / Message Chains** (judgement call from code review).

## What

`internal/plugintest/working_dir.go` gained an exported field:

```go
type WorkingDir struct {
	h        *Helper
	Coverage *semanticcoverage.Collector
	...
```

Step code then reaches through the working directory to reach the collector:

```go
wd.Coverage.PreApplyPlan()
wd.Coverage.ApplyCompleted(step.Destroy)
workingDir.Coverage.ImportPlan(plan)
workingDir.Coverage = testCaseWorkingDir.Coverage // import step
```

## Why this is a problem

- **Wrong owner.** `WorkingDir` manages a temp directory and runs
  Terraform. It has no use for the collector. It only carries it so that
  step code can find it, which is Feature Envy: the step code works on
  collector data but goes through `WorkingDir` to get there.
- **Package coupling.** `plugintest` now imports `semanticcoverage`. Every
  user of `plugintest`, including `helper/resource` tests and the
  `WorkingDir` unit tests, now depends on the coverage package, even
  though coverage is an opt-in experiment.
- **Manual propagation is easy to forget.** The import step creates a
  second `WorkingDir` and has to copy the field by hand. A future code
  path that creates a working directory and forgets the copy silently
  loses evidence. Nothing fails, and the report just undercounts.
- **Exported mutable field.** Any caller can swap or nil the collector
  partway through a test. Collector state is per test case, so a swap
  partway through a step corrupts that step's finding.
- **Message chains hide the dependency.** `wd.Coverage.X()` makes
  `WorkingDir` look like the real dependency. The real dependency is the
  collector.

## Fix

Pass the collector, or a `stepCoverage` wrapper, explicitly to the step
functions. Then remove the field from `WorkingDir`. See
`PLAN_FIX_SHOTGUN_SURGERY.md`, steps 1-3. Both smells are fixed by the
same change.

## Done when

- `internal/plugintest` has no import of `semanticcoverage`.
- `grep -rn '\.Coverage' helper internal/plugintest` returns nothing.
- Findings from an acceptance run are unchanged.
