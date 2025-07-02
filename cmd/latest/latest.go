package latest

import (
	"os"
	"os/exec"
	"time"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// NewLatestCmd creates the latest command
func NewLatestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "latest",
		Short: "Install and use the latest Go version",
		Long:  `Install the latest stable Go version and automatically switch to it.`,
		Run: func(cmd *cobra.Command, args []string) {
			if !common.CheckVersionManagerAvailable() {
				return
			}
			installLatest()
		},
	}
}

// installLatest installs and switches to the latest Go version
func installLatest() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	blue.Println("🚀 Installing latest Go version...")

	// Install latest version
	if !executeInstallLatest() {
		return
	}

	green.Println("✅ Latest version installed")

	// Switch to latest version
	blue.Println("� Switching to latest version...")
	switchToLatest()

	// Show current version
	showCurrentVersion(blue)
}

// executeInstallLatest performs the installation with progress bar
func executeInstallLatest() bool {
	red := color.New(color.FgRed)

	// Create progress bar for installation
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Installing latest Go"),
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
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	installCmd := getInstallCommand()
	if installCmd == nil {
		done <- true
		bar.Finish()
		red.Println("❌ No version manager available")
		return false
	}

	if err := installCmd.Run(); err != nil {
		done <- true
		bar.Finish()
		red.Println("❌ Error installing latest version")
		return false
	}

	done <- true
	bar.Finish()
	return true
}

// getInstallCommand returns the appropriate install command
func getInstallCommand() *exec.Cmd {
	if common.IsCommandAvailable("gobrew") {
		return exec.Command("gobrew", "install", "latest")
	}

	return nil
}

// switchToLatest switches to the latest installed version
func switchToLatest() {
	if common.IsCommandAvailable("gobrew") {
		useCmd := exec.Command("gobrew", "use", "latest")
		useCmd.Run()
	}
}

// showCurrentVersion displays the current Go version
func showCurrentVersion(blue *color.Color) {
	if goCmd := exec.Command("go", "version"); goCmd.Run() == nil {
		blue.Print("📋 Current version: ")
		goCmd.Stdout = os.Stdout
		goCmd.Run()
	}
}
