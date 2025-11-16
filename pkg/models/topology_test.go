package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNetworkTopologyJSON(t *testing.T) {
	// Create a sample topology
	topology := NetworkTopology{
		SubscriptionID: "test-subscription-123",
		ResourceGroup:  "test-rg",
		Timestamp:      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		VirtualNetworks: []VirtualNetwork{
			{
				ID:            "/subscriptions/test-sub/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1",
				Name:          "vnet1",
				ResourceGroup: "test-rg",
				Location:      "eastus",
				AddressSpace:  []string{"10.0.0.0/16"},
				DNSServers:    []string{"10.0.0.4", "10.0.0.5"},
				EnableDDoS:    true,
				Subnets: []Subnet{
					{
						ID:            "/subscriptions/test-sub/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
						Name:          "subnet1",
						AddressPrefix: "10.0.1.0/24",
						ServiceEndpoints: []string{"Microsoft.Storage", "Microsoft.KeyVault"},
						Delegations:      []string{},
						PrivateEndpoints: []string{},
					},
				},
				Peerings: []VNetPeering{
					{
						ID:                    "/subscriptions/test-sub/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/virtualNetworkPeerings/peering1",
						Name:                  "peering1",
						RemoteVNetID:          "/subscriptions/test-sub/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet2",
						RemoteVNetName:        "vnet2",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: false,
						AllowGatewayTransit:   false,
						UseRemoteGateways:     false,
					},
				},
			},
		},
		NSGs: []NetworkSecurityGroup{
			{
				ID:            "/subscriptions/test-sub/resourceGroups/test-rg/providers/Microsoft.Network/networkSecurityGroups/nsg1",
				Name:          "nsg1",
				ResourceGroup: "test-rg",
				Location:      "eastus",
				SecurityRules: []SecurityRule{
					{
						Name:                     "AllowSSH",
						Priority:                 100,
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "TCP",
						SourceAddressPrefix:      "10.0.0.0/8",
						SourcePortRange:          "*",
						DestinationAddressPrefix: "*",
						DestinationPortRange:     "22",
						Description:              "Allow SSH from internal network",
					},
				},
				Associations: NSGAssociations{
					Subnets:           []string{"/subscriptions/test-sub/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"},
					NetworkInterfaces: []string{},
				},
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.MarshalIndent(topology, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal topology to JSON: %v", err)
	}

	// Verify JSON is not empty
	if len(jsonData) == 0 {
		t.Error("JSON output is empty")
	}

	// Test JSON unmarshaling
	var decoded NetworkTopology
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal topology from JSON: %v", err)
	}

	// Verify key fields
	if decoded.SubscriptionID != topology.SubscriptionID {
		t.Errorf("SubscriptionID mismatch: got %s, want %s", decoded.SubscriptionID, topology.SubscriptionID)
	}

	if decoded.ResourceGroup != topology.ResourceGroup {
		t.Errorf("ResourceGroup mismatch: got %s, want %s", decoded.ResourceGroup, topology.ResourceGroup)
	}

	if len(decoded.VirtualNetworks) != len(topology.VirtualNetworks) {
		t.Errorf("VirtualNetworks count mismatch: got %d, want %d", len(decoded.VirtualNetworks), len(topology.VirtualNetworks))
	}

	if len(decoded.NSGs) != len(topology.NSGs) {
		t.Errorf("NSGs count mismatch: got %d, want %d", len(decoded.NSGs), len(topology.NSGs))
	}
}

func TestVirtualNetworkJSON(t *testing.T) {
	vnet := VirtualNetwork{
		ID:            "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
		Name:          "test-vnet",
		ResourceGroup: "test-rg",
		Location:      "westus2",
		AddressSpace:  []string{"10.0.0.0/16", "172.16.0.0/12"},
		DNSServers:    []string{},
		EnableDDoS:    false,
		Subnets:       []Subnet{},
		Peerings:      []VNetPeering{},
	}

	jsonData, err := json.Marshal(vnet)
	if err != nil {
		t.Fatalf("Failed to marshal VirtualNetwork: %v", err)
	}

	var decoded VirtualNetwork
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal VirtualNetwork: %v", err)
	}

	if decoded.Name != vnet.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, vnet.Name)
	}

	if len(decoded.AddressSpace) != len(vnet.AddressSpace) {
		t.Errorf("AddressSpace count mismatch: got %d, want %d", len(decoded.AddressSpace), len(vnet.AddressSpace))
	}
}

