package cmd

import (
	"fmt"

	"github.com/christian-doucette/time-to-go/internal/oled"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear [--i2c-bus YOUR_I2C_BUS]",
	Short: "Clears display for OLED monitor",
	Long:  "This command clears the display for an OLED monitor connected via I2C.",
	Run: func(cmd *cobra.Command, args []string) {
		i2cBus, _ := cmd.Flags().GetInt("i2c-bus")
		oled.ClearDisplay(fmt.Sprint(i2cBus))
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)

	clearCmd.Flags().Int("i2c-bus", 1, "Bus for the I2C connection")

}
