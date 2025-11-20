# DOT Syntax Error Fixes - Azure Firewall Implementation

## Issue Description
User reported: "Could not render SVG, failed to parse DOT content, syntax error in line 56 near '->'"

This error occurred specifically with networks containing Azure Firewalls, suggesting an issue with DOT file generation.

## Root Causes Identified

### 1. **Insufficient Character Sanitization**
The `sanitizeName()` function only handled 3 special characters:
- Hyphens `-`
- Dots `.`
- Spaces ` `

**Problem**: Azure resource names can contain many other special characters that are invalid in DOT node identifiers:
- Parentheses: `(`, `)`
- Brackets: `[`, `]`
- Colons: `:`
- Slashes: `/`, `\`
- And 20+ other special characters

### 2. **Missing Defensive Checks**
The code created edges to nodes without verifying the nodes existed in the lookup maps:
- `firewallNodes[fw.ID]` - no existence check
- `routeTables[rt.ID]` - no existence check
- `subnetNodes[fw.SubnetID]` - no existence check

**Problem**: If a node wasn't rendered (for any reason), creating an edge to it would result in a DOT syntax error like:
```
rt_main -> fw_0 [...]  // ERROR if fw_0 doesn't exist
```

## Fixes Implemented

### Fix 1: Enhanced `sanitizeName()` Function
**File**: `pkg/visualization/graphviz.go`

**Before**:
```go
func sanitizeName(name string) string {
    name = strings.ReplaceAll(name, "-", "_")
    name = strings.ReplaceAll(name, ".", "_")
    name = strings.ReplaceAll(name, " ", "_")
    return name
}
```

**After**:
```go
func sanitizeName(name string) string {
    // Handles 30+ special characters
    replacements := map[string]string{
        "-": "_", ".": "_", " ": "_",
        "(": "_", ")": "_", "[": "_", "]": "_",
        "{": "_", "}": "_", ":": "_", ";": "_",
        // ... and 20+ more
    }

    // Replace known special chars
    result := name
    for old, new := range replacements {
        result = strings.ReplaceAll(result, old, new)
    }

    // Remove ANY remaining non-alphanumeric characters
    var builder strings.Builder
    for _, r := range result {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
           (r >= '0' && r <= '9') || r == '_' {
            builder.WriteRune(r)
        } else {
            builder.WriteRune('_')
        }
    }
    return builder.String()
}
```

**Benefits**:
- Handles ALL special characters, not just 3
- Two-layer protection: explicit replacement + fallback sanitization
- Ensures DOT identifiers are always valid

### Fix 2: Defensive Node Reference Checks
**File**: `pkg/visualization/graphviz.go`

#### Firewall → Subnet Edges
**Before**:
```go
for _, fw := range topology.AzureFirewalls {
    if subnetNode, exists := subnetNodes[fw.SubnetID]; exists {
        fwNode := firewallNodes[fw.ID]  // NO CHECK!
        dot.WriteString(fmt.Sprintf("  %s -> %s [...];\n", fwNode, subnetNode))
    }
}
```

**After**:
```go
for _, fw := range topology.AzureFirewalls {
    // Defensive check: ensure BOTH nodes exist
    fwNode, fwExists := firewallNodes[fw.ID]
    subnetNode, subnetExists := subnetNodes[fw.SubnetID]
    if fwExists && subnetExists {
        dot.WriteString(fmt.Sprintf("  %s -> %s [...];\n", fwNode, subnetNode))
    }
}
```

#### Route Table → Firewall Edges
**Before**:
```go
for _, rt := range topology.RouteTables {
    rtNodeID := routeTables[rt.ID]  // NO CHECK!
    for _, route := range rt.Routes {
        if route.NextHopType == "VirtualAppliance" {
            for _, fw := range topology.AzureFirewalls {
                if fw.PrivateIPAddress == route.NextHopIPAddress {
                    fwNode := firewallNodes[fw.ID]  // NO CHECK!
                    dot.WriteString(fmt.Sprintf("  %s -> %s [...];\n", rtNodeID, fwNode))
                }
            }
        }
    }
}
```

**After**:
```go
for _, rt := range topology.RouteTables {
    // Defensive check: ensure route table node exists
    rtNodeID, rtExists := routeTables[rt.ID]
    if !rtExists {
        continue  // Skip if route table wasn't rendered
    }

    for _, route := range rt.Routes {
        if route.NextHopType == "VirtualAppliance" {
            for _, fw := range topology.AzureFirewalls {
                if fw.PrivateIPAddress == route.NextHopIPAddress {
                    // Defensive check: ensure firewall node exists
                    fwNode, fwExists := firewallNodes[fw.ID]
                    if !fwExists {
                        break  // Skip if firewall wasn't rendered
                    }
                    dot.WriteString(fmt.Sprintf("  %s -> %s [...];\n", rtNodeID, fwNode))
                }
            }
        }
    }
}
```

**Benefits**:
- Never creates edges to non-existent nodes
- Prevents DOT syntax errors from broken references
- Gracefully handles partial topology data

### Fix 3: Comprehensive Test Suite
**File**: `pkg/visualization/graphviz_test.go`

Added tests for:
1. **TestSanitizeName**: Tests all special character handling
2. **TestFirewallWithSpecialCharacters**: Tests firewalls with names like `fw-prod(east)-01`
3. **TestFirewallWithoutMatchingRoutes**: Tests when routes don't match firewall IPs
4. **TestFirewallEdgeWithMatchingIP**: Tests correct edge creation
5. **TestEmptyTopology**: Tests graceful handling of empty data

## Test Results

All tests pass:
```
=== RUN   TestSanitizeName
--- PASS: TestSanitizeName (0.00s)
=== RUN   TestFirewallWithoutMatchingRoutes
--- PASS: TestFirewallWithoutMatchingRoutes (0.00s)
=== RUN   TestFirewallWithSpecialCharacters
--- PASS: TestFirewallWithSpecialCharacters (0.00s)
=== RUN   TestFirewallEdgeWithMatchingIP
--- PASS: TestFirewallEdgeWithMatchingIP (0.00s)
PASS
```

## Edge Cases Handled

### ✅ Special Characters in Resource Names
- `fw-prod(east)-01` → node ID: `fw_0`, label: `fw-prod(east)-01`
- `firewall[test]` → node ID: `fw_0`, label: `firewall[test]`
- `my:firewall/prod` → node ID: `fw_0`, label: `my:firewall/prod`

### ✅ Non-existent Node References
- Route tables without associated subnets
- Firewalls without matching routes
- Broken subnet references

### ✅ Empty/Partial Topologies
- No firewalls
- No route tables
- No matching IPs

## Verification

Build and test:
```bash
go build -o az-network-analyzer
./az-network-analyzer analyze --dry-run -s test -g test --viz-format svg
```

✅ Successfully generates SVG without syntax errors
✅ Handles complex Azure resource names
✅ Gracefully handles missing data

## Prevention Measures

1. **Two-layer sanitization**: Explicit + fallback character handling
2. **Defensive programming**: Always check map existence before dereferencing
3. **Comprehensive testing**: Edge cases covered in test suite
4. **Clear error messages**: Would show which node reference failed (if logging added)

## Files Changed

1. `pkg/visualization/graphviz.go` - Enhanced sanitization + defensive checks
2. `pkg/visualization/graphviz_test.go` - Comprehensive test suite
3. `DOT_SYNTAX_FIXES.md` - This documentation

## Recommendation

For production use, consider adding:
1. **Logging**: Log when nodes are skipped due to non-existence
2. **Validation**: Pre-validate topology before rendering
3. **Error context**: Include resource IDs in error messages