func TestSecurityRuleJSON(t *testing.T) {
	rule := SecurityRule{
		Name:                     "DenyAllInbound",
		Priority:                 4096,
		Direction:                "Inbound",
		Access:                   "Deny",
		Protocol:                 "*",
		SourceAddressPrefix:      "*",
		SourcePortRange:          "*",
		DestinationAddressPrefix: "*",
		DestinationPortRange:     "*",
		Description:              "Deny all inbound traffic",
	}

	jsonData, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Failed to marshal SecurityRule: %v", err)
	}

	var decoded SecurityRule
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal SecurityRule: %v", err)
	}

	if decoded.Priority != rule.Priority {
		t.Errorf("Priority mismatch: got %d, want %d", decoded.Priority, rule.Priority)
	}

	if decoded.Direction != rule.Direction {
		t.Errorf("Direction mismatch: got %s, want %s", decoded.Direction, rule.Direction)
	}

	if decoded.Access != rule.Access {
		t.Errorf("Access mismatch: got %s, want %s", decoded.Access, rule.Access)
	}
}

func TestSubnetWithOptionalFields(t *testing.T) {
	nsgID := "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/nsg1"
	routeTableID := "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/routeTables/rt1"

	subnet := Subnet{
		ID:                   "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
		Name:                 "subnet-with-nsg",
		AddressPrefix:        "10.0.1.0/24",
		NetworkSecurityGroup: &nsgID,
		RouteTable:           &routeTableID,
		NATGateway:           nil,
		PrivateEndpoints:     []string{"pe1", "pe2"},
		ServiceEndpoints:     []string{"Microsoft.Storage"},
		Delegations:          []string{"Microsoft.Web/serverFarms"},
	}

	jsonData, err := json.Marshal(subnet)
	if err != nil {
		t.Fatalf("Failed to marshal Subnet: %v", err)
	}

	var decoded Subnet
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Subnet: %v", err)
	}

	if decoded.NetworkSecurityGroup == nil {
		t.Error("NetworkSecurityGroup should not be nil")
	} else if *decoded.NetworkSecurityGroup != nsgID {
		t.Errorf("NetworkSecurityGroup mismatch: got %s, want %s", *decoded.NetworkSecurityGroup, nsgID)
	}

	if decoded.RouteTable == nil {
		t.Error("RouteTable should not be nil")
	}

	if decoded.NATGateway != nil {
		t.Error("NATGateway should be nil")
	}

	if len(decoded.PrivateEndpoints) != 2 {
		t.Errorf("PrivateEndpoints count mismatch: got %d, want 2", len(decoded.PrivateEndpoints))
	}
}

func TestLoadBalancerJSON(t *testing.T) {
	lb := LoadBalancer{
		ID:            "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/loadBalancers/lb1",
		Name:          "test-lb",
		ResourceGroup: "test-rg",
		Location:      "eastus",
		SKU:           "Standard",
		Type:          "Public",
		FrontendIPConfigs: []FrontendIPConfig{
			{
				Name:              "frontend1",
				PrivateIPAddress:  "",
				PublicIPAddressID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/publicIPAddresses/pip1",
				SubnetID:          "",
			},
		},
		BackendAddressPools: []BackendAddressPool{
			{
				Name:             "backend1",
				BackendIPConfigs: []string{"nic1", "nic2"},
			},
		},
		LoadBalancingRules: []LoadBalancingRule{
			{
				Name:               "rule1",
				Protocol:           "TCP",
				FrontendPort:       80,
				BackendPort:        80,
				EnableFloatingIP:   false,
				IdleTimeoutMinutes: 4,
				LoadDistribution:   "Default",
			},
		},
		Probes: []Probe{
			{
				Name:              "probe1",
				Protocol:          "HTTP",
				Port:              80,
				IntervalInSeconds: 15,
				NumberOfProbes:    2,
				RequestPath:       "/health",
			},
		},
		InboundNATRules: []InboundNATRule{},
	}

	jsonData, err := json.Marshal(lb)
	if err != nil {
		t.Fatalf("Failed to marshal LoadBalancer: %v", err)
	}

	var decoded LoadBalancer
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal LoadBalancer: %v", err)
	}

	if decoded.SKU != lb.SKU {
		t.Errorf("SKU mismatch: got %s, want %s", decoded.SKU, lb.SKU)
	}

	if decoded.Type != lb.Type {
		t.Errorf("Type mismatch: got %s, want %s", decoded.Type, lb.Type)
	}

	if len(decoded.LoadBalancingRules) != 1 {
		t.Errorf("LoadBalancingRules count mismatch: got %d, want 1", len(decoded.LoadBalancingRules))
	}

	if decoded.LoadBalancingRules[0].FrontendPort != 80 {
		t.Errorf("FrontendPort mismatch: got %d, want 80", decoded.LoadBalancingRules[0].FrontendPort)
	}
}

