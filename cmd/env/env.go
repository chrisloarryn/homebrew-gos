package env

import (
	"github.com/spf13/cobra"
)

// CreateEnvCommand creates the env command
func CreateEnvCommand() *cobra.Command {
	envCmd := &cobra.Command{
		Use:   "env",
		Short: "Show or fix environment configuration",
		Long: `Show current environment configuration and optionally fix it.
	
Examples:
  gos env          # Show current environment
  gos env --fix    # Fix environment configuration
  gos env --export # Export current environment for sourcing
  gos env --check  # Run comprehensive environment validation`,
		Run: func(cmd *cobra.Command, args []string) {
			fix, _ := cmd.Flags().GetBool("fix")
			export, _ := cmd.Flags().GetBool("export")
			check, _ := cmd.Flags().GetBool("check")

			if export {
				ExportEnvironment()
			} else if fix {
				FixEnvironment()
			} else if check {
				ValidateEnvironment()
			} else {
				ShowDetailedEnvironment()
			}
		},
	}

	envCmd.Flags().Bool("fix", false, "Fix environment configuration")
	envCmd.Flags().Bool("export", false, "Export environment variables for sourcing")
	envCmd.Flags().Bool("check", false, "Run comprehensive environment validation")

	return envCmd
}
