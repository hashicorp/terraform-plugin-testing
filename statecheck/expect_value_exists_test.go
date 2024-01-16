// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestExpectValueExists_CheckState_ResourceNotFound(t *testing.T) {
	t.Parallel()

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
				ConfigStateChecks: r.ConfigStateChecks{
					statecheck.ExpectValueExists(
						"does_not_exist",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile(`does_not_exist - Resource not found in state`),
			},
		},
	})
}

func TestExpectValueExists_CheckState_AttributeNotFound(t *testing.T) {
	t.Parallel()

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
				ConfigStateChecks: r.ConfigStateChecks{
					statecheck.ExpectValueExists(
						"test_resource.one",
						tfjsonpath.New("does_not_exist"),
					),
				},
				ExpectError: regexp.MustCompile(`path not found: specified key does_not_exist not found in map at does_not_exist`),
			},
		},
	})
}

func TestExpectValueExists_CheckState_AttributeFound(t *testing.T) {
	t.Parallel()

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
				ConfigStateChecks: r.ConfigStateChecks{
					statecheck.ExpectValueExists(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
		},
	})
}
