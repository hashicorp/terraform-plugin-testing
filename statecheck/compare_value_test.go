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

func TestCompareValue_CheckState_Bool_ValuesSame_ValueDiffersError(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := statecheck.CompareValue(compare.ValuesSame())

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
					boolValuesDiffer.AddStateValue(
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
				ConfigStateChecks: []statecheck.StateCheck{
					boolValuesDiffer.AddStateValue(
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
				ConfigStateChecks: []statecheck.StateCheck{
					boolValuesDiffer.AddStateValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile(`expected values to be the same, but they differ: true != false`),
			},
		},
	})
}

func TestCompareValue_CheckState_Bool_ValuesSame(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := statecheck.CompareValue(compare.ValuesSame())

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
					boolValuesDiffer.AddStateValue(
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
				ConfigStateChecks: []statecheck.StateCheck{
					boolValuesDiffer.AddStateValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
		},
	})
}

func TestCompareValue_CheckState_Bool_ValuesDiffer_ValueSameError(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())

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
					boolValuesDiffer.AddStateValue(
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
				ConfigStateChecks: []statecheck.StateCheck{
					boolValuesDiffer.AddStateValue(
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
				ConfigStateChecks: []statecheck.StateCheck{
					boolValuesDiffer.AddStateValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
				ExpectError: regexp.MustCompile(`expected values to differ, but they are the same: false == false`),
			},
		},
	})
}

func TestCompareValue_CheckState_Bool_ValuesDiffer(t *testing.T) {
	t.Parallel()

	boolValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())

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
					boolValuesDiffer.AddStateValue(
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
				ConfigStateChecks: []statecheck.StateCheck{
					boolValuesDiffer.AddStateValue(
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
				ConfigStateChecks: []statecheck.StateCheck{
					boolValuesDiffer.AddStateValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
		},
	})
}
