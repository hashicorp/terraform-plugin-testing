## 1.12.0 (March 18, 2025)

NOTES:

* all: This Go module has been updated to Go 1.23 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.23 release notes](https://go.dev/doc/go1.23) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#454](https://github.com/hashicorp/terraform-plugin-testing/issues/454))

FEATURES:

* knownvalue: added function checks for custom validation of resource attribute or output values. ([#412](https://github.com/hashicorp/terraform-plugin-testing/issues/412))

ENHANCEMENTS:

* knownvalue: Updated the `ObjectExact` error message to report extra/missing attributes from the actual object. ([#451](https://github.com/hashicorp/terraform-plugin-testing/issues/451))
* plancheck: Improved the unknown value plan check error messages to include a known value if one exists. ([#450](https://github.com/hashicorp/terraform-plugin-testing/issues/450))

BUG FIXES:

* plancheck: Fixed bug with all unknown value plan checks where a valid path would return a "path not found" error. ([#450](https://github.com/hashicorp/terraform-plugin-testing/issues/450))

## 1.11.0 (November 19, 2024)

NOTES:

* all: This Go module has been updated to Go 1.22 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.22 release notes](https://go.dev/doc/go1.22) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#371](https://github.com/hashicorp/terraform-plugin-testing/issues/371))
* echoprovider: The `echoprovider` package is considered experimental and may be altered or removed in a subsequent release ([#389](https://github.com/hashicorp/terraform-plugin-testing/issues/389))

FEATURES:

* tfversion: Added `SkipIfNotAlpha` version check for testing experimental features of alpha Terraform builds. ([#388](https://github.com/hashicorp/terraform-plugin-testing/issues/388))
* echoprovider: Introduced new `echoprovider` package, which contains a v6 Terraform provider that can be used to test ephemeral resource data. ([#389](https://github.com/hashicorp/terraform-plugin-testing/issues/389))

## 1.10.0 (August 08, 2024)

NOTES:

* compare: The `compare` package is considered experimental and may be altered or removed in a subsequent release ([#330](https://github.com/hashicorp/terraform-plugin-testing/issues/330))
* statecheck: `CompareValue`, `CompareValueCollection`, and `CompareValuePairs` state checks are considered experimental and may be altered or removed in a subsequent release. ([#330](https://github.com/hashicorp/terraform-plugin-testing/issues/330))

FEATURES:

* compare: Introduced new `compare` package, which contains interfaces and implementations for value comparisons in state checks. ([#330](https://github.com/hashicorp/terraform-plugin-testing/issues/330))
* statecheck: Added `CompareValue` state check, which compares sequential values of the specified attribute at the given managed resource, or data source, using the supplied value comparer. ([#330](https://github.com/hashicorp/terraform-plugin-testing/issues/330))
* statecheck: Added `CompareValueCollection` state check, which compares each item in the specified collection (e.g., list, set) attribute, with the second specified attribute at the given managed resources, or data sources, using the supplied value comparer. ([#330](https://github.com/hashicorp/terraform-plugin-testing/issues/330))
* statecheck: Added `CompareValuePairs` state check, which compares the specified attributes at the given managed resources, or data sources, using the supplied value comparer. ([#330](https://github.com/hashicorp/terraform-plugin-testing/issues/330))

## 1.9.0 (July 09, 2024)

ENHANCEMENTS:

* knownvalue: Add `Int32Exact` check for int32 value testing. ([#356](https://github.com/hashicorp/terraform-plugin-testing/issues/356))
* knownvalue: Add `Float32Exact` check for float32 value testing. ([#356](https://github.com/hashicorp/terraform-plugin-testing/issues/356))

## 1.8.0 (May 17, 2024)

FEATURES:

* plancheck: Added  `ExpectDeferredChange` and `ExpectNoDeferredChanges` checks for experimental deferred action support. ([#331](https://github.com/hashicorp/terraform-plugin-testing/issues/331))
* tfversion: Added `SkipIfNotPrerelease` version check for testing experimental features of prerelease Terraform builds. ([#331](https://github.com/hashicorp/terraform-plugin-testing/issues/331))

ENHANCEMENTS:

* helper/acctest: Improve scope of IPv4/IPv6 random address generation in RandIpAddress() ([#305](https://github.com/hashicorp/terraform-plugin-testing/issues/305))
* knownvalue: Add `TupleExact`, `TuplePartial` and `TupleSizeExact` checks for dynamic value testing. ([#312](https://github.com/hashicorp/terraform-plugin-testing/issues/312))
* tfversion: Ensured Terraform CLI prerelease versions are considered semantically equal to patch versions in built-in checks to match the Terraform CLI versioning policy ([#303](https://github.com/hashicorp/terraform-plugin-testing/issues/303))
* helper/resource: Added `(TestCase).AdditionalCLIOptions` with `AllowDeferral` option for plan and apply commands. ([#331](https://github.com/hashicorp/terraform-plugin-testing/issues/331))

BUG FIXES:

* helper/resource: Fix panic in output state shimming when a tuple is present. ([#310](https://github.com/hashicorp/terraform-plugin-testing/issues/310))
* tfversion: Fixed `RequireBelow` ignoring equal versioning to fail a test ([#303](https://github.com/hashicorp/terraform-plugin-testing/issues/303))

## 1.7.0 (March 05, 2024)

NOTES:

* helper/resource: Error messages generated by the testing logic, which were updated for clarity in this release, are not protected by compatibility promises. While testing logic errors are usable in certain scenarios with `ErrorCheck` and `ExpectError` functionality, error messaging checks should be based on provider-controlled messaging or when appropriate to use other testing features such as `ExpectNonEmptyPlan` instead. ([#238](https://github.com/hashicorp/terraform-plugin-testing/issues/238))
* Numerical values in the plan are now represented as json.Number, not float64. Custom plan checks relying upon float64 representation may need altering ([#248](https://github.com/hashicorp/terraform-plugin-testing/issues/248))
* plancheck: Deprecated `ExpectNullOutputValue` and `ExpectNullOutputValueAtPath`. Use `ExpectKnownOutputValue` and `ExpectKnownOutputValueAtPath` with `knownvalue.Null` instead ([#275](https://github.com/hashicorp/terraform-plugin-testing/issues/275))
* plancheck: `ExpectKnownValue`, `ExpectKnownOutputValue` and `ExpectKnownOutputValueAtPath` plan checks are considered experimental and may be altered or removed in a subsequent release ([#276](https://github.com/hashicorp/terraform-plugin-testing/issues/276))
* statecheck: `ExpectKnownValue`, `ExpectKnownOutputValue` and `ExpectKnownOutputValueAtPath` state checks are considered experimental and may be altered or removed in a subsequent release ([#276](https://github.com/hashicorp/terraform-plugin-testing/issues/276))
* knownvalue: The `knownvalue` package is considered experimental and may be altered or removed in a subsequent release ([#276](https://github.com/hashicorp/terraform-plugin-testing/issues/276))
* all: This Go module has been updated to Go 1.21 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.21 release notes](https://go.dev/doc/go1.21) before upgrading. Any consumers building on earlier Go versions may experience errors ([#300](https://github.com/hashicorp/terraform-plugin-testing/issues/300))

FEATURES:

* plancheck: Added `ExpectKnownValue` plan check, which asserts that a given resource attribute has a defined type, and value ([#248](https://github.com/hashicorp/terraform-plugin-testing/issues/248))
* plancheck: Added `ExpectKnownOutputValue` plan check, which asserts that a given output value has a defined type, and value ([#248](https://github.com/hashicorp/terraform-plugin-testing/issues/248))
* plancheck: Added `ExpectKnownOutputValueAtPath` plan check, which asserts that a given output value at a specified path has a defined type, and value ([#248](https://github.com/hashicorp/terraform-plugin-testing/issues/248))
* knownvalue: Introduced new `knownvalue` package which contains types for working with plan checks and state checks ([#248](https://github.com/hashicorp/terraform-plugin-testing/issues/248))
* statecheck: Introduced new `statecheck` package with interface and built-in state check functionality ([#275](https://github.com/hashicorp/terraform-plugin-testing/issues/275))
* statecheck: Added `ExpectKnownValue` state check, which asserts that a given resource attribute has a defined type, and value ([#275](https://github.com/hashicorp/terraform-plugin-testing/issues/275))
* statecheck: Added `ExpectKnownOutputValue` state check, which asserts that a given output value has a defined type, and value ([#275](https://github.com/hashicorp/terraform-plugin-testing/issues/275))
* statecheck: Added `ExpectKnownOutputValueAtPath` plan check, which asserts that a given output value at a specified path has a defined type, and value ([#275](https://github.com/hashicorp/terraform-plugin-testing/issues/275))
* statecheck: Added `ExpectSensitiveValue` built-in state check, which asserts that a given attribute has a sensitive value ([#275](https://github.com/hashicorp/terraform-plugin-testing/issues/275))

BUG FIXES:

* helper/resource: Clarified error messaging from testing failures, especially when using `TestStep.PlanOnly: true` ([#238](https://github.com/hashicorp/terraform-plugin-testing/issues/238))
* helper/resource: Fix detection of provider block declaration in `Config`, `ConfigDirectory`, and `ConfigFile` ([#265](https://github.com/hashicorp/terraform-plugin-testing/issues/265))
* helper/resource: Fix detection of terraform block declaration in `Config`, `ConfigDirectory`, and `ConfigFile` ([#265](https://github.com/hashicorp/terraform-plugin-testing/issues/265))
* helper/resource: Fixed internal deferred test helpers to properly report file and line information in test failures. ([#292](https://github.com/hashicorp/terraform-plugin-testing/issues/292))

## 1.6.0 (December 04, 2023)

NOTES:

* all: This Go module has been updated to Go 1.20 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.20 release notes](https://go.dev/doc/go1.20) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#180](https://github.com/hashicorp/terraform-plugin-testing/issues/180))
* helper/resource: Configuration based `TestStep` now include post-apply plan checks for output changes in addition to resource changes. If this causes unexpected new test failures, most `output` configuration blocks can be likely be removed. Test steps involving resources and data sources should never need to use `output` configuration blocks as plan and state checks support working on resource and data source attributes values directly. ([#234](https://github.com/hashicorp/terraform-plugin-testing/issues/234))
* helper/resource: Implicit `terraform refresh` commands during each `TestStep` have been removed to fix plan check and performance issues, which can cause new test failures when testing schema changes (e.g. state upgrades) that have a final `TestStep` with `PlanOnly: true`. Remove `PlanOnly: true` from the final `TestStep` to fix affected tests which will ensure that updated schema changes are applied to the state before attempting to automatically destroy resources. ([#223](https://github.com/hashicorp/terraform-plugin-testing/issues/223))

FEATURES:

* plancheck: Added `ExpectUnknownOutputValue` built-in plan check, which asserts that a given output value at a specified address is unknown ([#220](https://github.com/hashicorp/terraform-plugin-testing/issues/220))
* plancheck: Added `ExpectUnknownOutputValueAtPath` built-in plan check, which asserts that a given output value at a specified address, and path is unknown ([#220](https://github.com/hashicorp/terraform-plugin-testing/issues/220))
* plancheck: Added `ExpectNullOutputValue` built-in plan check, which asserts that a given output value at a specified address is null ([#220](https://github.com/hashicorp/terraform-plugin-testing/issues/220))
* plancheck: Added `ExpectNullOutputValueAtPath` built-in plan check, which asserts that a given output value at a specified address, and path is null ([#220](https://github.com/hashicorp/terraform-plugin-testing/issues/220))

ENHANCEMENTS:

* helper/resource: Removed separate refresh commands, which increases testing performance ([#223](https://github.com/hashicorp/terraform-plugin-testing/issues/223))
* helper/resource: Automatically add `required_providers` configuration to `TestStep.Config` Terraform language configuration when using Terraform >= 1.0.* ([#216](https://github.com/hashicorp/terraform-plugin-testing/issues/216))

BUG FIXES:

* plancheck: Ensured `ExpectEmptyPlan` and `ExpectNonEmptyPlan` account for output changes ([#222](https://github.com/hashicorp/terraform-plugin-testing/issues/222))
* helper/resource: Ensured `TestStep.ExpectNonEmptyPlan` accounts for output changes with Terraform 0.14 and later ([#234](https://github.com/hashicorp/terraform-plugin-testing/issues/234))

## 1.5.1 (August 31, 2023)

BUG FIXES:

* helper/resource: Fix regression by allowing providers to be defined both at the `TestCase` level, and within `TestStep.Config` ([#177](https://github.com/hashicorp/terraform-plugin-testing/issues/177))

## 1.5.0 (August 31, 2023)

FEATURES:

* config: Introduced new `config` package which contains interfaces and helper functions for working with native Terraform configuration and variables ([#153](https://github.com/hashicorp/terraform-plugin-testing/issues/153))
* helper/resource: Added `TestStep.ConfigDirectory` to allow specifying a directory containing Terraform configuration for use during acceptance tests ([#153](https://github.com/hashicorp/terraform-plugin-testing/issues/153))
* helper/resource: Added `TestStep.ConfigFile` to allow specifying a file containing Terraform configuration for use during acceptance tests ([#153](https://github.com/hashicorp/terraform-plugin-testing/issues/153))
* helper/resource: Added `TestStep.ConfigVariables` to allow specifying Terraform variables for use with Terraform configuration during acceptance tests ([#153](https://github.com/hashicorp/terraform-plugin-testing/issues/153))
* helper/resource: Removed data resource and managed resource `id` attribute requirement ([#84](https://github.com/hashicorp/terraform-plugin-testing/issues/84))

ENHANCEMENTS:

* helper/resource: Added `TestStep` type `ImportStateVerifyIdentifierAttribute` field, which can override the default `id` attribute used for matching prior resource state with imported resource state ([#84](https://github.com/hashicorp/terraform-plugin-testing/issues/84))

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

