# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

resource "test_test" "test" {
  test = {
    terraform = {
      test = true
    }
  }
}