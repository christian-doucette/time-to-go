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

var (
	mtaApiKey, stopId string
	debug             bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "./time-to-go --mta-api-key YOUR_MTA_API_KEY --stop-id YOUR_STOP_ID [--bus YOUR_I2C_BUS] [--debug]",
	Short: "Display next subway times on OLED monitor via I2C",
	Long: `This command will pull the next subway arrival times for a specific stop from the MTA API. 	
If the debug option is included, it will print the arrival times to the terminal output.
If the debug option is not included, it will display the arrival times on an OLED display connected over I2C.`,
	Run: func(cmd *cobra.Command, args []string) {
		bus, _ := cmd.Flags().GetInt("bus")

		// gets the subway arrival data right now
		gtfsRaw := mta.CallAllSubwayRealtimeFeedApis(mtaApiKey)

		// parses out the next 5 arrival times for this stop
		arrivalTimes := gtfs.ExtractStopArrivalTimes(gtfsRaw, stopId, 4)

		if debug {
			// if debug option selected, prints arrival times to stdout
			fmt.Println(strings.Join(arrivalTimes, "\n"))

		} else {
			// otherwise displays arrival times to OLED display
			oled.DisplayTextLines(arrivalTimes, 13, 12, fmt.Sprint(bus))
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Flags().StringVar(&mtaApiKey, "mta-api-key", "", "API key used for calls to the MTA API")
	err := rootCmd.MarkFlagRequired("mta-api-key")
	if err != nil {
		os.Exit(1)
	}

	rootCmd.Flags().StringVar(&stopId, "stop-id", "", "Stop ID for the station")
	err = rootCmd.MarkFlagRequired("stop-id")
	if err != nil {
		os.Exit(1)
	}

	rootCmd.PersistentFlags().Int("bus", 1, "Bus for the I2C connection")

	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Print output to terminal instead of OLED display")

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
