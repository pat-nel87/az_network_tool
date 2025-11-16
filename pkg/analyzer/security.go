package analyzer

import (
	"fmt"
	"strings"

	"azure-network-analyzer/pkg/models"
)

// AnalyzeSecurityRisks performs security analysis on the entire topology
func AnalyzeSecurityRisks(topology *models.NetworkTopology) []SecurityFinding {
	findings := []SecurityFinding{}

	// Analyze NSG rules
	findings = append(findings, analyzeNSGRules(topology.NSGs)...)

	// Analyze subnet security
	findings = append(findings, analyzeSubnetSecurity(topology.VirtualNetworks)...)

	// Analyze gateway configurations
	findings = append(findings, analyzeGatewaySecurity(topology)...)

	return findings
}

// analyzeNSGRules checks NSG rules for security risks
func analyzeNSGRules(nsgs []models.NetworkSecurityGroup) []SecurityFinding {
	findings := []SecurityFinding{}

	for _, nsg := range nsgs {
		for _, rule := range nsg.SecurityRules {
			// Skip deny rules - they're generally good
			if rule.Access != "Allow" {
				continue
			}

			// Check for internet-exposed sensitive ports
			if isInternetSource(rule.SourceAddressPrefix) {
				findings = append(findings, checkSensitivePorts(nsg, rule)...)
			}

			// Check for overly permissive rules
			findings = append(findings, checkOverlyPermissive(nsg, rule)...)

			// Check for missing descriptions
			if rule.Description == "" {
				findings = append(findings, SecurityFinding{
					Severity:       SeverityLow,
					Category:       CategoryConfiguration,
					Resource:       nsg.Name,
					ResourceID:     nsg.ID,
					Rule:           rule.Name,
					Description:    fmt.Sprintf("Security rule '%s' has no description", rule.Name),
					Recommendation: "Add descriptive comments to all security rules for better maintainability",
				})
			}

			// Check for high priority allow rules (might override important deny rules)
			if rule.Priority < 200 && isWideOpen(rule) {
				findings = append(findings, SecurityFinding{
					Severity:       SeverityMedium,
					Category:       CategoryNSGRule,
					Resource:       nsg.Name,
					ResourceID:     nsg.ID,
					Rule:           rule.Name,
					Description:    fmt.Sprintf("High priority (%d) allow rule may override important deny rules", rule.Priority),
					Recommendation: "Review rule priority to ensure deny rules are not inadvertently bypassed",
				})
			}
		}
	}

	return findings
}

// checkSensitivePorts checks if sensitive ports are exposed to the internet
func checkSensitivePorts(nsg models.NetworkSecurityGroup, rule models.SecurityRule) []SecurityFinding {
	findings := []SecurityFinding{}

	// Map of sensitive ports and their risks
	sensitivePorts := map[string]struct {
		name     string
		severity string
	}{
		"22":    {"SSH", SeverityCritical},
		"3389":  {"RDP", SeverityCritical},
		"23":    {"Telnet", SeverityCritical},
		"21":    {"FTP", SeverityHigh},
		"445":   {"SMB", SeverityCritical},
		"1433":  {"SQL Server", SeverityCritical},
		"3306":  {"MySQL", SeverityCritical},
		"5432":  {"PostgreSQL", SeverityCritical},
		"27017": {"MongoDB", SeverityCritical},
		"6379":  {"Redis", SeverityHigh},
		"9200":  {"Elasticsearch", SeverityHigh},
	}

	// Check destination port range
	ports := parsePortRange(rule.DestinationPortRange)

	for _, port := range ports {
		if info, found := sensitivePorts[port]; found {
			findings = append(findings, SecurityFinding{
				Severity:   info.severity,
				Category:   CategoryNetworkExposure,
				Resource:   nsg.Name,
				ResourceID: nsg.ID,
				Rule:       rule.Name,
				Description: fmt.Sprintf("%s (port %s) is exposed to the internet via rule '%s'",
					info.name, port, rule.Name),
				Recommendation: fmt.Sprintf("Restrict %s access to specific IP addresses or use Azure Bastion/VPN for remote access", info.name),
			})
		}
	}

	// Check for all ports open
	if rule.DestinationPortRange == "*" || rule.DestinationPortRange == "0-65535" {
		findings = append(findings, SecurityFinding{
			Severity:       SeverityCritical,
			Category:       CategoryNetworkExposure,
			Resource:       nsg.Name,
			ResourceID:     nsg.ID,
			Rule:           rule.Name,
			Description:    fmt.Sprintf("All ports are exposed to the internet via rule '%s'", rule.Name),
			Recommendation: "Restrict to specific ports required for your application",
		})
	}

	return findings
}

// checkOverlyPermissive checks for overly permissive rules
func checkOverlyPermissive(nsg models.NetworkSecurityGroup, rule models.SecurityRule) []SecurityFinding {
	findings := []SecurityFinding{}

	// Check for any-to-any rules
	if isWideOpen(rule) {
		findings = append(findings, SecurityFinding{
			Severity:       SeverityHigh,
			Category:       CategoryNSGRule,
			Resource:       nsg.Name,
			ResourceID:     nsg.ID,
			Rule:           rule.Name,
			Description:    fmt.Sprintf("Rule '%s' allows traffic from any source to any destination on all ports", rule.Name),
			Recommendation: "Implement least-privilege access by restricting source, destination, and ports",
		})
	}

	// Check for wide port ranges
	if isWidePortRange(rule.DestinationPortRange) {
		findings = append(findings, SecurityFinding{
			Severity:       SeverityMedium,
			Category:       CategoryNSGRule,
			Resource:       nsg.Name,
			ResourceID:     nsg.ID,
			Rule:           rule.Name,
			Description:    fmt.Sprintf("Rule '%s' allows a wide range of ports (%s)", rule.Name, rule.DestinationPortRange),
			Recommendation: "Restrict to specific ports required for your application",
		})
	}

	return findings
}

