# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "random_password" "test" {
  length = 8

  numeric = false
}