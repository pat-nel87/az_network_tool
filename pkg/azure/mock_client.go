package azure

import (
	"context"
	"time"

	"azure-network-analyzer/pkg/models"
)

// MockAzureClient provides mock data for testing without Azure connectivity
type MockAzureClient struct {
	subscriptionID string
}

// NewMockAzureClient creates a new mock Azure client for testing
func NewMockAzureClient(subscriptionID string) *MockAzureClient {
	return &MockAzureClient{
		subscriptionID: subscriptionID,
	}
}

// GetVirtualNetworks returns mock VNet data
func (c *MockAzureClient) GetVirtualNetworks(ctx context.Context, resourceGroup string) ([]models.VirtualNetwork, error) {
	nsgID := "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/networkSecurityGroups/nsg-web"
	routeTableID := "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/routeTables/rt-main"
	natGatewayID := "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/natGateways/nat-outbound"

	return []models.VirtualNetwork{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub",
			Name:          "vnet-hub",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			AddressSpace:  []string{"10.0.0.0/16"},
			DNSServers:    []string{"10.0.0.4", "10.0.0.5"},
			EnableDDoS:    true,
			Subnets: []models.Subnet{
				{
					ID:                   "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/AzureFirewallSubnet",
					Name:                 "AzureFirewallSubnet",
					AddressPrefix:        "10.0.0.0/26",
					NetworkSecurityGroup: nil, // Firewall subnet doesn't need NSG
					RouteTable:           nil,
					NATGateway:           nil,
					PrivateEndpoints:     []string{},
					ServiceEndpoints:     []string{},
					Delegations:          []string{},
				},
				{
					ID:                   "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-web",
					Name:                 "subnet-web",
					AddressPrefix:        "10.0.1.0/24",
					NetworkSecurityGroup: &nsgID,
					RouteTable:           &routeTableID,
					NATGateway:           &natGatewayID,
					PrivateEndpoints:     []string{},
					ServiceEndpoints:     []string{"Microsoft.Storage", "Microsoft.KeyVault"},
					Delegations:          []string{},
				},
				{
					ID:                   "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-db",
					Name:                 "subnet-db",
					AddressPrefix:        "10.0.2.0/24",
					NetworkSecurityGroup: &nsgID,
					RouteTable:           nil,
					NATGateway:           nil,
					PrivateEndpoints:     []string{"/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/privateEndpoints/pe-sql"},
					ServiceEndpoints:     []string{"Microsoft.Sql"},
					Delegations:          []string{},
				},
				{
					ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/GatewaySubnet",
					Name:          "GatewaySubnet",
					AddressPrefix: "10.0.255.0/27",
					PrivateEndpoints: []string{},
					ServiceEndpoints: []string{},
					Delegations:      []string{},
				},
			},
			Peerings: []models.VNetPeering{
				{
					ID:                    "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/virtualNetworkPeerings/peer-to-spoke",
					Name:                  "peer-to-spoke",
					RemoteVNetID:          "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-spoke",
					RemoteVNetName:        "vnet-spoke",
					PeeringState:          "Connected",
					AllowVNetAccess:       true,
					AllowForwardedTraffic: true,
					AllowGatewayTransit:   true,
					UseRemoteGateways:     false,
				},
			},
		},
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-spoke",
			Name:          "vnet-spoke",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			AddressSpace:  []string{"10.1.0.0/16"},
			DNSServers:    []string{},
			EnableDDoS:    false,
			Subnets: []models.Subnet{
				{
					ID:               "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-spoke/subnets/subnet-app",
					Name:             "subnet-app",
					AddressPrefix:    "10.1.1.0/24",
					PrivateEndpoints: []string{},
					ServiceEndpoints: []string{},
					Delegations:      []string{"Microsoft.Web/serverFarms"},
				},
			},
			Peerings: []models.VNetPeering{
				{
					ID:                    "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-spoke/virtualNetworkPeerings/peer-to-hub",
					Name:                  "peer-to-hub",
					RemoteVNetID:          "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub",
					RemoteVNetName:        "vnet-hub",
					PeeringState:          "Connected",
					AllowVNetAccess:       true,
					AllowForwardedTraffic: false,
					AllowGatewayTransit:   false,
					UseRemoteGateways:     true,
				},
			},
		},
	}, nil
}

