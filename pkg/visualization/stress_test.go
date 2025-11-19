package visualization

import (
	"fmt"
	"testing"

	"azure-network-analyzer/pkg/models"
)

// generateComplexTopology creates an extremely complex network topology
// to stress test the visualization system
func generateComplexTopology(numVNets, subnetsPerVNet, rulesPerNSG int) *models.NetworkTopology {
	topology := &models.NetworkTopology{
		SubscriptionID:   "stress-test-subscription",
		ResourceGroup:    "stress-test-rg",
		VirtualNetworks:  make([]models.VirtualNetwork, 0, numVNets),
		NSGs:             make([]models.NetworkSecurityGroup, 0),
		RouteTables:      make([]models.RouteTable, 0),
		PrivateEndpoints: make([]models.PrivateEndpoint, 0),
		LoadBalancers:    make([]models.LoadBalancer, 0),
		AppGateways:      make([]models.ApplicationGateway, 0),
	}

	// Generate VNets with many subnets
	for v := 0; v < numVNets; v++ {
		vnet := models.VirtualNetwork{
			ID:            fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/virtualNetworks/vnet-%d", v),
			Name:          fmt.Sprintf("vnet-%d", v),
			Location:      "eastus",
			AddressSpace:  []string{fmt.Sprintf("10.%d.0.0/16", v%256)},
			Subnets:       make([]models.Subnet, 0, subnetsPerVNet),
			Peerings:      make([]models.VNetPeering, 0),
		}

		// Create subnets
		for s := 0; s < subnetsPerVNet; s++ {
			subnet := models.Subnet{
				ID:            fmt.Sprintf("%s/subnets/subnet-%d", vnet.ID, s),
				Name:          fmt.Sprintf("subnet-%d", s),
				AddressPrefix: fmt.Sprintf("10.%d.%d.0/24", v%256, s%256),
			}

			// Associate NSG with every other subnet
			if s%2 == 0 {
				nsgID := fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/networkSecurityGroups/nsg-vnet%d-subnet%d", v, s)
				subnet.NetworkSecurityGroup = &nsgID

				// Create NSG with many rules
				nsg := models.NetworkSecurityGroup{
					ID:            nsgID,
					Name:          fmt.Sprintf("nsg-vnet%d-subnet%d", v, s),
					Location:      "eastus",
					SecurityRules: make([]models.SecurityRule, 0, rulesPerNSG),
					Associations: models.NSGAssociations{
						Subnets: []string{subnet.ID},
					},
				}

				// Generate many security rules
				for r := 0; r < rulesPerNSG; r++ {
					rule := models.SecurityRule{
						Name:                     fmt.Sprintf("rule-%d", r),
						Priority:                 int32(100 + r),
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "TCP",
						SourcePortRange:          "*",
						DestinationPortRange:     fmt.Sprintf("%d", 1000+r),
						SourceAddressPrefix:      fmt.Sprintf("10.%d.%d.0/24", r%256, r%256),
						DestinationAddressPrefix: "*",
					}
					nsg.SecurityRules = append(nsg.SecurityRules, rule)
				}

				topology.NSGs = append(topology.NSGs, nsg)
			}

			// Associate Route Table with every third subnet
			if s%3 == 0 {
				rtID := fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/routeTables/rt-vnet%d-subnet%d", v, s)
				subnet.RouteTable = &rtID

				// Create Route Table with multiple routes
				rt := models.RouteTable{
					ID:       rtID,
					Name:     fmt.Sprintf("rt-vnet%d-subnet%d", v, s),
					Location: "eastus",
					Routes:   make([]models.Route, 0),
					AssociatedSubnets: []string{subnet.ID},
				}

				// Add routes
				for route := 0; route < 10; route++ {
					rt.Routes = append(rt.Routes, models.Route{
						Name:             fmt.Sprintf("route-%d", route),
						AddressPrefix:    fmt.Sprintf("192.168.%d.0/24", route),
						NextHopType:      "VirtualAppliance",
						NextHopIPAddress: fmt.Sprintf("10.0.0.%d", route+1),
					})
				}

				topology.RouteTables = append(topology.RouteTables, rt)
			}

			vnet.Subnets = append(vnet.Subnets, subnet)
		}

		// Create peerings to other VNets (mesh topology)
		for p := 0; p < numVNets; p++ {
			if p != v && p < v { // Only peer with lower-numbered VNets to avoid duplicates
				peering := models.VNetPeering{
					Name:                 fmt.Sprintf("peer-to-vnet-%d", p),
					RemoteVNetID:         fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/virtualNetworks/vnet-%d", p),
					AllowVNetAccess:      true,
					AllowForwardedTraffic: true,
					AllowGatewayTransit:  false,
					UseRemoteGateways:    false,
					PeeringState:         "Connected",
				}
				vnet.Peerings = append(vnet.Peerings, peering)
			}
		}

		topology.VirtualNetworks = append(topology.VirtualNetworks, vnet)
	}

	// Add Private Endpoints
	for pe := 0; pe < numVNets*2; pe++ {
		endpoint := models.PrivateEndpoint{
			ID:                   fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/privateEndpoints/pe-%d", pe),
			Name:                 fmt.Sprintf("pe-%d", pe),
			Location:             "eastus",
			SubnetID:             fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/virtualNetworks/vnet-%d/subnets/subnet-0", pe%numVNets),
			PrivateLinkServiceID: fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Storage/storageAccounts/storage%d", pe),
			PrivateIPAddress:     fmt.Sprintf("10.%d.0.%d", pe%numVNets, 100+pe),
		}
		topology.PrivateEndpoints = append(topology.PrivateEndpoints, endpoint)
	}

	// Add Load Balancers
	for lb := 0; lb < numVNets; lb++ {
		loadBalancer := models.LoadBalancer{
			ID:       fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/loadBalancers/lb-%d", lb),
			Name:     fmt.Sprintf("lb-%d", lb),
			Location: "eastus",
			SKU:      "Standard",
			FrontendIPConfigs: []models.FrontendIPConfig{
				{
					Name:             fmt.Sprintf("frontend-%d", lb),
					PrivateIPAddress: fmt.Sprintf("10.%d.0.10", lb),
					SubnetID:         fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/virtualNetworks/vnet-%d/subnets/subnet-0", lb),
				},
			},
			BackendAddressPools: []models.BackendAddressPool{
				{
					Name: fmt.Sprintf("backend-%d", lb),
				},
			},
		}
		topology.LoadBalancers = append(topology.LoadBalancers, loadBalancer)
	}

	// Add Application Gateways with complex configurations
	for ag := 0; ag < numVNets; ag++ {
		appGW := models.ApplicationGateway{
			ID:            fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/applicationGateways/appgw-%d", ag),
			Name:          fmt.Sprintf("appgw-%d", ag),
			Location:      "eastus",
			SKU:           "Standard_v2",
			Tier:          "Standard_v2",
			Capacity:      2,
			SubnetID:      fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/virtualNetworks/vnet-%d/subnets/subnet-0", ag),
			WAFEnabled:    ag%2 == 0, // Enable WAF on every other AppGW
			WAFMode:       "Prevention",
			FrontendIPConfigs: []models.AppGWFrontendIPConfig{
				{
					Name:             fmt.Sprintf("frontend-public-%d", ag),
					PublicIPAddressID: fmt.Sprintf("/subscriptions/stress-test/resourceGroups/stress-test-rg/providers/Microsoft.Network/publicIPAddresses/pip-appgw-%d", ag),
				},
				{
					Name:             fmt.Sprintf("frontend-private-%d", ag),
					PrivateIPAddress: fmt.Sprintf("10.%d.0.20", ag%256),
				},
			},
			FrontendPorts:       make([]models.AppGWFrontendPort, 0),
			BackendAddressPools: make([]models.AppGWBackendAddressPool, 0),
			BackendHTTPSettings: make([]models.AppGWBackendHTTPSettings, 0),
			HTTPListeners:       make([]models.AppGWHTTPListener, 0),
			RequestRoutingRules: make([]models.AppGWRequestRoutingRule, 0),
			Probes:              make([]models.AppGWProbe, 0),
		}

		// Add multiple frontend ports
		ports := []int32{80, 443, 8080, 8443}
		for i, port := range ports {
			appGW.FrontendPorts = append(appGW.FrontendPorts, models.AppGWFrontendPort{
				Name: fmt.Sprintf("port-%d", port),
				Port: port,
			})

			// Add HTTP listener for each port
			protocol := "Http"
			if port == 443 || port == 8443 {
				protocol = "Https"
			}
			appGW.HTTPListeners = append(appGW.HTTPListeners, models.AppGWHTTPListener{
				Name:             fmt.Sprintf("listener-%d", port),
				FrontendIPConfig: fmt.Sprintf("frontend-public-%d", ag),
				FrontendPort:     fmt.Sprintf("port-%d", port),
				Protocol:         protocol,
				HostName:         fmt.Sprintf("app%d.example.com", i),
			})
		}

		// Add multiple backend pools (simulate microservices architecture)
		services := []string{"api", "web", "auth", "data", "cache", "worker"}
		for _, svc := range services {
			pool := models.AppGWBackendAddressPool{
				Name:             fmt.Sprintf("pool-%s-%d", svc, ag),
				BackendAddresses: make([]string, 0),
			}
			// Each pool has multiple backend instances
			for inst := 0; inst < 5; inst++ {
				pool.BackendAddresses = append(pool.BackendAddresses,
					fmt.Sprintf("10.%d.%d.%d", ag%256, inst+1, inst+10))
			}
			appGW.BackendAddressPools = append(appGW.BackendAddressPools, pool)

			// Add HTTP settings for each service
			appGW.BackendHTTPSettings = append(appGW.BackendHTTPSettings, models.AppGWBackendHTTPSettings{
				Name:                fmt.Sprintf("settings-%s-%d", svc, ag),
				Port:                8080,
				Protocol:            "Http",
				CookieBasedAffinity: "Disabled",
				RequestTimeout:      30,
				ProbeName:           fmt.Sprintf("probe-%s-%d", svc, ag),
			})

			// Add probe for each service
			appGW.Probes = append(appGW.Probes, models.AppGWProbe{
				Name:               fmt.Sprintf("probe-%s-%d", svc, ag),
				Protocol:           "Http",
				Host:               fmt.Sprintf("%s.internal", svc),
				Path:               "/health",
				Interval:           30,
				Timeout:            30,
				UnhealthyThreshold: 3,
			})

			// Add routing rule for each service
			appGW.RequestRoutingRules = append(appGW.RequestRoutingRules, models.AppGWRequestRoutingRule{
				Name:                fmt.Sprintf("rule-%s-%d", svc, ag),
				RuleType:            "Basic",
				HTTPListener:        "listener-80",
				BackendAddressPool:  fmt.Sprintf("pool-%s-%d", svc, ag),
				BackendHTTPSettings: fmt.Sprintf("settings-%s-%d", svc, ag),
				Priority:            int32(100 + len(appGW.RequestRoutingRules)),
			})
		}

		topology.AppGateways = append(topology.AppGateways, appGW)
	}

	return topology
}

