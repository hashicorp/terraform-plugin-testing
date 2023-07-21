package config

import (
	"path/filepath"
	"strconv"
	"testing"
)

type TestStepConfigRequest struct {
	StepNumber int
}

type TestStepConfigFunc = func(TestStepConfigRequest) string

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

func ExecuteTestStepConfigFunc(f TestStepConfigFunc, r TestStepConfigRequest) string {
	if f != nil {
		return f(r)
	}

	return ""
}
