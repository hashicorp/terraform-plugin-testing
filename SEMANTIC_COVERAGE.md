# Semantic coverage (experimental)

> **Note:** This is a temporary home for this documentation. It will be
> migrated to the official Terraform plugin testing documentation later.

Semantic coverage reports which parts of a provider's schema the acceptance
tests exercise. For each resource, data source, ephemeral resource, list
resource, action, state store, and function, it reports whether a test used
it, which of its top-level attributes a test configured or found populated,
and which operations Terraform planned or completed.

Collection happens inside the test harness. While it is enabled, the harness
records evidence for each `TestStep`: it reads the configuration copied into
the step's working directory, and it inspects the step's plans and the state
after each operation. The `cmd/semantic-coverage-report` command then compares
the recorded findings with a `terraform providers schema -json` document.

## Usage

Collection is off by default. To turn it on, set
`TF_ACC_SEMANTIC_COVERAGE_DIR` to a directory when running acceptance tests.
The harness creates the directory if needed and writes one JSON finding per
executed test step.

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

Use a new, empty directory for each run. The report reads every `.json` file
under `-findings-dir`, including subdirectories, and counts an attribute as
covered if any of those findings covers it. Findings left over from an earlier
run would count toward the current one.

Running each test more than once in a run (`go test -count=N` with N > 1) is
not supported. After upgrading terraform-plugin-testing, run the tests again
and use the report from the same version: the report rejects findings whose
format version differs from its own.

When `TF_ACC_SEMANTIC_COVERAGE_DIR` is unset, nothing is recorded and the
harness runs no additional Terraform commands.

### Report flags

| Flag | Default | Description |
|---|---|---|
| `-findings-dir` | `findings` | Directory of findings to read |
| `-schema-file` | none, required | Output of `terraform providers schema -json` |
| `-json-out` | `coverage.json` | Path of the JSON report |
| `-text-out` | `coverage.txt` | Path of the text report |

The command exits with status 1 on any error, for example when `-schema-file`
is missing, the schema has no providers, a finding cannot be read or fails
validation, or the directory contains no findings. Uncovered features and
attributes do not change the exit status, and neither do report warnings.

## What a finding records

Each finding describes one test step. It holds the test's package and name,
the step number, the step's outcome, and the provider **features** that the
step used: resources, data sources, ephemeral resources, list resources,
actions, state stores, providers, and functions. All uses of one feature type
within a step are merged into one entry.

The outcome is one of:

| Outcome | Meaning |
|---|---|
| `success` | The step passed. |
| `expected_error` | The step failed with an error that matched `ExpectError`. |
| `failed` | The test failed during the step, including through `t.Fatal` or a panic. |

Skipped steps write no finding.

For each feature, the finding lists the following evidence.

**Configured attributes** are the top-level attributes and blocks set in the
step's HCL configuration. The collector reads the configuration files after
the harness writes them to the working directory, so `Config`, `ConfigFile`,
and `ConfigDirectory` are all read. Configuration that the harness generates
itself is not recorded.

**Populated attributes** are the top-level attributes with a non-null value in
the state produced by the step's apply, import, or refresh. The collector
reads only top-level attribute names, because keys further down can be
user-provided map keys. Resources in child modules record no populated
attributes: module configuration is not read, so a value set by the module
would otherwise look like one computed by the provider.

**Operations** are `create`, `update`, `replace`, `delete`, `read`, `import`,
`refresh`, `query`, and `invoke`. An operation is *planned* when it appears in
one of the step's plans. It is *completed* when Terraform performed it
successfully: by applying a plan, running an import command, refreshing, or
running a query. A no-op change records the feature but no operation. Steps
with `PlanOnly`, and import-block tests that only plan, record planned
operations only. A step with `Destroy: true` completes `delete`; the automatic
destroy after the last step is not recorded.

**Served providers** are recorded once per step rather than per feature. They
are the source addresses under which the harness served the provider
factories from the test binary. Providers that Terraform installs, such as
`ExternalProviders` or the released provider in an upgrade test, are not in
this list.

A finding can also carry warnings. Each one means that a source of evidence
for the step was unavailable, and the report repeats it:

| Warning | Meaning |
|---|---|
| `config_unavailable` | Some configuration could not be read or parsed. Features and configured attributes may be missing. |
| `plan_unavailable` | A plan could not be read. Planned and completed operations may be missing. |
| `state_unavailable` | The state after an operation could not be read. Populated attributes may be missing. |

## How the report counts coverage

The report applies these rules:

