# Improvement: Load Balancer and WAF Positioning in Graph Layout

## Problem
When a resource group contains more than 2 load balancers or WAFs, they are positioned far off to the right side of the graph. This creates several issues:
- Wastes the empty space in the lower left area of the graph
- Creates an unbalanced, spread-out visualization
- Makes it harder to see all resources without horizontal scrolling
- Reduces the professional appearance of the topology

## Proposed Solution
Optimize the positioning algorithm to utilize empty space more efficiently, particularly filling the lower left area of the graph with load balancers and WAFs when there are multiple instances.

## Requirements

### 1. Load Balancer/WAF Positioning Logic
- When multiple load balancers or WAFs exist (>2), distribute them to fill available empty space
- Prioritize filling the lower left area of the graph before extending to the right
- Maintain logical grouping (e.g., by resource group or subnet association) while optimizing space usage
- Ensure connections/edges to these resources remain clear and don't create excessive crossing

### 2. Layout Algorithm Considerations
- If using a hierarchical or force-directed layout, adjust constraints to prevent rightward bunching
- Consider vertical stacking or grid arrangement for multiple similar resources
- Balance between:
  - Compact layout (minimize whitespace)
  - Readability (don't overcrowd)
  - Logical grouping (related resources near each other)

### 3. Resource Group Handling
- When a resource group has many load balancers/WAFs, they should be arranged efficiently
- Consider sub-grouping or clustering visual techniques if beneficial
- Ensure resource group boundaries (if shown) remain clear

### 4. Edge Routing
- Connections from subnets/VNets to load balancers should remain clear
- Avoid excessive edge crossings that could result from repositioning
- Use edge routing techniques (orthogonal, curved, etc.) to maintain clarity

## Testing
Verify with topologies containing:
- 3-5 load balancers in a single resource group
- 5+ WAFs across multiple resource groups
- Mix of load balancers, WAFs, and other resources
- Both small and large overall topologies

## Success Criteria
- Lower left graph area is utilized efficiently
- Multiple load balancers/WAFs don't extend excessively to the right
- Graph maintains balanced appearance across width
- All resources remain visible without horizontal scrolling (in typical viewport)
- Resource relationships remain clear

---

# Fix: Legend and Private Links Table Alignment

## Problem
The legend is currently vertically centered relative to the page, while the private links table is aligned to the left bottom. This creates visual misalignment - the legend and table should be aligned at their tops for a cohesive, professional appearance.

## Expected Behavior
- The legend's top edge should align with the private links table's top edge
- Both elements should appear at the same vertical position below the graph
- Create a clean horizontal line where both elements begin

## Implementation Requirements

### 1. Alignment Fix
- Adjust CSS/layout code to align the top of the legend with the top of the private links table
- Both should be positioned at the same distance below the main graph
- Remove any vertical centering that's currently applied to the legend

### 2. Layout Structure
Current (problematic):
```
[          Graph          ]
                              
[Private Links Table]    [Legend] (centered)
```

Desired:
```
[          Graph          ]
                              
[Private Links Table]    [Legend] (top-aligned)
```

### 3. Spacing Considerations
- Maintain consistent top margin from the graph for both elements
- Ensure adequate spacing between the table and legend (horizontal gap)
- Both elements should feel like they're part of the same "row" below the graph

### 4. Responsive Behavior
- Alignment should hold across different screen sizes
- If elements stack on narrow screens, maintain proper spacing
- PDF/print output should also show proper alignment

## Testing
Verify:
- Legend and table tops are perfectly aligned
- Layout looks professional and balanced
- Alignment holds with different amounts of private link data
- Alignment holds when legend has varying numbers of entries
- Print/PDF output maintains alignment

## Success Criteria
- Legend and private links table top edges are aligned
- Visual appearance is clean and professional
- Layout feels cohesive rather than elements floating independently

## Files Likely Involved
- HTML report template
- CSS styling (flexbox/grid layout properties)
- Legend rendering code
- Private links table rendering code

## Priority
Medium - This is a polish issue that improves the professional appearance and visual coherence of the output.