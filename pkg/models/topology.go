package models

import "time"

// NetworkTopology represents the complete network topology for a resource group
type NetworkTopology struct {
	SubscriptionID   string                   `json:"subscriptionId"`
	ResourceGroup    string                   `json:"resourceGroup"`
	VirtualNetworks  []VirtualNetwork         `json:"virtualNetworks"`
	NSGs             []NetworkSecurityGroup   `json:"networkSecurityGroups"`
	PrivateEndpoints []PrivateEndpoint        `json:"privateEndpoints"`
	PrivateDNSZones  []PrivateDNSZone         `json:"privateDnsZones"`
	RouteTables      []RouteTable             `json:"routeTables"`
	NATGateways      []NATGateway             `json:"natGateways"`
	VPNGateways      []VPNGateway             `json:"vpnGateways"`
	ERCircuits       []ExpressRouteCircuit    `json:"expressRouteCircuits"`
	LoadBalancers    []LoadBalancer           `json:"loadBalancers"`
	AppGateways      []ApplicationGateway     `json:"applicationGateways"`
	AzureFirewalls   []AzureFirewall          `json:"azureFirewalls"`
	NetworkWatcher   *NetworkWatcherInsights  `json:"networkWatcher,omitempty"`
	Timestamp        time.Time                `json:"timestamp"`
}

// VirtualNetwork represents an Azure Virtual Network
type VirtualNetwork struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	ResourceGroup string        `json:"resourceGroup"`
	Location      string        `json:"location"`
	AddressSpace  []string      `json:"addressSpace"`
	Subnets       []Subnet      `json:"subnets"`
	Peerings      []VNetPeering `json:"peerings"`
	DNSServers    []string      `json:"dnsServers"`
	EnableDDoS    bool          `json:"enableDdosProtection"`
}

// Subnet represents a subnet within a virtual network
type Subnet struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	AddressPrefix        string   `json:"addressPrefix"`
	NetworkSecurityGroup *string  `json:"networkSecurityGroup,omitempty"` // NSG ID if associated
	RouteTable           *string  `json:"routeTable,omitempty"`           // Route table ID if associated
	NATGateway           *string  `json:"natGateway,omitempty"`           // NAT gateway ID if associated
	PrivateEndpoints     []string `json:"privateEndpoints"`               // List of private endpoint IDs
	ServiceEndpoints     []string `json:"serviceEndpoints"`
	Delegations          []string `json:"delegations"`
}

// NetworkSecurityGroup represents an Azure NSG
type NetworkSecurityGroup struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	ResourceGroup string          `json:"resourceGroup"`
	Location      string          `json:"location"`
	SecurityRules []SecurityRule  `json:"securityRules"`
	Associations  NSGAssociations `json:"associations"`
}

// SecurityRule represents a security rule within an NSG
type SecurityRule struct {
	Name                     string `json:"name"`
	Priority                 int32  `json:"priority"`
	Direction                string `json:"direction"` // Inbound/Outbound
	Access                   string `json:"access"`    // Allow/Deny
	Protocol                 string `json:"protocol"`
	SourceAddressPrefix      string `json:"sourceAddressPrefix"`
	SourcePortRange          string `json:"sourcePortRange"`
	DestinationAddressPrefix string `json:"destinationAddressPrefix"`
	DestinationPortRange     string `json:"destinationPortRange"`
	Description              string `json:"description"`
}

// NSGAssociations tracks what resources are associated with an NSG
type NSGAssociations struct {
	Subnets           []string `json:"subnets"`           // Subnet IDs
	NetworkInterfaces []string `json:"networkInterfaces"` // NIC IDs
}

// VNetPeering represents a peering connection between VNets
type VNetPeering struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	RemoteVNetID          string `json:"remoteVnetId"`
	RemoteVNetName        string `json:"remoteVnetName"`
	PeeringState          string `json:"peeringState"`
	AllowVNetAccess       bool   `json:"allowVnetAccess"`
	AllowForwardedTraffic bool   `json:"allowForwardedTraffic"`
	AllowGatewayTransit   bool   `json:"allowGatewayTransit"`
	UseRemoteGateways     bool   `json:"useRemoteGateways"`
}

