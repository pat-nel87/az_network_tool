package azure

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"azure-network-analyzer/pkg/models"
)

func TestMockClientGeneratesValidData(t *testing.T) {
	client := NewMockAzureClient("test-sub")
	ctx := context.Background()
	resourceGroup := "test-rg"

	t.Run("VNets have valid structure", func(t *testing.T) {
		vnets, err := client.GetVirtualNetworks(ctx, resourceGroup)
		if err != nil {
			t.Fatalf("GetVirtualNetworks failed: %v", err)
		}

		if len(vnets) == 0 {
			t.Fatal("Expected at least one VNet")
		}

		for _, vnet := range vnets {
			// Check required fields
			if vnet.ID == "" {
				t.Error("VNet ID is empty")
			}
			if vnet.Name == "" {
				t.Error("VNet Name is empty")
			}
			if vnet.ResourceGroup != resourceGroup {
				t.Errorf("VNet ResourceGroup mismatch: got %s, want %s", vnet.ResourceGroup, resourceGroup)
			}
			if vnet.Location == "" {
				t.Error("VNet Location is empty")
			}
			if len(vnet.AddressSpace) == 0 {
				t.Error("VNet has no address space")
			}

			// Validate address space format (basic CIDR check)
			for _, addr := range vnet.AddressSpace {
				if !strings.Contains(addr, "/") {
					t.Errorf("Invalid CIDR format: %s", addr)
				}
			}

			// Check subnets
			for _, subnet := range vnet.Subnets {
				if subnet.ID == "" {
					t.Error("Subnet ID is empty")
				}
				if subnet.Name == "" {
					t.Error("Subnet Name is empty")
				}
				if subnet.AddressPrefix == "" {
					t.Error("Subnet AddressPrefix is empty")
				}
				if !strings.Contains(subnet.AddressPrefix, "/") {
					t.Errorf("Invalid subnet CIDR format: %s", subnet.AddressPrefix)
				}

				// Verify subnet ID contains VNet ID
				if !strings.Contains(subnet.ID, vnet.Name) {
					t.Errorf("Subnet ID doesn't contain VNet name: %s", subnet.ID)
				}
			}

			// Check peerings
			for _, peering := range vnet.Peerings {
				if peering.ID == "" {
					t.Error("Peering ID is empty")
				}
				if peering.Name == "" {
					t.Error("Peering Name is empty")
				}
				if peering.RemoteVNetID == "" {
					t.Error("Peering RemoteVNetID is empty")
				}
				if peering.PeeringState == "" {
					t.Error("Peering State is empty")
				}
			}
		}
	})

	t.Run("NSGs have valid structure", func(t *testing.T) {
		nsgs, err := client.GetNetworkSecurityGroups(ctx, resourceGroup)
		if err != nil {
			t.Fatalf("GetNetworkSecurityGroups failed: %v", err)
		}

		if len(nsgs) == 0 {
			t.Fatal("Expected at least one NSG")
		}

		for _, nsg := range nsgs {
			if nsg.ID == "" {
				t.Error("NSG ID is empty")
			}
			if nsg.Name == "" {
				t.Error("NSG Name is empty")
			}
			if len(nsg.SecurityRules) == 0 {
				t.Error("NSG has no security rules")
			}

			// Check security rules
			for _, rule := range nsg.SecurityRules {
				if rule.Name == "" {
					t.Error("Rule Name is empty")
				}
				if rule.Priority < 100 || rule.Priority > 4096 {
					t.Errorf("Rule Priority out of range: %d", rule.Priority)
				}
				if rule.Direction != "Inbound" && rule.Direction != "Outbound" {
					t.Errorf("Invalid rule Direction: %s", rule.Direction)
				}
				if rule.Access != "Allow" && rule.Access != "Deny" {
					t.Errorf("Invalid rule Access: %s", rule.Access)
				}
			}
		}
	})

	t.Run("Private Endpoints have valid structure", func(t *testing.T) {
		endpoints, err := client.GetPrivateEndpoints(ctx, resourceGroup)
		if err != nil {
			t.Fatalf("GetPrivateEndpoints failed: %v", err)
		}

		for _, pe := range endpoints {
			if pe.ID == "" {
				t.Error("PrivateEndpoint ID is empty")
			}
			if pe.Name == "" {
				t.Error("PrivateEndpoint Name is empty")
			}
			if pe.SubnetID == "" {
				t.Error("PrivateEndpoint SubnetID is empty")
			}
			if pe.ConnectionState == "" {
				t.Error("PrivateEndpoint ConnectionState is empty")
			}
			if len(pe.GroupIDs) == 0 {
				t.Error("PrivateEndpoint has no GroupIDs")
			}
		}
	})

	t.Run("Route Tables have valid structure", func(t *testing.T) {
		routeTables, err := client.GetRouteTables(ctx, resourceGroup)
		if err != nil {
			t.Fatalf("GetRouteTables failed: %v", err)
		}

		for _, rt := range routeTables {
			if rt.ID == "" {
				t.Error("RouteTable ID is empty")
			}
			if rt.Name == "" {
				t.Error("RouteTable Name is empty")
			}

			for _, route := range rt.Routes {
				if route.Name == "" {
					t.Error("Route Name is empty")
				}
				if route.AddressPrefix == "" {
					t.Error("Route AddressPrefix is empty")
				}
				if route.NextHopType == "" {
					t.Error("Route NextHopType is empty")
				}
			}
		}
	})

	t.Run("Load Balancers have valid structure", func(t *testing.T) {
		lbs, err := client.GetLoadBalancers(ctx, resourceGroup)
		if err != nil {
			t.Fatalf("GetLoadBalancers failed: %v", err)
		}

		for _, lb := range lbs {
			if lb.ID == "" {
				t.Error("LoadBalancer ID is empty")
			}
			if lb.Name == "" {
				t.Error("LoadBalancer Name is empty")
			}
			if lb.SKU == "" {
				t.Error("LoadBalancer SKU is empty")
			}
			if lb.Type != "Public" && lb.Type != "Internal" {
				t.Errorf("Invalid LoadBalancer Type: %s", lb.Type)
			}

			// Check frontend configs
			if len(lb.FrontendIPConfigs) == 0 {
				t.Error("LoadBalancer has no FrontendIPConfigs")
			}

			// Check rules have valid ports
			for _, rule := range lb.LoadBalancingRules {
				if rule.FrontendPort <= 0 || rule.FrontendPort > 65535 {
					t.Errorf("Invalid FrontendPort: %d", rule.FrontendPort)
				}
				if rule.BackendPort <= 0 || rule.BackendPort > 65535 {
					t.Errorf("Invalid BackendPort: %d", rule.BackendPort)
				}
			}
		}
	})

	t.Run("VPN Gateways have valid structure", func(t *testing.T) {
		gws, err := client.GetVPNGateways(ctx, resourceGroup)
		if err != nil {
			t.Fatalf("GetVPNGateways failed: %v", err)
		}

		for _, gw := range gws {
			if gw.ID == "" {
				t.Error("VPNGateway ID is empty")
			}
			if gw.Name == "" {
				t.Error("VPNGateway Name is empty")
			}
			if gw.GatewayType == "" {
				t.Error("VPNGateway GatewayType is empty")
			}
			if gw.VpnType == "" {
				t.Error("VPNGateway VpnType is empty")
			}
			if gw.SKU == "" {
				t.Error("VPNGateway SKU is empty")
			}

			// Check BGP settings if present
			if gw.BGPSettings != nil {
				if gw.BGPSettings.ASN <= 0 {
					t.Error("BGP ASN should be positive")
				}
			}
		}
	})

	t.Run("Application Gateways have valid structure", func(t *testing.T) {
		appGws, err := client.GetApplicationGateways(ctx, resourceGroup)
		if err != nil {
			t.Fatalf("GetApplicationGateways failed: %v", err)
		}

		for _, ag := range appGws {
			if ag.ID == "" {
				t.Error("ApplicationGateway ID is empty")
			}
			if ag.Name == "" {
				t.Error("ApplicationGateway Name is empty")
			}
			if ag.SKU == "" {
				t.Error("ApplicationGateway SKU is empty")
			}
			if ag.Capacity <= 0 {
				t.Error("ApplicationGateway Capacity should be positive")
			}

			// WAF checks
			if ag.WAFEnabled && ag.WAFMode == "" {
				t.Error("WAF is enabled but mode is not set")
			}

			// Check frontend ports
			for _, port := range ag.FrontendPorts {
				if port.Port <= 0 || port.Port > 65535 {
					t.Errorf("Invalid FrontendPort: %d", port.Port)
				}
			}
		}
	})
}

