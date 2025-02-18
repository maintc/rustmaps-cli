package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate custom and procedural maps",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateGenerateFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		csv, _ := cmd.Flags().GetString("csv")
		savedConfig, _ := cmd.Flags().GetString("saved-config")
		seed, _ := cmd.Flags().GetString("seed")
		size, _ := cmd.Flags().GetInt("size")
		staging, _ := cmd.Flags().GetBool("staging")
		force, _ := cmd.Flags().GetBool("force")
		random, _ := cmd.Flags().GetBool("random")
		download, _ := cmd.Flags().GetBool("download")
		outputDir, _ := cmd.Flags().GetString("output-dir")

		if outputDir != "" {
			generator.OverrideDownloadsDir(logger, outputDir)
		}

		loadFromParams(csv, seed, size, savedConfig, staging, force, random)

		for {
			if !generator.Generate(logger) {
				break
			}
		}

		if download {
			now := time.Now()
			version := now.Format("2006-01-02_15-04-05")
			if err := generator.Download(logger, version); err != nil {
				fmt.Printf("Error downloading maps: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Maps downloaded to %s\n", filepath.Join(generator.GetDownloadsDir(), version))
		}
	},
}

func init() {
	generateCmd.Flags().StringP("csv", "c", "", "Path to the CSV file")
	generateCmd.Flags().StringP("saved-config", "S", "", "Saved config to use from RustMaps")
	generateCmd.Flags().StringP("seed", "s", "", "Seed to generate")
	generateCmd.Flags().IntP("size", "z", 0, "Size of the map to generate")
	generateCmd.Flags().BoolP("staging", "b", false, "Generate maps against staging branch")
	generateCmd.Flags().BoolP("force", "f", false, "Force generate even if map is already generated")
	generateCmd.Flags().BoolP("random", "r", false, "Randomly select the seed (size must be set)")
	generateCmd.Flags().BoolP("download", "d", false, "Download the generated custom maps (you can't download procedural maps)")
	generateCmd.Flags().StringP("output-dir", "o", "", "Output directory for downloaded maps")
}

// validateGenerateFlags checks mutual exclusivity and other flag rules
func validateGenerateFlags(cmd *cobra.Command) error {
	csv, _ := cmd.Flags().GetString("csv")
	savedConfig, _ := cmd.Flags().GetString("saved-config")
	seed, _ := cmd.Flags().GetString("seed")
	size, _ := cmd.Flags().GetInt("size")
	staging, _ := cmd.Flags().GetBool("staging")
	random, _ := cmd.Flags().GetBool("random")
	download, _ := cmd.Flags().GetBool("download")
	outputDir, _ := cmd.Flags().GetString("output-dir")

	// random can only be used with size
	if random && seed != "" {
		return fmt.Errorf("cannot use --random with --seed")
	}
	if random && csv != "" {
		return fmt.Errorf("cannot use --random with --csv")
	}
	if random && size == 0 {
		return fmt.Errorf("cannot use --random without --size")
	}

	// csv cannot be used with anything else
	if csv != "" && seed != "" {
		return fmt.Errorf("cannot use --csv with --seed")
	}
	if csv != "" && size != 0 {
		return fmt.Errorf("cannot use --csv with --size")
	}
	if csv != "" && savedConfig != "" {
		return fmt.Errorf("cannot use --csv with --saved-config")
	}
	if csv != "" && staging {
		return fmt.Errorf("cannot use --csv with --staging")
	}

	if outputDir != "" && !download {
		return fmt.Errorf("cannot use --output-dir without --download")
	}

	if !(csv != "" || (size != 0 && (seed != "" || random))) {
		return fmt.Errorf("must provide either --csv, or --size and --seed or --random")
	}

	return nil
}
