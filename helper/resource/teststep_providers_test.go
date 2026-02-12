// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestStepMergedConfig_TF_0_15(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testCase                TestCase
		testStep                TestStep
		configHasTerraformBlock bool
		configHasProviderBlock  bool
		expected                string
	}{
		"testcase-externalproviders-and-protov5providerfactories": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"localtest": nil,
				},
			},
			testStep: TestStep{
				Config: `
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
		"testcase-externalproviders-and-protov6providerfactories": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"localtest": nil,
				},
			},
			testStep: TestStep{
				Config: `
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
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"localtest": nil,
				},
			},
			testStep: TestStep{
				Config: `
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
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {},
				},
			},
			testStep: TestStep{
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
provider "test" {}

resource "test_test" "test" {}
`,
		},
		"testcase-externalproviders-source-and-versionconstraint": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
			},
			testStep: TestStep{
				Config: `
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
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source: "registry.terraform.io/hashicorp/test",
					},
				},
			},
			testStep: TestStep{
				Config: `
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
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						VersionConstraint: "1.2.3",
					},
				},
			},
			testStep: TestStep{
				Config: `
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
		"testcase-protov5providerfactories": {
			testCase: TestCase{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil,
				},
			},
			testStep: TestStep{
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
resource "test_test" "test" {}
`,
		},
		"testcase-protov6providerfactories": {
			testCase: TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil,
				},
			},
			testStep: TestStep{
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
resource "test_test" "test" {}
`,
		},
		"testcase-providerfactories": {
			testCase: TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil,
				},
			},
			testStep: TestStep{
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-and-protov5providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"localtest": nil,
				},
				Config: `
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
		"teststep-externalproviders-and-protov6providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"localtest": nil,
				},
				Config: `
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
		"teststep-externalproviders-and-providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"localtest": nil,
				},
				Config: `
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
provider "test" {}

resource "test_test" "test" {}
`,
			},
			configHasProviderBlock: true,
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
provider test {}

resource "test_test" "test" {}
`,
			},
			configHasProviderBlock: true,
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
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
			configHasTerraformBlock: true,
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {},
				},
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
provider "test" {}

resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-source-and-versionconstraint": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source: "registry.terraform.io/hashicorp/test",
					},
				},
				Config: `
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
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
		"teststep-protov5providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil,
				},
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
resource "test_test" "test" {}
`,
		},
		"teststep-protov6providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil,
				},
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
resource "test_test" "test" {}
`,
		},
		"teststep-providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil,
				},
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
resource "test_test" "test" {}
`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.testStep.mergedConfig(context.Background(), testCase.testCase, testCase.configHasTerraformBlock, testCase.configHasProviderBlock, tfversion.Version0_15_0)

			if err != nil {
				t.Errorf("cannot generate merged config: %s", err)
			}

			if diff := cmp.Diff(strings.TrimSpace(got), strings.TrimSpace(testCase.expected)); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStepMergedConfig_TF_1_6(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testCase                TestCase
		testStep                TestStep
		configHasTerraformBlock bool
		configHasProviderBlock  bool
		expected                string
	}{
		"testcase-externalproviders-and-protov5providerfactories": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"localtest": nil,
				},
			},
			testStep: TestStep{
				Config: `
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
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"testcase-externalproviders-and-protov6providerfactories": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"localtest": nil,
				},
			},
			testStep: TestStep{
				Config: `
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
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"testcase-externalproviders-and-providerfactories": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"localtest": nil,
				},
			},
			testStep: TestStep{
				Config: `
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
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"testcase-externalproviders-missing-source-and-versionconstraint": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {},
				},
			},
			testStep: TestStep{
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
provider "test" {}

resource "test_test" "test" {}
`,
		},
		"testcase-externalproviders-source-and-versionconstraint": {
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
			},
			testStep: TestStep{
				Config: `
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
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source: "registry.terraform.io/hashicorp/test",
					},
				},
			},
			testStep: TestStep{
				Config: `
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
			testCase: TestCase{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						VersionConstraint: "1.2.3",
					},
				},
			},
			testStep: TestStep{
				Config: `
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
		"testcase-protov5providerfactories": {
			testCase: TestCase{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil,
				},
			},
			testStep: TestStep{
				Config: `
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



resource "test_test" "test" {}
`,
		},
		"testcase-protov6providerfactories": {
			testCase: TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil,
				},
			},
			testStep: TestStep{
				Config: `
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



resource "test_test" "test" {}
`,
		},
		"testcase-providerfactories": {
			testCase: TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil,
				},
			},
			testStep: TestStep{
				Config: `
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



resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-and-protov5providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"localtest": nil,
				},
				Config: `
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
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"teststep-externalproviders-and-protov6providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"localtest": nil,
				},
				Config: `
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
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"teststep-externalproviders-and-providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"localtest": nil,
				},
				Config: `
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
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}


resource "externaltest_test" "test" {}

resource "localtest_test" "test" {}
`,
		},
		"teststep-externalproviders-config-with-provider-block-quoted": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
provider "test" {}

