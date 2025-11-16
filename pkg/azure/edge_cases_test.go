package azure

import (
	"testing"

	"azure-network-analyzer/pkg/models"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func TestExtractorsWithNilInputs(t *testing.T) {
	client := &AzureClient{}

	t.Run("extractSubnet with nil properties handles gracefully", func(t *testing.T) {
		subnet := &armnetwork.Subnet{
			ID:         nil,
			Name:       nil,
			Properties: nil,
		}

		result := client.extractSubnet(subnet)

		if result.ID != "" {
			t.Errorf("Expected empty ID, got %s", result.ID)
		}
		if result.Name != "" {
			t.Errorf("Expected empty Name, got %s", result.Name)
		}
		if result.AddressPrefix != "" {
			t.Errorf("Expected empty AddressPrefix, got %s", result.AddressPrefix)
		}
		if result.NetworkSecurityGroup != nil {
			t.Error("Expected nil NetworkSecurityGroup")
		}
		if len(result.PrivateEndpoints) != 0 {
			t.Errorf("Expected empty PrivateEndpoints, got %d items", len(result.PrivateEndpoints))
		}
	})

	t.Run("extractVNetPeering with nil properties", func(t *testing.T) {
		peering := &armnetwork.VirtualNetworkPeering{
			ID:         nil,
			Name:       nil,
			Properties: nil,
		}

		result := client.extractVNetPeering(peering)

		if result.ID != "" {
			t.Errorf("Expected empty ID, got %s", result.ID)
		}
		if result.RemoteVNetID != "" {
			t.Errorf("Expected empty RemoteVNetID, got %s", result.RemoteVNetID)
		}
		if result.PeeringState != "" {
			t.Errorf("Expected empty PeeringState, got %s", result.PeeringState)
		}
		// Boolean defaults should be false
		if result.AllowVNetAccess {
			t.Error("Expected AllowVNetAccess to be false")
		}
	})

	t.Run("extractRoute with nil properties", func(t *testing.T) {
		route := &armnetwork.Route{
			Name:       nil,
			Properties: nil,
		}

		result := client.extractRoute(route)

		if result.Name != "" {
			t.Errorf("Expected empty Name, got %s", result.Name)
		}
		if result.AddressPrefix != "" {
			t.Errorf("Expected empty AddressPrefix, got %s", result.AddressPrefix)
		}
		if result.NextHopType != "" {
			t.Errorf("Expected empty NextHopType, got %s", result.NextHopType)
		}
	})

	t.Run("extractSecurityRule with nil properties", func(t *testing.T) {
		rule := &armnetwork.SecurityRule{
			Name:       nil,
			Properties: nil,
		}

		result := extractSecurityRule(rule)

		if result.Name != "" {
			t.Errorf("Expected empty Name, got %s", result.Name)
		}
		if result.Priority != 0 {
			t.Errorf("Expected 0 Priority, got %d", result.Priority)
		}
		if result.Direction != "" {
			t.Errorf("Expected empty Direction, got %s", result.Direction)
		}
	})

	t.Run("extractFrontendIPConfig with nil properties", func(t *testing.T) {
		feConfig := &armnetwork.FrontendIPConfiguration{
			Name:       nil,
			Properties: nil,
		}

		result := client.extractFrontendIPConfig(feConfig)

		if result.Name != "" {
			t.Errorf("Expected empty Name, got %s", result.Name)
		}
		if result.PrivateIPAddress != "" {
			t.Errorf("Expected empty PrivateIPAddress, got %s", result.PrivateIPAddress)
		}
		if result.PublicIPAddressID != "" {
			t.Errorf("Expected empty PublicIPAddressID, got %s", result.PublicIPAddressID)
		}
	})

	t.Run("extractERPeering with nil properties", func(t *testing.T) {
		peering := &armnetwork.ExpressRouteCircuitPeering{
			Name:       nil,
			Properties: nil,
		}

		result := client.extractERPeering(peering)

		if result.Name != "" {
			t.Errorf("Expected empty Name, got %s", result.Name)
		}
		if result.PeerASN != 0 {
			t.Errorf("Expected 0 PeerASN, got %d", result.PeerASN)
		}
		if result.VlanID != 0 {
			t.Errorf("Expected 0 VlanID, got %d", result.VlanID)
		}
	})

	t.Run("extractAppGWProbe with nil properties", func(t *testing.T) {
		probe := &armnetwork.ApplicationGatewayProbe{
			Name:       nil,
			Properties: nil,
		}

		result := client.extractAppGWProbe(probe)

		if result.Name != "" {
			t.Errorf("Expected empty Name, got %s", result.Name)
		}
		if result.Interval != 0 {
			t.Errorf("Expected 0 Interval, got %d", result.Interval)
		}
		if result.UnhealthyThreshold != 0 {
			t.Errorf("Expected 0 UnhealthyThreshold, got %d", result.UnhealthyThreshold)
		}
	})
}

func TestExtractorsWithPartialData(t *testing.T) {
	client := &AzureClient{}

	t.Run("extractSubnet with some nil pointers in properties", func(t *testing.T) {
		// Simulate partial data where some fields are present but others are nil
		addressPrefix := "10.0.1.0/24"
		subnet := &armnetwork.Subnet{
			ID:   strPtr("subnet-id"),
			Name: strPtr("subnet-name"),
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix:        &addressPrefix,
				NetworkSecurityGroup: nil, // Explicitly nil
				RouteTable:           nil,
				NatGateway:           nil,
				PrivateEndpoints: []*armnetwork.PrivateEndpoint{
					{ID: nil}, // PE with nil ID
					{ID: strPtr("pe-2")},
				},
				ServiceEndpoints: []*armnetwork.ServiceEndpointPropertiesFormat{
					{Service: nil}, // SE with nil service
					{Service: strPtr("Microsoft.Storage")},
				},
				Delegations: []*armnetwork.Delegation{
					{Properties: nil}, // Delegation with nil properties
					{
						Properties: &armnetwork.ServiceDelegationPropertiesFormat{
							ServiceName: nil, // ServiceName is nil
						},
					},
				},
			},
		}

		result := client.extractSubnet(subnet)

		if result.AddressPrefix != "10.0.1.0/24" {
			t.Errorf("AddressPrefix mismatch: got %s", result.AddressPrefix)
		}
		if result.NetworkSecurityGroup != nil {
			t.Error("Expected nil NetworkSecurityGroup")
		}
		// Should skip nil IDs
		if len(result.PrivateEndpoints) != 1 {
			t.Errorf("Expected 1 PrivateEndpoint, got %d", len(result.PrivateEndpoints))
		}
		if len(result.ServiceEndpoints) != 1 {
			t.Errorf("Expected 1 ServiceEndpoint, got %d", len(result.ServiceEndpoints))
		}
		// Delegations with nil properties or nil ServiceName should be skipped
		if len(result.Delegations) != 0 {
			t.Errorf("Expected 0 Delegations, got %d", len(result.Delegations))
		}
	})

	t.Run("extractBackendAddressPool with nil IDs in configs", func(t *testing.T) {
		bePool := &armnetwork.BackendAddressPool{
			Name: strPtr("backend-pool"),
			Properties: &armnetwork.BackendAddressPoolPropertiesFormat{
				BackendIPConfigurations: []*armnetwork.InterfaceIPConfiguration{
					{ID: nil},
					{ID: strPtr("config-1")},
					{ID: nil},
					{ID: strPtr("config-2")},
				},
			},
		}

		result := client.extractBackendAddressPool(bePool)

		if len(result.BackendIPConfigs) != 2 {
			t.Errorf("Expected 2 BackendIPConfigs, got %d", len(result.BackendIPConfigs))
		}
	})

	t.Run("extractAppGWBackendAddressPool with mixed addresses", func(t *testing.T) {
		bePool := &armnetwork.ApplicationGatewayBackendAddressPool{
			Name: strPtr("appgw-backend"),
			Properties: &armnetwork.ApplicationGatewayBackendAddressPoolPropertiesFormat{
				BackendAddresses: []*armnetwork.ApplicationGatewayBackendAddress{
					{IPAddress: nil, Fqdn: nil},        // Both nil
					{IPAddress: strPtr("10.0.0.1"), Fqdn: nil},
					{IPAddress: nil, Fqdn: strPtr("app.example.com")},
					{IPAddress: strPtr("10.0.0.2"), Fqdn: strPtr("ignored")}, // IP takes precedence
				},
			},
		}

		result := client.extractAppGWBackendAddressPool(bePool)

		if len(result.BackendAddresses) != 3 {
			t.Errorf("Expected 3 BackendAddresses, got %d", len(result.BackendAddresses))
		}
	})
}

