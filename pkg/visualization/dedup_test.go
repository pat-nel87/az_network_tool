package visualization

import (
	"strings"
	"testing"

	"azure-network-analyzer/pkg/models"
)

func TestNATGatewayDeduplication(t *testing.T) {
	// Create a topology with 3 subnets sharing the same NAT Gateway
	sharedNATID := "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/natGateways/shared-nat"

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
						NATGateway:    &sharedNATID,
					},
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2",
						Name:          "subnet2",
						AddressPrefix: "10.0.2.0/24",
						NATGateway:    &sharedNATID,
					},
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet3",
						Name:          "subnet3",
						AddressPrefix: "10.0.3.0/24",
						NATGateway:    &sharedNATID,
					},
				},
			},
		},
	}

	// Generate DOT file
	dot := GenerateDOTFile(topology)

	// Count occurrences of NAT Gateway node declarations
	// Should only appear once: "nat_shared_nat [label="NAT Gateway\nshared-nat"..."
	natNodeCount := strings.Count(dot, `[label="NAT Gateway\nshared-nat"`)

	if natNodeCount != 1 {
		t.Errorf("Expected 1 NAT Gateway node, got %d", natNodeCount)
		t.Logf("DOT content:\n%s", dot)
	}

	// Verify all 3 subnets have edges to the NAT Gateway
	// Should have 3 edges: "subnet_0_X -> nat_shared_nat [style=solid, color=green, label="egress"]"
	edgeCount := strings.Count(dot, `-> nat_shared_nat [style=solid, color=green, label="egress"]`)

	if edgeCount != 3 {
		t.Errorf("Expected 3 edges to NAT Gateway, got %d", edgeCount)
		t.Logf("DOT content:\n%s", dot)
	}

	t.Logf("✓ NAT Gateway deduplication working correctly")
	t.Logf("  - Single NAT node created")
	t.Logf("  - All 3 subnets connected to shared NAT")
}

func TestMixedNATGateways(t *testing.T) {
	// Create a topology with some subnets sharing NAT, others with dedicated NAT
	sharedNATID := "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/natGateways/shared-nat"
	dedicatedNATID := "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/natGateways/dedicated-nat"

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
						NATGateway:    &sharedNATID,
					},
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2",
						Name:          "subnet2",
						AddressPrefix: "10.0.2.0/24",
						NATGateway:    &sharedNATID,
					},
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet3",
						Name:          "subnet3",
						AddressPrefix: "10.0.3.0/24",
						NATGateway:    &dedicatedNATID,
					},
				},
			},
		},
	}

	// Generate DOT file
	dot := GenerateDOTFile(topology)

	// Should have exactly 2 NAT Gateway nodes
	sharedNATCount := strings.Count(dot, `[label="NAT Gateway\nshared-nat"`)
	dedicatedNATCount := strings.Count(dot, `[label="NAT Gateway\ndedicated-nat"`)

	if sharedNATCount != 1 {
		t.Errorf("Expected 1 shared-nat node, got %d", sharedNATCount)
	}
	if dedicatedNATCount != 1 {
		t.Errorf("Expected 1 dedicated-nat node, got %d", dedicatedNATCount)
	}

	// Verify edge counts
	sharedEdges := strings.Count(dot, `-> nat_shared_nat [style=solid, color=green, label="egress"]`)
	dedicatedEdges := strings.Count(dot, `-> nat_dedicated_nat [style=solid, color=green, label="egress"]`)

	if sharedEdges != 2 {
		t.Errorf("Expected 2 edges to shared NAT, got %d", sharedEdges)
	}
	if dedicatedEdges != 1 {
		t.Errorf("Expected 1 edge to dedicated NAT, got %d", dedicatedEdges)
	}

	t.Logf("✓ Mixed NAT Gateway topology working correctly")
	t.Logf("  - 2 NAT nodes created (shared + dedicated)")
	t.Logf("  - Correct edge connections for each")
}
