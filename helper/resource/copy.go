package resource

import (
	"io"
	"os"
	"path"
)

// CopyFile copies a single file from src to dest.
func CopyFile(src, dest string) error {
	var srcFileInfo os.FileInfo

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}

	if srcFileInfo, err = os.Stat(src); err != nil {
		return err
	}

	return os.Chmod(dest, srcFileInfo.Mode())
}

// CopyDir recursively copies directories and files
// from src to dest.
func CopyDir(src string, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	dirEntries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		srcFilepath := path.Join(src, dirEntry.Name())
		destFilepath := path.Join(dest, dirEntry.Name())

		if dirEntry.IsDir() {
			if err = CopyDir(srcFilepath, destFilepath); err != nil {
				return err
			}
		} else {
			if err = CopyFile(srcFilepath, destFilepath); err != nil {
				return err
			}
		}
	}

	return nil
}
