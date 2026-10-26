// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"testing"
)

func TestSafeAddressHidesInstanceKeysContainingBrackets(t *testing.T) {
	t.Parallel()

	input := `module.test["secret]tail"].example_item.test["other-key"]`
	if got := safeAddress(input); got != "module.test[*].example_item.test[*]" {
		t.Fatalf("normalized address %q", got)
	}
}
