package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	mtaApiKey string
	stopId    string
	debug     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "./time-to-go",
	Short: "Display next subway times on OLED monitor via I2C",
	Long: `This command will pull the next subway arrival times for a specific stop from the MTA API. 	
If the debug option is included, it will print the arrival times to the terminal output.
If the debug option is not included, it will display the arrival times on an OLED display connected over I2C.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hi there")
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
