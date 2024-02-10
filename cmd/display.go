package cmd

import (
	"fmt"

	"github.com/christian-doucette/time-to-go/internal/gtfs"
	"github.com/christian-doucette/time-to-go/internal/mta"
	"github.com/spf13/cobra"
)

var (
	mtaApiKey, stopId string
)

// displayCmd represents the display command
var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		line := stopId[0]

		gtfsRaw := mta.CallRealtimeFeedApi(mtaApiKey, line)
		fmt.Println(gtfs.ExtractStopArrivalTimes(gtfsRaw, stopId, 5))
	},
}

func init() {
	displayCmd.Flags().StringVarP(&mtaApiKey, "mta-api-key", "k", "", "API key used for calls to the MTA API")
	displayCmd.Flags().StringVarP(&stopId, "stop-id", "s", "", "Stop ID for the station")

	rootCmd.AddCommand(displayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// displayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// displayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
