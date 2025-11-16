package analyzer

// AnalysisReport contains the results of topology and security analysis
type AnalysisReport struct {
	Summary          TopologySummary    `json:"summary"`
	SecurityFindings []SecurityFinding  `json:"security_findings"`
	OrphanedResources OrphanedResources `json:"orphaned_resources"`
	Recommendations  []string           `json:"recommendations"`
}

// TopologySummary provides statistics about the network topology
type TopologySummary struct {
	TotalVNets            int      `json:"total_vnets"`
	TotalSubnets          int      `json:"total_subnets"`
	TotalNSGs             int      `json:"total_nsgs"`
	TotalSecurityRules    int      `json:"total_security_rules"`
	TotalRouteTables      int      `json:"total_route_tables"`
	TotalRoutes           int      `json:"total_routes"`
	TotalPrivateEndpoints int      `json:"total_private_endpoints"`
	TotalPrivateDNSZones  int      `json:"total_private_dns_zones"`
	TotalNATGateways      int      `json:"total_nat_gateways"`
	TotalVPNGateways      int      `json:"total_vpn_gateways"`
	TotalERCircuits       int      `json:"total_er_circuits"`
	TotalLoadBalancers    int      `json:"total_load_balancers"`
	TotalAppGateways      int      `json:"total_app_gateways"`
	TotalIPAddressSpace   []string `json:"total_ip_address_space"`
	VNetPeeringCount      int      `json:"vnet_peering_count"`
	CrossRGDependencies   int      `json:"cross_rg_dependencies"`
}

// SecurityFinding represents a potential security issue
type SecurityFinding struct {
	Severity       string `json:"severity"`        // Critical, High, Medium, Low, Info
	Category       string `json:"category"`        // e.g., "NSG Rule", "Network Exposure"
	Resource       string `json:"resource"`        // Resource name (e.g., NSG name)
	ResourceID     string `json:"resource_id"`     // Full resource ID
	Rule           string `json:"rule"`            // Rule name if applicable
	Description    string `json:"description"`     // What the issue is
	Recommendation string `json:"recommendation"`  // How to fix it
}

// OrphanedResources contains resources that are not attached or used
type OrphanedResources struct {
	UnattachedNSGs       []string `json:"unattached_nsgs"`
	UnusedRouteTables    []string `json:"unused_route_tables"`
	UnusedNATGateways    []string `json:"unused_nat_gateways"`
	IsolatedSubnets      []string `json:"isolated_subnets"`       // Subnets with no NSG
	SubnetsWithoutRoutes []string `json:"subnets_without_routes"` // Subnets with no route table
}

// Severity levels
const (
	SeverityCritical = "Critical"
	SeverityHigh     = "High"
	SeverityMedium   = "Medium"
	SeverityLow      = "Low"
	SeverityInfo     = "Info"
)

// Security finding categories
const (
	CategoryNSGRule          = "NSG Rule"
	CategoryNetworkExposure  = "Network Exposure"
	CategoryMissingProtection = "Missing Protection"
	CategoryConfiguration    = "Configuration"
)
