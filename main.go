package main

import (
	"github.com/cristobalcontreras/gos/cmd"
)

// Build information. Populated at build-time via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version information in the root command
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
