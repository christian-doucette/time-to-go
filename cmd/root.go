package cmd

import (
	"os"

	"github.com/christian-doucette/time-to-go/internal/gtfs"
	"github.com/christian-doucette/time-to-go/internal/mta"
	"github.com/christian-doucette/time-to-go/internal/oled"
	"github.com/spf13/cobra"
)

var (
	mtaApiKey, stopId string
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
		line := stopId[0]

		gtfsRaw := mta.CallRealtimeFeedApi(mtaApiKey, line)
		arrivalTimes := gtfs.ExtractStopArrivalTimes(gtfsRaw, stopId, 5)
		oled.DisplayTextLines(arrivalTimes, 13, 12)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Flags().StringVar(&mtaApiKey, "mta-api-key", "", "API key used for calls to the MTA API")
	rootCmd.MarkFlagRequired("mta-api-key")

	rootCmd.Flags().StringVar(&stopId, "stop-id", "", "Stop ID for the station")
	rootCmd.MarkFlagRequired("stop-id")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
