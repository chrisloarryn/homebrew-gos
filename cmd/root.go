package cmd

import (
	"fmt"
	"os"

	"github.com/cristobalcontreras/gos/cmd/clean"
	defaultcmd "github.com/cristobalcontreras/gos/cmd/default"
	"github.com/cristobalcontreras/gos/cmd/env"
	"github.com/cristobalcontreras/gos/cmd/install"
	"github.com/cristobalcontreras/gos/cmd/latest"
	"github.com/cristobalcontreras/gos/cmd/list"
	"github.com/cristobalcontreras/gos/cmd/project"
	"github.com/cristobalcontreras/gos/cmd/reload"
	"github.com/cristobalcontreras/gos/cmd/remove"
	"github.com/cristobalcontreras/gos/cmd/setup"
	"github.com/cristobalcontreras/gos/cmd/status"
	"github.com/cristobalcontreras/gos/cmd/use"
	versioncmd "github.com/cristobalcontreras/gos/cmd/version"
	"github.com/spf13/cobra"
)

var (
	// Version information
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gos",
	Short: "A comprehensive Go version manager CLI",
	Long: `GOS is a powerful command-line tool for managing Go versions.
It provides functionality to install, switch, and manage multiple Go versions
using the 'g' version manager, along with comprehensive cleanup capabilities.

Features:
- Install and switch between Go versions
- Setup the 'g' version manager
- Deep clean Go installations
- Project-specific version management
- System status and diagnostics`,
	Example: `  gos install 1.21.5     # Install Go 1.21.5
  gos use 1.21.5          # Switch to Go 1.21.5
  gos setup               # Setup the 'g' version manager
  gos clean               # Deep clean Go installations
  gos status              # Show system status`,
	Version: getVersionString(),
}

// SetVersionInfo sets the version information
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = getVersionString()
}

func getVersionString() string {
	if version == "dev" {
		return fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	}
	return fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(install.NewInstallCmd())
	rootCmd.AddCommand(use.NewUseCmd())
	rootCmd.AddCommand(list.NewListCmd())
	rootCmd.AddCommand(remove.NewRemoveCmd())
	rootCmd.AddCommand(clean.NewCleanCmd())
	rootCmd.AddCommand(setup.NewSetupCmd())
	rootCmd.AddCommand(status.CreateStatusCommand())
	rootCmd.AddCommand(project.NewProjectCmd())
	rootCmd.AddCommand(latest.NewLatestCmd())
	rootCmd.AddCommand(reload.NewReloadCmd())
	rootCmd.AddCommand(defaultcmd.CreateDefaultCommand())
	rootCmd.AddCommand(env.CreateEnvCommand())
	rootCmd.AddCommand(versioncmd.NewVersionCmd())
}
