package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ava/hub/external/wiz"
	"ava/pkg/wire"

	"github.com/spf13/cobra"
)

const wizTimeout = 3 * time.Second

var wizCmd = &cobra.Command{
	Use:   "wiz",
	Short: "Control a WiZ bulb directly over UDP",
}

var wizStateCmd = &cobra.Command{
	Use:   "state <ip>",
	Short: "Read a bulb's identity and current state",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		light := wiz.New(args[0], wizTimeout)

		info, err := light.Identify(cmd.Context())
		if err != nil {
			return err
		}

		state, err := light.State(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Printf("mac        %s\nmodel      %s\nstate      %s\n", info.MAC, info.Model, describe(state))

		for _, capability := range light.Capabilities() {
			fmt.Printf("trait      %s\n", summarise(&capability))
		}

		return nil
	},
}

var wizOnCmd = &cobra.Command{
	Use:   "on <ip>",
	Short: "Switch a bulb on",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return set(cmd.Context(), args[0], wire.TraitPower, wire.Bool(true))
	},
}

var wizOffCmd = &cobra.Command{
	Use:   "off <ip>",
	Short: "Switch a bulb off",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return set(cmd.Context(), args[0], wire.TraitPower, wire.Bool(false))
	},
}

var wizDimCmd = &cobra.Command{
	Use:   "dim <ip> <10-100>",
	Short: "Set brightness, clamped to what the bulb accepts",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		percent, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("brightness must be a number: %w", err)
		}

		return set(cmd.Context(), args[0], wire.TraitBrightness, wire.Number(float64(percent)))
	},
}

var wizTempCmd = &cobra.Command{
	Use:   "temp <ip> <kelvin>",
	Short: "Set white temperature, clamped to 2200-6500K",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		kelvin, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("temperature must be a number: %w", err)
		}

		return set(cmd.Context(), args[0], wire.TraitColorTemp, wire.Number(float64(kelvin)))
	},
}

func init() {
	wizCmd.AddCommand(wizStateCmd, wizOnCmd, wizOffCmd, wizDimCmd, wizTempCmd)
}

func set(ctx context.Context, ip string, trait wire.Trait, value wire.Value) error {
	light := wiz.New(ip, wizTimeout)

	if _, err := light.Identify(ctx); err != nil {
		return err
	}

	return light.Apply(ctx, trait, value)
}

func summarise(capability *wire.Capability) string {
	out := fmt.Sprintf("%-16s %-6s %s", capability.Trait, capability.Kind, capability.Access)

	if capability.Min != nil && capability.Max != nil {
		out += fmt.Sprintf(" %g-%g%s", *capability.Min, *capability.Max, capability.Unit)
	}

	if len(capability.Values) > 0 {
		out += " " + strings.Join(capability.Values, "|")
	}

	return out
}