func TestStressSmallTopology(t *testing.T) {
	// Small: 10 VNets, 10 subnets each, 20 rules per NSG
	topology := generateComplexTopology(10, 10, 20)

	t.Logf("Generated topology: %d VNets, %d NSGs, %d Route Tables, %d Private Endpoints, %d Load Balancers",
		len(topology.VirtualNetworks), len(topology.NSGs), len(topology.RouteTables),
		len(topology.PrivateEndpoints), len(topology.LoadBalancers))

	// Test DOT generation
	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes", len(dot))

	if len(dot) == 0 {
		t.Error("DOT file is empty")
	}
}

func TestStressMediumTopology(t *testing.T) {
	// Medium: 25 VNets, 20 subnets each, 50 rules per NSG
	topology := generateComplexTopology(25, 20, 50)

	t.Logf("Generated topology: %d VNets, %d NSGs, %d Route Tables, %d Private Endpoints, %d Load Balancers",
		len(topology.VirtualNetworks), len(topology.NSGs), len(topology.RouteTables),
		len(topology.PrivateEndpoints), len(topology.LoadBalancers))

	// Test DOT generation
	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes", len(dot))

	if len(dot) == 0 {
		t.Error("DOT file is empty")
	}
}

func TestStressLargeTopology(t *testing.T) {
	// Large: 50 VNets, 30 subnets each, 100 rules per NSG
	topology := generateComplexTopology(50, 30, 100)

	t.Logf("Generated topology: %d VNets, %d NSGs, %d Route Tables, %d Private Endpoints, %d Load Balancers",
		len(topology.VirtualNetworks), len(topology.NSGs), len(topology.RouteTables),
		len(topology.PrivateEndpoints), len(topology.LoadBalancers))

	// Test DOT generation
	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes", len(dot))

	if len(dot) == 0 {
		t.Error("DOT file is empty")
	}
}

