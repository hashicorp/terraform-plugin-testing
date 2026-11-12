# Semantic coverage (experimental)

> **Note:** This is a temporary home for this documentation. It will be
> migrated to the official Terraform plugin testing documentation later.

Semantic coverage reports which parts of a provider's schema the acceptance
tests exercise. For each resource, data source, ephemeral resource, list
resource, action, state store, and function, it reports whether the tests
completed every operation of its lifecycle, which of its top-level attributes
a test configured or found populated, and which operations Terraform planned
or completed.

terraform-plugin-testing collects the evidence while the tests run. While
collection is enabled, plugin-testing records evidence for each `TestStep`: it
reads the configuration copied into the step's working directory, and it
inspects the step's plans and the state after each operation. The
`cmd/semantic-coverage-report` command then compares the recorded findings
with a `terraform providers schema -json` document.

## Usage

Collection is off by default. To turn it on, set
`TF_ACC_SEMANTIC_COVERAGE_DIR` to an empty or not yet existing directory when
running acceptance tests. plugin-testing creates the directory if needed and
writes one JSON finding per executed test step.

The report also needs the provider's schema, built from the same source as
the tests. In a provider repository, build the provider and point Terraform at
the binary with a `dev_overrides` CLI configuration. Terraform then uses that
binary without `terraform init`, and needs only a configuration that names
the provider. Replace `example-org/example` and `example` with the provider's
source address and type:

```sh
TF_ACC=1 \
TF_ACC_SEMANTIC_COVERAGE_DIR="$PWD/findings/run-1" \
  go test -count=1 ./internal/provider

go build -o "$PWD/bin/" .
mkdir -p schema
cat > dev.tfrc <<EOF
provider_installation {
  dev_overrides {
    "example-org/example" = "$PWD/bin"
  }
  direct {}
}
EOF
cat > schema/main.tf <<'EOF'
terraform {
  required_providers {
    example = {
      source = "example-org/example"
    }
  }
}
EOF
TF_CLI_CONFIG_FILE=dev.tfrc \
  terraform -chdir=schema providers schema -json > provider-schema.json

go run github.com/hashicorp/terraform-plugin-testing/cmd/semantic-coverage-report \
  -findings-dir findings/run-1 \
  -schema-file provider-schema.json \
  -json-out coverage.json -text-out coverage.txt
```

Terraform prints a warning about the development override to stderr. It
doesn't affect the JSON on stdout.

Use a new, empty findings directory for each run. The report reads every
`.json` file in the findings directory, including subdirectories, and counts
an attribute as covered if any of those findings covers it. Findings left
over from an earlier run would count toward the current one.

Running each test more than once in a run (`go test -count=N` with N > 1) is
not supported. After upgrading terraform-plugin-testing, run the tests again
and use the report from the same version: the report rejects findings whose
format version differs from its own.

While collecting, plugin-testing runs two additional `terraform show` commands
for each apply: one reads the saved plan, so that planned and completed
operations are recorded even for steps without plan checks, and one reads the
resulting state, so that populated attributes are recorded even for steps
without a `Check`. When `TF_ACC_SEMANTIC_COVERAGE_DIR` is unset, nothing is
recorded and these commands don't run.

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
the step number, and the provider **features** that the step used: resources,
data sources, ephemeral resources, list resources, actions, state stores,
providers, and functions. All uses of one feature type within a step are
merged into one entry.

Findings are written only for steps that pass. The automatic destroy after
the last step also writes a finding (see Operations). These steps write none:

- **Steps with `ExpectError`.** They test how the provider handles invalid
  input or a failing operation. A step that only checks a validation error
  shows nothing about whether the feature works with valid input.
- **Failed steps.** A failure can stop Terraform at any point, so the step's
  evidence is incomplete. Once a test has failed, its remaining steps write
  no finding either, even if the failure didn't stop the test.
- **Skipped steps.**

