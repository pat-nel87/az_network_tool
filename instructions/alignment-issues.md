# Improvement: Visualization Layout and Spacing Optimization

## Problem
The current visualization layout has inefficient use of space:
- The main network graph is centered, leaving excessive empty white space on the left side
- The private links table is bottom-aligned left, creating visual disconnect
- The legend is offset far to the right of the graph
- Overall layout feels unbalanced and wastes valuable screen/page real estate

## Proposed Solution
Optimize the layout to use space more efficiently and create a more cohesive, compact visualization.

## Requirements

### 1. Graph Positioning
- Move the main network topology graph closer to the left side of the page
- Reduce or eliminate excessive white space on the left
- Maintain reasonable left margin for readability (suggest 20-40px rather than center-alignment)
- Ensure the graph still has room to render properly without being cramped

### 2. Legend Positioning
- Position the legend adjacent to the graph rather than far to the right
- Options to consider:
  - Top-right corner of the visualization area
  - Right side, vertically aligned with the top of the graph
  - Below the graph, aligned with its right edge
- Ensure legend remains readable and doesn't overlap with graph elements
- Legend should feel integrated with the visualization, not floating separately

### 3. Private Links Table Layout
- Keep the table below the graph (current positioning is good)
- Consider aligning the table with the left edge of the graph for visual consistency
- Ensure table width is appropriate (doesn't stretch unnecessarily wide)
- Maintain clear visual separation between graph and table

### 4. Overall Layout Goals
- Create a more compact, professionally laid-out visualization
- Maximize information density without sacrificing readability
- Ensure layout works well for both:
  - Screen viewing (HTML reports)
  - Printed/PDF output
- Responsive behavior if viewing area is constrained

### 5. Configuration Flag: Exclude Private Links from Visualization

Add a command-line flag or configuration option to control private link inclusion:

**Flag Specification:**
- Flag name: `--exclude-private-links` or `--no-private-links` (or similar)
- Behavior: When enabled, completely omit private links from the visualization
- Scope: Should exclude both:
  - Private link nodes/edges from the graph (already being removed per previous issue)
  - The private links table from the output

**Use Cases:**
- Networks with many private links where the table isn't needed
- Simplified topology views focusing on core network structure
- Presentations where private link details aren't relevant
- Faster rendering for large topologies

**Implementation Notes:**
- Default behavior: Include private links table (backwards compatible)
- When flag is set: Skip private link data collection/processing entirely (performance optimization)
- Configuration file option: If you have a config file, add corresponding setting (e.g., `exclude_private_links: true`)
- Help text: Clearly document the flag in CLI help output

**Example Usage:**
```bash
# Include private links (default)
./network-analyzer analyze --subscription <sub-id>

# Exclude private links
./network-analyzer analyze --subscription <sub-id> --exclude-private-links
```

## Visual Design Considerations
- Consider a grid-based layout approach:
  - Graph occupies left 70-80% of width
  - Legend occupies remaining right space
  - Table spans full width below
- Ensure consistent margins and padding throughout
- Use visual hierarchy to guide the eye (graph → legend → table)

## Testing
Verify layout with:
- Small topology (few resources) - ensure it doesn't look lost on the page
- Large topology (many resources) - ensure everything fits without excessive scrolling
- Medium topology - ensure balanced, professional appearance
- Different screen sizes/resolutions
- PDF/print output quality
- Topologies with and without private links (test the exclude flag)

## Success Criteria
- No excessive empty white space on the left side
- Legend feels integrated with the graph, not floating
- Private links table visually connected to the overall layout
- Professional, polished appearance suitable for stakeholder presentations
- User has control over whether private links appear in output

## Files Likely Involved
- HTML report generation/templating code
- CSS styling for visualization layout
- Graph rendering positioning logic
- Report export/PDF generation
- CLI flag parsing and configuration handling
- Private links data processing (conditional logic based on flag)

## Priority
Medium-High - This is a polish issue that significantly improves the professional appearance and usability of the tool, making it more suitable for executive presentations and reports.