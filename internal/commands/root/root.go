package root

import (
	"fmt"
	"os"
	"strconv"

	"github.com/miscord-dev/epgstationctl/internal/config"
	"github.com/spf13/cobra"
)

var (
	serverURL    string
	timeout      int
	outputFlag   string
	noHeaderFlag bool
	verboseFlag  bool
	cfg          *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "epgstationctl",
	Short: "EPGStation CLI client",
	Long: `epgstationctl is a command-line interface for EPGStation.
It allows you to manage channels, programs, recordings, and reservations
from the command line.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Apply environment variables if flags are not explicitly set
		applyEnvDefaults(cmd)
		cfg = config.NewConfig(serverURL, timeout, outputFlag, noHeaderFlag, verboseFlag)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8888", "EPGStation server URL (env: EPGSTATIONCTL_SERVER)")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "Request timeout in seconds (env: EPGSTATIONCTL_TIMEOUT)")
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "table", "output format (table, json) (env: EPGSTATIONCTL_OUTPUT)")
	rootCmd.PersistentFlags().BoolVar(&noHeaderFlag, "no-header", false, "hide table headers (env: EPGSTATIONCTL_NO_HEADER)")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output (env: EPGSTATIONCTL_VERBOSE)")
}

func applyEnvDefaults(cmd *cobra.Command) {
	// Apply environment variables only if flags were not explicitly set
	if !cmd.PersistentFlags().Changed("server") {
		if envServer := os.Getenv("EPGSTATIONCTL_SERVER"); envServer != "" {
			serverURL = envServer
		}
	}

	if !cmd.PersistentFlags().Changed("timeout") {
		if envTimeout := os.Getenv("EPGSTATIONCTL_TIMEOUT"); envTimeout != "" {
			if timeoutVal, err := strconv.Atoi(envTimeout); err == nil {
				timeout = timeoutVal
			}
		}
	}

	if !cmd.PersistentFlags().Changed("output") {
		if envOutput := os.Getenv("EPGSTATIONCTL_OUTPUT"); envOutput != "" {
			outputFlag = envOutput
		}
	}

	if !cmd.PersistentFlags().Changed("no-header") {
		if envNoHeader := os.Getenv("EPGSTATIONCTL_NO_HEADER"); envNoHeader != "" {
			if noHeaderVal, err := strconv.ParseBool(envNoHeader); err == nil {
				noHeaderFlag = noHeaderVal
			}
		}
	}

	if !cmd.PersistentFlags().Changed("verbose") {
		if envVerbose := os.Getenv("EPGSTATIONCTL_VERBOSE"); envVerbose != "" {
			if verboseVal, err := strconv.ParseBool(envVerbose); err == nil {
				verboseFlag = verboseVal
			}
		}
	}
}

func GetConfig() *config.Config {
	return cfg
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
