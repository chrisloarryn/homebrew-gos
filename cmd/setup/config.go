package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// configureEnvironment sets up the environment variables for g  
func configureEnvironment() {
	homeDir := common.GetHomeDir()

	// Set environment variables
	setEnvironmentVariables(homeDir)

	yellow := color.New(color.FgYellow)
	yellow.Printf("  📝 Configuration added\n")
}

// generateGConfig generates the configuration script for g
func generateGConfig(homeDir, osName string) string {
	gPath := filepath.Join(homeDir, ".g")
	gBinPath := filepath.Join(gPath, "bin")

	config := fmt.Sprintf(`
# g (Go version manager) configuration
export G_HOME="%s"
export PATH="%s:$PATH"

# Go environment
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"
`, gPath, gBinPath)

	if osName == "Windows" {
		config = fmt.Sprintf(`
# g (Go version manager) configuration for Windows
set G_HOME=%s
set PATH=%s;%%PATH%%

# Go environment
set GOPATH=%%USERPROFILE%%\go
set PATH=%%GOPATH%%\bin;%%PATH%%
`, gPath, gBinPath)
	}

	return config
}

// setEnvironmentVariables adds g configuration to shell files
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
			if !common.HasConfigContent(file, "G_HOME") {
				common.AppendToFile(file, config)
			}
		}

	case "windows":
		// Windows - add to PowerShell profile
		profilePath := filepath.Join(homeDir, common.PowerShellProfile)
		if !common.HasConfigContent(profilePath, "G_HOME") {
			// Create directory if it doesn't exist
			os.MkdirAll(filepath.Dir(profilePath), 0755)
			common.AppendToFile(profilePath, config)
		}
	}
}

// createHelpScript creates a helper script with common commands
func createHelpScript() {
	homeDir := common.GetHomeDir()
	scriptPath := filepath.Join(homeDir, ".g", "go-help.sh")

	helpContent := `#!/bin/bash
# Go Version Manager Helper Script

echo "🚀 Go Version Manager Commands:"
echo ""
echo "📦 Installation:"
echo "   g install latest        # Install latest Go version"
echo "   g install 1.21.5        # Install specific version"
echo ""
echo "🔄 Version Management:"
echo "   g use latest            # Switch to latest version"
echo "   g use 1.21.5            # Switch to specific version"
echo "   g list                  # List installed versions"
echo "   g list-all              # List all available versions"
echo ""
echo "🗑️  Cleanup:"
echo "   g remove 1.21.5         # Remove specific version"
echo "   g clean                 # Clean cache"
echo ""
echo "📋 Information:"
echo "   g version               # Show g version"
echo "   go version              # Show current Go version"
echo "   which go                # Show Go path"
echo ""
echo "🛠️  gos Commands:"
echo "   gos status              # Show system status"
echo "   gos clean               # Deep clean"
echo "   gos env                 # Environment info"
echo ""
`

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
	gBin := filepath.Join(homeDir, ".g", "bin", "g")

	// Check if g binary exists
	if _, err := os.Stat(gBin); err == nil {
		green.Println("  ✅ 'g' binary found")
	} else {
		red.Println("  ❌ 'g' binary not found")
		return
	}

	// Check if g is working
	if common.IsCommandAvailable("g") {
		green.Println("  ✅ 'g' is available in PATH")
	} else {
		yellow.Println("  ⚠️  'g' not found in PATH (restart shell required)")
	}

	// Check if Go is installed
	if common.IsCommandAvailable("go") {
		green.Println("  ✅ Go is available")
	} else {
		yellow.Println("  ⚠️  Go not found (may need to install and use a version)")
	}
}
