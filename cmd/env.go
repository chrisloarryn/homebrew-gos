package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show or fix environment configuration",
	Long: `Show current environment configuration and optionally fix it.
	
Examples:
  gos env          # Show current environment
  gos env --fix    # Fix environment configuration
  gos env --export # Export current environment for sourcing`,
	Run: func(cmd *cobra.Command, args []string) {
		fix, _ := cmd.Flags().GetBool("fix")
		export, _ := cmd.Flags().GetBool("export")

		if export {
			exportEnvironment()
		} else if fix {
			fixEnvironment()
		} else {
			showDetailedEnvironment()
		}
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.Flags().Bool("fix", false, "Fix environment configuration")
	envCmd.Flags().Bool("export", false, "Export environment variables for sourcing")
}

func showDetailedEnvironment() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("🌍 Go Environment Configuration")
	fmt.Println("")

	// Expected values
	homeDir := os.Getenv("HOME")
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")

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
	requiredPaths := []string{
		filepath.Join(homeDir, ".g", "bin"),
		filepath.Join(homeDir, ".g", "go", "bin"),
		filepath.Join(homeDir, "go", "bin"),
	}

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
		"g directory": filepath.Join(homeDir, ".g"),
		"g bin directory": filepath.Join(homeDir, ".g", "bin"),
		"Go installation": expectedGoroot,
		"GOPATH": expectedGopath,
		"GOPATH bin": filepath.Join(expectedGopath, "bin"),
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

func fixEnvironment() {
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

	exportEnvironment()
	
	fmt.Println("")
	yellow.Println("📋 To apply changes immediately, run:")
	fmt.Println("  source ~/.zshrc")
	fmt.Println("  # or")
	fmt.Println("  eval $(gos env --export)")
}

func exportEnvironment() {
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
