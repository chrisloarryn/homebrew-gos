package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// GetHomeDir gets the user's home directory
func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

// UpdatePathForGobrew updates PATH for gobrew version manager
func UpdatePathForGobrew(homeDir string) {
	currentPath := os.Getenv("PATH")
	gobrewCurrentBin := filepath.Join(homeDir, GobrewDir, "current", "bin")
	gobrewBin := filepath.Join(homeDir, GobrewDir, "bin")

	pathSep := ";"
	if runtime.GOOS != "windows" {
		pathSep = ":"
	}

	// Remove any existing gobrew paths from PATH to avoid duplicates
	pathParts := strings.Split(currentPath, pathSep)
	var cleanPaths []string

	for _, part := range pathParts {
		if !strings.Contains(part, GobrewDir) {
			cleanPaths = append(cleanPaths, part)
		}
	}

	// Prepend gobrew paths
	finalPaths := []string{gobrewCurrentBin, gobrewBin}
	finalPaths = append(finalPaths, cleanPaths...)
	newPath := strings.Join(finalPaths, pathSep)

	os.Setenv("PATH", newPath)
}

// UpdatePathForGoEnvironment updates PATH with required Go paths
func UpdatePathForGoEnvironment() bool {
	currentPath := os.Getenv("PATH")
	homeDir := GetHomeDir()

	// Use gobrew paths
	requiredPaths := []string{
		filepath.Join(homeDir, ".gobrew", "current", "bin"),
		filepath.Join(homeDir, ".gobrew", "bin"),
		filepath.Join(homeDir, "go", "bin"),
	}

	newPaths := []string{}
	for _, reqPath := range requiredPaths {
		if !strings.Contains(currentPath, reqPath) {
			newPaths = append(newPaths, reqPath)
		}
	}

	if len(newPaths) > 0 {
		newPath := strings.Join(newPaths, ":") + ":" + currentPath
		os.Setenv("PATH", newPath)
		green := color.New(color.FgGreen)
		green.Printf("✅ Updated PATH with: %s\n", strings.Join(newPaths, ", "))
		return true
	}

	return false
}

// CleanGoPathsFromEnvironment removes all Go-related paths from PATH
func CleanGoPathsFromEnvironment() {
	currentPath := os.Getenv("PATH")
	pathSep := ":"
	if runtime.GOOS == "windows" {
		pathSep = ";"
	}

	pathParts := strings.Split(currentPath, pathSep)
	var cleanPaths []string

	goPatterns := []string{
		"/go/bin",           // Standard Go installation
		"/usr/local/go/bin", // System-wide Go installation
		"/.g/",              // g version manager
		"/gobrew/",          // gobrew version manager
		"/sdk/go",           // Go SDK installations
		"/golang/",          // Other Go installations
		"/versions/go",      // Versioned Go installations
		"google-cloud-sdk",  // Exclude Google Cloud SDK as it's not Go language
	}

	for _, part := range pathParts {
		shouldKeep := true

		// Check if this path contains any Go-related patterns
		for _, pattern := range goPatterns {
			if strings.Contains(part, pattern) {
				// Special case: keep Google Cloud SDK
				if pattern == "google-cloud-sdk" {
					continue
				}
				shouldKeep = false
				break
			}
		}

		if shouldKeep {
			cleanPaths = append(cleanPaths, part)
		}
	}

	// Set the cleaned PATH
	newPath := strings.Join(cleanPaths, pathSep)
	os.Setenv("PATH", newPath)
}

// UpdatePathForVersionManagerClean cleanly updates PATH for version managers
func UpdatePathForVersionManagerClean() {
	// First, clean all Go paths
	CleanGoPathsFromEnvironment()

	homeDir := GetHomeDir()
	currentPath := os.Getenv("PATH")
	pathSep := ":"
	if runtime.GOOS == "windows" {
		pathSep = ";"
	}

	var newGoPaths []string

	// Use gobrew paths
	gobrewCurrentBin := filepath.Join(homeDir, GobrewDir, "current", "bin")
	gobrewBin := filepath.Join(homeDir, GobrewDir, "bin")
	newGoPaths = append(newGoPaths, gobrewCurrentBin, gobrewBin)

	// Add user's Go bin for installed packages
	userGoBin := filepath.Join(homeDir, "go", "bin")
	newGoPaths = append(newGoPaths, userGoBin)

	// Prepend the new Go paths to the cleaned PATH
	if len(newGoPaths) > 0 {
		finalPath := strings.Join(newGoPaths, pathSep) + pathSep + currentPath
		os.Setenv("PATH", finalPath)
	}
}