resource "test_test" "test" {}
`,
			},
			configHasProviderBlock: true,
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
provider test {}

resource "test_test" "test" {}
`,
			},
			configHasProviderBlock: true,
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
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
			configHasTerraformBlock: true,
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {},
				},
				Config: `
resource "test_test" "test" {}
`,
			},
			expected: `
provider "test" {}

resource "test_test" "test" {}
`,
		},
		"teststep-externalproviders-source-and-versionconstraint": {
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source: "registry.terraform.io/hashicorp/test",
					},
				},
				Config: `
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
			testCase: TestCase{},
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						VersionConstraint: "1.2.3",
					},
				},
				Config: `
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
		"teststep-protov5providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil,
				},
				Config: `
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



resource "test_test" "test" {}
`,
		},
		"teststep-protov6providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil,
				},
				Config: `
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



resource "test_test" "test" {}
`,
		},
		"teststep-providerfactories": {
			testCase: TestCase{},
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil,
				},
				Config: `
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



resource "test_test" "test" {}
`,
		},
	}

	v, err := version.NewVersion("1.6.0")

	if err != nil {
		t.Errorf("error generating version: %s", err)
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.testStep.mergedConfig(context.Background(), testCase.testCase, testCase.configHasTerraformBlock, testCase.configHasProviderBlock, v)

			if err != nil {
				t.Errorf("cannot generate merged config: %s", err)
			}

			if diff := cmp.Diff(strings.TrimSpace(got), strings.TrimSpace(testCase.expected)); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStepProviderConfig_TF_0_15(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testStep          TestStep
		skipProviderBlock bool
		expected          string
	}{
		"externalproviders-and-protov5providerfactories": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"localtest": nil,
				},
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
`,
		},
		"externalproviders-and-protov6providerfactories": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"localtest": nil,
				},
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
`,
		},
		"externalproviders-and-providerfactories": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"localtest": nil,
				},
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
`,
		},
		"externalproviders-missing-source-and-versionconstraint": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {},
				},
			},
			expected: `provider "test" {}`,
		},
		"externalproviders-skip-provider-block": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
			},
			skipProviderBlock: true,
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}
`,
		},
		"externalproviders-source-and-versionconstraint": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
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
`,
		},
		"externalproviders-source": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source: "registry.terraform.io/hashicorp/test",
					},
				},
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
`,
		},
		"externalproviders-versionconstraint": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						VersionConstraint: "1.2.3",
					},
				},
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
`,
		},
		"protov5providerfactories": {
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil,
				},
			},
			expected: ``,
		},
		"protov6providerfactories": {
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil,
				},
			},
			expected: ``,
		},
		"providerfactories": {
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil,
				},
			},
			expected: ``,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.testStep.providerConfig(context.Background(), testCase.skipProviderBlock, tfversion.Version0_15_0)

			if err != nil {
				t.Errorf("cannot generate provider config: %s", err)
			}

			if diff := cmp.Diff(strings.TrimSpace(got), strings.TrimSpace(testCase.expected)); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStepProviderConfig_TF_1_6(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testStep          TestStep
		skipProviderBlock bool
		expected          string
	}{
		"externalproviders-and-protov5providerfactories": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"localtest": nil,
				},
			},
			expected: `
terraform {
  required_providers {
    externaltest = {
      source = "registry.terraform.io/hashicorp/externaltest"
      version = "1.2.3"
    }
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}
`,
		},
		"externalproviders-and-protov6providerfactories": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"localtest": nil,
				},
			},
			expected: `
terraform {
  required_providers {
    externaltest = {
      source = "registry.terraform.io/hashicorp/externaltest"
      version = "1.2.3"
    }
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}
`,
		},
		"externalproviders-and-providerfactories": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"externaltest": {
						Source:            "registry.terraform.io/hashicorp/externaltest",
						VersionConstraint: "1.2.3",
					},
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"localtest": nil,
				},
			},
			expected: `
terraform {
  required_providers {
    externaltest = {
      source = "registry.terraform.io/hashicorp/externaltest"
      version = "1.2.3"
    }
    localtest = {
      source = "registry.terraform.io/hashicorp/localtest"
    }
  }
}

provider "externaltest" {}
`,
		},
		"externalproviders-missing-source-and-versionconstraint": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {},
				},
			},
			expected: `provider "test" {}`,
		},
		"externalproviders-skip-provider-block": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
			},
			skipProviderBlock: true,
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
      version = "1.2.3"
    }
  }
}
`,
		},
		"externalproviders-source-and-versionconstraint": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source:            "registry.terraform.io/hashicorp/test",
						VersionConstraint: "1.2.3",
					},
				},
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
`,
		},
		"externalproviders-source": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						Source: "registry.terraform.io/hashicorp/test",
					},
				},
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
`,
		},
		"externalproviders-versionconstraint": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {
						VersionConstraint: "1.2.3",
					},
				},
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
`,
		},
		"protov5providerfactories": {
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil,
				},
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
    }
  }
}`,
		},
		"protov6providerfactories": {
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil,
				},
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
    }
  }
}`,
		},
		"providerfactories": {
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil,
				},
			},
			expected: `
terraform {
  required_providers {
    test = {
      source = "registry.terraform.io/hashicorp/test"
    }
  }
}`,
		},
	}

	v, err := version.NewVersion("1.6.0")

	if err != nil {
		t.Errorf("error generating version: %s", err)
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.testStep.providerConfig(context.Background(), testCase.skipProviderBlock, v)

			if err != nil {
				t.Errorf("cannot generate provider config: %s", err)
			}

			if diff := cmp.Diff(strings.TrimSpace(got), strings.TrimSpace(testCase.expected)); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTest_TestStep_ExternalProviders(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				Config: "# not empty",
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source: "registry.terraform.io/hashicorp/null",
					},
				},
			},
		},
	})
}

func TestTest_TestStep_ExternalProviders_DifferentProviders(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source: "registry.terraform.io/hashicorp/null",
					},
				},
			},
			{
				Config: `resource "random_pet" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"random": {
						Source: "registry.terraform.io/hashicorp/random",
					},
				},
			},
		},
	})
}

func TestTest_TestStep_ExternalProviders_DifferentVersions(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source:            "registry.terraform.io/hashicorp/null",
						VersionConstraint: "3.1.0",
					},
				},
			},
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source:            "registry.terraform.io/hashicorp/null",
						VersionConstraint: "3.1.1",
					},
				},
			},
		},
	})
}

func TestTest_TestStep_ExternalProviders_Error(t *testing.T) {
	t.Parallel()

	plugintest.TestExpectTFatal(t, func() {
		Test(&mockT{}, TestCase{
			Steps: []TestStep{
				{
					Config: "# not empty",
					ExternalProviders: map[string]ExternalProvider{
						"testnonexistent": {
							Source: "registry.terraform.io/hashicorp/testnonexistent",
						},
					},
				},
			},
		})
	})
}

