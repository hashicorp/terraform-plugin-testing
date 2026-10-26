// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Privacy rules
//
// Findings are written to a user-owned directory that outlives the test run
// and may be shared, so the collector records structure and never values:
//
//   - It does not persist configuration text, attribute values, plan or state
//     JSON, query result values, display names, or assertion matchers.
//   - From state, it reads only top-level attribute names and whether each
//     value is null or known. Nested map keys can be user data.
//   - Instance keys in resource addresses are replaced with [*] by
//     safeAddress, because they can contain arbitrary user strings.
//
// It does persist provider and feature names, normalized resource addresses,
// schema attribute names, known/null classifications, operations, test names,
// and test package names.

package semanticcoverage

import (
	"strings"
)

// Terraform instance keys can contain arbitrary strings, including ']'.
// Never persist their contents or the suffix of a malformed address.
func safeAddress(address string) string {
	var result strings.Builder
	for i := 0; i < len(address); i++ {
		if address[i] != '[' {
			result.WriteByte(address[i])
			continue
		}
		result.WriteString("[*]")
		quoted, escaped := false, false
		closed := false
		for i++; i < len(address); i++ {
			switch {
			case escaped:
				escaped = false
			case quoted && address[i] == '\\':
				escaped = true
			case address[i] == '"':
				quoted = !quoted
			case address[i] == ']' && !quoted:
				closed = true
			}
			if closed {
				break
			}
		}
		if !closed {
			break
		}
	}
	return result.String()
}
