package rustmaps

import (
	"fmt"

	"go.uber.org/zap"
)

func (g *Generator) ValidateAuthentication(log *zap.Logger) error {
	if g.config.APIKey == "" {
		log.Error("API key not set")
		return fmt.Errorf("API key not set")
	}

	if g.config.Tier == "" {
		log.Error("Tier not set")
		return fmt.Errorf("tier not set")
	}

	return nil
}

func (g *Generator) DetermineTier(log *zap.Logger) (string, bool) {
	limits, err := g.rmcli.GetLimits(log)
	if err != nil {
		log.Error("Error getting limits", zap.Error(err))
		return "", false
	}

	if tier, exists := tierLimits[limits.Data.Monthly.Allowed]; exists {
		return tier, true
	}

	log.Error("Invalid tier", zap.Int("allowed", limits.Data.Monthly.Allowed))
	return "", false
}
