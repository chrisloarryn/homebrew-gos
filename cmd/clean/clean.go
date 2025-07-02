package clean

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// NewCleanCmd creates the clean command
func NewCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Deep clean all Go installations",
		Long: `Perform a comprehensive cleanup of all Go installations, including:
- Go cache and modules
- Homebrew installations
- Manual system installations
- User directories with special permissions
- Shell configuration cleanup`,
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")
			DeepCleanGo(force)
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	// Add path subcommand
	pathCmd := &cobra.Command{
		Use:   "path",
		Short: "Clean PATH from Go installation conflicts",
		Long: `Clean your PATH environment variable from conflicting Go installations.
This will generate a script or modify your shell configuration to remove
multiple Go paths and keep only the version manager paths.`,
		Run: func(cmd *cobra.Command, args []string) {
			interactive, _ := cmd.Flags().GetBool("interactive")
			script, _ := cmd.Flags().GetBool("script")
			CleanPathConflicts(interactive, script)
		},
	}

	pathCmd.Flags().BoolP("interactive", "i", true, "Interactive mode with user prompts")
	pathCmd.Flags().BoolP("script", "s", false, "Generate cleanup script only")

	cmd.AddCommand(pathCmd)

	return cmd
}

// DeepCleanGo performs comprehensive Go cleanup
func DeepCleanGo(force bool) {
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)

	red.Println("🗑️  Complete Go system cleanup...")

	if !force {
		yellow.Println("\n⚠️  WARNING: This will remove ALL Go installations and configurations!")
		fmt.Print("Are you sure you want to continue? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "y" && input != "yes" {
			yellow.Println("Cleanup cancelled.")
			return
		}
	}

	// Create main progress bar for cleanup stages
	totalStages := 6
	mainBar := progressbar.NewOptions(totalStages,
		progressbar.OptionSetDescription("🧹 Deep cleanup progress"),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "█", SaucerHead: "█", SaucerPadding: "░", BarStart: "[", BarEnd: "]",
		}),
		progressbar.OptionShowCount(),
	)

	fmt.Println()

	blue.Println("▸ Cleaning existing Go cache and modules…")
	CleanGoCache()
	mainBar.Add(1)

	blue.Println("\n▸ Removing Homebrew installations…")
	CleanHomebrewGo()
	mainBar.Add(1)

	blue.Println("\n▸ Removing manual system installations…")
	CleanSystemGo()
	mainBar.Add(1)

	blue.Println("\n▸ Removing user directories with special permissions…")
	CleanUserDirectories()
	mainBar.Add(1)

	blue.Println("\n▸ Removing other managers and directories…")
	CleanOtherManagers()
	mainBar.Add(1)

	blue.Println("\n▸ Cleaning shell configuration…")
	CleanShellConfig()
	mainBar.Add(1)

	// Clear command hash
	exec.Command("hash", "-r").Run()

	mainBar.Finish()
	fmt.Println()
	green.Println("✅ Complete Go cleanup finished.")
	yellow.Println("📋 Backups of your configuration files were created.")
	yellow.Println("🔄 Run 'source ~/.zshrc' or open a new terminal.")
}

// CleanPathConflicts handles PATH cleanup for Go installation conflicts
func CleanPathConflicts(interactive, scriptOnly bool) {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	blue.Println("🧹 Go PATH Conflict Cleanup")
	fmt.Println()

	// Check for conflicts
	conflicts := detectPathConflicts()
	if len(conflicts) == 0 {
		green.Println("✅ No PATH conflicts detected!")
		fmt.Println("💡 Your PATH appears to be clean.")
		return
	}

	// Show detected conflicts
	yellow.Printf("⚠️  Detected %d conflicting Go paths in your environment:\n", len(conflicts))
	for _, conflict := range conflicts {
		fmt.Printf("  • %s\n", conflict)
	}
	fmt.Println()

	if scriptOnly {
		// Clean current session PATH directly
		cleanCurrentSessionPath()
		return
	}

	if interactive {
		// Ask user what they want to do
		fmt.Println("Choose an option:")
		fmt.Println("  1️⃣  Clean shell configuration files permanently")
		fmt.Println("  2️⃣  Clean current session PATH only")
		fmt.Println("  3️⃣  Show manual instructions")
		fmt.Println("  4️⃣  Cancel")
		fmt.Println()

		blue.Print("Enter your choice (1/2/3/4): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			cleanShellConfigFromPath()
		case "2":
			cleanCurrentSessionPath()
		case "3":
			showPathCleanupInstructions()
		case "4":
			fmt.Println("⏭️  Cleanup cancelled.")
		default:
			yellow.Println("Invalid choice. Showing manual instructions...")
			showPathCleanupInstructions()
		}
	} else {
		// Non-interactive: auto-clean
		cleanShellConfigFromPath()
	}
}

// detectPathConflicts finds conflicting Go paths in the current environment
func detectPathConflicts() []string {
	currentPath := os.Getenv("PATH")
	pathSep := ":"
	if os.Getenv("OS") == "Windows_NT" {
		pathSep = ";"
	}

	pathParts := strings.Split(currentPath, pathSep)
	var conflicts []string

	// Check each path part for conflicts
	for _, part := range pathParts {
		if isConflictingGoPath(part) {
			conflicts = append(conflicts, part)
		}
	}

	return removeDuplicates(conflicts)
}

// isConflictingGoPath checks if a path is a conflicting Go installation
func isConflictingGoPath(path string) bool {
	// Patterns for conflicting Go installations (not version managers)
	conflictPatterns := []string{
		"/usr/local/go/bin",
		"/sdk/go",
		"/golang/",
		"/opt/go/bin",
		"/opt/homebrew/bin/go",
		"/usr/bin/go",
		"/snap/go",
	}

	// Always skip gobrew version manager paths - these are NOT conflicts
	if strings.Contains(path, "/.gobrew/") {
		return false
	}

	// Check for conflicting patterns
	for _, pattern := range conflictPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	return false
}

