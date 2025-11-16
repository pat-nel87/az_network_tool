package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func TestSafeString(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty string",
			input:    strPtr(""),
			expected: "",
		},
		{
			name:     "non-empty string",
			input:    strPtr("test-value"),
			expected: "test-value",
		},
		{
			name:     "string with spaces",
			input:    strPtr("  test value  "),
			expected: "  test value  ",
		},
		{
			name:     "string with special characters",
			input:    strPtr("test/value-123_abc"),
			expected: "test/value-123_abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := safeString(tt.input)
			if result != tt.expected {
				t.Errorf("safeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractResourceName(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		expected   string
	}{
		{
			name:       "empty string",
			resourceID: "",
			expected:   "",
		},
		{
			name:       "simple resource ID",
			resourceID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
			expected:   "vnet1",
		},
		{
			name:       "subnet resource ID",
			resourceID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
			expected:   "subnet1",
		},
		{
			name:       "NSG resource ID",
			resourceID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
			expected:   "nsg1",
		},
		{
			name:       "resource name with hyphens",
			resourceID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/loadBalancers/lb-web-prod-001",
			expected:   "lb-web-prod-001",
		},
		{
			name:       "no slashes",
			resourceID: "simple-name",
			expected:   "simple-name",
		},
		{
			name:       "trailing slash",
			resourceID: "/subscriptions/sub1/",
			expected:   "",
		},
		{
			name:       "single slash",
			resourceID: "/",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractResourceName(tt.resourceID)
			if result != tt.expected {
				t.Errorf("extractResourceName(%q) = %q, want %q", tt.resourceID, result, tt.expected)
			}
		})
	}
}

func TestExtractVNetIDFromSubnet(t *testing.T) {
	tests := []struct {
		name     string
		subnetID string
		expected string
	}{
		{
			name:     "empty string",
			subnetID: "",
			expected: "",
		},
		{
			name:     "valid subnet ID",
			subnetID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
			expected: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
		},
		{
			name:     "subnet with complex name",
			subnetID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet-hub-prod/subnets/subnet-web-tier",
			expected: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet-hub-prod",
		},
		{
			name:     "GatewaySubnet",
			subnetID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/GatewaySubnet",
			expected: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
		},
		{
			name:     "no subnets segment",
			subnetID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
			expected: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
		},
		{
			name:     "malformed - subnets but no subnet name",
			subnetID: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/",
			expected: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVNetIDFromSubnet(tt.subnetID)
			if result != tt.expected {
				t.Errorf("extractVNetIDFromSubnet(%q) = %q, want %q", tt.subnetID, result, tt.expected)
			}
		})
	}
}