func TestTest_TestStep_ExternalProviders_NonHashiCorpNamespace(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ExternalProvider.Source is protocol version 6
		},
		Steps: []TestStep{
			{
				ExternalProviders: map[string]ExternalProvider{
					// This can be set to any provider outside the hashicorp namespace.
					// bflad/scaffoldingtest happens to be a published version of
					// terraform-provider-scaffolding-framework.
					"scaffoldingtest": {
						Source:            "registry.terraform.io/bflad/scaffoldingtest",
						VersionConstraint: "0.1.0",
					},
				},
				Config: `resource "scaffoldingtest_example" "test" {}`,
			},
		},
	})
}

func TestTest_TestStep_ExternalProvidersAndProviderFactories_NonHashiCorpNamespace(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ExternalProvider.Source is protocol version 6
		},
		Steps: []TestStep{
			{
				ExternalProviders: map[string]ExternalProvider{
					// This can be set to any provider outside the hashicorp namespace.
					// bflad/scaffoldingtest happens to be a published version of
					// terraform-provider-scaffolding-framework.
					"scaffoldingtest": {
						Source:            "registry.terraform.io/bflad/scaffoldingtest",
						VersionConstraint: "0.1.0",
					},
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"null": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return &schema.Provider{
							ResourcesMap: map[string]*schema.Resource{
								"null_resource": {
									CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
										d.SetId("test")
										return nil
									},
									DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									Schema: map[string]*schema.Schema{
										"triggers": {
											Elem:     &schema.Schema{Type: schema.TypeString},
											ForceNew: true,
											Optional: true,
											Type:     schema.TypeMap,
										},
									},
								},
							},
						}, nil
					},
				},
				Config: `
					resource "null_resource" "test" {}
					resource "scaffoldingtest_example" "test" {}
				`,
			},
		},
	})
}

func TestTest_TestStep_ExternalProviders_To_ProtoV6ProviderFactories(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		Steps: []TestStep{
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source:            "registry.terraform.io/hashicorp/null",
						VersionConstraint: "3.1.1",
					},
				},
			},
			{
				Config: `resource "null_resource" "test" {}`,
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"null": providerserver.NewProviderServer(testprovider.Provider{
						Resources: map[string]testprovider.Resource{
							"null_resource": {
								SchemaResponse: &resource.SchemaResponse{
									Schema: &tfprotov6.Schema{
										Block: &tfprotov6.SchemaBlock{
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:     "id",
													Type:     tftypes.String,
													Computed: true,
												},
												{
													Name:     "triggers",
													Type:     tftypes.Map{ElementType: tftypes.String},
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					}),
				},
			},
		},
	})
}

func TestTest_TestStep_ExternalProviders_To_ProviderFactories(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source:            "registry.terraform.io/hashicorp/null",
						VersionConstraint: "3.1.1",
					},
				},
			},
			{
				Config: `resource "null_resource" "test" {}`,
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"null": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return &schema.Provider{
							ResourcesMap: map[string]*schema.Resource{
								"null_resource": {
									CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
										d.SetId("test")
										return nil
									},
									DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									Schema: map[string]*schema.Schema{
										"triggers": {
											Elem:     &schema.Schema{Type: schema.TypeString},
											ForceNew: true,
											Optional: true,
											Type:     schema.TypeMap,
										},
									},
								},
							},
						}, nil
					},
				},
			},
		},
	})
}

func TestTest_TestStep_ExternalProviders_To_ProviderFactories_StateUpgraders(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source:            "registry.terraform.io/hashicorp/null",
						VersionConstraint: "3.1.1",
					},
				},
			},
			{
				Check:  TestCheckResourceAttr("null_resource.test", "id", "test-schema-version-1"),
				Config: `resource "null_resource" "test" {}`,
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"null": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return &schema.Provider{
							ResourcesMap: map[string]*schema.Resource{
								"null_resource": {
									CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
										d.SetId("test")
										return nil
									},
									DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									Schema: map[string]*schema.Schema{
										"triggers": {
											Elem:     &schema.Schema{Type: schema.TypeString},
											ForceNew: true,
											Optional: true,
											Type:     schema.TypeMap,
										},
									},
									SchemaVersion: 1, // null 3.1.3 is version 0
									StateUpgraders: []schema.StateUpgrader{
										{
											Type: cty.Object(map[string]cty.Type{
												"id":       cty.String,
												"triggers": cty.Map(cty.String),
											}),
											Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
												// null 3.1.3 sets the id attribute to a stringified random integer.
												// Double check that our resource wasn't created by this TestStep.
												id, ok := rawState["id"].(string)

												if !ok || id == "test" {
													return rawState, fmt.Errorf("unexpected rawState: %v", rawState)
												}

												rawState["id"] = "test-schema-version-1"

												return rawState, nil
											},
											Version: 0,
										},
									},
								},
							},
						}, nil
					},
				},
			},
		},
	})
}

func TestTest_TestStep_Taint(t *testing.T) {
	t.Parallel()

	var idOne, idTwo string

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				Check: ComposeAggregateTestCheckFunc(
					extractResourceAttr("test_resource.test", "id", &idOne),
				),
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": providerserver.NewProviderServer(testprovider.Provider{
						Resources: map[string]testprovider.Resource{
							"test_resource": {
								CreateResponse: &resource.CreateResponse{
									NewState: tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"id": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"id": tftypes.NewValue(tftypes.String, "test-value1"),
										},
									),
								},
								SchemaResponse: &resource.SchemaResponse{
									Schema: &tfprotov6.Schema{
										Block: &tfprotov6.SchemaBlock{
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:     "id",
													Type:     tftypes.String,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					}),
				},
			},
			{
				Taint:  []string{"test_resource.test"},
				Config: `resource "test_resource" "test" {}`,
				Check: ComposeAggregateTestCheckFunc(
					extractResourceAttr("test_resource.test", "id", &idTwo),
				),
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": providerserver.NewProviderServer(testprovider.Provider{
						Resources: map[string]testprovider.Resource{
							"test_resource": {
								CreateResponse: &resource.CreateResponse{
									NewState: tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"id": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"id": tftypes.NewValue(tftypes.String, "test-value2"),
										},
									),
								},
								SchemaResponse: &resource.SchemaResponse{
									Schema: &tfprotov6.Schema{
										Block: &tfprotov6.SchemaBlock{
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:     "id",
													Type:     tftypes.String,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					}),
				},
			},
		},
	})

	if idOne == idTwo {
		t.Errorf("taint is not causing destroy-create cycle, idOne == idTwo: %s == %s", idOne, idTwo)
	}
}

