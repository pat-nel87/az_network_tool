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

	// Deduplicate NAT Gateways, NSGs, and Route Tables
	natGateways := make(map[string]string) // resource ID -> node ID
	nsgs := make(map[string]string)
	routeTables := make(map[string]string)

	// Create clusters for each VNet
	for i, vnet := range topology.VirtualNetworks {
		vnetNodeID := fmt.Sprintf("vnet_%d", i)
		vnetNodes[vnet.ID] = vnetNodeID

		dot.WriteString(fmt.Sprintf("  subgraph cluster_%s {\n", sanitizeName(vnet.Name)))
		dot.WriteString(fmt.Sprintf("    label=\"VNet: %s\\n%s\";\n", vnet.Name, strings.Join(vnet.AddressSpace, "\\n")))
		dot.WriteString("    style=filled;\n")
		dot.WriteString("    color=lightblue;\n")
		dot.WriteString("    fillcolor=\"#e6f3ff\";\n")
		dot.WriteString("    fontsize=14;\n")
		dot.WriteString("    fontname=\"Helvetica-Bold\";\n\n")

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

			// Register NSG for deduplication
			if subnet.NetworkSecurityGroup != nil {
				nsgID := *subnet.NetworkSecurityGroup
				if _, exists := nsgs[nsgID]; !exists {
					nsgs[nsgID] = fmt.Sprintf("nsg_%s", sanitizeName(extractResourceName(nsgID)))
				}
			}

			// Register Route Table for deduplication
			if subnet.RouteTable != nil {
				rtID := *subnet.RouteTable
				if _, exists := routeTables[rtID]; !exists {
					routeTables[rtID] = fmt.Sprintf("rt_%s", sanitizeName(extractResourceName(rtID)))
				}
			}

			// Register NAT Gateway for deduplication
			if subnet.NATGateway != nil {
				natID := *subnet.NATGateway
				if _, exists := natGateways[natID]; !exists {
					natGateways[natID] = fmt.Sprintf("nat_%s", sanitizeName(extractResourceName(natID)))
				}
			}
		}

		dot.WriteString("  }\n\n")
	}

	// Render deduplicated NSGs (outside clusters)
	for nsgID, nsgNodeID := range nsgs {
		nsgName := extractResourceName(nsgID)
		dot.WriteString(fmt.Sprintf("  %s [label=\"NSG\\n%s\", fillcolor=\"#FFE4B5\", shape=octagon];\n",
			nsgNodeID, nsgName))
	}

	// Render deduplicated Route Tables (outside clusters)
	for rtID, rtNodeID := range routeTables {
		rtName := extractResourceName(rtID)
		dot.WriteString(fmt.Sprintf("  %s [label=\"Route Table\\n%s\", fillcolor=\"#DDA0DD\", shape=parallelogram];\n",
			rtNodeID, rtName))
	}

	// Render deduplicated NAT Gateways (outside clusters)
	for natID, natNodeID := range natGateways {
		natName := extractResourceName(natID)
		dot.WriteString(fmt.Sprintf("  %s [label=\"NAT Gateway\\n%s\", fillcolor=\"#98FB98\", shape=diamond];\n",
			natNodeID, natName))
	}

	// Connect subnets to their NSGs, Route Tables, and NAT Gateways
	for _, vnet := range topology.VirtualNetworks {
		for _, subnet := range vnet.Subnets {
			subnetNodeID := subnetNodes[subnet.ID]

			// Connect to NSG
			if subnet.NetworkSecurityGroup != nil {
				nsgNodeID := nsgs[*subnet.NetworkSecurityGroup]
				dot.WriteString(fmt.Sprintf("  %s -> %s [style=dashed, color=orange, label=\"protects\"];\n",
					nsgNodeID, subnetNodeID))
			}

			// Connect to Route Table
			if subnet.RouteTable != nil {
				rtNodeID := routeTables[*subnet.RouteTable]
				dot.WriteString(fmt.Sprintf("  %s -> %s [style=dotted, color=purple, label=\"routes\"];\n",
					subnetNodeID, rtNodeID))
			}

			// Connect to NAT Gateway
			if subnet.NATGateway != nil {
				natNodeID := natGateways[*subnet.NATGateway]
				dot.WriteString(fmt.Sprintf("  %s -> %s [style=solid, color=green, label=\"egress\"];\n",
					subnetNodeID, natNodeID))
			}
		}
	}
	dot.WriteString("\n")

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

	// Private Endpoints will be shown in a table instead of as nodes
	// (removed from graph for clarity)

	// Add Private Endpoints Table
	if len(topology.PrivateEndpoints) > 0 {
		dot.WriteString("\n  // Private Endpoints Table\n")
		dot.WriteString("  subgraph cluster_private_endpoints {\n")
		dot.WriteString("    label=\"Private Endpoints\";\n")
		dot.WriteString("    style=filled;\n")
		dot.WriteString("    fillcolor=\"#fff9e6\";\n")
		dot.WriteString("    fontsize=12;\n")
		dot.WriteString("    fontname=\"Helvetica-Bold\";\n")
		dot.WriteString("    node [shape=plaintext];\n")
		dot.WriteString("    pe_table [label=<\n")
		dot.WriteString("      <TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
		dot.WriteString("        <TR>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFB6C1\"><B>Name</B></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFB6C1\"><B>Target Service</B></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFB6C1\"><B>Subnet</B></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFB6C1\"><B>Private IP</B></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFB6C1\"><B>Status</B></TD>\n")
		dot.WriteString("        </TR>\n")

		for _, pe := range topology.PrivateEndpoints {
			targetName := extractResourceName(pe.PrivateLinkServiceID)
			subnetName := extractResourceName(pe.SubnetID)
			status := pe.ConnectionState
			if status == "" {
				status = "N/A"
			}
			dot.WriteString("        <TR>\n")
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\">%s</TD>\n", pe.Name))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\">%s</TD>\n", targetName))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\">%s</TD>\n", subnetName))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\">%s</TD>\n", pe.PrivateIPAddress))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\">%s</TD>\n", status))
			dot.WriteString("        </TR>\n")
		}

		dot.WriteString("      </TABLE>\n")
		dot.WriteString("    >];\n")
		dot.WriteString("  }\n\n")
	}

	// Add legend
	dot.WriteString("  // Legend\n")
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
