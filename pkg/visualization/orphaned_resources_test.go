package visualization

import (
	"strings"
	"testing"

	"azure-network-analyzer/pkg/models"
)

// TestOrphanedRouteTable tests that route tables without subnet associations are rendered
func TestOrphanedRouteTable(t *testing.T) {
	topology := &models.NetworkTopology{
		SubscriptionID: "test-sub",
		ResourceGroup:  "test-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1",
				Name:         "vnet1",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
						Name:          "subnet1",
						AddressPrefix: "10.0.1.0/24",
						// NOTE: No RouteTable association
					},
				},
			},
		},
		RouteTables: []models.RouteTable{
			{
				ID:   "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/routeTables/rt-orphaned",
				Name: "rt-orphaned",
				Routes: []models.Route{
					{
						Name:          "route1",
						AddressPrefix: "0.0.0.0/0",
						NextHopType:   "Internet",
					},
				},
				// NOTE: AssociatedSubnets is empty - this is an orphaned route table
				AssociatedSubnets: []string{},
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Check if orphaned route table is rendered
	if !strings.Contains(dot, "rt_orphaned") && !strings.Contains(dot, "Route Table") {
		t.Error("Orphaned route table should be rendered in the visualization")
	}
}

// TestOrphanedRouteTableWithFirewallRoute tests orphaned RT with firewall egress
func TestOrphanedRouteTableWithFirewallRoute(t *testing.T) {
	topology := &models.NetworkTopology{
		SubscriptionID: "test-sub",
		ResourceGroup:  "test-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1",
				Name:         "vnet1",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/AzureFirewallSubnet",
						Name:          "AzureFirewallSubnet",
						AddressPrefix: "10.0.1.0/24",
					},
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
						Name:          "subnet1",
						AddressPrefix: "10.0.2.0/24",
						// NOTE: No RouteTable association
					},
				},
			},
		},
		AzureFirewalls: []models.AzureFirewall{
			{
				ID:               "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/azureFirewalls/fw1",
				Name:             "fw1",
				PrivateIPAddress: "10.0.1.4",
				SubnetID:         "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/AzureFirewallSubnet",
			},
		},
		RouteTables: []models.RouteTable{
			{
				ID:   "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/routeTables/rt-orphaned-with-fw",
				Name: "rt-orphaned-with-fw",
				Routes: []models.Route{
					{
						Name:             "default-via-firewall",
						AddressPrefix:    "0.0.0.0/0",
						NextHopType:      "VirtualAppliance",
						NextHopIPAddress: "10.0.1.4", // Points to firewall
					},
				},
				AssociatedSubnets: []string{}, // Orphaned
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Should render the orphaned route table
	if !strings.Contains(dot, "rt_orphaned") {
		t.Error("Orphaned route table should be rendered even without subnet associations")
	}

	// Should show the route table -> firewall edge
	if !strings.Contains(dot, "egress via FW") {
		t.Error("Route table pointing to firewall should show egress edge even if orphaned")
	}
}

// TestOrphanedNSG tests if orphaned NSGs are rendered
func TestOrphanedNSG(t *testing.T) {
	topology := &models.NetworkTopology{
		SubscriptionID: "test-sub",
		ResourceGroup:  "test-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1",
				Name:         "vnet1",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
						Name:          "subnet1",
						AddressPrefix: "10.0.1.0/24",
						// NOTE: No NSG association
					},
				},
			},
		},
		NSGs: []models.NetworkSecurityGroup{
			{
				ID:   "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/networkSecurityGroups/nsg-orphaned",
				Name: "nsg-orphaned",
				Associations: models.NSGAssociations{
					Subnets:           []string{}, // No associations
					NetworkInterfaces: []string{},
				},
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Check if orphaned NSG is rendered
	// NSGs might be rendered differently - let's see
	hasNSG := strings.Contains(dot, "nsg_orphaned") || strings.Contains(dot, "NSG")

	t.Logf("Orphaned NSG rendered: %v", hasNSG)
	if !hasNSG {
		t.Log("Note: Orphaned NSG is not rendered - same issue as route tables")
	}
}
