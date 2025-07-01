package list

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// ListVersions lists all installed Go versions
func ListVersions() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("📋 Installed Go versions:")

	// Try different version managers based on availability
	if common.IsCommandAvailable("gobrew") {
		listVersionsWithGobrew()
	} else if common.IsCommandAvailable("g") {
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

// listVersionsWithGobrew lists versions using gobrew
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

// listVersionsWithG lists versions using g
func listVersionsWithG() {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	
	// Try to get the list of installed versions
	cmd := exec.Command("g", "list")
	output, err := cmd.Output()
	if err != nil {
		yellow.Println("  No versions installed via g")
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		yellow.Println("  No versions installed via g")
		return
	}

	fmt.Println("  Using g...")
	
	// Get current version using common helper
	currentVersion := common.GetCurrentGoVersion()

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "No") {
			if line == currentVersion || strings.Contains(line, "*") {
				green.Printf("  ✅ %s (current)\n", strings.ReplaceAll(line, "*", ""))
			} else {
				fmt.Printf("     %s\n", line)
			}
		}
	}
}

// listVersionsManually checks for manual Go installations
func listVersionsManually() bool {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	// Use common helper to get system Go info
	if version, goroot, found := common.GetSystemGoInfo(); found {
		green.Printf("  ✅ %s (system installation)\n", version)

		if goroot != "" {
			fmt.Printf("     Location: %s\n", goroot)
		}

		fmt.Println("")
		yellow.Println("💡 This appears to be a manual Go installation.")
		yellow.Println("   To manage multiple versions, consider installing a version manager:")
		fmt.Println("   gos setup               # Install version manager")

		return true
	}

	return false
}
