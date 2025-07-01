package use

import (
	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewUseCmd creates the use command
func NewUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [version]",
		Short: "Switch to a specific Go version",
		Long:  `Switch to a specific Go version that has been previously installed.`,
		Example: `  gos use 1.21.5        # Switch to Go 1.21.5
  gos use latest         # Switch to latest installed version`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !checkVersionManagerAvailable() {
				return
			}
			UseVersion(args[0])
		},
	}
	
	return cmd
}

// checkVersionManagerAvailable verifies if a version manager is available
func checkVersionManagerAvailable() bool {
	if common.IsCommandAvailable("gobrew") || common.IsCommandAvailable("g") {
		return true
	}

	color.Red("❌ Error: No version manager is installed.")
	color.Yellow("💡 Run first: gos setup")
	return false
}
