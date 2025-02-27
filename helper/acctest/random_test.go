// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acctest

import (
	"crypto/rsa"
	"net/netip"
	"regexp"
	"slices"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestRandIntRange(t *testing.T) {
	t.Parallel()

	v := RandInt()
	if vv := RandIntRange(v, v+1); vv != v {
		t.Errorf("expected RandIntRange(%d, %d) to return %d, got %d", v, v+1, v, vv)
	}
}

func TestRandIpAddress(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		s           string
		expected    *regexp.Regexp
		expectedErr string
	}{
		{
			s:        "0.0.0.0/0",
			expected: regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`),
		},
		{
			s:        "1.1.1.1/32",
			expected: regexp.MustCompile(`^1\.1\.1\.1$`),
		},
		{
			s:        "10.0.0.0/8",
			expected: regexp.MustCompile(`^10\.\d{1,3}\.\d{1,3}\.\d{1,3}$`),
		},
		{
			s:        "10.0.0.0/15",
			expected: regexp.MustCompile(`^10\.[01]\.\d{1,3}\.\d{1,3}$`),
		},
		{
			s:        "449d:e5f1:14b1:ddf3:8525:7e9e:4a0d:4a82/128",
			expected: regexp.MustCompile(`^449d:e5f1:14b1:ddf3:8525:7e9e:4a0d:4a82$`),
		},
		{
			s:        "2001:db8::/112",
			expected: regexp.MustCompile(`^2001:db8::[[:xdigit:]]{1,4}$`),
		},
		{
			s:           "abcdefg",
			expectedErr: "netip.ParsePrefix(\"abcdefg\"): no '/'",
		},
	}

	for i, tc := range testCases {
		v, err := RandIpAddress(tc.s)
		if err != nil {
			msg := err.Error()
			if tc.expectedErr == "" {
				t.Fatalf("expected test case %d to succeed but got error %q, ", i, msg)
			}
			if msg != tc.expectedErr {
				t.Fatalf("expected test case %d to fail with %q but got %q", i, tc.expectedErr, msg)
			}

			return
		}

		if !tc.expected.MatchString(v) {
			t.Errorf("expected test case %d to return %q but got %q", i, tc.expected, v)
		}

		if !netip.MustParsePrefix(tc.s).Contains(netip.MustParseAddr(v)) {
			t.Errorf("unexpected IP (%s) for prefix (%s)", v, tc.s)
		}
	}
}

func TestRandSSHKeyPair(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		comment string
	}{
		"comment": {
			comment: "test comment",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotPublicKey, gotPrivateKey, err := RandSSHKeyPair(testCase.comment)

			if err != nil {
				t.Fatalf("error generating SSH keypair: %s", err)
			}

			// Ensure public key is parsable and OpenSSH authorized keys
			// format.
			_, gotComment, _, _, err := ssh.ParseAuthorizedKey([]byte(gotPublicKey))

			if err != nil {
				t.Errorf("error parsing SSH public key: %s", err)
			}

			if gotComment != testCase.comment {
				t.Errorf("expected %q public key comment, got: %s", testCase.comment, gotComment)
			}

			// Ensure private key is parsable, has no passphrase, RSA, and
			// is 1024 bits for SDK version compatibility.
			gotPrivateKeyIface, err := ssh.ParseRawPrivateKey([]byte(gotPrivateKey))

			if err != nil {
				t.Errorf("error parsing SSH private key: %s", err)
			}

			rsaPrivateKey, ok := gotPrivateKeyIface.(*rsa.PrivateKey)

			if !ok {
				t.Fatalf("expected *rsa.PrivateKey SSH private key, got: %T", gotPrivateKeyIface)
			}

			if rsaPrivateKey.N.BitLen() != 1024 {
				t.Errorf("expected 1024 bit SSH private key, got: %d", rsaPrivateKey.N.BitLen())
			}
		})
	}
}

func TestInverseMask(t *testing.T) {
	t.Parallel()

	type testCase struct {
		prefixLen int
		byteLen   int
		expected  []byte
	}

	testCases := map[string]testCase{
		"0-bit_ipv4": {
			prefixLen: 0,
			byteLen:   4,
			expected:  []byte{255, 255, 255, 255},
		},
		"7-bit_ipv4": {
			prefixLen: 7,
			byteLen:   4,
			expected:  []byte{1, 255, 255, 255},
		},
		"8-bit_ipv4": {
			prefixLen: 8,
			byteLen:   4,
			expected:  []byte{0, 255, 255, 255},
		},
		"9-bit_ipv4": {
			prefixLen: 9,
			byteLen:   4,
			expected:  []byte{0, 127, 255, 255},
		},
		"27-bit_ipv4": {
			prefixLen: 27,
			byteLen:   4,
			expected:  []byte{0, 0, 0, 31},
		},
		"32-bit_ipv4": {
			prefixLen: 32,
			byteLen:   4,
			expected:  []byte{0, 0, 0, 0},
		},
		"32-bit_ipv6": {
			prefixLen: 32,
			byteLen:   16,
			expected:  []byte{0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		"64-bit_ipv6": {
			prefixLen: 64,
			byteLen:   16,
			expected:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		"128-bit_ipv6": {
			prefixLen: 128,
			byteLen:   16,
			expected:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result, err := inverseMask(tCase.prefixLen, tCase.byteLen)
			if err != nil {
				t.Fatal(err)
			}

			if slices.Compare(tCase.expected, result) != 0 {
				t.Fatalf("expected %v, got %v", tCase.expected, result)
			}
		})
	}
}
