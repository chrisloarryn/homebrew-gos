package install

import (
	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewInstallCmd creates the install command
func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [version]",
		Short: "Install a specific Go version",
		Long: `Install a specific Go version using gobrew.
If no version is specified, installs the latest stable version.`,
		Example: `  gos install 1.21.5    # Install Go 1.21.5
  gos install latest     # Install latest version
  gos install            # Install latest version (default)`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !checkVersionManagerAvailable() {
				return
			}

			version := "latest"
			if len(args) > 0 {
				version = args[0]
			}

			InstallVersion(version)
		},
	}

	return cmd
}

// checkVersionManagerAvailable verifies if gobrew is available
func checkVersionManagerAvailable() bool {
	if common.IsCommandAvailable("gobrew") {
		return true
	}

	color.Red("❌ Error: gobrew is not installed.")
	color.Yellow("💡 Run first: gos setup")
	return false
}
