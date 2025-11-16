# Azure Network Topology Analyzer - Claude Code Instructions

## Project Overview
Build a comprehensive CLI tool in Go that analyzes Azure network topology and generates detailed reports and visualizations. The tool should provide deep insights into VNets, subnets, NSGs, private endpoints, DNS zones, NAT gateways, route tables, VPN gateways, ExpressRoute, load balancers, and application gateways.

## Technology Stack
- **Language**: Go (1.21+)
- **CLI Framework**: Cobra
- **Azure SDK**: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/*
- **Authentication**: github.com/Azure/azure-sdk-for-go/sdk/azidentity
- **Visualization**: github.com/goccy/go-graphviz (for DOT file generation and rendering)
- **Output Formats**: JSON, Markdown, HTML, Graphviz DOT

## Project Structure
```
azure-network-analyzer/
├── cmd/
│   └── root.go              # Cobra root command setup
│   └── analyze.go           # Main analyze command
├── pkg/
│   ├── azure/
│   │   ├── client.go        # Azure client initialization
│   │   ├── vnets.go         # VNet and subnet queries
│   │   ├── nsgs.go          # NSG queries and analysis
│   │   ├── privatelink.go   # Private endpoints and DNS zones
│   │   ├── routing.go       # Route tables and NAT gateways
│   │   ├── gateways.go      # VPN and ExpressRoute gateways
│   │   ├── loadbalancers.go # Load balancers and App Gateways
│   │   └── networkwatcher.go # Network Watcher insights
│   ├── models/
│   │   └── topology.go      # Data models for network topology
│   ├── analyzer/
│   │   ├── analyzer.go      # Core analysis logic
│   │   └── security.go      # Security analysis (NSG risk assessment)
│   ├── visualization/
│   │   ├── graphviz.go      # DOT file generation
│   │   └── renderer.go      # SVG/PNG rendering
│   └── reporter/
│       ├── json.go          # JSON output
│       ├── markdown.go      # Markdown report generation
│       └── html.go          # HTML report generation
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Core Data Models

### Main Topology Model (pkg/models/topology.go)
```go
type NetworkTopology struct {
    SubscriptionID   string
    ResourceGroup    string
    VirtualNetworks  []VirtualNetwork
    NSGs             []NetworkSecurityGroup
    PrivateEndpoints []PrivateEndpoint
    PrivateDNSZones  []PrivateDNSZone
    RouteTables      []RouteTable
    NATGateways      []NATGateway
    VPNGateways      []VPNGateway
    ERCircuits       []ExpressRouteCircuit
    LoadBalancers    []LoadBalancer
    AppGateways      []ApplicationGateway
    NetworkWatcher   *NetworkWatcherInsights
    Timestamp        time.Time
}

type VirtualNetwork struct {
    ID               string
    Name             string
    ResourceGroup    string
    Location         string
    AddressSpace     []string
    Subnets          []Subnet
    Peerings         []VNetPeering
    DNSServers       []string
    EnableDDoS       bool
}

type Subnet struct {
    ID                  string
    Name                string
    AddressPrefix       string
    NetworkSecurityGroup *string // NSG ID if associated
    RouteTable          *string // Route table ID if associated
    NATGateway          *string // NAT gateway ID if associated
    PrivateEndpoints    []string // List of private endpoint IDs
    ServiceEndpoints    []string
    Delegations         []string
}

type NetworkSecurityGroup struct {
    ID              string
    Name            string
    ResourceGroup   string
    Location        string
    SecurityRules   []SecurityRule
    Associations    NSGAssociations
}

type SecurityRule struct {
    Name                     string
    Priority                 int32
    Direction                string // Inbound/Outbound
    Access                   string // Allow/Deny
    Protocol                 string
    SourceAddressPrefix      string
    SourcePortRange          string
    DestinationAddressPrefix string
    DestinationPortRange     string
    Description              string
}

type NSGAssociations struct {
    Subnets          []string // Subnet IDs
    NetworkInterfaces []string // NIC IDs
}

type VNetPeering struct {
    ID                      string
    Name                    string
    RemoteVNetID            string
    RemoteVNetName          string
    PeeringState            string
    AllowVNetAccess         bool
    AllowForwardedTraffic   bool
    AllowGatewayTransit     bool
    UseRemoteGateways       bool
}

type PrivateEndpoint struct {
    ID                    string
    Name                  string
    ResourceGroup         string
    Location              string
    SubnetID              string
    PrivateIPAddress      string
    PrivateLinkServiceID  string
    ConnectionState       string
    GroupIDs              []string
}

type PrivateDNSZone struct {
    ID                 string
    Name               string
    ResourceGroup      string
    VNetLinks          []VNetLink
    RecordSets         int
}

type VNetLink struct {
    ID                      string
    VNetID                  string
    VNetName                string
    RegistrationEnabled     bool
}

type RouteTable struct {
    ID                        string
    Name                      string
    ResourceGroup             string
    Location                  string
    Routes                    []Route
    DisableBGPRoutePropagation bool
    AssociatedSubnets         []string
}

type Route struct {
    Name              string
    AddressPrefix     string
    NextHopType       string // VirtualNetworkGateway, VNetLocal, Internet, VirtualAppliance, None
    NextHopIPAddress  string
}

type NATGateway struct {
    ID                  string
    Name                string
    ResourceGroup       string
    Location            string
    PublicIPAddresses   []string
    IdleTimeoutMinutes  int32
    AssociatedSubnets   []string
}

type VPNGateway struct {
    ID                  string
    Name                string
    ResourceGroup       string
    Location            string
    VNetID              string
    GatewayType         string // Vpn or ExpressRoute
    VpnType             string // RouteBased or PolicyBased
    SKU                 string
    ActiveActive        bool
    BGPSettings         *BGPSettings
    Connections         []VPNConnection
}

type BGPSettings struct {
    ASN               int64
    BGPPeeringAddress string
    PeerWeight        int32
}

type VPNConnection struct {
    ID                     string
    Name                   string
    ConnectionType         string // IPsec, Vnet2Vnet, ExpressRoute
    ConnectionStatus       string
    SharedKey              bool // Whether a shared key is configured
    EnableBGP              bool
    RemoteEntityID         string
}

type ExpressRouteCircuit struct {
    ID                    string
    Name                  string
    ResourceGroup         string
    Location              string
    ServiceProviderName   string
    PeeringLocation       string
    BandwidthInMbps       int32
    SKUTier               string
    SKUFamily             string
    CircuitProvisioningState string
    Peerings              []ERPeering
    Authorizations        []ERAuthorization
}

type ERPeering struct {
    Name                  string
    PeeringType           string // AzurePrivatePeering, AzurePublicPeering, MicrosoftPeering
    State                 string
    AzureASN              int32
    PeerASN               int64
    PrimaryPeerAddressPrefix string
    SecondaryPeerAddressPrefix string
    VlanID                int32
}

type ERAuthorization struct {
    Name                     string
    AuthorizationKey         bool // Whether key exists
    AuthorizationUseStatus   string
}

type LoadBalancer struct {
    ID                    string
    Name                  string
    ResourceGroup         string
    Location              string
    SKU                   string
    Type                  string // Public or Internal
    FrontendIPConfigs     []FrontendIPConfig
    BackendAddressPools   []BackendAddressPool
    LoadBalancingRules    []LoadBalancingRule
    Probes                []Probe
    InboundNATRules       []InboundNATRule
}

type FrontendIPConfig struct {
    Name               string
    PrivateIPAddress   string
    PublicIPAddressID  string
    SubnetID           string
}

type BackendAddressPool struct {
    Name          string
    BackendIPConfigs []string // NIC IDs
}

type LoadBalancingRule struct {
    Name                  string
    Protocol              string
    FrontendPort          int32
    BackendPort           int32
    EnableFloatingIP      bool
    IdleTimeoutMinutes    int32
    LoadDistribution      string
}

type Probe struct {
    Name            string
    Protocol        string
    Port            int32
    IntervalInSeconds int32
    NumberOfProbes  int32
    RequestPath     string // For HTTP/HTTPS
}

type InboundNATRule struct {
    Name                string
    Protocol            string
    FrontendPort        int32
    BackendPort         int32
    EnableFloatingIP    bool
}

type ApplicationGateway struct {
    ID                    string
    Name                  string
    ResourceGroup         string
    Location              string
    SKU                   string
    Tier                  string
    Capacity              int32
    SubnetID              string
    FrontendIPConfigs     []AppGWFrontendIPConfig
    FrontendPorts         []AppGWFrontendPort
    BackendAddressPools   []AppGWBackendAddressPool
    BackendHTTPSettings   []AppGWBackendHTTPSettings
    HTTPListeners         []AppGWHTTPListener
    RequestRoutingRules   []AppGWRequestRoutingRule
    Probes                []AppGWProbe
    WAFEnabled            bool
    WAFMode               string
}

type AppGWFrontendIPConfig struct {
    Name               string
    PrivateIPAddress   string
    PublicIPAddressID  string
}

type AppGWFrontendPort struct {
    Name string
    Port int32
}

type AppGWBackendAddressPool struct {
    Name          string
    BackendAddresses []string // IP addresses or FQDNs
}

type AppGWBackendHTTPSettings struct {
    Name                  string
    Port                  int32
    Protocol              string
    CookieBasedAffinity   string
    RequestTimeout        int32
    ProbeName             string
}

type AppGWHTTPListener struct {
    Name                  string
    FrontendIPConfig      string
    FrontendPort          string
    Protocol              string
    HostName              string
}

type AppGWRequestRoutingRule struct {
    Name                string
    RuleType            string
    HTTPListener        string
    BackendAddressPool  string
    BackendHTTPSettings string
    Priority            int32
}

type AppGWProbe struct {
    Name                string
    Protocol            string
    Host                string
    Path                string
    Interval            int32
    Timeout             int32
    UnhealthyThreshold  int32
}

type NetworkWatcherInsights struct {
    FlowLogsEnabled      bool
    FlowLogs             []FlowLog
    ConnectionMonitors   []ConnectionMonitor
    PacketCaptures       []PacketCapture
}

type FlowLog struct {
    ID                string
    NSGId             string
    StorageAccountID  string
    Enabled           bool
    RetentionDays     int32
    TrafficAnalytics  bool
}

type ConnectionMonitor struct {
    Name              string
    Source            string
    Destination       string
    MonitoringStatus  string
}

type PacketCapture struct {
    Name              string
    Target            string
    Status            string
}
```

## Implementation Guide

### Phase 1: Project Setup and Azure Client Initialization

1. **Initialize Go module and install dependencies**
```bash
go mod init github.com/yourusername/azure-network-analyzer
go get github.com/spf13/cobra@latest
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest
go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork@latest
go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources@latest
go get github.com/goccy/go-graphviz@latest
```

2. **Create main.go**
```go
package main

import "github.com/yourusername/azure-network-analyzer/cmd"

func main() {
    cmd.Execute()
}
```

3. **Implement pkg/azure/client.go**
- Create AzureClient struct that holds credential and subscription ID
- Initialize using DefaultAzureCredential (supports Azure CLI, managed identity, etc.)
- Create factory methods for each resource manager client:
  - VirtualNetworksClient
  - SubnetsClient
  - SecurityGroupsClient
  - PrivateEndpointsClient
  - PrivateDNSZonesClient
  - RouteTablesClient
  - NatGatewaysClient
  - VirtualNetworkGatewaysClient
  - ExpressRouteCircuitsClient
  - LoadBalancersClient
  - ApplicationGatewaysClient
  - NetworkWatchersClient
  - FlowLogsClient

```go
type AzureClient struct {
    cred           azidentity.DefaultAzureCredential
    subscriptionID string
}

func NewAzureClient(subscriptionID string) (*AzureClient, error) {
    // Initialize credential
    // Return client
}
```

### Phase 2: Data Collection Layer

4. **Implement pkg/azure/vnets.go**
- `GetVirtualNetworks(ctx, resourceGroup)` - List all VNets in RG
- `GetSubnets(ctx, resourceGroup, vnetName)` - Get all subnets for a VNet
- `GetVNetPeerings(ctx, resourceGroup, vnetName)` - Get peering information
- Handle pagination properly
- Extract all relevant properties including address spaces, DNS servers, DDoS protection

5. **Implement pkg/azure/nsgs.go**
- `GetNetworkSecurityGroups(ctx, resourceGroup)` - List all NSGs
- `GetNSGSecurityRules(ctx, resourceGroup, nsgName)` - Get all rules (default + custom)
- `GetNSGAssociations(ctx, nsg)` - Find which subnets/NICs are associated
- Parse and structure security rules properly

6. **Implement pkg/azure/privatelink.go**
- `GetPrivateEndpoints(ctx, resourceGroup)` - List all private endpoints
- `GetPrivateDNSZones(ctx, resourceGroup)` - List all private DNS zones
- `GetPrivateDNSZoneVNetLinks(ctx, resourceGroup, zoneName)` - Get VNet links for each zone
- Extract connection states, group IDs, and linked resources

7. **Implement pkg/azure/routing.go**
- `GetRouteTables(ctx, resourceGroup)` - List all route tables
- `GetRoutes(ctx, resourceGroup, routeTableName)` - Get routes in each table
- `GetNATGateways(ctx, resourceGroup)` - List NAT gateways
- `GetNATGatewayPublicIPs(ctx, natGateway)` - Get associated public IPs
- Map associations back to subnets

8. **Implement pkg/azure/gateways.go**
- `GetVPNGateways(ctx, resourceGroup)` - List VPN gateways
- `GetVPNConnections(ctx, resourceGroup, gatewayName)` - Get connections for each gateway
- `GetExpressRouteCircuits(ctx, resourceGroup)` - List ER circuits
- `GetERPeerings(ctx, resourceGroup, circuitName)` - Get peerings for each circuit
- `GetERAuthorizations(ctx, resourceGroup, circuitName)` - Get authorizations
- Extract BGP settings, connection states, SKUs

9. **Implement pkg/azure/loadbalancers.go**
- `GetLoadBalancers(ctx, resourceGroup)` - List all load balancers
- `GetApplicationGateways(ctx, resourceGroup)` - List all app gateways
- Extract all frontend/backend configs, rules, probes
- For App Gateway, get WAF configuration if enabled

10. **Implement pkg/azure/networkwatcher.go**
- `GetNetworkWatcher(ctx, resourceGroup)` - Find Network Watcher instance
- `GetFlowLogs(ctx, resourceGroup)` - List NSG flow logs
- `GetConnectionMonitors(ctx, resourceGroup)` - List connection monitors
- `GetPacketCaptures(ctx, resourceGroup)` - List packet captures
- Handle case where Network Watcher might not be enabled

### Phase 3: Analysis and Enrichment

11. **Implement pkg/analyzer/analyzer.go**
- `Analyze(topology *NetworkTopology) *AnalysisReport` - Main analysis function
- Generate summary statistics:
  - Total VNets, subnets, NSGs
  - Total IP address space used vs available
  - Count of each resource type
  - Cross-RG dependencies (peerings, connections)
- Build adjacency information for visualization
- Identify orphaned resources (NSGs not attached, route tables not used, etc.)

12. **Implement pkg/analyzer/security.go**
- `AnalyzeNSGRules(nsgs []NetworkSecurityGroup) []SecurityFinding`
- Flag potentially risky rules:
  - Rules allowing 0.0.0.0/0 or * as source
  - Wide port ranges (e.g., 0-65535)
  - RDP/SSH exposed to internet
  - High priority allow rules that might override deny rules
  - Rules with no description
- Rank findings by severity (Critical, High, Medium, Low, Info)
- Generate recommendations

```go
type SecurityFinding struct {
    Severity     string
    Category     string
    Resource     string // NSG name
    Rule         string // Rule name
    Description  string
    Recommendation string
}
```

### Phase 4: Visualization

13. **Implement pkg/visualization/graphviz.go**
- `GenerateDOTFile(topology *NetworkTopology) string`
- Create hierarchical graph:
  - Use clusters for VNets
  - Nodes for subnets, gateways, load balancers, private endpoints
  - Edges for:
    - VNet peerings (dashed lines)
    - Gateway connections (bold lines)
    - Private endpoint connections (dotted lines)
    - Subnet to load balancer backend pool relationships
- Color coding:
  - VNets: Light blue clusters
  - Subnets: Green boxes
  - NSGs: Yellow shields (attached to subnets)
  - Gateways: Purple diamonds
  - Load Balancers: Orange circles
  - Private Endpoints: Pink dots
- Add labels with key information (CIDR, SKU, status)

```go
// Example structure
func GenerateDOTFile(topology *NetworkTopology) string {
    var dot strings.Builder
    dot.WriteString("digraph NetworkTopology {\n")
    dot.WriteString("  rankdir=TB;\n")
    dot.WriteString("  node [shape=box];\n")
    
    // Create clusters for each VNet
    for _, vnet := range topology.VirtualNetworks {
        dot.WriteString(fmt.Sprintf("  subgraph cluster_%s {\n", sanitize(vnet.Name)))
        dot.WriteString(fmt.Sprintf("    label=\"%s\\n%s\";\n", vnet.Name, strings.Join(vnet.AddressSpace, ", ")))
        dot.WriteString("    color=lightblue;\n")
        
        // Add subnets as nodes
        for _, subnet := range vnet.Subnets {
            // Add subnet node
            // Add NSG shield if associated
            // Add connections to NAT gateway, route table
        }
        
        dot.WriteString("  }\n")
    }
    
    // Add peering edges
    // Add gateway connections
    // Add load balancer connections
    
    dot.WriteString("}\n")
    return dot.String()
}
```

14. **Implement pkg/visualization/renderer.go**
- `RenderSVG(dotContent string, outputPath string) error`
- `RenderPNG(dotContent string, outputPath string) error`
- Use go-graphviz to render DOT to image formats
- Handle errors gracefully (e.g., if Graphviz not installed)

### Phase 5: Reporting

15. **Implement pkg/reporter/json.go**
- `GenerateJSON(topology *NetworkTopology) ([]byte, error)`
- Marshal the entire topology structure to JSON
- Pretty print with proper indentation
- Include timestamp and metadata

16. **Implement pkg/reporter/markdown.go**
- `GenerateMarkdown(topology *NetworkTopology, analysis *AnalysisReport) string`
- Create hierarchical structure:

```markdown
# Azure Network Topology Report
**Resource Group:** <name>
**Subscription:** <id>
**Generated:** <timestamp>

## Executive Summary
- Total VNets: X
- Total Subnets: Y
- Total NSGs: Z
- Total Private Endpoints: A
...

## Virtual Networks

### VNet: <name>
- **Location:** <location>
- **Address Space:** <CIDR ranges>
- **DNS Servers:** <list or "Azure-provided">

#### Subnets
| Name | Address Prefix | NSG | Route Table | NAT Gateway | Service Endpoints |
|------|----------------|-----|-------------|-------------|-------------------|
| ... | ... | ... | ... | ... | ... |

#### Peerings
| Name | Remote VNet | Status | Gateway Transit | Forwarded Traffic |
|------|-------------|--------|-----------------|-------------------|
| ... | ... | ... | ... | ... |

## Network Security Groups

### NSG: <name>
**Associated with:**
- Subnets: <list>
- Network Interfaces: <list>

#### Security Rules
| Priority | Name | Direction | Access | Protocol | Source | Dest | Source Port | Dest Port |
|----------|------|-----------|--------|----------|--------|------|-------------|-----------|
| ... | ... | ... | ... | ... | ... | ... | ... | ... |

## Private Endpoints
...

## VPN Gateways
...

## ExpressRoute Circuits
...

## Load Balancers
...

## Application Gateways
...

## Security Findings
### Critical
- <finding>: <description>

### High
...

## Network Watcher Insights
...
```

17. **Implement pkg/reporter/html.go**
- `GenerateHTML(topology *NetworkTopology, analysis *AnalysisReport) string`
- Create rich HTML report with:
  - CSS for styling (embedded or inline)
  - Collapsible sections for each resource type
  - Tables with sortable columns (using simple JavaScript)
  - Color-coded security findings
  - Embedded SVG diagram if available
  - Navigation sidebar/table of contents
- Use template with clean, professional design
- Make it responsive for mobile viewing

### Phase 6: CLI Implementation

18. **Implement cmd/root.go**
```go
package cmd

import (
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "azure-network-analyzer",
    Short: "Analyze Azure network topology",
    Long:  `A comprehensive tool to analyze and visualize Azure network resources.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    // Global flags can be added here
}
```

19. **Implement cmd/analyze.go**
```go
var analyzeCmd = &cobra.Command{
    Use:   "analyze",
    Short: "Analyze network topology for a resource group",
    Long:  `Collect and analyze Azure network resources in a specified resource group.`,
    RunE:  runAnalyze,
}

var (
    subscriptionID string
    resourceGroup  string
    outputFormat   string
    outputPath     string
    includeViz     bool
    vizFormat      string
)

func init() {
    rootCmd.AddCommand(analyzeCmd)
    
    analyzeCmd.Flags().StringVarP(&subscriptionID, "subscription", "s", "", "Azure subscription ID")
    analyzeCmd.Flags().StringVarP(&resourceGroup, "resource-group", "g", "", "Resource group name")
    analyzeCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "markdown", "Output format (json|markdown|html)")
    analyzeCmd.Flags().StringVarP(&outputPath, "output", "f", "", "Output file path")
    analyzeCmd.Flags().BoolVar(&includeViz, "visualize", true, "Generate network topology diagram")
    analyzeCmd.Flags().StringVar(&vizFormat, "viz-format", "svg", "Visualization format (svg|png|dot)")
    
    analyzeCmd.MarkFlagRequired("subscription")
    analyzeCmd.MarkFlagRequired("resource-group")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
    ctx := context.Background()
    
    // 1. Initialize Azure client
    client, err := azure.NewAzureClient(subscriptionID)
    if err != nil {
        return fmt.Errorf("failed to create Azure client: %w", err)
    }
    
    // 2. Collect all network resources
    fmt.Println("Collecting network resources...")
    topology := &models.NetworkTopology{
        SubscriptionID: subscriptionID,
        ResourceGroup:  resourceGroup,
        Timestamp:      time.Now(),
    }
    
    // Collect VNets
    vnets, err := client.GetVirtualNetworks(ctx, resourceGroup)
    // ... error handling and progress indication
    topology.VirtualNetworks = vnets
    
    // Collect NSGs
    nsgs, err := client.GetNetworkSecurityGroups(ctx, resourceGroup)
    topology.NSGs = nsgs
    
    // ... collect all other resources
    
    // 3. Analyze topology
    fmt.Println("Analyzing topology...")
    analysisReport := analyzer.Analyze(topology)
    securityFindings := analyzer.AnalyzeNSGRules(topology.NSGs)
    
    // 4. Generate reports
    fmt.Println("Generating reports...")
    
    switch outputFormat {
    case "json":
        data, err := reporter.GenerateJSON(topology)
        // Write to file or stdout
    case "markdown":
        report := reporter.GenerateMarkdown(topology, analysisReport)
        // Write to file or stdout
    case "html":
        report := reporter.GenerateHTML(topology, analysisReport)
        // Write to file or stdout
    }
    
    // 5. Generate visualization if requested
    if includeViz {
        fmt.Println("Generating topology diagram...")
        dotContent := visualization.GenerateDOTFile(topology)
        
        switch vizFormat {
        case "svg":
            err = visualization.RenderSVG(dotContent, "topology.svg")
        case "png":
            err = visualization.RenderPNG(dotContent, "topology.png")
        case "dot":
            // Write DOT file directly
            err = os.WriteFile("topology.dot", []byte(dotContent), 0644)
        }
    }
    
    fmt.Println("Analysis complete!")
    return nil
}
```

### Phase 7: Error Handling and Polish

20. **Add comprehensive error handling**
- Wrap all Azure SDK calls with proper error checking
- Provide helpful error messages for common issues:
  - Authentication failures
  - Permission errors (suggest required RBAC roles: Reader on subscription/RG)
  - Resource not found errors
  - Network errors
- Use structured logging (consider adding zerolog or zap)

21. **Add progress indicators**
- Use spinner or progress bar for long-running operations
- Show which resource type is currently being collected
- Estimated completion percentage

22. **Add configuration file support**
- Support config file (YAML/JSON) for default settings
- Allow specifying multiple resource groups to analyze in sequence
- Store API pagination settings, timeouts

23. **Optimization**
- Use goroutines to parallelize resource collection where possible
- Implement connection pooling for Azure SDK clients
- Add caching for resource queries (with TTL)
- Allow incremental updates (compare with previous run)

### Phase 8: Testing

24. **Unit tests**
- Test data model marshaling/unmarshaling
- Test security analysis logic with known vulnerable configurations
- Test DOT file generation with sample topology
- Test report generation

25. **Integration tests**
- Create test fixtures with mock Azure responses
- Test end-to-end flow with mock Azure client
- Consider using azidentity's test helpers

26. **Manual testing checklist**
- Test with empty resource group
- Test with resource group containing each resource type
- Test with cross-RG VNet peerings
- Test with complex NSG rule sets
- Test with missing Network Watcher
- Test output formats (JSON, Markdown, HTML)
- Test visualization rendering
- Test authentication methods (Azure CLI, service principal)

### Phase 9: Documentation

27. **README.md**
```markdown
# Azure Network Topology Analyzer

A comprehensive CLI tool to analyze and visualize Azure network infrastructure.

## Features
- Collects all network resources from a resource group
- Generates detailed reports in multiple formats
- Creates network topology visualizations
- Analyzes NSG rules for security issues
- Supports VNets, Subnets, NSGs, Private Links, VPN/ER, Load Balancers, and more

## Installation
...

## Authentication
The tool uses Azure DefaultAzureCredential, supporting:
- Azure CLI (`az login`)
- Managed Identity
- Environment variables
- Service Principal

## Usage
...

## Required Azure Permissions
The service principal or user must have at minimum:
- Reader role on the subscription or resource group

## Output Examples
...

## Troubleshooting
...
```

28. **Add examples directory**
- Include sample outputs (JSON, Markdown, HTML)
- Include sample topology diagrams
- Include sample config files

### Phase 10: Build and Release

29. **Add Makefile**
```makefile
.PHONY: build test install clean

build:
	go build -o bin/azure-network-analyzer main.go

test:
	go test ./...

install:
	go install

clean:
	rm -rf bin/

lint:
	golangci-lint run

release:
	goreleaser release --snapshot --rm-dist
```

30. **Add GitHub Actions / CI/CD**
- Build and test on multiple platforms (Linux, macOS, Windows)
- Generate binaries for releases
- Run linters and security scanners

## Advanced Features (Optional Enhancements)

### Cross-Resource Group Analysis
- Add flag `--all-resource-groups` to analyze entire subscription
- Show cross-RG dependencies (peerings, ER connections)

### Diff Mode
- Add `--compare-with <previous-output.json>` flag
- Show what changed since last run
- Highlight new resources, deleted resources, configuration changes

### Export Formats
- Add Terraform export capability
- Generate ARM templates from discovered topology
- Export to draw.io format

### Compliance Checking
- Add compliance rules (CIS Azure Foundations Benchmark)
- Flag non-compliant configurations
- Generate compliance report

### Cost Estimation
- Integrate with Azure Pricing API
- Estimate monthly costs for network resources
- Show cost breakdown by resource type

### Interactive Mode
- Add TUI (terminal UI) using bubbletea or similar
- Allow interactive exploration of topology
- Drill down into specific resources

## Key Implementation Notes

1. **Azure SDK Pagination**: Many list operations return pagers - ensure you iterate through all pages
2. **Rate Limiting**: Be mindful of Azure API rate limits; add exponential backoff retry logic
3. **Resource Dependencies**: Some resources (like subnets) need parent resource info - build dependency graph
4. **Cross-Subscription Resources**: ExpressRoute and VNet peering can span subscriptions - handle gracefully
5. **Partial Failures**: If one resource type fails to collect, continue with others and report errors at end
6. **Large Topologies**: For very large environments (100+ VNets), consider streaming output and incremental visualization
7. **Sensitive Data**: Avoid logging or outputting secrets, connection strings, or keys

## Testing the Tool

### Test Resource Group Setup
Create a test resource group with:
- 2-3 VNets with various subnet configurations
- VNet peerings between them
- NSGs with diverse rule sets
- Private endpoint to Azure Storage/Key Vault
- Private DNS zone with VNet links
- Route table with custom routes
- NAT gateway
- VPN gateway with test connection
- Load balancer (standard SKU)
- Application gateway with WAF enabled

Run tool against this test RG and verify all resources are captured and visualized correctly.

## Success Criteria

- ✅ Tool successfully authenticates with Azure using DefaultAzureCredential
- ✅ Collects all specified resource types from a resource group
- ✅ Generates accurate JSON export of topology
- ✅ Generates readable Markdown report with all sections
- ✅ Generates formatted HTML report with styling
- ✅ Creates Graphviz DOT file representing topology
- ✅ Renders SVG/PNG visualization successfully
- ✅ Identifies and reports security findings from NSG analysis
- ✅ Handles errors gracefully with helpful messages
- ✅ Completes analysis in reasonable time (<2 minutes for typical RG)
- ✅ Documentation is clear and comprehensive

## Development Tips for Claude Code

1. Start with the foundation: models, Azure client initialization, and one simple resource type (e.g., VNets)
2. Test collection of one resource type end-to-end before moving to the next
3. Build reporting capability early - it helps validate data collection
4. Add visualization last - it's the most complex but builds on complete data
5. Use meaningful commit messages at each phase
6. Keep functions focused and small - easier to test and debug
7. Add TODOs for optimizations or nice-to-haves to implement later

Good luck building this comprehensive Azure network analysis tool! 🚀
