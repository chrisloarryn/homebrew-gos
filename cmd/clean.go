package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
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
		deepCleanGo(force)
	},
}

func init() {
	cleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

func deepCleanGo(force bool) {
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

	blue.Println("\n▸ Cleaning existing Go cache and modules…")
	cleanGoCache()

	blue.Println("\n▸ Removing Homebrew installations…")
	cleanHomebrewGo()

	blue.Println("\n▸ Removing manual system installations…")
	cleanSystemGo()

	blue.Println("\n▸ Removing user directories with special permissions…")
	cleanUserDirectories()

	blue.Println("\n▸ Removing other managers and directories…")
	cleanOtherManagers()

	blue.Println("\n▸ Cleaning shell configuration…")
	cleanShellConfig()

	// Clear command hash
	exec.Command("hash", "-r").Run()

	green.Println("\n✅ Complete Go cleanup finished.")
	yellow.Println("📋 Backups of your configuration files were created.")
	yellow.Println("🔄 Run 'source ~/.zshrc' or open a new terminal.")
}

func cleanGoCache() {
	// Try to clean with go command if available
	if _, err := exec.LookPath("go"); err == nil {
		fmt.Println("  Running go clean -modcache...")
		exec.Command("go", "clean", "-modcache").Run()
		
		fmt.Println("  Running go clean -cache...")
		exec.Command("go", "clean", "-cache").Run()
	}
}

func cleanHomebrewGo() {
	if _, err := exec.LookPath("brew"); err != nil {
		return
	}

	// Get list of Go formulas
	cmd := exec.Command("brew", "list", "--formula")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go") && (line == "go" || strings.Contains(line, "go@")) {
			fmt.Printf("  – brew uninstall %s\n", line)
			uninstallCmd := exec.Command("brew", "uninstall", "--ignore-dependencies", "--force", line)
			uninstallCmd.Run()
		}
	}
}

func cleanSystemGo() {
	// Remove manual system installations
	exec.Command("sudo", "rm", "-rf", "/usr/local/go").Run()
}

func cleanUserDirectories() {
	homeDir := os.Getenv("HOME")
	
	// Clean ~/go directory
	goDir := filepath.Join(homeDir, "go")
	if _, err := os.Stat(goDir); err == nil {
		fmt.Printf("  Fixing permissions in %s...\n", goDir)
		fixPermissions(goDir)
		if err := os.RemoveAll(goDir); err != nil {
			fmt.Printf("  Using sudo to remove %s...\n", goDir)
			exec.Command("sudo", "rm", "-rf", goDir).Run()
		}
	}

	// Clean Go cache directories
	cacheDirs := []string{
		filepath.Join(homeDir, ".cache", "go-build"),
		filepath.Join(homeDir, "Library", "Caches", "go-build"),
	}

	for _, cacheDir := range cacheDirs {
		if _, err := os.Stat(cacheDir); err == nil {
			fmt.Printf("  Removing cache: %s\n", cacheDir)
			fixPermissions(cacheDir)
			os.RemoveAll(cacheDir)
		}
	}
}

func cleanOtherManagers() {
	homeDir := os.Getenv("HOME")
	
	managerDirs := []string{
		filepath.Join(homeDir, "sdk"),
		filepath.Join(homeDir, ".gvm"),
		filepath.Join(homeDir, ".goenv"),
		filepath.Join(homeDir, ".g"),
	}

	for _, dir := range managerDirs {
		if _, err := os.Stat(dir); err == nil {
			// For sdk, only remove go* directories
			if strings.HasSuffix(dir, "sdk") {
				if entries, err := os.ReadDir(dir); err == nil {
					for _, entry := range entries {
						if strings.HasPrefix(entry.Name(), "go") {
							goSdkDir := filepath.Join(dir, entry.Name())
							os.RemoveAll(goSdkDir)
						}
					}
				}
			} else {
				os.RemoveAll(dir)
			}
		}
	}
}

func cleanShellConfig() {
	homeDir := os.Getenv("HOME")
	shellFiles := []string{
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bash_profile"),
		filepath.Join(homeDir, ".bashrc"),
	}

	for _, shellFile := range shellFiles {
		if _, err := os.Stat(shellFile); err == nil {
			cleanShellFile(shellFile)
		}
	}
}

func cleanShellFile(filename string) {
	// Create backup
	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("%s.backup.%s", filename, timestamp)
	
	input, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	
	os.WriteFile(backupFile, input, 0644)

	// Filter out Go-related lines
	lines := strings.Split(string(input), "\n")
	var filteredLines []string

	for _, line := range lines {
		if !containsGoConfig(line) {
			filteredLines = append(filteredLines, line)
		}
	}

	output := strings.Join(filteredLines, "\n")
	os.WriteFile(filename, []byte(output), 0644)
}

func containsGoConfig(line string) bool {
	goPatterns := []string{
		"go/bin",
		"GOPATH",
		"GOROOT",
		".gvm",
		".goenv",
		".g/bin",
		"Go Version Manager",
	}

	for _, pattern := range goPatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	return false
}

func fixPermissions(dir string) {
	// Recursively fix permissions
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			os.Chmod(path, 0755)
		} else {
			os.Chmod(path, 0644)
		}
		return nil
	})
}
