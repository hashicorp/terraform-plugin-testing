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
it.

```sh
TF_ACC=1 \
TF_ACC_SEMANTIC_COVERAGE_DIR="$PWD/findings/run-1" \
  go test -count=1 ./internal/provider

# Generate the schema with the same provider version used in the test run.
terraform -chdir=/path/to/initialized/provider/config providers schema -json \
  > provider-schema.json

go run github.com/hashicorp/terraform-plugin-testing/cmd/semantic-coverage-report \
  -findings findings/run-1 \
  -schema provider-schema.json \
  -json coverage.json -text coverage.txt
```

When `TF_ACC_SEMANTIC_COVERAGE_DIR` is unset, the collector is disabled and no
additional Terraform commands are run.

## Evidence model

- An input is covered when a completed or expected-error step configures its
  schema attribute in HCL or Terraform JSON configuration. The collector reads
  generated `Config` strings, `ConfigFile`, and `ConfigDirectory` at runtime.
- A computed attribute is covered when a completed step's resulting state has
  a non-null value for that attribute. Only top-level state attribute names are
  read: deeper maps can contain user-provided keys. Optional-computed attributes
  can be covered by either input or state evidence.
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

## Privacy

The collector does not persist configuration text, attribute values, plan or
state JSON, query result values, or assertion matchers. Instance keys in
resource addresses are replaced with `[*]`. It does persist provider and
feature names, normalized resource addresses, schema attribute names,
known/null classifications, operations, test names, and test package names.
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
- Modules referenced from external sources may be visible in plan actions or
  state, but not in the locally copied configuration files.
- Each plan observed while collecting costs one additional `terraform show`,
  and each completed apply one additional state read.

## Example result

Against `terraform-provider-local` (excluding the download-dependent upgrade
tests), 80 step findings yielded 9/9 active feature entities, 29/29 input
attributes, and 35/38 computed attributes. The three uncovered computed
attributes were ephemeral-resource outputs.
