// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"regexp"
)

var _ Check = stringRegularExpression{}

type stringRegularExpression struct {
	regex *regexp.Regexp
}

// CheckValue determines whether the passed value is of type string, and
// contains a sequence of bytes that match the regular expression supplied
// to StringRegularExpression.
func (v stringRegularExpression) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringRegularExpression check, got: %T", other)
	}

	if !v.regex.MatchString(otherVal) {
		return fmt.Errorf("expected regex match %s for StringRegularExpression check, got: %s", v.regex.String(), otherVal)
	}

	return nil
}

// String returns the string representation of the value.
func (v stringRegularExpression) String() string {
	return v.regex.String()
}

// StringRegularExpression returns a Check for asserting equality between the
// supplied regular expression and a value passed to the CheckValue method.
func StringRegularExpression(regex *regexp.Regexp) stringRegularExpression {
	return stringRegularExpression{
		regex: regex,
	}
}
