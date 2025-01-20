package rustmaps

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

// LoadCSV reads the currently selected map file
func (g *Generator) LoadCSV(log *zap.Logger, mapsPath string) error {

	if err := g.ValidateAuthentication(log); err != nil {
		log.Error("Error validating authentication", zap.Error(err))
		return err
	}

	if err := g.ValidateCSV(log, mapsPath); err != nil {
		log.Error("Error validating map file", zap.Error(err))
		return err
	}

	g.maps = nil
	g.target = mapsPath

	file, err := os.Open(g.target)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Allow variable number of fields per record
	reader.FieldsPerRecord = -1

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	g.maps = make([]*types.Map, 0, len(records))
	for _, record := range records {
		m := types.Map{
			Status: common.StatusPending, // Set default status
		}

		// Safely assign fields based on what's available
		if len(record) > 0 {
			m.Seed = record[0]
		}
		if len(record) > 1 {
			m.Size = parseInt(record[1])
		}
		if len(record) > 2 {
			m.SavedConfig = record[2]
		}
		if len(record) > 3 {
			m.Staging = parseBool(record[3])
		}
		if len(record) > 4 {
			m.MapID = record[4]
		}
		if len(record) > 5 {
			m.Status = record[5]
		}

		m.SetFilename()
		g.maps = append(g.maps, &m)
	}

	if len(g.maps) == 0 {
		log.Warn("No maps loaded")
		return fmt.Errorf("no maps loaded")
	}

	if (g.config.Tier == "Free" || g.config.Tier == "Supporter") && g.ContainCustomMaps() {
		log.Warn("Cannot generate custom maps with Free or Supporter tier")
		return fmt.Errorf("cannot generate custom maps with Free or Supporter tier")
	}

	return nil
}

func (g *Generator) ValidateCSV(log *zap.Logger, mapsPath string) error {

	if _, err := os.Stat(mapsPath); err != nil {
		log.Error("Error reading file", zap.Error(err))
		return err
	}

	// get first line in the file
	file, err := os.Open(mapsPath)
	if err != nil {
		log.Error("Error opening file", zap.Error(err))
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Error("Error reading CSV headers", zap.Error(err))
		return err
	}

	requiredColumns := []string{"seed", "size"}
	for _, req := range requiredColumns {
		found := false
		for _, col := range headers {
			if strings.TrimSpace(col) == req {
				found = true
				break
			}
		}
		if !found {
			log.Error("Missing required column", zap.String("column", req))
			return fmt.Errorf("missing required column: %s", req)
		}
	}

	return nil
}

// Import imports a CSV file containing map definitions
func (g *Generator) Import(log *zap.Logger, force bool) error {

	// if err := g.ValidateCSV(log, mapsPath); err != nil {
	// 	log.Error("Error validating map file", zap.Error(err))
	// 	return err
	// }

	for _, m := range g.maps {

		if m.Filename == "" {
			m.SetFilename()
		}

		path := filepath.Join(g.importsDir, m.Filename)

		if !force {
			if _, err := os.Stat(path); err == nil {
				log.Debug("Map file already exists, loading existing file", zap.String("path", path))

				// Read and deserialize the existing JSON file
				file, err := os.Open(path)
				if err != nil {
					log.Error("Error opening existing file", zap.Error(err), zap.String("path", path))
					return err
				}
				defer file.Close()

				var existingMap types.Map // Replace with the actual type of your map structure
				if err := json.NewDecoder(file).Decode(&existingMap); err != nil {
					log.Error("Error decoding existing map file", zap.Error(err), zap.String("path", path))
					return err
				}

				// Merge the fields from existingMap into m
				m.MergeFrom(existingMap)
				continue
			}
		}

		// export json to file
		file, err := os.Create(path)
		if err != nil {
			log.Error("Error creating file", zap.Error(err), zap.String("path", path))
			return err
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err := enc.Encode(m); err != nil {
			log.Error("Error encoding map", zap.Error(err), zap.String("map", m.String()), zap.String("path", path))
			return err
		}

		log.Debug("Exported map", zap.String("map", m.String()), zap.String("path", path))
	}

	return nil
}
