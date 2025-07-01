package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/cristobalcontreras/gos/cmd/common"
)

// ShowDetailedEnvironment displays detailed environment information
func ShowDetailedEnvironment() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("🌍 Go Environment Configuration")
	fmt.Println("")

	// Expected values based on OS and version manager
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

	// Check GOROOT
	actualGoroot := os.Getenv("GOROOT")
	if actualGoroot == expectedGoroot {
		green.Printf("✅ GOROOT: %s\n", actualGoroot)
	} else if actualGoroot == "" {
		red.Printf("❌ GOROOT: not set (should be: %s)\n", expectedGoroot)
	} else {
		yellow.Printf("⚠️  GOROOT: %s (expected: %s)\n", actualGoroot, expectedGoroot)
	}

	// Check GOPATH
	actualGopath := os.Getenv("GOPATH")
	if actualGopath == expectedGopath {
		green.Printf("✅ GOPATH: %s\n", actualGopath)
	} else if actualGopath == "" {
		red.Printf("❌ GOPATH: not set (should be: %s)\n", expectedGopath)
	} else {
		yellow.Printf("⚠️  GOPATH: %s (expected: %s)\n", actualGopath, expectedGopath)
	}

	// Check PATH
	path := os.Getenv("PATH")

	fmt.Println("\nPATH entries:")
	for _, reqPath := range requiredPaths {
		if strings.Contains(path, reqPath) {
			green.Printf("✅ %s\n", reqPath)
		} else {
			red.Printf("❌ %s (missing)\n", reqPath)
		}
	}

	// Check if directories exist
	fmt.Println("\nDirectories:")
	dirs := map[string]string{
		"GOPATH": expectedGopath,
		"GOPATH bin": filepath.Join(expectedGopath, "bin"),
	}

	if runtime.GOOS == "windows" && common.IsCommandAvailable("gobrew") {
		dirs["gobrew directory"] = filepath.Join(homeDir, ".gobrew")
		dirs["gobrew bin"] = filepath.Join(homeDir, ".gobrew", "bin")
		dirs["Go installation"] = expectedGoroot
	} else {
		dirs["g directory"] = filepath.Join(homeDir, ".g")
		dirs["g bin directory"] = filepath.Join(homeDir, ".g", "bin")
		dirs["Go installation"] = expectedGoroot
	}

	for name, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			green.Printf("✅ %s: %s\n", name, dir)
		} else {
			red.Printf("❌ %s: %s (missing)\n", name, dir)
		}
	}

	fmt.Println("")
	fmt.Println("💡 Use 'gos env --fix' to automatically fix configuration issues")
}
