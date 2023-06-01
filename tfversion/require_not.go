package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// RequireNot will fail the test if the Terraform CLI
// version matches the given version.
func RequireNot(version *version.Version) TerraformVersionCheck {
	return requireNotCheck{
		version: version,
	}
}

// requireNotCheck implements the TerraformVersionCheck interface
type requireNotCheck struct {
	version *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s requireNotCheck) CheckTerraformVersion(ctx context.Context, req CheckTFVersionRequest, resp *CheckTFVersionResponse) {

	if req.TerraformVersion.Equal(s.version) {
		resp.Error = fmt.Errorf("unexpected Terraform CLI version: %s", s.version)
	}
}
