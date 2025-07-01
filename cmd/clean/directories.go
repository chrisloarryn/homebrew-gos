package clean

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cristobalcontreras/gos/cmd/common"
)

// CleanUserDirectories removes user Go directories
func CleanUserDirectories() {
	homeDir := common.GetHomeDir()
	
	// Clean ~/go directory
	goDir := filepath.Join(homeDir, "go")
	if _, err := os.Stat(goDir); err == nil {
		fmt.Printf("  Fixing permissions in %s...\n", goDir)
		FixPermissions(goDir)
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
			FixPermissions(cacheDir)
			os.RemoveAll(cacheDir)
		}
	}
}

// CleanOtherManagers removes other Go version managers
func CleanOtherManagers() {
	homeDir := common.GetHomeDir()
	
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

// FixPermissions recursively fixes permissions for directories and files
func FixPermissions(dir string) {
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
