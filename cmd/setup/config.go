package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// configureEnvironment sets up the environment variables for gobrew
func configureEnvironment() {
	homeDir := common.GetHomeDir()

	// Set environment variables
	setEnvironmentVariables(homeDir)

	yellow := color.New(color.FgYellow)
	yellow.Printf("  📝 Configuration added\n")
}

// generateGobrewConfig generates the configuration script for gobrew
func generateGConfig(homeDir, osName string) string {
	gobrewPath := filepath.Join(homeDir, ".gobrew")
	gobrewBinPath := filepath.Join(gobrewPath, "bin")
	gobrewCurrentBinPath := filepath.Join(gobrewPath, "current", "bin")

	config := fmt.Sprintf(`
# gobrew (Go version manager) configuration
export PATH="%s:%s:$PATH"

# Go environment
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"
`, gobrewCurrentBinPath, gobrewBinPath)

	if osName == "Windows" {
		config = fmt.Sprintf(`
# gobrew (Go version manager) configuration for Windows
set PATH=%s;%s;%%PATH%%

# Go environment
set GOPATH=%%USERPROFILE%%\go
set PATH=%%GOPATH%%\bin;%%PATH%%
`, gobrewCurrentBinPath, gobrewBinPath)
	}

	return config
}

// setEnvironmentVariables adds gobrew configuration to shell files
func setEnvironmentVariables(homeDir string) {
	config := generateGConfig(homeDir, runtime.GOOS)

	switch runtime.GOOS {
	case "darwin", "linux":
		// Unix-like systems
		shellFiles := []string{
			filepath.Join(homeDir, common.ZshrcFile),
			filepath.Join(homeDir, common.BashrcFile),
			filepath.Join(homeDir, common.BashProfileFile),
		}

		for _, file := range shellFiles {
			if !common.HasConfigContent(file, "gobrew") {
				common.AppendToFile(file, config)
			}
		}

	case "windows":
		// Windows - add to PowerShell profile
		profilePath := filepath.Join(homeDir, common.PowerShellProfile)
		if !common.HasConfigContent(profilePath, "gobrew") {
			// Create directory if it doesn't exist
			os.MkdirAll(filepath.Dir(profilePath), 0755)
			common.AppendToFile(profilePath, config)
		}
	}
}

// createHelpScript creates a helper script with common commands
func createHelpScript() {
	homeDir := common.GetHomeDir()
	scriptPath := filepath.Join(homeDir, ".gobrew", "gobrew-help.sh")

	helpContent := `#!/bin/bash
# Go Version Manager Helper Script

echo "🚀 Go Version Manager Commands:"
echo ""
echo "📦 Installation:"
echo "   gobrew install latest   # Install latest Go version"
echo "   gobrew install 1.21.5   # Install specific version"
echo ""
echo "🔄 Version Management:"
echo "   gobrew use latest       # Switch to latest version"
echo "   gobrew use 1.21.5       # Switch to specific version"
echo "   gobrew ls               # List installed versions"
echo "   gobrew ls-remote        # List all available versions"
echo ""
echo "🗑️  Cleanup:"
echo "   gobrew uninstall 1.21.5 # Remove specific version"
echo "   gobrew clean            # Clean cache"
echo ""
echo "📋 Information:"
echo "   gobrew --version        # Show gobrew version"
echo "   go version              # Show current Go version"
echo "   which go                # Show Go path"
echo ""
echo "🛠️  gos Commands:"
echo "   gos status              # Show system status"
echo "   gos clean               # Deep clean"
echo "   gos env                 # Environment info"
echo ""
`

	// Create directory if it doesn't exist
	os.MkdirAll(filepath.Join(homeDir, ".gobrew"), 0755)
	common.WriteToFile(scriptPath, helpContent)

	// Make executable on Unix-like systems
	if runtime.GOOS != "windows" {
		os.Chmod(scriptPath, 0755)
	}
}

// verifyInstallation checks if the installation was successful
func verifyInstallation() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	homeDir := common.GetHomeDir()
	gobrewBin := filepath.Join(homeDir, ".gobrew", "bin", "gobrew")

	// Check if gobrew binary exists
	if _, err := os.Stat(gobrewBin); err == nil {
		green.Println("  ✅ 'gobrew' binary found")
	} else {
		red.Println("  ❌ 'gobrew' binary not found")
		return
	}

	// Check if gobrew is working
	if common.IsCommandAvailable("gobrew") {
		green.Println("  ✅ 'gobrew' is available in PATH")
	} else {
		yellow.Println("  ⚠️  'gobrew' not found in PATH (restart shell required)")
	}

	// Check if Go is installed
	if common.IsCommandAvailable("go") {
		green.Println("  ✅ Go is available")
	} else {
		yellow.Println("  ⚠️  Go not found (may need to install and use a version)")
	}
}