A run with failing tests therefore shows lower coverage, and the report
doesn't say why. Check that `go test` passed before you read the report.

For each feature, the finding lists the following evidence.

**Configured attributes** are the top-level attributes and blocks set in the
step's HCL configuration. The collector reads the configuration files after
plugin-testing writes them to the working directory, so `Config`,
`ConfigFile`, and `ConfigDirectory` are all read. Configuration that
plugin-testing generates itself is not recorded.

**Populated attributes** are the top-level attributes with a non-null value in
the state produced by the step's apply, import, or refresh. The collector
reads only top-level attribute names, because keys further down can be
user-provided map keys. Resources in child modules record no populated
attributes: the collector doesn't read module configuration, so it can't tell
whether a value was set by the module or computed by the provider.

**Operations** are the lifecycle events of a feature, named after the
provider protocol RPCs that implement them. The finding stores them per
feature in two lists. `planned` holds the operations that appear in one of
the step's plans. `completed` holds the operations that Terraform performed
successfully: plugin-testing adds an operation right after the Terraform
command that performs it returns without error. A no-op change records the
feature but no operation.

| Kind | Operations | Recorded when |
|---|---|---|
| provider | `configure` | A plan succeeds that has the provider in its configuration (inferred). |
| resource | `create`, `update`, `replace`, `delete` | A plan has the change, and an apply completes it. |
| resource | `read` | A refresh succeeds: a `RefreshState` step, or the refreshing plan that ends every config step. |
| resource | `import_id` | The import command, or a plan with an import block that has an `id`, succeeds. |
| resource | `import_identity` | A plan with an import block that has an `identity` succeeds. |
| data source | `read` | A refresh succeeds with the data source in state, or an apply completes a read that Terraform deferred. |
| ephemeral resource | `open`, `close` | A plan of a configuration that declares it succeeds (inferred). |
| list resource | `query` | A `Query` step succeeds. |
| action | `plan`, `invoke` | A plan has the invocation (`plan`), and an apply completes it (`invoke`). |
| state store | `configure`, `read`, `write`, `list_states` | A `StateStore` step passes. |
| state store | `delete_state` | A `StateStore` step without `DefaultWorkspaceOnly` passes. |
| state store | `lock`, `unlock` | A `StateStore` step with `VerifyStateStoreLock` passes. |
| function | `call` | A plan of a configuration that calls it succeeds (inferred). |

The collector doesn't observe RPCs, only Terraform's plans and state. The
operations marked "inferred" are not in either: a plan that succeeded must
have configured its providers, opened and closed its ephemeral resources, and
called its functions, so the collector records them for the features in the
step's configuration. Terraform imports a resource while planning an import
block, so a plan with an import completes it, although plugin-testing never
applies such a plan.

A step with `Destroy: true` completes `delete`. The automatic destroy after
the last step writes its own finding, which records a completed `delete` for
each resource in the state before the destroy. Its step number is the one
after the last step's, it has `post_test_destroy: true`, and the report names
it "post-test destroy". It is written only if the test hadn't failed and both
the destroy and `CheckDestroy` succeeded. It is separate from the last step's
finding, so that a failed cleanup can't drop the finding of a passed step.

**Served providers** are recorded once per step rather than per feature. They
are the source addresses under which plugin-testing served the provider
factories from the test binary. Providers that Terraform installs, such as
`ExternalProviders` or the released provider in an upgrade test, are not in
this list.

A finding can also carry warnings. Each one means that a source of evidence
for the step was unavailable, and the report repeats it. The codes don't
record the cause, and nothing logs it:

| Warning | Meaning |
|---|---|
| `config_unavailable` | Some configuration could not be read or parsed. Features and configured attributes may be missing. |
| `plan_unavailable` | A plan could not be read. Planned and completed operations may be missing. |
| `state_unavailable` | The state after an operation could not be read. Populated attributes may be missing. |

## How the report counts coverage

The report applies these rules:

