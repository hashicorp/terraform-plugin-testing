// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestTest_TestStep_ExpectError_NewConfig(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
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
											Required: true,
										},
									},
								},
							},
						},
						ValidateConfigResponse: &resource.ValidateConfigResponse{
							Diagnostics: []*tfprotov6.Diagnostic{
								{
									Severity: tfprotov6.DiagnosticSeverityError,
									Summary:  "Invalid Attribute Value",
									Detail:   "Diagnostic details",
								},
							},
						},
					},
				},
			}),
		},
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {
					id = "invalid-value"
				}`,
				ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value`),
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_EmptyPlanError(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		ExternalProviders: map[string]ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []TestStep{
			{
				Config:             `resource "terraform_data" "test" {}`,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("Expected a non-empty plan, but got an empty refresh plan"),
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_PreRefresh_ResourceChanges(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		ExternalProviders: map[string]ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []TestStep{
			{
				Config: `resource "terraform_data" "test" {
					# Never recommended for real world configurations, but tests
					# the intended behavior.
					input = timestamp()
				}`,
				ConfigPlanChecks: ConfigPlanChecks{
					// Verification of that the behavior is being caught pre
					// refresh. We want to ensure ExpectNonEmptyPlan allows test
					// to pass if pre refresh also has changes.
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("terraform_data.test", plancheck.ResourceActionUpdate),
					},
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_PostRefresh_OutputChanges(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipAbove(tfversion.Version0_14_0), // outputs before 0.14 always show as created
		},
		// Avoid our own validation that requires at least one provider config.
		ExternalProviders: map[string]ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []TestStep{
			{
				Config:             `output "test" { value = timestamp() }`,
				ExpectNonEmptyPlan: false, // compatibility compromise for 0.12 and 0.13
			},
		},
	})

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_14_0), // outputs before 0.14 always show as created
		},
		// Avoid our own validation that requires at least one provider config.
		ExternalProviders: map[string]ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []TestStep{
			{
				Config:             `output "test" { value = timestamp() }`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_PostRefresh_ResourceChanges(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"), // intentionally same
								},
							),
						},
						ReadResponse: &resource.ReadResponse{
							NewState: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id": tftypes.NewValue(tftypes.String, "not-test"), // intentionally different
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
											Required: true,
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {
					# Post create refresh intentionally changes configured value
					# which is an errant resource implementation. Create should
					# account for the correct post creation state, preventing an
					# immediate difference next Terraform run for practitioners.
					# This errant resource behavior verifies the expected
					# behavior of ExpectNonEmptyPlan for post refresh planning.
					id = "test"
				}`,
				ConfigPlanChecks: ConfigPlanChecks{
					// Verification of that the behavior is being caught post
					// refresh. We want to ensure ExpectNonEmptyPlan is being
					// triggered after the pre refresh plan shows no changes.
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("test_resource.test", plancheck.ResourceActionNoop),
					},
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func Test_NonEmptyPlan_PreRefresh_Error(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		ExternalProviders: map[string]ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []TestStep{
			{
				Config: `resource "terraform_data" "test" {
					# Never recommended for real world configurations, but tests
					# the intended behavior.
					input = timestamp()
				}`,
				ConfigPlanChecks: ConfigPlanChecks{
					// Verification of that the behavior is being caught pre
					// refresh.
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("terraform_data.test", plancheck.ResourceActionUpdate),
					},
				},
				ExpectNonEmptyPlan: false, // intentional
				ExpectError:        regexp.MustCompile("After applying this test step, the non-refresh plan was not empty."),
			},
		},
	})
}

func Test_NonEmptyPlan_PostRefresh_Error(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"), // intentionally same
								},
							),
						},
						ReadResponse: &resource.ReadResponse{
							NewState: tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"id": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"id": tftypes.NewValue(tftypes.String, "not-test"), // intentionally different
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
											Required: true,
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {
					# Post create refresh intentionally changes configured value
					# which is an errant resource implementation. Create should
					# account for the correct post creation state, preventing an
					# immediate difference next Terraform run for practitioners.
					# This errant resource behavior verifies the expected
					# behavior of ExpectNonEmptyPlan for post refresh planning.
					id = "test"
				}`,
				ConfigPlanChecks: ConfigPlanChecks{
					// Verification of that the behavior is being caught post
					// refresh.
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("test_resource.test", plancheck.ResourceActionNoop),
					},
				},
				ExpectNonEmptyPlan: false, // intentional
				ExpectError:        regexp.MustCompile("After applying this test step, the refresh plan was not empty."),
			},
		},
	})
}

