# Investigation & Feature: Azure Firewall Visualization Support

## Problem
Azure Firewall resources are not appearing in the network topology visualizations. After replacing NAT Gateways with Azure Firewall in the network topology, the firewall and its associated resources are not being rendered in the graphs.

## Investigation Required

### 1. Determine Current State
First, investigate why Azure Firewall is not showing up:

**Check Data Collection:**
- Does the Azure SDK/API client collect Azure Firewall resources?
- Review the data collection code to see if Azure Firewall is being queried
- Check if firewall data is present in the raw data structures but not rendered
- Verify API permissions/scopes include firewall resources

**Check Data Processing:**
- Is Azure Firewall data being parsed and stored in internal data structures?
- Review any filtering logic that might exclude firewall resources
- Check if there's a resource type whitelist/blacklist that excludes firewalls

**Check Visualization:**
- Is there rendering logic for Azure Firewall nodes?
- Check if firewall is mapped to a node type in the graph generation
- Verify if firewall appears in the legend definitions

**Likely Locations to Investigate:**
- Azure resource enumeration/discovery code
- Resource type handling/mapping
- Graph node creation logic
- Legend/node type definitions

### 2. Document Findings
Create a summary of what was found:
- Is Azure Firewall data being collected? (Yes/No)
- Is it being processed/stored? (Yes/No)
- Is it being rendered? (Yes/No)
- What is blocking the visualization? (missing collection, missing rendering, filtered out, etc.)

## Feature Requirements: Add Azure Firewall to Visualizations

Once the investigation is complete, implement full Azure Firewall support:

### 1. Data Collection
If not already collected, add Azure Firewall resource enumeration:
- Query Azure Firewall resources from Azure Resource Management API
- Collect relevant properties:
  - Firewall name
  - Resource group
  - Location/region
  - SKU/tier (Standard, Premium, Basic)
  - Associated VNet/subnet (typically "AzureFirewallSubnet")
  - Public IP address(es)
  - Private IP address
  - Firewall policy association (if applicable)
  - DNS settings
  - Threat intelligence mode
  - Provisioning state

### 2. Associated Resources
Azure Firewall requires and interacts with several related resources that should also be represented:

**Required Associated Resources:**
- **AzureFirewallSubnet**: The dedicated subnet where firewall is deployed (usually has this specific name)
- **Public IP Address(es)**: Firewall's public IP for outbound connectivity
- **Azure Firewall Policy**: If using firewall policies instead of classic rules
- **Route Tables**: UDRs directing traffic to the firewall (0.0.0.0/0 → firewall private IP)

**Optional/Related Resources:**
- **Azure Firewall Management Subnet**: If using forced tunneling
- **DDoS Protection Plan**: If associated with the VNet
- **NAT Rules, Network Rules, Application Rules**: Could be shown in a table similar to private links

### 3. Graph Visualization

**Node Representation:**
- Create a distinct node type for Azure Firewall
- Suggested icon/color: Shield or firewall icon, orange/red color scheme
- Node should display:
  - Firewall name
  - SKU tier (Standard/Premium/Basic)
  - Private IP address
  - Status indicator if available

**Positioning:**
- Firewall should be positioned near its associated subnet (AzureFirewallSubnet)
- Consider hub positioning if in a hub-and-spoke topology
- Avoid the rightward bunching issue that load balancers had

**Connections/Edges:**
- Edge from Azure Firewall to its subnet (AzureFirewallSubnet)
- Edges to public IP address(es)
- Edge to associated firewall policy (if exists)
- Edges from route tables that point to the firewall (show UDR relationships)
- Consider showing traffic flow direction (inbound/outbound) with arrow styles

### 4. Legend Integration
Add Azure Firewall to the legend:
- Icon/color matching the node representation
- Label: "Azure Firewall" with optional SKU indication
- Position appropriately in legend hierarchy (likely near load balancers, NAT gateways)

### 5. Hub-and-Spoke Considerations
Azure Firewall is commonly used in hub-and-spoke topologies:
- If firewall is in a hub VNet, highlight this relationship
- Show spoke VNets routing traffic through the hub firewall
- Consider adding visual indication of traffic flow through firewall (optional enhancement)

### 6. Firewall Rules/Policies Display
Similar to the private links table approach, consider displaying firewall configuration:

**Option A: Simple Approach**
- Just show the firewall node with basic info
- Users can inspect Azure Portal for rules

**Option B: Summary Table (Recommended)**
- Create a "Firewall Policies" or "Firewall Rules Summary" table
- Show high-level info:
  - Number of NAT rules
  - Number of network rules
  - Number of application rules
  - Threat intelligence mode
  - DNS proxy status
- Keep it summary-level to avoid overwhelming detail

**Option C: Detailed Tables**
- Full breakdown of rules (may be too detailed)
- Consider only if requested by users

### 7. Route Table Integration
Azure Firewall typically works with User-Defined Routes (UDRs):
- Ensure route tables pointing to firewall are visualized
- Show routes with destination 0.0.0.0/0 and next hop = firewall IP
- Visual indication that firewall is the egress point for certain subnets

## Testing

### Test Scenarios:
1. **Basic Firewall**: Single Azure Firewall in a VNet with AzureFirewallSubnet
2. **Hub-and-Spoke**: Firewall in hub VNet with multiple spoke VNets
3. **Firewall Policy**: Firewall associated with Azure Firewall Policy resource
4. **Multiple Public IPs**: Firewall with multiple public IP addresses
5. **Forced Tunneling**: Firewall with management subnet configuration
6. **Premium SKU**: Verify SKU/tier is displayed correctly

### Verification:
- Firewall appears in visualization
- All associations (subnet, IPs, policy) are shown
- Route tables referencing firewall display connections correctly
- Legend includes firewall
- Layout remains clean (no excessive spacing issues)
- Firewall node is visually distinct from other resource types

## Success Criteria
- Azure Firewall resources are automatically discovered and collected
- Firewall appears in network topology graphs with clear representation
- Associated resources (subnet, public IPs, policies) are connected
- Route table relationships to firewall are visualized
- Legend includes Azure Firewall with appropriate icon/color
- Layout accommodates firewall without creating whitespace issues
- Firewall integrates cleanly into hub-and-spoke topologies

## Documentation Updates
After implementation:
- Update README.md to mention Azure Firewall support in features list
- Add an example visualization showing Azure Firewall (good candidate for the examples gallery)
- Document any configuration flags related to firewall visualization

## Files Likely Involved
- Azure resource collection/enumeration code
- Resource type definitions and mappings
- Graph node generation logic
- Legend definition file
- Node styling/coloring configuration
- Edge/connection generation logic
- HTML report template (if adding firewall rules table)

## Priority
High - Azure Firewall is a critical security component in Azure networking, and its absence from visualizations is a significant gap for users implementing secure network architectures.

## Notes
- Azure Firewall is a premium security resource - its visualization should reflect its importance in the network architecture
- Consider future enhancements like showing firewall DNAT rules as connections to backend resources
- May want to add a flag to show/hide detailed firewall rule information (similar to private links flag)