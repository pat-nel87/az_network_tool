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
	// Topology demonstrating NAT gateway sharing across multiple subnets
	// This showcases the smart deduplication feature

	sharedNATID := "/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/natGateways/shared-nat-gateway"
	nsgID := "/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/networkSecurityGroups/web-nsg"

	topology := &models.NetworkTopology{
		SubscriptionID: "example-subscription-12345",
		ResourceGroup:  "nat-gateway-demo-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
				Name:         "prod-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:                   "/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/web-subnet",
						Name:                 "web-subnet",
						AddressPrefix:        "10.0.1.0/24",
						NetworkSecurityGroup: &nsgID,
						NATGateway:           &sharedNATID, // Shared NAT
					},
					{
						ID:                   "/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-subnet",
						Name:                 "app-subnet",
						AddressPrefix:        "10.0.2.0/24",
						NetworkSecurityGroup: &nsgID,
						NATGateway:           &sharedNATID, // Shared NAT
					},
					{
						ID:                   "/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/data-subnet",
						Name:                 "data-subnet",
						AddressPrefix:        "10.0.3.0/24",
						NATGateway:           &sharedNATID, // Shared NAT
					},
					{
						ID:            "/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/mgmt-subnet",
						Name:          "mgmt-subnet",
						AddressPrefix: "10.0.4.0/24",
						// No NAT gateway for management subnet
					},
				},
			},
		},
		NATGateways: []models.NATGateway{
			{
				ID:       sharedNATID,
				Name:     "shared-nat-gateway",
				Location: "eastus",
			},
		},
		NSGs: []models.NetworkSecurityGroup{
			{
				ID:       nsgID,
				Name:     "web-nsg",
				Location: "eastus",
				SecurityRules: []models.SecurityRule{
					{
						Name:                     "allow-http",
						Priority:                 100,
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "TCP",
						SourcePortRange:          "*",
						DestinationPortRange:     "80",
						SourceAddressPrefix:      "*",
						DestinationAddressPrefix: "*",
					},
				},
				Associations: models.NSGAssociations{
					Subnets: []string{
						"/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/web-subnet",
						"/subscriptions/example/resourceGroups/nat-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-subnet",
					},
				},
			},
		},
	}

	fmt.Println("Generating NAT Gateway sharing topology...")
	fmt.Printf("  - 1 VNet with 4 subnets\n")
	fmt.Printf("  - 3 subnets sharing 1 NAT Gateway (smart deduplication)\n")
	fmt.Printf("  - 1 mgmt subnet without NAT Gateway\n")
	fmt.Printf("  - 1 NSG shared across 2 subnets\n")
	fmt.Println()
	fmt.Println("Key Feature: NAT Gateway appears once with connections from 3 subnets")
	fmt.Println()

	// Generate DOT
	dotContent := visualization.GenerateDOTFile(topology)

	// Save DOT file
	dotFile := "docs/examples/nat-gateway-sharing.dot"
	err := os.WriteFile(dotFile, []byte(dotContent), 0644)
	if err != nil {
		fmt.Printf("Error writing DOT file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ DOT file saved: %s\n", dotFile)

	// Render SVG
	svgContent, err := visualization.RenderSVG(dotContent)
	if err != nil {
		fmt.Printf("Error rendering SVG: %v\n", err)
		os.Exit(1)
	}

	// Save SVG
	svgFile := "docs/examples/nat-gateway-sharing.svg"
	err = os.WriteFile(svgFile, svgContent, 0644)
	if err != nil {
		fmt.Printf("Error writing SVG file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ SVG file saved: %s (%.2f KB)\n", svgFile, float64(len(svgContent))/1024)
	fmt.Println("\nExample: NAT Gateway sharing (deduplication) complete!")
}
