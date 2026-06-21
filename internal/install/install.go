package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/core-file-privacy/core-file-privacy/internal/files"
)

type InstallOptions struct {
	Force      bool
	NoPath     bool
	PathOnly   bool
	Target     string
	Yes        bool
}

func Install(opts InstallOptions) error {
	shell := DetectShell()
	targetDir := opts.Target
	if targetDir == "" {
		targetDir = files.GetInstallDir()
	}

	installer := NewInstaller(targetDir, shell)
	installer.PrintHeader()

	binaryPath := filepath.Join(targetDir, getBinaryName())

	if !opts.Yes {
		fmt.Fprintf(os.Stderr, "Core File Privacy is not installed in your user tools directory.\n\n")
		fmt.Fprintf(os.Stderr, "Install it now?\nTarget: %s\n\n", binaryPath)

		confirmed, err := confirmInstall()
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, "Installation cancelled.")
			return nil
		}
	}

	if opts.PathOnly {
		return installPathOnly(installer, targetDir, shell)
	}

	totalSteps := 4
	if runtime.GOOS != "windows" && !opts.NoPath {
		totalSteps = 5
	}

	step := 0

	step++
	installer.PrintStep(step, totalSteps, "Creating install directory")
	if err := files.EnsureDir(targetDir); err != nil {
		installer.PrintDone()
		return fmt.Errorf("failed to create directory: %w", err)
	}
	installer.PrintDone()

	step++
	installer.PrintStep(step, totalSteps, "Copying binary")
	if err := copyCurrentBinary(binaryPath); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}
	installer.PrintDone()

	if runtime.GOOS != "windows" {
		step++
		installer.PrintStep(step, totalSteps, "Setting executable permissions")
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return fmt.Errorf("failed to set permissions: %w", err)
		}
		installer.PrintDone()
	}

	configFile := ""
	if !opts.NoPath {
		step++
		installer.PrintStep(step, totalSteps, "Updating PATH")
		if err := AddToPath(shell, targetDir); err != nil {
			installer.PrintDone()
			installer.PrintPathFailed()
			return nil
		}
		installer.PrintDone()

		configFile, _ = GetShellConfigFile(shell)
	}

	step++
	installer.PrintStep(step, totalSteps, "Verifying installation")
	if !files.FileExists(binaryPath) {
		return fmt.Errorf("verification failed: binary not found")
	}
	installer.PrintDone()

	installer.PrintSuccess(binaryPath, configFile)

	return nil
}

func installPathOnly(installer *Installer, targetDir, shell string) error {
	installer.PrintStep(1, 1, "Updating PATH")
	if err := AddToPath(shell, targetDir); err != nil {
		installer.PrintDone()
		installer.PrintPathFailed()
		return nil
	}
	installer.PrintDone()

	configFile, _ := GetShellConfigFile(shell)
	installer.PrintSuccess("", configFile)

	return nil
}

func copyCurrentBinary(destPath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	src, err := os.Open(execPath)
	if err != nil {
		return fmt.Errorf("failed to open source binary: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	return nil
}

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "cfp.exe"
	}
	return "cfp"
}

func confirmInstall() (bool, error) {
	fmt.Fprint(os.Stderr, "[Y/n] ")

	var response string
	fmt.Scanln(&response)

	if response == "" {
		return true, nil
	}

	response = strings.ToLower(response)
	return response == "y" || response == "yes", nil
}

func Uninstall(yes bool) error {
	installDir := files.GetInstallDir()
	binaryPath := files.GetInstallBinaryPath()

	if !files.FileExists(binaryPath) {
		return fmt.Errorf("cfp is not installed in %s", installDir)
	}

	if !yes {
		fmt.Fprintf(os.Stderr, "This will remove cfp from %s\n\n", installDir)
		fmt.Fprint(os.Stderr, "Continue? [y/N] ")

		var response string
		fmt.Scanln(&response)

		response = strings.ToLower(response)
		if response != "y" && response != "yes" {
			fmt.Fprintln(os.Stderr, "Uninstallation cancelled.")
			return nil
		}
	}

	if err := os.Remove(binaryPath); err != nil {
		return fmt.Errorf("failed to remove binary: %w", err)
	}

	entries, err := os.ReadDir(installDir)
	if err == nil && len(entries) == 0 {
		os.Remove(installDir)
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Uninstallation completed.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Note: PATH entries were not removed automatically.")
	fmt.Fprintln(os.Stderr, "You may want to remove them manually from your shell config.")

	return nil
}
