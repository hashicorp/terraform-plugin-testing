# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    random = {
      source = "registry.terraform.io/hashicorp/random"
      version = "3.2.0"
    }
  }
}

provider "random" {}

resource "random_password" "test" {
  length = 8

  numeric = false
}