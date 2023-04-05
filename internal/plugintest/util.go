// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func symlinkFile(src string, dest string) error {
	err := os.Symlink(src, dest)

	if err != nil {
		return fmt.Errorf("unable to symlink %q to %q: %w", src, dest, err)
	}

	srcInfo, err := os.Stat(src)

	if err != nil {
		return fmt.Errorf("unable to stat %q: %w", src, err)
	}

	err = os.Chmod(dest, srcInfo.Mode())

	if err != nil {
		return fmt.Errorf("unable to set %q permissions: %w", dest, err)
	}

	return nil
}

// symlinkDirectoriesOnly finds only the first-level child directories in srcDir
// and symlinks them into destDir.
// Unlike symlinkDir, this is done non-recursively in order to limit the number
// of file descriptors used.
func symlinkDirectoriesOnly(srcDir string, destDir string) error {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("unable to stat source directory %q: %w", srcDir, err)
	}

	err = os.MkdirAll(destDir, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("unable to make destination directory %q: %w", destDir, err)
	}

	dirEntries, err := os.ReadDir(srcDir)

	if err != nil {
		return fmt.Errorf("unable to read source directory %q: %w", srcDir, err)
	}

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}

		srcPath := filepath.Join(srcDir, dirEntry.Name())
		destPath := filepath.Join(destDir, dirEntry.Name())
		err := symlinkFile(srcPath, destPath)

		if err != nil {
			return fmt.Errorf("unable to symlink directory %q to %q: %w", srcPath, destPath, err)
		}
	}

	return nil
}

// CopyFile copies a single file from src to dest.
func CopyFile(src, dest string) error {
	var srcFileInfo os.FileInfo

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("unable to copy: %w", err)
	}

	if srcFileInfo, err = os.Stat(src); err != nil {
		return fmt.Errorf("unable to stat: %w", err)
	}

	return os.Chmod(dest, srcFileInfo.Mode())
}

// CopyDir recursively copies directories and files
// from src to dest.
func CopyDir(src, dest, baseDirName string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("unable to stat: %w", err)
	}

	if err = os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return fmt.Errorf("unable to create dir: %w", err)
	}

	dirEntries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("unable to read dir: %w", err)
	}

	for _, dirEntry := range dirEntries {
		srcFilepath := path.Join(src, dirEntry.Name())
		destFilepath := path.Join(dest, dirEntry.Name())

		if !strings.Contains(srcFilepath, baseDirName) {
			continue
		}

		if dirEntry.IsDir() {
			if err = CopyDir(srcFilepath, destFilepath, baseDirName); err != nil {
				return fmt.Errorf("unable to copy directory: %w", err)
			}
		} else {
			if err = CopyFile(srcFilepath, destFilepath); err != nil {
				return fmt.Errorf("unable to copy file: %w", err)
			}
		}
	}

	return nil
}