// PrivateEndpoint represents an Azure Private Endpoint
type PrivateEndpoint struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	ResourceGroup        string   `json:"resourceGroup"`
	Location             string   `json:"location"`
	SubnetID             string   `json:"subnetId"`
	PrivateIPAddress     string   `json:"privateIpAddress"`
	PrivateLinkServiceID string   `json:"privateLinkServiceId"`
	ConnectionState      string   `json:"connectionState"`
	GroupIDs             []string `json:"groupIds"`
}

// PrivateDNSZone represents an Azure Private DNS Zone
type PrivateDNSZone struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	ResourceGroup string     `json:"resourceGroup"`
	VNetLinks     []VNetLink `json:"vnetLinks"`
	RecordSets    int        `json:"recordSets"`
}

// VNetLink represents a link between a Private DNS Zone and a VNet
type VNetLink struct {
	ID                  string `json:"id"`
	VNetID              string `json:"vnetId"`
	VNetName            string `json:"vnetName"`
	RegistrationEnabled bool   `json:"registrationEnabled"`
}

// RouteTable represents an Azure Route Table
type RouteTable struct {
	ID                          string   `json:"id"`
	Name                        string   `json:"name"`
	ResourceGroup               string   `json:"resourceGroup"`
	Location                    string   `json:"location"`
	Routes                      []Route  `json:"routes"`
	DisableBGPRoutePropagation bool     `json:"disableBgpRoutePropagation"`
	AssociatedSubnets           []string `json:"associatedSubnets"`
}

// Route represents a route within a route table
type Route struct {
	Name             string `json:"name"`
	AddressPrefix    string `json:"addressPrefix"`
	NextHopType      string `json:"nextHopType"` // VirtualNetworkGateway, VNetLocal, Internet, VirtualAppliance, None
	NextHopIPAddress string `json:"nextHopIpAddress"`
}

// NATGateway represents an Azure NAT Gateway
type NATGateway struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	ResourceGroup      string   `json:"resourceGroup"`
	Location           string   `json:"location"`
	PublicIPAddresses  []string `json:"publicIpAddresses"`
	IdleTimeoutMinutes int32    `json:"idleTimeoutMinutes"`
	AssociatedSubnets  []string `json:"associatedSubnets"`
}

// VPNGateway represents an Azure VPN Gateway
type VPNGateway struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	ResourceGroup string          `json:"resourceGroup"`
	Location     string          `json:"location"`
	VNetID       string          `json:"vnetId"`
	GatewayType  string          `json:"gatewayType"` // Vpn or ExpressRoute
	VpnType      string          `json:"vpnType"`     // RouteBased or PolicyBased
	SKU          string          `json:"sku"`
	ActiveActive bool            `json:"activeActive"`
	BGPSettings  *BGPSettings    `json:"bgpSettings,omitempty"`
	Connections  []VPNConnection `json:"connections"`
}

// BGPSettings represents BGP configuration for a gateway
type BGPSettings struct {
	ASN               int64  `json:"asn"`
	BGPPeeringAddress string `json:"bgpPeeringAddress"`
	PeerWeight        int32  `json:"peerWeight"`
}

// VPNConnection represents a VPN connection
type VPNConnection struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ConnectionType   string `json:"connectionType"` // IPsec, Vnet2Vnet, ExpressRoute
	ConnectionStatus string `json:"connectionStatus"`
	SharedKey        bool   `json:"sharedKey"` // Whether a shared key is configured
	EnableBGP        bool   `json:"enableBgp"`
	RemoteEntityID   string `json:"remoteEntityId"`
}

// ExpressRouteCircuit represents an Azure ExpressRoute Circuit
type ExpressRouteCircuit struct {
	ID                       string            `json:"id"`
	Name                     string            `json:"name"`
	ResourceGroup            string            `json:"resourceGroup"`
	Location                 string            `json:"location"`
	ServiceProviderName      string            `json:"serviceProviderName"`
	PeeringLocation          string            `json:"peeringLocation"`
	BandwidthInMbps          int32             `json:"bandwidthInMbps"`
	SKUTier                  string            `json:"skuTier"`
	SKUFamily                string            `json:"skuFamily"`
	CircuitProvisioningState string            `json:"circuitProvisioningState"`
	Peerings                 []ERPeering       `json:"peerings"`
	Authorizations           []ERAuthorization `json:"authorizations"`
}

