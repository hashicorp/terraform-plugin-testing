// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

// OutputValueParams is supplied during instantiation of an output value plan check,
// and specifies the address, and optional value path for an output value.
//
// For example, if an output has been defined to point at a specific value:
//
//	resource "time_static" "one" {}
//
//	output "string_attribute" {
//	    value = time_static.one.rfc3339
//	}
//
// Then the value can be addressed directly, and does not require a valuePath:
//
//	plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
//	    OutputAddress: "string_attribute",
//	}),
//
// However, if an output has been specified to point at an object or a collection.
// For example:
//
//	resource "time_static" "one" {}
//
//	output "string_attribute" {
//	    value = time_static.one
//	}
//
// Then a specific value, for instance `rfc3339`, cannot be addressed directly, and requires a valuePath:
//
//	plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
//	    OutputAddress: "string_attribute",
//	    ValuePath: tfjsonpath.New("rfc3339"),
//	}),
type OutputValueParams struct {
	OutputAddress string
	ValuePath     tfjsonpath.Path
}
