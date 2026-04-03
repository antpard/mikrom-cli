package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the mikrom CLI version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isJSON() {
			data, _ := json.MarshalIndent(map[string]string{"version": Version}, "", "  ")
			fmt.Println(string(data))
			return nil
		}
		fmt.Printf("mikrom version %s\n", Version)
		return nil
	},
}
