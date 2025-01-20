package rustmaps

import (
	"encoding/json"
	"os"

	"github.com/maintc/rustmaps-cli/pkg/api"
)

// LoadConfig loads the configuration from disk
func (g *Generator) LoadConfig() error {
	data, err := os.ReadFile(g.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config if it doesn't exist
			return g.SaveConfig()
		}
		return err
	}

	if err := json.Unmarshal(data, &g.config); err != nil {
		return err
	}

	g.rmcli = api.NewRustMapsClient(g.config.APIKey)

	return nil
}

// SaveConfig saves the current configuration to disk
func (g *Generator) SaveConfig() error {
	data, err := json.MarshalIndent(g.config, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(g.configPath, data, 0644)
}
