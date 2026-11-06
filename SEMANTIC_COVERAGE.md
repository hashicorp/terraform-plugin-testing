# Semantic coverage (experimental)

> **Note:** This is a temporary home for this documentation. It will be
> migrated to the official Terraform plugin testing documentation later.

Semantic coverage reports which parts of a provider's schema the acceptance
tests exercise. For each resource, data source, ephemeral resource, list
resource, action, state store, and function, it reports whether a test used
it, which of its top-level attributes a test configured or found populated,
and which operations Terraform planned or completed.

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

Findings are written only for steps that pass. These steps write none:

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

**Operations** are `create`, `update`, `replace`, `delete`, `read`, `import`,
`refresh`, `query`, and `invoke`. The finding stores them per feature in two
lists. `planned` holds the operations that appear in one of the step's plans.
`completed` holds the operations that Terraform performed successfully:
plugin-testing adds an operation right after the Terraform command that
performs it returns without error, which is applying a plan, running an
import command, refreshing, or running a query. A no-op change records the
feature but no operation. Steps with `PlanOnly`, and import-block tests that
only plan, record planned operations only.

A step with `Destroy: true` completes `delete`. The automatic destroy after
the last step is not recorded, although it calls the provider's delete for
every resource left in state. It runs after the last step's finding is
written, so that a failed cleanup can't drop the finding of a passed step,
and it belongs to no step. As a result, `delete` shows as completed only for
tests that have a `Destroy: true` step.

**Served providers** are recorded once per step rather than per feature. They
are the source addresses under which plugin-testing served the provider
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

- A feature is covered when a step used it.
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
  attribute and one that leaves it to the provider.
- Deprecated attributes are counted like any other attribute.

Attributes are counted only for top-level schema attributes and blocks.
Functions take positional arguments rather than attributes, and the collector
doesn't record which arguments a call passes. A function is therefore covered
when a step calls it, and adds nothing to the attribute totals.

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
  name. Each entry has the steps that used it (`covered_by`), each attribute
  with the steps that configured it (`configured_by`) or found it computed
  (`computed_by`), and the steps that planned or completed each operation.
- `warnings` lists finding warnings, ambiguous matches, and provider mapping
  problems. It is omitted when there are none.

The text report prints the totals, one line per feature, one
`uncovered <attribute> (needs ...)` line per uncovered attribute, and the
warnings.

Operations are listed with the steps that performed them but are not turned
into a percentage. Some operations apply to every feature of a kind: `create`
and `delete` to resources, `query` to list resources, and `invoke` to actions.
Others depend on how the provider implements the feature, and the schema JSON
doesn't say whether they apply: `update` needs an attribute that can change in
place, `replace` needs an attribute that requires replacement, and `import`
needs import support. Ephemeral resources are opened, renewed, and closed,
which the collector doesn't record, and the `read` of a data source is
recorded only in some cases (see Limitations). The report therefore has no
complete list of applicable operations to divide by.

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
- A data source that Terraform reads while planning doesn't appear as a change
  in the plan, so its `read` is recorded only when Terraform defers the read
  to apply, for example because it depends on a resource that the same apply
  creates or changes.
- The automatic destroy after the last step is not recorded (see
  Operations).

## Example result

Running the tests of `terraform-provider-local` at commit `b3d30fe`, without
the upgrade tests that download released providers (`go test -skip Upgrade`),
produced 40 findings; the 18 steps with `ExpectError` wrote none. The schema
was generated with the `dev_overrides` recipe from [Usage](#usage). The report showed 9/9 features, 30/30
configurable attributes, and 39/42 computed attributes covered. The four
optional-computed permission attributes appear in both totals and were covered
both when set and when left unset. The three uncovered computed attributes
were outputs of an ephemeral resource. The completed operations were `create`
and `replace` for both file resources, `read` for the `local_command` data
source, and `invoke` for the `local_command` action. The `local_command` data
source records `read` only in the tests where it depends on a `local_file`
created in the same apply. The `local_file` and `local_sensitive_file` data
sources show no operations, because Terraform read them while planning.

The text report:

```text
Test-derived semantic coverage: 40 findings
Attribute scope: top-level schema attributes and blocks only; nested descendants are excluded.
Features 9/9; configurable 30/30; computed 39/42
registry.terraform.io/hashicorp/local action local_command: covered=true attributes=5/5 planned=[invoke] completed=[invoke]
registry.terraform.io/hashicorp/local data_source local_command: covered=true attributes=8/8 planned=[read] completed=[read]
registry.terraform.io/hashicorp/local data_source local_file: covered=true attributes=10/10 planned=[] completed=[]
registry.terraform.io/hashicorp/local data_source local_sensitive_file: covered=true attributes=10/10 planned=[] completed=[]
registry.terraform.io/hashicorp/local ephemeral_resource local_command: covered=true attributes=5/8 planned=[] completed=[]
  uncovered exit_code (needs a test that leaves it to the provider)
  uncovered stderr (needs a test that leaves it to the provider)
  uncovered stdout (needs a test that leaves it to the provider)
registry.terraform.io/hashicorp/local function direxists: covered=true attributes=0/0 planned=[] completed=[]
registry.terraform.io/hashicorp/local provider : covered=true attributes=0/0 planned=[] completed=[]
registry.terraform.io/hashicorp/local resource local_file: covered=true attributes=14/14 planned=[create replace] completed=[create replace]
registry.terraform.io/hashicorp/local resource local_sensitive_file: covered=true attributes=13/13 planned=[create replace] completed=[create replace]
```

An excerpt of the JSON report. It is complete except that `features` shows
only two of the nine entries; the full file has 2,868 lines.
<!-- <review> -->
<!-- also show a resource (ideally one with optional computed, computed and non computed attributes (or the best possible approximation) -->
<!-- if the provider has an action or other features, also show one of each -->
<!-- </review> -->

```json
{
  "version": 1,
  "attribute_scope": "top_level",
  "findings": 40,
  "totals": {
    "features": {
      "covered": 9,
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
      "kind": "data_source",
      "name": "local_file",
      "covered": true,
      "covered_by": [
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
          "path": "content_base64",
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
          "path": "content_base64sha256",
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
          "path": "content_base64sha512",
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
          "path": "content_md5",
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
          "path": "content_sha1",
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
          "path": "content_sha256",
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
          "path": "content_sha512",
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
        },
        {
          "path": "id",
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
        }
      ]
    },
    {
      "provider": "registry.terraform.io/hashicorp/local",
      "kind": "function",
      "name": "direxists",
      "covered": true,
      "covered_by": [
        {
          "package": "github.com/terraform-providers/terraform-provider-local/internal/provider",
          "test": "TestDirectoryExists_basic",
          "step": 1
        }
      ]
    }
  ]
}
```