func Test_ConfigPlanChecks_PreApply_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigPlanChecks.PreApply spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigPlanChecks.PreApply spy2 to be called at least once")
	}
}

func Test_ConfigPlanChecks_PreApply_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &planCheckSpy{
		err: errors.New("spy3 check failed"),
	}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						spy1,
						spy2,
						spy3,
					},
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 check failed)\n.*?(spy3 check failed)`),
			},
		},
	})
}

func Test_ConfigPlanChecks_PostApplyPreRefresh_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigPlanChecks.PostApplyPreRefresh spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigPlanChecks.PostApplyPreRefresh spy2 to be called at least once")
	}
}

func Test_ConfigPlanChecks_PostApplyPreRefresh_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &planCheckSpy{
		err: errors.New("spy3 check failed"),
	}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
						spy3,
					},
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 check failed)\n.*?(spy3 check failed)`),
			},
		},
	})
}

func Test_ConfigPlanChecks_PostApplyPostRefresh_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigPlanChecks.PostApplyPostRefresh spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigPlanChecks.PostApplyPostRefresh spy2 to be called at least once")
	}
}

func Test_ConfigPlanChecks_PostApplyPostRefresh_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &planCheckSpy{
		err: errors.New("spy3 check failed"),
	}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
						spy3,
					},
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 check failed)\n.*?(spy3 check failed)`),
			},
		},
	})
}

func Test_ConfigStateChecks_Called(t *testing.T) {
	t.Parallel()

	spy1 := &stateCheckSpy{}
	spy2 := &stateCheckSpy{}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					spy1,
					spy2,
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected ConfigStateChecks spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected ConfigStateChecks spy2 to be called at least once")
	}
}

func Test_ConfigStateChecks_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &stateCheckSpy{}
	spy2 := &stateCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &stateCheckSpy{
		err: errors.New("spy3 check failed"),
	}
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					spy1,
					spy2,
					spy3,
				},
				ExpectError: regexp.MustCompile(`.*?(spy2 check failed)\n.*?(spy3 check failed)`),
			},
		},
	})
}

func Test_PostApplyFunc_Called(t *testing.T) {
	t.Parallel()

	spy1 := &stateCheckSpy{}
	postFuncCalled := false
	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config: `resource "test_resource" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					spy1,
				},
				PostApplyFunc: func() {
					postFuncCalled = true
				},
			},
		},
	})

	if !postFuncCalled {
		t.Error("expected PostApplyFunc to be called at least once")
	}

	if !spy1.called {
		t.Error("expected ConfigStateChecks spy1 to be called at least once")
	}
}

// This regression test ensures that the combination of Config, Destroy, and Check never result in
// a "Saved Plan is Stale" error message, which occurs when the state serial does not match the plan.
//
// This can occur when the refresh that is only done to produce the shimmed "Check" state produces a new state serial.
// Running a fresh plan after refreshing solves that issue, which was introduced in: https://github.com/hashicorp/terraform-plugin-testing/pull/602
func Test_Destroy_Checks_Avoid_Stale_Plan(t *testing.T) {
	t.Parallel()

	UnitTest(t, TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
		},
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
									"id": tftypes.NewValue(tftypes.String, "test"),
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
		Steps: []TestStep{
			{
				Config:  `resource "test_resource" "test" {}`,
				Destroy: true,
				Check:   func(s *terraform.State) error { return nil },
			},
		},
	})
}
