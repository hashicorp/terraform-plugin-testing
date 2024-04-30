// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import "fmt"

type ValuesDifferAll struct{}

func (v ValuesDifferAll) CompareValues(values ...any) error {
	vals := map[any]struct{}{}

	for i := 0; i < len(values); i++ {
		if _, ok := vals[values[i]]; ok {
			return fmt.Errorf("expected values to differ, but value is duplicated: %v", values[i])
		}

		vals[values[i]] = struct{}{}
	}

	return nil
}
