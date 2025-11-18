package visualization

import (
	"strings"
	"testing"

	"azure-network-analyzer/pkg/models"
)

func TestPrivateEndpointsTable(t *testing.T) {
	// Create a topology with multiple private endpoints
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
					},
				},
			},
		},
		PrivateEndpoints: []models.PrivateEndpoint{
			{
				ID:                   "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/privateEndpoints/pe-storage",
				Name:                 "pe-storage",
				Location:             "eastus",
				SubnetID:             "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
				PrivateIPAddress:     "10.0.1.10",
				PrivateLinkServiceID: "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/myStorage",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/privateEndpoints/pe-sql",
				Name:                 "pe-sql",
				Location:             "eastus",
				SubnetID:             "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
				PrivateIPAddress:     "10.0.1.11",
				PrivateLinkServiceID: "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Sql/servers/mySqlServer",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/privateEndpoints/pe-keyvault",
				Name:                 "pe-keyvault",
				Location:             "eastus",
				SubnetID:             "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
				PrivateIPAddress:     "10.0.1.12",
				PrivateLinkServiceID: "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.KeyVault/vaults/myKeyVault",
				ConnectionState:      "Pending",
			},
		},
	}

	// Generate DOT file
	dot := GenerateDOTFile(topology)

	// Verify private endpoints are NOT rendered as nodes
	// Old behavior: pe_0, pe_1, pe_2 nodes with edges
	nodeCount := strings.Count(dot, `pe_0 [label=`)
	if nodeCount != 0 {
		t.Errorf("Private endpoints should not be rendered as nodes, found %d node declarations", nodeCount)
	}

	// Verify private endpoints table exists (as a node, not cluster)
	if !strings.Contains(dot, "pe_table") {
		t.Error("Private endpoints table node not found")
	}

	// Verify rank=sink is used to position table at bottom
	if !strings.Contains(dot, "rank=sink") {
		t.Error("Private endpoints table should use rank=sink for bottom positioning")
	}

	// Verify table header
	if !strings.Contains(dot, "Private Endpoints") {
		t.Error("Private endpoints table header not found")
	}

	// Verify table columns
	requiredColumns := []string{"Name", "Target Service", "Subnet", "Private IP", "Status"}
	for _, col := range requiredColumns {
		if !strings.Contains(dot, col) {
			t.Errorf("Table column '%s' not found", col)
		}
	}

	// Verify all 3 private endpoints are in the table
	if !strings.Contains(dot, "pe-storage") {
		t.Error("pe-storage not found in table")
	}
	if !strings.Contains(dot, "pe-sql") {
		t.Error("pe-sql not found in table")
	}
	if !strings.Contains(dot, "pe-keyvault") {
		t.Error("pe-keyvault not found in table")
	}

	// Verify target services are shown
	if !strings.Contains(dot, "myStorage") {
		t.Error("Storage account target not found in table")
	}
	if !strings.Contains(dot, "mySqlServer") {
		t.Error("SQL server target not found in table")
	}
	if !strings.Contains(dot, "myKeyVault") {
		t.Error("KeyVault target not found in table")
	}

	// Verify connection states
	if !strings.Contains(dot, "Approved") {
		t.Error("Approved status not found")
	}
	if !strings.Contains(dot, "Pending") {
		t.Error("Pending status not found")
	}

	t.Logf("✓ Private endpoints table rendering correctly")
	t.Logf("  - No private endpoint nodes in graph")
	t.Logf("  - Table with all 3 endpoints present")
	t.Logf("  - All columns and data displayed")
}

func TestNoPrivateEndpoints(t *testing.T) {
	// Topology with no private endpoints
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
					},
				},
			},
		},
		PrivateEndpoints: []models.PrivateEndpoint{},
	}

	// Generate DOT file
	dot := GenerateDOTFile(topology)

	// Verify no private endpoints table is created
	if strings.Contains(dot, "pe_table") {
		t.Error("Private endpoints table should not be present when there are no private endpoints")
	}

	t.Logf("✓ No private endpoints table when none exist")
}

func TestVNetLabels(t *testing.T) {
	// Test that VNet labels show actual names
	topology := &models.NetworkTopology{
		SubscriptionID: "test-sub",
		ResourceGroup:  "test-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet-eastus",
				Name:         "prod-vnet-eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets:      []models.Subnet{},
			},
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
				Name:         "hub-vnet",
				AddressSpace: []string{"10.1.0.0/16"},
				Subnets:      []models.Subnet{},
			},
		},
	}

	// Generate DOT file
	dot := GenerateDOTFile(topology)

	// Verify VNet names appear in labels
	if !strings.Contains(dot, "VNet: prod-vnet-eastus") {
		t.Error("VNet label 'prod-vnet-eastus' not found")
	}
	if !strings.Contains(dot, "VNet: hub-vnet") {
		t.Error("VNet label 'hub-vnet' not found")
	}

	// Verify bold font is applied
	if !strings.Contains(dot, "fontname=\"Helvetica-Bold\"") {
		t.Error("VNet labels should use bold font")
	}

	t.Logf("✓ VNet labels show actual names with proper formatting")
}
