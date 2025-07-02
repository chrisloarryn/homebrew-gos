package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// EnvironmentConfig holds the expected environment configuration
type EnvironmentConfig struct {
	ExpectedGoroot  string
	ExpectedGopath  string
	RequiredPaths   []string
	DirectoryChecks map[string]string
}

// ShowDetailedEnvironment displays detailed environment information
func ShowDetailedEnvironment() {
	blue := color.New(color.FgBlue)

	blue.Println("🌍 Go Environment Configuration")
	fmt.Println("")

	config := getEnvironmentConfig()

	checkGoEnvironmentVariables(config)
	checkPATHEntries(config)
	checkDirectories(config)

	fmt.Println("")
	fmt.Println("💡 Use 'gos env --fix' to automatically fix configuration issues")
}

// getEnvironmentConfig returns the expected environment configuration based on version manager preference
func getEnvironmentConfig() EnvironmentConfig {
	homeDir := common.GetHomeDir()

	// Prefer gobrew over .g in all systems
	if common.IsCommandAvailable("gobrew") {
		return getGobrewConfig(homeDir)
	}

	return getDefaultConfig(homeDir)
}

// getGobrewConfig returns configuration for gobrew on Windows
func getGobrewConfig(homeDir string) EnvironmentConfig {
	expectedGoroot := filepath.Join(homeDir, common.GobrewDir, "current", "go")
	expectedGopath := filepath.Join(homeDir, "go")

	return EnvironmentConfig{
		ExpectedGoroot: expectedGoroot,
		ExpectedGopath: expectedGopath,
		RequiredPaths: []string{
			filepath.Join(homeDir, common.GobrewDir, "bin"),
			filepath.Join(homeDir, common.GobrewDir, "current", "bin"),
			filepath.Join(homeDir, "go", "bin"),
		},
		DirectoryChecks: map[string]string{
			"GOPATH":           expectedGopath,
			"GOPATH bin":       filepath.Join(expectedGopath, "bin"),
			"gobrew directory": filepath.Join(homeDir, ".gobrew"),
			"gobrew bin":       filepath.Join(homeDir, ".gobrew", "bin"),
			"Go installation":  expectedGoroot,
		},
	}
}

// getDefaultConfig returns default configuration for g version manager
func getDefaultConfig(homeDir string) EnvironmentConfig {
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")

	return EnvironmentConfig{
		ExpectedGoroot: expectedGoroot,
		ExpectedGopath: expectedGopath,
		RequiredPaths: []string{
			filepath.Join(homeDir, ".g", "bin"),
			filepath.Join(homeDir, ".g", "go", "bin"),
			filepath.Join(homeDir, "go", "bin"),
		},
		DirectoryChecks: map[string]string{
			"GOPATH":          expectedGopath,
			"GOPATH bin":      filepath.Join(expectedGopath, "bin"),
			"g directory":     filepath.Join(homeDir, ".g"),
			"g bin directory": filepath.Join(homeDir, ".g", "bin"),
			"Go installation": expectedGoroot,
		},
	}
}

// checkGoEnvironmentVariables validates GOROOT and GOPATH settings
func checkGoEnvironmentVariables(config EnvironmentConfig) {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	// Check GOROOT
	actualGoroot := os.Getenv("GOROOT")
	if actualGoroot == config.ExpectedGoroot {
		green.Printf("✅ GOROOT: %s\n", actualGoroot)
	} else if actualGoroot == "" {
		red.Printf("❌ GOROOT: not set (should be: %s)\n", config.ExpectedGoroot)
	} else {
		yellow.Printf("⚠️  GOROOT: %s (expected: %s)\n", actualGoroot, config.ExpectedGoroot)
	}

	// Check GOPATH
	actualGopath := os.Getenv("GOPATH")
	if actualGopath == config.ExpectedGopath {
		green.Printf("✅ GOPATH: %s\n", actualGopath)
	} else if actualGopath == "" {
		red.Printf("❌ GOPATH: not set (should be: %s)\n", config.ExpectedGopath)
	} else {
		yellow.Printf("⚠️  GOPATH: %s (expected: %s)\n", actualGopath, config.ExpectedGopath)
	}
}

// checkPATHEntries validates required PATH entries
func checkPATHEntries(config EnvironmentConfig) {
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	path := os.Getenv("PATH")
	fmt.Println("\nPATH entries:")

	for _, reqPath := range config.RequiredPaths {
		if strings.Contains(path, reqPath) {
			green.Printf("✅ %s\n", reqPath)
		} else {
			red.Printf("❌ %s (missing)\n", reqPath)
		}
	}
}

// checkDirectories validates that required directories exist
func checkDirectories(config EnvironmentConfig) {
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	fmt.Println("\nDirectories:")

	for name, dir := range config.DirectoryChecks {
		if _, err := os.Stat(dir); err == nil {
			green.Printf("✅ %s: %s\n", name, dir)
		} else {
			red.Printf("❌ %s: %s (missing)\n", name, dir)
		}
	}
}
