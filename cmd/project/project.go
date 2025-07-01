package project

import (
	"os"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/cristobalcontreras/gos/cmd/use"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewProjectCmd creates the project command
func NewProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "project [version]",
		Short: "Configure Go version for current project",
		Long: `Configure a specific Go version for the current project by creating a .go-version file
and switching to that version.`,
		Example: `  gos project 1.21.5    # Configure project to use Go 1.21.5`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !common.CheckVersionManagerAvailable() {
				return
			}
			setupProjectVersion(args[0])
		},
	}
}

// setupProjectVersion configures a specific Go version for the current project
func setupProjectVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	blue.Printf("📁 Configuring version %s for this project...\n", version)

	// Create .go-version file
	goVersionFile := ".go-version"
	if err := os.WriteFile(goVersionFile, []byte(version), 0644); err != nil {
		color.Red("❌ Error creating .go-version file: %v", err)
		return
	}

	// Switch to that version
	use.UseVersion(version)

	green.Printf("✅ Project configured to use Go %s\n", version)
	blue.Printf("📄 File created: %s\n", goVersionFile)
}
