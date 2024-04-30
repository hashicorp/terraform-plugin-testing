// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import "fmt"

type ValuesSameAny struct{}

func (v ValuesSameAny) CompareValues(values ...any) error {
	if len(values) < 2 {
		return nil
	}

	vals := map[any]struct{}{}

	for i := 0; i < len(values); i++ {
		vals[values[i]] = struct{}{}
	}

	if len(vals) == len(values) {
		return fmt.Errorf("expected at least two values to be the same, but all values differ")
	}

	return nil
}
