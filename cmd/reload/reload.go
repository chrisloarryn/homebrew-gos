package reload

import (
	"fmt"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewReloadCmd creates the reload command
func NewReloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reload",
		Short: "Reload Go environment configuration",
		Long: `Reload Go environment configuration and verify it's working correctly.
	
This command will:
- Source shell configuration files
- Verify Go is available in PATH
- Check GOROOT and GOPATH settings
- Show current configuration status`,
		Run: func(cmd *cobra.Command, args []string) {
			reloadEnvironment()
		},
	}
}

// reloadEnvironment reloads and verifies the Go environment
func reloadEnvironment() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	blue.Println("🔄 Reloading Go environment...")

	// Update environment variables and PATH
	common.SetupGoEnvironment()

	// Verify Go installation
	fmt.Println("")
	blue.Println("🔍 Verifying Go installation...")

	if !common.VerifyGoInstallation() {
		return
	}

	// Verify GOROOT and GOPATH
	common.VerifyGoEnvironmentPaths()

	fmt.Println("")
	green.Println("🎉 Environment reload complete!")
	
	// Show helpful commands
	showUsefulCommands()
}

// showUsefulCommands displays helpful commands to the user
func showUsefulCommands() {
	blue := color.New(color.FgBlue)
	
	fmt.Println("")
	blue.Println("💡 Useful commands:")
	fmt.Println("  gos status        # Check overall status")
	fmt.Println("  gos env           # Show detailed environment")
	fmt.Println("  gos list          # List installed Go versions")
	fmt.Println("  go version        # Verify active Go version")
}
