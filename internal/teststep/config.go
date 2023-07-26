// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

const (
	rawConfigFileName     = "terraform_plugin_test.tf"
	rawConfigFileNameJSON = rawConfigFileName + ".json"
)

var (
	providerConfigBlockRegex  = regexp.MustCompile(`provider "?[a-zA-Z0-9_-]+"? {`)
	terraformConfigBlockRegex = regexp.MustCompile(`terraform {`)
)

type Config interface {
	HasConfiguration() bool
	HasConfigurationFiles() bool
	HasProviderBlock(context.Context) (bool, error)
	HasTerraformBlock(context.Context) (bool, error)
	Write(context.Context, string) error
}

type configuration struct {
	directory string
	raw       string
}

type ConfigurationRequest struct {
	Directory *string
	Raw       *string
}

func (c ConfigurationRequest) Validate() error {
	if c.Directory != nil && c.Raw != nil && *c.Directory != "" && *c.Raw != "" {
		return errors.New(`both "directory" and "raw" are populated, only one configuration option is allowed`)
	}

	return nil
}

func Configuration(req ConfigurationRequest) (configuration, error) {
	var config configuration

	if req.Directory != nil {
		config.directory = *req.Directory
	}

	if req.Raw != nil {
		config.raw = *req.Raw
	}

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

func (c configuration) HasConfigurationFiles() bool {
	return c.directory != ""
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

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
func (c configuration) HasTerraformBlock(ctx context.Context) (bool, error) {
	switch {
	case c.hasRaw(ctx):
		return terraformConfigBlockRegex.MatchString(c.raw), nil
	case c.hasDirectory(ctx):
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

func (c configuration) hasDirectory(_ context.Context) bool {
	return c.directory != ""
}

func (c configuration) hasRaw(_ context.Context) bool {
	return c.raw != ""
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

	return nil
}

func (c configuration) writeRaw(_ context.Context, dest string) error {
	outFilename := filepath.Join(dest, rawConfigFileName)
	rmFilename := filepath.Join(dest, rawConfigFileNameJSON)

	bCfg := []byte(c.raw)

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

func Pointer[T any](in T) *T {
	return &in
}
