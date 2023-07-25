## 1.4.0 (July 24, 2023)

FEATURES:

* tfjsonpath: Introduced new `tfjsonpath` package which contains methods that allow traversal of Terraform JSON data ([#154](https://github.com/hashicorp/terraform-plugin-testing/issues/154))
* plancheck: Added `ExpectUnknownValue` built-in plan check, which asserts that a given attribute has an unknown value ([#154](https://github.com/hashicorp/terraform-plugin-testing/issues/154))
* plancheck: Added `ExpectSensitiveValue` built-in plan check, which asserts that a given attribute has a sensitive value ([#154](https://github.com/hashicorp/terraform-plugin-testing/issues/154))

## 1.3.0 (June 13, 2023)

FEATURES:

* tfversion: Introduced new `tfversion` package with interface and built-in Terraform version check functionality ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `SkipAbove` built-in version check, which skips the test if the Terraform CLI version is above the given maximum. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `SkipBelow` built-in version check, which skips the test if the Terraform CLI version is below the given minimum. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `SkipBetween` built-in version check, which skips the test if the Terraform CLI version is between the given minimum (inclusive) and maximum (exclusive). ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `SkipIf` built-in version check, which skips the test if the Terraform CLI version matches the given version. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `RequireAbove` built-in version check, which fails the test if the Terraform CLI version is below the given maximum. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `RequireBelow` built-in version check, which fails the test if the Terraform CLI version is above the given minimum. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `RequireBetween` built-in version check, fails the test if the Terraform CLI version is outside the given minimum (exclusive) and maximum (inclusive). ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `RequireNot` built-in version check, which fails the test if the Terraform CLI version matches the given version. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `Any` built-in version check, which fails the test if none of the given sub-checks return a nil error and empty skip message. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))
* tfversion: Added `All` built-in version check, which fails or skips the test if any of the given sub-checks return a non-nil error or non-empty skip message. ([#128](https://github.com/hashicorp/terraform-plugin-testing/issues/128))

BUG FIXES:

* helper/resource: Fix path used when persisting working directory ([#113](https://github.com/hashicorp/terraform-plugin-testing/issues/113))

## 1.2.0 (March 22, 2023)

NOTES:

* This Go module has been updated to Go 1.19 per the [Go support policy](https://golang.org/doc/devel/release.html#policy). Any consumers building on earlier Go versions may experience errors. ([#91](https://github.com/hashicorp/terraform-plugin-testing/issues/91))
* helper/resource: Deprecated `PrefixedUniqueId()` and `UniqueId()`. Use the `github.com/hashicorp/terraform-plugin-sdk/v2/helper/id` package instead. ([#96](https://github.com/hashicorp/terraform-plugin-testing/issues/96))
* helper/resource: Deprecated `RetryContext()`, `StateChangeConf`, and associated `*Error` types. Use the `github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry` package instead. ([#96](https://github.com/hashicorp/terraform-plugin-testing/issues/96))
* helper/resource: Deprecated Terraform module-based `TestCheckFunc`, such as `TestCheckModuleResourceAttr`. Provider testing should always be possible within the root module of a Terraform configuration. Terraform module testing should be performed with Terraform core functionality or using tooling outside this Go module. ([#109](https://github.com/hashicorp/terraform-plugin-testing/issues/109))

FEATURES:

* plancheck: Introduced new `plancheck` package with interface and built-in plan check functionality ([#63](https://github.com/hashicorp/terraform-plugin-testing/issues/63))
* plancheck: Added `ExpectResourceAction` built-in plan check, which asserts that a given resource will have a specific resource change type in the plan ([#63](https://github.com/hashicorp/terraform-plugin-testing/issues/63))
* plancheck: Added `ExpectEmptyPlan` built-in plan check, which asserts that there are no resource changes in the plan ([#63](https://github.com/hashicorp/terraform-plugin-testing/issues/63))
* plancheck: Added `ExpectNonEmptyPlan` built-in plan check, which asserts that there is at least one resource change in the plan ([#63](https://github.com/hashicorp/terraform-plugin-testing/issues/63))

ENHANCEMENTS:

* helper/resource: Added plan check functionality to config and refresh modes with new fields `TestStep.ConfigPlanChecks` and `TestStep.RefreshPlanChecks` ([#63](https://github.com/hashicorp/terraform-plugin-testing/issues/63))

## 1.1.0 (February 06, 2023)

FEATURES:

* helper/resource: Added `TF_ACC_PERSIST_WORKING_DIR` environment variable to allow persisting of Terraform files generated during each test step ([#18](https://github.com/hashicorp/terraform-plugin-testing/issues/18))
* helper/resource: Added `TestCase` type `WorkingDir` field to allow specifying the base directory where testing files used by the testing module are generated ([#18](https://github.com/hashicorp/terraform-plugin-testing/issues/18))

## 1.0.0 (January 10, 2023)

NOTES:

* Same testing functionality as that of terraform-plugin-sdk v2.24.1, repacked in a standalone repository ([#24](https://github.com/hashicorp/terraform-plugin-testing/issues/24))

