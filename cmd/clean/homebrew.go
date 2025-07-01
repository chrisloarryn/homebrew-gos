package clean

import (
	"fmt"
	"os/exec"
	"strings"
)

// CleanHomebrewGo removes Go installations from Homebrew
func CleanHomebrewGo() {
	if _, err := exec.LookPath("brew"); err != nil {
		return
	}

	// Get list of Go formulas
	cmd := exec.Command("brew", "list", "--formula")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go") && (line == "go" || strings.Contains(line, "go@")) {
			fmt.Printf("  – brew uninstall %s\n", line)
			uninstallCmd := exec.Command("brew", "uninstall", "--ignore-dependencies", "--force", line)
			uninstallCmd.Run()
		}
	}
}
