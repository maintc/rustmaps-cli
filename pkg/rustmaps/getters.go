package rustmaps

import "github.com/maintc/rustmaps-cli/pkg/types"

func (g *Generator) GetDownloadsDir() string {
	return g.downloadsDir
}

func (g *Generator) GetImportDir() string {
	return g.importsDir
}

func (g *Generator) GetLogPath() string {
	return g.logPath
}

func (g *Generator) GetConfigPath() string {
	return g.configPath
}

func (g *Generator) GetMaps() []*types.Map {
	return g.maps
}
