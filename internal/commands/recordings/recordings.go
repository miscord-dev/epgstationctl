package recordings

import (
	"fmt"

	"github.com/miscord-dev/epgstationctl/internal/client"
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
	"github.com/miscord-dev/epgstationctl/internal/epgstation"
	"github.com/miscord-dev/epgstationctl/internal/output"
	"github.com/spf13/cobra"
)

var (
	offset    int
	limit     int
	halfWidth bool
)

var recordingsCmd = &cobra.Command{
	Use:   "recordings",
	Short: "Manage recordings",
	Long:  "Commands for managing EPGStation recordings",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current recordings",
	Long:  "List currently active recordings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetRecordingParams{
			IsHalfWidth: halfWidth,
		}

		if offset > 0 {
			offsetParam := epgstation.Offset(offset)
			params.Offset = &offsetParam
		}

		if limit > 0 {
			limitParam := epgstation.Limit(limit)
			params.Limit = &limitParam
		}

		recordings, err := client.GetRecordings(params)
		if err != nil {
			return fmt.Errorf("failed to get recordings: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(recordings.Records)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show recording status",
	Long:  "Show overall recording system status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetRecordingParams{
			IsHalfWidth: halfWidth,
		}

		recordings, err := client.GetRecordings(params)
		if err != nil {
			return fmt.Errorf("failed to get recordings: %w", err)
		}

		// Create status summary
		status := struct {
			TotalRecordings int `json:"total_recordings"`
			ActiveCount     int `json:"active_count"`
		}{
			TotalRecordings: recordings.Total,
			ActiveCount:     len(recordings.Records),
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(status)
	},
}

func init() {
	// Add flags to list command
	listCmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 24, "Number of recordings to retrieve")
	listCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	// Add flags to status command
	statusCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	recordingsCmd.AddCommand(listCmd)
	recordingsCmd.AddCommand(statusCmd)
	root.AddCommand(recordingsCmd)
}