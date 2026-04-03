package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spluca/mikrom-cli/internal/api"
)

var ippoolCmd = &cobra.Command{
	Use:   "ippool",
	Short: "Manage IP pools",
}

var ippoolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all IP pools",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()
		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")

		resp, err := newClient().ListIPPools(page, pageSize)
		if err != nil {
			return err
		}

		if isJSON() {
			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(resp.Items) == 0 {
			fmt.Println("No IP pools found")
			return nil
		}

		fmt.Printf("%-6s  %-20s  %-18s  %-15s  %-15s  %s\n", "ID", "NAME", "CIDR", "START", "END", "ACTIVE")
		for _, p := range resp.Items {
			fmt.Printf("%-6d  %-20s  %-18s  %-15s  %-15s  %v\n",
				p.ID, p.Name, p.CIDR, p.StartIP, p.EndIP, p.IsActive)
		}
		fmt.Printf("\nTotal: %d\n", resp.Total)
		return nil
	},
}

var ippoolGetCmd = &cobra.Command{
	Use:   "get <pool-id>",
	Short: "Get IP pool details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid pool ID %q: must be a number", args[0])
		}

		pool, err := newClient().GetIPPool(id)
		if err != nil {
			return err
		}

		if isJSON() {
			data, _ := json.MarshalIndent(pool, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		printIPPool(pool)
		return nil
	},
}

var ippoolCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new IP pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()
		name, _ := cmd.Flags().GetString("name")
		network, _ := cmd.Flags().GetString("network")
		cidr, _ := cmd.Flags().GetString("cidr")
		gateway, _ := cmd.Flags().GetString("gateway")
		startIP, _ := cmd.Flags().GetString("start-ip")
		endIP, _ := cmd.Flags().GetString("end-ip")

		pool, err := newClient().CreateIPPool(api.CreateIPPoolRequest{
			Name:    name,
			Network: network,
			CIDR:    cidr,
			Gateway: gateway,
			StartIP: startIP,
			EndIP:   endIP,
		})
		if err != nil {
			return err
		}

		if isJSON() {
			data, _ := json.MarshalIndent(pool, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("IP pool created: %d\n", pool.ID)
		printIPPool(pool)
		return nil
	},
}

var ippoolUpdateCmd = &cobra.Command{
	Use:   "update <pool-id>",
	Short: "Update an IP pool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid pool ID %q: must be a number", args[0])
		}

		req := api.UpdateIPPoolRequest{}
		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			req.Name = &name
		}
		if cmd.Flags().Changed("active") {
			active, _ := cmd.Flags().GetBool("active")
			req.IsActive = &active
		}

		pool, err := newClient().UpdateIPPool(id, req)
		if err != nil {
			return err
		}

		if isJSON() {
			data, _ := json.MarshalIndent(pool, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		printIPPool(pool)
		return nil
	},
}

var ippoolDeleteCmd = &cobra.Command{
	Use:   "delete <pool-id>",
	Short: "Delete an IP pool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid pool ID %q: must be a number", args[0])
		}

		if err := newClient().DeleteIPPool(id); err != nil {
			return err
		}

		fmt.Printf("IP pool %d deleted\n", id)
		return nil
	},
}

var ippoolStatsCmd = &cobra.Command{
	Use:   "stats <pool-id>",
	Short: "Show IP allocation stats for a pool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid pool ID %q: must be a number", args[0])
		}

		stats, err := newClient().GetIPPoolStats(id)
		if err != nil {
			return err
		}

		if isJSON() {
			data, _ := json.MarshalIndent(stats, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Pool:      %s (ID: %d)\n", stats.PoolName, stats.PoolID)
		fmt.Printf("Total:     %d\n", stats.Total)
		fmt.Printf("Allocated: %d\n", stats.Allocated)
		fmt.Printf("Available: %d\n", stats.Available)
		fmt.Printf("Usage:     %.2f%%\n", stats.UsagePercent)
		return nil
	},
}

