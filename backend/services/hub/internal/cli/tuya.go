package cli

import (
	"fmt"
	"sort"
	"time"

	"ava/hub/external/tuya"
	"ava/hub/internal/device"

	"github.com/spf13/cobra"
)

const tuyaTimeout = 5 * time.Second

var (
	tuyaIP      string
	tuyaID      string
	tuyaKey     string
	tuyaPercent int
)

var tuyaCmd = &cobra.Command{
	Use:   "tuya",
	Short: "Control a Tuya device over the local protocol",
	Long: `Speaks Tuya protocol 3.3 directly to a device on the LAN.

Every subcommand needs the device id from 'avactl discover' and the local key,
which lives in Tuya's cloud and has to be fetched once per device.`,
}

func newTuyaDevice() (*tuya.Device, error) {
	return tuya.New(&tuya.Config{
		ID:           tuyaID,
		IP:           tuyaIP,
		LocalKey:     tuyaKey,
		Capabilities: device.CapabilityBrightness | device.CapabilityColorTemp,
		Timeout:      tuyaTimeout,
	})
}

var tuyaStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Print every data point the device reports",
	RunE: func(cmd *cobra.Command, _ []string) error {
		plug, err := newTuyaDevice()
		if err != nil {
			return err
		}

		dps, err := plug.RawState(cmd.Context())
		if err != nil {
			return err
		}

		keys := make([]string, 0, len(dps))
		for key := range dps {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		fmt.Println("raw data points:")

		for _, key := range keys {
			fmt.Printf("  %-5s %v\n", key, dps[key])
		}

		return nil
	},
}

var tuyaOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Switch the device on",
	RunE: func(cmd *cobra.Command, _ []string) error {
		plug, err := newTuyaDevice()
		if err != nil {
			return err
		}

		return plug.SetPower(cmd.Context(), true)
	},
}

var tuyaOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Switch the device off",
	RunE: func(cmd *cobra.Command, _ []string) error {
		plug, err := newTuyaDevice()
		if err != nil {
			return err
		}

		return plug.SetPower(cmd.Context(), false)
	},
}

var tuyaDimCmd = &cobra.Command{
	Use:   "dim",
	Short: "Set brightness on a device that supports it",
	RunE: func(cmd *cobra.Command, _ []string) error {
		plug, err := newTuyaDevice()
		if err != nil {
			return err
		}

		return plug.SetBrightness(cmd.Context(), tuyaPercent)
	},
}

func init() {
	tuyaCmd.PersistentFlags().StringVar(&tuyaIP, "ip", "", "device address on the LAN")
	tuyaCmd.PersistentFlags().StringVar(&tuyaID, "id", "", "device id from discovery")
	tuyaCmd.PersistentFlags().StringVar(&tuyaKey, "key", "", "16 character local key")
	tuyaDimCmd.Flags().IntVar(&tuyaPercent, "percent", 50, "brightness, 0-100")

	for _, required := range []string{"ip", "id", "key"} {
		_ = tuyaCmd.MarkPersistentFlagRequired(required)
	}

	tuyaCmd.AddCommand(tuyaStateCmd, tuyaOnCmd, tuyaOffCmd, tuyaDimCmd)
}
