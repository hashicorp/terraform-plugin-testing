// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package echo contains a protocol v6 Terraform provider that can be used to transfer data from
// provider configuration to state via a managed resource. This is only meant for provider acceptance testing
// of data that is not stored to state, such as ephemeral resources.
//
// Example Usage:
//
//	ephemeral "examplecloud_thing" "this" {
//		name = "thing-one"
//	}
//
//	provider "echo" {
//		data = ephemeral.examplecloud_thing.this
//	}
//
//	resource "echo_test" "echo" {} // `echo_test.echo.data` will contain the ephemeral data from `ephemeral.examplecloud_thing.this`
package echo
