package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the 'g' Go version manager",
	Long: `Install and configure the 'g' Go version manager.
This will download and install 'g', configure environment variables,
and install the latest stable Go version.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupGoVersionManager()
	},
}

func setupGoVersionManager() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("🔧 Installing 'g' version manager for Go...")

	// Detect OS and architecture
	osName := runtime.GOOS
	arch := runtime.GOARCH

	switch osName {
	case "darwin":
		osName = "macOS"
	case "linux":
		osName = "Linux"
	case "windows":
		osName = "Windows"
	}

	if arch == "arm64" {
		if osName == "macOS" {
			fmt.Println("  Detected: Apple Silicon (M1/M2/M3)")
		} else {
			fmt.Println("  Detected: ARM64")
		}
	} else if arch == "amd64" {
		fmt.Println("  Detected: Intel x86_64")
	} else {
		fmt.Printf("  Detected: %s on %s\n", arch, osName)
	}

	blue.Println("\n▸ Downloading and installing 'g'...")

	// Create directory for g
	gDir := filepath.Join(os.Getenv("HOME"), ".g")
	if err := os.MkdirAll(gDir, 0755); err != nil {
		red.Printf("❌ Error creating .g directory: %v\n", err)
		return
	}

	// Try to install g using the install script
	if !installGWithScript() {
		yellow.Println("  ❌ Error installing 'g'. Trying alternative method...")
		if !installGManually() {
			red.Println("  ❌ Failed to install 'g'")
			return
		}
	}

	green.Println("  ✅ 'g' installed successfully")

	blue.Println("\n▸ Configuring PATH and environment variables...")

	// Configure environment variables
	configureEnvironment()

	blue.Println("\n▸ Installing latest stable Go version...")

	// Install latest Go version
	gBin := filepath.Join(os.Getenv("HOME"), ".g", "bin", "g")
	installCmd := exec.Command(gBin, "install", "latest")
	if err := installCmd.Run(); err != nil {
		yellow.Println("  ℹ️  Installing known specific version...")
		fallbackCmd := exec.Command(gBin, "install", "1.21.5")
		fallbackCmd.Run()
	} else {
		green.Println("  ✅ Go latest installed successfully")
	}

	// Create help script
	createHelpScript()

	green.Println("\n✅ Installation completed!")
	fmt.Println("")
	yellow.Println("📋 Next steps:")

	if osName == "Windows" {
		fmt.Println("1. Run: source ~/.bashrc  (or restart Git Bash/WSL)")
	} else {
		fmt.Println("1. Run: source ~/.zshrc  (or open a new terminal)")
	}

	fmt.Println("2. Verify: g --version")
	fmt.Println("3. Use: gos list  (to see installed versions)")
	fmt.Println("")
	yellow.Println("💡 To see all available commands:")
	fmt.Println("   ~/.g/go-help.sh")
	fmt.Println("")
	blue.Println("🚀 Quick examples:")
	fmt.Println("   gos install 1.21.5     # Install Go 1.21.5")
	fmt.Println("   gos use 1.21.5         # Switch to Go 1.21.5")
	fmt.Println("   gos list               # View installed versions")
}

func installGWithScript() bool {
	cmd := exec.Command("bash", "-c", "curl -sSL https://git.io/g-install | bash -s -- -y")
	return cmd.Run() == nil
}

func installGManually() bool {
	// Try git clone method
	if _, err := exec.LookPath("git"); err == nil {
		tmpDir := "/tmp/g"
		cloneCmd := exec.Command("git", "clone", "https://github.com/stefanmaric/g.git", tmpDir)
		if cloneCmd.Run() == nil {
			makeCmd := exec.Command("make", "install", fmt.Sprintf("PREFIX=%s", filepath.Join(os.Getenv("HOME"), ".g")))
			makeCmd.Dir = tmpDir
			if makeCmd.Run() == nil {
				os.RemoveAll(tmpDir)
				return true
			}
		}
		os.RemoveAll(tmpDir)
	}

	// Fallback: download directly
	gBinDir := filepath.Join(os.Getenv("HOME"), ".g", "bin")
	if err := os.MkdirAll(gBinDir, 0755); err != nil {
		return false
	}

	downloadCmd := exec.Command("curl", "-sSL", "https://raw.githubusercontent.com/stefanmaric/g/main/bin/g", "-o", filepath.Join(gBinDir, "g"))
	if downloadCmd.Run() != nil {
		return false
	}

	chmodCmd := exec.Command("chmod", "+x", filepath.Join(gBinDir, "g"))
	return chmodCmd.Run() == nil
}

func configureEnvironment() {
	osName := runtime.GOOS
	var shellFiles []string

	if osName == "windows" {
		shellFiles = []string{
			filepath.Join(os.Getenv("HOME"), ".bashrc"),
			filepath.Join(os.Getenv("HOME"), ".bash_profile"),
		}
	} else {
		shellFiles = []string{
			filepath.Join(os.Getenv("HOME"), ".zshrc"),
			filepath.Join(os.Getenv("HOME"), ".bashrc"),
			filepath.Join(os.Getenv("HOME"), ".bash_profile"),
		}
	}

	gConfig := `
