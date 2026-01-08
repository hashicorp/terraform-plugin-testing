package statestore_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"testing/fstest"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
)

// This state store implementation uses an in-memory map for storing state in a filesystem-like manner.
// This only works because the terraform-plugin-testing harness keeps provider instances running
// throughout a Go test with multiple Terraform CLI command calls by using a reattach configuration.
func exampleCloudStateStoreFS() *testprovider.StateStore {
	memFS := fstest.MapFS{}

	return &testprovider.StateStore{
		SchemaResponse: &statestore.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{Attributes: []*tfprotov6.SchemaAttribute{}},
			},
		},
		GetStatesFunc: func(ctx context.Context, req statestore.GetStatesRequest, resp *statestore.GetStatesResponse) {
			directories, err := memFS.ReadDir(".")
			if err != nil {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Error reading inmem filesystem",
					Detail:   err.Error(),
				})
				return
			}

			workspaces := make([]string, 0)
			for _, dir := range directories {
				workspaces = append(workspaces, dir.Name())
			}

			resp.StateIDs = workspaces
		},
		DeleteStateFunc: func(ctx context.Context, req statestore.DeleteStateRequest, resp *statestore.DeleteStateResponse) {
			for filePath := range memFS {
				if strings.HasPrefix(filePath, req.StateID) {
					delete(memFS, filePath)
				}
			}
		},
		LockStateFunc: func(ctx context.Context, req statestore.LockStateRequest, resp *statestore.LockStateResponse) {
			logging.HelperResourceDebug(ctx, "examplecloud_inmem: Lock support not implemented")
		},
		UnlockStateFunc: func(ctx context.Context, req statestore.UnlockStateRequest, resp *statestore.UnlockStateResponse) {
			logging.HelperResourceDebug(ctx, "examplecloud_inmem: Lock support not implemented")
		},
		ReadStateBytesFunc: func(ctx context.Context, req statestore.ReadStateBytesRequest, resp *statestore.ReadStateBytesResponse) {
			stateFilePath := filepath.Join(req.StateID, "terraform.tfstate")
			stateFile, err := memFS.Open(stateFilePath)
			if err != nil {
				// If there is no state file, Terraform will create one.
				if errors.Is(err, fs.ErrNotExist) {
					return
				}

				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  fmt.Sprintf("Error state %q at path %q", req.StateID, stateFilePath),
					Detail:   err.Error(),
				})
				return
			}

			stateBytes, err := io.ReadAll(stateFile)
			if err != nil {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  fmt.Sprintf("Error reading %q state bytes", req.StateID),
					Detail:   err.Error(),
				})
				return
			}

			resp.StateBytes = stateBytes
		},
		WriteStateBytesFunc: func(ctx context.Context, req statestore.WriteStateBytesRequest, resp *statestore.WriteStateBytesResponse) {
			stateFilePath := filepath.Join(req.StateID, "terraform.tfstate")
			memFS[stateFilePath] = &fstest.MapFile{
				Data:    req.StateBytes,
				Mode:    fs.ModePerm,
				ModTime: time.Now(),
			}
		},
	}
}
