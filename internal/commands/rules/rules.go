package rules

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

const (
	outputFormatJSON = "json"
	channelsAll      = "All"
	booleanFalse     = "false"
)

var (
	offset    int
	limit     int
	halfWidth bool
	// Create rule flags
	keyword      string
	searchName   bool
	searchDesc   bool
	enableGR     bool
	enableBS     bool
	enableCS     bool
	allowEndLack bool
	avoidDupe    bool
)

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage recording rules",
	Long:  "Commands for managing EPGStation recording rules",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recording rules",
	Long:  "List all recording rules with their status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		params := &epgstation.GetRulesParams{}

		if offset > 0 {
			params.Offset = &offset
		}

		if limit > 0 {
			params.Limit = &limit
		}

		rules, err := client.GetRules(params)
		if err != nil {
			return fmt.Errorf("failed to get rules: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case outputFormatJSON:
			formatter = output.NewJSONFormatter(nil)
			return formatter.Format(*rules)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
			return formatRulesAsTable(*rules, formatter)
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show <rule-id>",
	Short: "Show rule details",
	Long:  "Show detailed information about a specific recording rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		ruleId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid rule ID: %w", err)
		}

		rule, err := client.GetRule(ruleId)
		if err != nil {
			return fmt.Errorf("failed to get rule: %w", err)
		}

		var formatter output.Formatter
		switch cfg.Output.Format {
		case outputFormatJSON:
			formatter = output.NewJSONFormatter(nil)
			return formatter.Format(*rule)
		default:
			formatter = output.NewTableFormatter(nil, cfg.Output.NoHeader)
			return formatRuleDetailsAsTable(*rule, formatter)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <rule-id>",
	Short: "Delete a recording rule",
	Long:  "Delete a recording rule by its ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		ruleId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid rule ID: %w", err)
		}

		if err := client.DeleteRule(ruleId); err != nil {
			return fmt.Errorf("failed to delete rule: %w", err)
		}

		fmt.Printf("Rule %d deleted successfully\n", ruleId)
		return nil
	},
}

var enableCmd = &cobra.Command{
	Use:   "enable <rule-id>",
	Short: "Enable a recording rule",
	Long:  "Enable a disabled recording rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		ruleId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid rule ID: %w", err)
		}

		if err := client.EnableRule(ruleId); err != nil {
			return fmt.Errorf("failed to enable rule: %w", err)
		}

		fmt.Printf("Rule %d enabled successfully\n", ruleId)
		return nil
	},
}

var disableCmd = &cobra.Command{
	Use:   "disable <rule-id>",
	Short: "Disable a recording rule",
	Long:  "Disable a recording rule without deleting it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		ruleId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid rule ID: %w", err)
		}

		if err := client.DisableRule(ruleId); err != nil {
			return fmt.Errorf("failed to disable rule: %w", err)
		}

		fmt.Printf("Rule %d disabled successfully\n", ruleId)
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new recording rule",
	Long:  "Create a new recording rule with specified options",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := root.GetConfig()
		client, err := client.NewEPGStationClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create EPGStation client: %w", err)
		}

		if keyword == "" {
			return fmt.Errorf("keyword is required for creating a rule")
		}

		// Build the rule option
		skyEnabled := false
		ruleOption := epgstation.AddRuleOption{
			IsTimeSpecification: false,
			SearchOption: epgstation.RuleSearchOption{
				Keyword:     &keyword,
				Name:        &searchName,
				Description: &searchDesc,
				GR:          &enableGR,
				BS:          &enableBS,
				CS:          &enableCS,
				SKY:         &skyEnabled,
			},
			ReserveOption: epgstation.RuleReserveOption{
				Enable:         true,
				AllowEndLack:   allowEndLack,
				AvoidDuplicate: avoidDupe,
			},
		}

		result, err := client.CreateRule(ruleOption)
		if err != nil {
			return fmt.Errorf("failed to create rule: %w", err)
		}

		fmt.Printf("Rule created successfully with ID: %d\n", result.RuleId)
		return nil
	},
}

func init() {
	// Add flags to list command
	listCmd.Flags().IntVarP(&offset, "offset", "", 0, "Offset for pagination")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit number of results")
	listCmd.Flags().BoolVar(&halfWidth, "half-width", true, "Use half-width characters")

	// Add flags to create command
	createCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "Keyword to search for (required)")
	createCmd.Flags().BoolVar(&searchName, "search-name", true, "Search in program name")
	createCmd.Flags().BoolVar(&searchDesc, "search-description", false, "Search in program description")
	createCmd.Flags().BoolVar(&enableGR, "gr", true, "Enable for GR (terrestrial) channels")
	createCmd.Flags().BoolVar(&enableBS, "bs", true, "Enable for BS channels")
	createCmd.Flags().BoolVar(&enableCS, "cs", false, "Enable for CS channels")
	createCmd.Flags().BoolVar(&allowEndLack, "allow-end-lack", false, "Allow recordings that end early")
	createCmd.Flags().BoolVar(&avoidDupe, "avoid-duplicate", false, "Avoid duplicate recordings")

	rulesCmd.AddCommand(listCmd)
	rulesCmd.AddCommand(showCmd)
	rulesCmd.AddCommand(createCmd)
	rulesCmd.AddCommand(deleteCmd)
	rulesCmd.AddCommand(enableCmd)
	rulesCmd.AddCommand(disableCmd)
	root.AddCommand(rulesCmd)
}

