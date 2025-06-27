package programs

import (
	"fmt"
	"time"

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
	channelID int
	date      string
	days      int
	halfWidth bool
	limit     int
)

var programsCmd = &cobra.Command{
	Use:   "programs",
	Short: "Manage programs",
	Long:  "Commands for managing EPGStation programs and schedules",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List programs",
	Long:  "List programs from EPGStation schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetSchedulesParams{
			GR:          true,
			BS:          true,
			CS:          true,
			SKY:         true,
			IsHalfWidth: halfWidth,
		}

		// Parse date if provided
		if date != "" {
			t, err := time.Parse("2006-01-02", date)
			if err != nil {
				return fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
			}
			params.StartAt = epgstation.StartAt(t.Unix() * 1000) // Convert to milliseconds
		}

		// Add channel filter if provided
		if channelID > 0 {
			chID := channelID
			// Get specific channel schedule
			channelParams := &epgstation.GetSchedulesChannelIdParams{
				IsHalfWidth: halfWidth,
			}
			if date != "" {
				channelParams.StartAt = params.StartAt
			}
			if days > 0 {
				channelParams.Days = days
			}

			schedules, err := client.GetSchedulesChannelId(nil, chID, channelParams)
			if err != nil {
				return fmt.Errorf("failed to get channel schedule: %w", err)
			}

			var formatter output.Formatter
			switch cfg.Output.Format {
			case outputFormatJSON:
				formatter = output.NewJSONFormatter(nil)
				return formatter.Format(*schedules)
			default:
				formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
				// Custom formatting for Schedules - flatten to show programs with channel info
				return formatSchedulesAsTable(*schedules, formatter)
			}
		}

		schedules, err := client.GetSchedules(params)
		if err != nil {
			return fmt.Errorf("failed to get schedules: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case outputFormatJSON:
			formatter = output.NewJSONFormatter(nil)
			return formatter.Format(*schedules)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
			// Custom formatting for Schedules - flatten to show programs with channel info
			return formatSchedulesAsTable(*schedules, formatter)
		}
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show currently broadcasting programs",
	Long:  "Show programs currently being broadcast",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetSchedulesBroadcastingParams{
			IsHalfWidth: halfWidth,
		}

		schedules, err := client.GetBroadcastingPrograms(params)
		if err != nil {
			return fmt.Errorf("failed to get broadcasting programs: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case outputFormatJSON:
			formatter = output.NewJSONFormatter(nil)
			return formatter.Format(*schedules)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
			// Custom formatting for Schedules - flatten to show programs with channel info
			return formatSchedulesAsTable(*schedules, formatter)
		}
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search programs by keyword",
	Long:  "Search programs by keyword in title or description",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		searchKeyword := args[0]
		trueVal := true
		searchOption := epgstation.ScheduleSearchOption{
			Option: epgstation.RuleSearchOption{
				Keyword:     &searchKeyword,
				Name:        &trueVal,
				Description: &trueVal,
				GR:          &trueVal,
				BS:          &trueVal,
				CS:          &trueVal,
				SKY:         &trueVal,
			},
			IsHalfWidth: halfWidth,
		}

		if limit > 0 {
			limitFloat := float32(limit)
			searchOption.Limit = &limitFloat
		}

		programs, err := client.PostSchedulesSearch(nil, searchOption)
		if err != nil {
			return fmt.Errorf("failed to search programs: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case outputFormatJSON:
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(*programs)
	},
}

func init() {
	// Add flags to list command
	listCmd.Flags().IntVarP(&channelID, "channel", "c", 0, "Channel ID to filter")
	listCmd.Flags().StringVarP(&date, "date", "d", "", "Date to show programs for (YYYY-MM-DD)")
	listCmd.Flags().IntVar(&days, "days", 1, "Number of days to show")
	listCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	// Add flags to current command
	currentCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	// Add flags to search command
	searchCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")
	searchCmd.Flags().IntVarP(&limit, "limit", "l", 50, "Maximum number of results")

	programsCmd.AddCommand(listCmd)
	programsCmd.AddCommand(currentCmd)
	programsCmd.AddCommand(searchCmd)
	root.AddCommand(programsCmd)
}

// formatSchedulesAsTable formats schedule data in a user-friendly table format
func formatSchedulesAsTable(schedules epgstation.Schedules, formatter output.Formatter) error {
	if len(schedules) == 0 {
		fmt.Println("No programs found")
		return nil
	}

	// Create a flat list of programs with channel information
	type ProgramWithChannel struct {
		ChannelName string
		ChannelType string
		ProgramName string
		StartTime   string
		EndTime     string
		Description string
	}

	var programs []ProgramWithChannel
	for _, schedule := range schedules {
		channelName := schedule.Channel.Name

		for _, program := range schedule.Programs {
			startTime := formatUnixTime(int64(program.StartAt))
			endTime := formatUnixTime(int64(program.EndAt))

			description := ""
			if program.Description != nil {
				description = *program.Description
				if len(description) > 50 {
					description = description[:50] + "..."
				}
			}

			programs = append(programs, ProgramWithChannel{
				ChannelName: channelName,
				ChannelType: string(schedule.Channel.ChannelType),
				ProgramName: program.Name,
				StartTime:   startTime,
				EndTime:     endTime,
				Description: description,
			})
		}
	}

	return formatter.Format(programs)
}

// formatUnixTime converts Unix timestamp (in milliseconds) to readable time
func formatUnixTime(unixTimeMS int64) string {
	t := time.Unix(unixTimeMS/1000, 0)
	return t.Format("15:04")
}
