package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/cristobalcontreras/gos/cmd/common"
)

// ValidateEnvironment runs comprehensive environment validation
func ValidateEnvironment() {
	blue := color.New(color.FgBlue)
	blue.Println("🔍 Comprehensive Environment Validation")
	fmt.Println("")

	config := getValidationConfig()
	validationResult := &ValidationResult{}

	// Run all validation checks
	validateEnvironmentVariables(config, validationResult)
	validatePathConfiguration(config, validationResult)
	validateDirectoryStructure(config, validationResult)
	validateShellConfiguration(config, validationResult)
	hasVersionManager := validateVersionManager(validationResult)
	validateGoBinary(hasVersionManager, validationResult)

	// Display summary
	displayValidationSummary(validationResult)
}

// ValidationResult holds the results of environment validation
type ValidationResult struct {
	HasErrors   bool
	HasWarnings bool
}

// ValidationConfig extends EnvironmentConfig with validation-specific fields
type ValidationConfig struct {
	EnvironmentConfig
	HomeDir      string
	CurrentShell string
	ShellFile    string
}

// getValidationConfig returns the validation configuration
func getValidationConfig() *ValidationConfig {
	homeDir := common.GetHomeDir()
	envConfig := getEnvironmentConfig()
	
	return &ValidationConfig{
		EnvironmentConfig: envConfig,
		HomeDir:          homeDir,
		CurrentShell:     common.DetectCurrentShell(),
		ShellFile:        common.GetShellFileForCurrentShell(common.DetectCurrentShell(), homeDir),
	}
}

// validateEnvironmentVariables validates GOROOT and GOPATH
func validateEnvironmentVariables(config *ValidationConfig, result *ValidationResult) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	blue.Println("📋 Environment Variables:")

	// GOROOT validation
	actualGoroot := os.Getenv("GOROOT")
	if actualGoroot == config.ExpectedGoroot {
		green.Println("  ✅ GOROOT is correctly set")
	} else if actualGoroot == "" {
		red.Println("  ❌ GOROOT is not set")
		result.HasErrors = true
	} else {
		yellow.Printf("  ⚠️  GOROOT is set to %s (expected %s)\n", actualGoroot, config.ExpectedGoroot)
		result.HasWarnings = true
	}

	// GOPATH validation
	actualGopath := os.Getenv("GOPATH")
	if actualGopath == config.ExpectedGopath {
		green.Println("  ✅ GOPATH is correctly set")
	} else if actualGopath == "" {
		red.Println("  ❌ GOPATH is not set")
		result.HasErrors = true
	} else {
		yellow.Printf("  ⚠️  GOPATH is set to %s (expected %s)\n", actualGopath, config.ExpectedGopath)
		result.HasWarnings = true
	}
}

// validatePathConfiguration validates PATH environment variable
func validatePathConfiguration(config *ValidationConfig, result *ValidationResult) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Println("")
	blue.Println("🛤️  PATH Validation:")
	
	path := os.Getenv("PATH")
	pathMissing := 0

	for _, reqPath := range config.RequiredPaths {
		if strings.Contains(path, reqPath) {
			green.Printf("  ✅ %s is in PATH\n", reqPath)
		} else {
			yellow.Printf("  ⚠️  %s is missing from PATH\n", reqPath)
			pathMissing++
		}
	}

	if pathMissing > 0 {
		fmt.Printf("    💡 Run 'gos setup' or 'gos env --fix' to add missing PATH entries\n")
		result.HasWarnings = true
	}
}

// validateDirectoryStructure validates required directories
func validateDirectoryStructure(config *ValidationConfig, result *ValidationResult) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	fmt.Println("")
	blue.Println("📁 Directory Structure:")

	dirs := getRequiredDirectoriesForValidation(config)

	for name, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			green.Printf("  ✅ %s exists: %s\n", name, dir)
		} else {
			if strings.Contains(name, "GOPATH") {
				yellow.Printf("  ⚠️  %s missing: %s\n", name, dir)
				result.HasWarnings = true
			} else {
				red.Printf("  ❌ %s missing: %s\n", name, dir)
				result.HasErrors = true
			}
		}
	}
}

