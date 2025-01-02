// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testingiface

import (
	"sync"
	"testing"
)

// ExpectFailed provides a wrapper for test logic which should call any of the
// following:
//   - [testing.T.Error]
//   - [testing.T.Errorf]
//   - [testing.T.Fatal]
//   - [testing.T.Fatalf]
//
// If none of those were called, the real [testing.T.Fatal] is called to fail
// the test.
func ExpectFail(t *testing.T, logic func(*MockT)) {
	t.Helper()

	mockT := &MockT{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		logic(mockT)
	}()

	wg.Wait()

	if mockT.Failed() {
		return
	}

	t.Fatal("expected test failure")
}

// ExpectParallel provides a wrapper for test logic which should call the
// [testing.T.Parallel] method. If it doesn't, the real [testing.T.Fatal] is
// called.
func ExpectParallel(t *testing.T, logic func(*MockT)) {
	t.Helper()

	mockT := &MockT{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		logic(mockT)
	}()

	wg.Wait()

	if mockT.Failed() {
		t.Fatalf("unexpected test failure: %s", mockT.LastError())
	}

	if mockT.Skipped() {
		t.Fatalf("unexpected test skip: %s", mockT.LastSkipped())
	}

	if mockT.IsParallel() {
		return
	}

	t.Fatal("expected test parallel")
}

// ExpectPass provides a wrapper for test logic which should not call any of the
// following, which would mark the real test as passing:
//   - [testing.T.Skip]
//   - [testing.T.Skipf]
//
// If one of those were called, the real [testing.T.Fatal] is called to fail
// the test. This is only necessary to check for false positives with skipped
// tests.
func ExpectPass(t *testing.T, logic func(*MockT)) {
	t.Helper()

	mockT := &MockT{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		logic(mockT)
	}()

	wg.Wait()

	if mockT.Failed() {
		t.Fatalf("unexpected test failure: %s", mockT.LastError())
	}

	if mockT.Skipped() {
		t.Fatalf("unexpected test skip: %s", mockT.LastSkipped())
	}

	// test passed as expected
}

// ExpectSkip provides a wrapper for test logic which should call any of the
// following:
//   - [testing.T.Skip]
//   - [testing.T.Skipf]
//
// If none of those were called, the real [testing.T.Fatal] is called to fail
// the test.
func ExpectSkip(t *testing.T, logic func(*MockT)) {
	t.Helper()

	mockT := &MockT{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		logic(mockT)
	}()

	wg.Wait()

	if mockT.Failed() {
		t.Fatalf("unexpected test failure: %s", mockT.LastError())
	}

	if mockT.Skipped() {
		return
	}

	t.Fatal("test passed, expected test skip")
}