- A `failed` step contributes nothing. The report lists it as a warning.
- A feature is covered when a `success` or `expected_error` step used it.
- A configurable attribute (required, optional, or a block) is covered when
  such a step configured it.
- A computed attribute is covered when a `success` step found it populated
  and no use of that feature type in the step configured it. A configured
  value is copied into state, and configuration can't be paired with state per
  instance (`count`, `for_each`, modules), so any configuration of the
  attribute in the step disqualifies it.
- An optional-computed attribute counts once in the configurable total and
  once in the computed total. Covering both needs one step that sets the
  attribute and one that leaves it to the provider.
- An `expected_error` step keeps its configured attributes and planned
  operations, but its completed operations and populated attributes are
  dropped.
- Deprecated attributes are counted like any other attribute.

Attributes are counted only for top-level schema attributes and blocks.
Functions have no attributes, so they count only as features.

## Matching findings to the schema

Only evidence for providers served from the test binary counts, so the report
measures the provider code under test and not providers that Terraform
downloaded.

The harness serves a provider factory as
`registry.terraform.io/hashicorp/<name>` unless `TF_ACC_PROVIDER_HOST` or
`TF_ACC_PROVIDER_NAMESPACE` is set. The schema is usually generated under the
address the provider is published under, such as
`registry.terraform.io/example-org/example`. To bridge the two, the report maps
each served address to:

1. the schema provider with the same address, or otherwise
2. the only schema provider with the same type (the last address segment).

A served provider that matches neither, such as a provider that only the
tests use, is ignored. If no served provider in any finding maps to the
schema, the report prints a warning that lists the served addresses.

Configuration does not name a provider address. A resource or data source in
configuration is attributed to the served provider whose type matches the
type name prefix, which is the provider Terraform implies for it. For example,
`example_widget` is attributed to the served `example` provider. If several
schema providers define the same type name, the report warns and records
nothing for it.

## Report output

The JSON report has these top-level fields:

- `totals` holds covered and total counts for features, configurable
  attributes, and computed attributes.
- `features` lists every feature in the schema, sorted by provider, kind, and
  name. Each entry has the steps that used it, each attribute with the steps
  that configured it or found it computed, and the steps that planned or
  completed each operation.
- `warnings` lists finding warnings, failed steps, ambiguous matches, and
  provider mapping problems.

The text report prints the totals, one line per feature, one
`uncovered <attribute> (needs ...)` line per uncovered attribute, and the
warnings.

Operations are listed with the steps that performed them but are not turned
into a percentage. The schema doesn't say which operations apply to a
feature, so the report has nothing to divide by.

## Privacy

The collector does not write configuration text, attribute values, plan or
state JSON, query results, resource names, addresses, instance keys, or
assertion matchers. It does write provider source addresses, feature type
names, schema attribute names, operations, test names, and test package import
paths.

Each finding is created with mode `0600`, and the findings directory with mode
`0700` if the harness has to create it. Write findings only to a directory
that you control.

## Limitations

- Ephemeral resources and functions leave nothing in Terraform state. The
  report records that they were used and which attributes were configured,
  but their computed outputs are never covered.
- Nested attributes are not counted. Without a schema-aware walk of the
  state, nested attribute names can't be told apart from user map keys.
- Resource identity attributes are not counted. They appear only in state and
  have no configurable or computed distinction.
- Modules from external sources can appear in plans and state but not in the
  configuration files copied to the working directory. Their resources count
  as used, but their computed attributes are not covered.
- Terraform JSON configuration (a JSON `Config` string or a `.tf.json` file)
  is not read. Such steps get a `config_unavailable` warning and have no
  configuration evidence. Plan and state evidence is still recorded.
- Provider blocks and provider-defined function calls name a provider by its
  local name, which the report compares with the last segment of the source
  address. If `required_providers` assigns a different local name, as in
  `mycloud = { source = "hashicorp/aws" }`, that evidence is not attributed to
  the provider.
- The query file of a `Query` step stays in the working directory. Later steps
  in the same test therefore also record its list resources and their
  configured attributes.
- While collecting, each apply runs two additional `terraform show` commands:
  one for the saved plan and one for the resulting state.

## Example result

Running the tests of `terraform-provider-local` at commit `b3d30fe`, without
the upgrade tests that download released providers, produced 58 findings. The
report showed 9/9 features, 30/30 configurable attributes, and 39/42 computed
attributes covered. The four optional-computed permission attributes appear in
both totals and were covered both when set and when left unset. The three
uncovered computed attributes were outputs of an ephemeral resource. The
completed operations were `create` and `replace` for both file resources,
`read` for the `local_command` data source, and `invoke` for the
`local_command` action.
