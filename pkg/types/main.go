package types

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/maintc/rustmaps-cli/pkg/common"
)

// Config represents the application configuration
type Config struct {
	APIKey string `json:"api_key"`
	Tier   string `json:"tier"`
}

// Map represents a single map configuration
type Map struct {
	Seed        string `json:"seed"`
	Size        int    `json:"size"`
	SavedConfig string `json:"saved_config,omitempty"`
	Staging     bool   `json:"staging"`
	MapID       string `json:"map_id,omitempty"`
	Status      string `json:"status"`
	LastSync    string `json:"last_sync,omitempty"`
	Filename    string `json:"filename,omitempty"`
}

func NewMap(seed string, size int, savedConfig string, staging bool) *Map {
	m := &Map{
		Seed:        seed,
		Size:        size,
		SavedConfig: savedConfig,
		Staging:     staging,
	}
	m.Status = common.StatusPending
	return m
}

func (m *Map) SetFilename() {
	// export map to g.importsDir/m.seed_m.size_m.saved_config_staging
	filename := fmt.Sprintf("%s_%d", m.Seed, m.Size)

	if m.SavedConfig != "" {
		filename = fmt.Sprintf("%s_%s", filename, m.SavedConfig)
	}

	if m.Staging {
		filename = fmt.Sprintf("%s_staging", filename)
	}

	filename = fmt.Sprintf("%s.json", filename)
	m.Filename = filename
}

func (m *Map) ReportStatus(status string) {
	m.Status = status
	fmt.Println(m.String())
	m.MarkSynced()
}

func (m *Map) String() string {
	return fmt.Sprintf("Seed: %s | Size: %d | Config: '%s' | Status: '%s'", m.Seed, m.Size, m.SavedConfig, m.Status)
}

func (m *Map) ShouldSync() bool {
	// Parse the last sync time
	lastSyncTime, err := time.Parse(time.RFC3339, m.LastSync)
	if err != nil {
		// If parsing fails, we assume the map should sync
		return true
	}

	// Check if the current time is more than 5 minutes after the last sync time
	return time.Now().After(lastSyncTime.Add(5 * time.Minute))
}

func (m *Map) MarkSynced() {
	// Update the LastSync field with the current time in RFC3339 format
	m.LastSync = time.Now().Format(time.RFC3339)
}

func (m *Map) MergeFrom(other Map) {
	m.Seed = other.Seed
	m.Size = other.Size
	m.SavedConfig = other.SavedConfig
	m.Staging = other.Staging
	m.MapID = other.MapID
	m.Status = other.Status
	m.LastSync = other.LastSync
}

func (m *Map) SaveJSON(outputDir string) error {
	// Check if the Filename field is set
	if m.Filename == "" {
		return fmt.Errorf("map does not have its filename set")
	}

	outputPath := filepath.Join(outputDir, m.Filename)

	// Open the file for writing (create it if it doesn't exist)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create a JSON encoder and write the struct to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: Pretty-print the JSON
	if err := encoder.Encode(m); err != nil {
		return fmt.Errorf("failed to encode map to JSON: %w", err)
	}

	return nil
}

// CSVInfo represents a map and its entry count
type CSVInfo struct {
	Name  string
	Count int
}