// ERPeering represents an ExpressRoute peering
type ERPeering struct {
	Name                       string `json:"name"`
	PeeringType                string `json:"peeringType"` // AzurePrivatePeering, AzurePublicPeering, MicrosoftPeering
	State                      string `json:"state"`
	AzureASN                   int32  `json:"azureAsn"`
	PeerASN                    int64  `json:"peerAsn"`
	PrimaryPeerAddressPrefix   string `json:"primaryPeerAddressPrefix"`
	SecondaryPeerAddressPrefix string `json:"secondaryPeerAddressPrefix"`
	VlanID                     int32  `json:"vlanId"`
}

// ERAuthorization represents an ExpressRoute authorization
type ERAuthorization struct {
	Name                   string `json:"name"`
	AuthorizationKey       bool   `json:"authorizationKey"` // Whether key exists
	AuthorizationUseStatus string `json:"authorizationUseStatus"`
}

// LoadBalancer represents an Azure Load Balancer
type LoadBalancer struct {
	ID                  string               `json:"id"`
	Name                string               `json:"name"`
	ResourceGroup       string               `json:"resourceGroup"`
	Location            string               `json:"location"`
	SKU                 string               `json:"sku"`
	Type                string               `json:"type"` // Public or Internal
	FrontendIPConfigs   []FrontendIPConfig   `json:"frontendIpConfigs"`
	BackendAddressPools []BackendAddressPool `json:"backendAddressPools"`
	LoadBalancingRules  []LoadBalancingRule  `json:"loadBalancingRules"`
	Probes              []Probe              `json:"probes"`
	InboundNATRules     []InboundNATRule     `json:"inboundNatRules"`
}

// FrontendIPConfig represents a frontend IP configuration for a load balancer
type FrontendIPConfig struct {
	Name              string `json:"name"`
	PrivateIPAddress  string `json:"privateIpAddress"`
	PublicIPAddressID string `json:"publicIpAddressId"`
	SubnetID          string `json:"subnetId"`
}

// BackendAddressPool represents a backend address pool for a load balancer
type BackendAddressPool struct {
	Name             string   `json:"name"`
	BackendIPConfigs []string `json:"backendIpConfigs"` // NIC IDs
}

// LoadBalancingRule represents a load balancing rule
type LoadBalancingRule struct {
	Name               string `json:"name"`
	Protocol           string `json:"protocol"`
	FrontendPort       int32  `json:"frontendPort"`
	BackendPort        int32  `json:"backendPort"`
	EnableFloatingIP   bool   `json:"enableFloatingIp"`
	IdleTimeoutMinutes int32  `json:"idleTimeoutMinutes"`
	LoadDistribution   string `json:"loadDistribution"`
}

// Probe represents a health probe for a load balancer
type Probe struct {
	Name              string `json:"name"`
	Protocol          string `json:"protocol"`
	Port              int32  `json:"port"`
	IntervalInSeconds int32  `json:"intervalInSeconds"`
	NumberOfProbes    int32  `json:"numberOfProbes"`
	RequestPath       string `json:"requestPath"` // For HTTP/HTTPS
}

// InboundNATRule represents an inbound NAT rule for a load balancer
type InboundNATRule struct {
	Name             string `json:"name"`
	Protocol         string `json:"protocol"`
	FrontendPort     int32  `json:"frontendPort"`
	BackendPort      int32  `json:"backendPort"`
	EnableFloatingIP bool   `json:"enableFloatingIp"`
}

// ApplicationGateway represents an Azure Application Gateway
type ApplicationGateway struct {
	ID                  string                      `json:"id"`
	Name                string                      `json:"name"`
	ResourceGroup       string                      `json:"resourceGroup"`
	Location            string                      `json:"location"`
	SKU                 string                      `json:"sku"`
	Tier                string                      `json:"tier"`
	Capacity            int32                       `json:"capacity"`
	SubnetID            string                      `json:"subnetId"`
	FrontendIPConfigs   []AppGWFrontendIPConfig     `json:"frontendIpConfigs"`
	FrontendPorts       []AppGWFrontendPort         `json:"frontendPorts"`
	BackendAddressPools []AppGWBackendAddressPool   `json:"backendAddressPools"`
	BackendHTTPSettings []AppGWBackendHTTPSettings  `json:"backendHttpSettings"`
	HTTPListeners       []AppGWHTTPListener         `json:"httpListeners"`
	RequestRoutingRules []AppGWRequestRoutingRule   `json:"requestRoutingRules"`
	Probes              []AppGWProbe                `json:"probes"`
	WAFEnabled          bool                        `json:"wafEnabled"`
	WAFMode             string                      `json:"wafMode"`
}

