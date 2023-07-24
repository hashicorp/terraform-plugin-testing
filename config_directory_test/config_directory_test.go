// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config_directory_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTest_ConfigDirectory_StaticDirectory(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`fixtures/random_password_3.5.1`),
				Check:           resource.TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestNameDirectory(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(t),
				Check:           resource.TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestStepDirectory(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(t),
				Check:           resource.TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_StaticDirectory_MultipleFiles(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`fixtures/random_password_3.5.1_multiple_files`),
				Check:           resource.TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestNameDirectory_MultipleFiles(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(t),
				Check:           resource.TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

func TestTest_ConfigDirectory_TestStepDirectory_MultipleFiles(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(t),
				Check:           resource.TestCheckResourceAttrSet("random_password.test", "id"),
			},
		},
	})
}

// TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password resource. This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`fixtures/random_password_3.2.0`),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password resource. This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(t),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password resource. This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(t),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_MultipleFiles uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password resource. This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_StaticDirectory_AttributeDoesNotExist_MultipleFiles(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`fixtures/random_password_3.2.0_multiple_files`),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_MultipleFiles uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password resource. This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestNameDirectory_AttributeDoesNotExist_MultipleFiles(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(t),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

// TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_MultipleFiles uses Terraform
// configuration specifying a "numeric" attribute that was introduced in v3.3.0 of the
// random provider password resource. This test confirms that the TestCase ExternalProviders
// is being used when ConfigDirectory is set.
func TestTest_ConfigDirectory_TestStepDirectory_AttributeDoesNotExist_MultipleFiles(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(t),
				ExpectError:     regexp.MustCompile(`.*An argument named "numeric" is not expected here.`),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigDirectory_StaticDirectory(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory(`fixtures/random_id`),
				Check:           resource.TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigDirectory_TestNameDirectory(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(t),
				Check:           resource.TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}

func TestTest_TestStep_ProviderFactories_ConfigDirectory_TestStepDirectory(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(t),
				Check:           resource.TestCheckResourceAttrSet("random_id.test", "id"),
			},
		},
	})
}
