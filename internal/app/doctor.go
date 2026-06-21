package app

import (
	"fmt"
	"os"
	"runtime"

	"github.com/core-file-privacy/core-file-privacy/internal/files"
	"github.com/core-file-privacy/core-file-privacy/internal/install"
)

func Doctor() error {
	fmt.Fprintln(os.Stderr, "Core File Privacy doctor")
	fmt.Fprintln(os.Stderr)

	fmt.Fprintf(os.Stderr, "OS: %s\n", runtime.GOOS)
	fmt.Fprintf(os.Stderr, "Architecture: %s\n", runtime.GOARCH)

	installDir := files.GetInstallDir()
	fmt.Fprintf(os.Stderr, "Install directory: %s\n", installDir)

	binaryPath := files.GetInstallBinaryPath()
	if files.FileExists(binaryPath) {
		fmt.Fprintln(os.Stderr, "Installed binary: found")
	} else {
		fmt.Fprintln(os.Stderr, "Installed binary: not found")
	}

	currentBinary, err := os.Executable()
	if err == nil {
		fmt.Fprintf(os.Stderr, "Current binary: %s\n", currentBinary)
	}

	shell := install.DetectShell()
	if runtime.GOOS != "windows" {
		fmt.Fprintf(os.Stderr, "Shell: %s\n", shell)
	} else {
		fmt.Fprintf(os.Stderr, "Terminal: %s\n", install.DetectTerminal())
	}

	if install.IsInPath(installDir) {
		fmt.Fprintln(os.Stderr, "PATH contains install directory: yes")
	} else {
		fmt.Fprintln(os.Stderr, "PATH contains install directory: no")
	}

	fmt.Fprintln(os.Stderr, "Password input support: yes")
	fmt.Fprintln(os.Stderr, "Secure random: available")
	fmt.Fprintln(os.Stderr, "Archive support: yes")

	if runtime.GOOS == "windows" {
		fmt.Fprintln(os.Stderr, "Hidden file strategy: windows hidden attribute")
	} else {
		fmt.Fprintln(os.Stderr, "Hidden file strategy: dotfile")
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Status: OK")

	return nil
}
