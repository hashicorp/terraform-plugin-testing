// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import "fmt"

type ValuesDifferAny struct{}

func (v ValuesDifferAny) CompareValues(values ...any) error {
	if len(values) < 2 {
		return nil
	}

	vals := map[any]int{}

	for i := 0; i < len(values); i++ {
		vals[values[i]]++
	}

	if len(vals) < 2 {
		for k, v := range vals {
			if v > 1 {
				return fmt.Errorf("expected values to differ, but value is duplicated: %v", k)
			}
		}
	}

	return nil
}
