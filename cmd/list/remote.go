package list

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// ListRemoteVersions lists available remote Go versions
func ListRemoteVersions() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Println("🌐 Available versions:")

	// Try different version managers based on availability
	if common.IsCommandAvailable("gobrew") {
		listRemoteVersionsWithGobrew()
	} else if common.IsCommandAvailable("g") {
		listRemoteVersionsWithG()
	} else {
		yellow.Println("No version manager detected.")
		fmt.Println("")
		yellow.Println("💡 To install a version manager and browse remote versions:")
		fmt.Println("   gos setup               # Install version manager")
		fmt.Println("")
		yellow.Println("💡 You can also check manually at:")
		fmt.Println("   https://golang.org/dl/")
	}
}

// listRemoteVersionsWithGobrew lists remote versions using gobrew
func listRemoteVersionsWithGobrew() {
	yellow := color.New(color.FgYellow)

	// Create progress bar for fetching remote versions
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Fetching remote versions with gobrew"),
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
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	cmd := exec.Command("gobrew", "ls-remote")
	output, err := cmd.Output()
	if err != nil {
		done <- true
		bar.Finish()
		yellow.Println("  Could not get remote versions via gobrew")
		return
	}

	done <- true
	bar.Finish()

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Println("  Available versions from gobrew:")
	
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Available") {
			fmt.Printf("     %s\n", line)
			count++
			if count > 20 { // Limit output to first 20 versions
				fmt.Printf("     ... and %d more versions\n", len(lines)-count-1)
				fmt.Println("")
				fmt.Println("💡 Run 'gobrew ls-remote' to see all available versions")
				break
			}
		}
	}
}

// listRemoteVersionsWithG lists remote versions using g
func listRemoteVersionsWithG() {
	yellow := color.New(color.FgYellow)

	// Create progress bar for fetching remote versions
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Fetching remote versions with g"),
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
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	cmd := exec.Command("g", "list-all")
	output, err := cmd.Output()
	if err != nil {
		done <- true
		bar.Finish()
		yellow.Println("  Could not get remote versions via g")
		return
	}

	done <- true
	bar.Finish()

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Println("  Available versions from g:")
	
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Available") {
			fmt.Printf("     %s\n", line)
			count++
			if count > 20 { // Limit output to first 20 versions
				fmt.Printf("     ... and %d more versions\n", len(lines)-count-1)
				fmt.Println("")
				fmt.Println("💡 Run 'g list-all' to see all available versions")
				break
			}
		}
	}
}
