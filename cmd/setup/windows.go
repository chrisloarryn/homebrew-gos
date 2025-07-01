package setup

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// Constants for Windows installation
const (
	powerShellCommand = "-Command"
)

// installGForWindows handles Windows-specific installation
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

// installGobrew installs gobrew - best option for Windows
func installGobrew() bool {
	blue := color.New(color.FgBlue)
	blue.Println("  📥 Installing gobrew...")
	
	// Create progress bar
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Installing gobrew"),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "=", SaucerHead: ">", SaucerPadding: " ", BarStart: "[", BarEnd: "]",
		}),
	)
	
	go func() {
		for {
			bar.Add(1)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	var result bool
	// Try PowerShell installation
	if common.IsCommandAvailable("powershell") {
		installScript := "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.ps1'))"

		cmd := exec.Command("powershell", powerShellCommand, installScript)
		result = cmd.Run() == nil
		if result {
			bar.Finish()
			return true
		}
	}

	// Try with curl if available
	if common.IsCommandAvailable("curl") && common.IsCommandAvailable("bash") {
		result = common.ExecuteWithShell("curl -sL https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | bash") == nil
	} else {
		result = false
	}

	bar.Finish()
	return result
}

// installVoidintG installs voidint/g - alternative that supports Windows
func installVoidintG() bool {
	if common.IsCommandAvailable("powershell") {
		installScript := "iwr https://raw.githubusercontent.com/voidint/g/master/install.ps1 -useb | iex"
		return common.ExecuteWithShell(installScript) == nil
	}
	return false
}

// installGoDirectly installs Go directly as fallback
func installGoDirectly(homeDir string) bool {
	version := "1.21.6"
	goURL := fmt.Sprintf("https://golang.org/dl/go%s.windows-amd64.zip", version)

	blue := color.New(color.FgBlue)
	blue.Printf("  📥 Downloading Go %s for Windows...\n", version)

	// This would need implementation for downloading and extracting
	// For now, return false as it's complex to implement properly
	_ = goURL
	return false
}

// setupGoForWindows provides Windows-specific setup instructions
func setupGoForWindows() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	blue.Println("🪟 Windows Go Setup")
	fmt.Println("")

	yellow.Println("💡 Recommended options for Windows:")
	fmt.Println("")

	fmt.Println("1️⃣  Gobrew (Recommended):")
	fmt.Println("   • Best option for Windows")
	fmt.Println("   • Similar to 'g' but Windows-compatible")
	fmt.Println("   • Commands: gobrew install latest, gobrew use latest")
	fmt.Println("")

	fmt.Println("2️⃣  Package Managers:")
	fmt.Println("   🍺 Chocolatey: choco install golang")
	fmt.Println("   📦 Scoop: scoop install go")
	fmt.Println("   🍃 Winget: winget install GoLang.Go")
	fmt.Println("")

	fmt.Println("3️⃣  Official Installer:")
	fmt.Println("   🌐 Download from: https://golang.org/dl/")
	fmt.Println("")

	fmt.Println("4️⃣  WSL (Windows Subsystem for Linux):")
	fmt.Println("   🐧 Install WSL and use the Linux version of gos")
	fmt.Println("")

	// Try to install gobrew automatically
	blue.Println("🚀 Attempting automatic gobrew installation...")
	if installGobrew() {
		green.Println("✅ gobrew installed successfully!")
		fmt.Println("")
		yellow.Println("📋 Next steps:")
		fmt.Println("   gobrew install latest")
		fmt.Println("   gobrew use latest")
	} else {
		yellow.Println("❌ Automatic installation failed. Please use manual options above.")
	}
}
