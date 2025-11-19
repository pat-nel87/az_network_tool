# Fix: NAT Gateway Duplication in Graph Visualizations

## Problem
Currently, when multiple subnets share the same NAT Gateway, the visualization incorrectly renders a separate NAT Gateway node for each subnet, rather than having all subnets point to a single shared NAT Gateway node.

## Expected Behavior
- Each unique NAT Gateway (identified by its resource ID or name) should appear exactly once in the graph
- Multiple subnets using the same NAT Gateway should all have edges pointing to the same NAT Gateway node
- The NAT Gateway node should be positioned/rendered such that the connections from multiple subnets are clear

## Implementation Requirements

### 1. Deduplication Logic
- Before rendering NAT Gateway nodes, create a map/dictionary of unique NAT Gateways keyed by their resource ID
- When processing subnet-to-NAT-Gateway relationships, reference the deduplicated NAT Gateway nodes
- Ensure the NAT Gateway node contains complete information (not just data from the first subnet that references it)

### 2. Edge Connections
- Each subnet that uses the NAT Gateway should create an edge to the same NAT Gateway node
- Label edges appropriately to show the relationship (e.g., "uses NAT" or similar)
- Ensure edge routing is clear when multiple subnets connect to the same NAT Gateway

### 3. Visual Representation
- Consider the layout algorithm to position shared NAT Gateways optimally
- If using a hierarchical layout, NAT Gateways should be positioned to minimize edge crossings
- Ensure the shared nature of the NAT Gateway is visually apparent

## Testing
After implementation, verify with:
- A topology with 2+ subnets sharing the same NAT Gateway
- A topology with some subnets sharing NAT Gateways and others having dedicated ones
- Ensure node IDs are unique and edges correctly reference the deduplicated nodes

## Files Likely Involved
- Graph generation/rendering code (wherever nodes and edges are created)
- NAT Gateway data collection/processing logic
- Subnet-to-NAT-Gateway relationship mapping code