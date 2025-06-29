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

// Constants for repeated strings
const (
	bashrcFile        = ".bashrc"
	powerShellCommand = "-Command"
)

// getHomeDir returns the home directory for cross-platform compatibility
func getHomeDir() string {
	if runtime.GOOS == "windows" {
		if home := os.Getenv("USERPROFILE"); home != "" {
			return home
		}
		if home := os.Getenv("HOME"); home != "" {
			return home
		}
		if homedrive := os.Getenv("HOMEDRIVE"); homedrive != "" {
			if homepath := os.Getenv("HOMEPATH"); homepath != "" {
				return homedrive + homepath
			}
		}
	}
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return "."
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// executeWithShell runs a command with appropriate shell for the OS
func executeWithShell(command string) error {
	if runtime.GOOS == "windows" {
		// Try PowerShell first, then cmd
		if isCommandAvailable("powershell") {
			cmd := exec.Command("powershell", powerShellCommand, command)
			return cmd.Run()
		}
		cmd := exec.Command("cmd", "/C", command)
		return cmd.Run()
	}
	// Unix-like systems
	cmd := exec.Command("bash", "-c", command)
	return cmd.Run()
}

var setupCmd = &cobra.Command{
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

func init() {
	setupCmd.Flags().BoolP("force", "f", false, "Force reinstallation even if version managers are already installed")
}

func setupGoVersionManager(force bool) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
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

	// Check if Windows and use alternative setup
	if osName == "Windows" {
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
		return
	}

	blue.Println("\n▸ Downloading and installing 'g'...")

	// Create directory for g
	homeDir := getHomeDir()
	gDir := filepath.Join(homeDir, ".g")
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
	gBin := filepath.Join(homeDir, ".g", "bin", "g")
	if runtime.GOOS == "windows" {
		gBin += ".exe" // Add .exe extension for Windows
	}

	installCmd := exec.Command(gBin, "install", "latest")
	if err := installCmd.Run(); err != nil {
		yellow.Println("  ℹ️  Installing known specific version...")
		fallbackCmd := exec.Command(gBin, "install", "1.21.5")
		fallbackCmd.Run()
	} else {
		green.Println("  ✅ Go latest installed successfully")
	}

	// Auto-use the installed version
	blue.Println("\n▸ Activating installed Go version...")
	useCmd := exec.Command(gBin, "set", "latest")
	if err := useCmd.Run(); err != nil {
		// Try with specific version
		fallbackUseCmd := exec.Command(gBin, "set", "1.21.5")
		fallbackUseCmd.Run()
	}

	// Verify installation
	blue.Println("\n▸ Verifying installation...")
	verifyInstallation()

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
	if runtime.GOOS == "windows" {
		// For Windows, try with PowerShell and curl if available
		if isCommandAvailable("curl") {
			cmd := exec.Command("powershell", "-Command", "curl -sSL https://git.io/g-install | bash -s -- -y")
			return cmd.Run() == nil
		}
		// Fallback to manual installation for Windows
		return false
	}
	// Unix-like systems
	cmd := exec.Command("bash", "-c", "curl -sSL https://git.io/g-install | bash -s -- -y")
	return cmd.Run() == nil
}

func installGManually() bool {
	homeDir := getHomeDir()

	if runtime.GOOS == "windows" {
		return installGForWindows(homeDir)
	}

	// Unix-like systems
	return installGForUnix(homeDir)
}

func installGForWindows(homeDir string) bool {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	blue.Println("  💡 Windows detected - using alternative Go version managers...")
	fmt.Println("")

	// Option 1: Try to install gobrew (best option for Windows)
	blue.Println("  🔄 Attempting to install 'gobrew' (recommended for Windows)...")
	if installGobrew() {
		green.Println("  ✅ gobrew installed successfully!")
		fmt.Println("  📋 You can now use: gobrew use latest")
		return true
	}

	// Option 2: Try to install voidint/g (supports Windows)
	blue.Println("  🔄 Attempting to install 'voidint/g' (Windows compatible)...")
	if installVoidintG() {
		green.Println("  ✅ voidint/g installed successfully!")
		fmt.Println("  📋 You can now use: g install latest")
		return true
	}

	// Option 3: Manual Go installation
	blue.Println("  🔄 Attempting direct Go installation...")
	if installGoDirectly(homeDir) {
		green.Println("  ✅ Go installed directly!")
		return true
	}

	// If all fail, show manual options
	red.Println("  ❌ Automatic installation failed.")
	yellow.Println("  💡 Manual installation options:")
	fmt.Println("")
	fmt.Println("     🍺 Option 1 - Chocolatey:")
	fmt.Println("       choco install golang")
	fmt.Println("")
	fmt.Println("     📦 Option 2 - Scoop:")
	fmt.Println("       scoop install go")
	fmt.Println("")
	fmt.Println("     🌐 Option 3 - Official installer:")
	fmt.Println("       Download from: https://golang.org/dl/")
	fmt.Println("")
	fmt.Println("     🐧 Option 4 - WSL (Windows Subsystem for Linux):")
	fmt.Println("       Install WSL and use the Linux version of gos")
	fmt.Println("")

	return false
}

// Install gobrew - best option for Windows
func installGobrew() bool {
	// Try PowerShell installation
	if isCommandAvailable("powershell") {
		installScript := "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.ps1'))"

		cmd := exec.Command("powershell", powerShellCommand, installScript)
		if cmd.Run() == nil {
			return true
		}
	}

	// Try with curl if available
	if isCommandAvailable("curl") && isCommandAvailable("bash") {
		cmd := exec.Command("bash", "-c", "curl -sL https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | bash")
		if cmd.Run() == nil {
			return true
		}
	}

	return false
}

// Install voidint/g - alternative that supports Windows
func installVoidintG() bool {
	if isCommandAvailable("powershell") {
		installScript := "iwr https://raw.githubusercontent.com/voidint/g/master/install.ps1 -useb | iex"

		cmd := exec.Command("powershell", powerShellCommand, installScript)
		if cmd.Run() == nil {
			return true
		}
	}

	return false
}

// Direct Go installation as fallback
func installGoDirectly(homeDir string) bool {
	version := "1.21.6"
	goURL := fmt.Sprintf("https://golang.org/dl/go%s.windows-amd64.zip", version)

	blue := color.New(color.FgBlue)
	blue.Printf("  📥 Downloading Go %s for Windows...\n", version)

	return downloadAndInstallGo(goURL, version, homeDir)
}

func installGForUnix(homeDir string) bool {
	// Try git clone method
	if _, err := exec.LookPath("git"); err == nil {
		tmpDir := "/tmp/g"
		cloneCmd := exec.Command("git", "clone", "https://github.com/stefanmaric/g.git", tmpDir)
		if cloneCmd.Run() == nil {
			makeCmd := exec.Command("make", "install", fmt.Sprintf("PREFIX=%s", filepath.Join(homeDir, ".g")))
			makeCmd.Dir = tmpDir
			if makeCmd.Run() == nil {
				os.RemoveAll(tmpDir)
				return true
			}
		}
		os.RemoveAll(tmpDir)
	}

	// Fallback: download directly
	gBinDir := filepath.Join(homeDir, ".g", "bin")
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
	homeDir := getHomeDir()
	osName := runtime.GOOS
	var shellFiles []string

	if osName == "windows" {
		shellFiles = []string{
			filepath.Join(homeDir, bashrcFile),
			filepath.Join(homeDir, ".bash_profile"),
		}
	} else {
		shellFiles = []string{
			filepath.Join(homeDir, ".zshrc"),
			filepath.Join(homeDir, bashrcFile),
			filepath.Join(homeDir, ".bash_profile"),
		}
	}

	gConfig := generateGConfig(homeDir, osName)

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

	setEnvironmentVariables(homeDir)
}

func generateGConfig(homeDir, osName string) string {
	if osName == "windows" {
		return fmt.Sprintf(`
# === Go Version Manager (g) ===
export GOPATH=%s/go
export GOROOT=%s/.g/go
export PATH=%s/.g/bin:$GOROOT/bin:$GOPATH/bin:$PATH

# Ensure GOPATH bin directory exists
mkdir -p $GOPATH/bin

# Aliases for convenience
alias go-reload='source ~/.bashrc && go version'
alias go-help='cat ~/.g/go-help.sh'
`, homeDir, homeDir, homeDir)
	}

	return fmt.Sprintf(`
# === Go Version Manager (g) ===
export GOPATH=%s/go
export GOROOT=%s/.g/go
export PATH=%s/.g/bin:$GOROOT/bin:$GOPATH/bin:$PATH

# Ensure GOPATH bin directory exists
mkdir -p $GOPATH/bin

# Aliases for convenience
alias go-reload='source ~/.zshrc && go version'
alias go-help='cat ~/.g/go-help.sh'
`, homeDir, homeDir, homeDir)
}

func setEnvironmentVariables(homeDir string) {
	// Export for current session
	os.Setenv("GOPATH", filepath.Join(homeDir, "go"))
	os.Setenv("GOROOT", filepath.Join(homeDir, ".g", "go"))
	currentPath := os.Getenv("PATH")

	pathSep := ":"
	if runtime.GOOS == "windows" {
		pathSep = ";"
	}

	newPath := fmt.Sprintf("%s%s%s%s%s%s%s",
		filepath.Join(homeDir, ".g", "bin"),
		pathSep,
		filepath.Join(homeDir, ".g", "go", "bin"),
		pathSep,
		filepath.Join(homeDir, "go", "bin"),
		pathSep,
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
	homeDir := getHomeDir()
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

	helpFile := filepath.Join(homeDir, ".g", "go-help.sh")
	if err := os.WriteFile(helpFile, []byte(helpScript), 0755); err == nil {
		color.Blue("\n▸ Creating help script...")
	}
}

func verifyInstallation() {
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	homeDir := getHomeDir()

	// Check if 'g' is available
	gBin := filepath.Join(homeDir, ".g", "bin", "g")
	if runtime.GOOS == "windows" {
		gBin += ".exe"
	}

	if _, err := os.Stat(gBin); err != nil {
		red.Println("  ❌ 'g' binary not found")
		return
	}
	green.Println("  ✅ 'g' version manager installed")

	// Check if Go is available in PATH
	if _, err := exec.LookPath("go"); err != nil {
		yellow.Println("  ⚠️  Go not found in PATH (may need to restart shell)")
		if runtime.GOOS == "windows" {
			yellow.Println("      Please restart Git Bash/WSL or run: source ~/.bashrc")
		} else {
			yellow.Println("      Please run: source ~/.zshrc")
		}
		return
	}

	// Check Go version
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("  ✅ %s\n", version)
	}

	// Check GOROOT and GOPATH
	if goroot, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		expectedGoroot := filepath.Join(homeDir, ".g", "go")
		actualGoroot := strings.TrimSpace(string(goroot))
		if actualGoroot == expectedGoroot {
			green.Printf("  ✅ GOROOT correctly set: %s\n", actualGoroot)
		} else {
			yellow.Printf("  ⚠️  GOROOT: %s (expected: %s)\n", actualGoroot, expectedGoroot)
		}
	}

	if gopath, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		expectedGopath := filepath.Join(homeDir, "go")
		actualGopath := strings.TrimSpace(string(gopath))
		if actualGopath == expectedGopath {
			green.Printf("  ✅ GOPATH correctly set: %s\n", actualGopath)
		} else {
			yellow.Printf("  ⚠️  GOPATH: %s (expected: %s)\n", actualGopath, expectedGopath)
		}
	}
}

