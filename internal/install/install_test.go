package install

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDetectOS(t *testing.T) {
	detected := DetectOS()
	if detected != runtime.GOOS {
		t.Errorf("DetectOS() = %s, want %s", detected, runtime.GOOS)
	}
}

func TestDetectArch(t *testing.T) {
	detected := DetectArch()
	if detected != runtime.GOARCH {
		t.Errorf("DetectArch() = %s, want %s", detected, runtime.GOARCH)
	}
}

func TestGetShellConfigFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}

	home, _ := os.UserHomeDir()

	tests := []struct {
		shell    string
		contains string
	}{
		{"bash", ".bashrc"},
		{"zsh", ".zshrc"},
		{"fish", "config.fish"},
		{"sh", ".profile"},
	}

	for _, tt := range tests {
		configFile, err := GetShellConfigFile(tt.shell)
		if err != nil {
			t.Errorf("GetShellConfigFile(%s) failed: %v", tt.shell, err)
			continue
		}

		if !strings.Contains(configFile, tt.contains) {
			t.Errorf("GetShellConfigFile(%s) = %s, should contain %s", tt.shell, configFile, tt.contains)
		}

		if !strings.HasPrefix(configFile, home) {
			t.Errorf("config file should be in home directory")
		}
	}
}

func TestFormatPathLine(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}

	installDir := "/home/user/core-file-privacy"

	bashLine := FormatPathLine("bash", installDir)
	if !strings.Contains(bashLine, "export PATH") {
		t.Error("bash path line should contain export PATH")
	}
	if !strings.Contains(bashLine, installDir) {
		t.Error("bash path line should contain install dir")
	}

	zshLine := FormatPathLine("zsh", installDir)
	if !strings.Contains(zshLine, "export PATH") {
		t.Error("zsh path line should contain export PATH")
	}

	fishLine := FormatPathLine("fish", installDir)
	if !strings.Contains(fishLine, "fish_add_path") {
		t.Error("fish path line should contain fish_add_path")
	}
}

func TestIsInPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}

	tmpDir := t.TempDir()

	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", originalPath+":"+tmpDir)

	if !IsInPath(tmpDir) {
		t.Error("IsInPath should return true for directory in PATH")
	}

	if IsInPath("/nonexistent/path") {
		t.Error("IsInPath should return false for directory not in PATH")
	}
}

func TestAddToPathIdempotent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, ".zshrc")

	installDir := filepath.Join(tmpDir, "core-file-privacy")

	if err := os.WriteFile(configFile, []byte("# existing config\n"), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	if err := AddToPath("zsh", installDir); err != nil {
		t.Fatalf("first AddToPath failed: %v", err)
	}

	content1, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	if err := AddToPath("zsh", installDir); err != nil {
		t.Fatalf("second AddToPath failed: %v", err)
	}

	content2, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	if string(content1) != string(content2) {
		t.Error("AddToPath should be idempotent")
	}

	count := strings.Count(string(content2), installDir)
	if count != 1 {
		t.Errorf("install dir should appear exactly once, appeared %d times", count)
	}
}
