package client

import (
	"context"
	"net/http"
	"time"

	"github.com/miscord-dev/epgstationctl/internal/config"
	"github.com/miscord-dev/epgstationctl/internal/epgstation"
)

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var programs []epgstation.ScheduleProgramItem
	if err := parseJSONResponse(resp, &programs); err != nil {
		return nil, err
	}

	return &programs, nil
}
