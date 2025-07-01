package defaultcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// SetDefaultVersion sets a specific Go version as the default
func SetDefaultVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	blue.Printf("📌 Setting Go %s as default version...\n", version)

	var cmd *exec.Cmd
	homeDir := common.GetHomeDir()
	
	if runtime.GOOS == "windows" && common.IsCommandAvailable("gobrew") {
		cmd = exec.Command("gobrew", "use", version)
	} else if common.IsCommandAvailable("gobrew") {
		cmd = exec.Command("gobrew", "use", version)
	} else {
		// Try different possible locations for g
		gPaths := []string{
			filepath.Join(homeDir, ".g", "bin", "g"),
			filepath.Join(homeDir, "go", "bin", "g"),
			"/usr/local/bin/g",
		}
		
		var gPath string
		for _, path := range gPaths {
			if _, err := os.Stat(path); err == nil {
				gPath = path
				break
			}
		}
		
		if gPath == "" {
			red.Println("❌ No version manager available")
			return
		}
		
		cmd = exec.Command(gPath, "set", version)
	}

	if err := cmd.Run(); err != nil {
		red.Printf("❌ Error setting default version: %v\n", err)
		return
	}

	// Save the default version to a file for persistence
	defaultFile := filepath.Join(homeDir, ".gos-default")
	if err := os.WriteFile(defaultFile, []byte(version), 0644); err != nil {
		red.Printf("⚠️  Warning: Could not save default version: %v\n", err)
	}

	green.Printf("✅ Go %s is now the default version\n", version)
	
	// Verify the change
	fmt.Println()
	blue.Println("🔍 Verifying...")
	if goCmd := exec.Command("go", "version"); goCmd.Run() == nil {
		output, _ := goCmd.Output()
		fmt.Printf("  %s", output)
	}
}