//nolint:unparam
func extractResourceAttr(resourceName string, attributeName string, attributeValue *string) TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource name %s not found in state", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]

		if !ok {
			return fmt.Errorf("attribute %s not found in resource %s state", attributeName, resourceName)
		}

		*attributeValue = attrValue

		return nil
	}
}

func TestTest_TestStep_ProtoV5ProviderFactories(t *testing.T) {
	t.Parallel()

	UnitTest(&mockT{}, TestCase{
		Steps: []TestStep{
			{
				Config: "# not empty",
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": providerserver.NewProtov5ProviderServer(testprovider.Protov5Provider{}),
				},
			},
		},
	})
}

func TestTest_TestStep_ProtoV5ProviderFactories_Error(t *testing.T) {
	t.Parallel()

	plugintest.TestExpectTFatal(t, func() {
		UnitTest(&mockT{}, TestCase{
			Steps: []TestStep{
				{
					Config: "# not empty",
					ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
						"test": func() (tfprotov5.ProviderServer, error) { //nolint:unparam // required signature
							return nil, fmt.Errorf("test")
						},
					},
				},
			},
		})
	})
}

func TestTest_TestStep_ProtoV6ProviderFactories(t *testing.T) {
	t.Parallel()

	UnitTest(&mockT{}, TestCase{
		Steps: []TestStep{
			{
				Config: "# not empty",
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": providerserver.NewProviderServer(testprovider.Provider{}),
				},
			},
		},
	})
}

func TestTest_TestStep_ProtoV6ProviderFactories_Error(t *testing.T) {
	t.Parallel()

	plugintest.TestExpectTFatal(t, func() {
		UnitTest(&mockT{}, TestCase{
			Steps: []TestStep{
				{
					Config: "# not empty",
					ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
						"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
							return nil, fmt.Errorf("test")
						},
					},
				},
			},
		})
	})
}

