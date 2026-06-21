package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "0.1.0"
	BuildDate = "unknown"
	Commit    = "unknown"
)

func GetVersion() string {
	return Version
}

func GetFullVersion() string {
	return fmt.Sprintf("cfp version %s\nBuild date: %s\nCommit: %s\nOS/Arch: %s/%s",
		Version, BuildDate, Commit, runtime.GOOS, runtime.GOARCH)
}
