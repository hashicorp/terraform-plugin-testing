# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "random_password" "test" {
  length = var.length

  numeric = var.numeric
}