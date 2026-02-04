# Copyright IBM Corp. 2014, 2026
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
  length = var.length

  numeric = var.numeric
}

variable "length" {
  type = number
}

variable "numeric" {
  type = bool
}