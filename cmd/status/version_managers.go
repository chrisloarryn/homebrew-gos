package status

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// CheckVersionManagers displays information about available version managers
func CheckVersionManagers() {
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	// Check gobrew
	if _, err := exec.LookPath("gobrew"); err == nil {
		green.Print("  ✅ gobrew: ")
		if versionCmd := exec.Command("gobrew", "--version"); versionCmd.Run() == nil {
			versionCmd.Stdout = os.Stdout
			versionCmd.Run()
		} else {
			fmt.Println("installed")
		}
	}

	// Check if no version managers are installed
	if !common.IsCommandAvailable("gobrew") {
		red.Println("  ❌ No version managers installed")
		yellow.Println("  💡 Run: gos setup")
	}
}

// CheckInstalledVersions displays installed Go versions
func CheckInstalledVersions() {
	yellow := color.New(color.FgYellow)

	if common.IsCommandAvailable("gobrew") {
		// Show available versions using direct commands
		showInstalledVersions()
	} else {
		yellow.Println("  No version manager installed")
		yellow.Println("  💡 Run: gos setup")
	}
}

// showInstalledVersions displays installed Go versions using gobrew
func showInstalledVersions() {
	yellow := color.New(color.FgYellow)

	if common.IsCommandAvailable("gobrew") {
		showInstalledVersionsWithGobrew()
	} else {
		yellow.Println("  No version manager available")
	}
}

// showInstalledVersionsWithGobrew lists installed versions using gobrew
func showInstalledVersionsWithGobrew() {
	yellow := color.New(color.FgYellow)

	output, err := exec.Command("gobrew", "ls").Output()
	if err != nil {
		yellow.Printf("  Error listing gobrew versions: %v\n", err)
		return
	}

	lines := parseVersionOutput(output)
	if len(lines) == 0 {
		yellow.Println("  No Go versions installed")
		return
	}

	displayVersions(lines)
}

// parseVersionOutput parses the output from version manager commands
func parseVersionOutput(output []byte) []string {
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var validLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			validLines = append(validLines, line)
		}
	}

	return validLines
}

// displayVersions displays the parsed version lines with appropriate formatting
func displayVersions(lines []string) {
	green := color.New(color.FgGreen)

	for _, line := range lines {
		if strings.Contains(line, "*") {
			cleanLine := strings.Replace(line, "*", "", -1)
			green.Printf("  ✅ %s (current)\n", cleanLine)
		} else {
			fmt.Printf("  📦 %s\n", line)
		}
	}
}
