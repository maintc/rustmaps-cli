// Package rustmaps provides functionality for generating and managing Rust game maps
package rustmaps

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/maintc/rustmaps-cli/pkg/api"
	"github.com/maintc/rustmaps-cli/pkg/types"

	"go.uber.org/zap"
)

var (
	tierLimits = map[int]string{
		250:  "Free",
		500:  "Supporter",
		800:  "Premium",
		1000: "Organization 1",
		1750: "Organization 2",
	}
)

// Generator handles map generation and management
type Generator struct {
	config       types.Config
	maps         []*types.Map
	target       string
	rmcli        api.RustMapsClientBase
	configPath   string
	importsDir   string
	downloadsDir string
	logPath      string
	baseDir      string
	backoffTime  time.Duration
}

// NewGenerator creates a new Generator instance
func NewGenerator(baseDir *string) (*Generator, error) {
	g := &Generator{
		config:      types.Config{Tier: "Free"},
		backoffTime: 30 * time.Second,
	}

	if baseDir == nil {
		var err error
		baseDir = new(string)
		*baseDir, err = os.UserHomeDir()
		if err != nil {
			return nil, err
		}
	}

	g.baseDir = filepath.Join(*baseDir, ".rustmaps")
	return g, nil
}

// InitDirs initializes required directories
func (g *Generator) InitDirs() error {
	g.configPath = filepath.Join(g.baseDir, "config.json")
	g.importsDir = filepath.Join(g.baseDir, "imports")
	g.downloadsDir = filepath.Join(g.baseDir, "downloads")
	g.logPath = filepath.Join(g.baseDir, "generator.log")

	dirs := []string{g.baseDir, g.importsDir, g.downloadsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (g *Generator) CanGenerate(log *zap.Logger) bool {
	limits, err := g.rmcli.GetLimits(log)
	if err != nil {
		fmt.Printf("Error getting limits: %v\n", err)
		return false
	}

	var canGenerateConcurrent = limits.Data.Concurrent.Current < limits.Data.Concurrent.Allowed
	var canGenerateMonthly = limits.Data.Monthly.Current < limits.Data.Monthly.Allowed

	if !canGenerateConcurrent {
		fmt.Println("Cannot generate map: concurrent limit reached")
	}

	if !canGenerateMonthly {
		fmt.Println("Cannot generate map: monthly limit reached")
	}

	return canGenerateConcurrent && canGenerateMonthly
}

func (g *Generator) GetStatus(log *zap.Logger, m *types.Map) (*api.RustMapsStatusResponse, error) {
	status, err := g.rmcli.GetStatus(log, m)
	if err != nil {
		fmt.Printf("Error getting status: %v\n", err)
		return nil, err
	}

	return status, nil
}

func (g *Generator) SyncStatus(log *zap.Logger, m *types.Map) error {
	status, err := g.rmcli.GetStatus(log, m)
	if err != nil {
		fmt.Printf("Error getting status: %v\n", err)
		return err
	}

	m.ReportStatus(status.Meta.Status)
	m.SaveJSON(g.importsDir)
	return nil
}

func (g *Generator) AddMap(m *types.Map) {
	m.SetFilename()
	g.maps = append(g.maps, m)
}