// analyzeSubnetSecurity checks subnet-level security
func analyzeSubnetSecurity(vnets []models.VirtualNetwork) []SecurityFinding {
	findings := []SecurityFinding{}

	for _, vnet := range vnets {
		for _, subnet := range vnet.Subnets {
			// Check for subnets without NSG
			if subnet.NetworkSecurityGroup == nil {
				findings = append(findings, SecurityFinding{
					Severity:       SeverityHigh,
					Category:       CategoryMissingProtection,
					Resource:       fmt.Sprintf("%s/%s", vnet.Name, subnet.Name),
					ResourceID:     subnet.ID,
					Rule:           "",
					Description:    fmt.Sprintf("Subnet '%s' in VNet '%s' has no Network Security Group attached", subnet.Name, vnet.Name),
					Recommendation: "Attach an NSG to control inbound and outbound traffic",
				})
			}

			// Check for large subnets (might indicate poor network segmentation)
			if isLargeSubnet(subnet.AddressPrefix) {
				findings = append(findings, SecurityFinding{
					Severity:       SeverityInfo,
					Category:       CategoryConfiguration,
					Resource:       fmt.Sprintf("%s/%s", vnet.Name, subnet.Name),
					ResourceID:     subnet.ID,
					Rule:           "",
					Description:    fmt.Sprintf("Subnet '%s' has a large address space (%s)", subnet.Name, subnet.AddressPrefix),
					Recommendation: "Consider smaller subnets for better network segmentation and security isolation",
				})
			}
		}
	}

	return findings
}

// analyzeGatewaySecurity checks gateway configurations
func analyzeGatewaySecurity(topology *models.NetworkTopology) []SecurityFinding {
	findings := []SecurityFinding{}

	// Check VPN Gateway configurations
	for _, vpn := range topology.VPNGateways {
		// Check for basic SKU (limited features)
		if strings.Contains(strings.ToLower(vpn.SKU), "basic") {
			findings = append(findings, SecurityFinding{
				Severity:       SeverityMedium,
				Category:       CategoryConfiguration,
				Resource:       vpn.Name,
				ResourceID:     vpn.ID,
				Rule:           "",
				Description:    fmt.Sprintf("VPN Gateway '%s' uses Basic SKU with limited security features", vpn.Name),
				Recommendation: "Consider upgrading to VpnGw1 or higher for better performance and security features",
			})
		}
	}

	// Check Application Gateway WAF status
	for _, appgw := range topology.AppGateways {
		if !appgw.WAFEnabled {
			findings = append(findings, SecurityFinding{
				Severity:       SeverityHigh,
				Category:       CategoryMissingProtection,
				Resource:       appgw.Name,
				ResourceID:     appgw.ID,
				Rule:           "",
				Description:    fmt.Sprintf("Application Gateway '%s' does not have WAF enabled", appgw.Name),
				Recommendation: "Enable Web Application Firewall (WAF) to protect against common web vulnerabilities",
			})
		}
	}

	return findings
}

// Helper functions

func isInternetSource(source string) bool {
	return source == "*" || source == "0.0.0.0/0" || source == "Internet" || source == "Any"
}

func isWideOpen(rule models.SecurityRule) bool {
	return isInternetSource(rule.SourceAddressPrefix) &&
		(rule.DestinationAddressPrefix == "*" || rule.DestinationAddressPrefix == "0.0.0.0/0") &&
		(rule.DestinationPortRange == "*" || rule.DestinationPortRange == "0-65535")
}

func isWidePortRange(portRange string) bool {
	if portRange == "*" || portRange == "0-65535" {
		return true
	}
	// Check for ranges like "1-1000" or "80-8080"
	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) == 2 {
			// Simple heuristic: if range spans more than 100 ports, it's wide
			// A more sophisticated check would parse the numbers
			return len(portRange) > 5 // rough check
		}
	}
	return false
}

func parsePortRange(portRange string) []string {
	if portRange == "*" {
		return []string{"*"}
	}
	// Handle comma-separated ports
	if strings.Contains(portRange, ",") {
		return strings.Split(portRange, ",")
	}
	// Handle single port
	if !strings.Contains(portRange, "-") {
		return []string{portRange}
	}
	// For ranges, return the range string itself
	return []string{portRange}
}

func isLargeSubnet(addressPrefix string) bool {
	// Check for /16 or larger subnets
	if strings.Contains(addressPrefix, "/") {
		parts := strings.Split(addressPrefix, "/")
		if len(parts) == 2 {
			// /16 or smaller CIDR (larger subnet)
			if parts[1] == "16" || parts[1] == "15" || parts[1] == "14" ||
				parts[1] == "13" || parts[1] == "12" || parts[1] == "11" ||
				parts[1] == "10" || parts[1] == "9" || parts[1] == "8" {
				return true
			}
		}
	}
	return false
}