func TestMockTopologyDataIntegrity(t *testing.T) {
	topology := GenerateMockTopology("test-sub", "test-rg")

	t.Run("VNet peerings are reciprocal", func(t *testing.T) {
		// Build a map of VNet names to their peerings
		peeringMap := make(map[string][]string)
		for _, vnet := range topology.VirtualNetworks {
			for _, peering := range vnet.Peerings {
				peeringMap[vnet.Name] = append(peeringMap[vnet.Name], peering.RemoteVNetName)
			}
		}

		// Check reciprocity
		for vnetName, peers := range peeringMap {
			for _, peerName := range peers {
				// Check if peer has this vnet as a peer
				peerPeers, exists := peeringMap[peerName]
				if !exists {
					t.Errorf("VNet %s peers with %s, but %s has no peerings", vnetName, peerName, peerName)
					continue
				}

				found := false
				for _, p := range peerPeers {
					if p == vnetName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("VNet %s peers with %s, but %s doesn't peer back", vnetName, peerName, peerName)
				}
			}
		}
	})

	t.Run("NSG associations point to valid subnets", func(t *testing.T) {
		// Build a set of all subnet IDs
		subnetIDs := make(map[string]bool)
		for _, vnet := range topology.VirtualNetworks {
			for _, subnet := range vnet.Subnets {
				subnetIDs[subnet.ID] = true
			}
		}

		// Check NSG associations
		for _, nsg := range topology.NSGs {
			for _, subnetID := range nsg.Associations.Subnets {
				if !subnetIDs[subnetID] {
					t.Errorf("NSG %s associated with non-existent subnet: %s", nsg.Name, subnetID)
				}
			}
		}
	})

	t.Run("Private Endpoints reference valid subnets", func(t *testing.T) {
		subnetIDs := make(map[string]bool)
		for _, vnet := range topology.VirtualNetworks {
			for _, subnet := range vnet.Subnets {
				subnetIDs[subnet.ID] = true
			}
		}

		for _, pe := range topology.PrivateEndpoints {
			if !subnetIDs[pe.SubnetID] {
				t.Errorf("PrivateEndpoint %s references non-existent subnet: %s", pe.Name, pe.SubnetID)
			}
		}
	})

	t.Run("Route table associations are valid", func(t *testing.T) {
		subnetIDs := make(map[string]bool)
		for _, vnet := range topology.VirtualNetworks {
			for _, subnet := range vnet.Subnets {
				subnetIDs[subnet.ID] = true
			}
		}

		for _, rt := range topology.RouteTables {
			for _, subnetID := range rt.AssociatedSubnets {
				if !subnetIDs[subnetID] {
					t.Errorf("RouteTable %s associated with non-existent subnet: %s", rt.Name, subnetID)
				}
			}
		}
	})

	t.Run("NAT Gateway associations are valid", func(t *testing.T) {
		subnetIDs := make(map[string]bool)
		for _, vnet := range topology.VirtualNetworks {
			for _, subnet := range vnet.Subnets {
				subnetIDs[subnet.ID] = true
			}
		}

		for _, nat := range topology.NATGateways {
			for _, subnetID := range nat.AssociatedSubnets {
				if !subnetIDs[subnetID] {
					t.Errorf("NATGateway %s associated with non-existent subnet: %s", nat.Name, subnetID)
				}
			}
		}
	})

	t.Run("Subnet NSG references match NSG associations", func(t *testing.T) {
		// Build NSG association map
		nsgToSubnets := make(map[string]map[string]bool)
		for _, nsg := range topology.NSGs {
			nsgToSubnets[nsg.ID] = make(map[string]bool)
			for _, subnetID := range nsg.Associations.Subnets {
				nsgToSubnets[nsg.ID][subnetID] = true
			}
		}

		// Check subnet references
		for _, vnet := range topology.VirtualNetworks {
			for _, subnet := range vnet.Subnets {
				if subnet.NetworkSecurityGroup != nil {
					nsgID := *subnet.NetworkSecurityGroup
					if subnets, exists := nsgToSubnets[nsgID]; exists {
						if !subnets[subnet.ID] {
							t.Errorf("Subnet %s references NSG %s, but NSG doesn't list this subnet", subnet.Name, extractResourceName(nsgID))
						}
					}
				}
			}
		}
	})

	t.Run("Resource IDs follow Azure format", func(t *testing.T) {
		// Check all resource IDs follow /subscriptions/... pattern
		checkID := func(id, resourceType string) {
			if id == "" {
				return
			}
			if !strings.HasPrefix(id, "/subscriptions/") {
				t.Errorf("%s ID doesn't follow Azure format: %s", resourceType, id)
			}
			if !strings.Contains(id, "/resourceGroups/") {
				t.Errorf("%s ID missing resourceGroups segment: %s", resourceType, id)
			}
		}

		for _, vnet := range topology.VirtualNetworks {
			checkID(vnet.ID, "VNet")
			for _, subnet := range vnet.Subnets {
				checkID(subnet.ID, "Subnet")
			}
		}

		for _, nsg := range topology.NSGs {
			checkID(nsg.ID, "NSG")
		}

		for _, pe := range topology.PrivateEndpoints {
			checkID(pe.ID, "PrivateEndpoint")
		}
	})
}

