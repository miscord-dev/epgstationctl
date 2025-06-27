package channels

import (
	"fmt"

	"github.com/miscord-dev/epgstationctl/internal/client"
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
	"github.com/miscord-dev/epgstationctl/internal/output"
	"github.com/spf13/cobra"
)

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "Manage channels",
	Long:  "Commands for managing EPGStation channels",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all channels",
	Long:  "List all available channels from EPGStation",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		channels, err := client.GetChannels()
		if err != nil {
			return fmt.Errorf("failed to get channels: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(*channels)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <channel-id>",
	Short: "Show channel details",
	Long:  "Show detailed information about a specific channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		channels, err := client.GetChannels()
		if err != nil {
			return fmt.Errorf("failed to get channels: %w", err)
		}

		channelID := args[0]
		var targetChannel interface{}

		for _, channel := range *channels {
			if fmt.Sprintf("%d", channel.Id) == channelID {
				targetChannel = channel
				break
			}
		}

		if targetChannel == nil {
			return fmt.Errorf("channel with ID %s not found", channelID)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case "json":
			formatter = output.NewJSONFormatter(nil)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
		}

		return formatter.Format(targetChannel)
	},
}

func init() {
	channelsCmd.AddCommand(listCmd)
	channelsCmd.AddCommand(showCmd)
	root.AddCommand(channelsCmd)
}