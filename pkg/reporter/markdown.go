package reporter

import (
	"fmt"
	"strings"
	"time"

	"azure-network-analyzer/pkg/analyzer"
	"azure-network-analyzer/pkg/models"
)

// GenerateMarkdown creates a comprehensive Markdown report
func GenerateMarkdown(topology *models.NetworkTopology, analysis *analyzer.AnalysisReport) string {
	var md strings.Builder

	// Header
	md.WriteString("# Azure Network Topology Report\n\n")
	md.WriteString(fmt.Sprintf("**Subscription:** %s  \n", topology.SubscriptionID))
	md.WriteString(fmt.Sprintf("**Resource Group:** %s  \n", topology.ResourceGroup))
	md.WriteString(fmt.Sprintf("**Generated:** %s  \n\n", time.Now().Format("2006-01-02 15:04:05 MST")))

	// Executive Summary
	md.WriteString("## Executive Summary\n\n")
	md.WriteString(fmt.Sprintf("- **Total VNets:** %d\n", analysis.Summary.TotalVNets))
	md.WriteString(fmt.Sprintf("- **Total Subnets:** %d\n", analysis.Summary.TotalSubnets))
	md.WriteString(fmt.Sprintf("- **Total NSGs:** %d\n", analysis.Summary.TotalNSGs))
	md.WriteString(fmt.Sprintf("- **Total Security Rules:** %d\n", analysis.Summary.TotalSecurityRules))
	md.WriteString(fmt.Sprintf("- **VNet Peerings:** %d\n", analysis.Summary.VNetPeeringCount))
	md.WriteString(fmt.Sprintf("- **Security Findings:** %d\n", len(analysis.SecurityFindings)))

	// Count by severity
	critical, high, medium, low := countBySeverity(analysis.SecurityFindings)
	if critical > 0 || high > 0 {
		md.WriteString(fmt.Sprintf("  - Critical: %d, High: %d, Medium: %d, Low: %d\n", critical, high, medium, low))
	}
	md.WriteString("\n")

	// Security Findings Section
	if len(analysis.SecurityFindings) > 0 {
		md.WriteString("## Security Findings\n\n")

		if critical > 0 {
			md.WriteString("### Critical Issues\n\n")
			for _, f := range analysis.SecurityFindings {
				if f.Severity == analyzer.SeverityCritical {
					md.WriteString(fmt.Sprintf("#### %s\n", f.Description))
					md.WriteString(fmt.Sprintf("- **Resource:** %s\n", f.Resource))
					if f.Rule != "" {
						md.WriteString(fmt.Sprintf("- **Rule:** %s\n", f.Rule))
					}
					md.WriteString(fmt.Sprintf("- **Category:** %s\n", f.Category))
					md.WriteString(fmt.Sprintf("- **Recommendation:** %s\n\n", f.Recommendation))
				}
			}
		}

		if high > 0 {
			md.WriteString("### High Severity Issues\n\n")
			for _, f := range analysis.SecurityFindings {
				if f.Severity == analyzer.SeverityHigh {
					md.WriteString(fmt.Sprintf("#### %s\n", f.Description))
					md.WriteString(fmt.Sprintf("- **Resource:** %s\n", f.Resource))
					if f.Rule != "" {
						md.WriteString(fmt.Sprintf("- **Rule:** %s\n", f.Rule))
					}
					md.WriteString(fmt.Sprintf("- **Recommendation:** %s\n\n", f.Recommendation))
				}
			}
		}

		if medium > 0 {
			md.WriteString("### Medium Severity Issues\n\n")
			for _, f := range analysis.SecurityFindings {
				if f.Severity == analyzer.SeverityMedium {
					md.WriteString(fmt.Sprintf("- %s (%s)\n", f.Description, f.Resource))
				}
			}
			md.WriteString("\n")
		}

		if low > 0 {
			md.WriteString("### Low Severity Issues\n\n")
			for _, f := range analysis.SecurityFindings {
				if f.Severity == analyzer.SeverityLow {
					md.WriteString(fmt.Sprintf("- %s\n", f.Description))
				}
			}
			md.WriteString("\n")
		}
	}

	// Recommendations
	if len(analysis.Recommendations) > 0 {
		md.WriteString("## Recommendations\n\n")
		for i, rec := range analysis.Recommendations {
			md.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec))
		}
		md.WriteString("\n")
	}

	// Network Topology Details
	md.WriteString("## Network Topology\n\n")

	// Virtual Networks
	if len(topology.VirtualNetworks) > 0 {
		md.WriteString("### Virtual Networks\n\n")
		for _, vnet := range topology.VirtualNetworks {
			md.WriteString(fmt.Sprintf("#### %s\n", vnet.Name))
			md.WriteString(fmt.Sprintf("- **Location:** %s\n", vnet.Location))
			md.WriteString(fmt.Sprintf("- **Address Space:** %s\n", strings.Join(vnet.AddressSpace, ", ")))
			if len(vnet.DNSServers) > 0 {
				md.WriteString(fmt.Sprintf("- **DNS Servers:** %s\n", strings.Join(vnet.DNSServers, ", ")))
			}
			md.WriteString(fmt.Sprintf("- **DDoS Protection:** %v\n", vnet.EnableDDoS))

			if len(vnet.Subnets) > 0 {
				md.WriteString("\n**Subnets:**\n\n")
				md.WriteString("| Name | Address Prefix | NSG | Route Table | NAT Gateway |\n")
				md.WriteString("|------|----------------|-----|-------------|-------------|\n")
				for _, subnet := range vnet.Subnets {
					nsg := "-"
					if subnet.NetworkSecurityGroup != nil {
						nsg = extractName(*subnet.NetworkSecurityGroup)
					}
					rt := "-"
					if subnet.RouteTable != nil {
						rt = extractName(*subnet.RouteTable)
					}
					nat := "-"
					if subnet.NATGateway != nil {
						nat = extractName(*subnet.NATGateway)
					}
					md.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
						subnet.Name, subnet.AddressPrefix, nsg, rt, nat))
				}
			}

			if len(vnet.Peerings) > 0 {
				md.WriteString("\n**Peerings:**\n\n")
				for _, peer := range vnet.Peerings {
					md.WriteString(fmt.Sprintf("- %s → %s (State: %s)\n",
						peer.Name, peer.RemoteVNetName, peer.PeeringState))
				}
			}
			md.WriteString("\n")
		}
	}

	// Network Security Groups
	if len(topology.NSGs) > 0 {
		md.WriteString("### Network Security Groups\n\n")
		for _, nsg := range topology.NSGs {
			md.WriteString(fmt.Sprintf("#### %s\n", nsg.Name))
			md.WriteString(fmt.Sprintf("- **Location:** %s\n", nsg.Location))

			if len(nsg.SecurityRules) > 0 {
				md.WriteString("\n**Security Rules:**\n\n")
				md.WriteString("| Priority | Name | Direction | Access | Protocol | Source | Dest Port |\n")
				md.WriteString("|----------|------|-----------|--------|----------|--------|------------|\n")
				for _, rule := range nsg.SecurityRules {
					src := rule.SourceAddressPrefix
					if len(src) > 20 {
						src = src[:17] + "..."
					}
					md.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s | %s |\n",
						rule.Priority, rule.Name, rule.Direction, rule.Access,
						rule.Protocol, src, rule.DestinationPortRange))
				}
			}
			md.WriteString("\n")
		}
	}

	// Route Tables
	if len(topology.RouteTables) > 0 {
		md.WriteString("### Route Tables\n\n")
		for _, rt := range topology.RouteTables {
			md.WriteString(fmt.Sprintf("#### %s\n", rt.Name))
			md.WriteString(fmt.Sprintf("- **Location:** %s\n", rt.Location))
			md.WriteString(fmt.Sprintf("- **Disable BGP Propagation:** %v\n", rt.DisableBGPRoutePropagation))

			if len(rt.Routes) > 0 {
				md.WriteString("\n**Routes:**\n\n")
				md.WriteString("| Name | Address Prefix | Next Hop Type | Next Hop IP |\n")
				md.WriteString("|------|----------------|---------------|-------------|\n")
				for _, route := range rt.Routes {
					nextHopIP := route.NextHopIPAddress
					if nextHopIP == "" {
						nextHopIP = "-"
					}
					md.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
						route.Name, route.AddressPrefix, route.NextHopType, nextHopIP))
				}
			}
			md.WriteString("\n")
		}
	}

	// Private Endpoints
	if len(topology.PrivateEndpoints) > 0 {
		md.WriteString("### Private Endpoints\n\n")
		md.WriteString("| Name | Location | Target Resource | Subnet |\n")
		md.WriteString("|------|----------|-----------------|--------|\n")
		for _, pe := range topology.PrivateEndpoints {
			target := extractName(pe.PrivateLinkServiceID)
			subnet := extractName(pe.SubnetID)
			md.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				pe.Name, pe.Location, target, subnet))
		}
		md.WriteString("\n")
	}

	// Load Balancers
	if len(topology.LoadBalancers) > 0 {
		md.WriteString("### Load Balancers\n\n")
		for _, lb := range topology.LoadBalancers {
			md.WriteString(fmt.Sprintf("#### %s\n", lb.Name))
			md.WriteString(fmt.Sprintf("- **SKU:** %s\n", lb.SKU))
			md.WriteString(fmt.Sprintf("- **Frontend IPs:** %d\n", len(lb.FrontendIPConfigs)))
			md.WriteString(fmt.Sprintf("- **Backend Pools:** %d\n", len(lb.BackendAddressPools)))
			md.WriteString(fmt.Sprintf("- **Load Balancing Rules:** %d\n", len(lb.LoadBalancingRules)))
			md.WriteString("\n")
		}
	}

	// Application Gateways
	if len(topology.AppGateways) > 0 {
		md.WriteString("### Application Gateways\n\n")
		for _, appgw := range topology.AppGateways {
			md.WriteString(fmt.Sprintf("#### %s\n", appgw.Name))
			md.WriteString(fmt.Sprintf("- **SKU:** %s (Capacity: %d)\n", appgw.SKU, appgw.Capacity))
			md.WriteString(fmt.Sprintf("- **WAF Enabled:** %v\n", appgw.WAFEnabled))
			md.WriteString(fmt.Sprintf("- **HTTP Listeners:** %d\n", len(appgw.HTTPListeners)))
			md.WriteString(fmt.Sprintf("- **Backend Pools:** %d\n", len(appgw.BackendAddressPools)))
			md.WriteString("\n")
		}
	}

	// VPN Gateways
	if len(topology.VPNGateways) > 0 {
		md.WriteString("### VPN Gateways\n\n")
		for _, vpn := range topology.VPNGateways {
			md.WriteString(fmt.Sprintf("#### %s\n", vpn.Name))
			md.WriteString(fmt.Sprintf("- **Type:** %s\n", vpn.GatewayType))
			md.WriteString(fmt.Sprintf("- **VPN Type:** %s\n", vpn.VpnType))
			md.WriteString(fmt.Sprintf("- **SKU:** %s\n", vpn.SKU))
			if vpn.BGPSettings != nil {
				md.WriteString(fmt.Sprintf("- **BGP ASN:** %d\n", vpn.BGPSettings.ASN))
			}
			md.WriteString("\n")
		}
	}

	// Orphaned Resources
	hasOrphaned := len(analysis.OrphanedResources.UnattachedNSGs) > 0 ||
		len(analysis.OrphanedResources.UnusedRouteTables) > 0 ||
		len(analysis.OrphanedResources.IsolatedSubnets) > 0

	if hasOrphaned {
		md.WriteString("## Orphaned/Unused Resources\n\n")
		if len(analysis.OrphanedResources.UnattachedNSGs) > 0 {
			md.WriteString("### Unattached NSGs\n")
			for _, nsg := range analysis.OrphanedResources.UnattachedNSGs {
				md.WriteString(fmt.Sprintf("- %s\n", nsg))
			}
			md.WriteString("\n")
		}
		if len(analysis.OrphanedResources.UnusedRouteTables) > 0 {
			md.WriteString("### Unused Route Tables\n")
			for _, rt := range analysis.OrphanedResources.UnusedRouteTables {
				md.WriteString(fmt.Sprintf("- %s\n", rt))
			}
			md.WriteString("\n")
		}
		if len(analysis.OrphanedResources.IsolatedSubnets) > 0 {
			md.WriteString("### Subnets Without NSG\n")
			for _, subnet := range analysis.OrphanedResources.IsolatedSubnets {
				md.WriteString(fmt.Sprintf("- %s\n", subnet))
			}
			md.WriteString("\n")
		}
	}

	// Footer
	md.WriteString("---\n")
	md.WriteString("*Generated by Azure Network Topology Analyzer v1.0.0*\n")

	return md.String()
}

// Helper functions

func countBySeverity(findings []analyzer.SecurityFinding) (critical, high, medium, low int) {
	for _, f := range findings {
		switch f.Severity {
		case analyzer.SeverityCritical:
			critical++
		case analyzer.SeverityHigh:
			high++
		case analyzer.SeverityMedium:
			medium++
		case analyzer.SeverityLow:
			low++
		}
	}
	return
}

func extractName(resourceID string) string {
	if resourceID == "" {
		return ""
	}
	parts := strings.Split(resourceID, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return resourceID
}
