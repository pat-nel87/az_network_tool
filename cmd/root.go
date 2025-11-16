package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "azure-network-analyzer",
	Short: "Analyze Azure network topology",
	Long: `A comprehensive CLI tool to analyze and visualize Azure network resources.

This tool collects and analyzes Azure network topology including:
- Virtual Networks and Subnets
- Network Security Groups (NSGs)
- Private Endpoints and DNS Zones
- Route Tables and NAT Gateways
- VPN and ExpressRoute Gateways
- Load Balancers and Application Gateways
- Network Watcher insights

It generates detailed reports and network topology visualizations to help
understand your Azure network infrastructure.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
