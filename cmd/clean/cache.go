package clean

import (
	"fmt"
	"os/exec"
)

// CleanGoCache cleans Go cache and modules
func CleanGoCache() {
	// Try to clean with go command if available
	if _, err := exec.LookPath("go"); err == nil {
		fmt.Println("  Running go clean -modcache...")
		exec.Command("go", "clean", "-modcache").Run()
		
		fmt.Println("  Running go clean -cache...")
		exec.Command("go", "clean", "-cache").Run()
	}
}
