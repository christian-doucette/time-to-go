package cmd

import (
	"fmt"

	"github.com/christian-doucette/time-to-go/internal/oled"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear [--bus YOUR_I2C_BUS]",
	Short: "Clears display for OLED monitor",
	Long:  "This command clears the display for an OLED monitor connected via I2C.",
	Run: func(cmd *cobra.Command, args []string) {
		bus, _ := cmd.Flags().GetInt("bus")
		oled.ClearDisplay(fmt.Sprint(bus))
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
