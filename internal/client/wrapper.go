package client

import (
	"context"
	"net/http"
	"time"

	"github.com/miscord-dev/epgstationctl/internal/config"
	"github.com/miscord-dev/epgstationctl/internal/epgstation"
)

// RuleWithID represents a complete rule with ID information
// This fixes the issue where the generated Rule type doesn't include the ID field
type RuleWithID struct {
	ID          int  `json:"id"`
	ReservesCnt *int `json:"reservesCnt,omitempty"`
	epgstation.AddRuleOption
}

// RulesWithID represents the corrected rules response
type RulesWithID struct {
	Rules []RuleWithID `json:"rules"`
	Total int          `json:"total"`
}

type EPGStationClient struct {
	client *epgstation.Client
	config *config.Config
}

func NewEPGStationClient(cfg *config.Config) (*EPGStationClient, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Server.Timeout) * time.Second,
	}

	client, err := epgstation.NewClient(cfg.GetAPIBaseURL(), epgstation.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &EPGStationClient{
		client: client,
		config: cfg,
	}, nil
}

func (c *EPGStationClient) GetChannels() (*epgstation.ChannelItems, error) {
	resp, err := c.client.GetChannels(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var channels epgstation.ChannelItems
	if err := parseJSONResponse(resp, &channels); err != nil {
		return nil, err
	}

	return &channels, nil
}

func (c *EPGStationClient) GetSchedules(params *epgstation.GetSchedulesParams) (*epgstation.Schedules, error) {
	resp, err := c.client.GetSchedules(context.Background(), params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var schedules epgstation.Schedules
	if err := parseJSONResponse(resp, &schedules); err != nil {
		return nil, err
	}

	return &schedules, nil
}

func (c *EPGStationClient) GetBroadcastingPrograms(params *epgstation.GetSchedulesBroadcastingParams) (*epgstation.Schedules, error) {
	resp, err := c.client.GetSchedulesBroadcasting(context.Background(), params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var schedules epgstation.Schedules
	if err := parseJSONResponse(resp, &schedules); err != nil {
		return nil, err
	}

	return &schedules, nil
}

func (c *EPGStationClient) GetRecordings(params *epgstation.GetRecordingParams) (*epgstation.Records, error) {
	resp, err := c.client.GetRecording(context.Background(), params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var records epgstation.Records
	if err := parseJSONResponse(resp, &records); err != nil {
		return nil, err
	}

	return &records, nil
}

func (c *EPGStationClient) GetSchedulesChannelId(ctx interface{}, channelId epgstation.ChannelId, params *epgstation.GetSchedulesChannelIdParams) (*epgstation.Schedules, error) {
	resp, err := c.client.GetSchedulesChannelId(context.Background(), channelId, params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var schedules epgstation.Schedules
	if err := parseJSONResponse(resp, &schedules); err != nil {
		return nil, err
	}

	return &schedules, nil
}

func (c *EPGStationClient) PostSchedulesSearch(ctx interface{}, body epgstation.ScheduleSearchOption) (*[]epgstation.ScheduleProgramItem, error) {
	resp, err := c.client.PostSchedulesSearch(context.Background(), body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var programs []epgstation.ScheduleProgramItem
	if err := parseJSONResponse(resp, &programs); err != nil {
		return nil, err
	}

	return &programs, nil
}

// GetRules retrieves all recording rules
func (c *EPGStationClient) GetRules(params *epgstation.GetRulesParams) (*RulesWithID, error) {
	resp, err := c.client.GetRules(context.Background(), params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var rules RulesWithID
	if err := parseJSONResponse(resp, &rules); err != nil {
		return nil, err
	}

	return &rules, nil
}

// GetRule retrieves a specific recording rule by ID
func (c *EPGStationClient) GetRule(ruleId int) (*RuleWithID, error) {
	resp, err := c.client.GetRulesRuleId(context.Background(), ruleId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var rule RuleWithID
	if err := parseJSONResponse(resp, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// CreateRule creates a new recording rule
func (c *EPGStationClient) CreateRule(ruleOption epgstation.AddRuleOption) (*epgstation.AddedRule, error) {
	resp, err := c.client.PostRules(context.Background(), ruleOption)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return nil, handleErrorResponse(resp)
	}

	var result epgstation.AddedRule
	if err := parseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateRule updates an existing recording rule
func (c *EPGStationClient) UpdateRule(ruleId int, ruleOption epgstation.AddRuleOption) error {
	resp, err := c.client.PutRulesRuleId(context.Background(), ruleId, ruleOption)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// DeleteRule deletes a recording rule
func (c *EPGStationClient) DeleteRule(ruleId int) error {
	resp, err := c.client.DeleteRulesRuleId(context.Background(), ruleId)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// EnableRule enables a recording rule
func (c *EPGStationClient) EnableRule(ruleId int) error {
	resp, err := c.client.PutRulesRuleIdEnable(context.Background(), ruleId)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// DisableRule disables a recording rule
func (c *EPGStationClient) DisableRule(ruleId int) error {
	resp, err := c.client.PutRulesRuleIdDisable(context.Background(), ruleId)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// SearchRulesKeyword searches for keywords that can be used in rules
func (c *EPGStationClient) SearchRulesKeyword(params *epgstation.GetRulesKeywordParams) (*[]string, error) {
	resp, err := c.client.GetRulesKeyword(context.Background(), params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var keywords []string
	if err := parseJSONResponse(resp, &keywords); err != nil {
		return nil, err
	}

	return &keywords, nil
}
