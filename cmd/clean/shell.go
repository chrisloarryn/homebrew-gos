package clean

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cristobalcontreras/gos/cmd/common"
)

// CleanShellConfig cleans Go configuration from shell files
func CleanShellConfig() {
	homeDir := common.GetHomeDir()
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

// cleanShellFile removes Go-related lines from a shell configuration file
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

// containsGoConfig checks if a line contains Go-related configuration
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
