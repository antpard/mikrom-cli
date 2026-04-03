package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check API connectivity and health",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := newClient().Health()
		if err != nil {
			return fmt.Errorf("API unreachable: %w", err)
		}

		if isJSON() {
			data, _ := json.MarshalIndent(map[string]string{
				"status":  resp.Status,
				"api_url": cfg.APIURL,
			}, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("API URL: %s\n", cfg.APIURL)
		fmt.Printf("Status:  %s\n", resp.Status)
		return nil
	},
}
