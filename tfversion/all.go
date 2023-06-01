package tfversion

import (
	"context"
)

// All will return the first non-nil error or non-empty skip message
// if any of the given checks return a non-nil error or non-empty skip message.
// Otherwise, it will return a nil error and empty skip message (run the test)
func All(terraformVersionChecks ...TerraformVersionCheck) TerraformVersionCheck {
	return allCheck{
		terraformVersionChecks: terraformVersionChecks,
	}
}

// allCheck implements the TerraformVersionCheck interface
type allCheck struct {
	terraformVersionChecks []TerraformVersionCheck
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (a allCheck) CheckTerraformVersion(ctx context.Context, req CheckTFVersionRequest, resp *CheckTFVersionResponse) {

	for _, subCheck := range a.terraformVersionChecks {
		checkResp := CheckTFVersionResponse{}

		subCheck.CheckTerraformVersion(ctx, CheckTFVersionRequest{TerraformVersion: req.TerraformVersion}, &checkResp)

		if checkResp.Error != nil {
			resp.Error = checkResp.Error
			return
		}

		if checkResp.Skip != "" {
			resp.Skip = checkResp.Skip
			return
		}
	}
}
