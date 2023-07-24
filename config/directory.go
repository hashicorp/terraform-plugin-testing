// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path/filepath"
	"strconv"
	"testing"
)

// TestStepConfigFunc is the callback type used with acceptance tests to
// specify a string which identifies a directory containing Terraform
// configuration files.
type TestStepConfigFunc func(TestStepConfigRequest) string

// TestStepConfigRequest defines the request supplied to types
// implementing TestStepConfigFunc.
type TestStepConfigRequest struct {
	StepNumber int
}

// Exec executes TestStepConfigFunc if it is not nil, otherwise an
// empty string is returned.
func (f TestStepConfigFunc) Exec(req TestStepConfigRequest) string {
	if f != nil {
		return f(req)
	}

	return ""
}

// StaticDirectory is a helper function that returns the supplied
// directory when TestStepConfigFunc is executed.
func StaticDirectory(directory string) func(TestStepConfigRequest) string {
	return func(_ TestStepConfigRequest) string {
		return directory
	}
}

// TestNameDirectory returns the name of the test when TestStepConfigFunc
// is executed. This facilitates a convention of naming directories
// containing Terraform configuration files with the name of the test.
func TestNameDirectory(t *testing.T) func(TestStepConfigRequest) string { //nolint:paralleltest //Not a test
	return func(_ TestStepConfigRequest) string {
		return t.Name()
	}
}

// TestStepDirectory returns the name of the test suffixed with an
// OS specific separator and the test step number. This facilitates
// a convention of naming directories containing Terraform
// configuration files with the test step number and nesting of
// these files within a directory with the same name as the test.
func TestStepDirectory(t *testing.T) func(TestStepConfigRequest) string { //nolint:paralleltest //Not a test
	return func(req TestStepConfigRequest) string {
		return filepath.Join(t.Name(), strconv.Itoa(req.StepNumber))
	}
}
