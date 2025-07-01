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

	// Check g
	if _, err := exec.LookPath("g"); err == nil {
		green.Print("  ✅ g: ")
		if versionCmd := exec.Command("g", "--version"); versionCmd.Run() == nil {
			versionCmd.Stdout = os.Stdout
			versionCmd.Run()
		} else {
			fmt.Println("installed")
		}
	}

	// Check if no version managers are installed
	if !common.IsCommandAvailable("gobrew") && !common.IsCommandAvailable("g") {
		red.Println("  ❌ No version managers installed")
		yellow.Println("  💡 Run: gos setup")
	}
}

// CheckInstalledVersions displays installed Go versions
func CheckInstalledVersions() {
	yellow := color.New(color.FgYellow)

	if common.IsCommandAvailable("gobrew") || common.IsCommandAvailable("g") {
		// Show available versions using direct commands
		showInstalledVersions()
	} else {
		yellow.Println("  No version manager installed")
		yellow.Println("  💡 Run: gos setup")
	}
}

// showInstalledVersions displays installed Go versions using available version managers
func showInstalledVersions() {
	yellow := color.New(color.FgYellow)

	if common.IsCommandAvailable("gobrew") {
		showInstalledVersionsWithGobrew()
	} else if common.IsCommandAvailable("g") {
		showInstalledVersionsWithG()
	} else {
		yellow.Println("  No version manager available")
	}
}

// showInstalledVersionsWithGobrew lists installed versions using gobrew
func showInstalledVersionsWithGobrew() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	if output, err := exec.Command("gobrew", "ls").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					if strings.Contains(line, "*") {
						green.Printf("  ✅ %s (current)\n", strings.Replace(line, "*", "", -1))
					} else {
						fmt.Printf("  📦 %s\n", line)
					}
				}
			}
		} else {
			yellow.Println("  No Go versions installed")
		}
	} else {
		yellow.Printf("  Error listing gobrew versions: %v\n", err)
	}
}

// showInstalledVersionsWithG lists installed versions using g
func showInstalledVersionsWithG() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	if output, err := exec.Command("g", "ls").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					if strings.Contains(line, "*") {
						green.Printf("  ✅ %s (current)\n", strings.Replace(line, "*", "", -1))
					} else {
						fmt.Printf("  📦 %s\n", line)
					}
				}
			}
		} else {
			yellow.Println("  No Go versions installed")
		}
	} else {
		yellow.Printf("  Error listing g versions: %v\n", err)
	}
}
