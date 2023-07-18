// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var configProviderBlockRegex = regexp.MustCompile(`provider "?[a-zA-Z0-9_-]+"? {`)

type Config interface {
	GetRaw() string
	HasConfiguration() bool
	HasProviderBlock(context.Context) bool
	MergedConfig(context.Context, string, string) configuration
}

type configuration struct {
	directory string
	raw       string
}

type ConfigurationRequest struct {
	Directory string
	Raw       string
}

func Configuration(configRequest ConfigurationRequest) (configuration, error) {
	var populatedConfig []string
	var config configuration

	if configRequest.Directory != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "directory"))

		config = configuration{
			directory: configRequest.Directory,
		}
	}

	if configRequest.Raw != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "raw"))

		config = configuration{
			raw: configRequest.Raw,
		}
	}

	if len(populatedConfig) > 1 {
		return configuration{}, fmt.Errorf(
			"both %s are populated, only one configuration option is allowed",
			strings.Join(populatedConfig, " and "),
		)
	}

	return config, nil
}

func (c configuration) GetRaw() string {
	return c.raw
}

func (c configuration) HasConfiguration() bool {
	if c.directory != "" {
		return true
	}

	if c.raw != "" {
		return true
	}

	return false
}

// HasProviderBlock returns true if the Config has declared a provider
// configuration block, e.g. provider "examplecloud" {...}
//
// TODO: Need to handle configuration supplied through Directory or File.
func (c configuration) HasProviderBlock(_ context.Context) bool {
	return configProviderBlockRegex.MatchString(c.raw)
}

// MergedConfig prepends any necessary terraform configuration blocks to the
// TestStep Config.
//
// If there are ExternalProviders configurations in either the TestCase or
// TestStep, the terraform configuration block should be included with the
// step configuration to prevent errors with providers outside the
// registry.terraform.io hostname or outside the hashicorp namespace.
//
// TODO: Need to handle configuration supplied through Directory or File.
func (c configuration) MergedConfig(ctx context.Context, testCaseProviderConfig, testStepProviderConfig string) configuration {
	var config strings.Builder

	// Prevent issues with existing configurations containing the terraform
	// configuration block.
	if c.hasTerraformBlock(ctx) {
		return configuration{
			raw: c.raw,
		}
	}

	if testCaseProviderConfig != "" {
		config.WriteString(testCaseProviderConfig)
	} else {
		config.WriteString(testStepProviderConfig)
	}

	config.WriteString(c.raw)

	return configuration{
		raw: config.String(),
	}
}

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
//
// TODO: Need to handle configuration supplied through Directory or File.
func (c configuration) hasTerraformBlock(_ context.Context) bool {
	return strings.Contains(c.raw, "terraform {")
}
