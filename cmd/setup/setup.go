package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewSetupCmd creates the setup command
func NewSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup the 'g' Go version manager",
		Long: `Install and configure the 'g' Go version manager.
This will download and install 'g', configure environment variables,
and install the latest stable Go version.`,
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")
			setupGoVersionManager(force)
		},
	}
	
	cmd.Flags().BoolP("force", "f", false, "Force reinstallation even if version managers are already installed")
	return cmd
}

func setupGoVersionManager(force bool) {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("🔧 Setting up Go version manager...")

	// Check if any version manager is already installed (unless force is used)
	if !force && checkExistingInstallations() {
		return
	}

	if force {
		yellow.Println("⚡ Force flag detected - proceeding with installation...")
	}

	// Detect and display system information
	displaySystemInfo()

	// Handle Windows separately
	if runtime.GOOS == "windows" {
		handleWindowsSetup()
		return
	}

	// Unix-like systems setup
	if !performUnixSetup() {
		red.Println("❌ Setup failed")
		return
	}

	// Complete setup
	completeSetup()
}

// displaySystemInfo shows detected OS and architecture
func displaySystemInfo() {
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
}

// handleWindowsSetup manages Windows-specific setup
func handleWindowsSetup() {
	yellow := color.New(color.FgYellow)
	
	yellow.Println("\n⚠️  Windows detected.")
	yellow.Println("   The original 'g' version manager doesn't support Windows.")
	yellow.Println("   🚀 Using Windows-compatible alternatives...")

	fmt.Print("\n   Continue with Windows setup? (Y/n): ")
	var response string
	fmt.Scanln(&response)
	if response == "n" || response == "N" {
		yellow.Println("Installation cancelled.")
		return
	}

	setupGoForWindows()
}

// performUnixSetup handles Unix-like systems setup
func performUnixSetup() bool {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("\n▸ Downloading and installing 'g'...")

	// Create directory for g
	homeDir := common.GetHomeDir()
	gDir := filepath.Join(homeDir, ".g")
	if err := os.MkdirAll(gDir, 0755); err != nil {
		red.Printf("❌ Error creating .g directory: %v\n", err)
		return false
	}

	// Try to install g using the install script
	if !installGWithScript() {
		yellow.Println("  ❌ Error installing 'g'. Trying alternative method...")
		if !installGManually() {
			red.Println("  ❌ Failed to install 'g'")
			return false
		}
	}

	green.Println("  ✅ 'g' installed successfully")
	return true
}

// completeSetup finishes the setup process
func completeSetup() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	blue.Println("\n▸ Configuring PATH and environment variables...")
	configureEnvironment()

	blue.Println("\n▸ Installing latest stable Go version...")
	installLatestGo()

	blue.Println("\n▸ Activating installed Go version...")
	activateLatestGo()

	blue.Println("\n▸ Verifying installation...")
	verifyInstallation()

	createHelpScript()

	green.Println("\n✅ Installation completed!")
	displayNextSteps()
}

// installLatestGo installs the latest Go version
func installLatestGo() {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	
	homeDir := common.GetHomeDir()
	gBin := filepath.Join(homeDir, ".g", "bin", "g")
	if runtime.GOOS == "windows" {
		gBin += ".exe"
	}

	installCmd := exec.Command(gBin, "install", "latest")
	if err := installCmd.Run(); err != nil {
		yellow.Println("  ℹ️  Installing known specific version...")
		fallbackCmd := exec.Command(gBin, "install", "1.21.5")
		fallbackCmd.Run()
	} else {
		green.Println("  ✅ Go latest installed successfully")
	}
}

// activateLatestGo activates the installed Go version
func activateLatestGo() {
	homeDir := common.GetHomeDir()
	gBin := filepath.Join(homeDir, ".g", "bin", "g")
	if runtime.GOOS == "windows" {
		gBin += ".exe"
	}

	useCmd := exec.Command(gBin, "set", "latest")
	if err := useCmd.Run(); err != nil {
		// Try with specific version
		fallbackUseCmd := exec.Command(gBin, "set", "1.21.5")
		fallbackUseCmd.Run()
	}
}

// displayNextSteps shows the user what to do next
func displayNextSteps() {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	
	fmt.Println("")
	yellow.Println("📋 Next steps:")

	if runtime.GOOS == "windows" {
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
