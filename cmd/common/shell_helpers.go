package common

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// DetectCurrentShell detects the current shell being used
func DetectCurrentShell() string {
	// First check SHELL environment variable
	if shell := os.Getenv("SHELL"); shell != "" {
		return filepath.Base(shell)
	}
	
	// Check parent process name (works in most cases)
	if runtime.GOOS == "windows" {
		// On Windows, check common shells
		if os.Getenv("PSModulePath") != "" {
			return "powershell"
		}
		if os.Getenv("BASH") != "" || os.Getenv("BASH_VERSION") != "" {
			return "bash"
		}
		return "cmd"
	}
	
	// On Unix-like systems, check ZSH_VERSION or BASH_VERSION
	if os.Getenv("ZSH_VERSION") != "" {
		return "zsh"
	}
	if os.Getenv("BASH_VERSION") != "" {
		return "bash"
	}
	
	// Default fallback
	return "unknown"
}

// GetShellFileForCurrentShell returns the appropriate shell configuration file
func GetShellFileForCurrentShell(shell, homeDir string) string {
	switch shell {
	case "zsh":
		return ZshrcFile
	case "bash":
		// Check which bash file exists
		bashFiles := []string{BashrcFile, BashProfileFile, ProfileFile}
		for _, file := range bashFiles {
			if _, err := os.Stat(filepath.Join(homeDir, file)); err == nil {
				return file
			}
		}
		return BashrcFile // Default
	case "powershell":
		return PowerShellProfile
	case "cmd":
		return PowerShellProfile // PowerShell as alternative
	default:
		// Try to detect based on existing files
		if runtime.GOOS == "windows" {
			return PowerShellProfile
		}
		// Unix-like: check what exists
		candidates := []string{ZshrcFile, BashrcFile, BashProfileFile, ProfileFile}
		for _, file := range candidates {
			if _, err := os.Stat(filepath.Join(homeDir, file)); err == nil {
				return file
			}
		}
		return ZshrcFile // Default for Unix-like systems
	}
}

// ExecuteWithShell runs a command with appropriate shell for the OS
func ExecuteWithShell(command string) error {
	if runtime.GOOS == "windows" {
		// Try PowerShell first, then cmd
		if IsCommandAvailable("powershell") {
			cmd := exec.Command("powershell", "-Command", command)
			return cmd.Run()
		}
		cmd := exec.Command("cmd", "/C", command)
		return cmd.Run()
	}
	// Unix-like systems
	cmd := exec.Command("bash", "-c", command)
	return cmd.Run()
}