// setupGoForWindows provides an alternative setup for Windows using gobrew or voidint/g
func setupGoForWindows() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("🔧 Setting up Go for Windows...")
	fmt.Println("  Detected: Windows")

	homeDir := getHomeDir()

	blue.Println("\n▸ Installing Windows-compatible Go version manager...")

	// Try to install a Windows-compatible version manager
	if installGForWindows(homeDir) {
		green.Println("\n✅ Go version manager installed successfully!")

		blue.Println("\n▸ Installing latest Go version...")

		// Try to install Go with the installed manager
		if installGoWithManager() {
			green.Println("  ✅ Go installed successfully!")
		} else {
			yellow.Println("  ⚠️  Manual Go installation may be needed")
		}

		blue.Println("\n▸ Configuring environment...")
		configureWindowsEnvironmentForManager()

		green.Println("\n✅ Installation completed!")
		showWindowsPostInstallInstructions()
	} else {
		red.Println("\n❌ Automatic installation failed")
		showWindowsManualInstructions()
	}
}

func installGoWithManager() bool {
	// Try gobrew first
	if isCommandAvailable("gobrew") {
		cmd := exec.Command("gobrew", "use", "latest")
		return cmd.Run() == nil
	}

	// Try voidint/g
	if isCommandAvailable("g") {
		cmd := exec.Command("g", "install", "latest")
		return cmd.Run() == nil
	}

	return false
}

