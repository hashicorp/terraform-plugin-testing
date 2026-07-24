// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	testinginterface "github.com/mitchellh/go-testing-interface"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

type importStateVerifyRecordingT struct {
	testinginterface.RuntimeT
	fatalfMessages []string
}

func (t *importStateVerifyRecordingT) Fatalf(format string, args ...interface{}) {
	t.fatalfMessages = append(t.fatalfMessages, fmt.Sprintf(format, args...))
	t.RuntimeT.FailNow()
}

func TestImportCommand_ImportStateVerifyFailureIncludesStepNumber(t *testing.T) { //nolint:paralleltest
	testT := &importStateVerifyRecordingT{}

	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(testT, r.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories
			},
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
					Resources: map[string]testprovider.Resource{
						"examplecloud_thing": {
							CreateResponse: &resource.CreateResponse{
								NewState: tftypes.NewValue(
									tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"id": tftypes.String,
										},
									},
									map[string]tftypes.Value{
										"id": tftypes.NewValue(tftypes.String, "resource-test"),
									},
								),
							},
							ImportStateResponse: &resource.ImportStateResponse{
								State: tftypes.NewValue(
									tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"id": tftypes.String,
										},
									},
									map[string]tftypes.Value{
										"id": tftypes.NewValue(tftypes.String, "imported-resource-test"),
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
			Steps: []r.TestStep{
				{
					Config: `resource "examplecloud_thing" "test" {}`,
				},
				{
					ResourceName:      "examplecloud_thing.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	expected := "Step 2/2 error running import: Failed state verification, resource with ID imported-resource-test not found"
	actual := strings.Join(testT.fatalfMessages, "\n")

	if !strings.Contains(actual, expected) {
		t.Fatalf("expected fatal output to contain %q, got:\n%s", expected, actual)
	}
}
