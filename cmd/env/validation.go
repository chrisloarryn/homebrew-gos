package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/cristobalcontreras/gos/cmd/common"
)

// ValidateEnvironment runs comprehensive environment validation
func ValidateEnvironment() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("🔍 Comprehensive Environment Validation")
	fmt.Println("")

	hasErrors := false
	hasWarnings := false

	// Check basic environment variables
	homeDir := common.GetHomeDir()
	var expectedGoroot, expectedGopath string
	var requiredPaths []string

	if runtime.GOOS == "windows" {
		if common.IsCommandAvailable("gobrew") {
			expectedGoroot = filepath.Join(homeDir, ".gobrew", "current", "go")
			expectedGopath = filepath.Join(homeDir, "go")
			requiredPaths = []string{
				filepath.Join(homeDir, ".gobrew", "bin"),
				filepath.Join(homeDir, ".gobrew", "current", "bin"),
				filepath.Join(homeDir, "go", "bin"),
			}
		} else {
			expectedGoroot = filepath.Join(homeDir, ".g", "go")
			expectedGopath = filepath.Join(homeDir, "go")
			requiredPaths = []string{
				filepath.Join(homeDir, ".g", "bin"),
				filepath.Join(homeDir, ".g", "go", "bin"),
				filepath.Join(homeDir, "go", "bin"),
			}
		}
	} else {
		expectedGoroot = filepath.Join(homeDir, ".g", "go")
		expectedGopath = filepath.Join(homeDir, "go")
		requiredPaths = []string{
			filepath.Join(homeDir, ".g", "bin"),
			filepath.Join(homeDir, ".g", "go", "bin"),
			filepath.Join(homeDir, "go", "bin"),
		}
	}

	blue.Println("📋 Environment Variables:")
	
	// GOROOT validation
	actualGoroot := os.Getenv("GOROOT")
	if actualGoroot == expectedGoroot {
		green.Println("  ✅ GOROOT is correctly set")
	} else if actualGoroot == "" {
		red.Println("  ❌ GOROOT is not set")
		hasErrors = true
	} else {
		yellow.Printf("  ⚠️  GOROOT is set to %s (expected %s)\n", actualGoroot, expectedGoroot)
		hasWarnings = true
	}

	// GOPATH validation
	actualGopath := os.Getenv("GOPATH")
	if actualGopath == expectedGopath {
		green.Println("  ✅ GOPATH is correctly set")
	} else if actualGopath == "" {
		red.Println("  ❌ GOPATH is not set")
		hasErrors = true
	} else {
		yellow.Printf("  ⚠️  GOPATH is set to %s (expected %s)\n", actualGopath, expectedGopath)
		hasWarnings = true
	}

	// PATH validation
	fmt.Println("")
	blue.Println("🛤️  PATH Validation:")
	path := os.Getenv("PATH")

	pathMissing := 0
	for _, reqPath := range requiredPaths {
		if strings.Contains(path, reqPath) {
			green.Printf("  ✅ %s is in PATH\n", reqPath)
		} else {
			yellow.Printf("  ⚠️  %s is missing from PATH\n", reqPath)
			pathMissing++
		}
	}
	
	if pathMissing > 0 {
		fmt.Printf("    💡 Run 'gos setup' or 'gos env --fix' to add missing PATH entries\n")
		hasWarnings = true
	}

	// Directory structure validation
	fmt.Println("")
	blue.Println("📁 Directory Structure:")
	dirs := map[string]string{
		"GOPATH": expectedGopath,
		"GOPATH bin": filepath.Join(expectedGopath, "bin"),
	}

	// Only add version manager directories to critical check if we detect them
	if runtime.GOOS == "windows" && common.IsCommandAvailable("gobrew") {
		dirs["gobrew directory"] = filepath.Join(homeDir, ".gobrew")
		dirs["gobrew bin"] = filepath.Join(homeDir, ".gobrew", "bin")
		dirs["Go installation"] = expectedGoroot
	} else if common.IsGInstalled() {
		dirs["g directory"] = filepath.Join(homeDir, ".g")
		dirs["g bin directory"] = filepath.Join(homeDir, ".g", "bin")
		dirs["Go installation"] = expectedGoroot
	}

	for name, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			green.Printf("  ✅ %s exists: %s\n", name, dir)
		} else {
			if strings.Contains(name, "GOPATH") {
				yellow.Printf("  ⚠️  %s missing: %s\n", name, dir)
				hasWarnings = true
			} else {
				red.Printf("  ❌ %s missing: %s\n", name, dir)
				hasErrors = true
			}
		}
	}

	// Shell configuration validation
	fmt.Println("")
	blue.Println("🐚 Shell Configuration:")
	
	// Detect current shell
	currentShell := common.DetectCurrentShell()
	blue.Printf("  🔍 Detected shell: %s\n", currentShell)
	
	// Get shell file for current shell
	shellFile := common.GetShellFileForCurrentShell(currentShell, homeDir)
	
	if shellFile == "" {
		yellow.Println("  ⚠️  Could not determine shell configuration file")
		hasWarnings = true
	} else {
		fullPath := filepath.Join(homeDir, shellFile)
		if _, err := os.Stat(fullPath); err == nil {
			if hasGoConfig(fullPath) {
				green.Printf("  ✅ Go configuration found in %s\n", shellFile)
			} else {
				yellow.Printf("  ⚠️  %s exists but no Go configuration found\n", shellFile)
				fmt.Printf("    💡 Run 'gos setup' or 'gos env --fix' to add configuration\n")
				hasWarnings = true
			}
		} else {
			yellow.Printf("  ⚠️  Shell file %s does not exist\n", shellFile)
			fmt.Printf("    💡 Run 'gos setup' to create configuration\n")
			hasWarnings = true
		}
	}

	// Version manager validation
	fmt.Println("")
	blue.Println("🔧 Version Manager:")
	hasVersionManager := false
	if runtime.GOOS == "windows" && common.IsCommandAvailable("gobrew") {
		green.Println("  ✅ 'gobrew' version manager is available")
		hasVersionManager = true
		if versions := common.GetGobrewVersions(); len(versions) > 0 {
			green.Printf("  ✅ %d Go version(s) installed\n", len(versions))
		} else {
			yellow.Println("  ⚠️  No Go versions installed with gobrew")
			hasWarnings = true
		}
	} else if common.IsGInstalled() {
		green.Println("  ✅ 'g' version manager is available")
		hasVersionManager = true
		if versions := common.GetInstalledVersions(); len(versions) > 0 {
			green.Printf("  ✅ %d Go version(s) installed\n", len(versions))
		} else {
			yellow.Println("  ⚠️  No Go versions installed with g")
			hasWarnings = true
		}
	} else {
		yellow.Println("  ⚠️  No version manager found (gobrew or g)")
		fmt.Println("    💡 Run 'gos setup' to install a version manager")
		hasWarnings = true
	}

	// Go binary validation - only check if we have a version manager
	fmt.Println("")
	blue.Println("🐹 Go Binary:")
	if hasVersionManager {
		if goPath, err := exec.LookPath("go"); err == nil {
			green.Printf("  ✅ Go binary found: %s\n", goPath)
			
			// Check if go version works
			if output, err := exec.Command("go", "version").Output(); err == nil {
				version := strings.TrimSpace(string(output))
				green.Printf("  ✅ Go version: %s\n", version)
			} else {
				yellow.Println("  ⚠️  Go binary exists but 'go version' failed")
				hasWarnings = true
			}
		} else {
			yellow.Println("  ⚠️  Go binary not found in PATH")
			fmt.Println("    💡 Install a Go version with 'gos install latest'")
			hasWarnings = true
		}
	} else {
		yellow.Println("  ℹ️  Skipping Go binary check (no version manager)")
	}

	// Summary
	fmt.Println("")
	blue.Println("📊 Validation Summary:")
	if hasErrors {
		red.Println("  ❌ Environment has critical issues that need fixing")
		fmt.Println("  💡 Run 'gos env --fix' to attempt automatic fixes")
	} else if hasWarnings {
		yellow.Println("  ⚠️  Environment has minor issues")
		fmt.Println("  💡 Consider running 'gos env --fix' to optimize configuration")
	} else {
		green.Println("  ✅ Environment is properly configured!")
	}
}
