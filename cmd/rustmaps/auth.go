package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth [api-key]",
	Short: "Authenticate with RustMaps API",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := args[0]
		generator.SetApiKey(apiKey)
		tier, ok := generator.DetermineTier(logger)
		if !ok {
			fmt.Println("Provided API key is invalid, check logs for more info")
			os.Exit(1)
		}
		generator.SetTier(tier)
		if err := generator.SaveConfig(); err != nil {
			fmt.Println("Error saving config, check logs for more info")
			os.Exit(1)
		}
		fmt.Printf("API key verified: 🗺️ %s Subscriber\n", tier)
	},
}
