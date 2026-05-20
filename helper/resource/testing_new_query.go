// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource/query"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
)

func testStepNewQuery(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories) error {
	t.Helper()

	queryConfigRequest := teststep.ConfigurationRequest{
		Raw: &step.Config,
	}
	err := wd.SetQuery(ctx, teststep.Configuration(queryConfigRequest), step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting query config: %w", err)
	}

	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Init(ctx)
	})
	if err != nil {
		t.Fatalf("Error getting init: %s", err)
	}

	var queryOut []tfjson.LogMsg

	if step.GenerateConfig {
		// Terraform only populates ListResourceFoundData.Config in the JSON
		// stream when the -generate-config-out flag is supplied. Pass it a
		// scratch path inside an isolated temp directory: we discard the
		// file Terraform writes and assemble the final configuration from
		// the streamed messages via ConcatConfigs.
		generateOutDir, err := os.MkdirTemp("", "tf-query-generate-")
		if err != nil {
			return fmt.Errorf("Error creating temp dir for generate-config-out: %w", err)
		}
		generateOutPath := filepath.Join(generateOutDir, "generated.tf")

		err = runProviderCommand(ctx, t, wd, providers, func() error {
			var err error
			queryOut, err = wd.Query(ctx, tfexec.GenerateConfigOut(generateOutPath))
			return err
		})
		if err != nil {
			return err
		}

		found := make([]tfjson.ListResourceFoundData, 0, len(queryOut))
		for _, msg := range queryOut {
			if v, ok := msg.(tfjson.ListResourceFoundMessage); ok {
				found = append(found, v.ListResourceFound)
			}
		}
		concatConfig, err := ConcatConfigs(found, step.QueryConfigFilters)
		if err != nil {
			return fmt.Errorf("Error generating concatenated config: %w", err)
		}

		newConfigRequest := teststep.ConfigurationRequest{
			Raw: concatConfig,
		}
		err = wd.SetConfig(ctx, teststep.Configuration(newConfigRequest), step.ConfigVariables)
		if err != nil {
			return fmt.Errorf("error setting generated config: %w", err)
		}

		err = runProviderCommand(ctx, t, wd, providers, func() error {
			return wd.Init(ctx)
		})
		if err != nil {
			return fmt.Errorf("error running init against generated config: %w", err)
		}

		if err := runProviderCommandCreatePlan(ctx, t, wd, providers); err != nil {
			return fmt.Errorf("error creating plan from generated config: %w", err)
		}

		plan, err := runProviderCommandSavedPlan(ctx, t, wd, providers)
		if err != nil {
			return fmt.Errorf("error reading generated plan: %w", err)
		}

		if len(plan.ResourceChanges) == 0 {
			return fmt.Errorf("generated config plan: expected resource changes, got none")
		}

		// loop through plan changes and throw an error if not a no-op
		for _, rc := range plan.ResourceChanges {
			if !rc.Change.Actions.NoOp() {
				return fmt.Errorf("expected no-op from generated config testing, got action: %s for type: %s", rc.Change.Actions, rc.Address)
			}
		}
	} else {
		err = runProviderCommand(ctx, t, wd, providers, func() error {
			var err error
			queryOut, err = wd.Query(ctx)
			return err
		})
		if err != nil {
			return err
		}
	}

	return query.RunQueryChecks(ctx, t, queryOut, step.QueryResultChecks)
}

var resourceBlockHeaderRegex = regexp.MustCompile(`(resource\s*"[a-zA-Z0-9_-]+"\s*)"[a-zA-Z0-9_-]+"(\s*{)`)

func ConcatConfigs(found []tfjson.ListResourceFoundData, filters map[string]queryfilter.QueryFilter) (*string, error) {
	ctx := context.Background()

	configs := make([]string, 0, len(found))

	for name, filter := range filters {

		if name == "" {
			return nil, fmt.Errorf("filters map key must not be empty")
		}

		filteredResults := make([]tfjson.ListResourceFoundData, 0)

		for _, result := range found {
			keepResult := false

			resp := queryfilter.FilterQueryResponse{}
			filter.Filter(ctx, queryfilter.FilterQueryRequest{QueryItem: result}, &resp)

			if resp.Include {
				keepResult = true
			}

			if resp.Error != nil {
				return nil, resp.Error
			}

			if keepResult {
				filteredResults = append(filteredResults, result)
			}
		}

		if len(filteredResults) == 0 {
			return nil, fmt.Errorf("%s filter returned no results, filters must return exactly 1 result", name)
		}

		if len(filteredResults) > 1 {
			return nil, fmt.Errorf("%s filter returned %d results, filters must return exactly 1 result", name, len(filteredResults))
		}

		matchedItem := filteredResults[0]
		trimmed := strings.TrimSpace(matchedItem.Config)
		if trimmed == "" {
			continue
		}

		rewritten, ok := rewriteResourceName(trimmed, name)
		if !ok {
			return nil, fmt.Errorf("ConcatConfigs: could not locate a `resource` block to rename in Config for %q", matchedItem.Address)
		}
		configs = append(configs, rewritten)

	}

	concatenated := strings.Join(configs, "\n\n")

	return &concatenated, nil
}

// rewriteResourceName rewrites the resource name (second label) of the first
// `resource` block found in cfg to newName. The bool return value is false
// when no resource block header was located in cfg.
func rewriteResourceName(cfg, newName string) (string, bool) {
	loc := resourceBlockHeaderRegex.FindStringSubmatchIndex(cfg)
	if loc == nil {
		return cfg, false
	}

	// loc indexes: 0,1=full match; 2,3=group1 (`resource "<type>" `);
	// 4,5=group2 (`{`)
	prefix := cfg[loc[2]:loc[3]]
	suffix := cfg[loc[4]:loc[5]]

	return cfg[:loc[0]] + prefix + `"` + newName + `"` + suffix + cfg[loc[1]:], true
}