func TestTest_TestStep_ProtoV6ProviderFactories_To_ExternalProviders(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
		Steps: []TestStep{
			{
				Config: `resource "null_resource" "test" {}`,
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"null": providerserver.NewProviderServer(testprovider.Provider{
						Resources: map[string]testprovider.Resource{
							"null_resource": {
								CreateResponse: &resource.CreateResponse{
									NewState: tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"id":       tftypes.String,
												"triggers": tftypes.Map{ElementType: tftypes.String},
											},
										},
										map[string]tftypes.Value{
											"id":       tftypes.NewValue(tftypes.String, "test"),
											"triggers": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
										},
									),
								},
								SchemaResponse: &resource.SchemaResponse{
									Schema: &tfprotov6.Schema{
										Block: &tfprotov6.SchemaBlock{
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:     "id",
													Type:     tftypes.String,
													Computed: true,
												},
												{
													Name:     "triggers",
													Type:     tftypes.Map{ElementType: tftypes.String},
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					}),
				},
			},
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source: "registry.terraform.io/hashicorp/null",
					},
				},
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories(t *testing.T) {
	t.Parallel()

	UnitTest(&mockT{}, TestCase{
		Steps: []TestStep{
			{
				Config: "# not empty",
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return &schema.Provider{}, nil
					},
				},
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Error(t *testing.T) {
	t.Parallel()

	plugintest.TestExpectTFatal(t, func() {
		UnitTest(&mockT{}, TestCase{
			Steps: []TestStep{
				{
					Config: "# not empty",
					ProviderFactories: map[string]func() (*schema.Provider, error){
						"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
							return nil, fmt.Errorf("test")
						},
					},
				},
			},
		})
	})
}

func TestTest_TestStep_ProviderFactories_To_ExternalProviders(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				Config: `resource "null_resource" "test" {}`,
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"null": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return &schema.Provider{
							ResourcesMap: map[string]*schema.Resource{
								"null_resource": {
									CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
										d.SetId("test")
										return nil
									},
									DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									Schema: map[string]*schema.Schema{
										"triggers": {
											Elem:     &schema.Schema{Type: schema.TypeString},
											ForceNew: true,
											Optional: true,
											Type:     schema.TypeMap,
										},
									},
								},
							},
						}, nil
					},
				},
			},
			{
				Config: `resource "null_resource" "test" {}`,
				ExternalProviders: map[string]ExternalProvider{
					"null": {
						Source: "registry.terraform.io/hashicorp/null",
					},
				},
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Import_Inline(t *testing.T) {
	id := "none"

	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				Config: `resource "random_password" "test" { length = 12 }`,
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return &schema.Provider{
							ResourcesMap: map[string]*schema.Resource{
								"random_password": {
									DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
										return nil
									},
									Schema: map[string]*schema.Schema{
										"length": {
											Required: true,
											ForceNew: true,
											Type:     schema.TypeInt,
										},
										"result": {
											Type:      schema.TypeString,
											Computed:  true,
											Sensitive: true,
										},

										"id": {
											Computed: true,
											Type:     schema.TypeString,
										},
									},
									Importer: &schema.ResourceImporter{
										StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
											val := d.Id()

											d.SetId("none")

											err := d.Set("result", val)
											if err != nil {
												panic(err)
											}

											err = d.Set("length", len(val))
											if err != nil {
												panic(err)
											}

											return []*schema.ResourceData{d}, nil
										},
									},
								},
							},
						}, nil
					},
				},
				ResourceName:       "random_password.test",
				ImportState:        true,
				ImportStateId:      "Z=:cbrJE?Ltg",
				ImportStatePersist: true,
				ImportStateCheck: composeImportStateCheck(
					testCheckResourceAttrInstanceState(&id, "result", "Z=:cbrJE?Ltg"),
					testCheckResourceAttrInstanceState(&id, "length", "12"),
				),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Import_Inline_WithPersistMatch(t *testing.T) {
	var result1, result2 string

	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_password": {
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{
								"length": {
									Required: true,
									ForceNew: true,
									Type:     schema.TypeInt,
								},
								"result": {
									Type:      schema.TypeString,
									Computed:  true,
									Sensitive: true,
								},

								"id": {
									Computed: true,
									Type:     schema.TypeString,
								},
							},
							Importer: &schema.ResourceImporter{
								StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
									val := d.Id()

									d.SetId("none")

									err := d.Set("result", val)
									if err != nil {
										panic(err)
									}

									err = d.Set("length", len(val))
									if err != nil {
										panic(err)
									}

									return []*schema.ResourceData{d}, nil
								},
							},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				Config:             `resource "random_password" "test" { length = 12 }`,
				ResourceName:       "random_password.test",
				ImportState:        true,
				ImportStateId:      "Z=:cbrJE?Ltg",
				ImportStatePersist: true,
				ImportStateCheck: composeImportStateCheck(
					testExtractResourceAttrInstanceState("none", "result", &result1),
				),
			},
			{
				Config: `resource "random_password" "test" { length = 12 }`,
				Check: ComposeTestCheckFunc(
					testExtractResourceAttr("random_password.test", "result", &result2),
					testCheckAttributeValuesEqual(&result1, &result2),
				),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Import_Inline_WithoutPersist(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_password": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId("none")
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{
								"length": {
									Required: true,
									ForceNew: true,
									Type:     schema.TypeInt,
								},
								"result": {
									Type:      schema.TypeString,
									Computed:  true,
									Sensitive: true,
								},

								"id": {
									Computed: true,
									Type:     schema.TypeString,
								},
							},
							Importer: &schema.ResourceImporter{
								StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
									val := d.Id()

									d.SetId("none")

									err := d.Set("result", val)
									if err != nil {
										panic(err)
									}

									err = d.Set("length", len(val))
									if err != nil {
										panic(err)
									}

									return []*schema.ResourceData{d}, nil
								},
							},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				Config:             `resource "random_password" "test" { length = 12 }`,
				ResourceName:       "random_password.test",
				ImportState:        true,
				ImportStateId:      "Z=:cbrJE?Ltg",
				ImportStatePersist: false,
			},
			{
				Config: `resource "random_password" "test" { length = 12 }`,
				Check: ComposeTestCheckFunc(
					TestCheckNoResourceAttr("random_password.test", "result"),
				),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Import_External(t *testing.T) {
	id := "none"

	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config:             `resource "random_password" "test" { length = 12 }`,
				ResourceName:       "random_password.test",
				ImportState:        true,
				ImportStateId:      "Z=:cbrJE?Ltg",
				ImportStatePersist: true,
				ImportStateCheck: composeImportStateCheck(
					testCheckResourceAttrInstanceState(&id, "result", "Z=:cbrJE?Ltg"),
					testCheckResourceAttrInstanceState(&id, "length", "12"),
				),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Import_External_WithPersistMatch(t *testing.T) {
	var result1, result2 string

	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config:             `resource "random_password" "test" { length = 12 }`,
				ResourceName:       "random_password.test",
				ImportState:        true,
				ImportStateId:      "Z=:cbrJE?Ltg",
				ImportStatePersist: true,
				ImportStateCheck: composeImportStateCheck(
					testExtractResourceAttrInstanceState("none", "result", &result1),
				),
			},
			{
				Config: `resource "random_password" "test" { length = 12 }`,
				Check: ComposeTestCheckFunc(
					testExtractResourceAttr("random_password.test", "result", &result2),
					testCheckAttributeValuesEqual(&result1, &result2),
				),
			},
		},
	})
}

//nolint:paralleltest // Can't use t.Parallel with t.Setenv
func TestTest_TestStep_ProviderFactories_Import_External_WithPersistMatch_WithPersistWorkingDir(t *testing.T) {
	var result1, result2 string

	t.Setenv(plugintest.EnvTfAccPersistWorkingDir, "1")
	workingDir := t.TempDir()

	testSteps := []TestStep{
		{
			Config:             `resource "random_password" "test" { length = 12 }`,
			ResourceName:       "random_password.test",
			ImportState:        true,
			ImportStateId:      "Z=:cbrJE?Ltg",
			ImportStatePersist: true,
			ImportStateCheck: composeImportStateCheck(
				testExtractResourceAttrInstanceState("none", "result", &result1),
			),
		},
		{
			Config: `resource "random_password" "test" { length = 12 }`,
			Check: ComposeTestCheckFunc(
				testExtractResourceAttr("random_password.test", "result", &result2),
				testCheckAttributeValuesEqual(&result1, &result2),
			),
		},
	}

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		WorkingDir: workingDir,
		Steps:      testSteps,
	})

	for testStepIndex := range testSteps {
		dir := filepath.Join(workingDir, fmt.Sprintf("step_%s", strconv.Itoa(testStepIndex+1)))

		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			t.Errorf("cannot read dir: %s", dir)
		}

		var workingDirName string

		// Relies upon convention of a directory being created that is prefixed "work".
		for _, dirEntry := range dirEntries {
			if strings.HasPrefix(dirEntry.Name(), "work") && dirEntry.IsDir() {
				workingDirName = filepath.Join(dir, dirEntry.Name())
				break
			}
		}

		configPlanStateFiles := []string{
			"terraform_plugin_test.tf",
			"terraform.tfstate",
			"tfplan",
		}

		for _, file := range configPlanStateFiles {
			// Skip verifying plan for first test step as there is no plan file if the
			// resource does not already exist.
			if testStepIndex == 0 && file == "tfplan" {
				break
			}
			_, err = os.Stat(filepath.Join(workingDirName, file))
			if err != nil {
				t.Errorf("cannot stat %s in %s: %s", file, workingDirName, err)
			}
		}
	}
}

