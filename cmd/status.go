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
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	if _, err := exec.LookPath("go"); err != nil {
		yellow.Println("  ⚠️  Go is not available in PATH")
		yellow.Println("  💡 Try running: source ~/.zshrc")
		return
	}

	// Show version with better formatting
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("  ✅ %s\n", version)
	}

	// Show GOROOT with validation
	if output, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		goroot := strings.TrimSpace(string(output))
		expectedGoroot := filepath.Join(os.Getenv("HOME"), ".g", "go")
		if goroot == expectedGoroot {
			green.Printf("  ✅ GOROOT: %s\n", goroot)
		} else {
			blue.Printf("  ℹ️  GOROOT: %s\n", goroot)
		}
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
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	
	// Expected values
	expectedGoroot := filepath.Join(os.Getenv("HOME"), ".g", "go")
	expectedGopath := filepath.Join(os.Getenv("HOME"), "go")
	
	envVars := map[string]string{
		"GOROOT": expectedGoroot,
		"GOPATH": expectedGopath,
		"GOPROXY": "",
		"GOSUMDB": "",
		"GOMODCACHE": "",
	}
	
	for envVar, expected := range envVars {
		if value := os.Getenv(envVar); value != "" {
			if expected != "" && value == expected {
				green.Printf("  ✅ %s: %s\n", envVar, value)
			} else if expected != "" {
				yellow.Printf("  ⚠️  %s: %s (expected: %s)\n", envVar, value, expected)
			} else {
				blue.Printf("  ℹ️  %s: %s\n", envVar, value)
			}
		} else {
			if expected != "" {
				yellow.Printf("  ❌ %s: (not set, should be: %s)\n", envVar, expected)
			} else {
				fmt.Printf("  %s: (not set)\n", envVar)
			}
		}
	}

	// Show PATH entries related to Go
	fmt.Println("  PATH (Go-related entries):")
	path := os.Getenv("PATH")
	pathEntries := strings.Split(path, ":")
	hasGoBin := false
	hasGBin := false
	
	for _, entry := range pathEntries {
		if strings.Contains(entry, "go") || strings.Contains(entry, ".g") {
			if strings.Contains(entry, ".g/bin") {
				hasGBin = true
				green.Printf("    ✅ %s\n", entry)
			} else if strings.Contains(entry, "go/bin") {
				hasGoBin = true
				green.Printf("    ✅ %s\n", entry)
			} else {
				blue.Printf("    ℹ️  %s\n", entry)
			}
		}
	}
	
	if !hasGBin {
		yellow.Println("    ⚠️  ~/.g/bin not found in PATH")
	}
	if !hasGoBin {
		yellow.Println("    ⚠️  $GOPATH/bin not found in PATH")
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