// formatRulesAsTable formats rules data in a user-friendly table format
func formatRulesAsTable(rules client.RulesWithID, formatter output.Formatter) error {
	if len(rules.Rules) == 0 {
		fmt.Println("No recording rules found")
		return nil
	}

	// Create a flat list of rules with essential information
	type RuleInfo struct {
		ID       int
		Status   string
		Keyword  string
		Channels string
	}

	ruleList := make([]RuleInfo, 0, len(rules.Rules))
	for _, rule := range rules.Rules {
		status := "✗"
		if rule.ReserveOption.Enable {
			if rule.IsTimeSpecification {
				status = "⏰"
			} else {
				status = "✓"
			}
		}

		keyword := ""
		if rule.SearchOption.Keyword != nil {
			keyword = *rule.SearchOption.Keyword
		}

		var channelParts []string
		if rule.SearchOption.GR != nil && *rule.SearchOption.GR {
			channelParts = append(channelParts, "地上波")
		}
		if rule.SearchOption.BS != nil && *rule.SearchOption.BS {
			channelParts = append(channelParts, "BS")
		}
		if rule.SearchOption.CS != nil && *rule.SearchOption.CS {
			channelParts = append(channelParts, "CS")
		}
		
		channels := channelsAll
		if len(channelParts) > 0 {
			channels = strings.Join(channelParts, ",")
		}

		ruleList = append(ruleList, RuleInfo{
			ID:       rule.ID,
			Status:   status,
			Keyword:  keyword,
			Channels: channels,
		})
	}

	return formatter.Format(ruleList)
}

// formatRuleDetailsAsTable formats a single rule's details in a user-friendly table format
func formatRuleDetailsAsTable(rule client.RuleWithID, formatter output.Formatter) error {
	type RuleDetail struct {
		Field string
		Value string
	}

	var details []RuleDetail

	details = append(details, RuleDetail{"ID", fmt.Sprintf("%d", rule.ID)})

	if rule.ReservesCnt != nil {
		details = append(details, RuleDetail{"ReserveCount", fmt.Sprintf("%d", *rule.ReservesCnt)})
	} else {
		details = append(details, RuleDetail{"ReserveCount", "0"})
	}

	details = append(details, RuleDetail{"TimeSpecification", fmt.Sprintf("%t", rule.IsTimeSpecification)})

	// Search options
	if rule.SearchOption.Keyword != nil {
		details = append(details, RuleDetail{"Keyword", *rule.SearchOption.Keyword})
	}

	channels := []string{}
	if rule.SearchOption.GR != nil && *rule.SearchOption.GR {
		channels = append(channels, "GR")
	}
	if rule.SearchOption.BS != nil && *rule.SearchOption.BS {
		channels = append(channels, "BS")
	}
	if rule.SearchOption.CS != nil && *rule.SearchOption.CS {
		channels = append(channels, "CS")
	}
	if rule.SearchOption.SKY != nil && *rule.SearchOption.SKY {
		channels = append(channels, "SKY")
	}
	if len(channels) > 0 {
		details = append(details, RuleDetail{"Channels", fmt.Sprintf("%v", channels)})
	}

	searchName := booleanFalse
	if rule.SearchOption.Name != nil {
		searchName = fmt.Sprintf("%t", *rule.SearchOption.Name)
	}
	details = append(details, RuleDetail{"SearchName", searchName})

	searchDescription := booleanFalse
	if rule.SearchOption.Description != nil {
		searchDescription = fmt.Sprintf("%t", *rule.SearchOption.Description)
	}
	details = append(details, RuleDetail{"SearchDescription", searchDescription})

	searchExtended := booleanFalse
	if rule.SearchOption.Extended != nil {
		searchExtended = fmt.Sprintf("%t", *rule.SearchOption.Extended)
	}
	details = append(details, RuleDetail{"SearchExtended", searchExtended})

	// Reserve options
	details = append(details, RuleDetail{"RuleEnabled", fmt.Sprintf("%t", rule.ReserveOption.Enable)})
	details = append(details, RuleDetail{"AllowEndLack", fmt.Sprintf("%t", rule.ReserveOption.AllowEndLack)})
	details = append(details, RuleDetail{"AvoidDuplicate", fmt.Sprintf("%t", rule.ReserveOption.AvoidDuplicate)})

	return formatter.Format(details)
}
