package defaultcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// SetDefaultVersion sets a specific Go version as the default
func SetDefaultVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	blue.Printf("📌 Setting Go %s as default version...\n", version)

	// Use gobrew to set the default version
	if _, err := exec.LookPath("gobrew"); err != nil {
		red.Println("❌ gobrew is not installed")
		red.Println("💡 Run first: gos setup")
		return
	}

	cmd := exec.Command("gobrew", "use", version)

	homeDir := common.GetHomeDir()
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
