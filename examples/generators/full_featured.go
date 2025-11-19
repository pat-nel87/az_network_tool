package main

import (
	"fmt"
	"os"

	"azure-network-analyzer/pkg/models"
	"azure-network-analyzer/pkg/visualization"
)

func main() {
	// Comprehensive topology showcasing all features

	sharedNATID := "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/natGateways/prod-nat-gateway"
	nsg1ID := "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/networkSecurityGroups/web-nsg"
	nsg2ID := "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/networkSecurityGroups/data-nsg"
	rt1ID := "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/routeTables/firewall-rt"

	topology := &models.NetworkTopology{
		SubscriptionID: "example-subscription-12345",
		ResourceGroup:  "full-featured-demo-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
				Name:         "hub-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/GatewaySubnet",
						Name:          "GatewaySubnet",
						AddressPrefix: "10.0.0.0/27",
					},
					{
						ID:                   "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/shared-services",
						Name:                 "shared-services",
						AddressPrefix:        "10.0.1.0/24",
						NetworkSecurityGroup: &nsg1ID,
						RouteTable:           &rt1ID,
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                  "hub-to-prod",
						RemoteVNetID:          "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
						RemoteVNetName:        "prod-vnet",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: true,
						AllowGatewayTransit:   true,
					},
				},
			},
			{
				ID:           "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
				Name:         "prod-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.1.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:                   "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/web-tier",
						Name:                 "web-tier",
						AddressPrefix:        "10.1.1.0/24",
						NetworkSecurityGroup: &nsg1ID,
						NATGateway:           &sharedNATID,
					},
					{
						ID:                   "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-tier",
						Name:                 "app-tier",
						AddressPrefix:        "10.1.2.0/24",
						NetworkSecurityGroup: &nsg1ID,
						NATGateway:           &sharedNATID,
					},
					{
						ID:                   "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/data-tier",
						Name:                 "data-tier",
						AddressPrefix:        "10.1.3.0/24",
						NetworkSecurityGroup: &nsg2ID,
					},
					{
						ID:            "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/private-links",
						Name:          "private-links-subnet",
						AddressPrefix: "10.1.4.0/24",
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                "prod-to-hub",
						RemoteVNetID:        "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
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
				ID:          "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworkGateways/hub-vpn-gw",
				Name:        "hub-vpn-gateway",
				Location:    "eastus",
				VNetID:      "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
				GatewayType: "Vpn",
				VpnType:     "RouteBased",
				SKU:         "VpnGw2",
			},
		},
		LoadBalancers: []models.LoadBalancer{
			{
				ID:       "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/loadBalancers/web-lb",
				Name:     "web-lb-public",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Public",
			},
			{
				ID:       "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/loadBalancers/app-lb",
				Name:     "app-lb-internal",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Internal",
			},
		},
		AppGateways: []models.ApplicationGateway{
			{
				ID:         "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/applicationGateways/prod-appgw",
				Name:       "prod-app-gateway-waf",
				Location:   "eastus",
				SKU:        "WAF_v2",
				Tier:       "WAF_v2",
				Capacity:   2,
				WAFEnabled: true,
				WAFMode:    "Prevention",
			},
		},
		NATGateways: []models.NATGateway{
			{
				ID:       sharedNATID,
				Name:     "prod-nat-gateway",
				Location: "eastus",
			},
		},
		NSGs: []models.NetworkSecurityGroup{
			{
				ID:       nsg1ID,
				Name:     "web-nsg",
				Location: "eastus",
				Associations: models.NSGAssociations{
					Subnets: []string{
						"/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/shared-services",
						"/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/web-tier",
						"/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-tier",
					},
				},
			},
			{
				ID:       nsg2ID,
				Name:     "data-nsg",
				Location: "eastus",
				Associations: models.NSGAssociations{
					Subnets: []string{
						"/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/data-tier",
					},
				},
			},
		},
		RouteTables: []models.RouteTable{
			{
				ID:       rt1ID,
				Name:     "firewall-rt",
				Location: "eastus",
				AssociatedSubnets: []string{
					"/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/shared-services",
				},
			},
		},
		PrivateEndpoints: []models.PrivateEndpoint{
			{
				ID:                   "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/privateEndpoints/storage-pe",
				Name:                 "storage-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/private-links",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Storage/storageAccounts/prodstorage",
				PrivateIPAddress:     "10.1.4.10",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/privateEndpoints/sql-pe",
				Name:                 "sql-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/private-links",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Sql/servers/prodsql",
				PrivateIPAddress:     "10.1.4.11",
				ConnectionState:      "Approved",
			},
			{
				ID:                   "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/privateEndpoints/keyvault-pe",
				Name:                 "keyvault-pe",
				Location:             "eastus",
				SubnetID:             "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/private-links",
				PrivateLinkServiceID: "/subscriptions/example/resourceGroups/full-demo-rg/providers/Microsoft.KeyVault/vaults/prodkv",
				PrivateIPAddress:     "10.1.4.12",
				ConnectionState:      "Approved",
			},
		},
	}

	fmt.Println("Generating full-featured topology...")
	fmt.Printf("  - 2 VNets (hub-and-spoke)\n")
	fmt.Printf("  - VNet peerings\n")
	fmt.Printf("  - VPN Gateway\n")
	fmt.Printf("  - 2 Load Balancers\n")
	fmt.Printf("  - 1 Application Gateway with WAF\n")
	fmt.Printf("  - NAT Gateway (shared by 2 subnets)\n")
	fmt.Printf("  - 2 NSGs\n")
	fmt.Printf("  - 1 Route Table\n")
	fmt.Printf("  - 3 Private Endpoints\n")
	fmt.Println()

	// Generate DOT
	dotContent := visualization.GenerateDOTFile(topology)

	// Save files
	dotFile := "docs/examples/full-featured.dot"
	os.WriteFile(dotFile, []byte(dotContent), 0644)
	fmt.Printf("✓ DOT file saved: %s\n", dotFile)

	svgContent, _ := visualization.RenderSVG(dotContent)
	svgFile := "docs/examples/full-featured.svg"
	os.WriteFile(svgFile, svgContent, 0644)
	fmt.Printf("✓ SVG file saved: %s (%.2f KB)\n", svgFile, float64(len(svgContent))/1024)
	fmt.Println("\nExample: Full-featured topology complete!")
}
