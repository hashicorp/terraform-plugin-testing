// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path/filepath"
	"strconv"
)

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
//
// For example, given test code:
//
//	func TestExampleCloudThing_basic(t *testing.T) {
//	    resource.Test(t, resource.TestCase{
//	        Steps: []resource.TestStep{
//	            {
//	                ConfigDirectory: config.TestNameDirectory(),
//	            },
//	        },
//	    })
//	}
//
// The testing configurations will be expected in the
// testdata/TestExampleCloudThing_basic/ directory.
func TestNameDirectory() func(TestStepConfigRequest) string {
	return func(req TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName)
	}
}

// TestStepDirectory returns the name of the test suffixed
// with the test step number.
func TestStepDirectory() func(TestStepConfigRequest) string { //nolint:paralleltest //Not a test
	return func(req TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, strconv.Itoa(req.StepNumber))
	}
}
