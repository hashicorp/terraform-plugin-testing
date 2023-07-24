// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path/filepath"
	"strconv"
	"testing"
)

type TestStepConfigFunc func(TestStepConfigRequest) string

type TestStepConfigRequest struct {
	StepNumber int
}

func (f TestStepConfigFunc) Exec(req TestStepConfigRequest) string {
	if f != nil {
		return f(req)
	}

	return ""
}

func StaticDirectory(directory string) func(TestStepConfigRequest) string {
	return func(_ TestStepConfigRequest) string {
		return directory
	}
}

//nolint:paralleltest //Not a test
func TestNameDirectory(t *testing.T) func(TestStepConfigRequest) string {
	return func(_ TestStepConfigRequest) string {
		return t.Name()
	}
}

//nolint:paralleltest //Not a test
func TestStepDirectory(t *testing.T) func(TestStepConfigRequest) string {
	return func(req TestStepConfigRequest) string {
		return filepath.Join(t.Name(), strconv.Itoa(req.StepNumber))
	}
}