func configureWindowsEnvironmentForManager() {
	homeDir := getHomeDir()

	// Configure for gobrew
	if isCommandAvailable("gobrew") {
		configurePowerShellProfile(homeDir, "gobrew")
		configureBashProfile(homeDir, "gobrew")
		return
	}

	// Configure for voidint/g
	if isCommandAvailable("g") {
		configurePowerShellProfile(homeDir, "voidint-g")
		configureBashProfile(homeDir, "voidint-g")
		return
	}
}

func configurePowerShellProfile(homeDir, manager string) {
	profilePath := filepath.Join(homeDir, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
	profileDir := filepath.Dir(profilePath)
	os.MkdirAll(profileDir, 0755)

	var config string
	switch manager {
	case "gobrew":
		config = `
# === Go Configuration (gobrew) ===
$env:PATH = "$HOME\.gobrew\current\bin;$HOME\.gobrew\bin;" + $env:PATH
$env:GOROOT = "$HOME\.gobrew\current\go"
$env:GOPATH = "$HOME\go"

# Create GOPATH bin directory if it doesn't exist
if (!(Test-Path "$env:GOPATH\bin")) {
    New-Item -ItemType Directory -Path "$env:GOPATH\bin" -Force | Out-Null
}
`
	case "voidint-g":
		config = `
# === Go Configuration (voidint/g) ===
$env:GOROOT = "$HOME\.g\go"
$env:PATH = "$HOME\.g\bin;$env:GOROOT\bin;" + $env:PATH
$env:GOPATH = "$HOME\go"

# Create GOPATH bin directory if it doesn't exist
if (!(Test-Path "$env:GOPATH\bin")) {
    New-Item -ItemType Directory -Path "$env:GOPATH\bin" -Force | Out-Null
}
`
	}

	if config != "" && !hasGConfigPowerShell(profilePath) {
		appendToFile(profilePath, config)
		color.Green("  ✅ PowerShell profile configured for %s", manager)
	}
}

func configureBashProfile(homeDir, manager string) {
	bashrcPath := filepath.Join(homeDir, bashrcFile)

	var config string
	switch manager {
	case "gobrew":
		config = `
# === Go Configuration (gobrew) ===
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
export GOROOT="$HOME/.gobrew/current/go"
export GOPATH="$HOME/go"

# Ensure GOPATH bin directory exists
mkdir -p "$GOPATH/bin"
`
	case "voidint-g":
		config = `
# === Go Configuration (voidint/g) ===
export GOROOT="$HOME/.g/go"
export PATH="$HOME/.g/bin:$GOROOT/bin:$PATH"
export GOPATH="$HOME/go"

# Ensure GOPATH bin directory exists
mkdir -p "$GOPATH/bin"
`
	}

	if config != "" && !hasGConfig(bashrcPath) {
		appendToFile(bashrcPath, config)
		color.Green("  ✅ Bash profile configured for %s", manager)
	}
}

func showWindowsPostInstallInstructions() {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)

	fmt.Println("")
	yellow.Println("� Next steps:")

	if isCommandAvailable("gobrew") {
		fmt.Println("1. Restart your PowerShell/Terminal")
		fmt.Println("2. Verify: gobrew version")
		fmt.Println("3. Check Go: go version")
		fmt.Println("")
		blue.Println("🚀 Quick examples with gobrew:")
		fmt.Println("   gobrew use 1.21.6        # Switch to Go 1.21.6")
		fmt.Println("   gobrew use latest        # Use latest version")
		fmt.Println("   gobrew ls               # List installed versions")
		fmt.Println("   gobrew ls-remote        # List available versions")
	} else if isCommandAvailable("g") {
		fmt.Println("1. Restart your PowerShell/Terminal")
		fmt.Println("2. Verify: g version")
		fmt.Println("3. Check Go: go version")
		fmt.Println("")
		blue.Println("🚀 Quick examples with voidint/g:")
		fmt.Println("   g install 1.21.6        # Install Go 1.21.6")
		fmt.Println("   g use 1.21.6           # Switch to Go 1.21.6")
		fmt.Println("   g ls                   # List installed versions")
		fmt.Println("   g ls-remote            # List available versions")
	}

	fmt.Println("")
	yellow.Println("💡 You can also continue using gos commands:")
	fmt.Println("   gos status             # Show current status")
	fmt.Println("   gos list               # List versions (if compatible)")
}