func TestTest_TestStep_ProviderFactories_Import_External_WithoutPersistNonMatch(t *testing.T) {
	var result1, result2 string

	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config:             `resource "random_password" "test" { length = 12 }`,
				ResourceName:       "random_password.test",
				ImportState:        true,
				ImportStateId:      "Z=:cbrJE?Ltg",
				ImportStatePersist: false,
				ImportStateCheck: composeImportStateCheck(
					testExtractResourceAttrInstanceState("none", "result", &result1),
				),
			},
			{
				Config: `resource "random_password" "test" { length = 12 }`,
				Check: ComposeTestCheckFunc(
					testExtractResourceAttr("random_password.test", "result", &result2),
					testCheckAttributeValuesDiffer(&result1, &result2),
				),
			},
		},
	})
}

//nolint:paralleltest // Can't use t.Parallel with t.Setenv
func TestTest_TestStep_ProviderFactories_Import_External_WithoutPersistNonMatch_WithPersistWorkingDir(t *testing.T) {
	var result1, result2 string

	t.Setenv(plugintest.EnvTfAccPersistWorkingDir, "1")
	workingDir := t.TempDir()

	testSteps := []TestStep{
		{
			Config:             `resource "random_password" "test" { length = 12 }`,
			ResourceName:       "random_password.test",
			ImportState:        true,
			ImportStateId:      "Z=:cbrJE?Ltg",
			ImportStatePersist: false,
			ImportStateCheck: composeImportStateCheck(
				testExtractResourceAttrInstanceState("none", "result", &result1),
			),
		},
		{
			Config: `resource "random_password" "test" { length = 12 }`,
			Check: ComposeTestCheckFunc(
				testExtractResourceAttr("random_password.test", "result", &result2),
				testCheckAttributeValuesDiffer(&result1, &result2),
			),
		},
	}

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		WorkingDir: workingDir,
		Steps:      testSteps,
	})

	for testStepIndex := range testSteps {
		dir := filepath.Join(workingDir, fmt.Sprintf("step_%s", strconv.Itoa(testStepIndex+1)))

		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			t.Errorf("cannot read dir: %s", dir)
		}

		var workingDirName string

		// Relies upon convention of a directory being created that is prefixed "work".
		for _, dirEntry := range dirEntries {
			if strings.HasPrefix(dirEntry.Name(), "work") && dirEntry.IsDir() {
				workingDirName = filepath.Join(dir, dirEntry.Name())
				break
			}
		}

		configPlanStateFiles := []string{
			"terraform_plugin_test.tf",
			"terraform.tfstate",
			"tfplan",
		}

		for _, file := range configPlanStateFiles {
			// Skip verifying state and plan for first test step as ImportStatePersist is
			// false so the state is not persisted and there is no plan file if the
			// resource does not already exist.
			if testStepIndex == 0 && (file == "terraform.tfstate" || file == "tfplan") {
				break
			}
			_, err = os.Stat(filepath.Join(workingDirName, file))
			if err != nil {
				t.Errorf("cannot stat %s in %s: %s", file, workingDirName, err)
			}
		}
	}
}

func TestTest_TestStep_ProviderFactories_Refresh_Inline(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_password": {
							CreateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
								d.SetId("id")
								err := d.Set("min_special", 10)
								if err != nil {
									panic(err)
								}
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								err := d.Set("min_special", 2)
								if err != nil {
									panic(err)
								}
								return nil
							},
							Schema: map[string]*schema.Schema{
								"min_special": {
									Computed: true,
									Type:     schema.TypeInt,
								},

								"id": {
									Computed: true,
									Type:     schema.TypeString,
								},
							},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_password" "test" { }`,
				Check:  TestCheckResourceAttr("random_password.test", "min_special", "10"),
			},
			{
				RefreshState: true,
				Check:        TestCheckResourceAttr("random_password.test", "min_special", "2"),
			},
			{
				Config: `resource "random_password" "test" { }`,
				Check:  TestCheckResourceAttr("random_password.test", "min_special", "2"),
			},
		},
	})
}

//nolint:paralleltest // Can't use t.Parallel with t.Setenv
func TestTest_TestStep_ProviderFactories_CopyWorkingDir_EachTestStep(t *testing.T) {
	t.Setenv(plugintest.EnvTfAccPersistWorkingDir, "1")
	workingDir := t.TempDir()

	testSteps := []TestStep{
		{
			Config: `resource "random_password" "test" { }`,
		},
		{
			Config: `resource "random_password" "test" { }`,
		},
	}

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_password": {
							CreateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
								d.SetId("id")
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{
								"id": {
									Computed: true,
									Type:     schema.TypeString,
								},
							},
						},
					},
				}, nil
			},
		},
		WorkingDir: workingDir,
		Steps:      testSteps,
	})

	for k := range testSteps {
		dir := filepath.Join(workingDir, fmt.Sprintf("step_%s", strconv.Itoa(k+1)))

		_, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("cannot read dir: %s", dir)
		}
	}
}

