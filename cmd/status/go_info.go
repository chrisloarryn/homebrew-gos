package status

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// ShowCurrentGo displays information about the current Go installation
func ShowCurrentGo() {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	if _, err := exec.LookPath("go"); err != nil {
		yellow.Println("  ⚠️  Go is not available in PATH")
		yellow.Println("  💡 Try running: source ~/.zshrc")
		return
	}

	// Show version with better formatting
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("  ✅ %s\n", version)
	}

	// Show GOROOT with validation
	if output, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		goroot := strings.TrimSpace(string(output))
		expectedGoroot := filepath.Join(os.Getenv("HOME"), ".g", "go")
		if goroot == expectedGoroot {
			green.Printf("  ✅ GOROOT: %s\n", goroot)
		} else {
			blue.Printf("  ℹ️  GOROOT: %s\n", goroot)
		}
	}

	// Show GOPATH
	if output, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		fmt.Printf("  GOPATH: %s", string(output))
	}
}

// ShowDiskUsage displays disk usage information for Go directories
func ShowDiskUsage() {
	homeDir := os.Getenv("HOME")
	gDir := filepath.Join(homeDir, ".g")

	if _, err := os.Stat(gDir); err == nil {
		if output, err := exec.Command("du", "-sh", gDir).Output(); err == nil {
			fmt.Printf("  ~/.g directory: %s", string(output))
		} else {
			fmt.Println("  Could not calculate ~/.g directory size")
		}
	} else {
		fmt.Println("  ~/.g directory not found")
	}

	// Show Go workspace size if it exists
	goDir := filepath.Join(homeDir, "go")
	if _, err := os.Stat(goDir); err == nil {
		if output, err := exec.Command("du", "-sh", goDir).Output(); err == nil {
			fmt.Printf("  ~/go directory: %s", string(output))
		}
	}
}
