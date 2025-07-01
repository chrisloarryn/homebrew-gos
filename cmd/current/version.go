package current

import (
	"github.com/cristobalcontreras/gos/cmd/common"
	"github.com/spf13/cobra"
)

// NewCurrentCmd creates the current command to show Go version info
func NewCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show current Go version information",
		Long:  `Display detailed information about the currently active Go version, including GOROOT and GOPATH.`,
		Run: func(cmd *cobra.Command, args []string) {
			common.DisplayCurrentGoVersion()
		},
	}
}
