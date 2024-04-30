// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import "fmt"

type ValuesDiffer struct{}

func (v ValuesDiffer) CompareValues(values ...any) error {
	for i := 1; i < len(values); i++ {
		if values[i-1] == values[i] {
			return fmt.Errorf("expected values to differ, but they are the same: %v == %v", values[i-1], values[i])
		}
	}

	return nil
}
