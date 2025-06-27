package reserves

import (
	"fmt"
	"strconv"

	"github.com/miscord-dev/epgstationctl/internal/client"
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
	"github.com/miscord-dev/epgstationctl/internal/epgstation"
	"github.com/miscord-dev/epgstationctl/internal/output"
	"github.com/spf13/cobra"
)

const (
	outputFormatJSON = "json"
)

var (
	offset  int
	limit   int
	channel int
	status  string
	keyword string
	
	// Create reserve flags
	programID            int
	encodeMode           int
	encodeParentDir      string
	encodeDir            string
	allowEndLack         bool
)

var reservesCmd = &cobra.Command{
	Use:   "reserves",
	Short: "Manage recording reserves",
	Long:  "Commands for managing EPGStation recording reserves",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recording reserves",
	Long:  "List all recording reserves with optional filtering",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetReservesParams{}
		if offset > 0 {
			params.Offset = &offset
		}
		if limit > 0 {
			params.Limit = &limit
		}

		reserves, err := client.GetReserves(params)
		if err != nil {
			return fmt.Errorf("failed to get reserves: %w", err)
		}

		if cfg.Output.Format == outputFormatJSON {
			return output.PrintAsJSON(reserves)
		}

		return printReservesTable(reserves)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a manual reserve",
	Long:  "Create a manual recording reserve for a specific program",
	RunE: func(cmd *cobra.Command, args []string) error {
		if programID == 0 {
			return fmt.Errorf("program ID is required (use --program-id)")
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		reserveOption := epgstation.ManualReserveOption{}
		
		if allowEndLack {
			reserveOption.AllowEndLack = allowEndLack
		}

		// Add encode settings if specified
		if encodeMode > 0 || encodeParentDir != "" || encodeDir != "" {
			encodeOpt := &epgstation.ReserveEncodedOption{}
			if encodeMode > 0 {
				modeStr := strconv.Itoa(encodeMode)
				encodeOpt.Mode1 = &modeStr
			}
			if encodeParentDir != "" {
				encodeOpt.EncodeParentDirectoryName1 = &encodeParentDir
			}
			if encodeDir != "" {
				encodeOpt.Directory1 = &encodeDir
			}
			reserveOption.EncodeOption = encodeOpt
		}

		result, err := client.CreateReserve(reserveOption)
		if err != nil {
			return fmt.Errorf("failed to create reserve: %w", err)
		}

		if cfg.Output.Format == outputFormatJSON {
			return output.PrintAsJSON(result)
		}

		fmt.Printf("Reserve created successfully with ID: %d\n", result.ReserveId)
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing reserve",
	Long:  "Update the configuration of an existing recording reserve",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		updateOption := epgstation.EditManualReserveOption{}

		// Only include fields that were explicitly set
		if cmd.Flags().Changed("allow-end-lack") {
			updateOption.AllowEndLack = allowEndLack
		}

		// Handle encode settings for update
		if cmd.Flags().Changed("encode-mode") || cmd.Flags().Changed("encode-parent-directory") || cmd.Flags().Changed("encode-directory") {
			encodeOpt := &epgstation.ReserveEncodedOption{}
			if cmd.Flags().Changed("encode-mode") {
				modeStr := strconv.Itoa(encodeMode)
				encodeOpt.Mode1 = &modeStr
			}
			if cmd.Flags().Changed("encode-parent-directory") {
				encodeOpt.EncodeParentDirectoryName1 = &encodeParentDir
			}
			if cmd.Flags().Changed("encode-directory") {
				encodeOpt.Directory1 = &encodeDir
			}
			updateOption.EncodeOption = encodeOpt
		}

		err = client.UpdateReserve(reserveID, updateOption)
		if err != nil {
			return fmt.Errorf("failed to update reserve: %w", err)
		}

		fmt.Printf("Reserve %d updated successfully\n", reserveID)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a reserve",
	Long:  "Delete/cancel a recording reserve",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.DeleteReserve(reserveID)
		if err != nil {
			return fmt.Errorf("failed to delete reserve: %w", err)
		}

		fmt.Printf("Reserve %d deleted successfully\n", reserveID)
		return nil
	},
}

var unskipCmd = &cobra.Command{
	Use:   "unskip <id>",
	Short: "Remove skip status from a reserve",
	Long:  "Remove skip status from a recording reserve",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reserveID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid reserve ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.UnSkipReserve(reserveID)
		if err != nil {
			return fmt.Errorf("failed to unskip reserve: %w", err)
		}

		fmt.Printf("Reserve %d unskipped successfully\n", reserveID)
		return nil
	},
}

func printReservesTable(reserves *epgstation.Reserves) error {
	table := output.NewTable()
	table.SetHeader([]string{"ID", "Program", "Channel", "Start Time", "End Time", "Status", "Skip"})

	for _, reserve := range reserves.Reserves {
		status := "Scheduled"
		if reserve.IsConflict {
			status = "Conflict"
		}

		skipStatus := "No"
		if reserve.IsSkip {
			skipStatus = "Yes"
		}

		table.Append([]string{
			strconv.Itoa(reserve.Id),
			reserve.Name,
			"", // reserve.Program.ChannelName,
			output.FormatUnixTime(reserve.StartAt),
			output.FormatUnixTime(reserve.EndAt),
			status,
			skipStatus,
		})
	}

	table.Render()
	return nil
}

func init() {
	// List command flags
	listCmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	listCmd.Flags().IntVar(&limit, "limit", 0, "Limit for pagination")
	listCmd.Flags().IntVar(&channel, "channel", 0, "Filter by channel ID")
	listCmd.Flags().StringVar(&status, "status", "", "Filter by status")
	listCmd.Flags().StringVar(&keyword, "keyword", "", "Search by program title")

	// Create command flags
	createCmd.Flags().IntVar(&programID, "program-id", 0, "Program ID to reserve (required)")
	createCmd.Flags().IntVar(&encodeMode, "encode-mode", 0, "Encode mode ID")
	createCmd.Flags().StringVar(&encodeParentDir, "encode-parent-directory", "", "Encode parent directory")
	createCmd.Flags().StringVar(&encodeDir, "encode-directory", "", "Encode directory")
	createCmd.Flags().BoolVar(&allowEndLack, "allow-end-lack", false, "Allow recording to end early")
	createCmd.MarkFlagRequired("program-id")

	// Update command flags (optional)
	updateCmd.Flags().IntVar(&encodeMode, "encode-mode", 0, "Encode mode ID")
	updateCmd.Flags().StringVar(&encodeParentDir, "encode-parent-directory", "", "Encode parent directory")
	updateCmd.Flags().StringVar(&encodeDir, "encode-directory", "", "Encode directory")
	updateCmd.Flags().BoolVar(&allowEndLack, "allow-end-lack", false, "Allow recording to end early")

	// Add subcommands
	reservesCmd.AddCommand(listCmd)
	reservesCmd.AddCommand(createCmd)
	reservesCmd.AddCommand(updateCmd)
	reservesCmd.AddCommand(deleteCmd)
	reservesCmd.AddCommand(unskipCmd)

	// Register with root command
	root.AddCommand(reservesCmd)
}