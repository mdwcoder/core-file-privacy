package files

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetInstallDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return filepath.Join(home, "core-file-privacy")
}

func GetInstallBinaryPath() string {
	installDir := GetInstallDir()

	if runtime.GOOS == "windows" {
		return filepath.Join(installDir, "cfp.exe")
	}

	return filepath.Join(installDir, "cfp")
}

func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		return true
	}

	return info.Mode()&0111 != 0
}