// GetNetworkSecurityGroups returns mock NSG data
func (c *MockAzureClient) GetNetworkSecurityGroups(ctx context.Context, resourceGroup string) ([]models.NetworkSecurityGroup, error) {
	return []models.NetworkSecurityGroup{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/networkSecurityGroups/nsg-web",
			Name:          "nsg-web",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			SecurityRules: []models.SecurityRule{
				{
					Name:                     "AllowHTTP",
					Priority:                 100,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "TCP",
					SourceAddressPrefix:      "*",
					SourcePortRange:          "*",
					DestinationAddressPrefix: "*",
					DestinationPortRange:     "80",
					Description:              "Allow HTTP traffic",
				},
				{
					Name:                     "AllowHTTPS",
					Priority:                 110,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "TCP",
					SourceAddressPrefix:      "*",
					SourcePortRange:          "*",
					DestinationAddressPrefix: "*",
					DestinationPortRange:     "443",
					Description:              "Allow HTTPS traffic",
				},
				{
					Name:                     "AllowSSH",
					Priority:                 120,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "TCP",
					SourceAddressPrefix:      "0.0.0.0/0",
					SourcePortRange:          "*",
					DestinationAddressPrefix: "*",
					DestinationPortRange:     "22",
					Description:              "", // Missing description - security finding
				},
				{
					Name:                     "DenyAll",
					Priority:                 4096,
					Direction:                "Inbound",
					Access:                   "Deny",
					Protocol:                 "*",
					SourceAddressPrefix:      "*",
					SourcePortRange:          "*",
					DestinationAddressPrefix: "*",
					DestinationPortRange:     "*",
					Description:              "Deny all other inbound traffic",
				},
			},
			Associations: models.NSGAssociations{
				Subnets: []string{
					"/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-web",
					"/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-db",
				},
				NetworkInterfaces: []string{},
			},
		},
	}, nil
}

// GetPrivateEndpoints returns mock private endpoint data
func (c *MockAzureClient) GetPrivateEndpoints(ctx context.Context, resourceGroup string) ([]models.PrivateEndpoint, error) {
	return []models.PrivateEndpoint{
		{
			ID:                   "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/privateEndpoints/pe-sql",
			Name:                 "pe-sql",
			ResourceGroup:        resourceGroup,
			Location:             "eastus",
			SubnetID:             "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-db",
			PrivateIPAddress:     "10.0.2.10",
			PrivateLinkServiceID: "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Sql/servers/sql-server-prod",
			ConnectionState:      "Approved",
			GroupIDs:             []string{"sqlServer"},
		},
		{
			ID:                   "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/privateEndpoints/pe-storage",
			Name:                 "pe-storage",
			ResourceGroup:        resourceGroup,
			Location:             "eastus",
			SubnetID:             "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-db",
			PrivateIPAddress:     "10.0.2.11",
			PrivateLinkServiceID: "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Storage/storageAccounts/stprod",
			ConnectionState:      "Approved",
			GroupIDs:             []string{"blob"},
		},
	}, nil
}

// GetPrivateDNSZones returns mock private DNS zone data
func (c *MockAzureClient) GetPrivateDNSZones(ctx context.Context, resourceGroup string) ([]models.PrivateDNSZone, error) {
	return []models.PrivateDNSZone{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/privateDnsZones/privatelink.database.windows.net",
			Name:          "privatelink.database.windows.net",
			ResourceGroup: resourceGroup,
			RecordSets:    5,
			VNetLinks: []models.VNetLink{
				{
					ID:                  "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/privateDnsZones/privatelink.database.windows.net/virtualNetworkLinks/link-to-hub",
					VNetID:              "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub",
					VNetName:            "vnet-hub",
					RegistrationEnabled: false,
				},
			},
		},
	}, nil
}

// GetRouteTables returns mock route table data
func (c *MockAzureClient) GetRouteTables(ctx context.Context, resourceGroup string) ([]models.RouteTable, error) {
	return []models.RouteTable{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/routeTables/rt-main",
			Name:          "rt-main",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			Routes: []models.Route{
				{
					Name:             "route-to-internet",
					AddressPrefix:    "0.0.0.0/0",
					NextHopType:      "VirtualAppliance",
					NextHopIPAddress: "10.0.0.4", // Points to Azure Firewall
				},
				{
					Name:             "route-to-onprem",
					AddressPrefix:    "192.168.0.0/16",
					NextHopType:      "VirtualAppliance",
					NextHopIPAddress: "10.0.0.4", // Points to Azure Firewall
				},
			},
			DisableBGPRoutePropagation: false,
			AssociatedSubnets: []string{
				"/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-web",
			},
		},
	}, nil
}

// GetNATGateways returns mock NAT gateway data
func (c *MockAzureClient) GetNATGateways(ctx context.Context, resourceGroup string) ([]models.NATGateway, error) {
	return []models.NATGateway{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/natGateways/nat-outbound",
			Name:          "nat-outbound",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			PublicIPAddresses: []string{
				"/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/publicIPAddresses/pip-nat",
			},
			IdleTimeoutMinutes: 10,
			AssociatedSubnets: []string{
				"/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-web",
			},
		},
	}, nil
}

