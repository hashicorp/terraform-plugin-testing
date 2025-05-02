// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !testingtesting

package tfversion_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// buildTag filters tests to only those that this runner can run.
const buildTag = "testingtestingtesting"

// testEvent follows the test2json format reference:
// https://pkg.go.dev/cmd/test2json
type testEvent struct {
	// Action is one of a fixed set of action descriptions, including:
	// output, pass, skip, fail
	Action string

	// FailedBuild is set for Action == "fail" if the test failure was caused by a build failure.
	FailedBuild string

	// Output is set for Action == "output" and is a portion of the test's
	// output (standard output and standard error merged together).
	Output string

	// Test, if present, specifies the test, example, or benchmark function that caused the event.
	Test string
}

// testResult is this test runner's model of a test result
type testResult struct {
	Name string

	// ActualOutcome can have values: pass, fail, skip
	ActualOutcome string

	// ExpectedOutcome can have values: pass, fail, skip
	ExpectedOutcome string

	Output *strings.Builder
}

func (tr testResult) Green() bool {
	return tr.ActualOutcome == tr.ExpectedOutcome
}

func (tr testResult) Red() bool {
	return tr.ActualOutcome != tr.ExpectedOutcome
}

func (tr testResult) Panic() bool {
	return strings.HasPrefix(tr.Output.String(), "panic:")
}

func (tr testResult) Timeout() bool {
	return strings.HasPrefix(tr.Output.String(), "panic: test timed out")
}

// listTests reads the names and expected outcomes of tests to be run. It runs `go test -json -list -tags=[buildTag]`.
func listTests() (map[string]*testResult, error) {
	cmd := exec.Command("go", "test", "-json", "-list", ".", "-tags="+buildTag)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(output)
	scanner := bufio.NewScanner(reader)

	buffer := make([]byte, 1024*1024) // An arbitrary buffer size
	scanner.Buffer(buffer, 1024*1024) // An arbitrary max token size

	tests := make(map[string]*testResult, 10)
	for scanner.Scan() {
		event := testEvent{}
		err := json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		event.Output = strings.TrimSuffix(event.Output, "\n")

		if event.Action != "output" {
			continue
		}

		if strings.HasPrefix(event.Output, "Test") {
			var expectedOutcome string
			switch {
			case strings.Contains(event.Output, "ExpectFail"):
				expectedOutcome = "fail"

			case strings.Contains(event.Output, "ExpectSkip"):
				expectedOutcome = "skip"

			default:
				expectedOutcome = "pass"
			}

			tests[event.Output] = &testResult{Name: event.Output, ExpectedOutcome: expectedOutcome, Output: &strings.Builder{}}
		}
	}

	return tests, nil
}

// Run test files that have the buildTag that we're looking for.  Reads the
// output of the `go test -json` command. Reads a test name suffix to infer,
// for each test, its expected outcome: pass, fail, or skip. Compares actual
// outcome to expected outcome. Fails if any actual outcome does not match the
// corresponding expected outcome.
func TestTest(t *testing.T) {
	t.Parallel()

	testResults, err := listTests()
	if err != nil {
		t.Fatalf("failed to list tests: %v", err)
	}

	var output []byte
	cmd := exec.Command("go", "test", "-json", "-tags="+buildTag) // TODO: timeout?
	output, err = cmd.Output()

	switch err := err.(type) {
	default:
		t.Fatalf("TestTest failed to test: %v %s", err, string(output))

	// ExitError means that `go test` exited with an error. This includes
	// failing tests and failing builds. `nil` means that `go test` exited
	// without an error, which means that any tests that were found have
	// passed.
	case *exec.ExitError, nil:
		reader := bytes.NewReader(output)
		scanner := bufio.NewScanner(reader)

		buffer := make([]byte, 1024)
		scanner.Buffer(buffer, 1*1024*1024)

		event := testEvent{}
		for scanner.Scan() {
			if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
				t.Fatalf("failed to unmarshal JSON: %v", err)
			}

			switch event.Action {
			case "pass", "fail", "skip":
				testResults[event.Test].ActualOutcome = event.Action

			case "output":
				testResults[event.Test].Output.WriteString(event.Output)
			}

			if len(event.FailedBuild) > 0 {
				t.Fatalf("failed build: %s", event.FailedBuild)
			}
		}

		for _, testResult := range testResults {
			if testResult.ActualOutcome == "" {
				continue // TODO: revisit
			}
			t.Run(testResult.Name, func(subT *testing.T) {
				switch {
				// A `go test` timeout is a panic, so handle this case before
				// the general panic case.
				case testResult.Timeout():
					subT.Error("timeout")

				// Handle a failure to run a test
				case testResult.Panic():
					subT.Error("panic")

				// Handle a difference between expected and actual outcome of a completed test
				case testResult.Red():
					subT.Error("expected test to " + testResult.ExpectedOutcome + "; actual outcome: " + testResult.ActualOutcome)
				}
			})
		}
	}
}
