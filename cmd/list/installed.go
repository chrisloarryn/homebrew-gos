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

	// Use gobrew to list versions
	if _, err := exec.LookPath("gobrew"); err == nil {
		listVersionsWithGobrew()
	} else {
		yellow.Println("gobrew not detected.")
		fmt.Println("")
		yellow.Println("💡 To install gobrew:")
		fmt.Println("   gos setup               # Install gobrew")
		fmt.Println("")
		yellow.Println("💡 If Go is installed manually, check with:")
		fmt.Println("   go version              # Show current Go version")
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
