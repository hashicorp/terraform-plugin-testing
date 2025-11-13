// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hack

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"
)

// T is a replacement for github.com/mitchellh/go-testing-interface.T.
//
// RuntimeT and tshim are two implementations of T. RuntimeT can be used to
// test a test framework without the testing.T side effects of stopping
// goroutine execution. tshim can be used as an adapter around a standard
// testing.T.
type T interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Parallel()
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool

	// Not in mitchellh/go-testing-interface
	Attr(key, value string)
	Chdir(dir string)
	Context() context.Context
	Deadline() (deadline time.Time, ok bool)
	Output() io.Writer
	Setenv(key, value string)
	TempDir() string
}

var _ T = RuntimeT{}

type RuntimeT struct {
	fail bool
	skip bool
}

// Attr satisfies [T] and does nothing.
func (t RuntimeT) Attr(key string, value string) {
}

// Chdir satisfies [T] and does nothing.
func (r RuntimeT) Chdir(dir string) {
}

// Cleanup satisfies [T] and does nothing.
func (r RuntimeT) Cleanup(_ func()) {
}

// Context satisfies [T] and returns context.TODO()
func (r RuntimeT) Context() context.Context {
	return context.TODO()
}

// Deadline satisfies [T] and returns zero-values.
func (r RuntimeT) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Error is equivalent to Log followed by Fail.
func (t RuntimeT) Error(args ...interface{}) {
	t.Log(args)
	t.Fail()
}

// Errorf is equivalent to Logf followed by Fail.
func (t RuntimeT) Errorf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.Fail()
}

// Fail marks the function as having failed but continues execution.
func (t RuntimeT) Fail() {
	t.fail = true
}

// Failed reports whether the function has failed.
func (t RuntimeT) Failed() bool {
	return t.fail
}

// FailNow marks the function as having failed and stops its execution by
// calling panic().
//
// For compatibility, it mimics the string argument from
// mitchellg/go-testing-interface.
func (t RuntimeT) FailNow() {
	panic("testing.T failed, see logs for output")
}

// Fatal is equivalent to Log followed by FailNow.
func (t RuntimeT) Fatal(args ...interface{}) {
	t.Log(args)
	t.FailNow()
}

// Fatalf is equivalent to Logf followed by FailNow.
func (t RuntimeT) Fatalf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.FailNow()
}

// Log formats its arguments using default formatting, analogous to
// fmt.Println,, and records the text to standard output.
func (t RuntimeT) Log(args ...interface{}) {
	fmt.Println(fmt.Sprintf("%v", args))
}

// Logf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text to standard output.
func (t RuntimeT) Logf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Helper satisfies [T] and does nothing.
func (r RuntimeT) Helper() {
}

// Name satisfies [T] and returns an empty string.
func (r RuntimeT) Name() string {
	return ""
}

// Output satisfied [T] and returns io.Discard.
func (r RuntimeT) Output() io.Writer {
	return io.Discard
}

// Parallel satisfies [T] and does nothing.
func (r RuntimeT) Parallel() {
	panic("parallel not implemented") // TODO: Implement
}

// Setenv satisfies [T] and does nothing.
func (r RuntimeT) Setenv(key string, value string) {
}

// SkipNow marks the test as having been skipped.
//
// As a practical consideration, this does not stop execution in the way that
// [testing.T.SkipNow] does -- RuntimeT.Run does not run its function in a
// separate goroutine.
func (t RuntimeT) SkipNow() {
	t.Skip()
}

// Skipf is equivalent to Logf followed by SkipNow.
func (t RuntimeT) Skipf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.Skip()
}

// TempDir satisfies [T] and returns "/dev/null".
func (r RuntimeT) TempDir() string {
	return "/dev/null"
}

// Skip is equivalent to Log followed by SkipNow.
func (t RuntimeT) Skip(args ...interface{}) {
	t.Log(args)
	t.skip = true
}

// Skipped reports whether the test was skipped.
func (t RuntimeT) Skipped() bool {
	return t.skip
}

// tshim embeds a [testing.T] and satisfies [T].
type tshim struct {
	*testing.T
}

// Run runs f as a subtest of t.
//
// As practical consideration, this does nsot run its function in a separate
// goroutine and the name is not used.
func (t tshim) Run(name string, f func(T)) bool {
	return t.T.Run(name, func(t *testing.T) {
		f(tshim{t})
	})
}
