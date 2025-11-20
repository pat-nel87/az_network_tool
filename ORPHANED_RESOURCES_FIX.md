# Orphaned Route Table Visualization Fix

## Issue

**Q:** What if there is a route table with no subnet associations, would this be rendered?

**A:** Initially, **NO**. Orphaned route tables (not associated with any subnets) were not rendered in the visualization.

## Impact

This was problematic because:

1. **Invisible Route Tables**: Route tables configured but not yet attached to subnets wouldn't appear
2. **Lost Firewall Egress Information**: Route tables pointing to firewalls as next hop wouldn't show the egress relationship
3. **Incomplete Security Analysis**: Security analysts couldn't see all configured routing

## Example Scenario

```
Route Table: "rt-firewall-egress"
├─ Route: 0.0.0.0/0 → Firewall (10.0.1.4)
└─ Associated Subnets: [] (empty - orphaned)
```

**Before Fix**: This route table was invisible ❌
**After Fix**: Route table rendered with firewall egress edge ✅

## Root Cause

Route tables were only added to the rendering map when referenced by subnets:

```go
// OLD CODE - Only in subnet iteration
if subnet.RouteTable != nil {
    rtID := *subnet.RouteTable
    if _, exists := routeTables[rtID]; !exists {
        routeTables[rtID] = fmt.Sprintf("rt_%s", sanitizeName(extractResourceName(rtID)))
    }
}
```

Result: Orphaned route tables never made it into the `routeTables` map.

## Fix Implementation

**File**: `pkg/visualization/graphviz.go`

Added code to ensure ALL route tables from the topology are included:

```go
// Ensure ALL route tables are in the map (including orphaned ones)
// This is important for route tables that route to firewalls but aren't yet attached to subnets
for _, rt := range topology.RouteTables {
    if _, exists := routeTables[rt.ID]; !exists {
        routeTables[rt.ID] = fmt.Sprintf("rt_%s", sanitizeName(extractResourceName(rt.ID)))
    }
}

// Render deduplicated Route Tables (outside clusters)
for rtID, rtNodeID := range routeTables {
    rtName := extractResourceName(rtID)
    dot.WriteString(fmt.Sprintf("  %s [label=\"Route Table\\n%s\", fillcolor=\"#DDA0DD\", shape=parallelogram];\n",
        rtNodeID, rtName))
}
```

## Test Coverage

Added comprehensive tests in `orphaned_resources_test.go`:

### TestOrphanedRouteTable
Tests that orphaned route tables are rendered:
```go
RouteTables: []models.RouteTable{
    {
        ID:   "rt-orphaned",
        Name: "rt-orphaned",
        AssociatedSubnets: []string{}, // Empty - orphaned
    },
}
```
**Result**: ✅ Route table is rendered

### TestOrphanedRouteTableWithFirewallRoute
Tests that orphaned RTs with firewall routes show egress edges:
```go
RouteTables: []models.RouteTable{
    {
        ID:   "rt-orphaned-with-fw",
        Routes: []models.Route{
            {
                AddressPrefix:    "0.0.0.0/0",
                NextHopType:      "VirtualAppliance",
                NextHopIPAddress: "10.0.1.4", // Points to firewall
            },
        },
        AssociatedSubnets: []string{}, // Orphaned
    },
}
```
**Result**: ✅ Route table rendered + firewall egress edge created

### TestOrphanedNSG
Verified NSGs are already handled correctly:
**Result**: ✅ NSGs were already rendering orphaned resources properly

## Before vs After

### Before
```
Orphaned RT "rt-firewall-egress"
  ├─ Route: 0.0.0.0/0 → FW (10.0.1.4)
  └─ Not rendered ❌

Firewall: fw-hub (10.0.1.4)
  └─ Rendered ✅

Visualization:
  ┌─────────────┐
  │   Firewall  │  (no incoming edges from RT)
  └─────────────┘
```

### After
```
Orphaned RT "rt-firewall-egress"
  ├─ Route: 0.0.0.0/0 → FW (10.0.1.4)
  └─ Rendered ✅

Firewall: fw-hub (10.0.1.4)
  └─ Rendered ✅

Visualization:
  ┌──────────────┐    ┌─────────────┐
  │  Route Table │───▶│   Firewall  │
  │ (egress via FW)   │   fw-hub    │
  └──────────────┘    └─────────────┘
```

## Benefits

✅ **Complete Visibility**: All route tables visible, even if not attached
✅ **Firewall Egress Tracking**: Shows which RTs route through firewalls
✅ **Configuration Validation**: See route tables that are configured but not yet in use
✅ **Consistent with NSGs**: All resource types handle orphans the same way

## Test Results

All tests pass:
```bash
=== RUN   TestOrphanedRouteTable
--- PASS: TestOrphanedRouteTable (0.00s)
=== RUN   TestOrphanedRouteTableWithFirewallRoute
--- PASS: TestOrphanedRouteTableWithFirewallRoute (0.00s)
=== RUN   TestOrphanedNSG
--- PASS: TestOrphanedNSG (0.00s)
PASS
```

## Files Modified

1. `pkg/visualization/graphviz.go` - Added orphaned route table handling
2. `pkg/visualization/orphaned_resources_test.go` - Comprehensive test suite
3. `ORPHANED_RESOURCES_FIX.md` - This documentation

## Use Cases

This fix is particularly important for:

1. **Hub-and-Spoke Topologies**: Central route tables configured before spoke attachment
2. **Migration Scenarios**: Route tables prepared but not yet applied
3. **Security Auditing**: Seeing all routing configuration regardless of attachment
4. **Troubleshooting**: Understanding why traffic isn't routing (RT not attached)