func TestExtractSubnet(t *testing.T) {
	client := &AzureClient{}

	t.Run("nil properties", func(t *testing.T) {
		subnet := &armnetwork.Subnet{
			ID:   strPtr("/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
			Name: strPtr("subnet1"),
		}

		result := client.extractSubnet(subnet)

		if result.ID != "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1" {
			t.Errorf("ID mismatch: got %s", result.ID)
		}
		if result.Name != "subnet1" {
			t.Errorf("Name mismatch: got %s", result.Name)
		}
		if result.AddressPrefix != "" {
			t.Errorf("AddressPrefix should be empty: got %s", result.AddressPrefix)
		}
		if result.NetworkSecurityGroup != nil {
			t.Error("NetworkSecurityGroup should be nil")
		}
	})

	t.Run("full subnet with all associations", func(t *testing.T) {
		nsgID := "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/nsg1"
		rtID := "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/routeTables/rt1"
		natID := "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/natGateways/nat1"
		addressPrefix := "10.0.1.0/24"

		subnet := &armnetwork.Subnet{
			ID:   strPtr("/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
			Name: strPtr("subnet1"),
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: &addressPrefix,
				NetworkSecurityGroup: &armnetwork.SecurityGroup{
					ID: &nsgID,
				},
				RouteTable: &armnetwork.RouteTable{
					ID: &rtID,
				},
				NatGateway: &armnetwork.SubResource{
					ID: &natID,
				},
				PrivateEndpoints: []*armnetwork.PrivateEndpoint{
					{ID: strPtr("pe1")},
					{ID: strPtr("pe2")},
				},
				ServiceEndpoints: []*armnetwork.ServiceEndpointPropertiesFormat{
					{Service: strPtr("Microsoft.Storage")},
					{Service: strPtr("Microsoft.KeyVault")},
				},
				Delegations: []*armnetwork.Delegation{
					{
						Properties: &armnetwork.ServiceDelegationPropertiesFormat{
							ServiceName: strPtr("Microsoft.Web/serverFarms"),
						},
					},
				},
			},
		}

		result := client.extractSubnet(subnet)

		if result.AddressPrefix != "10.0.1.0/24" {
			t.Errorf("AddressPrefix mismatch: got %s", result.AddressPrefix)
		}
		if result.NetworkSecurityGroup == nil || *result.NetworkSecurityGroup != nsgID {
			t.Errorf("NetworkSecurityGroup mismatch")
		}
		if result.RouteTable == nil || *result.RouteTable != rtID {
			t.Errorf("RouteTable mismatch")
		}
		if result.NATGateway == nil || *result.NATGateway != natID {
			t.Errorf("NATGateway mismatch")
		}
		if len(result.PrivateEndpoints) != 2 {
			t.Errorf("PrivateEndpoints count mismatch: got %d", len(result.PrivateEndpoints))
		}
		if len(result.ServiceEndpoints) != 2 {
			t.Errorf("ServiceEndpoints count mismatch: got %d", len(result.ServiceEndpoints))
		}
		if len(result.Delegations) != 1 {
			t.Errorf("Delegations count mismatch: got %d", len(result.Delegations))
		}
	})
}

func TestExtractVNetPeering(t *testing.T) {
	client := &AzureClient{}

	t.Run("full peering configuration", func(t *testing.T) {
		peeringState := armnetwork.VirtualNetworkPeeringStateConnected
		peering := &armnetwork.VirtualNetworkPeering{
			ID:   strPtr("/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/virtualNetworkPeerings/peer1"),
			Name: strPtr("peer1"),
			Properties: &armnetwork.VirtualNetworkPeeringPropertiesFormat{
				RemoteVirtualNetwork: &armnetwork.SubResource{
					ID: strPtr("/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet2"),
				},
				PeeringState:              &peeringState,
				AllowVirtualNetworkAccess: boolPtr(true),
				AllowForwardedTraffic:     boolPtr(true),
				AllowGatewayTransit:       boolPtr(false),
				UseRemoteGateways:         boolPtr(false),
			},
		}

		result := client.extractVNetPeering(peering)

		if result.Name != "peer1" {
			t.Errorf("Name mismatch: got %s", result.Name)
		}
		if result.RemoteVNetName != "vnet2" {
			t.Errorf("RemoteVNetName mismatch: got %s", result.RemoteVNetName)
		}
		if result.PeeringState != "Connected" {
			t.Errorf("PeeringState mismatch: got %s", result.PeeringState)
		}
		if !result.AllowVNetAccess {
			t.Error("AllowVNetAccess should be true")
		}
		if !result.AllowForwardedTraffic {
			t.Error("AllowForwardedTraffic should be true")
		}
		if result.AllowGatewayTransit {
			t.Error("AllowGatewayTransit should be false")
		}
	})
}

func TestExtractRoute(t *testing.T) {
	client := &AzureClient{}

	t.Run("route with virtual appliance", func(t *testing.T) {
		nextHopType := armnetwork.RouteNextHopTypeVirtualAppliance
		route := &armnetwork.Route{
			Name: strPtr("route-to-internet"),
			Properties: &armnetwork.RoutePropertiesFormat{
				AddressPrefix:    strPtr("0.0.0.0/0"),
				NextHopType:      &nextHopType,
				NextHopIPAddress: strPtr("10.0.0.100"),
			},
		}

		result := client.extractRoute(route)

		if result.Name != "route-to-internet" {
			t.Errorf("Name mismatch: got %s", result.Name)
		}
		if result.AddressPrefix != "0.0.0.0/0" {
			t.Errorf("AddressPrefix mismatch: got %s", result.AddressPrefix)
		}
		if result.NextHopType != "VirtualAppliance" {
			t.Errorf("NextHopType mismatch: got %s", result.NextHopType)
		}
		if result.NextHopIPAddress != "10.0.0.100" {
			t.Errorf("NextHopIPAddress mismatch: got %s", result.NextHopIPAddress)
		}
	})
}

func TestExtractSecurityRule(t *testing.T) {
	t.Run("allow rule with all properties", func(t *testing.T) {
		priority := int32(100)
		direction := armnetwork.SecurityRuleDirectionInbound
		access := armnetwork.SecurityRuleAccessAllow
		protocol := armnetwork.SecurityRuleProtocolTCP

		rule := &armnetwork.SecurityRule{
			Name: strPtr("AllowHTTP"),
			Properties: &armnetwork.SecurityRulePropertiesFormat{
				Priority:                 &priority,
				Direction:                &direction,
				Access:                   &access,
				Protocol:                 &protocol,
				SourceAddressPrefix:      strPtr("*"),
				SourcePortRange:          strPtr("*"),
				DestinationAddressPrefix: strPtr("*"),
				DestinationPortRange:     strPtr("80"),
				Description:              strPtr("Allow HTTP traffic"),
			},
		}

		result := extractSecurityRule(rule)

		if result.Name != "AllowHTTP" {
			t.Errorf("Name mismatch: got %s", result.Name)
		}
		if result.Priority != 100 {
			t.Errorf("Priority mismatch: got %d", result.Priority)
		}
		if result.Direction != "Inbound" {
			t.Errorf("Direction mismatch: got %s", result.Direction)
		}
		if result.Access != "Allow" {
			t.Errorf("Access mismatch: got %s", result.Access)
		}
		if result.Protocol != "Tcp" {
			t.Errorf("Protocol mismatch: got %s", result.Protocol)
		}
		if result.DestinationPortRange != "80" {
			t.Errorf("DestinationPortRange mismatch: got %s", result.DestinationPortRange)
		}
	})
}

// Helper functions for tests
func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}
