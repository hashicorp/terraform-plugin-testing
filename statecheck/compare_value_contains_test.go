// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestCompareValueContains_CheckState_Map_ValuesSame_DifferError(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					map_attribute = {
						"a": "str2",
						"b": "str3",
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("map_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame{},
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str3 != str"),
			},
		},
	})
}

func TestCompareValueContains_CheckState_Map_ValuesSame(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					map_attribute = {
						"a": "str2",
						"b": test_resource.one.string_attribute,
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("map_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame{},
					),
				},
			},
		},
	})
}

func TestCompareValueContains_CheckState_Map_ValuesDiffer_SameError(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					map_attribute = {
						"a": test_resource.one.string_attribute,
						"b": test_resource.one.string_attribute,
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("map_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer{},
					),
				},
				ExpectError: regexp.MustCompile("expected values to differ, but they are the same: str == str"),
			},
		},
	})
}

func TestCompareValueContains_CheckState_Map_ValuesDiffer(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					map_attribute = {
						"a": test_resource.one.string_attribute,
						"b": "str2",
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("map_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer{},
					),
				},
			},
		},
	})
}

func TestCompareValueContains_CheckState_Set_ValuesSame_DifferError(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					set_attribute = [
						"str2",
						"str3"
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("set_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame{},
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str3 != str"),
			},
		},
	})
}

func TestCompareValueContains_CheckState_Set_ValuesSame(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					set_attribute = [
						"str2",
						test_resource.one.string_attribute
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("set_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame{},
					),
				},
			},
		},
	})
}

func TestCompareValueContains_CheckState_Set_ValuesDiffer_SameError(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					set_attribute = [
						test_resource.one.string_attribute,
						test_resource.one.string_attribute,
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("set_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer{},
					),
				},
				ExpectError: regexp.MustCompile("expected values to differ, but they are the same: str == str"),
			},
		},
	})
}

func TestCompareValueContains_CheckState_Set_ValuesDiffer(t *testing.T) {
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
					string_attribute = "str"
				}

				resource "test_resource" "two" {
					set_attribute = [
						test_resource.one.string_attribute,
						"str2"
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueContains(
						"test_resource.two",
						tfjsonpath.New("set_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer{},
					),
				},
			},
		},
	})
}
