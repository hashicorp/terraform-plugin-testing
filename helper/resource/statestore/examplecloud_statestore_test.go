// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing/fstest"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
)

var defaultLockFileName = "terraform.tfstate.tflock"
var defaultStateFileName = "terraform.tfstate"

// This state store implementation uses an in-memory map for storing state in a filesystem-like manner.
// This only works because the terraform-plugin-testing harness keeps provider instances running
// throughout a Go test with multiple Terraform CLI command calls by using a reattach configuration.
func exampleCloudValidStateStore() *testprovider.StateStore {
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
			lockFilePath := filepath.Join(req.StateID, defaultLockFileName)
			if lockFile, lockExists := memFS[lockFilePath]; lockExists {
				var lockData lockInfo
				err := json.Unmarshal(lockFile.Data, &lockData)
				if err != nil {
					resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "Error reading existing inmem filesystem lock",
						Detail:   err.Error(),
					})
					return
				}

				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Workspace is currently locked",
					Detail:   fmt.Sprintf("Workspace %q is currently locked by another client: \n\n%s", req.StateID, lockData),
				})
				return
			}

			lockInfo := newLockInfo(req)
			lockInfoBytes, err := json.Marshal(lockInfo)
			if err != nil {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Error creating lock file for inmem filesystem",
					Detail:   err.Error(),
				})
				return
			}

			memFS[lockFilePath] = &fstest.MapFile{
				Data:    lockInfoBytes,
				Mode:    fs.ModePerm,
				ModTime: time.Now(),
			}

			resp.LockID = lockInfo.ID
		},
		UnlockStateFunc: func(ctx context.Context, req statestore.UnlockStateRequest, resp *statestore.UnlockStateResponse) {
			lockFilePath := filepath.Join(req.StateID, defaultLockFileName)
			if _, lockExists := memFS[lockFilePath]; !lockExists {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Error unlocking inmem filesystem",
					Detail:   fmt.Sprintf("Workspace %q has already been unlocked.", req.StateID),
				})
				return
			}

			lockFile := memFS[lockFilePath]

			var lockData lockInfo
			err := json.Unmarshal(lockFile.Data, &lockData)
			if err != nil {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Error reading existing inmem filesystem lock",
					Detail:   err.Error(),
				})
				return
			}

			if lockData.ID != req.LockID {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Error unlocking inmem filesystem",
					Detail:   fmt.Sprintf("Workspace %q is currently locked by a different client, with lock ID %q. Terraform attempted to unlock with lock ID %q", req.StateID, lockData.ID, req.LockID),
				})
				return
			}

			delete(memFS, lockFilePath)
		},
		ReadStateBytesFunc: func(ctx context.Context, req statestore.ReadStateBytesRequest, resp *statestore.ReadStateBytesResponse) {
			stateFilePath := filepath.Join(req.StateID, defaultStateFileName)
			stateFile, err := memFS.Open(stateFilePath)
			if err != nil {
				// If there is no state file, Terraform will create one.
				if errors.Is(err, fs.ErrNotExist) {
					return
				}

				resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  fmt.Sprintf("Error reading state %q at path %q", req.StateID, stateFilePath),
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
			stateFilePath := filepath.Join(req.StateID, defaultStateFileName)
			memFS[stateFilePath] = &fstest.MapFile{
				Data:    req.StateBytes,
				Mode:    fs.ModePerm,
				ModTime: time.Now(),
			}
		},
	}
}

type lockInfo struct {
	ID        string
	Who       string
	Operation string
}

func (l *lockInfo) String() string {
	return fmt.Sprintf(`Lock Info:
  ID:        %s
  Operation: %s
  Who:       %s
`, l.ID, l.Operation, l.Who)
}

func newLockInfo(req statestore.LockStateRequest) lockInfo {
	rngSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, 16)
	rngSource.Read(buf)

	// This function will only return an error if the buffer length is != 16
	id, _ := uuid.FormatUUID(buf)

	var userName string
	if userInfo, err := user.Current(); err == nil {
		userName = userInfo.Username
	}
	host, _ := os.Hostname()

	return lockInfo{
		ID:        id,
		Who:       fmt.Sprintf("%s@%s", userName, host),
		Operation: req.Operation,
	}
}