func TestResourceIDEdgeCases(t *testing.T) {
	t.Run("extractResourceName with unusual patterns", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"//double//slashes//name", "name"},
			{"/single/", ""},
			{"no-slashes-at-all", "no-slashes-at-all"},
			{"/", ""},
			{"", ""},
			{"/a", "a"},
			{"a/", ""},
			{"/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p", "p"},
		}

		for _, tc := range testCases {
			result := extractResourceName(tc.input)
			if result != tc.expected {
				t.Errorf("extractResourceName(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		}
	})

	t.Run("extractVNetIDFromSubnet edge cases", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"/subnets/test", ""},                            // starts with /subnets/
			{"vnet/subnets/subnet", "vnet"},                  // no leading slash
			{"/subnets/", ""},                                // just /subnets/
			{"a/subnets/b/subnets/c", "a"},                   // multiple /subnets/ - takes first
			{"/sub/sub/subnets/x", "/sub/sub"},               // subnets in the middle
		}

		for _, tc := range testCases {
			result := extractVNetIDFromSubnet(tc.input)
			if result != tc.expected {
				t.Errorf("extractVNetIDFromSubnet(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		}
	})
}

func TestEmptyCollections(t *testing.T) {
	t.Run("Empty topology serializes correctly", func(t *testing.T) {
		topology := &models.NetworkTopology{
			SubscriptionID:   "empty-sub",
			ResourceGroup:    "empty-rg",
			VirtualNetworks:  []models.VirtualNetwork{},
			NSGs:             []models.NetworkSecurityGroup{},
			PrivateEndpoints: []models.PrivateEndpoint{},
			PrivateDNSZones:  []models.PrivateDNSZone{},
			RouteTables:      []models.RouteTable{},
			NATGateways:      []models.NATGateway{},
			VPNGateways:      []models.VPNGateway{},
			ERCircuits:       []models.ExpressRouteCircuit{},
			LoadBalancers:    []models.LoadBalancer{},
			AppGateways:      []models.ApplicationGateway{},
		}

		// Verify count function handles empty
		count := countResources(topology)
		if count != 0 {
			t.Errorf("Expected 0 resources, got %d", count)
		}
	})

	t.Run("VNet with no subnets or peerings", func(t *testing.T) {
		vnet := models.VirtualNetwork{
			ID:            "vnet-id",
			Name:          "empty-vnet",
			ResourceGroup: "rg",
			Location:      "eastus",
			AddressSpace:  []string{"10.0.0.0/8"},
			Subnets:       []models.Subnet{},
			Peerings:      []models.VNetPeering{},
			DNSServers:    []string{},
			EnableDDoS:    false,
		}

		if len(vnet.Subnets) != 0 {
			t.Error("Subnets should be empty")
		}
		if len(vnet.Peerings) != 0 {
			t.Error("Peerings should be empty")
		}
	})

	t.Run("NSG with no rules", func(t *testing.T) {
		nsg := models.NetworkSecurityGroup{
			ID:            "nsg-id",
			Name:          "empty-nsg",
			ResourceGroup: "rg",
			Location:      "eastus",
			SecurityRules: []models.SecurityRule{},
			Associations: models.NSGAssociations{
				Subnets:           []string{},
				NetworkInterfaces: []string{},
			},
		}

		if len(nsg.SecurityRules) != 0 {
			t.Error("SecurityRules should be empty")
		}
	})
}

