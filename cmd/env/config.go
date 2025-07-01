package env

import (
	"os"
	"strings"
)

// hasGoConfig checks if a shell file contains Go configuration
func hasGoConfig(filename string) bool {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false
	}
	
	contentStr := string(content)
	return strings.Contains(contentStr, "GOROOT") || 
		   strings.Contains(contentStr, "GOPATH") ||
		   strings.Contains(contentStr, "Go Version Manager")
}
