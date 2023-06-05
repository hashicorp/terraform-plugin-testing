// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCopyFile(t *testing.T) {
	t.Parallel()

	srcDir := t.TempDir()

	_, err := os.Create(filepath.Join(srcDir, "src.txt"))
	if err != nil {
		t.Fatalf("cannot create src file: %s", err)
	}

	destDir := t.TempDir()

	err = CopyFile(filepath.Join(srcDir, "src.txt"), filepath.Join(destDir, "src.txt"))
	if err != nil {
		t.Fatalf("cannot copy src file: %s", err)
	}

	srcDirEntries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("cannot read src dir: %s", srcDir)
	}

	destDirEntries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("cannot read dest dir: %s", srcDir)
	}

	if diff := cmp.Diff(srcDirEntries, destDirEntries, dirEntryComparer(t)); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

func TestCopyDir(t *testing.T) {
	t.Parallel()

	srcDir := t.TempDir()

	_, err := os.Create(filepath.Join(srcDir, "src.txt"))
	if err != nil {
		t.Fatalf("cannot create src file: %s", err)
	}

	err = CopyDir(srcDir, srcDir+"_1", "")
	if err != nil {
		t.Fatalf("cannot copy dir: %s", err)
	}
	defer os.RemoveAll(srcDir + "_1")

	srcDirEntries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("cannot read src dir: %s", srcDir)
	}

	destDirEntries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("cannot read dest dir: %s", srcDir)
	}

	if diff := cmp.Diff(srcDirEntries, destDirEntries, dirEntryComparer(t)); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

func dirEntryComparer(t *testing.T) cmp.Option {
	return cmp.Comparer(func(x, y []os.DirEntry) bool {
		if len(x) != len(y) {
			return false
		}

		for k, v := range x {
			if v.Type() != y[k].Type() {
				return false
			}

			if v.Name() != y[k].Name() {
				return false
			}

			vInfo, err := v.Info()
			if err != nil {
				t.Errorf("could not get FileInfo for v: %s", err)
			}

			ykInfo, err := y[k].Info()
			if err != nil {
				t.Errorf("could not get FileInfo for y[%d]: %s", k, err)
			}

			if vInfo.IsDir() != ykInfo.IsDir() {
				return false
			}

			if vInfo.Mode() != ykInfo.Mode() {
				return false
			}

			if vInfo.Name() != ykInfo.Name() {
				return false
			}

			if vInfo.Size() != ykInfo.Size() {
				return false
			}

			if vInfo.ModTime() != ykInfo.ModTime() {
				return false
			}
		}

		return true
	})
}
