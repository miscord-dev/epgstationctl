package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/miscord-dev/epgstationctl/internal/epgstation"
)

func handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response body", resp.StatusCode)
	}

	var apiError epgstation.Error
	if err := json.Unmarshal(body, &apiError); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("API error %d: %s", apiError.Code, apiError.Message)
}

func parseJSONResponse(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}
