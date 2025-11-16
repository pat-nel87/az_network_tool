package visualization

import (
	"fmt"
	"strings"

	"azure-network-analyzer/pkg/models"
)

// GenerateDOTFile creates a Graphviz DOT representation of the network topology
func GenerateDOTFile(topology *models.NetworkTopology) string {
	var dot strings.Builder

	dot.WriteString("digraph NetworkTopology {\n")
	dot.WriteString("  rankdir=TB;\n")
	dot.WriteString("  node [shape=box, style=filled];\n")
	dot.WriteString("  edge [fontsize=10];\n")
	dot.WriteString("  graph [fontname=\"Helvetica\", fontsize=12];\n")
	dot.WriteString("  node [fontname=\"Helvetica\", fontsize=10];\n")
	dot.WriteString("  edge [fontname=\"Helvetica\"];\n\n")

	// Title
	dot.WriteString(fmt.Sprintf("  labelloc=\"t\";\n"))
	dot.WriteString(fmt.Sprintf("  label=\"Azure Network Topology\\n%s / %s\";\n\n",
		topology.SubscriptionID, topology.ResourceGroup))

	// Track nodes for connections
	subnetNodes := make(map[string]string)
	vnetNodes := make(map[string]string)

	// Create clusters for each VNet
	for i, vnet := range topology.VirtualNetworks {
		vnetNodeID := fmt.Sprintf("vnet_%d", i)
		vnetNodes[vnet.ID] = vnetNodeID

		dot.WriteString(fmt.Sprintf("  subgraph cluster_%s {\n", sanitizeName(vnet.Name)))
		dot.WriteString(fmt.Sprintf("    label=\"%s\\n%s\";\n", vnet.Name, strings.Join(vnet.AddressSpace, "\\n")))
		dot.WriteString("    style=filled;\n")
		dot.WriteString("    color=lightblue;\n")
		dot.WriteString("    fillcolor=\"#e6f3ff\";\n\n")

		// Add subnets as nodes
		for j, subnet := range vnet.Subnets {
			subnetNodeID := fmt.Sprintf("subnet_%d_%d", i, j)
			subnetNodes[subnet.ID] = subnetNodeID

			// Determine subnet color based on associations
			color := "#90EE90" // Light green - default
			if subnet.NetworkSecurityGroup == nil {
				color = "#FFB6C1" // Light pink - no NSG (warning)
			}

			dot.WriteString(fmt.Sprintf("    %s [label=\"%s\\n%s\"", subnetNodeID, subnet.Name, subnet.AddressPrefix))
			dot.WriteString(fmt.Sprintf(", fillcolor=\"%s\"", color))
			dot.WriteString(", shape=box];\n")

			// Add NSG shield if associated
			if subnet.NetworkSecurityGroup != nil {
				nsgNodeID := fmt.Sprintf("nsg_%d_%d", i, j)
				nsgName := extractResourceName(*subnet.NetworkSecurityGroup)
				dot.WriteString(fmt.Sprintf("    %s [label=\"NSG\\n%s\", fillcolor=\"#FFE4B5\", shape=octagon];\n",
					nsgNodeID, nsgName))
				dot.WriteString(fmt.Sprintf("    %s -> %s [style=dashed, color=orange, label=\"protects\"];\n",
					nsgNodeID, subnetNodeID))
			}

			// Add Route Table connection
			if subnet.RouteTable != nil {
				rtNodeID := fmt.Sprintf("rt_%d_%d", i, j)
				rtName := extractResourceName(*subnet.RouteTable)
				dot.WriteString(fmt.Sprintf("    %s [label=\"RT\\n%s\", fillcolor=\"#DDA0DD\", shape=parallelogram];\n",
					rtNodeID, rtName))
				dot.WriteString(fmt.Sprintf("    %s -> %s [style=dotted, color=purple, label=\"routes\"];\n",
					subnetNodeID, rtNodeID))
			}

			// Add NAT Gateway connection
			if subnet.NATGateway != nil {
				natNodeID := fmt.Sprintf("nat_%d_%d", i, j)
				natName := extractResourceName(*subnet.NATGateway)
				dot.WriteString(fmt.Sprintf("    %s [label=\"NAT\\n%s\", fillcolor=\"#98FB98\", shape=diamond];\n",
					natNodeID, natName))
				dot.WriteString(fmt.Sprintf("    %s -> %s [style=solid, color=green, label=\"egress\"];\n",
					subnetNodeID, natNodeID))
			}
		}

		dot.WriteString("  }\n\n")
	}

	// Add VNet peering edges (outside clusters)
	for _, vnet := range topology.VirtualNetworks {
		for _, peering := range vnet.Peerings {
			fromNode := vnetNodes[vnet.ID]
			// Create a node for remote VNet if not in this topology
			toNode := fmt.Sprintf("remote_%s", sanitizeName(peering.RemoteVNetName))

			// Check if remote VNet is in our topology
			if remoteID, exists := vnetNodes[peering.RemoteVNetID]; exists {
				toNode = remoteID
			} else {
				// Create external VNet node
				dot.WriteString(fmt.Sprintf("  %s [label=\"%s\\n(External)\", fillcolor=\"#D3D3D3\", shape=box, style=\"filled,dashed\"];\n",
					toNode, peering.RemoteVNetName))
			}

			// Add peering edge
			peerColor := "green"
			if peering.PeeringState != "Connected" {
				peerColor = "red"
			}
			dot.WriteString(fmt.Sprintf("  %s -> %s [style=dashed, color=%s, label=\"peering\\n%s\", dir=both];\n",
				fromNode, toNode, peerColor, peering.PeeringState))
		}
	}

	// Add Load Balancers
	for i, lb := range topology.LoadBalancers {
		lbNodeID := fmt.Sprintf("lb_%d", i)
		dot.WriteString(fmt.Sprintf("  %s [label=\"LB\\n%s\\n%s\", fillcolor=\"#FFA500\", shape=ellipse];\n",
			lbNodeID, lb.Name, lb.SKU))

		// Connect to backend subnets (simplified - would need to parse backend pool IPs)
	}

	// Add Application Gateways
	for i, appgw := range topology.AppGateways {
		appgwNodeID := fmt.Sprintf("appgw_%d", i)
		wafLabel := ""
		if appgw.WAFEnabled {
			wafLabel = "\\n[WAF Enabled]"
		}
		dot.WriteString(fmt.Sprintf("  %s [label=\"AppGW\\n%s\\n%s%s\", fillcolor=\"#FF69B4\", shape=ellipse];\n",
			appgwNodeID, appgw.Name, appgw.SKU, wafLabel))
	}

	// Add VPN Gateways
	for i, vpn := range topology.VPNGateways {
		vpnNodeID := fmt.Sprintf("vpn_%d", i)
		dot.WriteString(fmt.Sprintf("  %s [label=\"VPN GW\\n%s\\n%s\", fillcolor=\"#9370DB\", shape=diamond];\n",
			vpnNodeID, vpn.Name, vpn.SKU))

		// Connect to VNet
		if vnetNode, exists := vnetNodes[vpn.VNetID]; exists {
			dot.WriteString(fmt.Sprintf("  %s -> %s [style=bold, color=purple, label=\"gateway\"];\n",
				vpnNodeID, vnetNode))
		}
	}

	// Add Private Endpoints
	for i, pe := range topology.PrivateEndpoints {
		peNodeID := fmt.Sprintf("pe_%d", i)
		targetName := extractResourceName(pe.PrivateLinkServiceID)
		dot.WriteString(fmt.Sprintf("  %s [label=\"PE\\n%s\\n→ %s\", fillcolor=\"#FFB6C1\", shape=point, width=0.3];\n",
			peNodeID, pe.Name, targetName))

		// Connect to subnet
		if subnetNode, exists := subnetNodes[pe.SubnetID]; exists {
			dot.WriteString(fmt.Sprintf("  %s -> %s [style=dotted, color=pink, label=\"private link\"];\n",
				subnetNode, peNodeID))
		}
	}

	// Add legend
	dot.WriteString("\n  // Legend\n")
	dot.WriteString("  subgraph cluster_legend {\n")
	dot.WriteString("    label=\"Legend\";\n")
	dot.WriteString("    style=filled;\n")
	dot.WriteString("    fillcolor=\"#f0f0f0\";\n")
	dot.WriteString("    node [shape=plaintext];\n")
	dot.WriteString("    legend [label=<\n")
	dot.WriteString("      <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#90EE90\">Subnet (with NSG)</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#FFB6C1\">Subnet (no NSG)</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#FFE4B5\">NSG</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#DDA0DD\">Route Table</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#9370DB\">VPN Gateway</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#FFA500\">Load Balancer</TD></TR>\n")
	dot.WriteString("      </TABLE>\n")
	dot.WriteString("    >];\n")
	dot.WriteString("  }\n")

	dot.WriteString("}\n")

	return dot.String()
}

// Helper functions

func sanitizeName(name string) string {
	// Replace characters that are invalid in DOT identifiers
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func extractResourceName(resourceID string) string {
	if resourceID == "" {
		return ""
	}
	parts := strings.Split(resourceID, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return resourceID
}
