# CLAUDE.md - Project Guidelines for Claude Code

## Project Overview

Azure Network Topology Analyzer - A CLI tool for analyzing Azure network infrastructure, identifying security risks, and generating topology reports with visualizations.

## Code Architecture

### Package Structure
- `cmd/` - CLI commands using Cobra framework
- `pkg/models/` - Data structures for Azure resources
- `pkg/azure/` - Azure SDK integration and API calls
- `pkg/analyzer/` - Security and topology analysis logic
- `pkg/reporter/` - Report generation (JSON, Markdown, HTML)
- `pkg/visualization/` - Graphviz diagram generation

### Key Patterns
- Lazy client initialization in `pkg/azure/client.go`
- Safe pointer dereferencing with `safeString()` helper
- Extractor methods convert Azure SDK types to our models
- Mock client for testing without Azure connectivity

## Code Style Guidelines

### Go Best Practices
- Follow Go idioms and effective Go guidelines
- Use meaningful variable names (no single letters except loops)
- Keep functions focused and small (<50 lines preferred)
- Return errors, don't panic
- Use table-driven tests

### Error Handling
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Check for nil pointers before dereferencing
- Log warnings for non-critical failures
- Return early on errors

### Azure SDK Patterns
- Always check if Azure SDK pointers are nil before accessing
- Use `safeString()` for string pointer conversion
- Extract resource names from Azure resource IDs
- Handle partial data gracefully (some fields may be nil)

## Security Requirements

### Code Security
- Never hardcode credentials or API keys
- Validate all Azure API responses
- Handle nil pointers to prevent panics
- No SQL/command injection vulnerabilities
- Sanitize user inputs

### Security Analysis Rules
- Flag any rule allowing 0.0.0.0/0 as source
- Critical: SSH (22), RDP (3389) exposed to internet
- Critical: Database ports exposed (1433, 3306, 5432, 27017)
- High: Subnets without NSG protection
- Medium: Wide port ranges (*, 0-65535)
- Low: Missing descriptions on security rules

## Testing Requirements

### Unit Tests
- Test all helper functions
- Test edge cases (nil inputs, empty collections, boundary values)
- Test security analysis logic with known vulnerable configs
- Mock Azure API responses for testing

### Test Patterns
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
    }{
        {"test case 1", input1, expected1},
        {"test case 2", input2, expected2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Adding New Features

### New Azure Resource Type
1. Add model struct in `pkg/models/topology.go`
2. Add field to `NetworkTopology` struct
3. Create extraction method in appropriate `pkg/azure/*.go` file
4. Add collection logic in `cmd/analyze.go`
5. Update analysis in `pkg/analyzer/`
6. Add to reporters in `pkg/reporter/`
7. Add to visualization in `pkg/visualization/graphviz.go`
8. Add tests for new functionality

### New Security Check
1. Add detection logic in `pkg/analyzer/security.go`
2. Use appropriate severity level (Critical/High/Medium/Low/Info)
3. Provide actionable recommendation
4. Add test case with vulnerable configuration
5. Update documentation

## Documentation

- Update README.md for user-facing changes
- Add comments for complex logic
- Document all public functions
- Keep CLAUDE.md updated with new patterns

## Build & Release

- Binaries built for: Linux (amd64/arm64), Windows (amd64), macOS (amd64/arm64)
- Docker image uses multi-stage Alpine build
- CI runs tests on all pushes to main
- Use semantic versioning for releases

## Common Commands

```bash
# Build
go build -o az-network-analyzer

# Test
go test ./... -v

# Format
go fmt ./...

# Lint
go vet ./...

# Dry run
./az-network-analyzer analyze --dry-run -s test -g test
```

## Important Files

- `main.go` - Entry point
- `cmd/analyze.go` - Main analysis orchestration
- `pkg/azure/client.go` - Azure client and helpers
- `pkg/analyzer/security.go` - Security risk detection
- `pkg/models/topology.go` - All data models
- `.github/workflows/ci.yml` - CI/CD pipeline

## Pull Request Guidelines

1. Include tests for new functionality
2. Update documentation as needed
3. Follow existing code patterns
4. Ensure all tests pass
5. Check for security implications
6. Keep PRs focused and small
