// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/checkpoint"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
)

// Config is used to configure the test helper. In most normal test programs
// the configuration is discovered automatically by an Init* function using
// DiscoverConfig, but this is exposed so that more complex scenarios can be
// implemented by direct configuration.
type Config struct {
	SourceDir          string
	TerraformExec      string
	execTempDir        string
	PreviousPluginExec string
}

// DiscoverConfig uses environment variables and other means to automatically
// discover a reasonable test helper configuration.
func DiscoverConfig(ctx context.Context, sourceDir string) (*Config, error) {
	tfVersion := strings.TrimPrefix(os.Getenv(EnvTfAccTerraformVersion), "v")
	tfPath := os.Getenv(EnvTfAccTerraformPath)

	tempDir := os.Getenv(EnvTfAccTempDir)
	tfDir, err := os.MkdirTemp(tempDir, "plugintest-terraform")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	var sources []src.Source
	// Upstream hc-install does not provide a String() method to identify the source, store details of each source to contextualize error during `installer.Ensure`
	var sourceDetails []string

	switch {
	case tfPath != "":
		sourceDetail := fmt.Sprintf("Terraform CLI source of exact path: %s", tfPath)
		sourceDetails = append(sourceDetails, sourceDetail)
		logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential %s", sourceDetail))

		sources = append(sources, &fs.AnyVersion{
			ExactBinPath: tfPath,
		})
	case tfVersion != "":
		tfVersion, err := version.NewVersion(tfVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid Terraform version: %w", err)
		}

		sourceDetail := fmt.Sprintf("Terraform CLI source of releases.hashicorp.com exact version %q for installation in: %s", tfVersion, tfDir)
		sourceDetails = append(sourceDetails, sourceDetail)
		logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential %s", sourceDetail))

		sources = append(sources, &releases.ExactVersion{
			InstallDir: tfDir,
			Product:    product.Terraform,
			Version:    tfVersion,
		})
	default:
		sourceDetail := "Terraform CLI source of local filesystem PATH lookup"
		sourceDetails = append(sourceDetails, sourceDetail)
		logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential %s", sourceDetail))

		sources = append(sources, &fs.AnyVersion{
			Product: &product.Terraform,
		})

		sourceDetail = fmt.Sprintf("Terraform CLI source of checkpoint.hashicorp.com latest version for installation in: %s", tfDir)
		sourceDetails = append(sourceDetails, sourceDetail)
		logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential %s", sourceDetail))

		sources = append(sources, &checkpoint.LatestVersion{
			InstallDir: tfDir,
			Product:    product.Terraform,
		})
	}

	installer := install.NewInstaller()
	tfExec, err := installer.Ensure(context.Background(), sources)
	if err != nil {
		return nil, fmt.Errorf("error ensuring Terraform CLI binary: %w -  attempted source(s): %v", err, sourceDetails)
	}

	ctx = logging.TestTerraformPathContext(ctx, tfExec)

	logging.HelperResourceDebug(ctx, "Found Terraform CLI")

	return &Config{
		SourceDir:     sourceDir,
		TerraformExec: tfExec,
		execTempDir:   tfDir,
	}, nil
}
