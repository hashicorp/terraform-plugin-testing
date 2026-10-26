# TODO: rename unclear semantic coverage collector APIs

Smell: **Mysterious Name** (judgement call from code review).

## What

Several names in `internal/semanticcoverage` don't say what
they do or what they cover:

| Current | Problem | Suggested |
|---|---|---|
| `Ready()` | It creates the output directory. "Ready" sounds like a query or a readiness check. The reviewer confirmed it needs a better name. | `EnsureDirectory()`, or fold it into `New`. `Finish` already calls `MkdirAll` itself, so check whether `Ready` is only there to fail early. If so, `New` returning an error is the cleanest option. |
| field `test` | Ambiguous. | `testName` |
| field `step` | Doesn't say whether it is an index. It is **1-based** (`Begin(stepNumber)` in `testing_new.go`), so `stepIndex` would be misleading. | `stepNumber`, matching `runNewTest` |
| fields `lastPlan` / `actions` | `lastPlan` holds the resource changes that `PreApplyPlan` picked for `ApplyCompleted` to mark completed; "last" doesn't say that. `actions` is the matching list of action invocations. | `lastPlan` → `pendingApplyChanges`; `actions` → `pendingApplyActions`. Keep the field comments. |
| `ApplyCompleted(destroy)` + `ProviderStepCompleted()` | They are always called together (`testing_new_config.go`). `ProviderStepCompleted` marks providers with the "step" operation, and its name doesn't say that. | Merge into `ApplySucceeded(destroy)`, or rename the second to `MarkProvidersUsed()` |
| `OperationForAddress(address, op)` | It reads like a getter, but it records a completed operation. | `RecordCompleted(address, op)` |
| `Plan(plan)` | It sounds like it creates a plan. It actually records the plan. | `ObservePlan(plan)`, matching `ObserveState` |
| `Parse(name, data)` / `Config(dir)` / `QueryConfig(dir)` | The noun names hide that these read files and record inputs. | `ObserveConfigFile`, `ObserveConfigDir`, `ObserveQueryConfig` |
| `states(f, values)` | It records known and null attribute names. | `recordAttributeStatus` |
| `module(m)` | It walks a state module recursively. | `observeStateModule` |
| `blocks(blocks)` | It walks HCL blocks and records features and inputs. | `observeHCLBlocks` |
| `add(xs, v)` | It is a set-insert into a slice. | Replace it with a set type, see `TODO_COLLECTOR_SIMPLIFY.md` item 2. Otherwise rename it to `appendUnique`. |
| `action(actions)` | It maps tfjson actions to an operation string. | `operationForActions` |

## Why this is a problem

- **Readers have to open the body to learn the contract.** The call sites
  in `helper/resource` are what most maintainers read. Lines like
  `coverage.ProviderStepCompleted()` or `coverage.Ready()` don't say
  what evidence is being recorded. That is the semantics the report
  depends on.
- **Inconsistent verbs blur the model.** Some methods record observations
  (`ObserveState`), some mark completion (`XCompleted`), and some do both
  under a noun name (`Plan`, `Config`). The evidence model in
  `SEMANTIC_COVERAGE.md` separates planned from completed evidence and
  inputs from state. The API should use the same split, so that a wrong
  call is visible at the call site.
- **Paired calls hint at a missing concept.** If two methods must always
  be called together, callers can forget one. That silently drops
  evidence.
- **Getter-looking mutators cause bugs.** Someone reading
  `OperationForAddress` might treat it as side-effect free and call it in
  a check path, which would record false coverage.

## Notes

- The package is `internal`, so renames are not breaking changes.
- If no honest name comes for a method, that is a sign it does two things.
  Split it.
