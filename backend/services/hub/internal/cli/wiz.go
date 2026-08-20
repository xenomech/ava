package cli

import (
	"fmt"
	"strconv"
	"time"

	"ava/hub/external/wiz"

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

		fmt.Printf("mac        %s\nmodel      %s\npower      %v\nbrightness %d%%\ntemp       %dK\ncaps       %s\n",
			info.MAC, info.Model, state.Power, state.Brightness, state.ColorTemp, light.Capabilities())

		return nil
	},
}

var wizOnCmd = &cobra.Command{
	Use:   "on <ip>",
	Short: "Switch a bulb on",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return wiz.New(args[0], wizTimeout).SetPower(cmd.Context(), true)
	},
}

var wizOffCmd = &cobra.Command{
	Use:   "off <ip>",
	Short: "Switch a bulb off",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return wiz.New(args[0], wizTimeout).SetPower(cmd.Context(), false)
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

		return wiz.New(args[0], wizTimeout).SetBrightness(cmd.Context(), percent)
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

		return wiz.New(args[0], wizTimeout).SetColorTemp(cmd.Context(), kelvin)
	},
}

func init() {
	wizCmd.AddCommand(wizStateCmd, wizOnCmd, wizOffCmd, wizDimCmd, wizTempCmd)
}
