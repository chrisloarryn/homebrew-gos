package common

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// GetSystemGoInfo returns information about system Go installation
func GetSystemGoInfo() (version string, goroot string, found bool) {
	// Check if Go is installed directly
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version = strings.TrimSpace(string(output))
		
		// Try to get GOROOT to see where it's installed
		if gorootOutput, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
			goroot = strings.TrimSpace(string(gorootOutput))
		}
		
		return version, goroot, true
	}
	
	return "", "", false
}

// SetupGoEnvironment sets up the Go environment variables and PATH
func SetupGoEnvironment() {
	green := color.New(color.FgGreen)
	
	homeDir := GetHomeDir()
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")

	// Set environment variables for current session
	os.Setenv("GOPATH", expectedGopath)
	os.Setenv("GOROOT", expectedGoroot)

	// Update PATH for current session
	if UpdatePathForGoEnvironment() {
		green.Println("✅ PATH updated for current session")
	}
}

// VerifyGoInstallation verifies that Go is available and shows version
func VerifyGoInstallation() bool {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	if _, err := exec.LookPath("go"); err != nil {
		red.Println("❌ Go not found in PATH")
		yellow.Println("💡 You may need to restart your terminal or run:")
		yellow.Println("   source ~/.zshrc")
		return false
	}

	// Show Go version
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("✅ %s\n", version)
	}
	
	return true
}

// VerifyGoEnvironmentPaths verifies GOROOT and GOPATH settings
func VerifyGoEnvironmentPaths() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	
	homeDir := GetHomeDir()
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")

	// Verify GOROOT
	if output, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		goroot := strings.TrimSpace(string(output))
		if goroot == expectedGoroot {
			green.Printf("✅ GOROOT: %s\n", goroot)
		} else {
			yellow.Printf("⚠️  GOROOT: %s (expected: %s)\n", goroot, expectedGoroot)
		}
	}

	// Verify GOPATH
	if output, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		gopath := strings.TrimSpace(string(output))
		if gopath == expectedGopath {
			green.Printf("✅ GOPATH: %s\n", gopath)
		} else {
			yellow.Printf("⚠️  GOPATH: %s (expected: %s)\n", gopath, expectedGopath)
		}
	}
}

// DisplayCurrentGoVersion displays the current Go version with GOROOT and GOPATH
func DisplayCurrentGoVersion() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("📍 Current Go version:")

	goCmd := exec.Command("go", "version")
	if err := goCmd.Run(); err != nil {
		yellow.Println("⚠️  Go is not available in PATH")
		return
	}

	goCmd.Stdout = os.Stdout
	goCmd.Run()

	// Show GOROOT and GOPATH
	if gorootCmd := exec.Command("go", "env", "GOROOT"); gorootCmd.Run() == nil {
		blue.Print("📂 GOROOT: ")
		gorootCmd.Stdout = os.Stdout
		gorootCmd.Run()
	}

	if gopathCmd := exec.Command("go", "env", "GOPATH"); gopathCmd.Run() == nil {
		blue.Print("📂 GOPATH: ")
		gopathCmd.Stdout = os.Stdout
		gopathCmd.Run()
	}
}

// GetGoVersion returns the current Go version as a string
func GetGoVersion() string {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "Go not available"
	}
	return strings.TrimSpace(string(output))
}

// GetGoEnvVar returns a specific Go environment variable
func GetGoEnvVar(envVar string) string {
	cmd := exec.Command("go", "env", envVar)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
