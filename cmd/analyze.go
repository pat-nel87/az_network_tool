package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"azure-network-analyzer/pkg/analyzer"
	"azure-network-analyzer/pkg/azure"
	"azure-network-analyzer/pkg/models"
	"azure-network-analyzer/pkg/reporter"
	"azure-network-analyzer/pkg/visualization"
	"github.com/spf13/cobra"
)

var (
	subscriptionID string
	resourceGroup  string
	outputFormat   string
	outputPath     string
	includeViz     bool
	vizFormat      string
	dryRun         bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze network topology for a resource group",
	Long: `Collect and analyze Azure network resources in a specified resource group.

This command will:
1. Connect to Azure using DefaultAzureCredential
2. Collect all network resources from the specified resource group
3. Analyze the topology and identify security findings
4. Generate reports in the specified format
5. Optionally create network topology visualizations`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&subscriptionID, "subscription", "s", "", "Azure subscription ID (required)")
	analyzeCmd.Flags().StringVarP(&resourceGroup, "resource-group", "g", "", "Resource group name (required)")
	analyzeCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "markdown", "Output format (json|markdown|html)")
	analyzeCmd.Flags().StringVarP(&outputPath, "output", "f", "", "Output file path (defaults to stdout)")
	analyzeCmd.Flags().BoolVar(&includeViz, "visualize", true, "Generate network topology diagram")
	analyzeCmd.Flags().StringVar(&vizFormat, "viz-format", "svg", "Visualization format (svg|png|dot)")
	analyzeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Use mock data instead of connecting to Azure (for testing)")

	analyzeCmd.MarkFlagRequired("subscription")
	analyzeCmd.MarkFlagRequired("resource-group")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	fmt.Println("Azure Network Topology Analyzer")
	fmt.Println("================================")
	fmt.Printf("Subscription: %s\n", subscriptionID)
	fmt.Printf("Resource Group: %s\n", resourceGroup)
	fmt.Printf("Output Format: %s\n", outputFormat)
	if dryRun {
		fmt.Println("Mode: DRY-RUN (using mock data)")
	}
	fmt.Println()

	var topology *models.NetworkTopology

	if dryRun {
		// Use mock data for testing
		fmt.Println("Generating mock network topology...")
		topology = azure.GenerateMockTopology(subscriptionID, resourceGroup)
		fmt.Printf("  - Generated %d VNets\n", len(topology.VirtualNetworks))
		fmt.Printf("  - Generated %d NSGs\n", len(topology.NSGs))
		fmt.Printf("  - Generated %d Private Endpoints\n", len(topology.PrivateEndpoints))
		fmt.Printf("  - Generated %d Private DNS Zones\n", len(topology.PrivateDNSZones))
		fmt.Printf("  - Generated %d Route Tables\n", len(topology.RouteTables))
		fmt.Printf("  - Generated %d NAT Gateways\n", len(topology.NATGateways))
		fmt.Printf("  - Generated %d VPN Gateways\n", len(topology.VPNGateways))
		fmt.Printf("  - Generated %d ExpressRoute Circuits\n", len(topology.ERCircuits))
		fmt.Printf("  - Generated %d Load Balancers\n", len(topology.LoadBalancers))
		fmt.Printf("  - Generated %d Application Gateways\n", len(topology.AppGateways))
		if topology.NetworkWatcher != nil {
			fmt.Printf("  - Generated Network Watcher insights\n")
		}
	} else {
		// 1. Initialize Azure client
		fmt.Println("Initializing Azure client...")
		client, err := azure.NewAzureClient(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to create Azure client: %w", err)
		}

		// 2. Collect all network resources
		fmt.Println("Collecting network resources...")
		topology = &models.NetworkTopology{
			SubscriptionID: subscriptionID,
			ResourceGroup:  resourceGroup,
			Timestamp:      time.Now(),
		}

		// Collect VNets
		fmt.Println("  - Collecting Virtual Networks...")
		vnets, err := client.GetVirtualNetworks(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get virtual networks: %w", err)
		}
		topology.VirtualNetworks = vnets
		fmt.Printf("    Found %d VNets\n", len(vnets))

		// Collect NSGs
		fmt.Println("  - Collecting Network Security Groups...")
		nsgs, err := client.GetNetworkSecurityGroups(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get NSGs: %w", err)
		}
		topology.NSGs = nsgs
		fmt.Printf("    Found %d NSGs\n", len(nsgs))

		// Collect Private Endpoints
		fmt.Println("  - Collecting Private Endpoints...")
		privateEndpoints, err := client.GetPrivateEndpoints(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get private endpoints: %w", err)
		}
		topology.PrivateEndpoints = privateEndpoints
		fmt.Printf("    Found %d Private Endpoints\n", len(privateEndpoints))

		// Collect Private DNS Zones
		fmt.Println("  - Collecting Private DNS Zones...")
		dnsZones, err := client.GetPrivateDNSZones(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get private DNS zones: %w", err)
		}
		topology.PrivateDNSZones = dnsZones
		fmt.Printf("    Found %d Private DNS Zones\n", len(dnsZones))

		// Collect Route Tables
		fmt.Println("  - Collecting Route Tables...")
		routeTables, err := client.GetRouteTables(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get route tables: %w", err)
		}
		topology.RouteTables = routeTables
		fmt.Printf("    Found %d Route Tables\n", len(routeTables))

		// Collect NAT Gateways
		fmt.Println("  - Collecting NAT Gateways...")
		natGateways, err := client.GetNATGateways(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get NAT gateways: %w", err)
		}
		topology.NATGateways = natGateways
		fmt.Printf("    Found %d NAT Gateways\n", len(natGateways))

		// Collect VPN Gateways
		fmt.Println("  - Collecting VPN Gateways...")
		vpnGateways, err := client.GetVPNGateways(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get VPN gateways: %w", err)
		}
		topology.VPNGateways = vpnGateways
		fmt.Printf("    Found %d VPN Gateways\n", len(vpnGateways))

		// Collect ExpressRoute Circuits
		fmt.Println("  - Collecting ExpressRoute Circuits...")
		erCircuits, err := client.GetExpressRouteCircuits(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get ExpressRoute circuits: %w", err)
		}
		topology.ERCircuits = erCircuits
		fmt.Printf("    Found %d ExpressRoute Circuits\n", len(erCircuits))

		// Collect Load Balancers
		fmt.Println("  - Collecting Load Balancers...")
		loadBalancers, err := client.GetLoadBalancers(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get load balancers: %w", err)
		}
		topology.LoadBalancers = loadBalancers
		fmt.Printf("    Found %d Load Balancers\n", len(loadBalancers))

		// Collect Application Gateways
		fmt.Println("  - Collecting Application Gateways...")
		appGateways, err := client.GetApplicationGateways(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("failed to get application gateways: %w", err)
		}
		topology.AppGateways = appGateways
		fmt.Printf("    Found %d Application Gateways\n", len(appGateways))

		// Collect Network Watcher insights
		fmt.Println("  - Collecting Network Watcher insights...")
		nwInsights, err := client.GetNetworkWatcherInsights(ctx, resourceGroup)
		if err != nil {
			fmt.Printf("    Warning: Could not get Network Watcher insights: %v\n", err)
		} else {
			topology.NetworkWatcher = nwInsights
			fmt.Println("    Network Watcher insights collected")
		}
	}

	fmt.Println()
	fmt.Printf("Collection complete! Total resources: %d\n", countResources(topology))

	// 3. Analyze topology
	fmt.Println("\nAnalyzing topology...")
	analysisReport := analyzer.Analyze(topology)

	// Display analysis results
	displayAnalysisResults(analysisReport)

	// 4. Generate reports
	fmt.Println("Generating reports...")
	var reportContent []byte
	var reportExt string

	switch outputFormat {
	case "json":
		fmt.Println("  Generating JSON report...")
		content, err := reporter.GenerateJSON(topology, analysisReport)
		if err != nil {
			return fmt.Errorf("failed to generate JSON report: %w", err)
		}
		reportContent = content
		reportExt = ".json"
	case "markdown":
		fmt.Println("  Generating Markdown report...")
		content := reporter.GenerateMarkdown(topology, analysisReport)
		reportContent = []byte(content)
		reportExt = ".md"
	case "html":
		fmt.Println("  Generating HTML report...")
		content := reporter.GenerateHTML(topology, analysisReport)
		reportContent = []byte(content)
		reportExt = ".html"
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	// Write report to file or display info
	if outputPath != "" {
		err := writeOutput(reportContent)
		if err != nil {
			return fmt.Errorf("failed to write report: %w", err)
		}
		fmt.Printf("  Report saved to: %s\n", outputPath)
	} else {
		// Auto-generate filename
		timestamp := time.Now().Format("20060102-150405")
		filename := fmt.Sprintf("network-report-%s-%s%s", resourceGroup, timestamp, reportExt)
		err := os.WriteFile(filename, reportContent, 0644)
		if err != nil {
			return fmt.Errorf("failed to write report: %w", err)
		}
		fmt.Printf("  Report saved to: %s\n", filename)
	}

	// 5. Generate visualization if requested
	if includeViz {
		fmt.Println("\nGenerating topology diagram...")

		// Generate DOT content
		dotContent := visualization.GenerateDOTFile(topology)

		timestamp := time.Now().Format("20060102-150405")
		var vizFilename string
		var vizContent []byte
		var err error

		switch vizFormat {
		case "svg":
			fmt.Println("  Rendering SVG...")
			vizContent, err = visualization.RenderSVG(dotContent)
			if err != nil {
				fmt.Printf("  Warning: Could not render SVG: %v\n", err)
				fmt.Println("  Falling back to DOT file...")
				vizContent = []byte(dotContent)
				vizFilename = fmt.Sprintf("network-topology-%s-%s.dot", resourceGroup, timestamp)
			} else {
				vizFilename = fmt.Sprintf("network-topology-%s-%s.svg", resourceGroup, timestamp)
			}
		case "png":
			fmt.Println("  Rendering PNG...")
			vizContent, err = visualization.RenderPNG(dotContent)
			if err != nil {
				fmt.Printf("  Warning: Could not render PNG: %v\n", err)
				fmt.Println("  Falling back to DOT file...")
				vizContent = []byte(dotContent)
				vizFilename = fmt.Sprintf("network-topology-%s-%s.dot", resourceGroup, timestamp)
			} else {
				vizFilename = fmt.Sprintf("network-topology-%s-%s.png", resourceGroup, timestamp)
			}
		case "dot":
			fmt.Println("  Generating DOT file...")
			vizContent = []byte(dotContent)
			vizFilename = fmt.Sprintf("network-topology-%s-%s.dot", resourceGroup, timestamp)
		default:
			return fmt.Errorf("unsupported visualization format: %s", vizFormat)
		}

		err = os.WriteFile(vizFilename, vizContent, 0644)
		if err != nil {
			return fmt.Errorf("failed to write visualization: %w", err)
		}
		fmt.Printf("  Diagram saved to: %s\n", vizFilename)
	}

	fmt.Println("\nAnalysis complete!")
	return nil
}

