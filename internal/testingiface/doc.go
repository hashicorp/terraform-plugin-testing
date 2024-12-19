// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package testingiface provides wrapper types compatible with the Go standard
// library [testing] package. These wrappers are necessary for implementing
// [testing] package helpers, since the Go standard library implementation is
// not extensible and the existing code in this Go module is built on directly
// interacting with [testing] functionality, such as calling [testing.T.Fatal].
//
// The [T] interface has all methods of the [testing.T] type and the [MockT]
// type is a lightweight mock implementation of the [T] interface. There are a
// collection of assertion helper functions such as:
//   - [ExpectFail]: That the test logic called the equivalent of
//     [testing.T.Error] or [testing.T.Fatal].
//   - [ExpectParallel]: That the test logic called the equivalent of
//     [testing.T.Parallel] and passed.
//   - [ExpectPass]: That the test logic did not call the equivalent of
//     [testing.T.Skip], since [testing] marks these tests as passing.
//   - [ExpectSkip]: That the test logic called the equivalent of
//     [testing.T.Skip].
//
// This code in this package is intentionally internal and should not be exposed
// in the Go module API. It is compatible with the Go 1.17 [testing] package.
// It replaces the archived github.com/mitchellh/go-testing-interface Go module,
// but is implemented with different approaches that enable calls to behave
// more closely to the Go standard library, such as calling [runtime.Goexit]
// when skipping, and preserving any error/skip messaging.
package testingiface