func TestApplicationGatewayJSON(t *testing.T) {
	appGW := ApplicationGateway{
		ID:            "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/applicationGateways/appgw1",
		Name:          "test-appgw",
		ResourceGroup: "test-rg",
		Location:      "eastus",
		SKU:           "WAF_v2",
		Tier:          "WAF_v2",
		Capacity:      2,
		SubnetID:      "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/appgw-subnet",
		WAFEnabled:    true,
		WAFMode:       "Prevention",
		FrontendIPConfigs: []AppGWFrontendIPConfig{
			{
				Name:              "appGwPublicFrontendIp",
				PrivateIPAddress:  "",
				PublicIPAddressID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/publicIPAddresses/appgw-pip",
			},
		},
		FrontendPorts: []AppGWFrontendPort{
			{Name: "port_80", Port: 80},
			{Name: "port_443", Port: 443},
		},
		BackendAddressPools:   []AppGWBackendAddressPool{},
		BackendHTTPSettings:   []AppGWBackendHTTPSettings{},
		HTTPListeners:         []AppGWHTTPListener{},
		RequestRoutingRules:   []AppGWRequestRoutingRule{},
		Probes:                []AppGWProbe{},
	}

	jsonData, err := json.Marshal(appGW)
	if err != nil {
		t.Fatalf("Failed to marshal ApplicationGateway: %v", err)
	}

	var decoded ApplicationGateway
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ApplicationGateway: %v", err)
	}

	if !decoded.WAFEnabled {
		t.Error("WAFEnabled should be true")
	}

	if decoded.WAFMode != "Prevention" {
		t.Errorf("WAFMode mismatch: got %s, want Prevention", decoded.WAFMode)
	}

	if decoded.Capacity != 2 {
		t.Errorf("Capacity mismatch: got %d, want 2", decoded.Capacity)
	}

	if len(decoded.FrontendPorts) != 2 {
		t.Errorf("FrontendPorts count mismatch: got %d, want 2", len(decoded.FrontendPorts))
	}
}

func TestExpressRouteCircuitJSON(t *testing.T) {
	er := ExpressRouteCircuit{
		ID:                       "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/er1",
		Name:                     "test-expressroute",
		ResourceGroup:            "test-rg",
		Location:                 "eastus",
		ServiceProviderName:      "Equinix",
		PeeringLocation:          "Washington DC",
		BandwidthInMbps:          1000,
		SKUTier:                  "Premium",
		SKUFamily:                "MeteredData",
		CircuitProvisioningState: "Enabled",
		Peerings: []ERPeering{
			{
				Name:                       "AzurePrivatePeering",
				PeeringType:                "AzurePrivatePeering",
				State:                      "Enabled",
				AzureASN:                   12076,
				PeerASN:                    65001,
				PrimaryPeerAddressPrefix:   "192.168.1.0/30",
				SecondaryPeerAddressPrefix: "192.168.2.0/30",
				VlanID:                     100,
			},
		},
		Authorizations: []ERAuthorization{
			{
				Name:                   "auth1",
				AuthorizationKey:       true,
				AuthorizationUseStatus: "Available",
			},
		},
	}

	jsonData, err := json.Marshal(er)
	if err != nil {
		t.Fatalf("Failed to marshal ExpressRouteCircuit: %v", err)
	}

	var decoded ExpressRouteCircuit
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ExpressRouteCircuit: %v", err)
	}

	if decoded.BandwidthInMbps != 1000 {
		t.Errorf("BandwidthInMbps mismatch: got %d, want 1000", decoded.BandwidthInMbps)
	}

	if len(decoded.Peerings) != 1 {
		t.Errorf("Peerings count mismatch: got %d, want 1", len(decoded.Peerings))
	}

	if decoded.Peerings[0].PeerASN != 65001 {
		t.Errorf("PeerASN mismatch: got %d, want 65001", decoded.Peerings[0].PeerASN)
	}
}

func TestNetworkWatcherInsightsJSON(t *testing.T) {
	insights := NetworkWatcherInsights{
		FlowLogsEnabled: true,
		FlowLogs: []FlowLog{
			{
				ID:               "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkWatchers/nw1/flowLogs/fl1",
				NSGId:            "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
				StorageAccountID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/sa1",
				Enabled:          true,
				RetentionDays:    30,
				TrafficAnalytics: true,
			},
		},
		ConnectionMonitors: []ConnectionMonitor{
			{
				Name:             "monitor1",
				Source:           "vm1",
				Destination:      "10.0.0.5",
				MonitoringStatus: "Running",
			},
		},
		PacketCaptures: []PacketCapture{
			{
				Name:   "capture1",
				Target: "vm1",
				Status: "Succeeded",
			},
		},
	}

	jsonData, err := json.Marshal(insights)
	if err != nil {
		t.Fatalf("Failed to marshal NetworkWatcherInsights: %v", err)
	}

	var decoded NetworkWatcherInsights
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal NetworkWatcherInsights: %v", err)
	}

	if !decoded.FlowLogsEnabled {
		t.Error("FlowLogsEnabled should be true")
	}

	if len(decoded.FlowLogs) != 1 {
		t.Errorf("FlowLogs count mismatch: got %d, want 1", len(decoded.FlowLogs))
	}

	if decoded.FlowLogs[0].RetentionDays != 30 {
		t.Errorf("RetentionDays mismatch: got %d, want 30", decoded.FlowLogs[0].RetentionDays)
	}
}
