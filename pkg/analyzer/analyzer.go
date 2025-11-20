package analyzer

import (
	"azure-network-analyzer/pkg/models"
)

// Analyze performs comprehensive analysis on the network topology
func Analyze(topology *models.NetworkTopology) *AnalysisReport {
	report := &AnalysisReport{
		Summary:           generateSummary(topology),
		SecurityFindings:  AnalyzeSecurityRisks(topology),
		OrphanedResources: findOrphanedResources(topology),
		Recommendations:   []string{},
	}

	// Generate high-level recommendations based on findings
	report.Recommendations = generateRecommendations(report)

	return report
}

// generateSummary creates statistics about the topology
func generateSummary(topology *models.NetworkTopology) TopologySummary {
	summary := TopologySummary{
		TotalVNets:            len(topology.VirtualNetworks),
		TotalNSGs:             len(topology.NSGs),
		TotalRouteTables:      len(topology.RouteTables),
		TotalPrivateEndpoints: len(topology.PrivateEndpoints),
		TotalPrivateDNSZones:  len(topology.PrivateDNSZones),
		TotalNATGateways:      len(topology.NATGateways),
		TotalVPNGateways:      len(topology.VPNGateways),
		TotalERCircuits:       len(topology.ERCircuits),
		TotalLoadBalancers:    len(topology.LoadBalancers),
		TotalAppGateways:      len(topology.AppGateways),
		TotalAzureFirewalls:   len(topology.AzureFirewalls),
		TotalIPAddressSpace:   []string{},
	}

	// Count subnets and collect address spaces
	for _, vnet := range topology.VirtualNetworks {
		summary.TotalSubnets += len(vnet.Subnets)
		summary.TotalIPAddressSpace = append(summary.TotalIPAddressSpace, vnet.AddressSpace...)
		summary.VNetPeeringCount += len(vnet.Peerings)
	}

	// Count security rules
	for _, nsg := range topology.NSGs {
		summary.TotalSecurityRules += len(nsg.SecurityRules)
	}

	// Count routes
	for _, rt := range topology.RouteTables {
		summary.TotalRoutes += len(rt.Routes)
	}

	// Count cross-RG dependencies (simplified: count peerings to different VNets)
	summary.CrossRGDependencies = summary.VNetPeeringCount / 2 // Each peering is counted twice

	return summary
}

// findOrphanedResources identifies resources that are not attached or used
func findOrphanedResources(topology *models.NetworkTopology) OrphanedResources {
	orphaned := OrphanedResources{
		UnattachedNSGs:       []string{},
		UnusedRouteTables:    []string{},
		UnusedNATGateways:    []string{},
		IsolatedSubnets:      []string{},
		SubnetsWithoutRoutes: []string{},
	}

	// Build maps of what's used
	usedNSGs := make(map[string]bool)
	usedRouteTables := make(map[string]bool)
	usedNATGateways := make(map[string]bool)

	// Check each subnet for associations
	for _, vnet := range topology.VirtualNetworks {
		for _, subnet := range vnet.Subnets {
			// Track NSG usage
			if subnet.NetworkSecurityGroup != nil {
				usedNSGs[*subnet.NetworkSecurityGroup] = true
			} else {
				// Subnet without NSG (potential security concern)
				orphaned.IsolatedSubnets = append(orphaned.IsolatedSubnets,
					vnet.Name+"/"+subnet.Name)
			}

			// Track Route Table usage
			if subnet.RouteTable != nil {
				usedRouteTables[*subnet.RouteTable] = true
			} else {
				orphaned.SubnetsWithoutRoutes = append(orphaned.SubnetsWithoutRoutes,
					vnet.Name+"/"+subnet.Name)
			}

			// Track NAT Gateway usage
			if subnet.NATGateway != nil {
				usedNATGateways[*subnet.NATGateway] = true
			}
		}
	}

	// Find unattached NSGs
	for _, nsg := range topology.NSGs {
		if !usedNSGs[nsg.ID] {
			orphaned.UnattachedNSGs = append(orphaned.UnattachedNSGs, nsg.Name)
		}
	}

	// Find unused Route Tables
	for _, rt := range topology.RouteTables {
		if !usedRouteTables[rt.ID] {
			orphaned.UnusedRouteTables = append(orphaned.UnusedRouteTables, rt.Name)
		}
	}

	// Find unused NAT Gateways
	for _, nat := range topology.NATGateways {
		if !usedNATGateways[nat.ID] {
			orphaned.UnusedNATGateways = append(orphaned.UnusedNATGateways, nat.Name)
		}
	}

	return orphaned
}

// generateRecommendations creates actionable recommendations based on findings
func generateRecommendations(report *AnalysisReport) []string {
	recommendations := []string{}

	// Count findings by severity
	criticalCount := 0
	highCount := 0
	for _, finding := range report.SecurityFindings {
		switch finding.Severity {
		case SeverityCritical:
			criticalCount++
		case SeverityHigh:
			highCount++
		}
	}

	// Add recommendations based on findings
	if criticalCount > 0 {
		recommendations = append(recommendations,
			"URGENT: Address critical security findings immediately to prevent potential breaches")
	}

	if highCount > 0 {
		recommendations = append(recommendations,
			"Review and remediate high-severity security findings within 24-48 hours")
	}

	// Check for orphaned resources
	if len(report.OrphanedResources.UnattachedNSGs) > 0 {
		recommendations = append(recommendations,
			"Consider removing unattached NSGs or attaching them to appropriate subnets")
	}

	if len(report.OrphanedResources.IsolatedSubnets) > 0 {
		recommendations = append(recommendations,
			"Attach NSGs to isolated subnets to improve security posture")
	}

	if len(report.OrphanedResources.UnusedRouteTables) > 0 {
		recommendations = append(recommendations,
			"Remove unused Route Tables to reduce configuration complexity")
	}

	// General recommendations
	if report.Summary.TotalVNets > 0 && report.Summary.VNetPeeringCount == 0 {
		recommendations = append(recommendations,
			"Consider VNet peering for connectivity between virtual networks if needed")
	}

	if report.Summary.TotalPrivateEndpoints == 0 && report.Summary.TotalVNets > 0 {
		recommendations = append(recommendations,
			"Consider using Private Endpoints for secure access to Azure PaaS services")
	}

	return recommendations
}
