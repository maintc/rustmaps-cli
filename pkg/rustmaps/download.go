package rustmaps

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"go.uber.org/zap"
)

func (g *Generator) OverrideDownloadsDir(log *zap.Logger, dir string) {
	g.downloadsDir = dir
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Error("Error creating downloads directory", zap.Error(err))
	}
}

// DownloadFile downloads a file using net/http
func (g *Generator) DownloadFile(log *zap.Logger, url, target string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Error creating request", zap.Error(err))
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error downloading file", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("Error downloading file", zap.String("status", resp.Status))
		return fmt.Errorf("error downloading file: %s", resp.Status)
	}

	file, err := os.Create(target)
	if err != nil {
		log.Error("Error creating file", zap.Error(err))
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Error("Error writing file", zap.Error(err))
		return err
	}

	return nil
}
func (g *Generator) Download(log *zap.Logger, version string) error {
	if len(g.maps) == 0 {
		log.Warn("No maps loaded")
		return fmt.Errorf("no maps loaded")
	}

	for _, m := range g.maps {
		if m.Status != common.StatusComplete {
			continue
		}

		if status, err := g.rmcli.GetStatus(log, m); err != nil {
			log.Error("Error downloading map", zap.String("seed", m.Seed), zap.Error(err))
			return err
		} else {
			if !status.Data.CanDownload {
				log.Warn("Cannot download map", zap.String("seed", m.Seed), zap.Int("size", m.Size))
				fmt.Println()
				stagingFlag := ""
				if m.Staging {
					stagingFlag = " -b"
				}
				fmt.Printf("But you can open it in the browser: `rustmaps open -s '%s' -z %d -S '%s'%s`\n", m.Seed, m.Size, m.SavedConfig, stagingFlag)
				fmt.Println()
				continue
			}
			downloadsDir := filepath.Join(g.downloadsDir, version)
			if err := os.MkdirAll(downloadsDir, 0755); err != nil {
				log.Error("Error creating downloads directory", zap.Error(err))
				return err
			}
			mapTarget := filepath.Join(downloadsDir, fmt.Sprintf("%s_%d_%s.map", m.Seed, m.Size, m.MapID))
			imageTarget := filepath.Join(downloadsDir, fmt.Sprintf("%s_%d_%s.png", m.Seed, m.Size, m.MapID))
			imageWithIconsTarget := filepath.Join(downloadsDir, fmt.Sprintf("%s_%d_%s_icons.png", m.Seed, m.Size, m.MapID))
			thumbnailTarget := filepath.Join(downloadsDir, fmt.Sprintf("%s_%d_%s_thumbnail.png", m.Seed, m.Size, m.MapID))
			fmt.Printf("Download URL: %s\n", status.Data.DownloadURL)
			if err := g.DownloadFile(log, status.Data.DownloadURL, mapTarget); err != nil {
				log.Error("Error downloading map", zap.String("seed", m.Seed), zap.Error(err))
				return err
			}
			if err := g.DownloadFile(log, status.Data.ImageURL, imageTarget); err != nil {
				log.Error("Error downloading image", zap.String("seed", m.Seed), zap.Error(err))
				return err
			}
			if err := g.DownloadFile(log, status.Data.ImageIconURL, imageWithIconsTarget); err != nil {
				log.Error("Error downloading image with icons", zap.String("seed", m.Seed), zap.Error(err))
				return err
			}
			if err := g.DownloadFile(log, status.Data.ThumbnailURL, thumbnailTarget); err != nil {
				log.Error("Error downloading thumbnail", zap.String("seed", m.Seed), zap.Error(err))
				return err
			}
		}
	}

	return nil
}
