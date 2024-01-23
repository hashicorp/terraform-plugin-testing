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

func TestExpectMatchingValues_CheckState_ResourceOneNotFound(t *testing.T) {
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
					statecheck.ExpectMatchingValues(
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

func TestExpectMatchingValues_CheckState_ResourceTwoNotFound(t *testing.T) {
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
					statecheck.ExpectMatchingValues(
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

func TestExpectMatchingValues_CheckState_AttributeOneNotFound(t *testing.T) {
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
					statecheck.ExpectMatchingValues(
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

func TestExpectMatchingValues_CheckState_AttributeTwoNotFound(t *testing.T) {
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
					statecheck.ExpectMatchingValues(
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

func TestExpectMatchingValues_CheckState_AttributeValuesNotEqualNil(t *testing.T) {
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
					statecheck.ExpectMatchingValues(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile("values are not equal: true != <nil>"),
			},
		},
	})
}

func TestExpectMatchingValues_CheckState_AttributeValuesNotEqual(t *testing.T) {
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
				
				resource "test_resource" "two" {
					bool_attribute = false
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectMatchingValues(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile("values are not equal: true != false"),
			},
		},
	})
}

func TestExpectMatchingValues_CheckState_AttributeValuesEqual_Bool(t *testing.T) {
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
				
				resource "test_resource" "two" {
					bool_attribute = true
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectMatchingValues(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
		},
	})
}

func TestExpectMatchingValues_CheckState_AttributeValuesEqual_List(t *testing.T) {
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
					list_attribute = [
						"value1",
						"value2"
					]
				}
				
				resource "test_resource" "two" {
					list_attribute = [
						"value1",
						"value2"
					]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectMatchingValues(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						"test_resource.two",
						tfjsonpath.New("list_attribute"),
					),
				},
			},
		},
	})
}
