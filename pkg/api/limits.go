package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RustMapsLimitsResponseMeta struct {
	Status     string   `json:"status"`
	StatusCode int      `json:"statusCode"`
	Errors     []string `json:"errors"`
}

type RustMapsLimitsResponseDataConcurrent struct {
	Current int `json:"current"`
	Allowed int `json:"allowed"`
}

type RustMapsLimitsResponseDataMonthly struct {
	Current int `json:"current"`
	Allowed int `json:"allowed"`
}

type RustMapsLimitsResponseData struct {
	Concurrent RustMapsLimitsResponseDataConcurrent `json:"concurrent"`
	Monthly    RustMapsLimitsResponseDataMonthly    `json:"monthly"`
}

type RustMapsLimitsResponse struct {
	Meta RustMapsLimitsResponseMeta `json:"meta"`
	Data RustMapsLimitsResponseData `json:"data"`
}

func (c *RustMapsClient) GetLimits(log *zap.Logger) (*RustMapsLimitsResponse, error) {
	c.rateLimiter.Wait()
	// Create a client with custom timeouts
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	log.Debug("GET /maps/limits - Getting API limits")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/maps/limits", c.ApiUrl), nil)
	if err != nil {
		log.Error("Error creating request", zap.Error(err))
		return nil, err
	}

	// Add headers
	req.Header.Set("X-API-Key", c.apiKey)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error making request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	log.Debug("Response status", zap.Int("status", resp.StatusCode))

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response", zap.Error(err))
		return nil, err
	}

	log.Debug("Response body", zap.String("body", string(body)))

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		log.Error("Unauthorized request")
		return nil, fmt.Errorf("unauthorized")
	}

	limits := &RustMapsLimitsResponse{}
	if err := json.Unmarshal(body, limits); err != nil {
		log.Error("Error unmarshalling response", zap.Error(err))
		return nil, err
	}
	return limits, nil
}
