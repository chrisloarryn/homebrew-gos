package version

import (
	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command to show Go version info
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show current Go version information",
		Long:  `Display detailed information about the currently active Go version, including GOROOT and GOPATH.`,
		Run: func(cmd *cobra.Command, args []string) {
			common.DisplayCurrentGoVersion()
		},
	}
}
