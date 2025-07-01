package common

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// Constants for repeated paths
const (
	UsrLocalBinG = "/usr/local/bin/g"
	ZshrcFile = ".zshrc"
	BashrcFile = ".bashrc"
	BashProfileFile = ".bash_profile"
	ProfileFile = ".profile"
	PowerShellProfile = "Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
	GobrewDir = ".gobrew"
)

// GetHomeDir gets the user's home directory
func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

// IsCommandAvailable checks if a command is available in PATH
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// CheckVersionManagerAvailable checks if any version manager is available
func CheckVersionManagerAvailable() bool {
	if IsCommandAvailable("gobrew") || IsCommandAvailable("g") {
		return true
	}

	color.Red("❌ Error: No version manager is installed.")
	color.Yellow("💡 Run first: gos setup")
	return false
}

// CheckGInstalled checks if the 'g' version manager is installed (legacy function)
func CheckGInstalled() bool {
	_, err := exec.LookPath("g")
	if err != nil {
		color.Red("❌ Error: The 'g' manager is not installed.")
		color.Yellow("💡 Run first: gos setup")
		return false
	}
	return true
}

// IsGInstalled checks if the 'g' version manager is installed with path detection
func IsGInstalled() bool {
	// Check for gobrew on Windows first
	if runtime.GOOS == "windows" && IsCommandAvailable("gobrew") {
		return true
	}
	
	homeDir := GetHomeDir()
	gPaths := []string{
		filepath.Join(homeDir, ".g", "bin", "g"),
		filepath.Join(homeDir, "go", "bin", "g"),
		UsrLocalBinG,
	}
	
	for _, path := range gPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	
	return IsCommandAvailable("g")
}

// GetInstalledVersions returns installed versions using 'g'
func GetInstalledVersions() []string {
	homeDir := GetHomeDir()
	gPaths := []string{
		filepath.Join(homeDir, ".g", "bin", "g"),
		filepath.Join(homeDir, "go", "bin", "g"),
		UsrLocalBinG,
	}
	
	var gPath string
	for _, path := range gPaths {
		if _, err := os.Stat(path); err == nil {
			gPath = path
			break
		}
	}
	
	if gPath == "" {
		return []string{}
	}
	
	cmd := exec.Command(gPath, "list")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	
	lines := strings.Split(string(output), "\n")
	var versions []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && line != "g" { // Skip empty lines and header
			versions = append(versions, line)
		}
	}
	
	return versions
}

// GetGobrewVersions returns installed versions using gobrew
func GetGobrewVersions() []string {
	if !IsCommandAvailable("gobrew") {
		return []string{}
	}
	
	cmd := exec.Command("gobrew", "ls")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	
	lines := strings.Split(string(output), "\n")
	var versions []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "=>") { // Skip empty lines and current indicator
			versions = append(versions, line)
		}
	}
	
	return versions
}

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
		return ".zshrc"
	case "bash":
		// Check which bash file exists
		bashFiles := []string{".bashrc", ".bash_profile", ".profile"}
		for _, file := range bashFiles {
			if _, err := os.Stat(filepath.Join(homeDir, file)); err == nil {
				return file
			}
		}
		return ".bashrc" // Default
	case "powershell":
		return "Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
	case "cmd":
		return "Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1" // PowerShell as alternative
	default:
		// Try to detect based on existing files
		if runtime.GOOS == "windows" {
			return "Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
		}
		// Unix-like: check what exists
		candidates := []string{".zshrc", ".bashrc", ".bash_profile", ".profile"}
		for _, file := range candidates {
			if _, err := os.Stat(filepath.Join(homeDir, file)); err == nil {
				return file
			}
		}
		return ".zshrc" // Default for Unix-like systems
	}
}

// UpdatePathForVersionManager updates the PATH environment variable to include the current Go version
func UpdatePathForVersionManager() {
	homeDir := GetHomeDir()
	
	if IsCommandAvailable("gobrew") {
		UpdatePathForGobrew(homeDir)
	} else if IsCommandAvailable("g") {
		UpdatePathForG(homeDir)
	}
}

