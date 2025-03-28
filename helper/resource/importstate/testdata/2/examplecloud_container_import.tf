# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "examplecloud_container" "test" {
  name     = "somevalue"
  location = "westeurope"
}

import {
  to = examplecloud_container.test
  id = "examplecloud_container.test"
}
