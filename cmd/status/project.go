package status

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

// ShowProjectConfig displays project-specific Go configuration
func ShowProjectConfig() {
	// Check for .go-version file
	if _, err := os.Stat(".go-version"); err == nil {
		if content, err := os.ReadFile(".go-version"); err == nil {
			version := strings.TrimSpace(string(content))
			color.Green("  ✅ .go-version found: %s", version)
		}
	} else {
		color.Yellow("  ℹ️  No .go-version file in current directory")
	}

	// Check for go.mod
	if _, err := os.Stat("go.mod"); err == nil {
		if content, err := os.ReadFile("go.mod"); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "go ") {
					version := strings.TrimPrefix(line, "go ")
					color.Green("  ✅ go.mod found, Go version: %s", version)
					break
				}
			}
		}
	} else {
		color.Yellow("  ℹ️  No go.mod file in current directory")
	}
}
