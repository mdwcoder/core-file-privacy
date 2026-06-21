//go:build !windows

package install

import (
	"fmt"
	"os"
	"strings"
)

func AddToPath(shell, installDir string) error {
	configFile, err := GetShellConfigFile(shell)
	if err != nil {
		return err
	}

	pathLine := FormatPathLine(shell, installDir)

	content, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if strings.Contains(string(content), installDir) {
		return nil
	}

	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(pathLine); err != nil {
		return fmt.Errorf("failed to write to config file: %w", err)
	}

	return nil
}

func IsInPath(installDir string) bool {
	path := os.Getenv("PATH")
	paths := strings.Split(path, ":")

	for _, p := range paths {
		if p == installDir {
			return true
		}
	}

	return false
}
