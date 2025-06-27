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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Go system status",
	Long: `Show comprehensive information about the current Go installation,
including version manager status, installed versions, and system configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		showStatus()
	},
}

func showStatus() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	blue.Println("📊 Go system status:")
	fmt.Println("")

	// Check 'g' manager
	blue.Println("🔧 'g' Manager:")
	if _, err := exec.LookPath("g"); err == nil {
		green.Print("  ✅ Installed: ")
		if versionCmd := exec.Command("g", "--version"); versionCmd.Run() == nil {
			versionCmd.Stdout = os.Stdout
			versionCmd.Run()
		} else {
			fmt.Println("unknown version")
		}
	} else {
		red.Println("  ❌ Not installed")
		yellow.Println("  💡 Run: gos setup")
	}

	fmt.Println("")

	// Current Go installation
	blue.Println("🐹 Current Go:")
	showCurrentGo()

	fmt.Println("")

	// Installed versions
	blue.Println("📦 Installed versions:")
	if _, err := exec.LookPath("g"); err == nil {
		listVersions()
	} else {
		yellow.Println("  'g' manager not installed")
	}

	fmt.Println("")

	// Disk space
	blue.Println("💾 Disk space:")
	showDiskUsage()

	fmt.Println("")

	// Environment variables
	blue.Println("🌍 Environment:")
	showEnvironment()

	fmt.Println("")

	// Project configuration
	blue.Println("📁 Project configuration:")
	showProjectConfig()
}

func showCurrentGo() {
	yellow := color.New(color.FgYellow)

	if _, err := exec.LookPath("go"); err != nil {
		yellow.Println("  ⚠️  Go is not available in PATH")
		return
	}

	// Show version
	fmt.Print("  Version: ")
	versionCmd := exec.Command("go", "version")
	versionCmd.Stdout = os.Stdout
	versionCmd.Run()

	// Show GOROOT
	if output, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		fmt.Printf("  GOROOT: %s", string(output))
	}

	// Show GOPATH
	if output, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		fmt.Printf("  GOPATH: %s", string(output))
	}
}

func showDiskUsage() {
	homeDir := os.Getenv("HOME")
	gDir := filepath.Join(homeDir, ".g")

	if _, err := os.Stat(gDir); err == nil {
		if output, err := exec.Command("du", "-sh", gDir).Output(); err == nil {
			fmt.Printf("  ~/.g directory: %s", string(output))
		} else {
			fmt.Println("  Could not calculate ~/.g directory size")
		}
	} else {
		fmt.Println("  ~/.g directory not found")
	}

	// Show Go workspace size if it exists
	goDir := filepath.Join(homeDir, "go")
	if _, err := os.Stat(goDir); err == nil {
		if output, err := exec.Command("du", "-sh", goDir).Output(); err == nil {
			fmt.Printf("  ~/go directory: %s", string(output))
		}
	}
}

func showEnvironment() {
	envVars := []string{"GOROOT", "GOPATH", "GOPROXY", "GOSUMDB", "GOMODCACHE"}
	
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			fmt.Printf("  %s: %s\n", envVar, value)
		} else {
			fmt.Printf("  %s: (not set)\n", envVar)
		}
	}

	// Show PATH entries related to Go
	fmt.Println("  PATH (Go-related entries):")
	path := os.Getenv("PATH")
	pathEntries := strings.Split(path, ":")
	
	for _, entry := range pathEntries {
		if strings.Contains(entry, "go") || strings.Contains(entry, ".g") {
			fmt.Printf("    %s\n", entry)
		}
	}
}

func showProjectConfig() {
	// Check for .go-version file
	if _, err := os.Stat(".go-version"); err == nil {
		if content, err := os.ReadFile(".go-version"); err == nil {
			version := strings.TrimSpace(string(content))
			color.Green("  ✅ .go-version found: %s", version)
		}
	} else {
		color.Yellow("  ℹ️  No .go-version file in current directory")
	}

	// Check for go.mod
	if _, err := os.Stat("go.mod"); err == nil {
		if content, err := os.ReadFile("go.mod"); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "go ") {
					version := strings.TrimPrefix(line, "go ")
					color.Green("  ✅ go.mod found, Go version: %s", version)
					break
				}
			}
		}
	} else {
		color.Yellow("  ℹ️  No go.mod file in current directory")
	}
}
