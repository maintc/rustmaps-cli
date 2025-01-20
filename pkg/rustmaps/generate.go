package rustmaps

import (
	"time"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"go.uber.org/zap"
)

func (g *Generator) Generate(log *zap.Logger) bool {

	if err := g.ValidateAuthentication(log); err != nil {
		log.Error("Error validating authentication", zap.Error(err))
		return false
	}

	if !(g.Pending() || g.Generating()) {
		log.Info("All maps are complete")
		return false
	}

	for _, m := range g.maps {
		if m.Status == common.StatusComplete && m.ShouldSync() {
			if err := g.SyncStatus(log, m); err != nil {
				log.Error("Error syncing status", zap.String("seed", m.Seed))
			}

			if m.Status == common.StatusNotFound {
				m.Status = common.StatusPending
			}
		}

		if m.Status == common.StatusGenerating {
			if err := g.SyncStatus(log, m); err != nil {
				log.Error("Error syncing status", zap.String("seed", m.Seed))
				continue
			}
		}
	}

	if !(g.Pending() && g.CanGenerate(log)) {
		time.Sleep(g.backoffTime)
		return true
	}

	for _, m := range g.maps {
		if m.Status == common.StatusPending {
			if m.SavedConfig == "" {
				g.rmcli.GenerateProcedural(log, m)
			} else {
				g.rmcli.GenerateCustom(log, m)
			}

			if err := m.SaveJSON(g.importsDir); err != nil {
				log.Error("Error saving map file", zap.Error(err))
			}

			break
		}
	}

	time.Sleep(2 * time.Second)
	return true
}
