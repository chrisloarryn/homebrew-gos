package defaultcmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
)

// ShowDefaultVersion displays the current default Go version
func ShowDefaultVersion() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	blue.Println("📌 Default Go version:")
	
	homeDir := common.GetHomeDir()
	
	// First check if we have a saved default version
	defaultFile := filepath.Join(homeDir, ".gos-default")
	if content, err := os.ReadFile(defaultFile); err == nil {
		version := strings.TrimSpace(string(content))
		if version != "" {
			green.Printf("  ✅ %s (saved default)\n", version)
			return
		}
	}
	
	if common.IsCommandAvailable("gobrew") {
		// gobrew: check current version
		cmd := exec.Command("gobrew", "use")
		if output, err := cmd.Output(); err == nil {
			version := strings.TrimSpace(string(output))
			if version != "" {
				green.Printf("  ✅ %s (via gobrew)\n", version)
				return
			}
		}
	} else {
		// Try different possible locations for g
		gPaths := []string{
			filepath.Join(homeDir, ".g", "bin", "g"),
			filepath.Join(homeDir, "go", "bin", "g"),
			"/usr/local/bin/g",
		}
		
		for _, gPath := range gPaths {
			if _, err := os.Stat(gPath); err == nil {
				// Check current active version in g using 'list' command
				cmd := exec.Command(gPath, "list")
				if output, err := cmd.Output(); err == nil {
					lines := strings.Split(string(output), "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if strings.HasPrefix(line, ">") {
							version := strings.TrimSpace(strings.TrimPrefix(line, ">"))
							green.Printf("  ✅ %s (via g)\n", version)
							return
						}
					}
				}
				break
			}
		}
		
		// Alternative: check symlink in .g directory
		goLink := filepath.Join(homeDir, ".g", "go")
		if target, err := os.Readlink(goLink); err == nil {
			version := filepath.Base(target)
			green.Printf("  ✅ %s (via g symlink)\n", version)
			return
		}
	}
	
	yellow.Println("  ⚠️  No default version set")
}
