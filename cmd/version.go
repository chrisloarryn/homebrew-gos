package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific Go version",
	Long: `Install a specific Go version using any available version manager (gobrew or g).
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
		if !checkVersionManagerAvailable() {
			return
		}
		useVersion(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Go versions",
	Long:  `List all Go versions that have been installed via any available version manager (gobrew, g) or manual installation.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	Long:    `Remove a specific Go version that has been installed via any available version manager.`,
	Example: `  gos remove 1.20.10    # Remove Go 1.20.10`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !checkVersionManagerAvailable() {
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
		if !checkVersionManagerAvailable() {
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
		if !checkVersionManagerAvailable() {
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

func checkVersionManagerAvailable() bool {
	if isCommandAvailable("gobrew") || isCommandAvailable("g") {
		return true
	}

	color.Red("❌ Error: No version manager is installed.")
	color.Yellow("💡 Run first: gos setup")
	return false
}

func installVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	blue.Printf("📦 Installing Go %s...\n", version)

	var cmd *exec.Cmd
	if isCommandAvailable("gobrew") {
		cmd = exec.Command("gobrew", "install", version)
	} else if isCommandAvailable("g") {
		cmd = exec.Command("g", "install", version)
	} else {
		red.Println("❌ No version manager available")
		return
	}

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

	var cmd *exec.Cmd
	if isCommandAvailable("gobrew") {
		cmd = exec.Command("gobrew", "use", version)
	} else if isCommandAvailable("g") {
		cmd = exec.Command("g", "set", version)
	} else {
		red.Println("❌ No version manager available")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		red.Printf("❌ Error switching to Go %s\n", version)
		yellow.Printf("💡 Is this version installed? Use: gos list\n")
		return
	}

	// Update PATH automatically after successful version switch
	updatePathForVersionManager()

	green.Printf("✅ Switched to Go %s\n", version)

	// Show current version
	if goCmd := exec.Command("go", "version"); goCmd.Run() == nil {
		blue.Print("📋 Current version: ")
		goCmd.Stdout = os.Stdout
		goCmd.Run()
	} else {
		yellow.Println("⚠️  Go not found in PATH. You may need to restart your terminal or run:")
		if isCommandAvailable("gobrew") {
			yellow.Println("   $env:PATH = \"$env:USERPROFILE\\.gobrew\\current\\bin;$env:PATH\"")
		}
	}
}

func listVersions() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("📋 Installed Go versions:")

	// Try different version managers based on availability
	if isCommandAvailable("gobrew") {
		listVersionsWithGobrew()
	} else if isCommandAvailable("g") {
		listVersionsWithG()
	} else {
		// Fallback: check for direct installations or show manual detection
		if !listVersionsManually() {
			yellow.Println("No version manager detected.")
			fmt.Println("")
			yellow.Println("💡 To install a version manager:")
			fmt.Println("   gos setup               # Install version manager")
			fmt.Println("")
			yellow.Println("💡 If Go is installed manually, check with:")
			fmt.Println("   go version              # Show current Go version")
		}
	}
}

func listVersionsWithGobrew() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Println("  Using gobrew...")

	cmd := exec.Command("gobrew", "ls")
	output, err := cmd.Output()
	if err != nil {
		yellow.Println("  No versions installed via gobrew")
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if strings.Contains(line, "*") || strings.Contains(line, "current") {
				green.Printf("  ✅ %s (current)\n", strings.ReplaceAll(line, "*", ""))
			} else {
				fmt.Printf("     %s\n", line)
			}
		}
	}
}

func listVersionsWithG() {
	cmd := exec.Command("g", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		yellow := color.New(color.FgYellow)
		yellow.Println("  No versions installed via g")
	}
}

func listVersionsManually() bool {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	// Check if Go is installed directly
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("  ✅ %s (system installation)\n", version)

		// Try to get GOROOT to see where it's installed
		if goroot, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
			fmt.Printf("     Location: %s\n", strings.TrimSpace(string(goroot)))
		}

		fmt.Println("")
		yellow.Println("💡 This appears to be a manual Go installation.")
		yellow.Println("   To manage multiple versions, consider installing a version manager:")
		fmt.Println("   gos setup               # Install version manager")

		return true
	}

	return false
}

func listRemoteVersions() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("🌐 Available versions:")

	// Try different version managers based on availability
	if isCommandAvailable("gobrew") {
		listRemoteVersionsWithGobrew()
	} else if isCommandAvailable("g") {
		listRemoteVersionsWithG()
	} else {
		yellow.Println("No version manager detected.")
		fmt.Println("")
		yellow.Println("💡 To install a version manager and browse remote versions:")
		fmt.Println("   gos setup               # Install version manager")
		fmt.Println("")
		yellow.Println("💡 You can also check manually at:")
		fmt.Println("   https://golang.org/dl/")
	}
}

