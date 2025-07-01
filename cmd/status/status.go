package status

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// CreateStatusCommand creates the status command
func CreateStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show Go system status",
		Long: `Show comprehensive information about the current Go installation,
including version manager status, installed versions, and system configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			ShowStatus()
		},
	}
}

// ShowStatus displays comprehensive Go system status
func ShowStatus() {
	blue := color.New(color.FgBlue)

	blue.Println("📊 Go system status:")
	fmt.Println("")

	// Check version managers
	blue.Println("🔧 Version Managers:")
	CheckVersionManagers()

	fmt.Println("")

	// Current Go installation
	blue.Println("🐹 Current Go:")
	ShowCurrentGo()

	fmt.Println("")

	// Installed versions
	blue.Println("📦 Installed versions:")
	CheckInstalledVersions()

	fmt.Println("")

	// Disk space
	blue.Println("💾 Disk space:")
	ShowDiskUsage()

	fmt.Println("")

	// Environment variables
	blue.Println("🌍 Environment:")
	ShowEnvironment()

	fmt.Println("")

	// Project configuration
	blue.Println("📁 Project configuration:")
	ShowProjectConfig()
}
