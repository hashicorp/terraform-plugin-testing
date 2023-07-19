// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var configProviderBlockRegex = regexp.MustCompile(`provider "?[a-zA-Z0-9_-]+"? {`)

type Config interface {
	HasConfiguration() bool
	HasDirectory() bool
	HasProviderBlock(context.Context) bool
	HasRaw(context.Context) bool
	Raw(context.Context) string
	WriteDirectory(context.Context, string) error
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

func (c configuration) HasDirectory() bool {
	return c.directory != ""
}

// HasProviderBlock returns true if the Config has declared a provider
// configuration block, e.g. provider "examplecloud" {...}
//
// TODO: Need to handle configuration supplied through Directory or File.
func (c configuration) HasProviderBlock(_ context.Context) bool {
	return configProviderBlockRegex.MatchString(c.raw)
}

func (c configuration) HasRaw(ctx context.Context) bool {
	return c.raw != ""
}

// Raw returns a string that assembles the Terraform configuration from
// the raw field, prefixed by either testCaseProviderConfig or
// testStepProviderConfig when the raw field itself does not contain a
// terraform block.
// TODO: Consider whether Raw and Directory should be dropped in favour
// of having configuration manage the writing of either raw configuration
// or directory files, respectively. This would allow for easier
// management of testCaseProviderConfig and testStepProviderConfig when
// handling the copying of files from the configuration.directory.
func (c configuration) Raw(ctx context.Context) string {
	var config strings.Builder

	// Prevent issues with existing configurations containing the terraform
	// configuration block.
	if strings.Contains(c.raw, "terraform {") {
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

// WriteDirectory copies the contents of c.directory to dest.
// TODO: Need to handle testCaseProviderConfig and testStepProviderConfig when
// copying files from c.directory to working directory for test
func (c configuration) WriteDirectory(ctx context.Context, dest string) error {
	// Copy all files from c.dir to dest
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
	containsTerraformConfig, err := c.filesContains(configDirectory, `terraform {`)

	if err != nil {
		return err
	}

	// Write contents of testCaseProviderConfig or testStepProviderConfig t dest.
	// TODO: Verify whether there are any naming collisions between the files in
	// c.directory and the files that are written below.
	if !containsTerraformConfig {
		if c.testCaseProviderConfig != "" {
			path := filepath.Join(dest, "testCaseProviderConfig.tf")

			err := os.WriteFile(path, []byte(c.testCaseProviderConfig), 0700)

			if err != nil {
				return err
			}
		} else {
			path := filepath.Join(dest, "testStepProviderConfig.tf")

			err := os.WriteFile(path, []byte(c.testStepProviderConfig), 0700)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c configuration) filesContains(dir, find string) (bool, error) {
	dirEntries, err := os.ReadDir(dir)

	if err != nil {
		return false, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		path := filepath.Join(dir, dirEntry.Name())

		contains, err := c.fileContains(path, find)

		if err != nil {
			return false, err
		}

		if contains {
			return true, nil
		}
	}

	return false, nil
}

func (c configuration) fileContains(path, find string) (bool, error) {
	f, err := os.ReadFile(path)

	if err != nil {
		return false, err
	}

	return strings.Contains(string(f), find), nil
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
