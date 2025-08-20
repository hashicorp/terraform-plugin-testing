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

func TestCompareValue_CheckQuery_NoQueryValues(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := querycheck.CompareValue(compare.ValuesSame())

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
				ConfigQueryChecks: []querycheck.QueryCheck{
					// No query values have been added
					boolValuesDiffer,
				},
				ExpectError: regexp.MustCompile(`resource addresses index out of bounds: 0`),
			},
		},
	})
}

func TestCompareValue_CheckQuery_ValuesSame_ValueDiffersError(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := querycheck.CompareValue(compare.ValuesSame())

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
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}
				`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = false
				}
				`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile(`expected values to be the same, but they differ: true != false`),
			},
		},
	})
}

func TestCompareValue_CheckQuery_ValuesSame(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := querycheck.CompareValue(compare.ValuesSame())

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
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}
				`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
		},
	})
}

func TestCompareValue_CheckQuery_ValuesDiffer(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := querycheck.CompareValue(compare.ValuesDiffer())

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
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = false
				}
				`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}
				`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					boolValuesDiffer.AddQueryValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
		},
	})
}
