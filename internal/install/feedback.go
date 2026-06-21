package install

import (
	"fmt"
	"os"
)

type Installer struct {
	targetDir string
	shell     string
	os        string
}

func NewInstaller(targetDir, shell string) *Installer {
	return &Installer{
		targetDir: targetDir,
		shell:     shell,
		os:        DetectOS(),
	}
}

func (i *Installer) PrintHeader() {
	fmt.Fprintln(os.Stderr, "Core File Privacy installer")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "OS detected: %s\n", i.os)

	if i.os != "windows" {
		fmt.Fprintf(os.Stderr, "Shell detected: %s\n", i.shell)
	} else {
		fmt.Fprintf(os.Stderr, "Terminal: %s\n", DetectTerminal())
	}

	fmt.Fprintf(os.Stderr, "Install target: %s\n", i.targetDir)
	fmt.Fprintln(os.Stderr)
}

func (i *Installer) PrintStep(step, total int, message string) {
	fmt.Fprintf(os.Stderr, "[%d/%d] %s... ", step, total, message)
}

func (i *Installer) PrintDone() {
	fmt.Fprintln(os.Stderr, "done")
}

func (i *Installer) PrintSuccess(binaryPath string, configFile string) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Installation completed.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Program installed at:")
	fmt.Fprintf(os.Stderr, "  %s\n", binaryPath)

	if configFile != "" {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "PATH updated in:")
		fmt.Fprintf(os.Stderr, "  %s\n", configFile)
	}

	if i.os != "windows" {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Shell detected:")
		fmt.Fprintf(os.Stderr, "  %s\n", i.shell)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Restart your terminal or run:")
		fmt.Fprintf(os.Stderr, "  source %s\n", configFile)
	} else {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "PATH updated:")
		fmt.Fprintln(os.Stderr, "  User PATH")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Restart PowerShell/CMD/Windows Terminal and run:")
		fmt.Fprintln(os.Stderr, "  cfp help")
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Try:")
	fmt.Fprintln(os.Stderr, "  cfp help")
}

func (i *Installer) PrintPathFailed() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Installation completed, but PATH could not be updated automatically.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Add this directory to your PATH manually:")
	fmt.Fprintf(os.Stderr, "  %s\n", i.targetDir)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Then restart your terminal.")
}
