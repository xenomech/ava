package cli

import (
	"fmt"
	"strings"
	"time"

	"ava/hub/external/wiz"
	"ava/pkg/wire"

	"github.com/spf13/cobra"
)

var discoverTimeout time.Duration

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Find WiZ devices on this network",
	Long: `Broadcasts on the LAN and prints whatever answers.

WiZ bulbs reply to a UDP probe on 38899.`,
	RunE: runDiscover,
}

func init() {
	discoverCmd.Flags().DurationVar(&discoverTimeout, "timeout", 6*time.Second, "how long to listen")
}

func runDiscover(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	fmt.Printf("listening for %s ...\n\n", discoverTimeout)

	lights, err := wiz.Discover(ctx, discoverTimeout)

	fmt.Printf("%-6s %-16s %-20s %-8s %s\n", "VENDOR", "IP", "ID", "VERSION", "NOTE")

	for at := range lights {
		fmt.Printf("%-6s %-16s %-20s %-8s %s\n",
			"wiz", lights[at].Info.IP, lights[at].Info.MAC, "-", describe(lights[at].State))
	}

	if len(lights) == 0 {
		fmt.Println("(nothing answered)")
	}

	return err
}

func describe(state wire.State) string {
	parts := make([]string, 0, len(state))

	for _, trait := range []wire.Trait{wire.TraitPower, wire.TraitBrightness, wire.TraitColorTemp} {
		if value, ok := state.Get(trait); ok {
			parts = append(parts, fmt.Sprintf("%s=%s", trait, value))
		}
	}

	return strings.Join(parts, " ")
}