// AppGWFrontendIPConfig represents a frontend IP configuration for an Application Gateway
type AppGWFrontendIPConfig struct {
	Name              string `json:"name"`
	PrivateIPAddress  string `json:"privateIpAddress"`
	PublicIPAddressID string `json:"publicIpAddressId"`
}

// AppGWFrontendPort represents a frontend port for an Application Gateway
type AppGWFrontendPort struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

// AppGWBackendAddressPool represents a backend address pool for an Application Gateway
type AppGWBackendAddressPool struct {
	Name             string   `json:"name"`
	BackendAddresses []string `json:"backendAddresses"` // IP addresses or FQDNs
}

// AppGWBackendHTTPSettings represents backend HTTP settings for an Application Gateway
type AppGWBackendHTTPSettings struct {
	Name                string `json:"name"`
	Port                int32  `json:"port"`
	Protocol            string `json:"protocol"`
	CookieBasedAffinity string `json:"cookieBasedAffinity"`
	RequestTimeout      int32  `json:"requestTimeout"`
	ProbeName           string `json:"probeName"`
}

// AppGWHTTPListener represents an HTTP listener for an Application Gateway
type AppGWHTTPListener struct {
	Name             string `json:"name"`
	FrontendIPConfig string `json:"frontendIpConfig"`
	FrontendPort     string `json:"frontendPort"`
	Protocol         string `json:"protocol"`
	HostName         string `json:"hostName"`
}

// AppGWRequestRoutingRule represents a request routing rule for an Application Gateway
type AppGWRequestRoutingRule struct {
	Name                string `json:"name"`
	RuleType            string `json:"ruleType"`
	HTTPListener        string `json:"httpListener"`
	BackendAddressPool  string `json:"backendAddressPool"`
	BackendHTTPSettings string `json:"backendHttpSettings"`
	Priority            int32  `json:"priority"`
}

// AppGWProbe represents a health probe for an Application Gateway
type AppGWProbe struct {
	Name               string `json:"name"`
	Protocol           string `json:"protocol"`
	Host               string `json:"host"`
	Path               string `json:"path"`
	Interval           int32  `json:"interval"`
	Timeout            int32  `json:"timeout"`
	UnhealthyThreshold int32  `json:"unhealthyThreshold"`
}

// AzureFirewall represents an Azure Firewall
type AzureFirewall struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	ResourceGroup     string   `json:"resourceGroup"`
	Location          string   `json:"location"`
	SKU               string   `json:"sku"` // Standard, Premium, Basic
	SubnetID          string   `json:"subnetId"`
	PrivateIPAddress  string   `json:"privateIpAddress"`
	PublicIPAddresses []string `json:"publicIpAddresses"`
	FirewallPolicyID  string   `json:"firewallPolicyId,omitempty"`
	ThreatIntelMode   string   `json:"threatIntelMode"`
	DNSProxyEnabled   bool     `json:"dnsProxyEnabled"`
	ProvisioningState string   `json:"provisioningState"`
}

// NetworkWatcherInsights contains Network Watcher related information
type NetworkWatcherInsights struct {
	FlowLogsEnabled    bool                `json:"flowLogsEnabled"`
	FlowLogs           []FlowLog           `json:"flowLogs"`
	ConnectionMonitors []ConnectionMonitor `json:"connectionMonitors"`
	PacketCaptures     []PacketCapture     `json:"packetCaptures"`
}

// FlowLog represents an NSG flow log configuration
type FlowLog struct {
	ID               string `json:"id"`
	NSGId            string `json:"nsgId"`
	StorageAccountID string `json:"storageAccountId"`
	Enabled          bool   `json:"enabled"`
	RetentionDays    int32  `json:"retentionDays"`
	TrafficAnalytics bool   `json:"trafficAnalytics"`
}

// ConnectionMonitor represents a Network Watcher connection monitor
type ConnectionMonitor struct {
	Name             string `json:"name"`
	Source           string `json:"source"`
	Destination      string `json:"destination"`
	MonitoringStatus string `json:"monitoringStatus"`
}

// PacketCapture represents a Network Watcher packet capture
type PacketCapture struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Status string `json:"status"`
}
