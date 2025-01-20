package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func openInBrowser(m *types.Map) {
	// Add your callback logic here
	status, err := generator.GetStatus(logger, m)
	if err != nil {
		fmt.Printf("Error getting status: %v\n", err)
		os.Exit(1)
	}
	if status.Data.URL == "" {
		fmt.Println("No URL found")
		os.Exit(1)
	}
	if m.SavedConfig == "" && m.Staging {
		// add query param staging for procedural maps
		status.Data.URL += "?staging=true"
	}
	if err = browser.OpenURL(status.Data.URL); err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open generated maps in the browser",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateOpenFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		csv, _ := cmd.Flags().GetString("csv")
		savedConfig, _ := cmd.Flags().GetString("saved-config")
		seed, _ := cmd.Flags().GetString("seed")
		size, _ := cmd.Flags().GetInt("size")
		staging, _ := cmd.Flags().GetBool("staging")
		force, _ := cmd.Flags().GetBool("force")
		random, _ := cmd.Flags().GetBool("random")

		loadFromParams(csv, seed, size, savedConfig, staging, force, random)

		maps := generator.GetMaps()
		if len(maps) == 0 {
			fmt.Println("No maps were loaded")
			os.Exit(1)
		}

		if len(maps) == 1 {
			openInBrowser(maps[0])
			os.Exit(0)
		}

		var items []string
		for _, m := range maps {
			if m.Status == common.StatusComplete {
				items = append(items, m.String())
			}
		}

		if len(items) == 0 {
			fmt.Printf("Loaded %d maps, but none are complete", len(maps))
			os.Exit(1)
		}

		// Create the prompt
		prompt := promptui.Select{
			Label: "Select a map",
			Items: items,
			Templates: &promptui.SelectTemplates{
				// Customize how items are displayed
				Label:    "{{ . }}?",
				Active:   "🗺️  {{ . | cyan }}", // Display active item with an icon and color
				Inactive: "  {{ . | faint }}",
				Selected: "📍  You selected: {{ . | green }}",
			},
		}

		for {
			// Run the prompt
			_, result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed: %v\n", err)
				os.Exit(1)
			}

			// Find the selected map
			for _, m := range maps {
				if result == m.String() {
					openInBrowser(m)
					break
				}
			}
		}
	},
}

func init() {
	openCmd.Flags().StringP("csv", "c", "", "Path to the CSV file")
	openCmd.Flags().StringP("saved-config", "S", "", "Saved config to use from RustMaps")
	openCmd.Flags().StringP("seed", "s", "", "Seed to open")
	openCmd.Flags().IntP("size", "z", 0, "Size of the map to open")
	openCmd.Flags().BoolP("staging", "b", false, "Open maps against staging branch")
}

// validateOpenFlags checks mutual exclusivity and other flag rules
func validateOpenFlags(cmd *cobra.Command) error {
	csv, _ := cmd.Flags().GetString("csv")
	savedConfig, _ := cmd.Flags().GetString("saved-config")
	seed, _ := cmd.Flags().GetString("seed")
	size, _ := cmd.Flags().GetInt("size")
	staging, _ := cmd.Flags().GetBool("staging")

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

	if !(csv != "" || (size != 0 && seed != "")) {
		return fmt.Errorf("must provide either --csv, or --size and --seed with or without --staging")
	}

	return nil
}
