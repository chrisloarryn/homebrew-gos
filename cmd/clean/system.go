package clean

import (
	"os/exec"
)

// CleanSystemGo removes manual system installations
func CleanSystemGo() {
	// Remove manual system installations
	exec.Command("sudo", "rm", "-rf", "/usr/local/go").Run()
}
