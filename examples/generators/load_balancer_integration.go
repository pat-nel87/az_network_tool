package main

import (
	"fmt"
	"os"

	"azure-network-analyzer/pkg/models"
	"azure-network-analyzer/pkg/visualization"
)

func main() {
	// Topology showcasing load balancers and application gateways

	topology := &models.NetworkTopology{
		SubscriptionID: "example-subscription-12345",
		ResourceGroup:  "load-balancer-demo-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/example/resourceGroups/lb-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet",
				Name:         "app-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/example/resourceGroups/lb-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/web-tier",
						Name:          "web-tier",
						AddressPrefix: "10.0.1.0/24",
					},
					{
						ID:            "/subscriptions/example/resourceGroups/lb-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/app-tier",
						Name:          "app-tier",
						AddressPrefix: "10.0.2.0/24",
					},
					{
						ID:            "/subscriptions/example/resourceGroups/lb-demo-rg/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/appgw-subnet",
						Name:          "appgw-subnet",
						AddressPrefix: "10.0.3.0/24",
					},
				},
			},
		},
		LoadBalancers: []models.LoadBalancer{
			{
				ID:       "/subscriptions/example/resourceGroups/lb-demo-rg/providers/Microsoft.Network/loadBalancers/public-lb",
				Name:     "public-web-lb",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Public",
			},
			{
				ID:       "/subscriptions/example/resourceGroups/lb-demo-rg/providers/Microsoft.Network/loadBalancers/internal-lb",
				Name:     "internal-app-lb",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Internal",
			},
		},
		AppGateways: []models.ApplicationGateway{
			{
				ID:         "/subscriptions/example/resourceGroups/lb-demo-rg/providers/Microsoft.Network/applicationGateways/app-gateway",
				Name:       "prod-app-gateway",
				Location:   "eastus",
				SKU:        "WAF_v2",
				Tier:       "WAF_v2",
				Capacity:   2,
				WAFEnabled: true,
				WAFMode:    "Prevention",
			},
		},
	}

	fmt.Println("Generating load balancer integration topology...")
	fmt.Printf("  - 1 VNet with 3 tiers\n")
	fmt.Printf("  - 2 Load Balancers (1 public, 1 internal)\n")
	fmt.Printf("  - 1 Application Gateway with WAF\n")
	fmt.Println()

	// Generate DOT
	dotContent := visualization.GenerateDOTFile(topology)

	// Save files
	dotFile := "docs/examples/load-balancer-integration.dot"
	os.WriteFile(dotFile, []byte(dotContent), 0644)
	fmt.Printf("✓ DOT file saved: %s\n", dotFile)

	svgContent, _ := visualization.RenderSVG(dotContent)
	svgFile := "docs/examples/load-balancer-integration.svg"
	os.WriteFile(svgFile, svgContent, 0644)
	fmt.Printf("✓ SVG file saved: %s (%.2f KB)\n", svgFile, float64(len(svgContent))/1024)
	fmt.Println("\nExample: Load balancer integration complete!")
}
