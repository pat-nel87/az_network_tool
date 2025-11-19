# Improvement: Private Links Display

## Problem
Private links can quickly overwhelm graph visualizations, making the topology difficult to read and navigate. When networks contain many private links, the visual clutter obscures the primary network structure.

## Proposed Solution
Display private links in a tabular format instead of as nodes in the graph visualization.

## Requirements

### 1. Remove Private Links from Graph
- Private link nodes should not be rendered in the main network topology visualization
- Remove edges connecting resources to private links from the graph
- Ensure removal of private links doesn't break the overall graph layout or connectivity

### 2. Create Private Links Table
Display private links in a separate table with the following information:
- **Resource/Service Name**: The Azure resource or service the private link connects to
- **Private Link Name**: The name/identifier of the private link itself
- **Associated Resource**: What resource/subnet is using this private link
- **Connection Status**: If available (e.g., Approved, Pending, Rejected)
- **Private IP Address**: If applicable

### 3. Table Presentation
- Position the table in a logical location (e.g., below the graph, in a separate section, or in a collapsible panel)
- Make the table sortable and/or filterable if there are many entries
- Consider grouping by VNet or subnet if that provides clarity
- Ensure the table is included in generated reports

### 4. Handling Edge Cases
- If a network has zero private links, don't show an empty table (or show "No private links found")
- For networks with 1-3 private links, the table is still preferable to maintain consistency

## Success Criteria
- Graph visualizations remain clean and readable even with 20+ private links
- Users can quickly identify what private links exist and what they connect to
- The table provides sufficient detail without being overwhelming

---

# Fix: VNet Peering Diagram Labels

## Problem
VNet peering diagrams currently show generic labels like "vnet_0", "vnet_1", etc., which makes it unclear which actual VNet is being referenced. Users need to see the actual VNet names to understand the peering relationships.

## Expected Behavior
- Peering diagram nodes should be labeled with the actual VNet name (e.g., "prod-vnet-eastus", "hub-vnet")
- Labels should be the VNet's display name or resource name as it appears in Azure
- The label should be clear and readable in the visualization

## Implementation Requirements

### 1. Label Mapping
- Replace generic "vnet_0" style labels with actual VNet names
- Ensure the mapping between VNet data and graph nodes preserves the correct name
- If VNet names are very long, consider truncation with tooltip showing full name (optional enhancement)

### 2. Consistency
- Use the same naming convention throughout the visualization
- If showing VNet names in the main topology, use identical formatting in the peering diagram
- Ensure subscription or resource group context is clear if VNets have duplicate names across subscriptions

### 3. Visual Clarity
- Ensure text is legible (appropriate font size, no overlapping)
- If using abbreviations or truncation, provide a way to see the full name
- Consider showing additional context like resource group in parentheses if helpful: "prod-vnet (rg-production)"

## Testing
Verify with:
- Single peering relationship between two VNets
- Complex hub-and-spoke topology with multiple peerings
- VNets with long names
- VNets with similar names across different resource groups

## Files Likely Involved
- Peering diagram generation code
- VNet data structure/model where names are stored
- Graph node labeling logic

---

# Priority
Both issues are high priority as they significantly impact usability:
1. Private links table (high impact on visual clarity)
2. VNet peering labels (essential for understanding topology)

Fixing these will make the application production-ready and significantly improve user experience.