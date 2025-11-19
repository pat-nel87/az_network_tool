package main

import (
	"fmt"
	"os"

	"azure-network-analyzer/pkg/models"
	"azure-network-analyzer/pkg/visualization"
)

func main() {
	// Topology showcasing subnet-level details with NSGs and Route Tables

	nsg1ID := "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/networkSecurityGroups/web-nsg"
	nsg2ID := "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/networkSecurityGroups/db-nsg"
	rt1ID := "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/routeTables/firewall-rt"

	topology := &models.NetworkTopology{
		SubscriptionID: "example-subscription-12345",
		ResourceGroup:  "subnet-details-demo-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet",
				Name:         "app-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:                   "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/frontend",
						Name:                 "frontend",
						AddressPrefix:        "10.0.1.0/24",
						NetworkSecurityGroup: &nsg1ID,
						RouteTable:           &rt1ID,
					},
					{
						ID:                   "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/backend",
						Name:                 "backend",
						AddressPrefix:        "10.0.2.0/24",
						NetworkSecurityGroup: &nsg1ID,
						RouteTable:           &rt1ID,
					},
					{
						ID:                   "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/database",
						Name:                 "database",
						AddressPrefix:        "10.0.3.0/24",
						NetworkSecurityGroup: &nsg2ID,
					},
					{
						ID:            "/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/AzureBastionSubnet",
						Name:          "AzureBastionSubnet",
						AddressPrefix: "10.0.255.0/27",
						// No NSG or route table for Bastion
					},
				},
			},
		},
		NSGs: []models.NetworkSecurityGroup{
			{
				ID:       nsg1ID,
				Name:     "web-nsg",
				Location: "eastus",
				SecurityRules: []models.SecurityRule{
					{
						Name:                     "allow-https",
						Priority:                 100,
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "TCP",
						DestinationPortRange:     "443",
					},
					{
						Name:                     "allow-http",
						Priority:                 110,
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "TCP",
						DestinationPortRange:     "80",
					},
				},
				Associations: models.NSGAssociations{
					Subnets: []string{
						"/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/frontend",
						"/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/backend",
					},
				},
			},
			{
				ID:       nsg2ID,
				Name:     "db-nsg",
				Location: "eastus",
				SecurityRules: []models.SecurityRule{
					{
						Name:                     "allow-sql",
						Priority:                 100,
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "TCP",
						DestinationPortRange:     "1433",
						SourceAddressPrefix:      "10.0.2.0/24",
					},
				},
				Associations: models.NSGAssociations{
					Subnets: []string{
						"/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/database",
					},
				},
			},
		},
		RouteTables: []models.RouteTable{
			{
				ID:       rt1ID,
				Name:     "firewall-rt",
				Location: "eastus",
				Routes: []models.Route{
					{
						Name:             "to-firewall",
						AddressPrefix:    "0.0.0.0/0",
						NextHopType:      "VirtualAppliance",
						NextHopIPAddress: "10.0.100.4",
					},
					{
						Name:          "local-vnet",
						AddressPrefix: "10.0.0.0/16",
						NextHopType:   "VnetLocal",
					},
				},
				AssociatedSubnets: []string{
					"/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/frontend",
					"/subscriptions/example/resourceGroups/subnet-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/backend",
				},
			},
		},
	}

	fmt.Println("Generating subnet-level details topology...")
	fmt.Printf("  - 1 VNet with 4 subnets\n")
	fmt.Printf("  - 2 NSGs with security rules\n")
	fmt.Printf("  - 1 Route Table with custom routes\n")
	fmt.Printf("  - Subnet associations clearly shown\n")
	fmt.Println()

	// Generate DOT
	dotContent := visualization.GenerateDOTFile(topology)

	// Save files
	dotFile := "docs/examples/subnet-details.dot"
	os.WriteFile(dotFile, []byte(dotContent), 0644)
	fmt.Printf("✓ DOT file saved: %s\n", dotFile)

	svgContent, _ := visualization.RenderSVG(dotContent)
	svgFile := "docs/examples/subnet-details.svg"
	os.WriteFile(svgFile, svgContent, 0644)
	fmt.Printf("✓ SVG file saved: %s (%.2f KB)\n", svgFile, float64(len(svgContent))/1024)
	fmt.Println("\nExample: Subnet-level details complete!")
}