// GetVPNGateways returns mock VPN gateway data
func (c *MockAzureClient) GetVPNGateways(ctx context.Context, resourceGroup string) ([]models.VPNGateway, error) {
	return []models.VPNGateway{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworkGateways/vpn-gateway",
			Name:          "vpn-gateway",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			VNetID:        "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub",
			GatewayType:   "Vpn",
			VpnType:       "RouteBased",
			SKU:           "VpnGw2",
			ActiveActive:  false,
			BGPSettings: &models.BGPSettings{
				ASN:               65515,
				BGPPeeringAddress: "10.0.255.30",
				PeerWeight:        0,
			},
			Connections: []models.VPNConnection{
				{
					ID:               "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/connections/vpn-to-onprem",
					Name:             "vpn-to-onprem",
					ConnectionType:   "IPsec",
					ConnectionStatus: "Connected",
					SharedKey:        true,
					EnableBGP:        true,
					RemoteEntityID:   "52.168.1.100",
				},
			},
		},
	}, nil
}

// GetExpressRouteCircuits returns mock ExpressRoute circuit data
func (c *MockAzureClient) GetExpressRouteCircuits(ctx context.Context, resourceGroup string) ([]models.ExpressRouteCircuit, error) {
	return []models.ExpressRouteCircuit{}, nil // No ER circuits in mock
}

// GetLoadBalancers returns mock load balancer data
func (c *MockAzureClient) GetLoadBalancers(ctx context.Context, resourceGroup string) ([]models.LoadBalancer, error) {
	return []models.LoadBalancer{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/loadBalancers/lb-web",
			Name:          "lb-web",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			SKU:           "Standard",
			Type:          "Public",
			FrontendIPConfigs: []models.FrontendIPConfig{
				{
					Name:              "frontend-public",
					PrivateIPAddress:  "",
					PublicIPAddressID: "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/publicIPAddresses/pip-lb",
					SubnetID:          "",
				},
			},
			BackendAddressPools: []models.BackendAddressPool{
				{
					Name:             "backend-web-servers",
					BackendIPConfigs: []string{"nic-web-1", "nic-web-2", "nic-web-3"},
				},
			},
			LoadBalancingRules: []models.LoadBalancingRule{
				{
					Name:               "rule-http",
					Protocol:           "TCP",
					FrontendPort:       80,
					BackendPort:        80,
					EnableFloatingIP:   false,
					IdleTimeoutMinutes: 4,
					LoadDistribution:   "Default",
				},
				{
					Name:               "rule-https",
					Protocol:           "TCP",
					FrontendPort:       443,
					BackendPort:        443,
					EnableFloatingIP:   false,
					IdleTimeoutMinutes: 4,
					LoadDistribution:   "Default",
				},
			},
			Probes: []models.Probe{
				{
					Name:              "probe-http",
					Protocol:          "HTTP",
					Port:              80,
					IntervalInSeconds: 15,
					NumberOfProbes:    2,
					RequestPath:       "/health",
				},
			},
			InboundNATRules: []models.InboundNATRule{},
		},
	}, nil
}

// GetAzureFirewalls returns mock Azure Firewall data
func (c *MockAzureClient) GetAzureFirewalls(ctx context.Context, resourceGroup string) ([]models.AzureFirewall, error) {
	return []models.AzureFirewall{
		{
			ID:               "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/azureFirewalls/fw-hub",
			Name:             "fw-hub",
			ResourceGroup:    resourceGroup,
			Location:         "eastus",
			SKU:              "Premium",
			SubnetID:         "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/AzureFirewallSubnet",
			PrivateIPAddress: "10.0.0.4",
			PublicIPAddresses: []string{
				"/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/publicIPAddresses/pip-firewall",
			},
			ThreatIntelMode:   "Alert",
			DNSProxyEnabled:   true,
			ProvisioningState: "Succeeded",
		},
	}, nil
}

