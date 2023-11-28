// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// mergedConfig prepends any necessary terraform configuration blocks to the
// TestStep Config.
//
// If there are ExternalProviders configurations in either the TestCase or
// TestStep, the terraform configuration block should be included with the
// step configuration to prevent errors with providers outside the
// registry.terraform.io hostname or outside the hashicorp namespace.
// This is only necessary when using TestStep.Config.
//
// When TestStep.ConfigDirectory is used, the expectation is that the
// Terraform configuration files will specify a terraform configuration
// block and/or provider blocks as necessary.
func (s TestStep) mergedConfig(ctx context.Context, testCase TestCase, configHasTerraformBlock, configHasProviderBlock bool) string {
	var config strings.Builder

	// Prevent issues with existing configurations containing the terraform
	// configuration block.
	if configHasTerraformBlock {
		config.WriteString(s.Config)

		return config.String()
	}

	if testCase.hasProviders(ctx) {
		config.WriteString(s.providerConfigTestCase(ctx, configHasProviderBlock, testCase))
	} else {
		config.WriteString(s.providerConfig(ctx, configHasProviderBlock))
	}

	config.WriteString(s.Config)

	return config.String()
}

// providerConfig takes the list of providers in a TestStep and returns a
// config with only empty provider blocks. This is useful for Import, where no
// config is provided, but the providers must be defined.
func (s TestStep) providerConfig(_ context.Context, skipProviderBlock bool) string {
	var providerBlocks, requiredProviderBlocks strings.Builder

	for name, externalProvider := range s.ExternalProviders {
		if !skipProviderBlock {
			providerBlocks.WriteString(fmt.Sprintf("provider %q {}\n", name))
		}

		if externalProvider.Source == "" && externalProvider.VersionConstraint == "" {
			continue
		}

		requiredProviderBlocks.WriteString(fmt.Sprintf("    %s = {\n", name))

		if externalProvider.Source != "" {
			requiredProviderBlocks.WriteString(fmt.Sprintf("      source = %q\n", externalProvider.Source))
		}

		if externalProvider.VersionConstraint != "" {
			requiredProviderBlocks.WriteString(fmt.Sprintf("      version = %q\n", externalProvider.VersionConstraint))
		}

		requiredProviderBlocks.WriteString("    }\n")
	}

	for name := range s.ProviderFactories {
		requiredProviderBlocks.WriteString(addTerraformBlockSource(name, s.Config))
	}

	for name := range s.ProtoV5ProviderFactories {
		requiredProviderBlocks.WriteString(addTerraformBlockSource(name, s.Config))
	}

	for name := range s.ProtoV6ProviderFactories {
		requiredProviderBlocks.WriteString(addTerraformBlockSource(name, s.Config))
	}

	if requiredProviderBlocks.Len() > 0 {
		return fmt.Sprintf(`
terraform {
  required_providers {
%[1]s
  }
}

%[2]s
`, strings.TrimSuffix(requiredProviderBlocks.String(), "\n"), providerBlocks.String())
	}

	providerBlocksStr := providerBlocks.String()

	return providerBlocksStr
}

func (s TestStep) providerConfigTestCase(_ context.Context, skipProviderBlock bool, testCase TestCase) string {
	var providerBlocks, requiredProviderBlocks strings.Builder

	providerNames := make(map[string]struct{}, len(testCase.Providers))

	for name := range testCase.Providers {
		providerNames[name] = struct{}{}
	}

	for name := range testCase.ProviderFactories {
		delete(providerNames, name)
	}

	// [BF] The Providers field handling predates the logic being moved to this
	//      method. It's not entirely clear to me at this time why this field
	//      is being used and not the others, but leaving it here just in case
	//      it does have a special purpose that wasn't being unit tested prior.
	for name := range providerNames {
		providerBlocks.WriteString(fmt.Sprintf("provider %q {}\n", name))

		requiredProviderBlocks.WriteString(fmt.Sprintf("    %s = {\n", name))

		requiredProviderBlocks.WriteString("    }\n")
	}

	for name, externalProvider := range testCase.ExternalProviders {
		if !skipProviderBlock {
			providerBlocks.WriteString(fmt.Sprintf("provider %q {}\n", name))
		}

		if externalProvider.Source == "" && externalProvider.VersionConstraint == "" {
			continue
		}

		requiredProviderBlocks.WriteString(fmt.Sprintf("    %s = {\n", name))

		if externalProvider.Source != "" {
			requiredProviderBlocks.WriteString(fmt.Sprintf("      source = %q\n", externalProvider.Source))
		}

		if externalProvider.VersionConstraint != "" {
			requiredProviderBlocks.WriteString(fmt.Sprintf("      version = %q\n", externalProvider.VersionConstraint))
		}

		requiredProviderBlocks.WriteString("    }\n")
	}

	for name := range testCase.ProviderFactories {
		providerFactoryBlocks := addTerraformBlockSource(name, s.Config)

		if len(providerFactoryBlocks) > 0 {
			requiredProviderBlocks.WriteString(providerFactoryBlocks)
		}
	}

	for name := range testCase.ProtoV5ProviderFactories {
		protov5ProviderFactoryBlocks := addTerraformBlockSource(name, s.Config)

		if len(protov5ProviderFactoryBlocks) > 0 {
			requiredProviderBlocks.WriteString(protov5ProviderFactoryBlocks)
		}
	}

	for name := range testCase.ProtoV6ProviderFactories {
		protov6ProviderFactoryBlocks := addTerraformBlockSource(name, s.Config)

		if len(protov6ProviderFactoryBlocks) > 0 {
			requiredProviderBlocks.WriteString(addTerraformBlockSource(name, s.Config))
		}
	}

	if requiredProviderBlocks.Len() > 0 {
		return fmt.Sprintf(`
terraform {
  required_providers {
%[1]s
  }
}

%[2]s
`, strings.TrimSuffix(requiredProviderBlocks.String(), "\n"), providerBlocks.String())
	}

	return providerBlocks.String()
}

func addTerraformBlockSource(name, config string) string {
	var js json.RawMessage

	// Do not process JSON.
	if err := json.Unmarshal([]byte(config), &js); err == nil {
		return ""
	}

	var providerBlocks strings.Builder

	providerBlocks.WriteString(fmt.Sprintf("    %s = {\n", name))
	providerBlocks.WriteString(fmt.Sprintf("      source = %q\n", getProviderAddr(name)))
	providerBlocks.WriteString("    }\n")

	return providerBlocks.String()
}
