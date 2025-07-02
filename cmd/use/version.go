package use

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// UseVersion switches to a specific Go version
func UseVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	blue.Printf("🔄 Switching to Go %s...\n", version)

	cmd := createVersionSwitchCommand(version, blue, red)
	if cmd == nil {
		return
	}

	if !executeVersionSwitch(cmd, version, blue, red, yellow) {
		return
	}

	green.Printf("✅ Version switch command completed\n")

	// Update PATH and verify installation
	performPostSwitchVerification(version, blue, green, yellow)
}

// createVersionSwitchCommand creates the command for version switching using gobrew
func createVersionSwitchCommand(version string, blue, red *color.Color) *exec.Cmd {
	blue.Println("  Using gobrew...")
	cmd := exec.Command("gobrew", "use", version)

	// Set output to see what's happening
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// executeVersionSwitch executes the version switch command
func executeVersionSwitch(cmd *exec.Cmd, version string, blue, red, yellow *color.Color) bool {
	blue.Println("  Executing version switch...")
	if err := cmd.Run(); err != nil {
		red.Printf("❌ Error switching to Go %s: %v\n", version, err)
		yellow.Printf("💡 Is this version installed? Use: gos list\n")
		return false
	}
	return true
}

// performPostSwitchVerification handles PATH update and verification
func performPostSwitchVerification(version string, blue, green, yellow *color.Color) {
	// Clean all Go paths and update PATH for version manager
	blue.Println("  Cleaning Go paths and updating PATH...")
	common.UpdatePathForVersionManagerClean()

	// Show current version and PATH update instructions
	blue.Println("\n📋 Verifying installation...")

	// Use the version manager's Go binary directly for verification
	goPath := getVersionManagerGoPath()

	if goPath != "" {
		verifyWithDirectPath(goPath, green, yellow)
	} else {
		verifyWithPathResolution(version, green, yellow)
	}
}

// getVersionManagerGoPath returns the path to gobrew's Go binary
func getVersionManagerGoPath() string {
	homeDir := common.GetHomeDir()
	return filepath.Join(homeDir, ".gobrew", "current", "bin", "go")
}

// verifyWithDirectPath verifies Go version using direct path
func verifyWithDirectPath(goPath string, green, yellow *color.Color) {
	// Check version manager's Go version
	var vmVersion string
	if output, err := exec.Command(goPath, "version").Output(); err == nil {
		vmVersion = strings.TrimSpace(string(output))
		green.Printf("✅ Version manager: %s\n", vmVersion)
	} else {
		yellow.Printf("⚠️  Error checking version manager: %v\n", err)
		showPathUpdateInstructions(yellow)
		return
	}

	// Check system's Go version (from PATH)
	if output, err := exec.Command("go", "version").Output(); err == nil {
		systemVersion := strings.TrimSpace(string(output))
		if vmVersion == systemVersion {
			green.Printf("✅ System version: %s\n", systemVersion)
		} else {
			yellow.Printf("⚠️  System version: %s\n", systemVersion)
			yellow.Println("⚠️  Version mismatch! Multiple Go installations detected in PATH.")
			showPathCleanupInstructions(yellow)
		}
	} else {
		yellow.Println("⚠️  Go binary not found in system PATH")
		showPathUpdateInstructions(yellow)
	}
}

// verifyWithPathResolution verifies Go version using PATH resolution
func verifyWithPathResolution(version string, green, yellow *color.Color) {
	if output, err := exec.Command("go", "version").Output(); err == nil {
		currentVersion := strings.TrimSpace(string(output))
		if strings.Contains(currentVersion, version) {
			green.Printf("✅ Current version: %s\n", currentVersion)
		} else {
			yellow.Printf("⚠️  Version mismatch - found: %s\n", currentVersion)
			yellow.Printf("   Expected: %s\n", version)
			showPathUpdateInstructions(yellow)
		}
	} else {
		yellow.Println("⚠️  Go binary not found in PATH")
		showPathUpdateInstructions(yellow)
	}
}

// showPathUpdateInstructions displays instructions for updating PATH
func showPathUpdateInstructions(yellow *color.Color) {
	yellow.Println("⚠️  PATH needs to be updated for this terminal session.")
	yellow.Println("💡 To use the new Go version immediately, run:")
	if runtime.GOOS == "windows" {
		yellow.Println("   $env:PATH = \"$env:USERPROFILE\\.gobrew\\current\\bin;$env:USERPROFILE\\.gobrew\\bin;$env:PATH\"")
	} else {
		yellow.Println("   export PATH=\"$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH\"")
	}
	fmt.Println()
	color.New(color.FgBlue).Println("🔄 Or simply open a new terminal window.")
}

// showPathCleanupInstructions displays instructions for cleaning Go paths from PATH
func showPathCleanupInstructions(yellow *color.Color) {
	yellow.Println("💡 To fix PATH conflicts, you have several options:")
	fmt.Println()

	// Offer interactive cleanup
	fmt.Println("Would you like automated help with PATH cleanup? (y/n)")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		common.PromptUserForPathCleanup()
	} else {
		// Show manual instructions
		yellow.Println("📋 Manual cleanup instructions:")
		yellow.Println()
		yellow.Println("   🧹 Clean your shell configuration files (~/.zshrc, ~/.bashrc, etc.):")
		yellow.Println("   Remove or comment out lines that add Go paths like:")
		yellow.Println("     - export PATH=/usr/local/go/bin:$PATH")
		yellow.Println("     - export PATH=$HOME/sdk/go*/bin:$PATH")
		yellow.Println("     - export PATH=$HOME/go/bin:$PATH (keep this one)")
		yellow.Println()
		yellow.Println("   ✅ Keep only the gobrew paths:")
		yellow.Println("     export PATH=\"$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$HOME/go/bin:$PATH\"")
		yellow.Println()
		yellow.Println("   🔄 After editing, restart your terminal or run: source ~/.zshrc")
		fmt.Println()
	}
}
