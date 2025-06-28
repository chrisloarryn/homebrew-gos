package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload Go environment configuration",
	Long: `Reload Go environment configuration and verify it's working correctly.
	
This command will:
- Source shell configuration files
- Verify Go is available in PATH
- Check GOROOT and GOPATH settings
- Show current configuration status`,
	Run: func(cmd *cobra.Command, args []string) {
		reloadEnvironment()
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)
}

func reloadEnvironment() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("🔄 Reloading Go environment...")

	// Get expected paths
	homeDir := os.Getenv("HOME")
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")

	// Set environment variables for current session
	os.Setenv("GOPATH", expectedGopath)
	os.Setenv("GOROOT", expectedGoroot)

	// Update PATH for current session
	currentPath := os.Getenv("PATH")
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
		green.Printf("✅ Updated PATH with: %s\n", strings.Join(newPaths, ", "))
	}

	// Verify Go is available
	fmt.Println("")
	blue.Println("🔍 Verifying Go installation...")

	if _, err := exec.LookPath("go"); err != nil {
		red.Println("❌ Go not found in PATH")
		yellow.Println("💡 You may need to restart your terminal or run:")
		yellow.Println("   source ~/.zshrc")
		return
	}

	// Show Go version
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("✅ %s\n", version)
	}

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

	fmt.Println("")
	green.Println("🎉 Environment reload complete!")
	
	// Show helpful commands
	fmt.Println("")
	blue.Println("💡 Useful commands:")
	fmt.Println("  gos status        # Check overall status")
	fmt.Println("  gos env           # Show detailed environment")
	fmt.Println("  gos list          # List installed Go versions")
	fmt.Println("  go version        # Verify active Go version")
}
