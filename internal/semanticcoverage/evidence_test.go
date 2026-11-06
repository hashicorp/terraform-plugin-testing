// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestConfigurationEvidence(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		files        map[string]string
		want         []Feature
		wantWarnings []Warning
	}{
		"blocks declare features": {
			files: map[string]string{"main.tf": `
provider "example" { region = "x" }
resource "example_item" "a" { name = "x" }
data "example_item" "b" { id = "x" }
ephemeral "example_secret" "c" { length = 1 }
action "example_action" "d" {
  config { message = "x" }
}
`},
			want: []Feature{
				{Kind: KindAction, Name: "example_action", Configured: []string{"message"}},
				{Kind: KindDataSource, Name: "example_item", Configured: []string{"id"}},
				{Kind: KindEphemeralResource, Name: "example_secret", Configured: []string{"length"}},
				{Kind: KindProvider, Name: "example", Configured: []string{"region"}},
				{Kind: KindResource, Name: "example_item", Configured: []string{"name"}},
			},
		},
		"uses of one type merge": {
			files: map[string]string{"main.tf": `
resource "example_item" "a" { name = "x" }
resource "example_item" "b" { tags = {} }
`},
			want: []Feature{{Kind: KindResource, Name: "example_item", Configured: []string{"name", "tags"}}},
		},
		"only top-level attributes and blocks are configured": {
			files: map[string]string{"main.tf": `
resource "example_item" "a" {
  settings {
    mode = "x"
  }
}
`},
			want: []Feature{{Kind: KindResource, Name: "example_item", Configured: []string{"settings"}}},
		},
		"dynamic blocks are named by their label": {
			files: map[string]string{"main.tf": `
resource "example_item" "a" {
  dynamic "rule" {
    for_each = []
    content {}
  }
}
`},
			want: []Feature{{Kind: KindResource, Name: "example_item", Configured: []string{"rule"}}},
		},
		"blocks in a feature are schema blocks": {
			files: map[string]string{"main.tf": `
resource "example_item" "a" {
  data "example_other" "b" {}
}
`},
			want: []Feature{{Kind: KindResource, Name: "example_item", Configured: []string{"data"}}},
		},
		"check and terraform blocks declare features": {
			files: map[string]string{"main.tf": `
check "health" {
  data "example_item" "a" { id = "x" }
}
terraform {
  state_store "example_store" {
    provider "example" {}
    path = "x"
  }
}
`},
			want: []Feature{
				{Kind: KindDataSource, Name: "example_item", Configured: []string{"id"}},
				{Kind: KindProvider, Name: "example"},
				{Kind: KindStateStore, Name: "example_store", Configured: []string{"path", "provider"}},
			},
		},
		"provider functions are distinct per provider": {
			files: map[string]string{"main.tf": `
output "a" { value = provider::example::parse("x") }
output "b" { value = provider::other::parse(provider::example::parse("x")) }
`},
			want: []Feature{
				{Kind: KindFunction, Name: "parse", Provider: "example"},
				{Kind: KindFunction, Name: "parse", Provider: "other"},
			},
		},
		"list resources in query files, without providers": {
			files: map[string]string{"main.tfquery.hcl": `
provider "example" {}
list "example_item" "a" {
  provider = example
  config { filter = "x" }
}
`},
			want: []Feature{{Kind: KindListResource, Name: "example_item", Configured: []string{"filter"}}},
		},
		"JSON configuration is not read": {
			files: map[string]string{
				"a.tf.json": `{"resource": {"example_item": {"a": {"name": "x"}}}}`,
				"b.tf":      `resource "example_other" "b" {}`,
			},
			want:         []Feature{{Kind: KindResource, Name: "example_other"}},
			wantWarnings: []Warning{WarningConfigUnavailable},
		},
		"unparsable configuration": {
			files:        map[string]string{"main.tf": `resource "example_item" {`},
			wantWarnings: []Warning{WarningConfigUnavailable},
		},
		"other files are not read": {
			files: map[string]string{"terraform.tfvars": `name = "x"`},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := configFiles(t, test.files)
			got := stepFinding(t, func(s *Step) { s.ObserveConfigDir(dir) })
			if diff := cmp.Diff(test.want, got.Features); diff != "" {
				t.Errorf("features mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantWarnings, got.Warnings); diff != "" {
				t.Errorf("warnings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOperationEvidence(t *testing.T) {
	t.Parallel()

	create := plan(itemChange(tfjson.ActionCreate))
	noOp := plan(itemChange(tfjson.ActionNoop))
	item := func(planned, completed []Operation) []Feature {
		return []Feature{{Kind: KindResource, Name: "example_item", Provider: exampleProvider, Planned: planned, Completed: completed}}
	}

	tests := map[string]struct {
		record       func(s *Step)
		want         []Feature
		wantWarnings []Warning
	}{
		"a plan plans its operations": {
			record: func(s *Step) { s.ObservePlan(create) },
			want:   item([]Operation{OperationCreate}, nil),
		},
		"an applied plan completes its operations": {
			record: func(s *Step) {
				s.ObservePlan(create)
				s.ApplyCompleted(create)
				s.ObservePlan(noOp)
			},
			want: item([]Operation{OperationCreate}, []Operation{OperationCreate}),
		},
		"a no-op change has no operation": {
			record: func(s *Step) {
				s.ObservePlan(noOp)
				s.ApplyCompleted(noOp)
			},
			want: item(nil, nil),
		},
		"a destroy completes a delete": {
			record: func(s *Step) {
				destroy := plan(itemChange(tfjson.ActionDelete))
				s.ObservePlan(destroy)
				s.ApplyCompleted(destroy)
			},
			want: item([]Operation{OperationDelete}, []Operation{OperationDelete}),
		},
		// Terraform imports while planning.
		"an import block with an ID imports by ID": {
			record: func(s *Step) {
				p := plan(itemChange(tfjson.ActionNoop))
				p.ResourceChanges[0].Change.Importing = &tfjson.Importing{ID: "x"}
				s.ObservePlan(p)
			},
			want: item([]Operation{OperationImportID}, []Operation{OperationImportID}),
		},
		"an import block with an identity imports by identity": {
			record: func(s *Step) {
				p := plan(itemChange(tfjson.ActionNoop))
				p.ResourceChanges[0].Change.Importing = &tfjson.Importing{Identity: map[string]any{"id": "x"}}
				s.ObservePlan(p)
			},
			want: item([]Operation{OperationImportIdentity}, []Operation{OperationImportIdentity}),
		},
		"a deferred data source read": {
			record: func(s *Step) {
				c := itemChange(tfjson.ActionRead)
				c.Mode = tfjson.DataResourceMode
				s.ObservePlan(plan(c))
			},
			want: []Feature{{Kind: KindDataSource, Name: "example_item", Provider: exampleProvider, Planned: []Operation{OperationRead}}},
		},
		"an applied action invocation": {
			record: func(s *Step) {
				p := &tfjson.Plan{ActionInvocations: []*tfjson.ActionInvocation{{Address: "action.example_action.a", Type: "example_action", ProviderName: exampleProvider}}}
				s.ObservePlan(p)
				s.ApplyCompleted(p)
			},
			// A plan that lists an invocation planned the action.
			want: []Feature{{Kind: KindAction, Name: "example_action", Provider: exampleProvider, Planned: []Operation{OperationInvoke}, Completed: []Operation{OperationInvoke, OperationPlan}}},
		},
		"a plan configures its providers": {
			record: func(s *Step) {
				s.ObservePlan(&tfjson.Plan{Config: &tfjson.Config{ProviderConfigs: map[string]*tfjson.ProviderConfig{
					"example": {Name: "example", FullName: exampleProvider},
				}}})
			},
			want: []Feature{{Kind: KindProvider, Name: exampleProvider, Completed: []Operation{OperationConfigure}}},
		},
		"an unreadable plan": {
			record: func(s *Step) {
				s.ObservePlan(nil)
				s.ApplyCompleted(nil)
			},
			wantWarnings: []Warning{WarningPlanUnavailable},
		},
		"an import command imports by ID": {
			record: func(s *Step) {
				imported := state(resource("example_item", map[string]any{"id": "x"}))
				s.ObserveState(imported)
				s.ImportCompleted(imported, "example_item.test")
				s.ImportCompleted(imported, "example_item.other")
			},
			want: []Feature{{Kind: KindResource, Name: "example_item", Provider: exampleProvider, Populated: []string{"id"}, Completed: []Operation{OperationImportID}}},
		},
		"a refresh reads the resources and data sources in its state": {
			record: func(s *Step) {
				data := resource("example_lookup", nil)
				data.Mode = tfjson.DataResourceMode
				s.ObserveState(state(resource("example_other", nil)))
				s.RefreshCompleted(state(resource("example_item", nil), data))
			},
			want: []Feature{
				{Kind: KindDataSource, Name: "example_lookup", Provider: exampleProvider, Completed: []Operation{OperationRead}},
				{Kind: KindResource, Name: "example_item", Provider: exampleProvider, Completed: []Operation{OperationRead}},
				{Kind: KindResource, Name: "example_other", Provider: exampleProvider},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := stepFinding(t, test.record)
			if diff := cmp.Diff(test.want, got.Features); diff != "" {
				t.Errorf("features mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantWarnings, got.Warnings); diff != "" {
				t.Errorf("warnings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlanningInfersEphemeralResourcesAndFunctions(t *testing.T) {
	t.Parallel()

	dir := configDir(t, `
ephemeral "example_secret" "a" {}
output "b" { value = provider::example::parse("x") }
`)
	tests := map[string]struct {
		plan *tfjson.Plan
		want []Feature
	}{
		"a plan opens and closes ephemeral resources and calls functions": {
			plan: plan(),
			want: []Feature{
				{Kind: KindEphemeralResource, Name: "example_secret", Completed: []Operation{OperationClose, OperationOpen}},
				{Kind: KindFunction, Name: "parse", Provider: "example", Completed: []Operation{OperationCall}},
			},
		},
		"without a plan they are only declared": {
			want: []Feature{
				{Kind: KindEphemeralResource, Name: "example_secret"},
				{Kind: KindFunction, Name: "parse", Provider: "example"},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := stepFinding(t, func(s *Step) {
				s.ObserveConfigDir(dir)
				if test.plan != nil {
					s.ObservePlan(test.plan)
				}
			})
			if diff := cmp.Diff(test.want, got.Features); diff != "" {
				t.Errorf("features mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStateStoreEvidence(t *testing.T) {
	t.Parallel()

	dir := configDir(t, `
terraform {
  state_store "example_store" {
    provider "example" {}
  }
}
`)
	tests := map[string]struct {
		record func(s *Step)
		want   []Operation
	}{
		"default workspace only": {
			record: func(s *Step) { s.StateStoreCompleted(false) },
			want:   []Operation{OperationConfigure, OperationListStates, OperationRead, OperationWrite},
		},
		"workspaces": {
			record: func(s *Step) { s.StateStoreCompleted(true) },
			want:   []Operation{OperationConfigure, OperationDeleteState, OperationListStates, OperationRead, OperationWrite},
		},
		"verified lock": {
			record: func(s *Step) { s.StateStoreLockVerified() },
			want:   []Operation{OperationLock, OperationUnlock},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := stepFinding(t, func(s *Step) {
				s.ObserveConfigDir(dir)
				test.record(s)
			})
			want := []Feature{
				{Kind: KindProvider, Name: "example"},
				{Kind: KindStateStore, Name: "example_store", Configured: []string{"provider"}, Completed: test.want},
			}
			if diff := cmp.Diff(want, got.Features); diff != "" {
				t.Errorf("features mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestQueryCompletesListResources(t *testing.T) {
	t.Parallel()

	dir := configFiles(t, map[string]string{"main.tfquery.hcl": `list "example_item" "a" {}`})
	got := stepFinding(t, func(s *Step) {
		s.ObserveQueryConfig(filepath.Join(dir, "main.tfquery.hcl"))
		s.ObserveState(state(resource("example_item", nil)))
		s.QueryCompleted()
	})
	want := []Feature{
		{Kind: KindListResource, Name: "example_item", Completed: []Operation{OperationQuery}},
		{Kind: KindResource, Name: "example_item", Provider: exampleProvider},
	}
	if diff := cmp.Diff(want, got.Features); diff != "" {
		t.Errorf("features mismatch (-want +got):\n%s", diff)
	}
}

func TestStateEvidence(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		record       func(s *Step)
		want         []Feature
		wantWarnings []Warning
	}{
		"non-null top-level attributes are populated": {
			record: func(s *Step) {
				s.ObserveState(state(resource("example_item", map[string]any{"id": "x", "unset": nil, "tags": map[string]any{"key": "x"}})))
			},
			want: []Feature{{Kind: KindResource, Name: "example_item", Provider: exampleProvider, Populated: []string{"id", "tags"}}},
		},
		"resources in child modules have no populated attributes": {
			record: func(s *Step) {
				r := resource("example_item", map[string]any{"id": "x"})
				r.Address = "module.child." + r.Address
				st := state()
				st.Values.RootModule.ChildModules = []*tfjson.StateModule{{Address: "module.child", Resources: []*tfjson.StateResource{r}}}
				s.ObserveState(st)
			},
			want: []Feature{{Kind: KindResource, Name: "example_item", Provider: exampleProvider}},
		},
		"empty state": {
			record: func(s *Step) { s.ObserveState(&tfjson.State{}) },
		},
		"unreadable state": {
			record:       func(s *Step) { s.ObserveState(nil) },
			wantWarnings: []Warning{WarningStateUnavailable},
		},
		"plan values are not state": {
			record: func(s *Step) {
				p := plan(itemChange(tfjson.ActionCreate))
				p.ResourceChanges[0].Change.After = map[string]any{"id": "x"}
				s.ObservePlan(p)
			},
			want: []Feature{{Kind: KindResource, Name: "example_item", Provider: exampleProvider, Planned: []Operation{OperationCreate}}},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := stepFinding(t, test.record)
			if diff := cmp.Diff(test.want, got.Features); diff != "" {
				t.Errorf("features mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantWarnings, got.Warnings); diff != "" {
				t.Errorf("warnings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfigurationAndStateOfOneTypeMerge(t *testing.T) {
	t.Parallel()

	dir := configDir(t, `resource "example_item" "a" { name = "x" }`)
	got := stepFinding(t, func(s *Step) {
		s.ObserveConfigDir(dir)
		s.ObserveState(state(resource("example_item", map[string]any{"id": "x", "name": "x"})))
	})
	want := []Feature{{Kind: KindResource, Name: "example_item", Provider: exampleProvider, Configured: []string{"name"}, Populated: []string{"id", "name"}}}
	if diff := cmp.Diff(want, got.Features); diff != "" {
		t.Errorf("features mismatch (-want +got):\n%s", diff)
	}
}