- A feature is covered when steps completed every operation that its kind
  requires (see below). It doesn't matter whether that was in one step or in
  several. Attribute coverage is counted separately and doesn't affect it.
- A configurable attribute (required, optional, or a block) is covered when
  a step configured it.
- A computed attribute is covered when a step found it populated and didn't
  configure it for any use of that feature type. Terraform copies
  configured values into state, so a populated value shows that the provider
  computed it only if the step left the attribute unset. The check is per
  type, not per resource instance, because configuration and state can't be
  matched up per instance under `count`, `for_each`, or modules.
- An optional-computed attribute counts once in the configurable total and
  once in the computed total. Covering both needs one step that sets the
  attribute and one that leaves it unset and finds it populated.
- Deprecated attributes are counted like any other attribute.

Attributes are counted only for top-level schema attributes and blocks.
Functions take positional arguments rather than attributes, and the collector
doesn't record which arguments a call passes. A function therefore adds
nothing to the attribute totals.

A feature needs these completed operations to be covered:

| Kind | Required operations |
|---|---|
| provider | `configure` |
| resource | `create`, `read`, `update`, `delete`, `import_id`, and `import_identity` if the schema has a resource identity schema for it |
| data source | `read` |
| ephemeral resource | `open`, `close` |
| list resource | `query` |
| action | `plan`, `invoke` |
| state store | `configure`, `read`, `write`, `list_states`, `delete_state` |
| function | `call` |

A planned operation doesn't count. `replace`, `lock`, and `unlock` are
reported but not required.

Every resource needs `update` and `import_id`, although some resources
support neither: a resource whose attributes all require replacement never
updates in place, and a resource can lack import support. The schema JSON
doesn't say whether either applies, so such a resource is never covered. Its
`missing_operations` show why.

## Matching findings to the schema

Only evidence for providers served from the test binary counts, so the report
measures the provider code under test and not providers that Terraform
downloaded.

plugin-testing serves a provider factory as
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

- `version` is the format version of the report.
- `attribute_scope` is always `top_level`: only top-level attributes and
  blocks are counted.
- `findings` is the number of findings read.
- `totals` holds covered and total counts for features, configurable
  attributes, and computed attributes.
- `features` lists every feature in the schema, sorted by provider, kind, and
  name. Each entry has the steps whose evidence names it (`used_by`), the
  required operations that no step completed (`missing_operations`), each
  attribute with the steps that configured it (`configured_by`) or found it
  computed (`computed_by`), and the steps that planned or completed each
  operation.
- `warnings` lists finding warnings, ambiguous matches, and provider mapping
  problems. It is omitted when there are none.

The text report prints the totals and one line per feature. Under a feature,
it prints a `missing operations [...]` line if required operations are
missing, and one `uncovered <attribute> (needs ...)` line per uncovered
attribute. The warnings come last.

## Privacy

The collector does not write configuration text, attribute values, plan or
state JSON, query results, resource names, addresses, instance keys, or
assertion matchers. It does write provider source addresses, feature type
names, schema attribute names, operations, test names, and test package import
paths.

Each finding is created with mode `0600`, and the findings directory with mode
`0700` if plugin-testing has to create it. Write findings only to a directory
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
- Provider `configure`, ephemeral resource `open` and `close`, and function
  `call` are inferred from successful plans, not observed (see Operations).
  A provider that Terraform doesn't need for the plan, or an ephemeral
  resource or function call that the plan doesn't evaluate, still counts.
- Renewing an ephemeral resource is not recorded. It happens only when the
  resource's lease expires during an operation.
- The read of a resource is recorded only from steps that refresh. A config
  step with the `NoRefresh` plan option, or a destroy step with
  `PreventPostDestroyRefresh`, records none.

## Example result

