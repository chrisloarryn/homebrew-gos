package status

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// ShowEnvironment displays environment variables and PATH information
func ShowEnvironment() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	
	// Expected values
	expectedGoroot := filepath.Join(os.Getenv("HOME"), ".g", "go")
	expectedGopath := filepath.Join(os.Getenv("HOME"), "go")
	
	envVars := map[string]string{
		"GOROOT": expectedGoroot,
		"GOPATH": expectedGopath,
		"GOPROXY": "",
		"GOSUMDB": "",
		"GOMODCACHE": "",
	}
	
	for envVar, expected := range envVars {
		if value := os.Getenv(envVar); value != "" {
			if expected != "" && value == expected {
				green.Printf("  ✅ %s: %s\n", envVar, value)
			} else if expected != "" {
				yellow.Printf("  ⚠️  %s: %s (expected: %s)\n", envVar, value, expected)
			} else {
				blue.Printf("  ℹ️  %s: %s\n", envVar, value)
			}
		} else {
			if expected != "" {
				yellow.Printf("  ❌ %s: (not set, should be: %s)\n", envVar, expected)
			} else {
				fmt.Printf("  %s: (not set)\n", envVar)
			}
		}
	}

	// Show PATH entries related to Go
	showGoPathEntries()
}

// showGoPathEntries displays Go-related PATH entries
func showGoPathEntries() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)

	fmt.Println("  PATH (Go-related entries):")
	path := os.Getenv("PATH")
	pathEntries := strings.Split(path, ":")
	hasGoBin := false
	hasGBin := false
	
	for _, entry := range pathEntries {
		if strings.Contains(entry, "go") || strings.Contains(entry, ".g") {
			if strings.Contains(entry, ".g/bin") {
				hasGBin = true
				green.Printf("    ✅ %s\n", entry)
			} else if strings.Contains(entry, "go/bin") {
				hasGoBin = true
				green.Printf("    ✅ %s\n", entry)
			} else {
				blue.Printf("    ℹ️  %s\n", entry)
			}
		}
	}
	
	if !hasGBin {
		yellow.Println("    ⚠️  ~/.g/bin not found in PATH")
	}
	if !hasGoBin {
		yellow.Println("    ⚠️  $GOPATH/bin not found in PATH")
	}
}
