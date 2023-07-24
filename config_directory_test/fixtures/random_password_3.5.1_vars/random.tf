# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    random = {
      source = "registry.terraform.io/hashicorp/random"
      version = "3.5.1"
    }
  }
}

provider "random" {}

resource "random_password" "test" {
  length = var.length

  numeric = var.numeric
}