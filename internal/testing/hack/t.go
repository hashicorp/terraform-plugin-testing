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

// BaseT is a replacement for github.com/mitchellh/go-testing-interface.T.
//
// RuntimeT and StandardT are two implementations of BaseT. MetaT can be used
// to test a testing.T-based test framework without the side effects of
// stopping goroutine execution. StandardT can be used as an adapter around a
// standard testing.T.
//
// Precedent for StandardT: the unexported tshim type in
// github.com/rogpeppe/go-internal/testscript.
type BaseT interface {
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

var _ BaseT = MetaT{}

type MetaT struct {
	fail bool
	skip bool
}

// Attr satisfies [BaseT] and does nothing.
func (t MetaT) Attr(key string, value string) {
}

// Chdir satisfies [BaseT] and does nothing.
func (r MetaT) Chdir(dir string) {
}

// Cleanup satisfies [BaseT] and does nothing.
func (r MetaT) Cleanup(_ func()) {
}

// Context satisfies [BaseT] and returns context.TODO()
func (r MetaT) Context() context.Context {
	return context.TODO()
}

// Deadline satisfies [BaseT] and returns zero-values.
func (r MetaT) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Error is equivalent to Log followed by Fail.
func (t MetaT) Error(args ...interface{}) {
	t.Log(args)
	t.Fail()
}

// Errorf is equivalent to Logf followed by Fail.
func (t MetaT) Errorf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.Fail()
}

// Fail marks the function as having failed but continues execution.
func (t MetaT) Fail() {
	t.fail = true
}

// Failed reports whether the function has failed.
func (t MetaT) Failed() bool {
	return t.fail
}

// FailNow marks the function as having failed and stops its execution by
// calling panic().
//
// For compatibility, it mimics the string argument from
// mitchellg/go-testing-interface.
func (t MetaT) FailNow() {
	panic("testing.T failed, see logs for output")
}

// Fatal is equivalent to Log followed by FailNow.
func (t MetaT) Fatal(args ...interface{}) {
	t.Log(args)
	t.FailNow()
}

// Fatalf is equivalent to Logf followed by FailNow.
func (t MetaT) Fatalf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.FailNow()
}

// Log formats its arguments using default formatting, analogous to
// fmt.Println,, and records the text to standard output.
func (t MetaT) Log(args ...interface{}) {
	fmt.Println(fmt.Sprintf("%v", args))
}

// Logf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text to standard output.
func (t MetaT) Logf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Helper satisfies [BaseT] and does nothing.
func (r MetaT) Helper() {
}

// Name satisfies [BaseT] and returns an empty string.
func (r MetaT) Name() string {
	return ""
}

// Output satisfied [BaseT] and returns io.Discard.
func (r MetaT) Output() io.Writer {
	return io.Discard
}

// Parallel satisfies [BaseT] and does nothing.
func (r MetaT) Parallel() {
	panic("parallel not implemented") // TODO: Implement
}

// Setenv satisfies [BaseT] and does nothing.
func (r MetaT) Setenv(key string, value string) {
}

// SkipNow marks the test as having been skipped.
//
// As a practical consideration, this does not stop execution in the way that
// [testing.T.SkipNow] does -- RuntimeT.Run does not run its function in a
// separate goroutine.
func (t MetaT) SkipNow() {
	t.Skip()
}

// Skipf is equivalent to Logf followed by SkipNow.
func (t MetaT) Skipf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.Skip()
}

// TempDir satisfies [BaseT] and returns "/dev/null".
func (r MetaT) TempDir() string {
	return "/dev/null"
}

// Skip is equivalent to Log followed by SkipNow.
func (t MetaT) Skip(args ...interface{}) {
	t.Log(args)
	t.skip = true
}

// Skipped reports whether the test was skipped.
func (t MetaT) Skipped() bool {
	return t.skip
}

var _ BaseT = StandardT{}

// StandardT embeds a [testing.T] and satisfies [BaseT].
type StandardT struct {
	*testing.T
}

// Run runs f as a subtest of t.
//
// As practical consideration, this does nsot run its function in a separate
// goroutine and the name is not used.
func (t StandardT) Run(name string, f func(BaseT)) bool {
	return t.T.Run(name, func(t *testing.T) {
		f(StandardT{t})
	})
}
