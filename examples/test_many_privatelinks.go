package main

import (
	"fmt"
	"os"

	"azure-network-analyzer/pkg/models"
	"azure-network-analyzer/pkg/visualization"
)

func main() {
	// Shared resources
	sharedNATID := "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/natGateways/prod-nat-gateway"
	nsg1ID := "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/networkSecurityGroups/prod-nsg-web"
	nsg2ID := "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/networkSecurityGroups/prod-nsg-data"
	rt1ID := "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/routeTables/prod-rt-firewall"

	// Create a realistic production topology with 25 private endpoints and other resources
	topology := &models.NetworkTopology{
		SubscriptionID: "demo-subscription-abc123",
		ResourceGroup:  "production-enterprise-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
				Name:         "hub-vnet-eastus",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/GatewaySubnet",
						Name:                 "GatewaySubnet",
						AddressPrefix:        "10.0.1.0/27",
						NetworkSecurityGroup: &nsg1ID,
					},
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/AzureFirewallSubnet",
						Name:                 "AzureFirewallSubnet",
						AddressPrefix:        "10.0.2.0/24",
						RouteTable:           &rt1ID,
					},
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/shared-services",
						Name:                 "shared-services-subnet",
						AddressPrefix:        "10.0.3.0/24",
						NetworkSecurityGroup: &nsg1ID,
						NATGateway:           &sharedNATID,
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                  "hub-to-prod-spoke",
						RemoteVNetID:          "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
						RemoteVNetName:        "prod-vnet-eastus",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: true,
					},
					{
						Name:                  "hub-to-services-spoke",
						RemoteVNetID:          "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/services-vnet",
						RemoteVNetName:        "services-vnet-westus",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: true,
					},
				},
			},
			{
				ID:           "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
				Name:         "prod-vnet-eastus",
				Location:     "eastus",
				AddressSpace: []string{"10.1.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/web-tier",
						Name:                 "web-tier-subnet",
						AddressPrefix:        "10.1.1.0/24",
						NetworkSecurityGroup: &nsg1ID,
						NATGateway:           &sharedNATID,
					},
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-tier",
						Name:                 "app-tier-subnet",
						AddressPrefix:        "10.1.2.0/24",
						NetworkSecurityGroup: &nsg1ID,
						NATGateway:           &sharedNATID,
					},
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/data-tier",
						Name:                 "data-tier-subnet",
						AddressPrefix:        "10.1.3.0/24",
						NetworkSecurityGroup: &nsg2ID,
						RouteTable:           &rt1ID,
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                  "prod-to-hub",
						RemoteVNetID:          "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
						RemoteVNetName:        "hub-vnet-eastus",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: true,
					},
				},
			},
			{
				ID:           "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/services-vnet",
				Name:         "services-vnet-westus",
				Location:     "westus",
				AddressSpace: []string{"10.2.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/private-links",
						Name:                 "private-links-subnet",
						AddressPrefix:        "10.2.1.0/24",
						NetworkSecurityGroup: &nsg2ID,
					},
					{
						ID:                   "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/integration",
						Name:                 "integration-subnet",
						AddressPrefix:        "10.2.2.0/24",
						NetworkSecurityGroup: &nsg1ID,
						RouteTable:           &rt1ID,
					},
				},
				Peerings: []models.VNetPeering{
					{
						Name:                  "services-to-hub",
						RemoteVNetID:          "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
						RemoteVNetName:        "hub-vnet-eastus",
						PeeringState:          "Connected",
						AllowVNetAccess:       true,
						AllowForwardedTraffic: true,
					},
				},
			},
		},
		VPNGateways: []models.VPNGateway{
			{
				ID:          "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworkGateways/hub-vpn-gw",
				Name:        "hub-vpn-gateway",
				Location:    "eastus",
				VNetID:      "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet",
				GatewayType: "Vpn",
				VpnType:     "RouteBased",
				SKU:         "VpnGw2",
			},
		},
		LoadBalancers: []models.LoadBalancer{
			{
				ID:       "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/loadBalancers/prod-web-lb",
				Name:     "prod-web-lb-public",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Public",
			},
			{
				ID:       "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/loadBalancers/prod-app-lb",
				Name:     "prod-app-lb-internal",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Internal",
			},
			{
				ID:       "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/loadBalancers/prod-db-lb",
				Name:     "prod-db-lb-internal",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Internal",
			},
			{
				ID:       "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/loadBalancers/prod-api-lb",
				Name:     "prod-api-lb-public",
				Location: "eastus",
				SKU:      "Standard",
				Type:     "Public",
			},
			{
				ID:       "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/loadBalancers/prod-services-lb",
				Name:     "prod-services-lb-internal",
				Location: "westus",
				SKU:      "Standard",
				Type:     "Internal",
			},
		},
		AppGateways: []models.ApplicationGateway{
			{
				ID:         "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/applicationGateways/prod-appgw-waf-1",
				Name:       "prod-appgw-waf-primary",
				Location:   "eastus",
				SKU:        "WAF_v2",
				Tier:       "WAF_v2",
				Capacity:   3,
				WAFEnabled: true,
				WAFMode:    "Prevention",
			},
			{
				ID:         "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/applicationGateways/prod-appgw-waf-2",
				Name:       "prod-appgw-waf-secondary",
				Location:   "eastus",
				SKU:        "WAF_v2",
				Tier:       "WAF_v2",
				Capacity:   2,
				WAFEnabled: true,
				WAFMode:    "Detection",
			},
			{
				ID:         "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/applicationGateways/prod-appgw-waf-3",
				Name:       "prod-appgw-waf-regional",
				Location:   "westus",
				SKU:        "WAF_v2",
				Tier:       "WAF_v2",
				Capacity:   2,
				WAFEnabled: true,
				WAFMode:    "Prevention",
			},
			{
				ID:         "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/applicationGateways/services-appgw-waf",
				Name:       "services-appgw-waf",
				Location:   "westus",
				SKU:        "WAF_v2",
				Tier:       "WAF_v2",
				Capacity:   2,
				WAFEnabled: true,
				WAFMode:    "Prevention",
			},
			{
				ID:         "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/applicationGateways/staging-appgw",
				Name:       "staging-appgw-standard",
				Location:   "eastus",
				SKU:        "Standard_v2",
				Tier:       "Standard_v2",
				Capacity:   1,
				WAFEnabled: false,
			},
			{
				ID:         "/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/applicationGateways/dev-appgw-waf",
				Name:       "dev-appgw-waf",
				Location:   "eastus",
				SKU:        "WAF_v2",
				Tier:       "WAF_v2",
				Capacity:   1,
				WAFEnabled: true,
				WAFMode:    "Detection",
			},
		},
		PrivateEndpoints: []models.PrivateEndpoint{},
	}

	// Azure service types for variety
	services := []struct {
		prefix string
		service string
	}{
		{"storage", "Microsoft.Storage/storageAccounts"},
		{"sql", "Microsoft.Sql/servers"},
		{"cosmos", "Microsoft.DocumentDB/databaseAccounts"},
		{"keyvault", "Microsoft.KeyVault/vaults"},
		{"servicebus", "Microsoft.ServiceBus/namespaces"},
		{"eventhub", "Microsoft.EventHub/namespaces"},
		{"redis", "Microsoft.Cache/Redis"},
		{"acr", "Microsoft.ContainerRegistry/registries"},
		{"mysql", "Microsoft.DBforMySQL/servers"},
		{"postgres", "Microsoft.DBforPostgreSQL/servers"},
		{"appconfig", "Microsoft.AppConfiguration/configurationStores"},
		{"synapse", "Microsoft.Synapse/workspaces"},
		{"cognitiveservices", "Microsoft.CognitiveServices/accounts"},
		{"searchservice", "Microsoft.Search/searchServices"},
		{"webpubsub", "Microsoft.SignalRService/webPubSub"},
		{"apim", "Microsoft.ApiManagement/service"},
		{"datafactory", "Microsoft.DataFactory/factories"},
		{"batch", "Microsoft.Batch/batchAccounts"},
		{"iot", "Microsoft.Devices/IotHubs"},
		{"monitor", "Microsoft.Monitor/accounts"},
		{"purview", "Microsoft.Purview/accounts"},
		{"backup", "Microsoft.RecoveryServices/vaults"},
		{"automation", "Microsoft.Automation/automationAccounts"},
		{"signalr", "Microsoft.SignalRService/SignalR"},
		{"mlworkspace", "Microsoft.MachineLearningServices/workspaces"},
	}

	// Generate 25 private endpoints
	subnets := []string{
		"/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/data-subnet",
		"/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-subnet",
		"/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/services-vnet/subnets/pls-subnet",
	}

	statuses := []string{"Approved", "Approved", "Approved", "Approved", "Pending"}

	for i := 0; i < 25; i++ {
		svc := services[i%len(services)]
		subnet := subnets[i%len(subnets)]
		status := statuses[i%len(statuses)]

		pe := models.PrivateEndpoint{
			ID:                   fmt.Sprintf("/subscriptions/demo/resourceGroups/prod-rg/providers/Microsoft.Network/privateEndpoints/pe-%s-%d", svc.prefix, i+1),
			Name:                 fmt.Sprintf("pe-%s-%d", svc.prefix, i+1),
			Location:             "eastus",
			SubnetID:             subnet,
			PrivateIPAddress:     fmt.Sprintf("10.0.%d.%d", (i/250)+1, (i%250)+10),
			PrivateLinkServiceID: fmt.Sprintf("/subscriptions/demo/resourceGroups/prod-rg/providers/%s/%s-%d", svc.service, svc.prefix, i+1),
			ConnectionState:      status,
		}
		topology.PrivateEndpoints = append(topology.PrivateEndpoints, pe)
	}

	fmt.Println("Generating realistic production topology visualization...")
	fmt.Printf("  - 3 VNets (hub-and-spoke)\n")
	fmt.Printf("  - 8 Subnets\n")
	fmt.Printf("  - 25 Private Endpoints\n")
	fmt.Printf("  - 1 VPN Gateway\n")
	fmt.Printf("  - 5 Load Balancers\n")
	fmt.Printf("  - 6 Application Gateways (5 with WAF enabled)\n")
	fmt.Printf("  - Shared NAT Gateway (3 subnets)\n")
	fmt.Printf("  - 2 NSGs (shared across subnets)\n")
	fmt.Printf("  - 1 Route Table (shared)\n")
	fmt.Println()

	// Generate DOT
	dotContent := visualization.GenerateDOTFile(topology)

	// Save DOT file
	dotFile := "/tmp/many-privatelinks-demo.dot"
	err := os.WriteFile(dotFile, []byte(dotContent), 0644)
	if err != nil {
		fmt.Printf("Error writing DOT file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ DOT file saved: %s (%.2f KB)\n", dotFile, float64(len(dotContent))/1024)

	// Render SVG
	fmt.Println("Rendering SVG...")
	svgContent, err := visualization.RenderSVG(dotContent)
	if err != nil {
		fmt.Printf("Error rendering SVG: %v\n", err)
		fmt.Println("DOT file saved, you can render it manually with: dot -Tsvg " + dotFile + " -o output.svg")
		os.Exit(1)
	}

	// Save SVG
	svgFile := "/tmp/many-privatelinks-demo.svg"
	err = os.WriteFile(svgFile, svgContent, 0644)
	if err != nil {
		fmt.Printf("Error writing SVG file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ SVG file saved: %s (%.2f KB)\n", svgFile, float64(len(svgContent))/1024)
	fmt.Println()
	fmt.Println("View the visualization:")
	fmt.Printf("  open %s\n", svgFile)
}
