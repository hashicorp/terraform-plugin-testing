// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

// TODO: doc
type DeferredReason string

const (
	// TODO: doc
	DeferredReasonResourceConfigUnknown DeferredReason = "resource_config_unknown"
	// TODO: doc
	DeferredReasonProviderConfigUnknown DeferredReason = "provider_config_unknown"
	// TODO: doc
	DeferredReasonAbsentPrereq DeferredReason = "absent_prereq"
)
