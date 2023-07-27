// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
)

var _ Config = configurationDirectory{}

type configurationDirectory struct {
	directory string
}

func (c configurationDirectory) HasConfigurationFiles() bool {
	return true
}

// HasProviderBlock returns true if the Config has declared a provider
// configuration block, e.g. provider "examplecloud" {...}
func (c configurationDirectory) HasProviderBlock(ctx context.Context) (bool, error) {
	pwd, err := os.Getwd()

	if err != nil {
		return false, err
	}

	configDirectory := filepath.Join(pwd, c.directory)

	contains, err := filesContains(configDirectory, providerConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
func (c configurationDirectory) HasTerraformBlock(ctx context.Context) (bool, error) {
	pwd, err := os.Getwd()

	if err != nil {
		return false, err
	}

	configDirectory := filepath.Join(pwd, c.directory)

	contains, err := filesContains(configDirectory, terraformConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

func (c configurationDirectory) Write(ctx context.Context, dest string) error {
	// Copy all files from c.directory to dest
	pwd, err := os.Getwd()

	if err != nil {
		return err
	}

	configDirectory := filepath.Join(pwd, c.directory)

	err = copyFiles(configDirectory, dest)

	if err != nil {
		return err
	}

	return nil
}
