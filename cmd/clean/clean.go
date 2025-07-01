package clean

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// NewCleanCmd creates the clean command
func NewCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Deep clean all Go installations",
		Long: `Perform a comprehensive cleanup of all Go installations, including:
- Go cache and modules
- Homebrew installations
- Manual system installations
- User directories with special permissions
- Shell configuration cleanup`,
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")
			DeepCleanGo(force)
		},
	}
	
	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	return cmd
}

// DeepCleanGo performs comprehensive Go cleanup
func DeepCleanGo(force bool) {
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)

	red.Println("🗑️  Complete Go system cleanup...")

	if !force {
		yellow.Println("\n⚠️  WARNING: This will remove ALL Go installations and configurations!")
		fmt.Print("Are you sure you want to continue? (y/N): ")
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input != "y" && input != "yes" {
			yellow.Println("Cleanup cancelled.")
			return
		}
	}

	// Create main progress bar for cleanup stages
	totalStages := 6
	mainBar := progressbar.NewOptions(totalStages,
		progressbar.OptionSetDescription("🧹 Deep cleanup progress"),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "█", SaucerHead: "█", SaucerPadding: "░", BarStart: "[", BarEnd: "]",
		}),
		progressbar.OptionShowCount(),
	)

	fmt.Println()
	
	blue.Println("▸ Cleaning existing Go cache and modules…")
	CleanGoCache()
	mainBar.Add(1)

	blue.Println("\n▸ Removing Homebrew installations…")
	CleanHomebrewGo()
	mainBar.Add(1)

	blue.Println("\n▸ Removing manual system installations…")
	CleanSystemGo()
	mainBar.Add(1)

	blue.Println("\n▸ Removing user directories with special permissions…")
	CleanUserDirectories()
	mainBar.Add(1)

	blue.Println("\n▸ Removing other managers and directories…")
	CleanOtherManagers()
	mainBar.Add(1)

	blue.Println("\n▸ Cleaning shell configuration…")
	CleanShellConfig()
	mainBar.Add(1)

	// Clear command hash
	exec.Command("hash", "-r").Run()

	mainBar.Finish()
	fmt.Println()
	green.Println("✅ Complete Go cleanup finished.")
	yellow.Println("📋 Backups of your configuration files were created.")
	yellow.Println("🔄 Run 'source ~/.zshrc' or open a new terminal.")
}
