# Semantic coverage (experimental)

> **Note:** This is a temporary home for this documentation. It will be
> migrated to the official Terraform plugin testing documentation later.

Semantic coverage reports which parts of a provider's schema the acceptance
tests actually exercise: which resources, data sources, ephemeral resources,
list resources, actions, state stores, and functions are used, which of their
attributes are configured or populated, and which operations were planned or
completed.

It is derived from the tests themselves. While enabled, the test harness
observes each `TestStep`: it reads the configuration copied into the step's
working directory and inspects Terraform plans and post-operation state. The
`cmd/semantic-coverage-report` command compares those findings with a
`terraform providers schema -json` document.

## Usage

Collection is opt-in. Set `TF_ACC_SEMANTIC_COVERAGE_DIR` to a directory when
running acceptance tests; one JSON finding is written per executed test step.
Use an empty directory per run, because the report reads every JSON file under
it and combines all findings: a feature or attribute is covered if any finding
covers it. Findings from other runs would therefore count as coverage of this one.
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

Each finding describes one test step: its test name, step number, and
outcome, and the provider **features** (resources, data sources, ephemeral
resources, list resources, actions, state stores, providers, and functions)
that the step used. All uses of one feature type in a step are merged.

- **Outcome.** `success` means the step passed. `expected_error` means it
  failed with an error matching `ExpectError`. `failed` means the test failed
  during the step. Skipped steps write no finding.
- **Configured** attributes are the top-level attributes and blocks that the
  step's HCL configuration sets. The collector reads generated `Config`
  strings, `ConfigFile`, and `ConfigDirectory` at runtime. Configuration that
  the test harness generates itself is not recorded.
- **Populated** attributes are the top-level attributes with a non-null value
  in the state that the step's apply, import, or refresh produced. Only
  top-level state attribute names are read: deeper maps can contain
  user-provided keys.
- **Operations** are `create`, `update`, `replace`, `delete`, `read`,
  `import`, `refresh`, `query`, and `invoke`. An operation is **planned** when
  it appears in one of the step's plans, and **completed** when Terraform
  performed it successfully: an applied plan, an import command, a refresh,
  or a query. A no-op change has no operation. `PlanOnly` steps and
  import-block tests that only plan contribute planned operations only.
  Explicit `Destroy` steps complete `delete`; the automatic cleanup after the
  test is not recorded.
- A finding can carry warnings. Each warning means that one source of
  evidence for the step was unavailable:

  | Warning | Meaning |
  |---|---|
  | `config_unavailable` | Some configuration could not be read or parsed. Features and configured attributes may be missing. |
  | `plan_unavailable` | A plan could not be read. Planned and completed operations may be missing. |
  | `state_unavailable` | The state after an operation could not be read. Populated attributes may be missing. |

## Coverage rules

The report applies these rules to the findings:

- Failed steps contribute no coverage.
- A feature is **covered** when it appears in a finding that didn't fail.
- A configurable attribute (required, optional, or a block) is covered when a
  step that didn't fail configures it.
- A computed attribute is covered when a successful step finds it populated
  without configuring it. A configured value is echoed into state, and
  configuration and state can't be paired per instance (`count`, `for_each`,
  modules), so a step that configures the attribute on any feature of the
  same type doesn't count.
- Optional-computed attributes count in both the configurable and the computed
  total, because both behaviours need a test: one step that sets the attribute
  and one that leaves it to the provider.
- A matched `ExpectError` keeps configured attributes and planned operations,
  but can't claim completed operations or populated attributes.

## Privacy

The collector does not persist configuration text, attribute values, plan or
state JSON, query results, resource names, addresses, instance keys, or
assertion matchers. It does persist provider and feature type names, schema
attribute names, operations, test names, and test package import paths.
Write findings to a trusted private directory.

## Report

The report contains schema-backed feature and attribute totals, plus planned
and completed operation lists with step references. Operation totals are
observations, not percentages, because the schema does not declare which
lifecycle cases apply to a feature. The schema document must describe the same
provider selection as the test run.

## Current limits

- Ephemeral-resource and function results are not stored in Terraform state,
  so their use and configured attributes are identified, but their computed
  outputs are not covered.
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
- The query file of a `Query` step stays in the working directory, so later
  steps in the same test also record its list resources and their configured
  attributes.
- Each apply costs one additional `terraform show` for the plan and one for
  the resulting state while collecting.

## Example result

Against `terraform-provider-local` at commit `b3d30fe` (excluding the
download-dependent upgrade tests), 58 findings yielded 9/9 covered features,
30/30 configurable attributes, and 39/42 computed attributes. The four
optional-computed permission attributes count in both totals and were covered
both set and unset. The three uncovered computed attributes were
ephemeral-resource outputs. Completed operations were `create` and `replace`
for both file resources, `read` for the `local_command` data source, and
`invoke` for the `local_command` action.