func showWindowsManualInstructions() {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)

	fmt.Println("")
	yellow.Println("📋 Manual installation options:")
	fmt.Println("")

	blue.Println("🍺 Option 1 - Chocolatey (recommended):")
	fmt.Println("   # Install Chocolatey first (if not installed):")
	fmt.Println("   Set-ExecutionPolicy Bypass -Scope Process -Force")
	fmt.Println("   [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072")
	fmt.Println("   iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))")
	fmt.Println("")
	fmt.Println("   # Then install Go:")
	fmt.Println("   choco install golang")
	fmt.Println("")

	blue.Println("📦 Option 2 - Scoop:")
	fmt.Println("   # Install Scoop first (if not installed):")
	fmt.Println("   Set-ExecutionPolicy RemoteSigned -Scope CurrentUser")
	fmt.Println("   irm get.scoop.sh | iex")
	fmt.Println("")
	fmt.Println("   # Then install Go:")
	fmt.Println("   scoop install go")
	fmt.Println("")

	blue.Println("🌐 Option 3 - Official installer:")
	fmt.Println("   1. Visit: https://golang.org/dl/")
	fmt.Println("   2. Download the Windows installer (.msi)")
	fmt.Println("   3. Run the installer")
	fmt.Println("")

	blue.Println("� Option 4 - WSL (Windows Subsystem for Linux):")
	fmt.Println("   1. Install WSL: wsl --install")
	fmt.Println("   2. Use the Linux version of gos in WSL")
	fmt.Println("")

	yellow.Println("💡 After manual installation, you can use gos commands:")
	fmt.Println("   gos status             # Check current setup")
	fmt.Println("   gos env                # Show environment info")
}

