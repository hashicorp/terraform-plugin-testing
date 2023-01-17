package resource

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("cannot create src dir: %s", err)
	}
	defer os.RemoveAll(srcDir)

	_, err = os.Create(filepath.Join(srcDir, "src.txt"))
	if err != nil {
		t.Fatalf("cannot create src file: %s", err)
	}

	destDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("cannot create dest dir: %s", err)
	}
	defer os.RemoveAll(destDir)

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

	if !Equal(t, srcDirEntries, destDirEntries) {
		t.Fatalf("dir entries differ: %v, %v", srcDirEntries, destDirEntries)
	}
}

func TestCopyDir(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("cannot create src dir: %s", err)
	}
	defer os.RemoveAll(srcDir)

	_, err = os.Create(filepath.Join(srcDir, "src.txt"))
	if err != nil {
		t.Fatalf("cannot create src file: %s", err)
	}

	err = CopyDir(srcDir, srcDir+"_1")
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

	if !Equal(t, srcDirEntries, destDirEntries) {
		t.Fatalf("dir entries differ: %v, %v", srcDirEntries, destDirEntries)
	}
}

func Equal(t *testing.T, a, b []os.DirEntry) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v.Type() != b[i].Type() {
			return false
		}

		if v.Name() != b[i].Name() {
			return false
		}

		aInfo, err := v.Info()
		if err != nil {
			t.Fatalf("cannot get file info: %s", err)
		}

		bInfo, err := b[i].Info()
		if err != nil {
			t.Fatalf("cannot get file info: %s", err)
		}

		if aInfo.Mode() != bInfo.Mode() {
			return false
		}
	}

	return true
}
