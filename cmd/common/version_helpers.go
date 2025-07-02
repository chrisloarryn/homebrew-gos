package common

import (
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// IsCommandAvailable checks if a command is available in PATH
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// CheckVersionManagerAvailable checks if gobrew version manager is available
func CheckVersionManagerAvailable() bool {
	if IsCommandAvailable("gobrew") {
		return true
	}

	color.Red("❌ Error: No version manager is installed.")
	color.Yellow("💡 Run first: gos setup")
	return false
}

// IsGInstalled checks if the gobrew version manager is installed
func IsGInstalled() bool {
	return IsCommandAvailable("gobrew")
}

// GetGobrewVersions returns installed versions using gobrew
func GetGobrewVersions() []string {
	if _, err := exec.LookPath("gobrew"); err != nil {
		return []string{}
	}

	cmd := exec.Command("gobrew", "ls")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(string(output), "\n")
	var versions []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "=>") { // Skip empty lines and current indicator
			versions = append(versions, line)
		}
	}

	return versions
}

// GetCurrentGoVersion returns the currently active Go version
func GetCurrentGoVersion() string {
	if IsCommandAvailable("gobrew") {
		return getCurrentGoVersionWithGobrew()
	}

	// Fallback: check system Go
	if output, err := exec.Command("go", "version").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}

	return ""
}

// getCurrentGoVersionWithGobrew gets current version using gobrew
func getCurrentGoVersionWithGobrew() string {
	if _, err := exec.LookPath("gobrew"); err != nil {
		return ""
	}

	cmd := exec.Command("gobrew", "ls")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "*") || strings.Contains(line, "current") {
			return strings.ReplaceAll(line, "*", "")
		}
	}
	return ""
}
