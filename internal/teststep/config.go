// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	rawConfigFileName              = "terraform_plugin_test.tf"
	rawConfigFileNameJSON          = rawConfigFileName + ".json"
	testCaseProviderConfigFileName = "test_case_provider_config.tf"
	testStepProviderConfigFileName = "test_step_provider_config.tf"
)

var (
	providerConfigBlockRegex  = regexp.MustCompile(`provider "?[a-zA-Z0-9_-]+"? {`)
	terraformConfigBlockRegex = regexp.MustCompile(`terraform {`)
)

type Config interface {
	HasConfiguration() bool
	HasProviderBlock(context.Context) (bool, error)
	Write(context.Context, string) error
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
func (c configuration) HasProviderBlock(ctx context.Context) (bool, error) {
	switch {
	case c.hasRaw(ctx):
		return providerConfigBlockRegex.MatchString(c.raw), nil
	case c.hasDirectory(ctx):
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

	return false, nil
}

func (c configuration) Write(ctx context.Context, dest string) error {
	switch {
	case c.directory != "":
		err := c.writeDirectory(ctx, dest)

		if err != nil {
			return err
		}
	case c.raw != "":
		err := c.writeRaw(ctx, dest)

		if err != nil {
			return err
		}
	}

	return nil
}

func copyFiles(path string, dstPath string) error {
	infos, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, info := range infos {
		srcPath := filepath.Join(path, info.Name())

		if info.IsDir() {
			continue
		} else {
			err = copyFile(srcPath, dstPath)

			if err != nil {
				return err
			}
		}

	}
	return nil
}

func copyFile(path string, dstPath string) error {
	srcF, err := os.Open(path)

	if err != nil {
		return err
	}

	defer srcF.Close()

	di, err := os.Stat(dstPath)

	if err != nil {
		return err
	}

	if di.IsDir() {
		_, file := filepath.Split(path)
		dstPath = filepath.Join(dstPath, file)
	}

	dstF, err := os.Create(dstPath)

	if err != nil {
		return err
	}

	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}

func filesContains(dir string, find *regexp.Regexp) (bool, error) {
	dirEntries, err := os.ReadDir(dir)

	if err != nil {
		return false, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		path := filepath.Join(dir, dirEntry.Name())

		contains, err := fileContains(path, find)

		if err != nil {
			return false, err
		}

		if contains {
			return true, nil
		}
	}

	return false, nil
}

func fileContains(path string, find *regexp.Regexp) (bool, error) {
	f, err := os.ReadFile(path)

	if err != nil {
		return false, err
	}

	return find.MatchString(string(f)), nil
}

func fileExists(dir, fileName string) (bool, error) {
	infos, err := os.ReadDir(dir)

	if err != nil {
		return false, err
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		} else {
			if fileName == info.Name() {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c configuration) hasDirectory(_ context.Context) bool {
	return c.directory != ""
}

func (c configuration) hasRaw(_ context.Context) bool {
	return c.raw != ""
}

// prepareRaw returns a string that assembles the Terraform configuration from
// the raw field, prefixed by either testCaseProviderConfig or
// testStepProviderConfig when the raw field itself does not contain a
// terraform block.
func (c configuration) prepareRaw(_ context.Context) string {
	var config strings.Builder

	// Prevent issues with existing configurations containing the terraform
	// configuration block.
	if terraformConfigBlockRegex.MatchString(c.raw) {
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

// writeDirectory copies the contents of c.directory to dest.
func (c configuration) writeDirectory(_ context.Context, dest string) error {
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

	// Determine whether any of the files in configDirectory contain terraform block.
	containsTerraformConfig, err := filesContains(configDirectory, terraformConfigBlockRegex)

	if err != nil {
		return err
	}

	// Write contents of testCaseProviderConfig or testStepProviderConfig to dest.
	if !containsTerraformConfig {
		if c.testCaseProviderConfig != "" {
			path := filepath.Join(dest, testCaseProviderConfigFileName)

			configFileExists, err := fileExists(configDirectory, testCaseProviderConfigFileName)

			if err != nil {
				return err
			}

			if configFileExists {
				return fmt.Errorf("%s already exists in %s, ", testCaseProviderConfigFileName, configDirectory)
			}

			err = os.WriteFile(path, []byte(c.testCaseProviderConfig), 0700)

			if err != nil {
				return err
			}
		} else {
			path := filepath.Join(dest, testStepProviderConfigFileName)

			configFileExists, err := fileExists(configDirectory, testStepProviderConfigFileName)

			if err != nil {
				return err
			}

			if configFileExists {
				return fmt.Errorf("%s already exists in %s, ", testStepProviderConfigFileName, configDirectory)
			}

			err = os.WriteFile(path, []byte(c.testStepProviderConfig), 0700)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c configuration) writeRaw(ctx context.Context, dest string) error {
	outFilename := filepath.Join(dest, rawConfigFileName)
	rmFilename := filepath.Join(dest, rawConfigFileNameJSON)

	bCfg := []byte(c.prepareRaw(ctx))

	if json.Valid(bCfg) {
		outFilename, rmFilename = rmFilename, outFilename
	}

	if err := os.Remove(rmFilename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to remove %q: %w", rmFilename, err)
	}

	err := os.WriteFile(outFilename, bCfg, 0700)

	if err != nil {
		return err
	}

	return nil
}
