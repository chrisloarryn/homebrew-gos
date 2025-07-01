package use

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// UseVersion switches to a specific Go version
func UseVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	blue.Printf("🔄 Switching to Go %s...\n", version)

	var cmd *exec.Cmd
	if common.IsCommandAvailable("gobrew") {
		cmd = exec.Command("gobrew", "use", version)
	} else if common.IsCommandAvailable("g") {
		cmd = exec.Command("g", "set", version)
	} else {
		red.Println("❌ No version manager available")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		red.Printf("❌ Error switching to Go %s\n", version)
		yellow.Printf("💡 Is this version installed? Use: gos list\n")
		return
	}

	green.Printf("✅ Switched to Go %s\n", version)

	// Update PATH automatically after successful version switch
	common.UpdatePathForVersionManager()

	// Show current version and PATH update instructions
	blue.Println("\n📋 Verifying installation...")
	if goCmd := exec.Command("go", "version"); goCmd.Run() == nil {
		blue.Print("✅ Current version: ")
		goCmd.Stdout = os.Stdout
		goCmd.Run()
		fmt.Println()
	} else {
		showPathUpdateInstructions(yellow)
	}
}

// showPathUpdateInstructions displays instructions for updating PATH
func showPathUpdateInstructions(yellow *color.Color) {
	yellow.Println("⚠️  PATH needs to be updated for this terminal session.")
	yellow.Println("💡 To use the new Go version immediately, run:")
	if common.IsCommandAvailable("gobrew") {
		if runtime.GOOS == "windows" {
			yellow.Println("   $env:PATH = \"$env:USERPROFILE\\.gobrew\\current\\bin;$env:USERPROFILE\\.gobrew\\bin;$env:PATH\"")
		} else {
			yellow.Println("   export PATH=\"$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH\"")
		}
	} else if common.IsCommandAvailable("g") {
		if runtime.GOOS == "windows" {
			yellow.Println("   $env:PATH = \"$env:USERPROFILE\\.g\\go\\bin;$env:USERPROFILE\\.g\\bin;$env:PATH\"")
		} else {
			yellow.Println("   export PATH=\"$HOME/.g/go/bin:$HOME/.g/bin:$PATH\"")
		}
	}
	fmt.Println()
	color.New(color.FgBlue).Println("🔄 Or simply open a new terminal window.")
}
