package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/christian-doucette/time-to-go/internal/gtfs"
	"github.com/christian-doucette/time-to-go/internal/mta"
	"github.com/christian-doucette/time-to-go/internal/oled"
	"github.com/spf13/cobra"
)

// subwayCmd represents the subway command
var subwayCmd = &cobra.Command{
	Use:   "subway --stop-id YOUR_STOP_ID [--bus YOUR_I2C_BUS] [--debug] [--least-minutes-ahead LEAST_MINUTES_AHEAD]",
	Short: "Display next subway times on OLED monitor via I2C",
	Long: `This command will pull the next subway arrival times for a specific stop from the MTA API. 	
If the debug option is included, it will print the arrival times to the terminal output.
If the debug option is not included, it will display the arrival times on an OLED display connected over I2C.`,
	Run: func(cmd *cobra.Command, args []string) {
		i2cBus, _ := cmd.Flags().GetInt("i2c-bus")
		stopId, _ := cmd.Flags().GetString("stop-id")
		debug, _ := cmd.Flags().GetBool("debug")
		leastMinutesAhead, _ := cmd.Flags().GetInt64("least-minutes-ahead")

		// gets the subway arrival data right now
		gtfsRaw := mta.CallAllSubwayRealtimeFeedApis()

		// parses out the next 4 arrival times for this stop
		arrivalTimes := gtfs.ExtractSubwayStopArrivalTimes(gtfsRaw, stopId, 4, leastMinutesAhead)

		if debug {
			// if debug option selected, prints arrival times to stdout
			fmt.Println(strings.Join(arrivalTimes, "\n"))

		} else {
			// otherwise displays arrival times to OLED display
			oled.DisplayTextLines(arrivalTimes, 13, 12, fmt.Sprint(i2cBus))
		}
	},
}

func init() {
	rootCmd.AddCommand(subwayCmd)

	subwayCmd.Flags().String("stop-id", "", "Stop ID for the subway station")
	err := subwayCmd.MarkFlagRequired("stop-id")
	if err != nil {
		os.Exit(1)
	}

	subwayCmd.Flags().Int("i2c-bus", 1, "Bus for the I2C connection")
	subwayCmd.Flags().Bool("debug", false, "Print output to terminal instead of OLED display")
	subwayCmd.Flags().Int64("least-minutes-ahead", 0, "Minimum number of minutes into the future you want to see subways for")
}
