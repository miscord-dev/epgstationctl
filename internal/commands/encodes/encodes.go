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

const (
	outputFormatJSON = "json"
)

var (
	offset     int
	limit      int
	isEncoding bool
	
	// Create encode flags
	recordedID           int
	encodeMode           int
	encodeParentDir      string
	encodeDir            string
	removeOriginal       bool
)

var encodesCmd = &cobra.Command{
	Use:   "encodes",
	Short: "Manage encoding jobs",
	Long:  "Commands for managing EPGStation encoding jobs",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List encoding jobs",
	Long:  "List all encoding jobs with their status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		encodes, err := client.GetEncodes(&epgstation.GetEncodeParams{IsHalfWidth: false}) // Assuming IsHalfWidth is required
		if err != nil {
			return fmt.Errorf("failed to get encodes: %w", err)
		}

		if cfg.Output.Format == outputFormatJSON {
			return output.PrintAsJSON(encodes)
		}

		return printEncodesTable(encodes)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a manual encode job",
	Long:  "Create a manual encoding job for a recorded program",
	RunE: func(cmd *cobra.Command, args []string) error {
		if recordedID == 0 {
			return fmt.Errorf("recorded ID is required (use --recorded-id)")
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		encodeOption := epgstation.AddManualEncodeProgramOption{
			RecordedId: recordedID,
			Mode:       strconv.Itoa(encodeMode), // Mode is string
		}

		if encodeParentDir != "" {
			encodeOption.ParentDir = &encodeParentDir
		}
		if encodeDir != "" {
			encodeOption.Directory = &encodeDir
		}
		encodeOption.RemoveOriginal = removeOriginal // RemoveOriginal is bool, not *bool

		result, err := client.CreateEncode(encodeOption)
		if err != nil {
			return fmt.Errorf("failed to create encode: %w", err)
		}

		if cfg.Output.Format == outputFormatJSON {
			return output.PrintAsJSON(result)
		}

		fmt.Printf("Encode job created successfully with ID: %d\n", result.EncodeId)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an encode job",
	Long:  "Delete/cancel an encoding job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		encodeID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid encode ID: %s", args[0])
		}

		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		err = client.DeleteEncode(encodeID)
		if err != nil {
			return fmt.Errorf("failed to delete encode: %w", err)
		}

		fmt.Printf("Encode job %d deleted successfully\n", encodeID)
		return nil
	},
}

func printEncodesTable(encodes *epgstation.EncodeInfo) error {
	table := output.NewTable()
	table.SetHeader([]string{"ID", "Name", "Mode", "Status", "Progress", "Started"})

	for _, encode := range encodes.RunningItems {
		progress := "N/A"
		if encode.Percent != nil {
			progress = fmt.Sprintf("%.1f%%", *encode.Percent)
		}

		table.Append([]string{
			strconv.Itoa(encode.Id),
			encode.Recorded.Name,
			encode.Mode,
			"Encoding",
			progress,
			output.FormatUnixTime(encode.Recorded.StartAt),
		})
	}

	for _, encode := range encodes.WaitItems {
		table.Append([]string{
			strconv.Itoa(encode.Id),
			encode.Recorded.Name,
			encode.Mode,
			"Waiting",
			"N/A",
			output.FormatUnixTime(encode.Recorded.StartAt),
		})
	}

	table.Render()
	return nil
}

func init() {
	// Create command flags
	createCmd.Flags().IntVar(&recordedID, "recorded-id", 0, "Recorded program ID to encode (required)")
	createCmd.Flags().IntVar(&encodeMode, "encode-mode", 0, "Encode mode ID")
	createCmd.Flags().StringVar(&encodeParentDir, "encode-parent-directory", "", "Encode parent directory")
	createCmd.Flags().StringVar(&encodeDir, "encode-directory", "", "Encode directory")
	createCmd.Flags().BoolVar(&removeOriginal, "remove-original", false, "Remove original file after encoding")
	createCmd.MarkFlagRequired("recorded-id")
	createCmd.MarkFlagRequired("encode-mode") // Make encode-mode required

	// Add subcommands
	encodesCmd.AddCommand(listCmd)
	encodesCmd.AddCommand(createCmd)
	encodesCmd.AddCommand(deleteCmd)

	// Register with root command
	root.AddCommand(encodesCmd)
}