Running the tests of `terraform-provider-local` at commit `b3d30fe`, without
the upgrade tests that download released providers (`go test -skip Upgrade`),
produced 73 findings: 40 for test steps and 33 for post-test destroys. The 18
steps with `ExpectError` wrote none. The schema was generated with the
`dev_overrides` recipe from [Usage](#usage).

The report showed 7/9 features, 30/30 configurable attributes, and 39/42
computed attributes covered:

- The two file resources are not covered. Their tests create, read, replace,
  and delete them, but never update or import them: every change to a file
  resource forces replacement, and neither supports import. Both are missing
  `update` and `import_id`, as described in
  [How the report counts coverage](#how-the-report-counts-coverage).
- The four optional-computed permission attributes appear in both totals and
  were covered both when set and when left unset.
- The three uncovered computed attributes are outputs of the ephemeral
  resource. Ephemeral resources leave nothing in state (see
  [Limitations](#limitations)).
- The `local_file` and `local_sensitive_file` data sources were each used in
  one test, whose refreshing plan recorded their `read`.
- `delete` of the file resources was completed only by post-test destroys.
  No test has a `Destroy: true` step.

The text report:

```text
Test-derived semantic coverage: 73 findings
Attribute scope: top-level schema attributes and blocks only; nested descendants are excluded.
Features 7/9; configurable 30/30; computed 39/42
registry.terraform.io/hashicorp/local action local_command: covered=true attributes=5/5 planned=[invoke] completed=[invoke plan]
registry.terraform.io/hashicorp/local data_source local_command: covered=true attributes=8/8 planned=[read] completed=[read]
registry.terraform.io/hashicorp/local data_source local_file: covered=true attributes=10/10 planned=[] completed=[read]
registry.terraform.io/hashicorp/local data_source local_sensitive_file: covered=true attributes=10/10 planned=[] completed=[read]
registry.terraform.io/hashicorp/local ephemeral_resource local_command: covered=true attributes=5/8 planned=[] completed=[close open]
  uncovered exit_code (needs a test that leaves it unset and finds it populated)
  uncovered stderr (needs a test that leaves it unset and finds it populated)
  uncovered stdout (needs a test that leaves it unset and finds it populated)
registry.terraform.io/hashicorp/local function direxists: covered=true attributes=0/0 planned=[] completed=[call]
registry.terraform.io/hashicorp/local provider : covered=true attributes=0/0 planned=[] completed=[configure]
registry.terraform.io/hashicorp/local resource local_file: covered=false attributes=14/14 planned=[create replace] completed=[create delete read replace]
  missing operations [update import_id]
registry.terraform.io/hashicorp/local resource local_sensitive_file: covered=false attributes=13/13 planned=[create replace] completed=[create delete read replace]
  missing operations [update import_id]
```

An excerpt of the JSON report. It shows one feature of each kind, with some of
their attributes: for `local_file`, a configurable attribute (`filename`), an
optional-computed one (`file_permission`), and a computed one (`id`). Each
list of steps is cut to its first entry. The full file has 3504 lines.

```json
{
  "version": 1,
  "attribute_scope": "top_level",
  "findings": 73,
  "totals": {
    "features": {
      "covered": 7,
      "total": 9
    },
    "configurable": {
      "covered": 30,
      "total": 30
    },
    "computed": {
      "covered": 39,
      "total": 42
    }
  },
  "features": [
    {
      "provider": "registry.terraform.io/hashicorp/local",
      "kind": "action",
      "name": "local_command",
      "covered": true,
      "used_by": [
        {
          "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
          "test": "TestLocalCommandAction_absolute_path_bash",
          "step": 1
        }
      ],
      "attributes": [
        {
          "path": "command",
          "configurable": true,
          "computed": false,
          "covered": true,
          "configured_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalCommandAction_absolute_path_bash",
              "step": 1
            }
          ]
        }
      ],
      "planned": {
        "invoke": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandAction_absolute_path_bash",
            "step": 1
          }
        ]
      },
      "completed": {
        "invoke": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandAction_absolute_path_bash",
            "step": 1
          }
        ],
        "plan": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandAction_absolute_path_bash",
            "step": 1
          }
        ]
      }
    },
    {
      "provider": "registry.terraform.io/hashicorp/local",
      "kind": "data_source",
      "name": "local_file",
      "covered": true,
      "used_by": [
        {
          "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
          "test": "TestLocalFileDataSource",
          "step": 1
        }
      ],
      "attributes": [
        {
          "path": "content",
          "configurable": false,
          "computed": true,
          "covered": true,
          "computed_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalFileDataSource",
              "step": 1
            }
          ]
        },
        {
          "path": "filename",
          "configurable": true,
          "computed": false,
          "covered": true,
          "configured_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalFileDataSource",
              "step": 1
            }
          ]
        }
      ],
      "completed": {
        "read": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalFileDataSource",
            "step": 1
          }
        ]
      }
    },
    {
      "provider": "registry.terraform.io/hashicorp/local",
      "kind": "ephemeral_resource",
      "name": "local_command",
      "covered": true,
      "used_by": [
        {
          "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
          "test": "TestLocalCommandEphemeral_absolute_path_with_working_directory",
          "step": 1
        }
      ],
      "attributes": [
        {
          "path": "command",
          "configurable": true,
          "computed": false,
          "covered": true,
          "configured_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalCommandEphemeral_absolute_path_with_working_directory",
              "step": 1
            }
          ]
        },
        {
          "path": "stdout",
          "configurable": false,
          "computed": true,
          "covered": false
        }
      ],
      "completed": {
        "close": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandEphemeral_absolute_path_with_working_directory",
            "step": 1
          }
        ],
        "open": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandEphemeral_absolute_path_with_working_directory",
            "step": 1
          }
        ]
      }
    },
    {
      "provider": "registry.terraform.io/hashicorp/local",
      "kind": "function",
      "name": "direxists",
      "covered": true,
      "used_by": [
        {
          "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
          "test": "TestDirectoryExists_basic",
          "step": 1
        }
      ],
      "completed": {
        "call": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestDirectoryExists_basic",
            "step": 1
          }
        ]
      }
    },
    {
      "provider": "registry.terraform.io/hashicorp/local",
      "kind": "provider",
      "name": "",
      "covered": true,
      "used_by": [
        {
          "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
          "test": "TestDirectoryExists_basic",
          "step": 1
        }
      ],
      "completed": {
        "configure": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestDirectoryExists_basic",
            "step": 1
          }
        ]
      }
    },
    {
      "provider": "registry.terraform.io/hashicorp/local",
      "kind": "resource",
      "name": "local_file",
      "covered": false,
      "used_by": [
        {
          "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
          "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
          "step": 1
        }
      ],
      "missing_operations": [
        "update",
        "import_id"
      ],
      "attributes": [
        {
          "path": "file_permission",
          "configurable": true,
          "computed": true,
          "covered": true,
          "configured_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalFile_Permissions",
              "step": 2
            }
          ],
          "computed_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
              "step": 1
            }
          ]
        },
        {
          "path": "filename",
          "configurable": true,
          "computed": false,
          "covered": true,
          "configured_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
              "step": 1
            }
          ]
        },
        {
          "path": "id",
          "configurable": false,
          "computed": true,
          "covered": true,
          "computed_by": [
            {
              "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
              "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
              "step": 1
            }
          ]
        }
      ],
      "planned": {
        "create": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
            "step": 1
          }
        ],
        "replace": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalFile_Basic",
            "step": 2
          }
        ]
      },
      "completed": {
        "create": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
            "step": 1
          }
        ],
        "delete": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
            "step": 2,
            "post_test_destroy": true
          }
        ],
        "read": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalCommandDataSource_absolute_path_with_working_directory",
            "step": 1
          }
        ],
        "replace": [
          {
            "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
            "test": "TestLocalFile_Basic",
            "step": 2
          }
        ]
      }
    }
  ]
}
```
