package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific Go version",
	Long: `Install a specific Go version using the 'g' version manager.
If no version is specified, installs the latest stable version.`,
	Example: `  gos install 1.21.5    # Install Go 1.21.5
  gos install latest     # Install latest version
  gos install            # Install latest version (default)`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !checkGInstalled() {
			return
		}

		version := "latest"
		if len(args) > 0 {
			version = args[0]
		}

		installVersion(version)
	},
}

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a specific Go version",
	Long:  `Switch to a specific Go version that has been previously installed.`,
	Example: `  gos use 1.21.5        # Switch to Go 1.21.5
  gos use latest         # Switch to latest installed version`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !checkGInstalled() {
			return
		}
		useVersion(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Go versions",
	Long:  `List all Go versions that have been installed via the 'g' version manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !checkGInstalled() {
			return
		}

		remote, _ := cmd.Flags().GetBool("remote")
		if remote {
			listRemoteVersions()
		} else {
			listVersions()
		}
	},
}

var removeCmd = &cobra.Command{
	Use:     "remove [version]",
	Short:   "Remove a specific Go version",
	Long:    `Remove a specific Go version that has been installed via the 'g' version manager.`,
	Example: `  gos remove 1.20.10    # Remove Go 1.20.10`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !checkGInstalled() {
			return
		}
		removeVersion(args[0])
	},
}

var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Install and use the latest Go version",
	Long:  `Install the latest stable Go version and automatically switch to it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !checkGInstalled() {
			return
		}
		installLatest()
	},
}

var projectCmd = &cobra.Command{
	Use:   "project [version]",
	Short: "Configure Go version for current project",
	Long: `Configure a specific Go version for the current project by creating a .go-version file
and switching to that version.`,
	Example: `  gos project 1.21.5    # Configure project to use Go 1.21.5`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !checkGInstalled() {
			return
		}
		setupProjectVersion(args[0])
	},
}

func init() {
	listCmd.Flags().BoolP("remote", "r", false, "List available remote versions")
}

// Helper functions
func checkGInstalled() bool {
	_, err := exec.LookPath("g")
	if err != nil {
		color.Red("❌ Error: The 'g' manager is not installed.")
		color.Yellow("💡 Run first: gos setup")
		return false
	}
	return true
}

func installVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	blue.Printf("📦 Installing Go %s...\n", version)

	cmd := exec.Command("g", "install", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		red.Printf("❌ Error installing Go %s\n", version)
		return
	}

	green.Printf("✅ Go %s installed successfully\n", version)
}

func useVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	blue.Printf("🔄 Switching to Go %s...\n", version)

	cmd := exec.Command("g", "use", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		red.Printf("❌ Error switching to Go %s\n", version)
		yellow.Printf("💡 Is this version installed? Use: gos list\n")
		return
	}

	green.Printf("✅ Switched to Go %s\n", version)

	// Show current version
	if goCmd := exec.Command("go", "version"); goCmd.Run() == nil {
		blue.Print("📋 Current version: ")
		goCmd.Stdout = os.Stdout
		goCmd.Run()
	}
}

func listVersions() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("📋 Installed Go versions:")

	cmd := exec.Command("g", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		yellow.Println("No versions installed")
	}
}

func listRemoteVersions() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("🌐 Available versions (latest 10):")

	cmd := exec.Command("g", "list-all")
	output, err := cmd.Output()
	if err != nil {
		yellow.Println("Could not get remote versions")
		return
	}

	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		if line != "" && count < 10 {
			fmt.Println(line)
			count++
		}
	}
}

func removeVersion(version string) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	yellow.Printf("🗑️  Removing Go %s...\n", version)

	cmd := exec.Command("g", "remove", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		red.Printf("❌ Error removing Go %s\n", version)
		return
	}

	green.Printf("✅ Go %s removed successfully\n", version)
}

func installLatest() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	blue.Println("🚀 Installing latest Go version...")

	installCmd := exec.Command("g", "install", "latest")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		red.Println("❌ Error installing latest version")
		return
	}

	green.Println("✅ Latest version installed")

	// Switch to latest
	useCmd := exec.Command("g", "use", "latest")
	useCmd.Run()

	// Show current version
	if goCmd := exec.Command("go", "version"); goCmd.Run() == nil {
		blue.Print("📋 Current version: ")
		goCmd.Stdout = os.Stdout
		goCmd.Run()
	}
}

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
	useVersion(version)

	green.Printf("✅ Project configured to use Go %s\n", version)
	blue.Printf("📄 File created: %s\n", goVersionFile)
}
