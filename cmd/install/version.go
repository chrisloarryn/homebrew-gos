package install

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// InstallVersion installs a specific Go version
func InstallVersion(version string) {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	blue.Printf("📦 Installing Go %s...\n", version)

	// Create progress bar for installation
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(fmt.Sprintf("Installing Go %s", version)),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowCount(),
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
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	cmd := exec.Command("gobrew", "install", version)

	if err := cmd.Run(); err != nil {
		done <- true
		bar.Finish()
		red.Printf("❌ Error installing Go %s\n", version)
		return
	}

	done <- true
	bar.Finish()
	green.Printf("✅ Go %s installed successfully\n", version)
}
