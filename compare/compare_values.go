// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

type ValueComparer interface {
	CompareValues(values ...any) error
}
