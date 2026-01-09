// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func RequiredBoolAttribute(name string) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.Bool,
		Required: true,
	}
}

func OptionalComputedListAttribute(name string, elementType tftypes.Type) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.List{ElementType: elementType},
		Optional: true,
		Computed: true,
	}
}

func RequiredListAttribute(name string, elementType tftypes.Type) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.List{ElementType: elementType},
		Required: true,
	}
}

func RequiredNumberAttribute(name string) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.Number,
		Required: true,
	}
}

func OptionalNumberAttribute(name string) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:     name,
		Type:     tftypes.Number,
		Optional: true,
	}
}

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
