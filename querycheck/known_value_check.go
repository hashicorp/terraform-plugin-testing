// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func NewKnownValueCheck(path tfjsonpath.Path, check knownvalue.Check) KnownValueCheck {
	return KnownValueCheck{
		Path:       path,
		KnownValue: check,
	}
}