# === Go Version Manager (g) ===
export GOPATH=$HOME/go
export GOROOT=$HOME/.g/go
export PATH=$HOME/.g/bin:$GOROOT/bin:$GOPATH/bin:$PATH
`

	configAdded := false
	for _, shellFile := range shellFiles {
		if _, err := os.Stat(shellFile); err == nil || shellFile == shellFiles[0] {
			if !hasGConfig(shellFile) {
				if err := appendToFile(shellFile, gConfig); err == nil {
					color.Green("  ✅ Configuration added to %s", shellFile)
					configAdded = true
				}
			} else {
				color.Yellow("  ℹ️  Configuration already exists in %s", shellFile)
				configAdded = true
			}
			break
		}
	}

	if !configAdded {
		defaultShell := shellFiles[0]
		if err := writeToFile(defaultShell, gConfig); err == nil {
			color.Green("  ✅ Configuration created in %s", defaultShell)
		}
	}

	// Export for current session
	os.Setenv("GOPATH", filepath.Join(os.Getenv("HOME"), "go"))
	os.Setenv("GOROOT", filepath.Join(os.Getenv("HOME"), ".g", "go"))
	currentPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s:%s:%s:%s",
		filepath.Join(os.Getenv("HOME"), ".g", "bin"),
		filepath.Join(os.Getenv("HOME"), ".g", "go", "bin"),
		filepath.Join(os.Getenv("HOME"), "go", "bin"),
		currentPath)
	os.Setenv("PATH", newPath)
}

func hasGConfig(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Go Version Manager (g)") {
			return true
		}
	}
	return false
}

func appendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func writeToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func createHelpScript() {
	helpScript := `#!/bin/bash
# Useful commands for the 'g' version manager

echo "🐹 Go Version Manager - Useful commands:"
echo ""
echo "📦 Installation:"
echo "  gos install latest        # Install latest version"
echo "  gos install 1.21.5        # Install specific version"
echo "  gos install 1.20.x        # Install latest 1.20.x"
echo ""
echo "🔄 Version switching:"
echo "  gos use 1.21.5            # Switch to specific version"
echo "  gos use latest            # Switch to latest installed"
echo ""
echo "📋 Information:"
echo "  gos list                  # List installed versions"
echo "  gos list --remote         # List all available versions"
echo "  gos status                # Show current version and status"
echo "  go version                # Confirm active Go version"
echo ""
echo "🗑️  Cleanup:"
echo "  gos remove 1.20.10        # Remove specific version"
echo "  gos clean                 # Deep clean all Go installations"
echo ""
echo "💡 Usage examples:"
echo "  gos install 1.21.5 && gos use 1.21.5"
echo "  gos project 1.21.5       # Set version for current project"
echo ""
`

	helpFile := filepath.Join(os.Getenv("HOME"), ".g", "go-help.sh")
	if err := os.WriteFile(helpFile, []byte(helpScript), 0755); err == nil {
		color.Blue("\n▸ Creating help script...")
	}
}