// getRequiredDirectoriesForValidation returns the map of required directories
func getRequiredDirectoriesForValidation(config *ValidationConfig) map[string]string {
	dirs := map[string]string{
		"GOPATH":     config.ExpectedGopath,
		"GOPATH bin": filepath.Join(config.ExpectedGopath, "bin"),
	}

	// Add version manager directories if detected
	if runtime.GOOS == "windows" && common.IsCommandAvailable("gobrew") {
		dirs["gobrew directory"] = filepath.Join(config.HomeDir, ".gobrew")
		dirs["gobrew bin"] = filepath.Join(config.HomeDir, ".gobrew", "bin")
		dirs["Go installation"] = config.ExpectedGoroot
	} else if common.IsGInstalled() {
		dirs["g directory"] = filepath.Join(config.HomeDir, ".g")
		dirs["g bin directory"] = filepath.Join(config.HomeDir, ".g", "bin")
		dirs["Go installation"] = config.ExpectedGoroot
	}

	return dirs
}

// validateShellConfiguration validates shell configuration files
func validateShellConfiguration(config *ValidationConfig, result *ValidationResult) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Println("")
	blue.Println("🐚 Shell Configuration:")
	blue.Printf("  🔍 Detected shell: %s\n", config.CurrentShell)

	if config.ShellFile == "" {
		yellow.Println("  ⚠️  Could not determine shell configuration file")
		result.HasWarnings = true
		return
	}

	fullPath := filepath.Join(config.HomeDir, config.ShellFile)
	if _, err := os.Stat(fullPath); err == nil {
		if hasGoConfig(fullPath) {
			green.Printf("  ✅ Go configuration found in %s\n", config.ShellFile)
		} else {
			yellow.Printf("  ⚠️  %s exists but no Go configuration found\n", config.ShellFile)
			fmt.Printf("    💡 Run 'gos setup' or 'gos env --fix' to add configuration\n")
			result.HasWarnings = true
		}
	} else {
		yellow.Printf("  ⚠️  Shell file %s does not exist\n", config.ShellFile)
		fmt.Printf("    💡 Run 'gos setup' to create configuration\n")
		result.HasWarnings = true
	}
}

// validateVersionManager validates version manager availability
func validateVersionManager(result *ValidationResult) bool {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	fmt.Println("")
	blue.Println("🔧 Version Manager:")

	if runtime.GOOS == "windows" && common.IsCommandAvailable("gobrew") {
		return validateGobrewManager(result)
	} else if common.IsGInstalled() {
		return validateGManager(result)
	}

	yellow.Println("  ⚠️  No version manager found (gobrew or g)")
	fmt.Println("    💡 Run 'gos setup' to install a version manager")
	result.HasWarnings = true
	return false
}

// validateGobrewManager validates gobrew version manager
func validateGobrewManager(result *ValidationResult) bool {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	green.Println("  ✅ 'gobrew' version manager is available")
	
	if versions := common.GetGobrewVersions(); len(versions) > 0 {
		green.Printf("  ✅ %d Go version(s) installed\n", len(versions))
	} else {
		yellow.Println("  ⚠️  No Go versions installed with gobrew")
		result.HasWarnings = true
	}
	
	return true
}

// validateGManager validates g version manager
func validateGManager(result *ValidationResult) bool {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	green.Println("  ✅ 'g' version manager is available")
	
	if versions := common.GetInstalledVersions(); len(versions) > 0 {
		green.Printf("  ✅ %d Go version(s) installed\n", len(versions))
	} else {
		yellow.Println("  ⚠️  No Go versions installed with g")
		result.HasWarnings = true
	}
	
	return true
}

// validateGoBinary validates Go binary availability
func validateGoBinary(hasVersionManager bool, result *ValidationResult) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Println("")
	blue.Println("🐹 Go Binary:")

	if !hasVersionManager {
		yellow.Println("  ℹ️  Skipping Go binary check (no version manager)")
		return
	}

	goPath, err := exec.LookPath("go")
	if err != nil {
		yellow.Println("  ⚠️  Go binary not found in PATH")
		fmt.Println("    💡 Install a Go version with 'gos install latest'")
		result.HasWarnings = true
		return
	}

	green.Printf("  ✅ Go binary found: %s\n", goPath)

	// Check if go version works
	if output, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(output))
		green.Printf("  ✅ Go version: %s\n", version)
	} else {
		yellow.Println("  ⚠️  Go binary exists but 'go version' failed")
		result.HasWarnings = true
	}
}

// displayValidationSummary displays the final validation summary
func displayValidationSummary(result *ValidationResult) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	fmt.Println("")
	blue.Println("📊 Validation Summary:")

	if result.HasErrors {
		red.Println("  ❌ Environment has critical issues that need fixing")
		fmt.Println("  💡 Run 'gos env --fix' to attempt automatic fixes")
	} else if result.HasWarnings {
		yellow.Println("  ⚠️  Environment has minor issues")
		fmt.Println("  💡 Consider running 'gos env --fix' to optimize configuration")
	} else {
		green.Println("  ✅ Environment is properly configured!")
	}
}
