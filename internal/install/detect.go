package install

import (
	"os"
	"runtime"
	"strings"
)

func DetectOS() string {
	return runtime.GOOS
}

func DetectArch() string {
	return runtime.GOARCH
}

func DetectShell() string {
	if runtime.GOOS == "windows" {
		return detectWindowsTerminal()
	}

	return detectUnixShell()
}

func detectUnixShell() string {
	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		return parts[len(parts)-1]
	}

	return "unknown"
}

func detectWindowsTerminal() string {
	if os.Getenv("WT_SESSION") != "" {
		return "windows-terminal"
	}

	if os.Getenv("PSModulePath") != "" {
		return "powershell"
	}

	if os.Getenv("CMDLINE") != "" {
		return "cmd"
	}

	return "unknown"
}

func DetectTerminal() string {
	if runtime.GOOS != "windows" {
		return "terminal"
	}

	return detectWindowsTerminal()
}
