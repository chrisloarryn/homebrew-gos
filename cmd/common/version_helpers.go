package common

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// IsCommandAvailable checks if a command is available in PATH
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// CheckVersionManagerAvailable checks if any version manager is available
func CheckVersionManagerAvailable() bool {
	if IsCommandAvailable("gobrew") || IsCommandAvailable("g") {
		return true
	}

	color.Red("❌ Error: No version manager is installed.")
	color.Yellow("💡 Run first: gos setup")
	return false
}

// CheckGInstalled checks if the 'g' version manager is installed (legacy function)
func CheckGInstalled() bool {
	_, err := exec.LookPath("g")
	if err != nil {
		color.Red("❌ Error: The 'g' manager is not installed.")
		color.Yellow("💡 Run first: gos setup")
		return false
	}
	return true
}

// IsGInstalled checks if the 'g' version manager is installed with path detection
func IsGInstalled() bool {
	// Check for gobrew on Windows first
	if runtime.GOOS == "windows" && IsCommandAvailable("gobrew") {
		return true
	}
	
	homeDir := GetHomeDir()
	gPaths := []string{
		filepath.Join(homeDir, ".g", "bin", "g"),
		filepath.Join(homeDir, "go", "bin", "g"),
		UsrLocalBinG,
	}
	
	for _, path := range gPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	
	return IsCommandAvailable("g")
}

// GetInstalledVersions returns installed versions using 'g'
func GetInstalledVersions() []string {
	homeDir := GetHomeDir()
	gPaths := []string{
		filepath.Join(homeDir, ".g", "bin", "g"),
		filepath.Join(homeDir, "go", "bin", "g"),
		UsrLocalBinG,
	}
	
	var gPath string
	for _, path := range gPaths {
		if _, err := os.Stat(path); err == nil {
			gPath = path
			break
		}
	}
	
	if gPath == "" {
		return []string{}
	}
	
	cmd := exec.Command(gPath, "list")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	
	lines := strings.Split(string(output), "\n")
	var versions []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && line != "g" { // Skip empty lines and header
			versions = append(versions, line)
		}
	}
	
	return versions
}

// GetGobrewVersions returns installed versions using gobrew
func GetGobrewVersions() []string {
	if !IsCommandAvailable("gobrew") {
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
	} else if IsCommandAvailable("g") {
		return getCurrentGoVersionWithG()
	}
	
	// Fallback: check system Go
	if output, err := exec.Command("go", "version").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	
	return ""
}

// getCurrentGoVersionWithGobrew gets current version using gobrew
func getCurrentGoVersionWithGobrew() string {
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

// getCurrentGoVersionWithG gets current version using g
func getCurrentGoVersionWithG() string {
	// Try to get current version
	if currentCmd := exec.Command("g", "which"); currentCmd != nil {
		if currentOutput, currentErr := currentCmd.Output(); currentErr == nil {
			return strings.TrimSpace(string(currentOutput))
		}
	}
	
	// If that doesn't work, check symlink
	homeDir := GetHomeDir()
	goLink := filepath.Join(homeDir, ".g", "go")
	if target, err := os.Readlink(goLink); err == nil {
		return filepath.Base(target)
	}
	
	return ""
}

// FindGPath returns the path to the g version manager executable
func FindGPath() string {
	homeDir := GetHomeDir()
	gPaths := []string{
		filepath.Join(homeDir, ".g", "bin", "g"),
		filepath.Join(homeDir, "go", "bin", "g"),
		UsrLocalBinG,
	}
	
	for _, path := range gPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return ""
}