// countResources returns the total number of resources in the topology
func countResources(topology *models.NetworkTopology) int {
	count := len(topology.VirtualNetworks)
	count += len(topology.NSGs)
	count += len(topology.PrivateEndpoints)
	count += len(topology.PrivateDNSZones)
	count += len(topology.RouteTables)
	count += len(topology.NATGateways)
	count += len(topology.VPNGateways)
	count += len(topology.ERCircuits)
	count += len(topology.LoadBalancers)
	count += len(topology.AppGateways)
	return count
}

// writeOutput writes content to the specified output path or stdout
func writeOutput(content []byte) error {
	if outputPath == "" {
		_, err := os.Stdout.Write(content)
		return err
	}

	return os.WriteFile(outputPath, content, 0644)
}

// displayAnalysisResults shows the analysis results in the console
func displayAnalysisResults(report *analyzer.AnalysisReport) {
	// Display summary
	fmt.Println("\n--- TOPOLOGY SUMMARY ---")
	fmt.Printf("Virtual Networks: %d\n", report.Summary.TotalVNets)
	fmt.Printf("Subnets: %d\n", report.Summary.TotalSubnets)
	fmt.Printf("Network Security Groups: %d\n", report.Summary.TotalNSGs)
	fmt.Printf("Security Rules: %d\n", report.Summary.TotalSecurityRules)
	fmt.Printf("Route Tables: %d\n", report.Summary.TotalRouteTables)
	fmt.Printf("Routes: %d\n", report.Summary.TotalRoutes)
	fmt.Printf("Private Endpoints: %d\n", report.Summary.TotalPrivateEndpoints)
	fmt.Printf("NAT Gateways: %d\n", report.Summary.TotalNATGateways)
	fmt.Printf("VPN Gateways: %d\n", report.Summary.TotalVPNGateways)
	fmt.Printf("ExpressRoute Circuits: %d\n", report.Summary.TotalERCircuits)
	fmt.Printf("Load Balancers: %d\n", report.Summary.TotalLoadBalancers)
	fmt.Printf("Application Gateways: %d\n", report.Summary.TotalAppGateways)
	if len(report.Summary.TotalIPAddressSpace) > 0 {
		fmt.Printf("Address Spaces: %v\n", report.Summary.TotalIPAddressSpace)
	}
	fmt.Printf("VNet Peerings: %d\n", report.Summary.VNetPeeringCount)

	// Display security findings
	if len(report.SecurityFindings) > 0 {
		fmt.Println("\n--- SECURITY FINDINGS ---")

		// Group by severity
		critical := []analyzer.SecurityFinding{}
		high := []analyzer.SecurityFinding{}
		medium := []analyzer.SecurityFinding{}
		low := []analyzer.SecurityFinding{}
		info := []analyzer.SecurityFinding{}

		for _, finding := range report.SecurityFindings {
			switch finding.Severity {
			case analyzer.SeverityCritical:
				critical = append(critical, finding)
			case analyzer.SeverityHigh:
				high = append(high, finding)
			case analyzer.SeverityMedium:
				medium = append(medium, finding)
			case analyzer.SeverityLow:
				low = append(low, finding)
			case analyzer.SeverityInfo:
				info = append(info, finding)
			}
		}

		fmt.Printf("\nTotal: %d findings\n", len(report.SecurityFindings))
		fmt.Printf("  Critical: %d | High: %d | Medium: %d | Low: %d | Info: %d\n",
			len(critical), len(high), len(medium), len(low), len(info))

		// Display critical and high findings in detail
		if len(critical) > 0 {
			fmt.Println("\n[CRITICAL]")
			for _, f := range critical {
				fmt.Printf("  * %s\n", f.Description)
				fmt.Printf("    Resource: %s", f.Resource)
				if f.Rule != "" {
					fmt.Printf(" | Rule: %s", f.Rule)
				}
				fmt.Println()
				fmt.Printf("    Recommendation: %s\n", f.Recommendation)
			}
		}

		if len(high) > 0 {
			fmt.Println("\n[HIGH]")
			for _, f := range high {
				fmt.Printf("  * %s\n", f.Description)
				fmt.Printf("    Resource: %s", f.Resource)
				if f.Rule != "" {
					fmt.Printf(" | Rule: %s", f.Rule)
				}
				fmt.Println()
				fmt.Printf("    Recommendation: %s\n", f.Recommendation)
			}
		}

		if len(medium) > 0 {
			fmt.Println("\n[MEDIUM]")
			for _, f := range medium {
				fmt.Printf("  * %s\n", f.Description)
				fmt.Printf("    Resource: %s\n", f.Resource)
			}
		}

		if len(low) > 0 {
			fmt.Println("\n[LOW]")
			for _, f := range low {
				fmt.Printf("  * %s\n", f.Description)
			}
		}
	} else {
		fmt.Println("\n--- SECURITY FINDINGS ---")
		fmt.Println("No security issues found!")
	}

	// Display orphaned resources
	hasOrphaned := len(report.OrphanedResources.UnattachedNSGs) > 0 ||
		len(report.OrphanedResources.UnusedRouteTables) > 0 ||
		len(report.OrphanedResources.UnusedNATGateways) > 0 ||
		len(report.OrphanedResources.IsolatedSubnets) > 0

	if hasOrphaned {
		fmt.Println("\n--- ORPHANED/UNUSED RESOURCES ---")
		if len(report.OrphanedResources.UnattachedNSGs) > 0 {
			fmt.Printf("Unattached NSGs: %v\n", report.OrphanedResources.UnattachedNSGs)
		}
		if len(report.OrphanedResources.UnusedRouteTables) > 0 {
			fmt.Printf("Unused Route Tables: %v\n", report.OrphanedResources.UnusedRouteTables)
		}
		if len(report.OrphanedResources.UnusedNATGateways) > 0 {
			fmt.Printf("Unused NAT Gateways: %v\n", report.OrphanedResources.UnusedNATGateways)
		}
		if len(report.OrphanedResources.IsolatedSubnets) > 0 {
			fmt.Printf("Subnets without NSG: %v\n", report.OrphanedResources.IsolatedSubnets)
		}
	}

	// Display recommendations
	if len(report.Recommendations) > 0 {
		fmt.Println("\n--- RECOMMENDATIONS ---")
		for i, rec := range report.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec)
		}
	}

	fmt.Println()
}
