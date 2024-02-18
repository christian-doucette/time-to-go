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
	Use:   "main",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		bus, _ := cmd.Flags().GetInt("bus")
		line := stopId[0]

		gtfsRaw := mta.CallRealtimeFeedApi(mtaApiKey, line)
		arrivalTimes := gtfs.ExtractStopArrivalTimes(gtfsRaw, stopId, 5)

		if debug {
			fmt.Println(strings.Join(arrivalTimes, "\n"))

		} else {
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
