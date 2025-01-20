package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"

	"go.uber.org/zap"
)

type RustMapsGenerateResponseMeta struct {
	Status     string   `json:"status"`
	StatusCode int      `json:"statusCode"`
	Errors     []string `json:"errors"`
}

type RustMapsGenerateResponseData struct {
	MapID                string    `json:"mapId"`
	QueuePosition        int       `json:"queuePosition"`
	State                string    `json:"state"`
	CurrentStep          string    `json:"currentStep"`
	LastGeneratorPingUtc time.Time `json:"lastGeneratorPingUtc"`
}

type RustMapsGenerateResponse struct {
	Meta RustMapsGenerateResponseMeta `json:"meta"`
	Data RustMapsGenerateResponseData `json:"data"`
}

type RustMapsGenerateProceduralRequest struct {
	Size    int    `json:"size"`
	Seed    string `json:"seed"`
	Staging bool   `json:"staging"`
}

type RustMapsGenerateCustomRequest struct {
	MapParameters RustMapsGenerateProceduralRequest `json:"mapParameters"`
	ConfigName    string                            `json:"configName"`
}

func (c *RustMapsClient) GenerateCustom(log *zap.Logger, m *types.Map) (*RustMapsGenerateResponse, error) {
	c.rateLimiter.Wait()
	// Create a client with custom timeouts
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	data := RustMapsGenerateCustomRequest{
		MapParameters: RustMapsGenerateProceduralRequest{
			Size:    m.Size,
			Seed:    m.Seed,
			Staging: m.Staging,
		},
		ConfigName: m.SavedConfig,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create request
	log.Debug("Generating custom map", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/maps/custom/saved-config", c.ApiUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error making request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	log.Debug("Response status", zap.Int("status", resp.StatusCode), zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response", zap.Error(err))
		return nil, err
	}

	log.Debug("Response body", zap.String("body", string(body)), zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		log.Error("Unauthorized request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusUnauthorized)
		return nil, fmt.Errorf(common.StatusUnauthorized)
	case http.StatusForbidden:
		log.Error("Forbidden request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusForbidden)
		return nil, fmt.Errorf(common.StatusForbidden)
	case http.StatusConflict:
		log.Debug("Map already generating", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusGenerating)
		return nil, nil
	}

	var generateResponse RustMapsGenerateResponse
	if err := json.Unmarshal(body, &generateResponse); err != nil {
		return nil, err
	}

	m.MapID = generateResponse.Data.MapID
	switch resp.StatusCode {
	case http.StatusOK:
		log.Debug("Map generated", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusComplete)
		return &generateResponse, nil
	case 201:
		log.Debug("Map generating", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusGenerating)
		return &generateResponse, nil
	case http.StatusBadRequest:
		log.Error("Bad request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))
		if len(generateResponse.Meta.Errors) > 0 && generateResponse.Meta.Errors[0] == "Staging is not enabled" {
			m.ReportStatus(common.StatusStagingNotEnabled)
		} else {
			m.ReportStatus(common.StatusBadRequest)
		}
		return nil, fmt.Errorf(common.StatusBadRequest)
	}

	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func (c *RustMapsClient) GenerateProcedural(log *zap.Logger, m *types.Map) (*RustMapsGenerateResponse, error) {
	c.rateLimiter.Wait()
	// Create a client with custom timeouts
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	data := RustMapsGenerateProceduralRequest{
		Size:    m.Size,
		Seed:    m.Seed,
		Staging: m.Staging,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create request
	log.Debug("Generating procedural map", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/maps", c.ApiUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error making request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	log.Debug("Response status", zap.Int("status", resp.StatusCode), zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response", zap.Error(err))
		return nil, err
	}

	log.Debug("Response body", zap.String("body", string(body)), zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))

	switch resp.StatusCode {
	case http.StatusBadRequest:
		log.Error("Bad request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusBadRequest)
		return nil, fmt.Errorf(common.StatusBadRequest)
	case http.StatusUnauthorized:
		log.Error("Unauthorized request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusUnauthorized)
		return nil, fmt.Errorf(common.StatusUnauthorized)
	case http.StatusForbidden:
		log.Error("Forbidden request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusForbidden)
		return nil, fmt.Errorf(common.StatusForbidden)
	case http.StatusConflict:
		log.Debug("Map already generating", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusGenerating)
		return nil, nil
	}

	var generateResponse RustMapsGenerateResponse
	if err := json.Unmarshal(body, &generateResponse); err != nil {
		return nil, err
	}

	m.MapID = generateResponse.Data.MapID

	switch resp.StatusCode {
	case http.StatusOK:
		log.Debug("Map generated", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusComplete)
		return &generateResponse, nil
	case 201:
		log.Debug("Map generating", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		m.ReportStatus(common.StatusGenerating)
		return &generateResponse, nil
	}

	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
