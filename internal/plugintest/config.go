// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/checkpoint"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"

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

var versionLocks map[string]*sync.Mutex
var versionLocksMutex sync.Mutex

func init() {
	versionLocks = make(map[string]*sync.Mutex)
}

// DiscoverConfig uses environment variables and other means to automatically
// discover a reasonable test helper configuration.
func DiscoverConfig(ctx context.Context, sourceDir string) (*Config, error) {
	tfVersion := strings.TrimPrefix(os.Getenv(EnvTfAccTerraformVersion), "v")
	return DiscoverConfigWithExactVersion(ctx, sourceDir, tfVersion)
}

func DiscoverConfigWithExactVersion(ctx context.Context, sourceDir string, exactVersion string) (*Config, error) {
	tfPath := os.Getenv(EnvTfAccTerraformPath)

	tempDir := os.Getenv(EnvTfAccTempDir)
	tfDir, err := os.MkdirTemp(tempDir, "plugintest-terraform")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	installLogger := tfsdklog.GetSDKSubsystemLogger(ctx, logging.SubsystemInstall)

	var sources []src.Source
	switch {
	case tfPath != "":
		logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential Terraform CLI source of exact path: %s", tfPath))

		sources = append(sources, &fs.AnyVersion{
			ExactBinPath: tfPath,
		})

	case exactVersion != "":
		tfVersion, err := version.NewVersion(exactVersion)

		if err != nil {
			return nil, fmt.Errorf("invalid Terraform version: %w", err)
		}

		logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential Terraform CLI source of releases.hashicorp.com exact version %q for installation in: %s", tfVersion, tfDir))

		if tempDir == "" {
			tempDir = os.TempDir()
		}
		pathParts := []string{
			strings.TrimRight(tempDir, string(os.PathSeparator)),
			"plugintest-terraform",
			strconv.Itoa(os.Getpid()),
			tfVersion.String(),
		}
		tfDir := strings.Join(pathParts, string(os.PathSeparator))
		if err := os.MkdirAll(tfDir, 0o700); err != nil {
			return nil, fmt.Errorf("failed to create temporary directory for Terraform CLI: %w", err)
		}
		findSource := &fs.ExactVersion{
			ExtraPaths: []string{tfDir},
			Product:    product.Terraform,
			Version:    tfVersion,
		}
		releasesSource := &releases.ExactVersion{
			InstallDir: tfDir,
			Product:    product.Terraform,
			Version:    tfVersion,
		}
		sources = append(sources, findSource, releasesSource)

	default:
		logging.HelperResourceTrace(ctx, "Adding potential Terraform CLI source of local filesystem PATH lookup")
		logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential Terraform CLI source of checkpoint.hashicorp.com latest version for installation in: %s", tfDir))

		sources = append(sources, &fs.AnyVersion{
			Product: &product.Terraform,
		})
		sources = append(sources, &checkpoint.LatestVersion{
			InstallDir: tfDir,
			Product:    product.Terraform,
		})
	}

	stdlibLogger := installLogger.StandardLogger(&hclog.StandardLoggerOptions{
		InferLevels: true,
	})

	installer := install.NewInstaller()
	installer.SetLogger(stdlibLogger)

	var installerMutex *sync.Mutex
	if len(exactVersion) > 0 {
		normalizedTFVersion := version.Must(version.NewVersion(exactVersion)).String()
		versionLocksMutex.Lock()
		if _, ok := versionLocks[normalizedTFVersion]; !ok {
			versionLocks[normalizedTFVersion] = &sync.Mutex{}
		}
		installerMutex = versionLocks[normalizedTFVersion]
		versionLocksMutex.Unlock()
	}

	if installerMutex != nil {
		installerMutex.Lock()
	}
	tfExec, err := installer.Ensure(context.Background(), sources)
	if err != nil {
		return nil, fmt.Errorf("failed to find or install Terraform CLI from %+v: %w", sources, err)
	}
	if installerMutex != nil {
		installerMutex.Unlock()
	}

	ctx = logging.TestTerraformPathContext(ctx, tfExec)

	logging.HelperResourceDebug(ctx, "Found Terraform CLI")

	return &Config{
		SourceDir:     sourceDir,
		TerraformExec: tfExec,
		execTempDir:   tfDir,
	}, nil
}

func TFExactVersion(ctx context.Context, v string) error {
	tfVersion, err := version.NewVersion(v)
	var sources []src.Source

	if err != nil {
		return fmt.Errorf("invalid Terraform version: %w", err)
	}

	tempDir := os.TempDir()
	pathParts := []string{
		strings.TrimRight(tempDir, string(os.PathSeparator)),
		"plugintest-terraform",
		strconv.Itoa(os.Getpid()),
		tfVersion.String(),
	}
	tfDir := strings.Join(pathParts, string(os.PathSeparator))
	if err := os.MkdirAll(tfDir, 0o700); err != nil {
		return fmt.Errorf("failed to create temporary directory for Terraform CLI: %w", err)
	}

	logging.HelperResourceTrace(ctx, fmt.Sprintf("Adding potential Terraform CLI source of releases.hashicorp.com exact version %q for installation in: %s", tfVersion, tfDir))

	findSource := &fs.ExactVersion{
		ExtraPaths: []string{tfDir},
		Product:    product.Terraform,
		Version:    tfVersion,
	}
	releasesSource := &releases.ExactVersion{
		InstallDir: tfDir,
		Product:    product.Terraform,
		Version:    tfVersion,
	}
	sources = append(sources, findSource, releasesSource)

	installLogger := tfsdklog.GetSDKSubsystemLogger(ctx, logging.SubsystemInstall)
	stdlibLogger := installLogger.StandardLogger(&hclog.StandardLoggerOptions{
		InferLevels: true,
	})

	installer := install.NewInstaller()
	installer.SetLogger(stdlibLogger)

	var installerMutex *sync.Mutex
	normalizedTFVersion := tfVersion.String()
	versionLocksMutex.Lock()
	if _, ok := versionLocks[normalizedTFVersion]; !ok {
		versionLocks[normalizedTFVersion] = &sync.Mutex{}
	}
	installerMutex = versionLocks[normalizedTFVersion]
	versionLocksMutex.Unlock()

	installerMutex.Lock()
	tfExec, err := installer.Ensure(context.Background(), sources)
	if err != nil {
		return fmt.Errorf("failed to find or install Terraform CLI from %+v: %w", sources, err)
	}
	installerMutex.Unlock()

	ctx = logging.TestTerraformPathContext(ctx, tfExec)

	logging.HelperResourceDebug(ctx, "Found Terraform CLI")
	return nil
}
