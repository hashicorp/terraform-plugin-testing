# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

resource "random_password" "test" {
  length = 8

  numeric = false
}