package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// FixEnvironment fixes environment configuration issues
func FixEnvironment() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	blue.Println("🔧 Fixing Go environment configuration...")

	homeDir := os.Getenv("HOME")
	expectedGopath := filepath.Join(homeDir, "go")

	// Create GOPATH directory if it doesn't exist
	if err := os.MkdirAll(filepath.Join(expectedGopath, "bin"), 0755); err == nil {
		green.Printf("✅ Created GOPATH directory: %s\n", expectedGopath)
	}

	// Add to shell configuration
	shellFiles := []string{
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bashrc"),
	}

	for _, shellFile := range shellFiles {
		if _, err := os.Stat(shellFile); err == nil {
			yellow.Printf("💡 Please add this to %s and restart your shell:\n", shellFile)
			break
		}
	}

	ExportEnvironment()
	
	fmt.Println("")
	yellow.Println("📋 To apply changes immediately, run:")
	fmt.Println("  source ~/.zshrc")
	fmt.Println("  # or")
	fmt.Println("  eval $(gos env --export)")
}

// ExportEnvironment exports environment variables for sourcing
func ExportEnvironment() {
	homeDir := os.Getenv("HOME")
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")
	
	fmt.Printf("export GOROOT=%s\n", expectedGoroot)
	fmt.Printf("export GOPATH=%s\n", expectedGopath)
	
	// Build PATH
	requiredPaths := []string{
		filepath.Join(homeDir, ".g", "bin"),
		filepath.Join(homeDir, ".g", "go", "bin"),
		filepath.Join(homeDir, "go", "bin"),
	}
	
	currentPath := os.Getenv("PATH")
	newPaths := []string{}
	
	for _, reqPath := range requiredPaths {
		if !strings.Contains(currentPath, reqPath) {
			newPaths = append(newPaths, reqPath)
		}
	}
	
	if len(newPaths) > 0 {
		fmt.Printf("export PATH=%s:$PATH\n", strings.Join(newPaths, ":"))
	}
}
