package reserves

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/miscord-dev/epgstationctl/internal/client"
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
	"github.com/miscord-dev/epgstationctl/internal/epgstation"
	"github.com/miscord-dev/epgstationctl/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Common flags
	offset    int
	limit     int
	halfWidth bool

	// List flags
	reserveType string
	ruleId      int
	verbose     bool
	full        bool
	columns     string

	// Create/Update flags
	programId    int64
	encodeMode1  string
	encodeMode2  string
	encodeMode3  string
	parentDir    string
	directory    string
	allowEndLack bool
	tags         string

	// Status flags
	startAt string
	endAt   string
)

var reservesCmd = &cobra.Command{
	Use:   "reserves",
	Short: "Manage reservations",
	Long:  "Commands for managing EPGStation reservations",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List reservations",
	Long: `List all reservations with optional filtering.

Column Display Options:
  Default: Shows essential columns (Id, Name, StartAt, EndAt, ChannelId, IsConflict, IsSkip, IsOverlap)
  --verbose: Shows additional columns (includes RuleId, ProgramId, Description)
  --full: Shows all available columns
  --columns: Custom column selection (e.g., --columns=Id,Name,StartAt,ChannelId)

Available columns: Id, Name, StartAt, EndAt, ChannelId, RuleId, ProgramId, IsConflict, 
IsSkip, IsOverlap, Description, AllowEndLack, EncodeMode1, EncodeMode2, EncodeMode3, 
Directory, ParentDirectoryName, Tags, and more.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetReservesParams{
			IsHalfWidth: halfWidth,
		}

		if offset > 0 {
			params.Offset = &offset
		}

		if limit > 0 {
			params.Limit = &limit
		}

		if reserveType != "" {
			var typeValue epgstation.GetReserveType = reserveType
			params.Type = &typeValue
		}

		if ruleId > 0 {
			params.RuleId = &ruleId
		}

		reserves, err := client.GetReserves(params)
		if err != nil {
			return fmt.Errorf("failed to get reserves: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			// Define column sets
			defaultColumns := []string{"Id", "Name", "StartAt", "EndAt", "ChannelId", "IsConflict", "IsSkip", "IsOverlap"}
			verboseColumns := []string{"Id", "Name", "StartAt", "EndAt", "ChannelId", "RuleId", "ProgramId", "IsConflict", "IsSkip", "IsOverlap", "Description"}

			var selectedColumns []string

			if columns != "" {
				// Custom columns specified
				selectedColumns = strings.Split(strings.ReplaceAll(columns, " ", ""), ",")
			} else if full {
				// Full output - use default formatter (all columns)
				formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
				return formatter.Format(reserves.Reserves)
			} else if verbose {
				// Verbose output
				selectedColumns = verboseColumns
			} else {
				// Default minimal output
				selectedColumns = defaultColumns
			}

			formatter = output.NewTableFormatterWithColumns(nil, cfg.Output.NoHeader, selectedColumns)
		}

		return formatter.Format(reserves.Reserves)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <reserve-id>",
	Short: "Show detailed reservation information",
	Long:  "Show detailed information for a specific reservation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		reserve, err := client.GetReserve(reserveId)
		if err != nil {
			return fmt.Errorf("failed to get reserve: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(reserve)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new reservation",
	Long:  "Create a new manual reservation",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		// Build reserve option - note that ManualReserveOption doesn't have ProgramId
		// This will need to be handled differently or use timeSpecifiedOption
		reserveOption := epgstation.ManualReserveOption{}

		// Build encode option
		encodeOption := epgstation.ReserveEncodedOption{}

		if encodeMode1 != "" {
			encodeOption.Mode1 = &encodeMode1
		}
		if encodeMode2 != "" {
			encodeOption.Mode2 = &encodeMode2
		}
		if encodeMode3 != "" {
			encodeOption.Mode3 = &encodeMode3
		}

		// Build save option with directory settings
		saveOption := epgstation.ReserveSaveOption{}
		if parentDir != "" {
			saveOption.ParentDirectoryName = &parentDir
		}
		if directory != "" {
			saveOption.Directory = &directory
		}

		// ManualReserveOption is actually an alias for EditManualReserveOption
		reserveOption.AllowEndLack = allowEndLack
		reserveOption.SaveOption = &saveOption
		reserveOption.EncodeOption = &encodeOption

		result, err := client.CreateReserve(reserveOption)
		if err != nil {
			return fmt.Errorf("failed to create reserve: %w", err)
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

var updateCmd = &cobra.Command{
	Use:   "update <reserve-id>",
	Short: "Update an existing reservation",
	Long:  "Update settings for an existing reservation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		// Build encode option
		encodeOption := epgstation.ReserveEncodedOption{}

		if encodeMode1 != "" {
			encodeOption.Mode1 = &encodeMode1
		}
		if encodeMode2 != "" {
			encodeOption.Mode2 = &encodeMode2
		}
		if encodeMode3 != "" {
			encodeOption.Mode3 = &encodeMode3
		}

		// Build save option with directory settings
		saveOption := epgstation.ReserveSaveOption{}
		if parentDir != "" {
			saveOption.ParentDirectoryName = &parentDir
		}
		if directory != "" {
			saveOption.Directory = &directory
		}

		// Build update option
		updateOption := epgstation.EditManualReserveOption{
			AllowEndLack: allowEndLack,
			SaveOption:   &saveOption,
			EncodeOption: &encodeOption,
		}

		err = client.UpdateReserve(reserveId, updateOption)
		if err != nil {
			return fmt.Errorf("failed to update reserve: %w", err)
		}

		fmt.Printf("Reserve %d updated successfully\n", reserveId)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <reserve-id>",
	Short: "Delete a reservation",
	Long:  "Delete an existing reservation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.DeleteReserve(reserveId)
		if err != nil {
			return fmt.Errorf("failed to delete reserve: %w", err)
		}

		fmt.Printf("Reserve %d deleted successfully\n", reserveId)
		return nil
	},
}

var unskipCmd = &cobra.Command{
	Use:   "unskip <reserve-id>",
	Short: "Remove skip status from reservation",
	Long:  "Remove skip status from a reservation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.UnskipReserve(reserveId)
		if err != nil {
			return fmt.Errorf("failed to unskip reserve: %w", err)
		}

		fmt.Printf("Reserve %d unskipped successfully\n", reserveId)
		return nil
	},
}

var removeOverlapCmd = &cobra.Command{
	Use:   "remove-overlap <reserve-id>",
	Short: "Remove overlap status from reservation",
	Long:  "Remove overlap status from a reservation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.RemoveOverlapReserve(reserveId)
		if err != nil {
			return fmt.Errorf("failed to remove overlap from reserve: %w", err)
		}

		fmt.Printf("Overlap removed from reserve %d successfully\n", reserveId)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show reservation status summary",
	Long:  "Show reservation status summary with counts by category",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		counts, err := client.GetReserveCounts()
		if err != nil {
			return fmt.Errorf("failed to get reserve counts: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(counts)
	},
}

var updateSystemCmd = &cobra.Command{
	Use:   "update-system",
	Short: "Trigger reservation system update",
	Long:  "Trigger a reservation system update",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.UpdateReserveSystem()
		if err != nil {
			return fmt.Errorf("failed to update reserve system: %w", err)
		}

		fmt.Println("Reserve system update triggered successfully")
		return nil
	},
}

func init() {
	// List command flags
	listCmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Number of reserves to retrieve")
	listCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")
	listCmd.Flags().StringVar(&reserveType, "type", "", "Filter by reserve type (all, normal, conflict, skip, overlap)")
	listCmd.Flags().IntVar(&ruleId, "rule-id", 0, "Filter by rule ID")
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show more columns (includes RuleId, ProgramId, Description)")
	listCmd.Flags().BoolVar(&full, "full", false, "Show all available columns")
	listCmd.Flags().StringVar(&columns, "columns", "", "Comma-separated list of columns to display (e.g., 'Id,Name,StartAt,EndAt')")

	// Show command flags
	showCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	// Create command flags
	createCmd.Flags().Int64Var(&programId, "program-id", 0, "Program ID to reserve")
	createCmd.Flags().StringVar(&encodeMode1, "encode-mode1", "", "Primary encode mode")
	createCmd.Flags().StringVar(&encodeMode2, "encode-mode2", "", "Secondary encode mode")
	createCmd.Flags().StringVar(&encodeMode3, "encode-mode3", "", "Tertiary encode mode")
	createCmd.Flags().StringVar(&parentDir, "parent-dir", "", "Parent directory for saving")
	createCmd.Flags().StringVar(&directory, "directory", "", "Directory for saving")
	createCmd.Flags().BoolVar(&allowEndLack, "allow-end-lack", false, "Allow recording to end early")
	createCmd.Flags().StringVar(&tags, "tags", "", "Tags for the reservation")

	// Update command flags
	updateCmd.Flags().StringVar(&encodeMode1, "encode-mode1", "", "Primary encode mode")
	updateCmd.Flags().StringVar(&encodeMode2, "encode-mode2", "", "Secondary encode mode")
	updateCmd.Flags().StringVar(&encodeMode3, "encode-mode3", "", "Tertiary encode mode")
	updateCmd.Flags().StringVar(&parentDir, "parent-dir", "", "Parent directory for saving")
	updateCmd.Flags().StringVar(&directory, "directory", "", "Directory for saving")
	updateCmd.Flags().BoolVar(&allowEndLack, "allow-end-lack", false, "Allow recording to end early")
	updateCmd.Flags().StringVar(&tags, "tags", "", "Tags for the reservation")

	// Status command flags
	statusCmd.Flags().StringVar(&startAt, "start-at", "", "Start date for filtering (YYYY-MM-DD format)")
	statusCmd.Flags().StringVar(&endAt, "end-at", "", "End date for filtering (YYYY-MM-DD format)")

	// Add subcommands
	reservesCmd.AddCommand(listCmd)
	reservesCmd.AddCommand(showCmd)
	reservesCmd.AddCommand(createCmd)
	reservesCmd.AddCommand(updateCmd)
	reservesCmd.AddCommand(deleteCmd)
	reservesCmd.AddCommand(unskipCmd)
	reservesCmd.AddCommand(removeOverlapCmd)
	reservesCmd.AddCommand(statusCmd)
	reservesCmd.AddCommand(updateSystemCmd)

	// Register with root command
	root.AddCommand(reservesCmd)
}
