package visualization

import (
	"strings"
	"testing"

	"azure-network-analyzer/pkg/models"
)

// TestFirewallWithPublicIPEgress tests that firewalls with public IPs show Internet egress
func TestFirewallWithPublicIPEgress(t *testing.T) {
	rtID := "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/routeTables/rt1"
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
						RouteTable:    &rtID,
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
				PublicIPAddresses: []string{
					"/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/fw-pip",
				},
			},
		},
		RouteTables: []models.RouteTable{
			{
				ID:   rtID,
				Name: "rt1",
				Routes: []models.Route{
					{
						Name:             "default-route",
						AddressPrefix:    "0.0.0.0/0",
						NextHopType:      "VirtualAppliance",
						NextHopIPAddress: "10.0.1.4",
					},
				},
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Should have Internet node
	if !strings.Contains(dot, "internet [") {
		t.Error("DOT file should contain Internet node for firewall egress")
	}

	// Should have edge from firewall to Internet
	if !strings.Contains(dot, "fw_0 -> internet") {
		t.Error("DOT file should show edge from firewall to Internet")
	}

	// Should indicate public IP egress
	if !strings.Contains(dot, "Public IP egress") {
		t.Error("DOT file should indicate firewall uses public IP for Internet egress")
	}
}

// TestFirewallWithoutPublicIP tests that firewalls without public IPs don't show Internet egress
func TestFirewallWithoutPublicIP(t *testing.T) {
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
				},
			},
		},
		AzureFirewalls: []models.AzureFirewall{
			{
				ID:                "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/azureFirewalls/fw1",
				Name:              "fw1",
				PrivateIPAddress:  "10.0.1.4",
				SubnetID:          "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/AzureFirewallSubnet",
				PublicIPAddresses: []string{}, // NO PUBLIC IP
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Should NOT have edge from firewall to Internet (no public IP)
	if strings.Contains(dot, "fw_0 -> internet") {
		t.Error("DOT file should not show Internet egress for firewall without public IP")
	}
}

// TestMultipleFirewallsWithPublicIPs tests multiple firewalls with Internet egress
func TestMultipleFirewallsWithPublicIPs(t *testing.T) {
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
				},
			},
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet2",
				Name:         "vnet2",
				AddressSpace: []string{"10.1.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet2/subnets/AzureFirewallSubnet",
						Name:          "AzureFirewallSubnet",
						AddressPrefix: "10.1.1.0/24",
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
				PublicIPAddresses: []string{
					"/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/fw1-pip",
				},
			},
			{
				ID:               "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/azureFirewalls/fw2",
				Name:             "fw2",
				PrivateIPAddress: "10.1.1.4",
				SubnetID:         "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet2/subnets/AzureFirewallSubnet",
				PublicIPAddresses: []string{
					"/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/fw2-pip",
				},
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Should have single Internet node (shared)
	internetCount := strings.Count(dot, "internet [label=\"Internet\"")
	if internetCount != 1 {
		t.Errorf("Should have exactly 1 Internet node, got %d", internetCount)
	}

	// Should have edges from both firewalls to Internet
	if !strings.Contains(dot, "fw_0 -> internet") {
		t.Error("DOT file should show edge from fw_0 to Internet")
	}
	if !strings.Contains(dot, "fw_1 -> internet") {
		t.Error("DOT file should show edge from fw_1 to Internet")
	}
}