// Helper functions for Windows Go installation

func downloadAndInstallGo(url, version, goDir string) bool {
	// Create temporary directory
	tempDir := filepath.Join(os.TempDir(), "gos-install")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return false
	}
	defer os.RemoveAll(tempDir)

	zipFile := filepath.Join(tempDir, fmt.Sprintf("go%s.zip", version))

	// Download Go zip file
	if !downloadFile(url, zipFile) {
		return false
	}

	// Extract zip file
	extractDir := filepath.Join(goDir, fmt.Sprintf("go%s", version))
	return extractZip(zipFile, extractDir)
}

func downloadFile(url, filepath string) bool {
	// Try with PowerShell Invoke-WebRequest first
	if isCommandAvailable("powershell") {
		cmd := exec.Command("powershell", powerShellCommand,
			fmt.Sprintf("Invoke-WebRequest -Uri '%s' -OutFile '%s'", url, filepath))
		if cmd.Run() == nil {
			return true
		}
	}

	// Try with curl if available
	if isCommandAvailable("curl") {
		cmd := exec.Command("curl", "-L", "-o", filepath, url)
		if cmd.Run() == nil {
			return true
		}
	}

	return false
}

func extractZip(src, dest string) bool {
	if isCommandAvailable("powershell") {
		cmd := exec.Command("powershell", powerShellCommand,
			fmt.Sprintf("Expand-Archive -Path '%s' -DestinationPath '%s' -Force", src, dest))
		return cmd.Run() == nil
	}
	return false
}

