//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"azure-network-analyzer/pkg/models"
	"azure-network-analyzer/pkg/visualization"
)

func main() {
	// Topology showcasing private endpoints in table format

	topology := &models.NetworkTopology{
		SubscriptionID: "example-subscription-12345",
		ResourceGroup:  "private-links-demo-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet",
				Name:         "services-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/private-endpoints",
						Name:          "private-endpoints-subnet",
						AddressPrefix: "10.0.1.0/24",
					},
					{
						ID:            "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/app-subnet",
						Name:          "app-subnet",
						AddressPrefix: "10.0.2.0/24",
					},
				},
			},
		},
		PrivateEndpoints: []models.PrivateEndpoint{
			{
				ID:                   "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/privateEndpoints/storage-pe",
				Name:                 "storage-account-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/private-endpoints",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Storage/storageAccounts/prodstorageacct",
				PrivateIPAddress:     "10.0.1.10",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/privateEndpoints/sql-pe",
				Name:                 "sql-database-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/private-endpoints",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Sql/servers/prodsqlserver",
				PrivateIPAddress:     "10.0.1.11",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/privateEndpoints/keyvault-pe",
				Name:                 "keyvault-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/private-endpoints",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.KeyVault/vaults/prodkeyvault",
				PrivateIPAddress:     "10.0.1.12",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/privateEndpoints/cosmos-pe",
				Name:                 "cosmosdb-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/private-endpoints",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.DocumentDB/databaseAccounts/prodcosmosdb",
				PrivateIPAddress:     "10.0.1.13",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/privateEndpoints/acr-pe",
				Name:                 "container-registry-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/private-endpoints",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/pl-demo-rg/providers/Microsoft.ContainerRegistry/registries/prodacr",
				PrivateIPAddress:     "10.0.1.14",
				ConnectionState:      "Approved",
			},
		},
	}

	fmt.Println("Generating private links topology...")
	fmt.Printf("  - 1 VNet with dedicated private endpoints subnet\n")
	fmt.Printf("  - 5 Private Endpoints (Storage, SQL, KeyVault, CosmosDB, ACR)\n")
	fmt.Printf("  - Private links shown in clean table format\n")
	fmt.Println()

	// Generate DOT
	dotContent := visualization.GenerateDOTFile(topology)

	// Save files
	dotFile := "docs/examples/private-links-example.dot"
	os.WriteFile(dotFile, []byte(dotContent), 0644)
	fmt.Printf("✓ DOT file saved: %s\n", dotFile)

	svgContent, _ := visualization.RenderSVG(dotContent)
	svgFile := "docs/examples/private-links-example.svg"
	os.WriteFile(svgFile, svgContent, 0644)
	fmt.Printf("✓ SVG file saved: %s (%.2f KB)\n", svgFile, float64(len(svgContent))/1024)
	fmt.Println("\nExample: Private links table complete!")
}
