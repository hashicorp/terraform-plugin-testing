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

	"github.com/hashicorp/terraform-plugin-testing/config"
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
	HasConfigurationFiles() bool
	HasProviderBlock(context.Context) (bool, error)
	HasTerraformBlock(context.Context) (bool, error)
	Write(context.Context, string) error
}

type PrepareConfigurationRequest struct {
	Directory             config.TestStepConfigFunc
	File                  config.TestStepConfigFunc
	Raw                   string
	TestStepConfigRequest config.TestStepConfigRequest
}

func (p PrepareConfigurationRequest) Exec() ConfigurationRequest {
	directory := Pointer(p.Directory.Exec(p.TestStepConfigRequest))
	file := Pointer(p.File.Exec(p.TestStepConfigRequest))
	raw := Pointer(p.Raw)

	return ConfigurationRequest{
		Directory: directory,
		File:      file,
		Raw:       raw,
	}
}

type ConfigurationRequest struct {
	Directory *string
	File      *string
	Raw       *string
}

func (c ConfigurationRequest) Validate() error {
	var configSet []string

	if c.Directory != nil && *c.Directory != "" {
		configSet = append(configSet, "directory")
	}

	if c.File != nil && *c.File != "" {
		configSet = append(configSet, "file")
	}

	if c.Raw != nil && *c.Raw != "" {
		configSet = append(configSet, "raw")
	}

	if len(configSet) > 1 {
		configSetStr := strings.Join(configSet, `, `)

		i := strings.LastIndex(configSetStr, ", ")

		if i != -1 {
			configSetStr = configSetStr[:i] + " and " + configSetStr[i+len(", "):]
		}

		return fmt.Errorf(`%s are populated, only one of "directory", "file", or "raw"  is allowed`, configSetStr)
	}

	return nil
}

func Configuration(req ConfigurationRequest) (Config, error) {
	if req.Directory != nil && *req.Directory != "" {
		return configurationDirectory{
			directory: *req.Directory,
		}, nil
	}

	if req.File != nil && *req.File != "" {
		return configurationFile{
			file: *req.File,
		}, nil
	}

	if req.Raw != nil && *req.Raw != "" {
		return configurationString{
			raw: *req.Raw,
		}, nil
	}

	return nil, nil
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

func Pointer[T any](in T) *T {
	return &in
}
