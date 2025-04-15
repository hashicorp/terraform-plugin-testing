// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func ComputedStringAttribute(name string) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.String,
		Computed: true,
	}
}

func OptionalStringAttribute(name string) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.String,
		Optional: true,
	}
}

func RequiredStringAttribute(name string) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.String,
		Required: true,
	}
}
