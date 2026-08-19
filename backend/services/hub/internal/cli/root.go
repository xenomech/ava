package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "avactl",
	Short: "Talk to smart devices on this network",
	Long: `avactl (` + Version + `)

Discover and drive WiZ and Tuya devices directly over the LAN, without the
server or the hub. Useful for confirming a device answers before pairing a hub,
and for debugging one that has stopped responding.`,
	Version:       Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(wizCmd)
	rootCmd.AddCommand(tuyaCmd)
}