// UpdatePathForGobrew updates PATH for gobrew version manager
func UpdatePathForGobrew(homeDir string) {
	currentPath := os.Getenv("PATH")
	gobrewCurrentBin := filepath.Join(homeDir, ".gobrew", "current", "bin")
	gobrewBin := filepath.Join(homeDir, ".gobrew", "bin")
	
	pathSep := ";"
	if runtime.GOOS != "windows" {
		pathSep = ":"
	}

	// Remove any existing gobrew paths from PATH to avoid duplicates
	pathParts := strings.Split(currentPath, pathSep)
	var cleanPaths []string

	for _, part := range pathParts {
		if !strings.Contains(part, ".gobrew") {
			cleanPaths = append(cleanPaths, part)
		}
	}

	// Prepend gobrew paths
	finalPaths := []string{gobrewCurrentBin, gobrewBin}
	finalPaths = append(finalPaths, cleanPaths...)
	newPath := strings.Join(finalPaths, pathSep)

	os.Setenv("PATH", newPath)
}

// UpdatePathForG updates PATH for g version manager
func UpdatePathForG(homeDir string) {
	currentPath := os.Getenv("PATH")
	gGoBin := filepath.Join(homeDir, ".g", "go", "bin")
	gBin := filepath.Join(homeDir, ".g", "bin")
	
	pathSep := ";"
	if runtime.GOOS != "windows" {
		pathSep = ":"
	}

	// Remove any existing g paths from PATH to avoid duplicates
	pathParts := strings.Split(currentPath, pathSep)
	var cleanPaths []string

	for _, part := range pathParts {
		if !strings.Contains(part, ".g") {
			cleanPaths = append(cleanPaths, part)
		}
	}

	// Prepend g paths
	finalPaths := []string{gGoBin, gBin}
	finalPaths = append(finalPaths, cleanPaths...)
	newPath := strings.Join(finalPaths, pathSep)

	os.Setenv("PATH", newPath)
}

// GetCurrentGoVersion returns the currently active Go version
func GetCurrentGoVersion() string {
	if IsCommandAvailable("gobrew") {
		return getCurrentGoVersionWithGobrew()
	} else if IsCommandAvailable("g") {
		return getCurrentGoVersionWithG()
	}
	
	// Fallback: check system Go
	if output, err := exec.Command("go", "version").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	
	return ""
}

// getCurrentGoVersionWithGobrew gets current version using gobrew
func getCurrentGoVersionWithGobrew() string {
	cmd := exec.Command("gobrew", "ls")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "*") || strings.Contains(line, "current") {
			return strings.ReplaceAll(line, "*", "")
		}
	}
	return ""
}

// getCurrentGoVersionWithG gets current version using g
func getCurrentGoVersionWithG() string {
	// Try to get current version
	if currentCmd := exec.Command("g", "which"); currentCmd != nil {
		if currentOutput, currentErr := currentCmd.Output(); currentErr == nil {
			return strings.TrimSpace(string(currentOutput))
		}
	}
	
	// If that doesn't work, check symlink
	homeDir := GetHomeDir()
	goLink := filepath.Join(homeDir, ".g", "go")
	if target, err := os.Readlink(goLink); err == nil {
		return filepath.Base(target)
	}
	
	return ""
}

// GetSystemGoInfo returns information about system Go installation
func GetSystemGoInfo() (version string, goroot string, found bool) {
	// Check if Go is installed directly
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version = strings.TrimSpace(string(output))
		
		// Try to get GOROOT to see where it's installed
		if gorootOutput, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
			goroot = strings.TrimSpace(string(gorootOutput))
		}
		
		return version, goroot, true
	}
	
	return "", "", false
}

// FindGPath returns the path to the g version manager executable
func FindGPath() string {
	homeDir := GetHomeDir()
	gPaths := []string{
		filepath.Join(homeDir, ".g", "bin", "g"),
		filepath.Join(homeDir, "go", "bin", "g"),
		UsrLocalBinG,
	}
	
	for _, path := range gPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return ""
}

