// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actioncheck_test

import (
	"github.com/hashicorp/terraform-plugin-testing/actioncheck"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Example of how to use action checks in a test case
func ExampleExpectProgressMessageContains() {
	// This is how you would use action checks in a real test
	_ = resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: `
					action "aws_lambda_invoke" "test" {
						config {
							function_name = "my-function"
							payload       = "{\"key\":\"value\"}"
							log_type      = "Tail"
						}
					}
				`,
				ActionChecks: []actioncheck.ActionCheck{
					// Check that the action produces a message containing log output
					resource.TestCheckProgressMessageContains("aws_lambda_invoke.test", "Lambda function logs:"),

					// Check that we get exactly 2 progress messages (success + logs)
					resource.TestCheckProgressMessageCount("aws_lambda_invoke.test", 2),

					// Check that messages appear in the expected sequence
					resource.TestCheckProgressMessageSequence("aws_lambda_invoke.test", []string{
						"invoked successfully",
						"Lambda function logs:",
					}),
				},
			},
		},
	}
}