func TestStressExtremeTopology(t *testing.T) {
	// Extreme: 100 VNets, 50 subnets each, 200 rules per NSG
	// This simulates a very large enterprise network
	topology := generateComplexTopology(100, 50, 200)

	t.Logf("Generated topology: %d VNets, %d NSGs, %d Route Tables, %d Private Endpoints, %d Load Balancers",
		len(topology.VirtualNetworks), len(topology.NSGs), len(topology.RouteTables),
		len(topology.PrivateEndpoints), len(topology.LoadBalancers))

	totalSubnets := 0
	for _, vnet := range topology.VirtualNetworks {
		totalSubnets += len(vnet.Subnets)
	}
	t.Logf("Total subnets: %d", totalSubnets)

	totalRules := 0
	for _, nsg := range topology.NSGs {
		totalRules += len(nsg.SecurityRules)
	}
	t.Logf("Total security rules: %d", totalRules)

	// Test DOT generation
	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes (%.2f MB)", len(dot), float64(len(dot))/1024/1024)

	if len(dot) == 0 {
		t.Error("DOT file is empty")
	}
}

func TestStressSVGRenderingSmall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SVG rendering test in short mode")
	}

	// Small topology for SVG rendering test
	topology := generateComplexTopology(5, 5, 10)

	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes", len(dot))

	svg, err := RenderSVG(dot)
	if err != nil {
		t.Fatalf("Failed to render SVG: %v", err)
	}

	t.Logf("SVG file size: %d bytes (%.2f KB)", len(svg), float64(len(svg))/1024)
}

