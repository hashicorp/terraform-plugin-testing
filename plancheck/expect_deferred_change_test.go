// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_ExpectDeferredChange_Reason_Match(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_9_0),
			tfversion.SkipIfNotAlpha(),
		},
		AdditionalCLIOptions: &r.AdditionalCLIOptions{
			Plan:  r.PlanOptions{AllowDeferral: true},
			Apply: r.ApplyOptions{AllowDeferral: true},
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						PlanChangeFunc: func(ctx context.Context, req resource.PlanChangeRequest, resp *resource.PlanChangeResponse) {
							resp.Deferred = &tfprotov6.Deferred{
								Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
							}
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
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectDeferredChange("test_resource.test", plancheck.DeferredReasonResourceConfigUnknown),
					},
				},
			},
		},
	})
}

func Test_ExpectDeferredChange_Reason_NoMatch(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_9_0),
			tfversion.SkipIfNotAlpha(),
		},
		AdditionalCLIOptions: &r.AdditionalCLIOptions{
			Plan:  r.PlanOptions{AllowDeferral: true},
			Apply: r.ApplyOptions{AllowDeferral: true},
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						PlanChangeFunc: func(ctx context.Context, req resource.PlanChangeRequest, resp *resource.PlanChangeResponse) {
							resp.Deferred = &tfprotov6.Deferred{
								Reason: tfprotov6.DeferredReasonAbsentPrereq,
							}
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
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectDeferredChange("test_resource.test", plancheck.DeferredReasonProviderConfigUnknown),
					},
				},
				ExpectError: regexp.MustCompile(`expected "provider_config_unknown", got deferred reason: "absent_prereq"`),
			},
		},
	})
}

func Test_ExpectDeferredChange_NoDeferredChanges(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_9_0),
			tfversion.SkipIfNotAlpha(),
		},
		AdditionalCLIOptions: &r.AdditionalCLIOptions{
			Plan:  r.PlanOptions{AllowDeferral: true},
			Apply: r.ApplyOptions{AllowDeferral: true},
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
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
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectDeferredChange("test_resource.test", plancheck.DeferredReasonProviderConfigUnknown),
					},
				},
				ExpectError: regexp.MustCompile(`No deferred changes found for resource`),
			},
		},
	})
}
