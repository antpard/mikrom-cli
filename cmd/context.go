package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage connection contexts (API URL + credentials)",
}

var contextListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contexts",
	RunE: func(cmd *cobra.Command, args []string) error {
		active := cfg.ActiveContext()

		type row struct {
			Name          string `json:"name"`
			APIURL        string `json:"api_url"`
			Authenticated bool   `json:"authenticated"`
			Active        bool   `json:"active"`
		}

		// Always include the implicit "default" context from the flat fields.
		rows := []row{}
		if len(cfg.Contexts) == 0 {
			rows = append(rows, row{
				Name:          "default",
				APIURL:        cfg.APIURL,
				Authenticated: cfg.Token != "",
				Active:        true,
			})
		} else {
			for name, entry := range cfg.Contexts {
				rows = append(rows, row{
					Name:          name,
					APIURL:        entry.APIURL,
					Authenticated: entry.Token != "",
					Active:        name == active,
				})
			}
		}

		if isJSON() {
			data, _ := json.MarshalIndent(rows, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("  %-16s  %-35s  %s\n", "NAME", "API URL", "AUTH")
		for _, r := range rows {
			marker := " "
			if r.Active {
				marker = "*"
			}
			auth := "no"
			if r.Authenticated {
				auth = "yes"
			}
			fmt.Printf("%s %-16s  %-35s  %s\n", marker, r.Name, r.APIURL, auth)
		}
		return nil
	},
}

var contextShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the active context",
	RunE: func(cmd *cobra.Command, args []string) error {
		active := cfg.ActiveContext()

		if isJSON() {
			data, _ := json.MarshalIndent(map[string]any{
				"name":          active,
				"api_url":       cfg.APIURL,
				"authenticated": cfg.Token != "",
			}, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Name:          %s\n", active)
		fmt.Printf("API URL:       %s\n", cfg.APIURL)
		auth := "no"
		if cfg.Token != "" {
			auth = "yes"
		}
		fmt.Printf("Authenticated: %s\n", auth)
		return nil
	},
}

var contextUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch to a named context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.UseContext(args[0]); err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Switched to context %q\n", args[0])
		return nil
	},
}

var contextAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add or update a named context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("api-url")
		tok, _ := cmd.Flags().GetString("token")

		cfg.AddContext(args[0], url, tok)
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Context %q saved (api-url: %s)\n", args[0], url)
		return nil
	},
}

var contextRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a named context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.RemoveContext(args[0]); err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Context %q removed\n", args[0])
		return nil
	},
}

func init() {
	contextAddCmd.Flags().String("api-url", "http://localhost:8080", "API URL for this context")
	contextAddCmd.Flags().String("token", "", "Authentication token for this context")
	contextAddCmd.MarkFlagRequired("api-url")

	contextCmd.AddCommand(contextListCmd)
	contextCmd.AddCommand(contextShowCmd)
	contextCmd.AddCommand(contextUseCmd)
	contextCmd.AddCommand(contextAddCmd)
	contextCmd.AddCommand(contextRemoveCmd)
}