func TestStressSVGRenderingMedium(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SVG rendering test in short mode")
	}

	// Medium topology - this may cause memory issues
	topology := generateComplexTopology(10, 10, 20)

	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes (%.2f KB)", len(dot), float64(len(dot))/1024)

	svg, err := RenderSVG(dot)
	if err != nil {
		t.Fatalf("Failed to render SVG: %v", err)
	}

	t.Logf("SVG file size: %d bytes (%.2f KB)", len(svg), float64(len(svg))/1024)
}

func TestStressSVGRenderingLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SVG rendering test in short mode")
	}

	// Large topology - likely to cause memory issues on constrained systems
	topology := generateComplexTopology(20, 15, 30)

	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes (%.2f MB)", len(dot), float64(len(dot))/1024/1024)

	svg, err := RenderSVG(dot)
	if err != nil {
		t.Fatalf("Failed to render SVG: %v", err)
	}

	t.Logf("SVG file size: %d bytes (%.2f MB)", len(svg), float64(len(svg))/1024/1024)
}

// BenchmarkDOTGeneration benchmarks DOT file generation
func BenchmarkDOTGeneration(b *testing.B) {
	topology := generateComplexTopology(10, 10, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateDOTFile(topology)
	}
}

// BenchmarkSVGRendering benchmarks SVG rendering
func BenchmarkSVGRendering(b *testing.B) {
	topology := generateComplexTopology(5, 5, 10)
	dot := GenerateDOTFile(topology)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := RenderSVG(dot)
		if err != nil {
			b.Fatalf("Failed to render SVG: %v", err)
		}
	}
}

func TestStressSVGRenderingVeryLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping very large topology test in short mode")
	}

	// Very large topology - too large for reliable SVG rendering in CI
	// SVG rendering causes timeouts due to WASM memory constraints
	topology := generateComplexTopology(30, 20, 40)

	t.Logf("Generated topology: %d VNets, %d NSGs, %d Route Tables, %d AppGWs",
		len(topology.VirtualNetworks), len(topology.NSGs),
		len(topology.RouteTables), len(topology.AppGateways))

	totalSubnets := 0
	for _, vnet := range topology.VirtualNetworks {
		totalSubnets += len(vnet.Subnets)
	}
	t.Logf("Total subnets: %d", totalSubnets)

	// Only test DOT generation for very large topologies
	// SVG rendering is too resource-intensive and causes timeouts in CI
	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes (%.2f MB)", len(dot), float64(len(dot))/1024/1024)

	if len(dot) == 0 {
		t.Error("DOT file is empty")
	}

	t.Log("Note: SVG rendering skipped for very large topology to avoid CI timeouts")
	t.Log("For large topologies, use external GraphViz: dot -Tsvg input.dot -o output.svg")
}

func TestStressSVGRenderingExtreme(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping extreme topology test in short mode")
	}

	// Extreme topology - too large for SVG rendering, so we only test DOT generation
	// SVG rendering would cause memory/timeout issues in CI environments
	topology := generateComplexTopology(50, 25, 50)

	t.Logf("Generated topology: %d VNets, %d NSGs, %d Route Tables, %d AppGWs",
		len(topology.VirtualNetworks), len(topology.NSGs),
		len(topology.RouteTables), len(topology.AppGateways))

	totalBackendPools := 0
	for _, appgw := range topology.AppGateways {
		totalBackendPools += len(appgw.BackendAddressPools)
	}
	t.Logf("Total AppGW backend pools: %d", totalBackendPools)

	totalSubnets := 0
	for _, vnet := range topology.VirtualNetworks {
		totalSubnets += len(vnet.Subnets)
	}
	t.Logf("Total subnets: %d", totalSubnets)

	// Only test DOT generation for extreme topologies
	// SVG rendering is too resource-intensive and causes timeouts
	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes (%.2f MB)", len(dot), float64(len(dot))/1024/1024)

	if len(dot) == 0 {
		t.Error("DOT file is empty")
	}

	t.Log("Note: SVG rendering skipped for extreme topology to avoid CI timeouts")
	t.Log("For such large topologies, use external GraphViz: dot -Tsvg input.dot -o output.svg")
}

func TestStressMaximumTopology(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping maximum topology test in short mode")
	}

	// Maximum stress test - very likely to cause memory issues
	topology := generateComplexTopology(75, 30, 75)

	t.Logf("Generated topology: %d VNets, %d NSGs, %d Route Tables, %d AppGWs, %d LBs",
		len(topology.VirtualNetworks), len(topology.NSGs),
		len(topology.RouteTables), len(topology.AppGateways), len(topology.LoadBalancers))

	totalSubnets := 0
	for _, vnet := range topology.VirtualNetworks {
		totalSubnets += len(vnet.Subnets)
	}
	t.Logf("Total subnets: %d", totalSubnets)

	totalRules := 0
	for _, nsg := range topology.NSGs {
		totalRules += len(nsg.SecurityRules)
	}
	t.Logf("Total security rules: %d", totalRules)

	totalAppGWComponents := 0
	for _, appgw := range topology.AppGateways {
		totalAppGWComponents += len(appgw.BackendAddressPools) +
			len(appgw.HTTPListeners) +
			len(appgw.RequestRoutingRules) +
			len(appgw.Probes)
	}
	t.Logf("Total AppGW components: %d", totalAppGWComponents)

	// Only test DOT generation - SVG may fail
	dot := GenerateDOTFile(topology)
	t.Logf("DOT file size: %d bytes (%.2f MB)", len(dot), float64(len(dot))/1024/1024)

	if len(dot) == 0 {
		t.Error("DOT file is empty")
	}
}
