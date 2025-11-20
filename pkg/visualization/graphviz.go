package visualization

import (
	"fmt"
	"strings"

	"azure-network-analyzer/pkg/models"
)

// VisualizationOptions controls what is included in the visualization
type VisualizationOptions struct {
	ExcludePrivateLinks bool // If true, omit private endpoints from the visualization
}

// GenerateDOTFile creates a Graphviz DOT representation of the network topology
func GenerateDOTFile(topology *models.NetworkTopology) string {
	return GenerateDOTFileWithOptions(topology, VisualizationOptions{})
}

// GenerateDOTFileWithOptions creates a Graphviz DOT representation with custom options
func GenerateDOTFileWithOptions(topology *models.NetworkTopology, opts VisualizationOptions) string {
	var dot strings.Builder

	dot.WriteString("digraph NetworkTopology {\n")
	// Layout optimization attributes
	dot.WriteString("  rankdir=TB;\n")
	dot.WriteString("  margin=0.2;          // Reduce overall margin\n")
	dot.WriteString("  pad=0.2;             // Reduce padding\n")
	dot.WriteString("  ranksep=1.0;         // Increased spacing between ranks to avoid overlaps\n")
	dot.WriteString("  nodesep=0.6;         // Increased spacing between nodes\n")
	dot.WriteString("  splines=polyline;    // Use polyline for better edge routing around obstacles\n")
	dot.WriteString("  concentrate=true;    // Merge edges where possible\n")
	dot.WriteString("  compound=true;       // Allow edges to/from clusters\n")

	// Default node and edge styles
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

		// Create a VNet header node for peering connections
		dot.WriteString(fmt.Sprintf("    %s [label=\"%s\", shape=box, style=\"filled,bold\", fillcolor=\"#d0e7ff\", fontsize=12, fontname=\"Helvetica-Bold\"];\n",
			vnetNodeID, vnet.Name))

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

	// Add Load Balancers - grouped for efficient layout
	if len(topology.LoadBalancers) > 0 {
		dot.WriteString("\n  // Load Balancers (grouped for efficient placement)\n")

		// Group load balancers to fill lower-left space efficiently
		if len(topology.LoadBalancers) <= 3 {
			// Few load balancers - put them on same rank
			dot.WriteString("  { rank=same;\n")
			for i, lb := range topology.LoadBalancers {
				lbNodeID := fmt.Sprintf("lb_%d", i)
				dot.WriteString(fmt.Sprintf("    %s [label=\"LB\\n%s\\n%s\", fillcolor=\"#FFA500\", shape=ellipse];\n",
					lbNodeID, lb.Name, lb.SKU))
			}
			dot.WriteString("  }\n")
		} else {
			// Many load balancers - distribute across multiple ranks for vertical stacking
			for i, lb := range topology.LoadBalancers {
				lbNodeID := fmt.Sprintf("lb_%d", i)
				dot.WriteString(fmt.Sprintf("  %s [label=\"LB\\n%s\\n%s\", fillcolor=\"#FFA500\", shape=ellipse];\n",
					lbNodeID, lb.Name, lb.SKU))
			}
			// Create invisible edges to control vertical stacking
			for i := 0; i < len(topology.LoadBalancers)-1; i++ {
				dot.WriteString(fmt.Sprintf("  lb_%d -> lb_%d [style=invis];\n", i, i+1))
			}
		}
	}

	// Add Application Gateways - grouped for efficient layout
	if len(topology.AppGateways) > 0 {
		dot.WriteString("\n  // Application Gateways (grouped for efficient placement)\n")

		if len(topology.AppGateways) <= 3 {
			// Few app gateways - put them on same rank
			dot.WriteString("  { rank=same;\n")
			for i, appgw := range topology.AppGateways {
				appgwNodeID := fmt.Sprintf("appgw_%d", i)
				wafLabel := ""
				if appgw.WAFEnabled {
					wafLabel = "\\n[WAF Enabled]"
				}
				dot.WriteString(fmt.Sprintf("    %s [label=\"AppGW\\n%s\\n%s%s\", fillcolor=\"#FF69B4\", shape=ellipse];\n",
					appgwNodeID, appgw.Name, appgw.SKU, wafLabel))
			}
			dot.WriteString("  }\n")
		} else {
			// Many app gateways - distribute for vertical stacking
			for i, appgw := range topology.AppGateways {
				appgwNodeID := fmt.Sprintf("appgw_%d", i)
				wafLabel := ""
				if appgw.WAFEnabled {
					wafLabel = "\\n[WAF Enabled]"
				}
				dot.WriteString(fmt.Sprintf("  %s [label=\"AppGW\\n%s\\n%s%s\", fillcolor=\"#FF69B4\", shape=ellipse];\n",
					appgwNodeID, appgw.Name, appgw.SKU, wafLabel))
			}
			// Create invisible edges to control vertical stacking
			for i := 0; i < len(topology.AppGateways)-1; i++ {
				dot.WriteString(fmt.Sprintf("  appgw_%d -> appgw_%d [style=invis];\n", i, i+1))
			}
		}
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

	// Add Azure Firewalls - grouped for efficient layout
	firewallNodes := make(map[string]string) // firewall ID -> node ID
	if len(topology.AzureFirewalls) > 0 {
		dot.WriteString("\n  // Azure Firewalls (grouped for efficient placement)\n")

		if len(topology.AzureFirewalls) <= 3 {
			// Few firewalls - put them on same rank
			dot.WriteString("  { rank=same;\n")
			for i, fw := range topology.AzureFirewalls {
				fwNodeID := fmt.Sprintf("fw_%d", i)
				firewallNodes[fw.ID] = fwNodeID
				dot.WriteString(fmt.Sprintf("    %s [label=\"Firewall\\n%s\\n%s\\n%s\", fillcolor=\"#FF6B6B\", shape=hexagon];\n",
					fwNodeID, fw.Name, fw.SKU, fw.PrivateIPAddress))
			}
			dot.WriteString("  }\n")
		} else {
			// Many firewalls - distribute for vertical stacking
			for i, fw := range topology.AzureFirewalls {
				fwNodeID := fmt.Sprintf("fw_%d", i)
				firewallNodes[fw.ID] = fwNodeID
				dot.WriteString(fmt.Sprintf("  %s [label=\"Firewall\\n%s\\n%s\\n%s\", fillcolor=\"#FF6B6B\", shape=hexagon];\n",
					fwNodeID, fw.Name, fw.SKU, fw.PrivateIPAddress))
			}
			// Create invisible edges to control vertical stacking
			for i := 0; i < len(topology.AzureFirewalls)-1; i++ {
				dot.WriteString(fmt.Sprintf("  fw_%d -> fw_%d [style=invis];\n", i, i+1))
			}
		}

		// Connect firewalls to their subnets
		for _, fw := range topology.AzureFirewalls {
			// Defensive check: ensure both nodes exist before creating edge
			fwNode, fwExists := firewallNodes[fw.ID]
			subnetNode, subnetExists := subnetNodes[fw.SubnetID]
			if fwExists && subnetExists {
				dot.WriteString(fmt.Sprintf("  %s -> %s [style=bold, color=\"#FF6B6B\", label=\"protects\"];\n",
					fwNode, subnetNode))
			}
		}

		// Connect route tables to firewalls (when routes use firewall as next hop)
		for _, rt := range topology.RouteTables {
			// Defensive check: ensure route table node exists
			rtNodeID, rtExists := routeTables[rt.ID]
			if !rtExists {
				continue // Skip if route table node doesn't exist
			}

			for _, route := range rt.Routes {
				// Check if route uses a Virtual Appliance (firewall) as next hop
				if route.NextHopType == "VirtualAppliance" && route.NextHopIPAddress != "" {
					// Find the firewall with matching private IP
					for _, fw := range topology.AzureFirewalls {
						if fw.PrivateIPAddress == route.NextHopIPAddress {
							// Defensive check: ensure firewall node exists
							fwNode, fwExists := firewallNodes[fw.ID]
							if !fwExists {
								break // Skip if firewall node doesn't exist
							}

							// Create edge showing route -> firewall for egress
							routeLabel := route.AddressPrefix
							if routeLabel == "0.0.0.0/0" {
								routeLabel = "default route"
							}
							// Use penwidth for emphasis and respect layout constraints
							dot.WriteString(fmt.Sprintf("  %s -> %s [style=bold, color=\"#FF6B6B\", penwidth=2.0, label=\"%s\\negress via FW\"];\n",
								rtNodeID, fwNode, routeLabel))
							break
						}
					}
				}
			}
		}
		dot.WriteString("\n")
	}

	// Bottom section: Legend and Private Links Table (aligned horizontally)
	dot.WriteString("\n  // Bottom section - Legend and Private Links Table (top-aligned)\n")
	dot.WriteString("  {\n")
	dot.WriteString("    rank=sink;\n") // Both at bottom, tops aligned
	dot.WriteString("    node [shape=plaintext];\n\n")

	// Add legend (left side)
	dot.WriteString("    legend [label=<\n")
	dot.WriteString("      <TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\" BGCOLOR=\"#f0f0f0\">\n")
	dot.WriteString("        <TR><TD COLSPAN=\"2\" BGCOLOR=\"#d0d0d0\"><B>Legend</B></TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#90EE90\" WIDTH=\"20\">  </TD><TD ALIGN=\"LEFT\">Subnet (with NSG)</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#FFB6C1\">  </TD><TD ALIGN=\"LEFT\">Subnet (no NSG)</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#FFE4B5\">  </TD><TD ALIGN=\"LEFT\">NSG</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#DDA0DD\">  </TD><TD ALIGN=\"LEFT\">Route Table</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#9370DB\">  </TD><TD ALIGN=\"LEFT\">VPN Gateway</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#FFA500\">  </TD><TD ALIGN=\"LEFT\">Load Balancer</TD></TR>\n")
	dot.WriteString("        <TR><TD BGCOLOR=\"#FF6B6B\">  </TD><TD ALIGN=\"LEFT\">Azure Firewall</TD></TR>\n")
	dot.WriteString("      </TABLE>\n")
	dot.WriteString("    >];\n\n")

	// Add Private Endpoints Table (center/right side) - unless excluded
	if len(topology.PrivateEndpoints) > 0 && !opts.ExcludePrivateLinks {
		dot.WriteString("    pe_table [label=<\n")
		dot.WriteString("      <TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"3\">\n")
		dot.WriteString("        <TR><TD COLSPAN=\"5\" BGCOLOR=\"#FFB6C1\"><FONT POINT-SIZE=\"11\"><B>Private Endpoints</B></FONT></TD></TR>\n")
		dot.WriteString("        <TR>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFE4E1\"><FONT POINT-SIZE=\"9\"><B>Name</B></FONT></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFE4E1\"><FONT POINT-SIZE=\"9\"><B>Target Service</B></FONT></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFE4E1\"><FONT POINT-SIZE=\"9\"><B>Subnet</B></FONT></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFE4E1\"><FONT POINT-SIZE=\"9\"><B>Private IP</B></FONT></TD>\n")
		dot.WriteString("          <TD BGCOLOR=\"#FFE4E1\"><FONT POINT-SIZE=\"9\"><B>Status</B></FONT></TD>\n")
		dot.WriteString("        </TR>\n")

		for _, pe := range topology.PrivateEndpoints {
			targetName := extractResourceName(pe.PrivateLinkServiceID)
			subnetName := extractResourceName(pe.SubnetID)
			status := pe.ConnectionState
			if status == "" {
				status = "N/A"
			}
			dot.WriteString("        <TR>\n")
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\"><FONT POINT-SIZE=\"8\">%s</FONT></TD>\n", pe.Name))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\"><FONT POINT-SIZE=\"8\">%s</FONT></TD>\n", targetName))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\"><FONT POINT-SIZE=\"8\">%s</FONT></TD>\n", subnetName))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\"><FONT POINT-SIZE=\"8\">%s</FONT></TD>\n", pe.PrivateIPAddress))
			dot.WriteString(fmt.Sprintf("          <TD ALIGN=\"LEFT\"><FONT POINT-SIZE=\"8\">%s</FONT></TD>\n", status))
			dot.WriteString("        </TR>\n")
		}

		dot.WriteString("      </TABLE>\n")
		dot.WriteString("    >];\n\n")

		// Add invisible edge to control spacing between legend and table
		dot.WriteString("    legend -> pe_table [style=invis, minlen=2];\n")
	}

	dot.WriteString("  }\n")

	dot.WriteString("}\n")

	return dot.String()
}

// Helper functions

func sanitizeName(name string) string {
	// Replace characters that are invalid in DOT identifiers
	// DOT identifiers can only contain: letters, digits, underscores
	// and must not start with a digit (but we allow it for simplicity)

	replacements := map[string]string{
		"-":  "_",
		".":  "_",
		" ":  "_",
		"(":  "_",
		")":  "_",
		"[":  "_",
		"]":  "_",
		"{":  "_",
		"}":  "_",
		":":  "_",
		";":  "_",
		",":  "_",
		"<":  "_",
		">":  "_",
		"\"": "_",
		"'":  "_",
		"/":  "_",
		"\\": "_",
		"|":  "_",
		"!":  "_",
		"@":  "_",
		"#":  "_",
		"$":  "_",
		"%":  "_",
		"^":  "_",
		"&":  "_",
		"*":  "_",
		"+":  "_",
		"=":  "_",
		"~":  "_",
		"`":  "_",
		"?":  "_",
	}

	result := name
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	// Remove any remaining non-alphanumeric characters except underscores
	var builder strings.Builder
	for _, r := range result {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			builder.WriteRune(r)
		} else {
			builder.WriteRune('_')
		}
	}

	return builder.String()
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