func configureWindowsEnvironment(goDir, version string) {
	homeDir := getHomeDir()

	// Create PowerShell profile if it doesn't exist
	profilePath := filepath.Join(homeDir, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
	profileDir := filepath.Dir(profilePath)
	os.MkdirAll(profileDir, 0755)

	goRoot := filepath.Join(goDir, fmt.Sprintf("go%s", version), "go")
	goPath := filepath.Join(homeDir, "go")

	psConfig := fmt.Sprintf(`
# === Go Configuration (gos) ===
$env:GOROOT = "%s"
$env:GOPATH = "%s"
$env:PATH = "$env:GOROOT\bin;$env:GOPATH\bin;" + $env:PATH

# Create GOPATH bin directory if it doesn't exist
if (!(Test-Path "$env:GOPATH\bin")) {
    New-Item -ItemType Directory -Path "$env:GOPATH\bin" -Force | Out-Null
}
`, goRoot, goPath)

	if !hasGConfigPowerShell(profilePath) {
		appendToFile(profilePath, psConfig)
		color.Green("  ✅ PowerShell profile configured")
	}

	// Also create .bashrc for Git Bash users
	bashrcPath := filepath.Join(homeDir, bashrcFile)
	bashConfig := fmt.Sprintf(`
# === Go Configuration (gos) ===
export GOROOT="%s"
export GOPATH="%s"
export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"

# Ensure GOPATH bin directory exists
mkdir -p "$GOPATH/bin"
`, goRoot, goPath)

	if !hasGConfig(bashrcPath) {
		appendToFile(bashrcPath, bashConfig)
		color.Green("  ✅ Bash profile configured")
	}
}

func hasGConfigPowerShell(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Go Configuration (gos)") {
			return true
		}
	}
	return false
}

func createVersionFile(goDir, version string) {
	versionFile := filepath.Join(goDir, "current_version")
	os.WriteFile(versionFile, []byte(version), 0644)
}

// checkExistingInstallations verifies if any Go version manager is already installed
func checkExistingInstallations() bool {
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	var installedManagers []string

	// Check for gobrew
	if isCommandAvailable("gobrew") {
		if output, err := exec.Command("gobrew", "version").Output(); err == nil {
			version := strings.TrimSpace(string(output))
			installedManagers = append(installedManagers, fmt.Sprintf("gobrew (%s)", version))
		}
	}

	// Check for voidint/g
	if isCommandAvailable("g") {
		if output, err := exec.Command("g", "version").Output(); err == nil {
			version := strings.TrimSpace(string(output))
			installedManagers = append(installedManagers, fmt.Sprintf("voidint/g (%s)", version))
		}
	}

	// Check for original g (Unix systems)
	homeDir := getHomeDir()
	if runtime.GOOS != "windows" {
		gBin := filepath.Join(homeDir, ".g", "bin", "g")
		if _, err := os.Stat(gBin); err == nil {
			installedManagers = append(installedManagers, "g (original)")
		}
	}

	// Check for gvm
	if isCommandAvailable("gvm") {
		installedManagers = append(installedManagers, "gvm")
	}

	// Check if Go is already installed
	goInstalled := false
	var goVersion string
	if output, err := exec.Command("go", "version").Output(); err == nil {
		goVersion = strings.TrimSpace(string(output))
		goInstalled = true
	}

	if len(installedManagers) > 0 || goInstalled {
		green.Println("✅ Existing Go setup detected!")
		fmt.Println("")

		if goInstalled {
			green.Printf("  🐹 Go is installed: %s\n", goVersion)
		}

		if len(installedManagers) > 0 {
			blue.Println("  📦 Version managers found:")
			for _, manager := range installedManagers {
				fmt.Printf("    • %s\n", manager)
			}
		}

		fmt.Println("")
		yellow.Println("💡 Your Go development environment is already configured!")
		fmt.Println("   No additional setup needed.")
		fmt.Println("")
		fmt.Println("   You can use:")
		if goInstalled {
			fmt.Println("   • go version              # Check current Go version")
		}
		for _, manager := range installedManagers {
			if strings.Contains(manager, "gobrew") {
				fmt.Println("   • gobrew ls               # List installed versions")
				fmt.Println("   • gobrew use latest       # Switch to latest version")
			} else if strings.Contains(manager, "voidint/g") {
				fmt.Println("   • g ls                    # List installed versions")
				fmt.Println("   • g install latest        # Install latest version")
			} else if strings.Contains(manager, "gvm") {
				fmt.Println("   • gvm list                # List installed versions")
				fmt.Println("   • gvm use latest          # Switch to latest version")
			}
		}
		fmt.Println("   • gos status              # Show detailed status")
		fmt.Println("")

		yellow.Printf("   To force reinstallation, use: gos setup --force\n")
		return true
	}

	return false
}
