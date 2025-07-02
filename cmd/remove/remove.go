package remove

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// NewRemoveCmd creates the remove command
func NewRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [version]",
		Short:   "Remove a specific Go version",
		Long:    `Remove a specific Go version that has been installed via any available version manager.`,
		Example: `  gos remove 1.20.10    # Remove Go 1.20.10`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !common.CheckVersionManagerAvailable() {
				return
			}
			removeVersion(args[0])
		},
	}
}

// removeVersion removes a specific Go version
func removeVersion(version string) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	yellow.Printf("🗑️  Removing Go %s...\n", version)

	// Create progress bar for removal
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(fmt.Sprintf("Removing Go %s", version)),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "=", SaucerHead: ">", SaucerPadding: " ", BarStart: "[", BarEnd: "]",
		}),
	)

	// Start progress bar in goroutine
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				bar.Add(1)
				time.Sleep(150 * time.Millisecond)
			}
		}
	}()

	var cmd *exec.Cmd

	if common.IsCommandAvailable("gobrew") {
		cmd = exec.Command("gobrew", "uninstall", version)
	} else {
		done <- true
		bar.Finish()
		red.Println("❌ No version manager available")
		return
	}

	if err := cmd.Run(); err != nil {
		done <- true
		bar.Finish()
		red.Printf("❌ Error removing Go %s\n", version)
		return
	}

	done <- true
	bar.Finish()
	green.Printf("✅ Go %s removed successfully\n", version)
}
