package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// installGWithScript attempts to install g using the official install script
func installGWithScript() bool {
	blue := color.New(color.FgBlue)
	blue.Println("  📥 Downloading g installer...")
	
	// Create progress bar for installation
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Installing g"),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
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
	if runtime.GOOS == "windows" {
		// For Windows, try with PowerShell and curl if available
		if common.IsCommandAvailable("curl") {
			result = common.ExecuteWithShell("curl -sSL https://git.io/g-install | bash -s -- -y") == nil
		} else {
			result = false
		}
	} else {
		// Unix-like systems
		result = common.ExecuteWithShell("curl -sSL https://git.io/g-install | bash -s -- -y") == nil
	}
	
	bar.Finish()
	return result
}

// installGManually attempts manual installation of g
func installGManually() bool {
	homeDir := common.GetHomeDir()

	if runtime.GOOS == "windows" {
		return installGForWindows(homeDir)
	}

	// Unix-like systems
	return installGForUnix(homeDir)
}

// installGForUnix installs g for Unix-like systems
func installGForUnix(homeDir string) bool {
	red := color.New(color.FgRed)
	blue := color.New(color.FgBlue)
	
	blue.Println("  🔄 Manual installation for Unix-like system...")
	
	// Create .g directory structure
	gDir := filepath.Join(homeDir, ".g")
	binDir := filepath.Join(gDir, "bin")
	
	if err := os.MkdirAll(binDir, 0755); err != nil {
		red.Printf("❌ Error creating directories: %v\n", err)
		return false
	}
	
	// Try to download g binary directly
	if common.IsCommandAvailable("curl") {
		downloadURL := "https://github.com/stefanmaric/g/releases/latest/download/g-$(uname -s)-$(uname -m)"
		gBinary := filepath.Join(binDir, "g")
		
		downloadCmd := fmt.Sprintf("curl -L %s -o %s && chmod +x %s", downloadURL, gBinary, gBinary)
		if common.ExecuteWithShell(downloadCmd) == nil {
			return true
		}
	}
	
	// Try wget as fallback
	if common.IsCommandAvailable("wget") {
		downloadURL := "https://github.com/stefanmaric/g/releases/latest/download/g-$(uname -s)-$(uname -m)"
		gBinary := filepath.Join(binDir, "g")
		
		downloadCmd := fmt.Sprintf("wget %s -O %s && chmod +x %s", downloadURL, gBinary, gBinary)
		if common.ExecuteWithShell(downloadCmd) == nil {
			return true
		}
	}
	
	return false
}

// checkExistingInstallations checks if version managers are already installed
func checkExistingInstallations() bool {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	
	hasG := common.IsCommandAvailable("g")
	hasGobrew := common.IsCommandAvailable("gobrew")
	
	if hasG || hasGobrew {
		green.Println("✅ Version manager already detected:")
		if hasG {
			fmt.Println("  • 'g' is available")
		}
		if hasGobrew {
			fmt.Println("  • 'gobrew' is available")
		}
		fmt.Println("")
		yellow.Println("💡 Use --force to reinstall anyway")
		fmt.Println("   Example: gos setup --force")
		return true
	}
	
	return false
}
