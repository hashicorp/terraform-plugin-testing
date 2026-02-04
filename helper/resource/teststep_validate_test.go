// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestTestStepHasExternalProviders(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testStep TestStep
		expected bool
	}{
		"none": {
			testStep: TestStep{},
			expected: false,
		},
		"externalproviders": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.testStep.hasExternalProviders()

			if got != test.expected {
				t.Errorf("expected %t, got %t", test.expected, got)
			}
		})
	}
}

func TestTestStepHasProviders(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testStep TestStep
		expected bool
	}{
		"none": {
			testStep: TestStep{},
			expected: false,
		},
		"externalproviders": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
			},
			expected: true,
		},
		"protov5providerfactories": {
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			expected: true,
		},
		"protov6providerfactories": {
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			expected: true,
		},
		"providerfactories": {
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil, // does not need to be real
				},
			},
			expected: true,
		},
	}

	var stepIndex int

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := test.testStep.hasProviders(context.Background(), stepIndex, "TestTestStepHasProviders")

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if got != test.expected {
				t.Errorf("expected %t, got %t", test.expected, got)
			}
		})

		stepIndex++
	}
}

func TestTestStepValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testStep                TestStep
		testStepConfig          string
		testStepConfigDirectory string
		testStepConfigFile      string
		testStepValidateRequest testStepValidateRequest
		expectedError           error
	}{
		"config-and-importstate-and-refreshstate-missing": {
			testStep:                TestStep{},
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("TestStep missing Config or ConfigDirectory or ConfigFile or ImportState or RefreshState"),
		},
		"config-and-refreshstate-both-set": {
			testStep: TestStep{
				RefreshState: true,
			},
			testStepConfig: "# not empty",
			expectedError:  fmt.Errorf("TestStep cannot have Config or ConfigDirectory or ConfigFile and RefreshState"),
		},
		"config-directory-and-refreshstate-both-set": {
			testStep: TestStep{
				RefreshState: true,
			},
			testStepConfigDirectory: "# not empty",
			expectedError:           fmt.Errorf("TestStep cannot have Config or ConfigDirectory or ConfigFile and RefreshState"),
		},
		"config-file-and-refreshstate-both-set": {
			testStep: TestStep{
				RefreshState: true,
			},
			testStepConfigFile: "# not empty",
			expectedError:      fmt.Errorf("TestStep cannot have Config or ConfigDirectory or ConfigFile and RefreshState"),
		},
		"refreshstate-first-step": {
			testStep: TestStep{
				RefreshState: true,
			},
			testStepValidateRequest: testStepValidateRequest{
				StepNumber: 1,
			},
			expectedError: fmt.Errorf("TestStep cannot have RefreshState as first step"),
		},
		"importstate-and-refreshstate-both-true": {
			testStep: TestStep{
				ImportState:  true,
				RefreshState: true,
			},
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("TestStep cannot have ImportState and RefreshState in same step"),
		},
		"destroy-and-refreshstate-both-true": {
			testStep: TestStep{
				Destroy:      true,
				RefreshState: true,
			},
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("TestStep cannot have RefreshState and Destroy"),
		},
		"externalproviders-overlapping-providerfactories": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfig:          "# not empty",
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("TestStep provider \"test\" set in both ExternalProviders and ProviderFactories"),
		},
		"externalproviders-overlapping-providerfactories-config-directory": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigDirectory: "# not empty",
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("TestStep provider \"test\" set in both ExternalProviders and ProviderFactories"),
		},
		"externalproviders-overlapping-providerfactories-config-file": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigFile:      "# not empty",
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("TestStep provider \"test\" set in both ExternalProviders and ProviderFactories"),
		},
		"externalproviders-testcase-config-directory": {
			testStep:                TestStep{},
			testStepConfigDirectory: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasExternalProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified within the terraform configuration files when using TestStep.Config"),
		},
		"externalproviders-testcase-config-file": {
			testStep:           TestStep{},
			testStepConfigFile: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasExternalProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified within the terraform configuration files when using TestStep.Config"),
		},
		"externalproviders-teststep-config-directory": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
			},
			testStepConfigDirectory: "# not empty",
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("Providers must only be specified within the terraform configuration files when using TestStep.Config"),
		},
		"externalproviders-teststep-config-file": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
			},
			testStepConfigFile:      "# not empty",
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("Providers must only be specified within the terraform configuration files when using TestStep.Config"),
		},
		"externalproviders-testcase-providers": {
			testStep: TestStep{
				ExternalProviders: map[string]ExternalProvider{
					"test": {}, // does not need to be real
				},
			},
			testStepConfig: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"importstate-missing-resourcename": {
			testStep: TestStep{
				ImportState: true,
			},
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("TestStep ImportState must be specified with ImportStateId, ImportStateIdFunc, or ResourceName"),
		},
		// This test has been added to verify that providers can be defined
		// both within the TestStep.Config and at the TestCase level.
		// The regression was reported in
		// https://github.com/hashicorp/terraform-plugin-testing/issues/176
		"config-providers-testcase-providers": {
			testStep: TestStep{
				Config: "provider abc {",
			},
			testStepConfig: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
		},
		"protov5providerfactories-testcase-providers": {
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfig: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"protov5providerfactories-testcase-providers-config-directory": {
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigDirectory: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"protov5providerfactories-testcase-providers-config-file": {
			testStep: TestStep{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigFile: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"protov6providerfactories-testcase-providers": {
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfig: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"protov6providerfactories-testcase-providers-config-directory": {
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigDirectory: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"protov6providerfactories-testcase-providers-config-file": {
			testStep: TestStep{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigFile: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"providerfactories-testcase-providers": {
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfig: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"providerfactories-testcase-providers-config-directory": {
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigDirectory: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"providerfactories-testcase-providers-config-file": {
			testStep: TestStep{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": nil, // does not need to be real
				},
			},
			testStepConfigFile: "# not empty",
			testStepValidateRequest: testStepValidateRequest{
				TestCaseHasProviders: true,
			},
			expectedError: fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level"),
		},
		"configplanchecks-preapply-not-config-mode": {
			testStep: TestStep{
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{&planCheckSpy{}},
				},
				RefreshState: true,
			},
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep ConfigPlanChecks.PreApply must only be specified with Config"),
		},
		"configplanchecks-preapply-not-planonly": {
			testStep: TestStep{
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{&planCheckSpy{}},
				},
				PlanOnly: true,
			},
			testStepConfig:          "# not empty",
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep ConfigPlanChecks.PreApply cannot be run with PlanOnly"),
		},
		"configplanchecks-preapply-not-planonly-config-directory": {
			testStep: TestStep{
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{&planCheckSpy{}},
				},
				PlanOnly: true,
			},
			testStepConfigDirectory: "testdata/fixtures/random_id",
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep ConfigPlanChecks.PreApply cannot be run with PlanOnly"),
		},
		"configplanchecks-preapply-not-planonly-config-file": {
			testStep: TestStep{
				ConfigPlanChecks: ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{&planCheckSpy{}},
				},
				PlanOnly: true,
			},
			testStepConfigFile:      "testdata/fixtures/random_id/random.tf",
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep ConfigPlanChecks.PreApply cannot be run with PlanOnly"),
		},
		"configplanchecks-postapplyprerefresh-not-config-mode": {
			testStep: TestStep{
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{&planCheckSpy{}},
				},
				RefreshState: true,
			},
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep ConfigPlanChecks.PostApplyPreRefresh must only be specified with Config"),
		},
		"configplanchecks-postapplypostrefresh-not-config-mode": {
			testStep: TestStep{
				ConfigPlanChecks: ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{&planCheckSpy{}},
				},
				RefreshState: true,
			},
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep ConfigPlanChecks.PostApplyPostRefresh must only be specified with Config"),
		},
		"configstatechecks-not-config-mode": {
			testStep: TestStep{
				ConfigStateChecks: []statecheck.StateCheck{
					&stateCheckSpy{},
				},
				RefreshState: true,
			},
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep ConfigStateChecks must only be specified with Config"),
		},
		"refreshplanchecks-postrefresh-not-refresh-mode": {
			testStep: TestStep{
				RefreshPlanChecks: RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{&planCheckSpy{}},
				},
			},
			testStepConfig:          "# not empty",
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep RefreshPlanChecks.PostRefresh must only be specified with RefreshState"),
		},
		"refreshplanchecks-postrefresh-not-refresh-mode-config-directory": {
			testStep: TestStep{
				RefreshPlanChecks: RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{&planCheckSpy{}},
				},
			},
			testStepConfigDirectory: "testdata/fixtures/random_id",
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep RefreshPlanChecks.PostRefresh must only be specified with RefreshState"),
		},
		"refreshplanchecks-postrefresh-not-refresh-mode-config-file": {
			testStep: TestStep{
				RefreshPlanChecks: RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{&planCheckSpy{}},
				},
			},
			testStepConfigFile:      "testdata/fixtures/random_id/random.tf",
			testStepValidateRequest: testStepValidateRequest{TestCaseHasProviders: true},
			expectedError:           errors.New("TestStep RefreshPlanChecks.PostRefresh must only be specified with RefreshState"),
		},
		"state-store-mode-missing-config": {
			testStep: TestStep{
				StateStore: true,
			},
			testStepValidateRequest: testStepValidateRequest{},
			expectedError:           fmt.Errorf("TestStep missing Config or ConfigDirectory or ConfigFile or ImportState or RefreshState"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			configRequest := teststep.PrepareConfigurationRequest{
				Directory:             func(config.TestStepConfigRequest) string { return test.testStepConfigDirectory },
				File:                  func(config.TestStepConfigRequest) string { return test.testStepConfigFile },
				Raw:                   test.testStepConfig,
				TestStepConfigRequest: config.TestStepConfigRequest{},
			}.Exec()

			testStepConfig := teststep.Configuration(configRequest)

			testStepValidateRequest := test.testStepValidateRequest
			testStepValidateRequest.StepConfiguration = testStepConfig

			err := test.testStep.validate(context.Background(), testStepValidateRequest)

			if err != nil {
				if test.expectedError == nil {
					t.Fatalf("unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), test.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", test.expectedError, err)
				}
			}

			if err == nil && test.expectedError != nil {
				t.Errorf("expected error: %s", test.expectedError)
			}
		})
	}
}