func TestTest_TestStep_ProviderFactories_RefreshWithPlanModifier_Inline(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_password": {
							CustomizeDiff: customdiff.All(
								func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
									special, ok := d.Get("special").(bool)
									if !ok {
										return fmt.Errorf("unexpected type %T for 'special' key", d.Get("special"))
									}

									if special == true {
										err := d.SetNew("special", false)
										if err != nil {
											panic(err)
										}
									}
									return nil
								},
							),
							CreateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
								d.SetId("id")
								err := d.Set("special", false)
								if err != nil {
									panic(err)
								}
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								t := getTimeForTest()
								if t.After(time.Now().Add(time.Hour * 1)) {
									err := d.Set("special", true)
									if err != nil {
										panic(err)
									}
								}
								return nil
							},
							Schema: map[string]*schema.Schema{
								"special": {
									Computed: true,
									Type:     schema.TypeBool,
									ForceNew: true,
								},

								"id": {
									Computed: true,
									Type:     schema.TypeString,
								},
							},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_password" "test" { }`,
				Check:  TestCheckResourceAttr("random_password.test", "special", "false"),
			},
			{
				PreConfig:          setTimeForTest(time.Now().Add(time.Hour * 2)),
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check:              TestCheckResourceAttr("random_password.test", "special", "true"),
			},
			{
				PreConfig: setTimeForTest(time.Now()),
				Config:    `resource "random_password" "test" { }`,
				Check:     TestCheckResourceAttr("random_password.test", "special", "false"),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Import_Inline_With_Data_Source(t *testing.T) {
	var id string

	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"http": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					DataSourcesMap: map[string]*schema.Resource{
						"http": {
							ReadContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
								url, ok := d.Get("url").(string)
								if !ok {
									return diag.Errorf("unexpected type %T for 'url' key", d.Get("url"))
								}

								responseHeaders := map[string]string{
									"headerOne":   "one",
									"headerTwo":   "two",
									"headerThree": "three",
									"headerFour":  "four",
								}
								if err := d.Set("response_headers", responseHeaders); err != nil {
									return append(diags, diag.Errorf("Error setting HTTP response headers: %s", err)...)
								}

								d.SetId(url)

								return diags
							},
							Schema: map[string]*schema.Schema{
								"url": {
									Type:     schema.TypeString,
									Required: true,
								},
								"response_headers": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				}, nil
			},
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_string": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId("none")
								err := d.Set("length", 4)
								if err != nil {
									panic(err)
								}
								err = d.Set("result", "none")
								if err != nil {
									panic(err)
								}
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{
								"length": {
									Required: true,
									ForceNew: true,
									Type:     schema.TypeInt,
								},
								"result": {
									Type:      schema.TypeString,
									Computed:  true,
									Sensitive: true,
								},

								"id": {
									Computed: true,
									Type:     schema.TypeString,
								},
							},
							Importer: &schema.ResourceImporter{
								StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
									val := d.Id()

									d.SetId(val)

									err := d.Set("result", val)
									if err != nil {
										panic(err)
									}

									err = d.Set("length", len(val))
									if err != nil {
										panic(err)
									}

									return []*schema.ResourceData{d}, nil
								},
							},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				Config: `data "http" "example" {
							url = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
						}

						resource "random_string" "example" {
							length = length(data.http.example.response_headers)
						}`,
				Check: extractResourceAttr("random_string.example", "id", &id),
			},
			{
				Config: `data "http" "example" {
							url = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
						}

						resource "random_string" "example" {
							length = length(data.http.example.response_headers)
						}`,
				ResourceName: "random_string.example",
				ImportState:  true,
				ImportStateCheck: composeImportStateCheck(
					testCheckResourceAttrInstanceState(&id, "length", "4"),
				),
				ImportStateVerify: true,
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_Import_External_With_Data_Source(t *testing.T) {
	var id string

	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"null": {
				Source: "registry.terraform.io/hashicorp/null",
			},
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config: `
					data "null_data_source" "values" {
						inputs = {
							length = 12
						}
					}

					resource "random_string" "example" {
						length = data.null_data_source.values.outputs["length"]
					}
				`,
				Check: extractResourceAttr("random_string.example", "id", &id),
			},
			{
				ResourceName:      "random_string.example",
				ImportState:       true,
				ImportStateCheck:  testCheckResourceAttrInstanceState(&id, "length", "12"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestTest_ConfigDirectory_StaticDirectory(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.5.1`),
				Check:           TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_StaticDirectory_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.5.1_vars`),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_StaticDirectory_VarsMissing(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.5.1_vars`),
				Check:           TestCheckResourceAttrSet("random_password.test", "id"),
				ExpectError:     regexp.MustCompile(`.*Error: No value for required variable`)},
		},
	})
}

func TestTest_ConfigDirectory_TestNameDirectory(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				Check:           TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestNameDirectory_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestStepDirectory(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check:           TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_StepNotHardcoded uses a multistep test
// to prove that the test step number is not hardcoded and to show that the
// configuration files that are copied from the test step directory in test step 1
// are removed prior to running test step 2.
func TestTest_ConfigDirectory_TestStepDirectory_StepNotHardcoded(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check:           TestCheckResourceAttrPtr("random_password.test", "length", teststep.Pointer("9")),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestStepDirectory_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_StaticDirectory_MultipleFiles(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.5.1_multiple_files`),
				Check:           TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_StaticDirectory_MultipleFiles_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.5.1_multiple_files_vars`),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestNameDirectory_MultipleFiles(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				Check:           TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestNameDirectory_MultipleFiles_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestStepDirectory_MultipleFiles(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check:           TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_MultipleFiles_StepNotHardcoded uses a
// multistep test to prove that the test step number is not hardcoded, and to show
// that the configuration files that are copied from the test step directory in test
// step 1 are removed prior to running test step 2.
func TestTest_ConfigDirectory_TestStepDirectory_MultipleFiles_StepNotHardcoded(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check:           TestCheckResourceAttrPtr("random_password.test", "length", teststep.Pointer("9")),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestStepDirectory_MultipleFiles_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_MultipleFiles_Vars_StepNotHardcoded uses a
// multistep test to prove that the test step number is not hardcoded, and to show
// that the configuration files that are copied from the test step directory in test
// step 1 are removed prior to running test step 2.
func TestTest_ConfigDirectory_TestStepDirectory_MultipleFiles_Vars_StepNotHardcoded(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(9),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrPtr("random_password.test", "length", teststep.Pointer("9")),
			},
		},
	})
}

// TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.2.0`),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.2.0_vars`),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_MultipleFiles uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_MultipleFiles(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.2.0_multiple_files`),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_MultipleFiles_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_MultipleFiles_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_password_3.2.0_multiple_files_vars`),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_MultipleFiles uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_MultipleFiles(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_MultipleFiles_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_MultipleFiles_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_MultipleFiles uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_MultipleFiles(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_MultipleFiles_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_MultipleFiles_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigDirectory_StaticDirectory(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_id": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId(time.Now().String())
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`testdata/fixtures/random_id`),
				Check:           TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigDirectory_TestNameDirectory(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_id": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId(time.Now().String())
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				Check:           TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigDirectory_TestStepDirectory(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_id": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId(time.Now().String())
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check:           TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func TestTest_ConfigFile_StaticFile(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.StaticFile(`testdata/fixtures/random_password_3.5.1/random.tf`),
				Check:      TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigFile_StaticFile_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.StaticFile(`testdata/fixtures/random_password_3.5.1_vars_single_file/random.tf`),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigFile_StaticFile_VarsMissing(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile:  config.StaticFile(`testdata/fixtures/random_password_3.5.1_vars_single_file/random.tf`),
				Check:       TestCheckResourceAttrSet("random_password.test", "id"),
				ExpectError: regexp.MustCompile(`.*Error: No value for required variable`)},
		},
	})
}

func TestTest_ConfigFile_TestNameFile(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.TestNameFile("random.tf"),
				Check:      TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigFile_TestNameFile_Vars(t *testing.T) {
	t.Setenv("TF_ACC_LOG", "TRACE")
	t.Setenv("TF_ACC_LOG_PATH", "./config_test_out.log")

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.TestNameFile("random.tf"),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigFile_TestStepFile(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.TestStepFile("random.tf"),
				Check:      TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigFile_TestStepFile_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.TestStepFile("random.tf"),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				Check: TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

// TestTest_ConfigFile_StaticFile_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigFile_StaticFile_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile:  config.StaticFile(`testdata/fixtures/random_password_3.2.0/random.tf`),
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigFile_StaticFile_AttributeDoesNotExist_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigFile_StaticFile_AttributeDoesNotExist_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.StaticFile(`testdata/fixtures/random_password_3.2.0_vars_single_file/random.tf`),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigFile_TestNameFile_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigFile_TestNameFile_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile:  config.TestNameFile("random.tf"),
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigFile_TestNameFile_AttributeDoesNotExist_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigFile_TestNameFile_AttributeDoesNotExist_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.TestNameFile("random.tf"),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigFile_TestStepFile_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigFile_TestStepFile_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile:  config.TestStepFile("random.tf"),
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigFile_TestStepFile_AttributeDoesNotExist_Vars uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password  This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigFile_TestStepFile_AttributeDoesNotExist_Vars(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		Steps: []TestStep{
			{
				ConfigFile: config.TestStepFile("random.tf"),
				ConfigVariables: config.Variables{
					"length":  config.IntegerVariable(8),
					"numeric": config.BoolVariable(false),
				},
				ExpectError: regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigFile_StaticFile(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_id": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId(time.Now().String())
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				ConfigFile: config.StaticFile(`testdata/fixtures/random_id/random.tf`),
				Check:      TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigFile_TestNameFile(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_id": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId(time.Now().String())
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				ConfigFile: config.TestNameFile("random.tf"),
				Check:      TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigFile_TestStepFile(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"random": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{
						"random_id": {
							CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
								d.SetId(time.Now().String())
								return nil
							},
							DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
								return nil
							},
							Schema: map[string]*schema.Schema{},
						},
					},
				}, nil
			},
		},
		Steps: []TestStep{
			{
				ConfigFile: config.TestStepFile("random.tf"),
				Check:      TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func setTimeForTest(t time.Time) func() {
	return func() {
		getTimeForTest = func() time.Time {
			return t
		}
	}
}

var getTimeForTest = func() time.Time {
	return time.Now()
}

func composeImportStateCheck(fs ...ImportStateCheckFunc) ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		for i, f := range fs {
			if err := f(s); err != nil {
				return fmt.Errorf("check %d/%d error: %s", i+1, len(fs), err)
			}
		}

		return nil
	}
}

//nolint:unparam // Generic test function
func testExtractResourceAttrInstanceState(id, attributeName string, attributeValue *string) ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != id {
				continue
			}

			if attrVal, ok := v.Attributes[attributeName]; ok {
				*attributeValue = attrVal

				return nil
			}
		}

		return fmt.Errorf("attribute %s not found in instance state", attributeName)
	}
}

func testCheckResourceAttrInstanceState(id *string, attributeName, attributeValue string) ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != *id {
				continue
			}

			if attrVal, ok := v.Attributes[attributeName]; ok {
				if attrVal != attributeValue {
					return fmt.Errorf("expected: %s got: %s", attributeValue, attrVal)
				}

				return nil
			}
		}

		return fmt.Errorf("attribute %s not found in instance state", attributeName)
	}
}

//nolint:unparam // Generic test function
func testExtractResourceAttr(resourceName string, attributeName string, attributeValue *string) TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource name %s not found in state", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]

		if !ok {
			return fmt.Errorf("attribute %s not found in resource %s state", attributeName, resourceName)
		}

		*attributeValue = attrValue

		return nil
	}
}

func testCheckAttributeValuesEqual(i *string, j *string) TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) != testStringValue(j) {
			return fmt.Errorf("attribute values are different, got %s and %s", testStringValue(i), testStringValue(j))
		}

		return nil
	}
}

func testCheckAttributeValuesDiffer(i *string, j *string) TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) == testStringValue(j) {
			return fmt.Errorf("attribute values are the same")
		}

		return nil
	}
}

func testStringValue(sPtr *string) string {
	if sPtr == nil {
		return ""
	}

	return *sPtr
}
