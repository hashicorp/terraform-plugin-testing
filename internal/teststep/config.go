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
	GetRaw(context.Context) string
	HasConfiguration() bool
	HasProviderBlock(context.Context) bool
}

type configuration struct {
	directory              string
	raw                    string
	testCaseProviderConfig string
	testStepProviderConfig string
}

type ConfigurationRequest struct {
	Directory              string
	Raw                    string
	TestCaseProviderConfig string
	TestStepProviderConfig string
}

func Configuration(req ConfigurationRequest) (configuration, error) {
	var populatedConfig []string
	var config configuration

	if req.Directory != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "directory"))

		config = configuration{
			directory: req.Directory,
		}
	}

	if req.Raw != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "raw"))

		config = configuration{
			raw: req.Raw,
		}
	}

	if len(populatedConfig) > 1 {
		return configuration{}, fmt.Errorf(
			"both %s are populated, only one configuration option is allowed",
			strings.Join(populatedConfig, " and "),
		)
	}

	config.testCaseProviderConfig = req.TestCaseProviderConfig
	config.testStepProviderConfig = req.TestStepProviderConfig

	return config, nil
}

func (c configuration) GetRaw(ctx context.Context) string {
	var config strings.Builder

	// Prevent issues with existing configurations containing the terraform
	// configuration block.
	if c.hasTerraformBlock(ctx) {
		return c.raw
	}

	if c.testCaseProviderConfig != "" {
		config.WriteString(c.testCaseProviderConfig)
	} else {
		config.WriteString(c.testStepProviderConfig)
	}

	config.WriteString(c.raw)

	return config.String()
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

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
//
// TODO: Need to handle configuration supplied through Directory or File.
func (c configuration) hasTerraformBlock(_ context.Context) bool {
	return strings.Contains(c.raw, "terraform {")
}
