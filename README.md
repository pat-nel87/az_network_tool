# Azure Network Topology Analyzer

A CLI tool for analyzing Azure network infrastructure, identifying security risks, and generating comprehensive topology reports.

## Mission

Provide cloud engineers and security teams with deep visibility into their Azure network architecture by:

- **Discovering** all network resources (VNets, NSGs, Private Endpoints, Route Tables, Gateways, Load Balancers)
- **Analyzing** security configurations to identify risks and misconfigurations
- **Visualizing** network topology and relationships between resources
- **Reporting** findings in actionable formats (JSON, Markdown, HTML)

## Features

### Current (Phase 1-2)
- [x] Core data models for all Azure network resource types
- [x] Azure SDK integration with lazy client initialization
- [x] CLI framework with Cobra
- [x] Mock/dry-run mode for testing without Azure connectivity
- [x] Comprehensive test suite with edge case coverage

### Planned
- [ ] Security analysis (NSG rule risks, open ports, overly permissive rules)
- [ ] Topology analysis (isolated subnets, peering issues, routing anomalies)
- [ ] Network visualization with Graphviz
- [ ] Multi-format reporting (JSON, Markdown, HTML)
- [ ] Network Watcher integration (flow logs, connection monitors)

## Quick Start

```bash
# Build
go build -o az-network-analyzer

# Dry-run with mock data (no Azure required)
./az-network-analyzer analyze --dry-run

# Analyze real Azure resources
./az-network-analyzer analyze \
  --subscription "your-subscription-id" \
  --resource-group "your-resource-group"
```

## Requirements

- Go 1.21+
- Azure CLI (for authentication)
- Azure subscription with Reader access

## Project Structure

```
.
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go
│   └── analyze.go
├── pkg/
│   ├── models/            # Data structures for all resource types
│   │   └── topology.go
│   └── azure/             # Azure SDK client and operations
│       ├── client.go      # Core client, helpers, extractors
│       ├── vnets.go       # Virtual Network operations
│       ├── nsgs.go        # Network Security Group operations
│       ├── privatelink.go # Private Endpoint/DNS operations
│       ├── routing.go     # Route Table/NAT Gateway operations
│       ├── gateways.go    # VPN/ExpressRoute operations
│       ├── loadbalancers.go # Load Balancer/App Gateway operations
│       └── mock_client.go # Mock data for testing
├── main.go
└── go.mod
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Check coverage
go test ./... -cover
```

## Authentication

Uses Azure DefaultAzureCredential which supports:
- Azure CLI (`az login`)
- Environment variables
- Managed Identity
- Visual Studio Code credentials

## License

MIT

## Contributing

This project is under active development. See the instructions folder for the complete implementation roadmap.
