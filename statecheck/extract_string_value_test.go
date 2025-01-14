// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestExtractStringValue_Basic(t *testing.T) {
	t.Parallel()

	// targetVar will be set to the extracted value.
	var targetVar string

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
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExtractStringValue(
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						&targetVar,
					),
				},
			},
		},
	})

	t.Run("check_target_var", func(t *testing.T) {
		if err := testAccAssertStringEquals("str", targetVar); err != nil {
			t.Errorf("extracted value does not match expected value: %v", err)
		}
	})
}

// testAccAssertStringEquals compares the expected and target string values.
func testAccAssertStringEquals(expected string, targetVar string) error {
	if targetVar != expected {
		return fmt.Errorf("expected targetVar to be %v, got %v", expected, targetVar)
	}
	return nil
}
