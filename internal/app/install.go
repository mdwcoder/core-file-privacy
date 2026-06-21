package app

import (
	"github.com/core-file-privacy/core-file-privacy/internal/install"
)

type InstallOptions struct {
	Force    bool
	NoPath   bool
	PathOnly bool
	Target   string
	Yes      bool
}

func Install(opts InstallOptions) error {
	return install.Install(install.InstallOptions{
		Force:    opts.Force,
		NoPath:   opts.NoPath,
		PathOnly: opts.PathOnly,
		Target:   opts.Target,
		Yes:      opts.Yes,
	})
}

type UninstallOptions struct {
	Yes bool
}

func Uninstall(opts UninstallOptions) error {
	return install.Uninstall(opts.Yes)
}
