// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestCompareValuePairs_CheckState_ValuesSame_DifferError(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(anTestProvider),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
					float_attribute = 1.234
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
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

func TestCompareValuePairs_CheckState_ValuesSame(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(anTestProvider),
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
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

func TestCompareValuePairs_CheckState_ValuesDiffer_SameError(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(anTestProvider),
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
					statecheck.CompareValuePairs(
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

func TestCompareValuePairs_CheckState_ValuesDiffer(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(anTestProvider),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
					float_attribute = 1.234
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
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