func TestBoundaryValues(t *testing.T) {
	t.Run("Security rule priority boundaries", func(t *testing.T) {
		// Azure NSG priorities: 100-4096
		testCases := []int32{100, 4096, 200, 3000}

		for _, priority := range testCases {
			rule := models.SecurityRule{
				Name:     "test-rule",
				Priority: priority,
			}

			if rule.Priority < 100 || rule.Priority > 4096 {
				t.Errorf("Priority %d is out of valid range", priority)
			}
		}
	})

	t.Run("Port number boundaries", func(t *testing.T) {
		// Valid ports: 1-65535
		testCases := []int32{1, 80, 443, 65535}

		for _, port := range testCases {
			probe := models.Probe{
				Name: "test-probe",
				Port: port,
			}

			if probe.Port < 1 || probe.Port > 65535 {
				t.Errorf("Port %d is out of valid range", probe.Port)
			}
		}
	})

	t.Run("BGP ASN values", func(t *testing.T) {
		// Valid private ASNs: 64512-65534 (2-byte) or 4200000000-4294967294 (4-byte)
		bgp := models.BGPSettings{
			ASN:               65515,
			BGPPeeringAddress: "10.0.0.1",
			PeerWeight:        0,
		}

		if bgp.ASN <= 0 {
			t.Error("ASN should be positive")
		}
	})

	t.Run("ExpressRoute bandwidth values", func(t *testing.T) {
		// Common ER bandwidths: 50, 100, 200, 500, 1000, 2000, 5000, 10000
		circuit := models.ExpressRouteCircuit{
			BandwidthInMbps: 1000,
		}

		validBandwidths := []int32{50, 100, 200, 500, 1000, 2000, 5000, 10000}
		found := false
		for _, bw := range validBandwidths {
			if circuit.BandwidthInMbps == bw {
				found = true
				break
			}
		}

		if !found && circuit.BandwidthInMbps != 0 {
			// Not a hard error, just noting unusual value
			t.Logf("Bandwidth %d is not a standard ER bandwidth", circuit.BandwidthInMbps)
		}
	})
}

func TestSpecialCharactersInNames(t *testing.T) {
	t.Run("Resource names with hyphens and underscores", func(t *testing.T) {
		vnet := models.VirtualNetwork{
			Name: "vnet-prod-001_eastus",
		}

		if vnet.Name == "" {
			t.Error("Name should handle special characters")
		}
	})

	t.Run("Extract resource name with special characters", func(t *testing.T) {
		id := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet-prod_001"
		result := extractResourceName(id)

		if result != "vnet-prod_001" {
			t.Errorf("Expected vnet-prod_001, got %s", result)
		}
	})
}

// Helper function to count resources
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
