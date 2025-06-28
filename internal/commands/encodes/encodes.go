package encodes

import (
	"fmt"
	"strconv"

	"github.com/miscord-dev/epgstationctl/internal/client"
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
	"github.com/miscord-dev/epgstationctl/internal/epgstation"
	"github.com/miscord-dev/epgstationctl/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Common flags
	halfWidth bool

	// Add command flags
	recordedId        int64
	sourceVideoFileId int64
	mode              string
	parentDir         string
	directory         string
	removeOriginal    bool
	saveSameDirectory bool
)

var encodesCmd = &cobra.Command{
	Use:   "encodes",
	Short: "Manage encoding jobs",
	Long:  "Commands for managing EPGStation encoding jobs",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List encoding jobs",
	Long:  "List running and queued encoding jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetEncodeParams{
			IsHalfWidth: halfWidth,
		}

		encodes, err := client.GetEncodes(params)
		if err != nil {
			return fmt.Errorf("failed to get encodes: %w", err)
		}

		// Create a combined list of running and waiting items for display
		type EncodeDisplayItem struct {
			ID      int      `json:"id"`
			Mode    string   `json:"mode"`
			Status  string   `json:"status"`
			Percent *float32 `json:"percent,omitempty"`
			Name    string   `json:"name,omitempty"`
			Log     *string  `json:"log,omitempty"`
		}

		var displayItems []EncodeDisplayItem

		// Add running items
		for _, item := range encodes.RunningItems {
			displayItem := EncodeDisplayItem{
				ID:      item.Id,
				Mode:    item.Mode,
				Status:  "running",
				Percent: item.Percent,
				Log:     item.Log,
			}
			displayItem.Name = item.Recorded.Name
			displayItems = append(displayItems, displayItem)
		}

		// Add waiting items
		for _, item := range encodes.WaitItems {
			displayItem := EncodeDisplayItem{
				ID:     item.Id,
				Mode:   item.Mode,
				Status: "waiting",
				Log:    item.Log,
			}
			displayItem.Name = item.Recorded.Name
			displayItems = append(displayItems, displayItem)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(displayItems)
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new encoding job",
	Long:  "Add a new manual encoding job for recorded content",
	RunE: func(cmd *cobra.Command, args []string) error {
		if recordedId <= 0 {
			return fmt.Errorf("recorded-id is required and must be positive")
		}

		if mode == "" {
			return fmt.Errorf("mode is required")
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		// Build encode option
		encodeOption := epgstation.AddManualEncodeProgramOption{
			RecordedId: epgstation.RecordedId(recordedId),
			Mode:       mode,
		}

		if sourceVideoFileId > 0 {
			encodeOption.SourceVideoFileId = epgstation.VideoFileId(sourceVideoFileId)
		} else {
			// SourceVideoFileId is required, use 0 as default
			encodeOption.SourceVideoFileId = epgstation.VideoFileId(0)
		}

		if parentDir != "" {
			encodeOption.ParentDir = &parentDir
		}

		if directory != "" {
			encodeOption.Directory = &directory
		}

		encodeOption.RemoveOriginal = removeOriginal

		if saveSameDirectory {
			encodeOption.IsSaveSameDirectory = &saveSameDirectory
		}

		result, err := client.CreateEncode(encodeOption)
		if err != nil {
			return fmt.Errorf("failed to add encode: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(result)
	},
}

var cancelCmd = &cobra.Command{
	Use:   "cancel <encode-id>",
	Short: "Cancel an encoding job",
	Long:  "Cancel a running or queued encoding job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		encodeId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid encode ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.CancelEncode(encodeId)
		if err != nil {
			return fmt.Errorf("failed to cancel encode: %w", err)
		}

		fmt.Printf("Encode job %d cancelled successfully\n", encodeId)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show encoding system status",
	Long:  "Show encoding system status with running and queued job counts",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetEncodeParams{
			IsHalfWidth: halfWidth,
		}

		encodes, err := client.GetEncodes(params)
		if err != nil {
			return fmt.Errorf("failed to get encodes: %w", err)
		}

		// Create status summary
		status := struct {
			RunningCount int `json:"running_count"`
			WaitingCount int `json:"waiting_count"`
			TotalCount   int `json:"total_count"`
		}{
			RunningCount: 0,
			WaitingCount: 0,
		}

		status.RunningCount = len(encodes.RunningItems)
		status.WaitingCount = len(encodes.WaitItems)

		status.TotalCount = status.RunningCount + status.WaitingCount

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
	// List command flags
	listCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	// Add command flags
	addCmd.Flags().Int64Var(&recordedId, "recorded-id", 0, "Recorded program ID (required)")
	addCmd.Flags().Int64Var(&sourceVideoFileId, "source-file-id", 0, "Source video file ID")
	addCmd.Flags().StringVar(&mode, "mode", "", "Encode mode (required)")
	addCmd.Flags().StringVar(&parentDir, "parent-dir", "", "Parent directory for output")
	addCmd.Flags().StringVar(&directory, "directory", "", "Directory for output")
	addCmd.Flags().BoolVar(&removeOriginal, "remove-original", false, "Remove original file after encoding")
	addCmd.Flags().BoolVar(&saveSameDirectory, "save-same-dir", false, "Save in same directory as original")

	// Mark required flags
	addCmd.MarkFlagRequired("recorded-id")
	addCmd.MarkFlagRequired("mode")

	// Status command flags
	statusCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	// Add subcommands
	encodesCmd.AddCommand(listCmd)
	encodesCmd.AddCommand(addCmd)
	encodesCmd.AddCommand(cancelCmd)
	encodesCmd.AddCommand(statusCmd)

	// Register with root command
	root.AddCommand(encodesCmd)
}
