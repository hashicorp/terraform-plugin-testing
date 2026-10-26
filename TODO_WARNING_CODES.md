# TODO: warnings become a typed enum with three codes

Source: review discussion of the warning codes.

## Decisions

- A warning means "this step's evidence of type X may be incomplete". There is one code per evidence type, and the codes don't say why. To debug the collector itself, use the logs or the tests, not the codes.
- Add `type Warning string` with constants and `Valid()` in `finding.go`, next to `Outcome`. `Finding.Warnings` becomes `[]Warning`. The report rejects unknown codes, the same way it does for outcomes.
- `step_skipped` is no longer a warning. It becomes `DiscardStep()`, see `TODO_FINISH_STEP.md`. After that, `Warning()` has no side effects.

## Mapping

| new code | replaces |
|---|---|
| `WarningConfigUnavailable` = `config_unavailable` | `config_directory_unavailable`, `config_file_unavailable`, `query_config_unavailable`, `config_parse_failed`, `config_body_unsupported`, `function_parse_failed` |
| `WarningPlanUnavailable` = `plan_unavailable` | `plan_unavailable`, `pre_apply_plan_unavailable` |
| `WarningStateUnavailable` = `state_unavailable` | `state_unavailable` |

## Changes

1. Add the type and the constants. Change `Warning(code string)` to `Warning(code Warning)`, and make the internal `add(&c.warnings, "...")` calls use the constants as well.
2. Replace the codes at every emit site. Find them with `rg 'Warning\("|warnings, "'` in `internal/semanticcoverage` and `helper/resource`.
3. In the report, change `output.Warnings` to `[]semanticcoverage.Warning`, and validate the codes in the findings reader.
4. Add a short table of the three codes and their meaning to `SEMANTIC_COVERAGE.md`.
5. Update the collector tests that assert on the old codes.

## Order

Do this after `TODO_FINISH_STEP.md`, because that one removes `step_skipped` first. If the file split (`PLAN_FIX_COLLECTOR_SPLIT.md`) is already done, the type goes into `finding.go`.

## Format version

Changing the codes is an incompatible format change. Bump `FindingVersion` once, together with `TODO_FINISH_STEP.md`.
