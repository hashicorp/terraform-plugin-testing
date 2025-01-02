// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testingiface

// T is the interface that contains all the methods of the Go standard library
// [*testing.T] type as of Go 1.17.
//
// For complete backwards compatibility, it explicitly does not include the
// Deadline and Run methods to match the prior
// github.com/mitchellh/go-testing-interface.T interface. If either of those
// methods are needed, it should be relatively safe to add to this interface
// under the guise that this internal interface should match the [*testing.T]
// implementation.
type T interface {
	Cleanup(func())

	// Excluded to match the prior github.com/mitchellh/go-testing-interface.T
	// interface for complete backwards compatibility. It is relatively safe to
	// introduce if necessary in the future though.
	// Deadline() (deadline time.Time, ok bool)

	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	Failed() bool
	FailNow()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Parallel()

	// Excluded to match the prior github.com/mitchellh/go-testing-interface.T
	// interface for complete backwards compatibility. It is relatively safe to
	// introduce if necessary in the future though.
	Run(name string, f func(T)) bool

	Setenv(key string, value string)
	Skip(args ...any)
	Skipf(format string, args ...any)
	SkipNow()
	Skipped() bool
	TempDir() string
}
