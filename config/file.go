// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path/filepath"
	"strconv"
)

// StaticFile is a helper function that returns the supplied
// file when TestStepConfigFunc is executed.
func StaticFile(file string) func(TestStepConfigRequest) string {
	return func(_ TestStepConfigRequest) string {
		return file
	}
}

// TestNameFile returns the name of the test suffixed with the supplied
// file name when TestStepConfigFunc is executed (e.g., "testdata/TestExampleCloudThing_basic/test.tf.
//
// For example, given test code:
//
//	func TestExampleCloudThing_basic(t *testing.T) {
//	    resource.Test(t, resource.TestCase{
//	        Steps: []resource.TestStep{
//	            {
//	                ConfigFile: config.TestNameFile("test.tf"),
//	            },
//	        },
//	    })
//	}
//
// The testing configuration will be expected in the
// testdata/TestExampleCloudThing_basic/test.tf file.
func TestNameFile(file string) func(TestStepConfigRequest) string {
	return func(req TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, file)
	}
}

// TestStepFile returns the name of the test suffixed
// with the test step number and the supplied file name.
func TestStepFile(file string) func(TestStepConfigRequest) string { //nolint:paralleltest //Not a test
	return func(req TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, strconv.Itoa(req.StepNumber), file)
	}
}
