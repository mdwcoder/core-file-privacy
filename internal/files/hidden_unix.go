//go:build !windows

package files

import (
	"os"
	"path/filepath"
	"strings"
)

func HideFile(path string) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	if strings.HasPrefix(base, ".") {
		return nil
	}

	hiddenPath := filepath.Join(dir, "."+base)

	if err := os.Rename(path, hiddenPath); err != nil {
		return err
	}

	return nil
}

func IsHidden(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".")
}

func UnhideFile(path string) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	if !strings.HasPrefix(base, ".") {
		return nil
	}

	unhiddenPath := filepath.Join(dir, strings.TrimPrefix(base, "."))

	if err := os.Rename(path, unhiddenPath); err != nil {
		return err
	}

	return nil
}

func GetHiddenName(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	if !strings.HasPrefix(base, ".") {
		return filepath.Join(dir, "."+base)
	}

	return path
}
