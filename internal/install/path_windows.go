//go:build windows

package install

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func AddToPath(shell, installDir string) error {
	script := fmt.Sprintf(`
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*%s*") {
	$newPath = $currentPath + ";%s"
	[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
	Write-Host "PATH updated"
} else {
	Write-Host "PATH already contains install directory"
}
`, installDir, installDir)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update PATH: %w", err)
	}

	return nil
}

func IsInPath(installDir string) bool {
	path := os.Getenv("PATH")
	paths := strings.Split(path, ";")

	for _, p := range paths {
		if strings.EqualFold(p, installDir) {
			return true
		}
	}

	return false
}
