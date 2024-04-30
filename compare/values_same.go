// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import "fmt"

type ValuesSame struct{}

func (v ValuesSame) CompareValues(values ...any) error {
	for i := 1; i < len(values); i++ {
		if values[i-1] != values[i] {
			return fmt.Errorf("expected values to be the same, but they differ: %v != %v", values[i-1], values[i])
		}
	}

	return nil
}
