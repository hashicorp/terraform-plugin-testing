// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testingiface

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var _ T = (*MockT)(nil)

// MockT is a lightweight, mock implementation of the non-extensible Go
// standard library [*testing.T] implementation. This type should only be used
// in the unit testing of functionality within this Go module and never exposed
// in the Go module API.
//
// This type intentionally is not feature-complete and only includes the methods
// necessary for the existing code in this Go module.
type MockT struct {
	// TestName can be set to the name of the test, such as calling the real
	// [testing.T.Name]. This is returned in the [Name] method.
	TestName string

	isHelper    bool
	isFailed    bool
	isSkipped   bool
	isParallel  bool
	lastError   string
	lastSkipped string
}

// T interface implementations

func (t *MockT) Cleanup(func()) {
	panic("not implemented")
}

func (t *MockT) Deadline() (deadline time.Time, ok bool) {
	panic("not implemented")
}

func (t *MockT) Error(args ...any) {
	t.lastError = fmt.Sprintln(args...)
	t.Log(args...)
	t.Fail()
}

func (t *MockT) Errorf(format string, args ...any) {
	t.lastError = fmt.Sprintf(format, args...)
	t.Logf(format, args...)
	t.Fail()
}

func (t *MockT) Fail() {
	t.isFailed = true
}

func (t *MockT) Failed() bool {
	return t.isFailed
}

func (t *MockT) FailNow() {
	t.Fail()
	runtime.Goexit()
}

func (t *MockT) Fatal(args ...any) {
	t.lastError = fmt.Sprintln(args...)
	t.Log(args...)
	t.FailNow()
}

func (t *MockT) Fatalf(format string, args ...any) {
	t.lastError = fmt.Sprintf(format, args...)
	t.Log(args...)
	t.FailNow()
}

func (t *MockT) Helper() {
	t.isHelper = true
}

func (t *MockT) Log(args ...any) {
	if args == nil {
		return
	}

	fmt.Fprintln(os.Stdout, args...)
}

func (t *MockT) Logf(format string, args ...any) {
	if format == "" {
		return
	}

	fmt.Fprintf(os.Stdout, format, args...)
}

func (t *MockT) Name() string {
	return t.TestName
}

func (t *MockT) Parallel() {
	t.isParallel = true
}

func (t *MockT) Run(name string, f func(t T)) bool {
	t.Log("Running subtest:", name)
	defer t.Log("Finished subtest:", name)

	f(t)
	return !t.isFailed
}

func (t *MockT) Setenv(key string, value string) {
	panic("not implemented")
}

func (t *MockT) Skip(args ...any) {
	t.lastSkipped = fmt.Sprintln(args...)
	t.Log(args...)
	t.SkipNow()
}

func (t *MockT) Skipf(format string, args ...any) {
	t.lastSkipped = fmt.Sprintf(format, args...)
	t.Logf(format, args...)
	t.SkipNow()
}

func (t *MockT) SkipNow() {
	t.isSkipped = true
	runtime.Goexit()
}

func (t *MockT) Skipped() bool {
	return t.isSkipped
}

func (t *MockT) TempDir() string {
	panic("not implemented")
}

// Custom methods

func (t *MockT) IsParallel() bool {
	return t.isParallel
}

func (t *MockT) LastError() string {
	return t.lastError
}

func (t *MockT) LastSkipped() string {
	return t.lastSkipped
}
