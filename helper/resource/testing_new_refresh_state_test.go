// Copyright (c) HashiCorp, Inc.
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
)

func Test_RefreshPlanChecks_PostRefresh_Called(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{}
	UnitTest(t, TestCase{
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
			},
			{
				RefreshState: true,
				RefreshPlanChecks: RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						spy1,
						spy2,
					},
				},
			},
		},
	})

	if !spy1.called {
		t.Error("expected RefreshPlanChecks.PostRefresh spy1 to be called at least once")
	}

	if !spy2.called {
		t.Error("expected RefreshPlanChecks.PostRefresh spy2 to be called at least once")
	}
}

func Test_RefreshPlanChecks_PostRefresh_Errors(t *testing.T) {
	t.Parallel()

	spy1 := &planCheckSpy{}
	spy2 := &planCheckSpy{
		err: errors.New("spy2 check failed"),
	}
	spy3 := &planCheckSpy{
		err: errors.New("spy3 check failed"),
	}
	UnitTest(t, TestCase{
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
			},
			{
				RefreshState: true,
				RefreshPlanChecks: RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
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
