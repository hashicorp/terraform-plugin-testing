// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestExtractBoolValue_Basic(t *testing.T) {
	t.Parallel()

	// targetVar will be set to the extracted value.
	var targetVar bool

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExtractBoolValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						&targetVar,
					),
				},
			},
		},
	})

	t.Run("CheckTargetVar", func(t *testing.T) {
		if err := testAccAssertBoolEquals(true, targetVar); err != nil {
			t.Errorf("Error in testAccAssertBoolEquals: %v", err)
		}
	})
}

func TestExtractBoolValue_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	// targetVar will be set to the extracted value.
	var targetVar bool

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExtractBoolValue(
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						&targetVar,
					),
				},
				ExpectError: regexp.MustCompile(`invalid type for attribute \'float_attribute\' in \'test_resource\.one\'. Expected: bool, Got: json\.Number`),
			},
		},
	})
}

func TestExtractBoolValue_Null(t *testing.T) {
	t.Parallel()

	// targetVar will be set to the extracted value.
	var targetVar bool

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = null
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExtractBoolValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						&targetVar,
					),
				},
				ExpectError: regexp.MustCompile(`nil: result for attribute \'bool_attribute\' in \'test_resource.one\'`),
			},
		},
	})
}

func TestExtractBoolValue_ResourceNotFound(t *testing.T) {
	t.Parallel()

	// targetVar will be set to the extracted value.
	var targetVar bool

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExtractBoolValue(
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
						&targetVar,
					),
				},
				ExpectError: regexp.MustCompile("test_resource.two - Resource not found in state"),
			},
		},
	})
}

// testAccAssertBoolEquals compares the expected and target bool values.
func testAccAssertBoolEquals(expected bool, targetVar bool) error {
	if targetVar != expected {
		return fmt.Errorf("expected targetVar to be %v, got %v", expected, targetVar)
	}
	return nil
}
