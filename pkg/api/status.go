package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"

	"go.uber.org/zap"
)

type RustMapsStatusResponseMeta struct {
	Status     string   `json:"status"`
	StatusCode int      `json:"statusCode"`
	Errors     []string `json:"errors"`
}

type RustMapsStatusResponseDataCoordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type RustMapsStatusResponseDataMonuments struct {
	Type         string                                `json:"type"`
	Coordinates  RustMapsStatusResponseDataCoordinates `json:"coordinates"`
	NameOverride string                                `json:"nameOverride"`
}

type RustMapsStatusResponseDataBiomePercentages struct {
	S float64 `json:"s"`
	D float64 `json:"d"`
	F float64 `json:"f"`
	T float64 `json:"t"`
	J float64 `json:"j"`
}

type RustMapsStatusResponseData struct {
	ID                  string                                     `json:"id"`
	Type                string                                     `json:"type"`
	Seed                int                                        `json:"seed"`
	Size                int                                        `json:"size"`
	SaveVersion         int                                        `json:"saveVersion"`
	URL                 string                                     `json:"url"`
	RawImageURL         string                                     `json:"rawImageUrl"`
	ImageURL            string                                     `json:"imageUrl"`
	ImageIconURL        string                                     `json:"imageIconUrl"`
	ThumbnailURL        string                                     `json:"thumbnailUrl"`
	IsStaging           bool                                       `json:"isStaging"`
	IsCustomMap         bool                                       `json:"isCustomMap"`
	CanDownload         bool                                       `json:"canDownload"`
	DownloadURL         string                                     `json:"downloadUrl"`
	TotalMonuments      int                                        `json:"totalMonuments"`
	Monuments           []RustMapsStatusResponseDataMonuments      `json:"monuments"`
	LandPercentageOfMap int                                        `json:"landPercentageOfMap"`
	BiomePercentages    RustMapsStatusResponseDataBiomePercentages `json:"biomePercentages"`
	Islands             int                                        `json:"islands"`
	Mountains           int                                        `json:"mountains"`
	IceLakes            int                                        `json:"iceLakes"`
	Rivers              int                                        `json:"rivers"`
	Lakes               int                                        `json:"lakes"`
	Canyons             int                                        `json:"canyons"`
	Oases               int                                        `json:"oases"`
	BuildableRocks      int                                        `json:"buildableRocks"`
}

type RustMapsStatusResponse struct {
	Meta RustMapsStatusResponseMeta `json:"meta"`
	Data RustMapsStatusResponseData `json:"data"`
}

func (c *RustMapsClient) GetStatus(log *zap.Logger, m *types.Map) (*RustMapsStatusResponse, error) {
	c.rateLimiter.Wait()
	// Create a client with custom timeouts
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var endpoint = m.MapID
	if endpoint == "" {
		endpoint = fmt.Sprintf("%d/%s", m.Size, m.Seed)
	}

	// Create request
	log.Debug("Getting map status", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/maps/%s", c.ApiUrl, endpoint), nil)
	if err != nil {
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

	log.Debug("Response status", zap.Int("status", resp.StatusCode), zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response", zap.Error(err))
		return nil, err
	}

	log.Debug("Response body", zap.String("body", string(body)), zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.String("config", m.SavedConfig), zap.Bool("staging", m.Staging))

	status := &RustMapsStatusResponse{}
	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.Unmarshal(body, status); err != nil {
			return nil, err
		}
		status.Meta.Status = common.StatusComplete
		return status, nil
	case http.StatusUnauthorized:
		log.Error("Unauthorized request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		status.Meta.Status = common.StatusUnauthorized
		status.Meta.StatusCode = http.StatusUnauthorized
	case http.StatusForbidden:
		log.Error("Forbidden request", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		status.Meta.Status = common.StatusForbidden
		status.Meta.StatusCode = http.StatusForbidden
	case http.StatusNotFound:
		log.Error("Map not found", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		status.Meta.Status = common.StatusNotFound
		status.Meta.StatusCode = http.StatusNotFound
	case http.StatusConflict:
		log.Debug("Map generating", zap.String("seed", m.Seed), zap.Int("size", m.Size), zap.Bool("staging", m.Staging))
		status.Meta.Status = common.StatusGenerating
		status.Meta.StatusCode = http.StatusConflict
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return status, nil
}