// GetApplicationGateways returns mock application gateway data
func (c *MockAzureClient) GetApplicationGateways(ctx context.Context, resourceGroup string) ([]models.ApplicationGateway, error) {
	return []models.ApplicationGateway{
		{
			ID:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/applicationGateways/appgw-web",
			Name:          "appgw-web",
			ResourceGroup: resourceGroup,
			Location:      "eastus",
			SKU:           "WAF_v2",
			Tier:          "WAF_v2",
			Capacity:      2,
			SubnetID:      "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/vnet-hub/subnets/subnet-appgw",
			WAFEnabled:    true,
			WAFMode:       "Prevention",
			FrontendIPConfigs: []models.AppGWFrontendIPConfig{
				{
					Name:              "appGwPublicFrontendIp",
					PrivateIPAddress:  "",
					PublicIPAddressID: "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/publicIPAddresses/pip-appgw",
				},
			},
			FrontendPorts: []models.AppGWFrontendPort{
				{Name: "port_80", Port: 80},
				{Name: "port_443", Port: 443},
			},
			BackendAddressPools: []models.AppGWBackendAddressPool{
				{
					Name:             "backend-pool",
					BackendAddresses: []string{"10.0.1.10", "10.0.1.11"},
				},
			},
			BackendHTTPSettings: []models.AppGWBackendHTTPSettings{
				{
					Name:                "http-settings",
					Port:                80,
					Protocol:            "Http",
					CookieBasedAffinity: "Disabled",
					RequestTimeout:      30,
					ProbeName:           "health-probe",
				},
			},
			HTTPListeners: []models.AppGWHTTPListener{
				{
					Name:             "http-listener",
					FrontendIPConfig: "appGwPublicFrontendIp",
					FrontendPort:     "port_80",
					Protocol:         "Http",
					HostName:         "",
				},
			},
			RequestRoutingRules: []models.AppGWRequestRoutingRule{
				{
					Name:                "rule-basic",
					RuleType:            "Basic",
					HTTPListener:        "http-listener",
					BackendAddressPool:  "backend-pool",
					BackendHTTPSettings: "http-settings",
					Priority:            100,
				},
			},
			Probes: []models.AppGWProbe{
				{
					Name:               "health-probe",
					Protocol:           "Http",
					Host:               "localhost",
					Path:               "/health",
					Interval:           30,
					Timeout:            30,
					UnhealthyThreshold: 3,
				},
			},
		},
	}, nil
}

// GetNetworkWatcherInsights returns mock Network Watcher insights
func (c *MockAzureClient) GetNetworkWatcherInsights(ctx context.Context, resourceGroup string) (*models.NetworkWatcherInsights, error) {
	return &models.NetworkWatcherInsights{
		FlowLogsEnabled: true,
		FlowLogs: []models.FlowLog{
			{
				ID:               "/subscriptions/" + c.subscriptionID + "/resourceGroups/NetworkWatcherRG/providers/Microsoft.Network/networkWatchers/NetworkWatcher_eastus/flowLogs/fl-nsg-web",
				NSGId:            "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/networkSecurityGroups/nsg-web",
				StorageAccountID: "/subscriptions/" + c.subscriptionID + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Storage/storageAccounts/stflowlogs",
				Enabled:          true,
				RetentionDays:    30,
				TrafficAnalytics: true,
			},
		},
		ConnectionMonitors: []models.ConnectionMonitor{
			{
				Name:             "monitor-web-to-db",
				Source:           "vm-web-01",
				Destination:      "10.0.2.10",
				MonitoringStatus: "Running",
			},
		},
		PacketCaptures: []models.PacketCapture{},
	}, nil
}

// GenerateMockTopology generates a complete mock topology for testing
func GenerateMockTopology(subscriptionID, resourceGroup string) *models.NetworkTopology {
	client := NewMockAzureClient(subscriptionID)
	ctx := context.Background()

	vnets, _ := client.GetVirtualNetworks(ctx, resourceGroup)
	nsgs, _ := client.GetNetworkSecurityGroups(ctx, resourceGroup)
	privateEndpoints, _ := client.GetPrivateEndpoints(ctx, resourceGroup)
	dnsZones, _ := client.GetPrivateDNSZones(ctx, resourceGroup)
	routeTables, _ := client.GetRouteTables(ctx, resourceGroup)
	natGateways, _ := client.GetNATGateways(ctx, resourceGroup)
	vpnGateways, _ := client.GetVPNGateways(ctx, resourceGroup)
	erCircuits, _ := client.GetExpressRouteCircuits(ctx, resourceGroup)
	loadBalancers, _ := client.GetLoadBalancers(ctx, resourceGroup)
	appGateways, _ := client.GetApplicationGateways(ctx, resourceGroup)
	azureFirewalls, _ := client.GetAzureFirewalls(ctx, resourceGroup)
	nwInsights, _ := client.GetNetworkWatcherInsights(ctx, resourceGroup)

	return &models.NetworkTopology{
		SubscriptionID:   subscriptionID,
		ResourceGroup:    resourceGroup,
		VirtualNetworks:  vnets,
		NSGs:             nsgs,
		PrivateEndpoints: privateEndpoints,
		PrivateDNSZones:  dnsZones,
		RouteTables:      routeTables,
		NATGateways:      natGateways,
		VPNGateways:      vpnGateways,
		ERCircuits:       erCircuits,
		LoadBalancers:    loadBalancers,
		AppGateways:      appGateways,
		AzureFirewalls:   azureFirewalls,
		NetworkWatcher:   nwInsights,
		Timestamp:        time.Now(),
	}
}
