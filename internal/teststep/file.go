// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
)

var _ Config = configurationFile{}

type configurationFile struct {
	file string
}

func (c configurationFile) HasConfigurationFiles() bool {
	return true
}

// HasProviderBlock returns true if the Config has declared a provider
// configuration block, e.g. provider "examplecloud" {...}
func (c configurationFile) HasProviderBlock(ctx context.Context) (bool, error) {
	pwd, err := os.Getwd()

	if err != nil {
		return false, err
	}

	configFile := filepath.Join(pwd, c.file)

	contains, err := fileContains(configFile, providerConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
func (c configurationFile) HasTerraformBlock(ctx context.Context) (bool, error) {
	pwd, err := os.Getwd()

	if err != nil {
		return false, err
	}

	configFile := filepath.Join(pwd, c.file)

	contains, err := fileContains(configFile, terraformConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

func (c configurationFile) Write(ctx context.Context, dest string) error {
	// Copy file from c.file to dest
	pwd, err := os.Getwd()

	if err != nil {
		return err
	}

	configFile := filepath.Join(pwd, c.file)

	err = copyFile(configFile, dest)

	if err != nil {
		return err
	}

	return nil
}
