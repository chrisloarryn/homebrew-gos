package common

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// GetHomeDir gets the user's home directory
func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

// UpdatePathForVersionManager updates the PATH environment variable to include the current Go version
func UpdatePathForVersionManager() {
	homeDir := GetHomeDir()
	
	if IsCommandAvailable("gobrew") {
		UpdatePathForGobrew(homeDir)
	} else if IsCommandAvailable("g") {
		UpdatePathForG(homeDir)
	}
}

// UpdatePathForGobrew updates PATH for gobrew version manager
func UpdatePathForGobrew(homeDir string) {
	currentPath := os.Getenv("PATH")
	gobrewCurrentBin := filepath.Join(homeDir, GobrewDir, "current", "bin")
	gobrewBin := filepath.Join(homeDir, GobrewDir, "bin")
	
	pathSep := ";"
	if runtime.GOOS != "windows" {
		pathSep = ":"
	}

	// Remove any existing gobrew paths from PATH to avoid duplicates
	pathParts := strings.Split(currentPath, pathSep)
	var cleanPaths []string

	for _, part := range pathParts {
		if !strings.Contains(part, GobrewDir) {
			cleanPaths = append(cleanPaths, part)
		}
	}

	// Prepend gobrew paths
	finalPaths := []string{gobrewCurrentBin, gobrewBin}
	finalPaths = append(finalPaths, cleanPaths...)
	newPath := strings.Join(finalPaths, pathSep)

	os.Setenv("PATH", newPath)
}

// UpdatePathForG updates PATH for g version manager
func UpdatePathForG(homeDir string) {
	currentPath := os.Getenv("PATH")
	gGoBin := filepath.Join(homeDir, ".g", "go", "bin")
	gBin := filepath.Join(homeDir, ".g", "bin")
	
	pathSep := ";"
	if runtime.GOOS != "windows" {
		pathSep = ":"
	}

	// Remove any existing g paths from PATH to avoid duplicates
	pathParts := strings.Split(currentPath, pathSep)
	var cleanPaths []string

	for _, part := range pathParts {
		if !strings.Contains(part, ".g") {
			cleanPaths = append(cleanPaths, part)
		}
	}

	// Prepend g paths
	finalPaths := []string{gGoBin, gBin}
	finalPaths = append(finalPaths, cleanPaths...)
	newPath := strings.Join(finalPaths, pathSep)

	os.Setenv("PATH", newPath)
}

// UpdatePathForGoEnvironment updates PATH with required Go paths
func UpdatePathForGoEnvironment() bool {
	currentPath := os.Getenv("PATH")
	homeDir := GetHomeDir()
	
	requiredPaths := []string{
		filepath.Join(homeDir, ".g", "bin"),
		filepath.Join(homeDir, ".g", "go", "bin"),
		filepath.Join(homeDir, "go", "bin"),
	}

	newPaths := []string{}
	for _, reqPath := range requiredPaths {
		if !strings.Contains(currentPath, reqPath) {
			newPaths = append(newPaths, reqPath)
		}
	}

	if len(newPaths) > 0 {
		newPath := strings.Join(newPaths, ":") + ":" + currentPath
		os.Setenv("PATH", newPath)
		green := color.New(color.FgGreen)
		green.Printf("✅ Updated PATH with: %s\n", strings.Join(newPaths, ", "))
		return true
	}
	
	return false
}

// GetExpectedGoRoot returns the expected GOROOT path
func GetExpectedGoRoot() string {
	return filepath.Join(GetHomeDir(), ".g", "go")
}

// GetExpectedGoPath returns the expected GOPATH path
func GetExpectedGoPath() string {
	return filepath.Join(GetHomeDir(), "go")
}
