package cli

import (
	"fmt"
	"time"

	"ava/hub/external/tuya"
	"ava/hub/external/wiz"

	"github.com/spf13/cobra"
)

var discoverTimeout time.Duration

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Find WiZ and Tuya devices on this network",
	Long: `Broadcasts on the LAN and prints whatever answers.

WiZ bulbs reply to a UDP probe on 38899. Tuya devices announce themselves on
6666 and 6667; only protocol 3.3 can be controlled, so anything newer is listed
with a warning rather than silently ignored.`,
	RunE: runDiscover,
}

func init() {
	discoverCmd.Flags().DurationVar(&discoverTimeout, "timeout", 6*time.Second, "how long to listen")
}

func runDiscover(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	fmt.Printf("listening for %s ...\n\n", discoverTimeout)

	lights, wizErr := wiz.Discover(ctx, discoverTimeout)
	plugs, tuyaErr := tuya.Discover(ctx, discoverTimeout)

	fmt.Printf("%-6s %-16s %-20s %-8s %s\n", "VENDOR", "IP", "ID", "VERSION", "NOTE")

	for at := range lights {
		fmt.Printf("%-6s %-16s %-20s %-8s on=%v %d%% %dK\n",
			"wiz", lights[at].Info.IP, lights[at].Info.MAC, "-",
			lights[at].State.Power, lights[at].State.Brightness, lights[at].State.ColorTemp)
	}

	for at := range plugs {
		note := "supported"
		if !plugs[at].Supported() {
			note = "UNSUPPORTED protocol - only 3.3 works"
		}

		fmt.Printf("%-6s %-16s %-20s %-8s %s\n",
			"tuya", plugs[at].Info.IP, plugs[at].Info.ID, plugs[at].Version, note)
	}

	if len(lights) == 0 && len(plugs) == 0 {
		fmt.Println("(nothing answered)")
	}

	if wizErr != nil {
		return wizErr
	}

	return tuyaErr
}