// PromptUserForPathCleanup asks user if they want to clean their PATH and provides options
func PromptUserForPathCleanup() {
	fmt.Println()
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	yellow.Println("🧹 PATH cleanup detected multiple Go installations!")
	fmt.Println()
	fmt.Println("Would you like to:")
	fmt.Println("  1️⃣  Generate a cleanup script to run")
	fmt.Println("  2️⃣  Show manual cleanup instructions")
	fmt.Println("  3️⃣  Skip for now")
	fmt.Println()

	blue.Print("Enter your choice (1/2/3): ")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		generateCleanupScript()
	case "2":
		showManualCleanupInstructions()
	case "3":
		fmt.Println("⏭️  Skipping PATH cleanup for now.")
	default:
		green.Println("ℹ️  Invalid choice. Showing manual instructions...")
		showManualCleanupInstructions()
	}
}

// generateCleanupScript creates a script that user can source to clean their PATH
func generateCleanupScript() {
	homeDir := GetHomeDir()
	scriptPath := filepath.Join(homeDir, "clean-go-path.sh")

	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	// Use gobrew paths
	vmPaths := `export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$HOME/go/bin"`

	scriptContent := `#!/bin/bash
# Go PATH Cleanup Script
# Generated by gos - Go Version Manager CLI

echo "🧹 Cleaning Go paths from current session..."

# Save current PATH
export ORIGINAL_PATH="$PATH"

# Remove Go-related paths
NEW_PATH=""
IFS=':' read -ra PATH_ARRAY <<< "$PATH"
for path in "${PATH_ARRAY[@]}"; do
    # Skip Go-related paths
    if [[ "$path" == *"/go/bin"* ]] || \
       [[ "$path" == *"/usr/local/go/bin"* ]] || \
       [[ "$path" == *"/.g/"* ]] || \
       [[ "$path" == *"/gobrew/"* ]] || \
       [[ "$path" == *"/sdk/go"* ]] || \
       [[ "$path" == *"/golang/"* ]] || \
       [[ "$path" == *"/versions/go"* ]]; then
        # Skip unless it's Google Cloud SDK (keep it)
        if [[ "$path" == *"google-cloud-sdk"* ]]; then
            if [ -z "$NEW_PATH" ]; then
                NEW_PATH="$path"
            else
                NEW_PATH="$NEW_PATH:$path"
            fi
        fi
        continue
    fi

    # Add non-Go paths
    if [ -z "$NEW_PATH" ]; then
        NEW_PATH="$path"
    else
        NEW_PATH="$NEW_PATH:$path"
    fi
done

# Set cleaned PATH
export PATH="$NEW_PATH"

# Add version manager paths
` + vmPaths + `:$PATH"

echo "✅ PATH cleaned and version manager paths added!"
echo "🔍 Current Go version:"
go version 2>/dev/null || echo "❌ Go not found in PATH - you may need to install a version first"

echo ""
echo "💡 To make this permanent, add this line to your ~/.zshrc or ~/.bashrc:"
echo "` + vmPaths + `:$PATH\""
echo ""
echo "🔄 To restore original PATH: export PATH=\"$ORIGINAL_PATH\""
`

	// Write the script
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		yellow.Printf("❌ Error creating script: %v\n", err)
		fmt.Println("📋 Here are the manual commands instead:")
		showManualCleanupInstructions()
		return
	}

	green.Printf("✅ Cleanup script created: %s\n", scriptPath)
	fmt.Println()
	blue.Println("🚀 To clean your PATH in this session, run:")
	fmt.Printf("   source %s\n", scriptPath)
	fmt.Println()
	blue.Println("🔧 To make it permanent, add this to your ~/.zshrc or ~/.bashrc:")
	pathLine := strings.TrimSuffix(strings.Replace(vmPaths, `export PATH="`, "", 1), `"`)
	fmt.Printf("   %s:$PATH\n", pathLine)
	fmt.Println()
	yellow.Println("⚠️  Remember to restart your terminal or run 'source ~/.zshrc' after editing your shell config!")
}

// showManualCleanupInstructions displays detailed manual cleanup steps
func showManualCleanupInstructions() {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	fmt.Println()
	blue.Println("📋 Manual PATH Cleanup Instructions:")
	fmt.Println()

	yellow.Println("1️⃣  Edit your shell configuration file:")
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		fmt.Println("   nano ~/.zshrc     # for zsh")
		fmt.Println("   nano ~/.bashrc    # for bash")
	} else {
		fmt.Println("   Edit your shell profile file")
	}
	fmt.Println()

	yellow.Println("2️⃣  Remove or comment out these types of lines:")
	fmt.Println("   # export PATH=/usr/local/go/bin:$PATH")
	fmt.Println("   # export PATH=$HOME/sdk/go1.xx.x/bin:$PATH")
	fmt.Println("   # export PATH=/opt/go/bin:$PATH")
	fmt.Println()

	yellow.Println("3️⃣  Add only the version manager path:")
	green.Println("   export PATH=\"$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$HOME/go/bin:$PATH\"")
	fmt.Println()

	yellow.Println("4️⃣  Save the file and reload your shell:")
	fmt.Println("   source ~/.zshrc   # or source ~/.bashrc")
	fmt.Println()

	blue.Println("🔍 After cleanup, verify with:")
	fmt.Println("   go version")
	fmt.Println("   which go")
	fmt.Println()
}
