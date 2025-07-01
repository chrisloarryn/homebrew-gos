package defaultcmd

import (
	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/spf13/cobra"
)

// CreateDefaultCommand creates the default command
func CreateDefaultCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "default [version]",
		Short: "Set a default Go version",
		Long: `Set a specific Go version as the default. This version will be used when no project-specific version is configured.
If no version is specified, shows the current default version.`,
		Example: `  gos default 1.21.5       # Set Go 1.21.5 as default
  gos default               # Show current default version`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !common.CheckVersionManagerAvailable() {
				return
			}
			
			if len(args) == 0 {
				ShowDefaultVersion()
			} else {
				SetDefaultVersion(args[0])
			}
		},
	}
}