var ippoolAllStatsCmd = &cobra.Command{
	Use:   "all-stats",
	Short: "Show IP allocation stats for all pools",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		stats, err := newClient().GetAllPoolStats()
		if err != nil {
			return err
		}

		if isJSON() {
			data, _ := json.MarshalIndent(stats, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(stats) == 0 {
			fmt.Println("No IP pools found")
			return nil
		}

		fmt.Printf("%-6s  %-20s  %-8s  %-10s  %-10s  %s\n", "ID", "NAME", "TOTAL", "ALLOCATED", "AVAILABLE", "USAGE%")
		for _, s := range stats {
			fmt.Printf("%-6d  %-20s  %-8d  %-10d  %-10d  %.2f%%\n",
				s.PoolID, s.PoolName, s.Total, s.Allocated, s.Available, s.UsagePercent)
		}
		return nil
	},
}

var ippoolSuggestRangeCmd = &cobra.Command{
	Use:   "suggest-range",
	Short: "Suggest a usable IP range for a given CIDR",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()
		cidr, _ := cmd.Flags().GetString("cidr")

		result, err := newClient().SuggestIPRange(cidr)
		if err != nil {
			return err
		}

		if isJSON() {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("CIDR:              %s\n", result.CIDR)
		fmt.Printf("Network address:   %s\n", result.NetworkAddress)
		fmt.Printf("Broadcast address: %s\n", result.BroadcastAddress)
		fmt.Printf("First usable IP:   %s\n", result.FirstUsableIP)
		fmt.Printf("Last usable IP:    %s\n", result.LastUsableIP)
		fmt.Printf("Total hosts:       %d\n", result.TotalHosts)
		fmt.Printf("Suggested start:   %s\n", result.SuggestedStart)
		fmt.Printf("Suggested end:     %s\n", result.SuggestedEnd)
		return nil
	},
}

func printIPPool(p *api.IPPool) {
	fmt.Printf("ID:       %d\n", p.ID)
	fmt.Printf("Name:     %s\n", p.Name)
	fmt.Printf("Network:  %s\n", p.Network)
	fmt.Printf("CIDR:     %s\n", p.CIDR)
	fmt.Printf("Gateway:  %s\n", p.Gateway)
	fmt.Printf("Start:    %s\n", p.StartIP)
	fmt.Printf("End:      %s\n", p.EndIP)
	fmt.Printf("Active:   %v\n", p.IsActive)
}

func init() {
	ippoolListCmd.Flags().Int("page", 1, "Page number")
	ippoolListCmd.Flags().Int("page-size", 20, "Items per page")

	ippoolCreateCmd.Flags().String("name", "", "Pool name")
	ippoolCreateCmd.Flags().String("network", "", "Network address (e.g. 10.100.0.0)")
	ippoolCreateCmd.Flags().String("cidr", "", "Network CIDR (e.g. 10.100.0.0/24)")
	ippoolCreateCmd.Flags().String("gateway", "", "Gateway IP")
	ippoolCreateCmd.Flags().String("start-ip", "", "First usable IP")
	ippoolCreateCmd.Flags().String("end-ip", "", "Last usable IP")
	ippoolCreateCmd.MarkFlagRequired("name")
	ippoolCreateCmd.MarkFlagRequired("network")
	ippoolCreateCmd.MarkFlagRequired("cidr")
	ippoolCreateCmd.MarkFlagRequired("gateway")
	ippoolCreateCmd.MarkFlagRequired("start-ip")
	ippoolCreateCmd.MarkFlagRequired("end-ip")

	ippoolUpdateCmd.Flags().String("name", "", "New pool name")
	ippoolUpdateCmd.Flags().Bool("active", true, "Set pool active/inactive")

	ippoolSuggestRangeCmd.Flags().String("cidr", "", "CIDR to suggest a range for")
	ippoolSuggestRangeCmd.MarkFlagRequired("cidr")

	ippoolCmd.AddCommand(ippoolListCmd)
	ippoolCmd.AddCommand(ippoolGetCmd)
	ippoolCmd.AddCommand(ippoolCreateCmd)
	ippoolCmd.AddCommand(ippoolUpdateCmd)
	ippoolCmd.AddCommand(ippoolDeleteCmd)
	ippoolCmd.AddCommand(ippoolStatsCmd)
	ippoolCmd.AddCommand(ippoolAllStatsCmd)
	ippoolCmd.AddCommand(ippoolSuggestRangeCmd)
}