func listRemoteVersionsWithGobrew() {
	yellow := color.New(color.FgYellow)

	fmt.Println("  Using gobrew...")

	cmd := exec.Command("gobrew", "ls-remote")
	output, err := cmd.Output()
	if err != nil {
		yellow.Println("  Could not get remote versions via gobrew")
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && count < 15 { // Show more versions since they're useful
			fmt.Printf("     %s\n", line)
			count++
		}
	}

	if count >= 15 {
		fmt.Println("     ... (more versions available)")
		fmt.Println("")
		yellow.Println("💡 To see all versions: gobrew ls-remote")
	}
}

func listRemoteVersionsWithG() {
	yellow := color.New(color.FgYellow)

	cmd := exec.Command("g", "ls-remote")
	output, err := cmd.Output()
	if err != nil {
		// Try alternative command
		cmd = exec.Command("g", "list-all")
		output, err = cmd.Output()
		if err != nil {
			yellow.Println("  Could not get remote versions")
			return
		}
	}

	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && count < 15 {
			fmt.Printf("     %s\n", line)
			count++
		}
	}

	if count >= 15 {
		fmt.Println("     ... (more versions available)")
		fmt.Println("")
		yellow.Println("💡 To see all versions: g ls-remote")
	}
}

func removeVersion(version string) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	yellow.Printf("🗑️  Removing Go %s...\n", version)

	var cmd *exec.Cmd
	if isCommandAvailable("gobrew") {
		cmd = exec.Command("gobrew", "uninstall", version)
	} else if isCommandAvailable("g") {
		cmd = exec.Command("g", "remove", version)
	} else {
		red.Println("❌ No version manager available")
		return
	}

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

	var installCmd *exec.Cmd
	if isCommandAvailable("gobrew") {
		installCmd = exec.Command("gobrew", "install", "latest")
	} else if isCommandAvailable("g") {
		installCmd = exec.Command("g", "install", "latest")
	} else {
		red.Println("❌ No version manager available")
		return
	}

	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		red.Println("❌ Error installing latest version")
		return
	}

	green.Println("✅ Latest version installed")

	// Switch to latest
	var useCmd *exec.Cmd
	if isCommandAvailable("gobrew") {
		useCmd = exec.Command("gobrew", "use", "latest")
	} else if isCommandAvailable("g") {
		useCmd = exec.Command("g", "set", "latest")
	}

	if useCmd != nil {
		useCmd.Run()
	}

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

func getCurrentVersion() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("📍 Current Go version:")

	goCmd := exec.Command("go", "version")
	if err := goCmd.Run(); err != nil {
		yellow.Println("⚠️  Go is not available in PATH")
		return
	}

	goCmd.Stdout = os.Stdout
	goCmd.Run()

	// Show GOROOT and GOPATH
	if gorootCmd := exec.Command("go", "env", "GOROOT"); gorootCmd.Run() == nil {
		blue.Print("📂 GOROOT: ")
		gorootCmd.Stdout = os.Stdout
		gorootCmd.Run()
	}

	if gopathCmd := exec.Command("go", "env", "GOPATH"); gopathCmd.Run() == nil {
		blue.Print("📂 GOPATH: ")
		gopathCmd.Stdout = os.Stdout
		gopathCmd.Run()
	}
}

// updatePathForVersionManager updates the PATH environment variable to include the current Go version
func updatePathForVersionManager() {
	homeDir := getHomeDir()
	currentPath := os.Getenv("PATH")

	var newGoBin string
	var gobreWBin string

	if isCommandAvailable("gobrew") {
		// For gobrew, add both gobrew bin and current Go bin to PATH
		newGoBin = filepath.Join(homeDir, ".gobrew", "current", "bin")
		gobreWBin = filepath.Join(homeDir, ".gobrew", "bin")

		// Remove any existing gobrew paths from PATH to avoid duplicates
		pathSep := ";"
		if runtime.GOOS != "windows" {
			pathSep = ":"
		}

		pathParts := strings.Split(currentPath, pathSep)
		var cleanPaths []string

		for _, part := range pathParts {
			// Skip existing gobrew paths to avoid duplicates
			if !strings.Contains(part, ".gobrew") {
				cleanPaths = append(cleanPaths, part)
			}
		}

		// Prepend gobrew paths
		finalPaths := []string{newGoBin, gobreWBin}
		finalPaths = append(finalPaths, cleanPaths...)
		newPath := strings.Join(finalPaths, pathSep)

		os.Setenv("PATH", newPath)

	} else if isCommandAvailable("g") {
		// For g version manager
		newGoBin = filepath.Join(homeDir, ".g", "go", "bin")
		gBin := filepath.Join(homeDir, ".g", "bin")

		pathSep := ";"
		if runtime.GOOS != "windows" {
			pathSep = ":"
		}

		pathParts := strings.Split(currentPath, pathSep)
		var cleanPaths []string

		for _, part := range pathParts {
			// Skip existing g paths to avoid duplicates
			if !strings.Contains(part, ".g") {
				cleanPaths = append(cleanPaths, part)
			}
		}

		// Prepend g paths
		finalPaths := []string{newGoBin, gBin}
		finalPaths = append(finalPaths, cleanPaths...)
		newPath := strings.Join(finalPaths, pathSep)

		os.Setenv("PATH", newPath)
	}
}
