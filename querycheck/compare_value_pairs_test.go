// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestCompareValuePairs_CheckQuery_ValuesSame_DifferError(t *testing.T) {
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
					float_attribute = 1.234
				}
				`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.CompareValuePairs(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: true != 1.234"),
			},
		},
	})
}

func TestCompareValuePairs_CheckQuery_ValuesSame(t *testing.T) {
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
				}
				`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.CompareValuePairs(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValuePairs_CheckQuery_ValuesDiffer_SameError(t *testing.T) {
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
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.CompareValuePairs(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
						compare.ValuesDiffer(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to differ, but they are the same: true == true"),
			},
		},
	})
}

func TestCompareValuePairs_CheckQuery_ValuesDiffer(t *testing.T) {
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
					float_attribute = 1.234
				}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.CompareValuePairs(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}
