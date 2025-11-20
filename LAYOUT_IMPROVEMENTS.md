# Visualization Layout Improvements

## Overview
Implemented comprehensive layout optimizations to improve the visual appearance and space efficiency of network topology visualizations.

## Changes Implemented

### 1. ✅ Layout Optimization Attributes
Added GraphViz layout attributes to reduce wasted space and improve overall layout:

```graphviz
digraph NetworkTopology {
  rankdir=TB;
  margin=0.2;          // Reduced overall margin
  pad=0.2;             // Reduced padding
  ranksep=0.75;        // Optimized spacing between ranks
  nodesep=0.5;         // Optimized spacing between nodes
  splines=ortho;       // Orthogonal lines for cleaner look
  concentrate=true;    // Merge edges where possible
```

**Benefits:**
- Eliminated excessive white space on the left side
- More compact, professional layout
- Better use of screen/page real estate
- Cleaner edge routing with orthogonal splines

### 2. ✅ CLI Flag: --exclude-private-links
Added new command-line flag to control private endpoint visibility:

```bash
# Include private endpoints (default)
az-network-analyzer analyze -s <sub-id> -g <rg>

# Exclude private endpoints from visualization
az-network-analyzer analyze -s <sub-id> -g <rg> --exclude-private-links
```

**Use Cases:**
- Simplified topology views focusing on core network structure
- Networks with many private links where the table adds clutter
- Presentations where private link details aren't relevant
- Faster rendering for large topologies

**Implementation:**
- Added `VisualizationOptions` struct to control rendering
- `GenerateDOTFile()` maintains backward compatibility
- `GenerateDOTFileWithOptions()` accepts custom options
- Private links table is completely omitted when flag is set

### 3. ✅ Legend Repositioning
Moved legend from far right to bottom-left of the diagram:

**Before:** Legend was in a cluster, positioned far to the right (especially with many resources)
**After:** Legend positioned at bottom-left using `rank=sink`, integrated with the visualization

```graphviz
// Legend positioned at bottom-left
legend [shape=plaintext, label=<
  <TABLE ...>
    ...
  </TABLE>
>];
{ rank=sink; legend; }
```

**Benefits:**
- Legend feels integrated with the visualization, not floating
- Consistently positioned regardless of topology size
- Better space utilization - doesn't push other elements around
- More professional, balanced appearance

### 4. ✅ Private Links Table Alignment
Private links table remains bottom-center with improved positioning:

```graphviz
// Private Endpoints Table (bottom of diagram)
{
  rank=sink;  // Force to bottom
  node [shape=plaintext];
  pe_table [label=<TABLE>...];
}
```

**Features:**
- Compact table format (8pt data, 9pt headers, 11pt title)
- Shows: Name, Target Service, Subnet, Private IP, Status
- Centered at bottom for visual balance
- Only rendered when private endpoints exist AND not excluded

### 5. ✅ Visual Layout Hierarchy
The new layout creates a clear visual hierarchy:

1. **Title** (top-center)
2. **Main Network Graph** (left-aligned, reduced margins)
3. **Legend** (bottom-left)
4. **Private Endpoints Table** (bottom-center, if present)

## Testing Results

Tested across different topology sizes:

### Small Topology (few resources)
- ✓ Legend properly positioned at bottom-left
- ✓ No excessive white space
- ✓ Professional appearance

### Medium Topology (standard production network)
- ✓ Well-balanced layout
- ✓ Legend integrated with graph
- ✓ Private links table readable

### Large Topology (25+ private endpoints, many resources)
- ✓ Legend stays bottom-left (not pushed right)
- ✓ Graph well-aligned to left
- ✓ Private links table compact and readable
- ✓ All elements visible without excessive scrolling

### Exclude Private Links
- ✓ Legend still bottom-left
- ✓ Table completely omitted
- ✓ Cleaner visualization for presentations

## Visual Comparison

### Before
- ❌ Excessive left margin/white space
- ❌ Legend far to the right, floating
- ❌ No control over private link visibility
- ❌ Wasted screen real estate

### After
- ✅ Efficient left alignment
- ✅ Legend integrated at bottom-left
- ✅ Optional private link exclusion
- ✅ Professional, compact layout
- ✅ Better for presentations and reports

## Files Modified

1. **pkg/visualization/graphviz.go**
   - Added `VisualizationOptions` struct
   - Added `GenerateDOTFileWithOptions()` function
   - Optimized layout attributes (margin, pad, ranksep, nodesep, splines)
   - Repositioned legend to bottom-left
   - Made private links table conditional

2. **cmd/analyze.go**
   - Added `--exclude-private-links` flag
   - Pass options to visualization generation
   - Display message when private links are excluded

## Usage Examples

```bash
# Standard visualization (optimized layout, includes private links)
az-network-analyzer analyze -s my-sub -g my-rg --viz-format svg

# Exclude private links for cleaner presentation
az-network-analyzer analyze -s my-sub -g my-rg --viz-format svg --exclude-private-links

# Generate different formats with optimized layout
az-network-analyzer analyze -s my-sub -g my-rg --viz-format png
az-network-analyzer analyze -s my-sub -g my-rg --viz-format pdf

# Docker usage
docker run --rm \
  -v $(pwd)/output:/output \
  -e AZURE_TENANT_ID=$AZURE_TENANT_ID \
  -e AZURE_CLIENT_ID=$AZURE_CLIENT_ID \
  -e AZURE_CLIENT_SECRET=$AZURE_CLIENT_SECRET \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze -s $SUB_ID -g $RG --exclude-private-links -f /output/topology.md
```

## Performance Impact

- **Rendering Speed:** Slightly faster due to optimized layout algorithm
- **File Size:** Comparable (optimized attributes may reduce DOT file size slightly)
- **Memory:** No significant change
- **With --exclude-private-links:** Faster for large topologies (skips private endpoint processing)

## Backward Compatibility

✅ **Fully Backward Compatible**
- `GenerateDOTFile()` function unchanged - existing code works as-is
- Default behavior includes private links (same as before)
- New features opt-in via CLI flag

## Future Enhancements

Potential improvements for consideration:
- Additional layout algorithms (dot, neato, fdp, circo, twopi)
- Custom legend content via configuration
- Adjustable spacing parameters via flags
- Layout presets (compact, standard, spacious)

## Success Criteria

All requirements from `instructions/alignment-issues.md` met:

✅ Graph positioned closer to left side
✅ Reduced excessive white space
✅ Legend integrated with visualization (bottom-left)
✅ Private links table well-positioned (bottom-center)
✅ Professional, compact layout
✅ CLI flag to exclude private links
✅ Tested across small, medium, and large topologies
✅ Suitable for presentations and reports
