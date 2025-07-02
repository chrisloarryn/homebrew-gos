package list

import (
	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed Go versions",
		Long:  `List all Go versions that have been installed via gobrew.`,
		Run: func(cmd *cobra.Command, args []string) {
			remote, _ := cmd.Flags().GetBool("remote")
			if remote {
				ListRemoteVersions()
			} else {
				ListVersions()
			}
		},
	}

	cmd.Flags().BoolP("remote", "r", false, "List available remote versions")
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