// removeDuplicates removes duplicate entries from a slice
func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// cleanShellConfigFromPath automatically cleans shell configuration files
func cleanShellConfigFromPath() {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	blue.Println("🔧 Automatically cleaning shell configuration files...")
	fmt.Println()

	// First do the existing shell cleanup
	CleanShellConfig()

	// Then add the proper version manager path
	addVersionManagerPathToShell()

	fmt.Println()
	green.Println("✅ Shell configuration cleaned successfully!")
	yellow.Println("🔄 Please restart your terminal or run 'source ~/.zshrc' to apply changes.")
}

// addVersionManagerPathToShell adds the correct version manager path to shell config
func addVersionManagerPathToShell() {
	homeDir := common.GetHomeDir()

	var shellFile string
	var pathLine string

	// Determine shell file and path
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		shellFile = filepath.Join(homeDir, ".zshrc")
	} else if strings.Contains(shell, "bash") {
		shellFile = filepath.Join(homeDir, ".bashrc")
		if _, err := os.Stat(shellFile); os.IsNotExist(err) {
			shellFile = filepath.Join(homeDir, ".bash_profile")
		}
	} else {
		// Default to .zshrc on modern systems
		shellFile = filepath.Join(homeDir, ".zshrc")
	}

	// Use gobrew as the version manager
	pathLine = `export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$HOME/go/bin:$PATH"`

	// Read existing content
	content := ""
	if data, err := os.ReadFile(shellFile); err == nil {
		content = string(data)
	}

	// Check if path line already exists
	if strings.Contains(content, pathLine) {
		fmt.Printf("✅ Version manager path already configured in %s\n", filepath.Base(shellFile))
		return
	}

	// Add the path line
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += "\n# Go Version Manager PATH\n"
	content += pathLine + "\n"

	// Write back
	if err := os.WriteFile(shellFile, []byte(content), 0644); err != nil {
		yellow := color.New(color.FgYellow)
		yellow.Printf("⚠️  Could not write to %s: %v\n", shellFile, err)
		yellow.Println("Please add this line manually:")
		fmt.Println("  " + pathLine)
	} else {
		green := color.New(color.FgGreen)
		green.Printf("✅ Added version manager path to %s\n", filepath.Base(shellFile))
	}
}

// cleanCurrentSessionPath cleans the PATH of the current session directly
func cleanCurrentSessionPath() {
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("🧹 Cleaning PATH for current session...")
	fmt.Println()

	// Get current PATH
	currentPath := os.Getenv("PATH")
	pathSep := ":"
	if os.Getenv("OS") == "Windows_NT" {
		pathSep = ";"
	}

	pathParts := strings.Split(currentPath, pathSep)
	var cleanedParts []string

	// Use gobrew version manager paths
	homeDir := common.GetHomeDir()
	vmPaths := []string{
		filepath.Join(homeDir, ".gobrew", "current", "bin"),
		filepath.Join(homeDir, ".gobrew", "bin"),
		filepath.Join(homeDir, "go", "bin"),
	}

	// Add version manager paths first
	for _, vmPath := range vmPaths {
		cleanedParts = append(cleanedParts, vmPath)
	}

	// Process each path part
	for _, part := range pathParts {
		// Skip empty parts
		if part == "" {
			continue
		}

		// Skip conflicting Go installations
		if isConflictingGoPath(part) {
			yellow.Printf("  Removing conflicting path: %s\n", part)
			continue
		}

		// Skip if it's already a version manager path (avoid duplicates)
		isVMPath := false
		for _, vmPath := range vmPaths {
			if part == vmPath {
				isVMPath = true
				break
			}
		}
		if isVMPath {
			continue
		}

		// Add non-conflicting paths
		cleanedParts = append(cleanedParts, part)
	}

	// Build the new PATH
	newPath := strings.Join(cleanedParts, pathSep)

	// Set the new PATH for the current process
	os.Setenv("PATH", newPath)

	fmt.Println()
	green.Println("✅ PATH cleaned for current session!")

	// Test Go version
	blue.Println("� Testing Go installation:")
	if output, err := exec.Command("go", "version").Output(); err == nil {
		green.Printf("✅ %s\n", strings.TrimSpace(string(output)))
	} else {
		yellow.Println("⚠️  Go not found in cleaned PATH")
	}

	fmt.Println()
	yellow.Println("� This change is temporary for this session only.")
	yellow.Println("💡 For permanent changes, use option 1 (Clean shell configuration files).")
	fmt.Println()
}

// showPathCleanupInstructions shows manual cleanup instructions
func showPathCleanupInstructions() {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	fmt.Println()
	blue.Println("📋 Manual PATH Cleanup Instructions:")
	fmt.Println()

	yellow.Println("1️⃣  Edit your shell configuration file:")
	fmt.Println("   nano ~/.zshrc     # for zsh")
	fmt.Println("   nano ~/.bashrc    # for bash")
	fmt.Println()

	yellow.Println("2️⃣  Remove these conflicting lines:")
	fmt.Println("   # export PATH=/usr/local/go/bin:$PATH")
	fmt.Println("   # export PATH=$HOME/sdk/go*/bin:$PATH")
	fmt.Println("   # export PATH=/opt/go/bin:$PATH")
	fmt.Println()

	yellow.Println("3️⃣  Add only the gobrew version manager path:")
	green.Println("   export PATH=\"$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$HOME/go/bin:$PATH\"")
	fmt.Println()

	yellow.Println("4️⃣  Save and reload:")
	fmt.Println("   source ~/.zshrc")
	fmt.Println()
}