func TestMockTopologyJSONRoundTrip(t *testing.T) {
	original := GenerateMockTopology("test-sub", "test-rg")

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(original, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal topology: %v", err)
	}

	// Verify JSON is not empty and reasonably sized
	if len(jsonData) < 1000 {
		t.Error("JSON output seems too small")
	}

	// Deserialize back
	var decoded models.NetworkTopology
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal topology: %v", err)
	}

	// Verify key properties preserved
	if decoded.SubscriptionID != original.SubscriptionID {
		t.Errorf("SubscriptionID mismatch after round-trip")
	}

	if decoded.ResourceGroup != original.ResourceGroup {
		t.Errorf("ResourceGroup mismatch after round-trip")
	}

	if len(decoded.VirtualNetworks) != len(original.VirtualNetworks) {
		t.Errorf("VirtualNetworks count mismatch: got %d, want %d", len(decoded.VirtualNetworks), len(original.VirtualNetworks))
	}

	if len(decoded.NSGs) != len(original.NSGs) {
		t.Errorf("NSGs count mismatch: got %d, want %d", len(decoded.NSGs), len(original.NSGs))
	}

	if len(decoded.PrivateEndpoints) != len(original.PrivateEndpoints) {
		t.Errorf("PrivateEndpoints count mismatch: got %d, want %d", len(decoded.PrivateEndpoints), len(original.PrivateEndpoints))
	}

	// Deep check - verify first VNet's subnets are preserved
	if len(original.VirtualNetworks) > 0 && len(decoded.VirtualNetworks) > 0 {
		origVNet := original.VirtualNetworks[0]
		decVNet := decoded.VirtualNetworks[0]

		if len(decVNet.Subnets) != len(origVNet.Subnets) {
			t.Errorf("First VNet subnets count mismatch: got %d, want %d", len(decVNet.Subnets), len(origVNet.Subnets))
		}

		if len(decVNet.Peerings) != len(origVNet.Peerings) {
			t.Errorf("First VNet peerings count mismatch: got %d, want %d", len(decVNet.Peerings), len(origVNet.Peerings))
		}
	}

	// Check NSG rules preserved
	if len(original.NSGs) > 0 && len(decoded.NSGs) > 0 {
		origNSG := original.NSGs[0]
		decNSG := decoded.NSGs[0]

		if len(decNSG.SecurityRules) != len(origNSG.SecurityRules) {
			t.Errorf("First NSG rules count mismatch: got %d, want %d", len(decNSG.SecurityRules), len(origNSG.SecurityRules))
		}
	}
}

