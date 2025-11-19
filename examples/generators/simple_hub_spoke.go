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
	// Simple hub-and-spoke topology
	// 1 hub VNet with gateway subnet
	// 2 spoke VNets (production and development)

	topology := &models.NetworkTopology{
		SubscriptionID: "example-subscription-12345",
		ResourceGroup:  "example-hub-spoke-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
				Name:         "hub-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/GatewaySubnet",
						Name:          "GatewaySubnet",
						AddressPrefix: "10.0.0.0/27",
					},
					{
						ID:            "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/shared-services",
						Name:          "shared-services",
						AddressPrefix: "10.0.1.0/24",
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                  "hub-to-prod",
						RemoteVNetID:          "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
						RemoteVNetName:        "prod-vnet",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: true,
						AllowGatewayTransit:   true,
					},
					{
						Name:                  "hub-to-dev",
						RemoteVNetID:          "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/dev-vnet",
						RemoteVNetName:        "dev-vnet",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: true,
						AllowGatewayTransit:   true,
					},
				},
			},
			{
				ID:           "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
				Name:         "prod-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.1.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/web-tier",
						Name:          "web-tier",
						AddressPrefix: "10.1.1.0/24",
					},
					{
						ID:            "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-tier",
						Name:          "app-tier",
						AddressPrefix: "10.1.2.0/24",
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                "prod-to-hub",
						RemoteVNetID:        "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
						RemoteVNetName:      "hub-vnet",
						PeeringState:        "Connected",
						AllowVNetAccess:     true,
						UseRemoteGateways:   true,
					},
				},
			},
			{
				ID:           "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/dev-vnet",
				Name:         "dev-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.2.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/dev-vnet/subnets/dev-subnet",
						Name:          "dev-subnet",
						AddressPrefix: "10.2.1.0/24",
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                "dev-to-hub",
						RemoteVNetID:        "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
						RemoteVNetName:      "hub-vnet",
						PeeringState:        "Connected",
						AllowVNetAccess:     true,
						UseRemoteGateways:   true,
					},
				},
			},
		},
		VPNGateways: []models.VPNGateway{
			{
				ID:          "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworkGateways/hub-vpn-gateway",
				Name:        "hub-vpn-gateway",
				Location:    "eastus",
				VNetID:      "/subscriptions/example/resourceGroups/hub-spoke-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
				GatewayType: "Vpn",
				VpnType:     "RouteBased",
				SKU:         "VpnGw1",
			},
		},
	}

	fmt.Println("Generating simple hub-spoke topology...")
	fmt.Printf("  - 1 Hub VNet (hub-vnet)\n")
	fmt.Printf("  - 2 Spoke VNets (prod-vnet, dev-vnet)\n")
	fmt.Printf("  - VNet peerings (bidirectional)\n")
	fmt.Printf("  - VPN Gateway in hub\n")
	fmt.Println()

	// Generate DOT
	dotContent := visualization.GenerateDOTFile(topology)

	// Save DOT file
	dotFile := "docs/examples/simple-hub-spoke.dot"
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
	svgFile := "docs/examples/simple-hub-spoke.svg"
	err = os.WriteFile(svgFile, svgContent, 0644)
	if err != nil {
		fmt.Printf("Error writing SVG file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ SVG file saved: %s (%.2f KB)\n", svgFile, float64(len(svgContent))/1024)
	fmt.Println("\nExample: Simple hub-and-spoke architecture complete!")
}
