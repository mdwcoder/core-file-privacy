package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetShellConfigFile(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	if runtime.GOOS == "windows" {
		return "", fmt.Errorf("Windows does not use shell config files")
	}

	switch shell {
	case "bash":
		bashrc := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(bashrc); err == nil {
			return bashrc, nil
		}
		return filepath.Join(home, ".bash_profile"), nil

	case "zsh":
		return filepath.Join(home, ".zshrc"), nil

	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish"), nil

	case "sh":
		return filepath.Join(home, ".profile"), nil

	default:
		return filepath.Join(home, ".profile"), nil
	}
}

func FormatPathLine(shell, installDir string) string {
	if runtime.GOOS == "windows" {
		return ""
	}

	switch shell {
	case "fish":
		return fmt.Sprintf("\n# Core File Privacy\nfish_add_path \"%s\"\n", installDir)
	default:
		return fmt.Sprintf("\n# Core File Privacy\nexport PATH=\"%s:$PATH\"\n", installDir)
	}
}