// SetupGoEnvironment sets up the Go environment variables and PATH
func SetupGoEnvironment() {
	green := color.New(color.FgGreen)
	
	homeDir := GetHomeDir()
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")

	// Set environment variables for current session
	os.Setenv("GOPATH", expectedGopath)
	os.Setenv("GOROOT", expectedGoroot)

	// Update PATH for current session
	if UpdatePathForGoEnvironment() {
		green.Println("✅ PATH updated for current session")
	}
}

// UpdatePathForGoEnvironment updates PATH with required Go paths
func UpdatePathForGoEnvironment() bool {
	currentPath := os.Getenv("PATH")
	homeDir := GetHomeDir()
	
	requiredPaths := []string{
		filepath.Join(homeDir, ".g", "bin"),
		filepath.Join(homeDir, ".g", "go", "bin"),
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

// VerifyGoInstallation verifies that Go is available and shows version
func VerifyGoInstallation() bool {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	if _, err := exec.LookPath("go"); err != nil {
		red.Println("❌ Go not found in PATH")
		yellow.Println("💡 You may need to restart your terminal or run:")
		yellow.Println("   source ~/.zshrc")
		return false
	}

	// Show Go version
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("✅ %s\n", version)
	}
	
	return true
}

// VerifyGoEnvironmentPaths verifies GOROOT and GOPATH settings
func VerifyGoEnvironmentPaths() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	
	homeDir := GetHomeDir()
	expectedGoroot := filepath.Join(homeDir, ".g", "go")
	expectedGopath := filepath.Join(homeDir, "go")

	// Verify GOROOT
	if output, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		goroot := strings.TrimSpace(string(output))
		if goroot == expectedGoroot {
			green.Printf("✅ GOROOT: %s\n", goroot)
		} else {
			yellow.Printf("⚠️  GOROOT: %s (expected: %s)\n", goroot, expectedGoroot)
		}
	}

	// Verify GOPATH
	if output, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		gopath := strings.TrimSpace(string(output))
		if gopath == expectedGopath {
			green.Printf("✅ GOPATH: %s\n", gopath)
		} else {
			yellow.Printf("⚠️  GOPATH: %s (expected: %s)\n", gopath, expectedGopath)
		}
	}
}

// GetExpectedGoRoot returns the expected GOROOT path
func GetExpectedGoRoot() string {
	return filepath.Join(GetHomeDir(), ".g", "go")
}

// GetExpectedGoPath returns the expected GOPATH path
func GetExpectedGoPath() string {
	return filepath.Join(GetHomeDir(), "go")
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

// AppendToFile appends content to a file
func AppendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = file.WriteString(content)
	return err
}

// WriteToFile writes content to a file
func WriteToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// HasConfigContent checks if a file contains specific configuration content
func HasConfigContent(filename, content string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), content) {
			return true
		}
	}
	return false
}

// DisplayCurrentGoVersion displays the current Go version with GOROOT and GOPATH
func DisplayCurrentGoVersion() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("📍 Current Go version:")

	goCmd := exec.Command("go", "version")
	if err := goCmd.Run(); err != nil {
		yellow.Println("⚠️  Go is not available in PATH")
		return
	}

	goCmd.Stdout = os.Stdout
	goCmd.Run()

	// Show GOROOT and GOPATH
	if gorootCmd := exec.Command("go", "env", "GOROOT"); gorootCmd.Run() == nil {
		blue.Print("📂 GOROOT: ")
		gorootCmd.Stdout = os.Stdout
		gorootCmd.Run()
	}

	if gopathCmd := exec.Command("go", "env", "GOPATH"); gopathCmd.Run() == nil {
		blue.Print("📂 GOPATH: ")
		gopathCmd.Stdout = os.Stdout
		gopathCmd.Run()
	}
}

// GetGoVersion returns the current Go version as a string
func GetGoVersion() string {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "Go not available"
	}
	return strings.TrimSpace(string(output))
}

// GetGoEnvVar returns a specific Go environment variable
func GetGoEnvVar(envVar string) string {
	cmd := exec.Command("go", "env", envVar)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
