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

func TestExpectContains_CheckState_ResourceOneNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}
				
				resource "test_resource" "two" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectContains(
						"does_not_exist_one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile("does_not_exist_one - Resource not found in state"),
			},
		},
	})
}

func TestExpectContains_CheckState_ResourceTwoNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}
				
				resource "test_resource" "two" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectContains(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"does_not_exist_two",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile("does_not_exist_two - Resource not found in state"),
			},
		},
	})
}

func TestExpectContains_CheckState_AttributeOneNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}
				
				resource "test_resource" "two" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectContains(
						"test_resource.one",
						tfjsonpath.New("does_not_exist_one"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile("path not found: specified key does_not_exist_one not found in map at does_not_exist_one"),
			},
		},
	})
}

func TestExpectContains_CheckState_AttributeTwoNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}
				
				resource "test_resource" "two" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectContains(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("does_not_exist_two"),
					),
				},
				ExpectError: regexp.MustCompile("path not found: specified key does_not_exist_two not found in map at does_not_exist_two"),
			},
		},
	})
}

func TestExpectContains_CheckState_AttributeOneNotSet(t *testing.T) {
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
				
				resource "test_resource" "two" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectContains(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile("test_resource.one.bool_attribute is not a collection"),
			},
		},
	})
}

func TestExpectContains_CheckState_NotFound(t *testing.T) {
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
					string_attribute = "value1"
				}
				
				resource "test_resource" "two" {
					set_attribute = [
						test_resource.one.string_attribute,
						"value2"
					]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectContains(
						"test_resource.two",
						tfjsonpath.New("set_attribute"),
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile(`value of test_resource\.one\.bool_attribute is not found in value of test_resource\.two\.set_attribute`),
			},
		},
	})
}

func TestExpectContains_CheckState_Found(t *testing.T) {
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
					string_attribute = "value1"
				}
				
				resource "test_resource" "two" {
					set_attribute = [
						test_resource.one.string_attribute,
						"value2"
					]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectContains(
						"test_resource.two",
						tfjsonpath.New("set_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
					),
				},
			},
		},
	})
}
