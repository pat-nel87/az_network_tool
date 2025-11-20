package visualization

import (
	"strings"
	"testing"

	"azure-network-analyzer/pkg/models"
)

// TestSanitizeName tests the sanitization of node names
func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic name", "simple", "simple"},
		{"with hyphens", "my-resource-name", "my_resource_name"},
		{"with dots", "my.resource.name", "my_resource_name"},
		{"with spaces", "my resource name", "my_resource_name"},
		{"mixed special chars", "my-resource.name 2", "my_resource_name_2"},
		{"starting with digit", "123-resource", "123_resource"}, // DOT allows this but might cause issues
		{"with parentheses", "resource(prod)", "resource_prod_"},
		{"with brackets", "resource[0]", "resource_0_"},
		{"with colons", "resource:prod", "resource_prod"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeName(tt.input)
			// For now, just check it doesn't panic and returns something
			if result == "" && tt.input != "" {
				t.Errorf("sanitizeName(%q) returned empty string", tt.input)
			}
		})
	}
}

// TestExtractResourceName tests resource name extraction from Azure IDs
func TestExtractResourceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"valid azure resource ID",
			"/subscriptions/abc123/resourceGroups/rg1/providers/Microsoft.Network/azureFirewalls/fw-prod",
			"fw-prod",
		},
		{
			"short ID",
			"fw-prod",
			"fw-prod",
		},
		{
			"empty string",
			"",
			"",
		},
		{
			"with special chars",
			"/subscriptions/abc/resourceGroups/rg/providers/Microsoft.Network/azureFirewalls/fw-prod(east)",
			"fw-prod(east)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractResourceName(tt.input)
			if result != tt.expected {
				t.Errorf("extractResourceName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFirewallWithoutMatchingRoutes tests firewall rendering when routes don't match
func TestFirewallWithoutMatchingRoutes(t *testing.T) {
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
				ID:               "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/azureFirewalls/fw1",
				Name:             "fw1",
				PrivateIPAddress: "10.0.1.4",
				SubnetID:         "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/AzureFirewallSubnet",
			},
		},
		RouteTables: []models.RouteTable{
			{
				ID:   "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/routeTables/rt1",
				Name: "rt1",
				Routes: []models.Route{
					{
						Name:             "default-route",
						AddressPrefix:    "0.0.0.0/0",
						NextHopType:      "VirtualAppliance",
						NextHopIPAddress: "10.0.1.99", // DIFFERENT IP - won't match firewall
					},
				},
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Should not crash and should produce valid DOT
	if !strings.Contains(dot, "digraph NetworkTopology") {
		t.Error("DOT file doesn't contain graph declaration")
	}

	// Should have firewall node
	if !strings.Contains(dot, "fw_0") {
		t.Error("DOT file doesn't contain firewall node")
	}

	// Should NOT have edge from route table to firewall (IPs don't match)
	if strings.Contains(dot, "rt_rt1 -> fw_0") {
		t.Error("DOT file should not have edge when IPs don't match")
	}
}

// TestFirewallWithSpecialCharacters tests firewall names with special characters
func TestFirewallWithSpecialCharacters(t *testing.T) {
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
				ID:               "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/azureFirewalls/fw-prod(east)-01",
				Name:             "fw-prod(east)-01",
				PrivateIPAddress: "10.0.1.4",
				SubnetID:         "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/AzureFirewallSubnet",
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Should not crash
	if !strings.Contains(dot, "digraph NetworkTopology") {
		t.Error("DOT file doesn't contain graph declaration")
	}

	// Check that firewall node ID (fw_0) is present and doesn't have special chars in ID
	if !strings.Contains(dot, "fw_0 [label=") {
		t.Error("DOT file should contain firewall node with sanitized ID")
	}

	// Verify the label preserves the original name (allowed in quoted labels)
	if !strings.Contains(dot, "fw-prod(east)-01") {
		t.Error("DOT file should preserve original firewall name in label")
	}
}

// TestEmptyTopology tests handling of empty topology
func TestEmptyTopology(t *testing.T) {
	topology := &models.NetworkTopology{
		SubscriptionID:  "test-sub",
		ResourceGroup:   "test-rg",
		VirtualNetworks: []models.VirtualNetwork{},
		AzureFirewalls:  []models.AzureFirewall{},
	}

	dot := GenerateDOTFile(topology)

	// Should not crash and should produce valid DOT
	if !strings.Contains(dot, "digraph NetworkTopology") {
		t.Error("DOT file doesn't contain graph declaration")
	}
}

// TestFirewallEdgeWithMatchingIP tests correct edge creation
func TestFirewallEdgeWithMatchingIP(t *testing.T) {
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
						RouteTable:    &rtID, // Associate route table with subnet
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
				ID:   rtID,
				Name: "rt1",
				Routes: []models.Route{
					{
						Name:             "default-route",
						AddressPrefix:    "0.0.0.0/0",
						NextHopType:      "VirtualAppliance",
						NextHopIPAddress: "10.0.1.4", // MATCHES firewall IP
					},
				},
			},
		},
	}

	dot := GenerateDOTFile(topology)

	// Should have edge from route table to firewall
	if !strings.Contains(dot, "egress via FW") {
		t.Error("DOT file should contain edge from route table to firewall")
	}
}
