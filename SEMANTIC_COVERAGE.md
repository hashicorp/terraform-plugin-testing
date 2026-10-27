# Semantic coverage (experimental)

> **Note:** This is a temporary home for this documentation. It will be
> migrated to the official Terraform plugin testing documentation later.

Semantic coverage reports which parts of a provider's schema the acceptance
tests actually exercise: which resources, data sources, ephemeral resources,
list resources, actions, state stores, and functions are used, which of their
attributes are configured or populated, and which lifecycle operations were
planned or completed.

It is derived from the tests themselves. While enabled, the test harness
observes each `TestStep`: it reads the configuration copied into the step's
working directory and inspects Terraform plans and post-operation state. The
`cmd/semantic-coverage-report` command compares those findings with a
`terraform providers schema -json` document.

## Usage

Collection is opt-in. Set `TF_ACC_SEMANTIC_COVERAGE_DIR` to a directory when
running acceptance tests; one JSON finding is written per executed test step.
Use an empty directory per run, because the report reads every JSON file under
it and combines all findings: a schema element is covered if any finding covers
it. Findings from other runs would therefore count as coverage of this one.
Running tests repeatedly within one run (`go test -count=N` with N > 1) is not
supported. Re-run the tests after upgrading terraform-plugin-testing, because
the report rejects findings written by an incompatible version.

```sh
TF_ACC=1 \
TF_ACC_SEMANTIC_COVERAGE_DIR="$PWD/findings/run-1" \
  go test -count=1 ./internal/provider

# Generate the schema with the same provider version used in the test run.
terraform -chdir=/path/to/initialized/provider/config providers schema -json \
  > provider-schema.json

go run github.com/hashicorp/terraform-plugin-testing/cmd/semantic-coverage-report \
  -findings-dir findings/run-1 \
  -schema-file provider-schema.json \
  -json-out coverage.json -text-out coverage.txt
```

When `TF_ACC_SEMANTIC_COVERAGE_DIR` is unset, the collector is disabled and no
additional Terraform commands are run.

## Evidence model

- A configurable attribute (required, optional, or a block) is covered when a
  completed or expected-error step configures it in HCL configuration. The
  collector reads generated `Config` strings, `ConfigFile`, and
  `ConfigDirectory` at runtime.
- A computed attribute is covered when a successful step's resulting state has
  a non-null value for it that the step didn't configure. Only top-level state
  attribute names are read: deeper maps can contain user-provided keys.
- Optional-computed attributes count in both the configurable and the computed
  total, because both behaviours need a test: one step that sets the attribute
  and one that leaves it to the provider. Configuration and state can't be
  paired per instance (`count`, `for_each`, modules), so state evidence for
  such an attribute only counts in a step where no resource of the same type
  configures it.
- Resource actions in saved plans are **planned** evidence. **Completed**
  create, update, replace, delete, or no-op requires a successful apply of the
  matching pre-apply plan. `PlanOnly` never contributes completed apply
  evidence. Explicit `Destroy` steps can contribute completed destroy;
  automatic cleanup does not. Command imports contribute completed import;
  import-block tests that only plan contribute planned import.
- Action invocations from plans and successful apply, list-resource query
  results, provider-defined function references, and state-store step
  completion are recorded where the test harness exposes them.
- Each finding includes its test name, step number, and outcome
  (`success`, `expected_error`, `failed`). A matched `ExpectError`
  retains input and plan evidence but cannot claim a completed operation or
  computed state value. Failed steps contribute no coverage. Skipped steps
  write no finding.
- A finding can carry warnings. Each warning means that one type of evidence
  for the step may be incomplete:

  | Warning | Meaning |
  |---|---|
  | `config_unavailable` | Some configuration could not be read or parsed. Input and feature evidence may be missing. |
  | `plan_unavailable` | A plan could not be read. Planned and completed operation evidence may be missing. |
  | `state_unavailable` | The post-operation state could not be read. Known and null attribute evidence may be missing. |

## Privacy

The collector does not persist configuration text, attribute values, plan or
state JSON, query result values, or assertion matchers. Instance keys in
resource addresses are replaced with `[*]`. It does persist provider and
feature names, normalized resource addresses, schema attribute names,
known/null classifications, operations, test names, and test package import
paths.
Write findings to a trusted private directory.

## Report

The report contains schema-backed entity and attribute totals, plus planned and
completed operation lists with step references. Operation totals are
observations, not percentages, because the schema does not declare which
lifecycle cases apply to a feature. The schema document must describe the same
provider selection as the test run.

## Current limits

- Ephemeral-resource and function results are not stored in Terraform state,
  so their use and inputs are identified but their computed outputs are not
  covered.
- Attribute totals include only top-level schema attributes and blocks. Nested
  attributes are not counted until a schema-aware walker can distinguish
  nested schema names from user map keys.
- Resource identity attributes are not counted. They are only observed in
  state and have no configurable/computed distinction (only required/optional
  for import).
- Modules referenced from external sources may be visible in plan actions or
  state, but not in the locally copied configuration files.
- Terraform JSON configuration (a JSON `Config` string or a `.tf.json` file)
  is not read. Steps that use it get a `config_unavailable` warning and have no
  configuration evidence; plan and state evidence is still recorded.
- Provider blocks and provider-defined function calls name a provider by its
  local name, which the report matches against the type segment of the
  provider's source address. If `required_providers` gives a provider a
  different local name (for example `mycloud = { source = "hashicorp/aws" }`),
  that evidence is not attributed to the provider.
- Each plan observed while collecting costs one additional `terraform show`,
  and each completed apply one additional state read.

## Example result

Against `terraform-provider-local` (excluding the download-dependent upgrade
tests), 58 step findings yielded 9/9 active feature entities, 30/30
configurable attributes, and 39/42 computed attributes. The four
optional-computed permission attributes count in both totals and were covered
both set and unset. The three uncovered computed attributes were
ephemeral-resource outputs.