// Import the models package type for the round-trip test
type NetworkTopology = struct {
	SubscriptionID   string
	ResourceGroup    string
	VirtualNetworks  []struct {
		ID            string
		Name          string
		ResourceGroup string
		Location      string
		AddressSpace  []string
		Subnets       []struct {
			ID                   string
			Name                 string
			AddressPrefix        string
			NetworkSecurityGroup *string
			RouteTable           *string
			NATGateway           *string
			PrivateEndpoints     []string
			ServiceEndpoints     []string
			Delegations          []string
		}
		Peerings []struct {
			ID                    string
			Name                  string
			RemoteVNetID          string
			RemoteVNetName        string
			PeeringState          string
			AllowVNetAccess       bool
			AllowForwardedTraffic bool
			AllowGatewayTransit   bool
			UseRemoteGateways     bool
		}
		DNSServers []string
		EnableDDoS bool
	}
	NSGs []struct {
		ID            string
		Name          string
		ResourceGroup string
		Location      string
		SecurityRules []struct {
			Name                     string
			Priority                 int32
			Direction                string
			Access                   string
			Protocol                 string
			SourceAddressPrefix      string
			SourcePortRange          string
			DestinationAddressPrefix string
			DestinationPortRange     string
			Description              string
		}
		Associations struct {
			Subnets           []string
			NetworkInterfaces []string
		}
	}
	PrivateEndpoints []struct {
		ID                   string
		Name                 string
		ResourceGroup        string
		Location             string
		SubnetID             string
		PrivateIPAddress     string
		PrivateLinkServiceID string
		ConnectionState      string
		GroupIDs             []string
	}
}
