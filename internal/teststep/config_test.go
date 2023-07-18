// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfigHasProviderBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config   configuration
		expected bool
	}{
		"no-config": {
			config:   configuration{},
			expected: false,
		},
		"provider-meta-attribute": {
			config: configuration{
				raw: `
resource "test_test" "test" {
  provider = test.test
}
`,
			},
			expected: false,
		},
		"provider-object-attribute": {
			config: configuration{
				raw: `
resource "test_test" "test" {
  test = {
	provider = {
	  test = true
	}
  }
}
`,
			},
			expected: false,
		},
		"provider-string-attribute": {
			config: configuration{
				raw: `
resource "test_test" "test" {
  test = {
	provider = "test"
  }
}
`,
			},
			expected: false,
		},
		"provider-block-quoted-with-attributes": {
			config: configuration{
				raw: `
provider "test" {
  test = true
}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"provider-block-unquoted-with-attributes": {
			config: configuration{
				raw: `
provider test {
  test = true
}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"provider-block-quoted-without-attributes": {
			config: configuration{
				raw: `
provider "test" {}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
		"provider-block-unquoted-without-attributes": {
			config: configuration{
				raw: `
provider test {}

resource "test_test" "test" {}
`,
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.config.HasProviderBlock(context.Background())

			if testCase.expected != got {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}

func TestStepMergedConfig(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testCaseProviderConfig string
		testStepProviderConfig string
		config                 configuration
		expected               string
	}{
		"testcase-externalproviders-and-protov5providerfactories": {
			testCaseProviderConfig: `
terraform {
 required_providers {
   externaltest = {
     source = "registry.terraform.io/hashicorp/externaltest"
     version = "1.2.3"
   }
 }
}

provider "externaltest" {}

`,
			config: configuration{
				raw: `
resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
			},
			expected: `
terraform {
 required_providers {
   externaltest = {
     source = "registry.terraform.io/hashicorp/externaltest"
     version = "1.2.3"
   }
 }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"testcase-externalproviders-and-providerfactories": {
			testCaseProviderConfig: `
terraform {
  required_providers {
    externaltest = {
      source = "registry.terraform.io/hashicorp/externaltest"
      version = "1.2.3"
    }
  }
}

provider "externaltest" {}

`,
			config: configuration{
				raw: `
resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    externaltest = {
      source = "registry.terraform.io/hashicorp/externaltest"
      version = "1.2.3"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"testcase-externalproviders-missing-source-and-versionconstraint": {
			testCaseProviderConfig: `
provider "test" {}

`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
provider "test" {}


resource "test_test" "test" {}
`,
		},
		"testcase-externalproviders-source-and-versionconstraint": {
			testCaseProviderConfig: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

provider "test" {}

`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

provider "test" {}


resource "test_test" "test" {}
`,
		},
		"testcase-externalproviders-source": {
			testCaseProviderConfig: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
    }
  }
}

provider "test" {}

`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
    }
  }
}

provider "test" {}


resource "test_test" "test" {}
`,
		},
		"testcase-externalproviders-versionconstraint": {
			testCaseProviderConfig: `
terraform {
  required_providers {
    test = {
      version = "1.2.3"
    }
  }
}

provider "test" {}

`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      version = "1.2.3"
    }
  }
}

provider "test" {}


resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders": {
			testStepProviderConfig: `
terraform {
  required_providers {
    externaltest = {
      source = "registry.terraform.io/hashicorp/externaltest"
      version = "1.2.3"
    }
  }
}

provider "externaltest" {}

`,
			config: configuration{
				raw: `
resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    externaltest = {
      source = "registry.terraform.io/hashicorp/externaltest"
      version = "1.2.3"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"teststep-externalproviders-config-with-provider-block-quoted": {
			testStepProviderConfig: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

`,
			config: configuration{
				raw: `
provider "test" {}

resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}


provider "test" {}

resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-config-with-provider-block-unquoted": {
			testStepProviderConfig: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

`,
			config: configuration{
				raw: `
provider test {}

resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}


provider test {}

resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-config-with-terraform-block": {
			testStepProviderConfig: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

provider "test" {}

`,
			config: configuration{
				raw: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-missing-source-and-versionconstraint": {
			testStepProviderConfig: `
provider "test" {}
`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
provider "test" {}

resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-source-and-versionconstraint": {
			testStepProviderConfig: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

provider "test" {}

`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}

provider "test" {}


resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-source": {
			testStepProviderConfig: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
    }
  }
}

provider "test" {}

`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
    }
  }
}

provider "test" {}


resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-versionconstraint": {
			testStepProviderConfig: `
terraform {
  required_providers {
    test = {
      version = "1.2.3"
    }
  }
}

provider "test" {}

`,
			config: configuration{
				raw: `
resource "test_test" "test" {}
`,
			},
			expected: `
terraform {
  required_providers {
    test = {
      version = "1.2.3"
    }
  }
}

provider "test" {}


resource "test_test" "test" {}
`,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.config.MergedConfig(context.Background(), testCase.testCaseProviderConfig, testCase.testStepProviderConfig).raw

			if diff := cmp.Diff(strings.TrimSpace(got), strings.TrimSpace(testCase.expected)); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
