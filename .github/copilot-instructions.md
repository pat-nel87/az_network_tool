# GitHub Copilot Instructions for Azure Network Topology Analyzer

## Project Overview

This is an Azure Network Topology Analyzer - a comprehensive CLI tool built in Go for analyzing Azure network infrastructure, identifying security risks, and generating topology reports with visualizations.

**Key Technologies:**
- Language: Go 1.23+
- CLI Framework: Cobra
- Azure SDK: github.com/Azure/azure-sdk-for-go
- Visualization: github.com/goccy/go-graphviz
- Output Formats: JSON, Markdown, HTML, Graphviz (DOT/SVG/PNG)

## Code Architecture

### Package Structure
```
cmd/              - CLI commands using Cobra framework
pkg/models/       - Data structures for Azure resources
pkg/azure/        - Azure SDK integration and API calls
pkg/analyzer/     - Security and topology analysis logic
pkg/reporter/     - Report generation (JSON, Markdown, HTML)
pkg/visualization/ - Graphviz diagram generation
```

### Key Design Patterns
- **Lazy client initialization** in `pkg/azure/client.go`
- **Safe pointer dereferencing** using `safeString()` helper function
- **Extractor methods** to convert Azure SDK types to our internal models
- **Mock client** pattern for testing without Azure connectivity

## Coding Standards

### Go Best Practices
- Follow standard Go idioms and effective Go guidelines
- Use meaningful variable names (avoid single letters except in loops)
- Keep functions focused and small (prefer <50 lines)
- Return errors, don't panic
- Use table-driven tests

### Error Handling
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Always check for nil pointers before dereferencing
- Log warnings for non-critical failures
- Return early on errors

### Azure SDK Patterns
**CRITICAL:** Azure SDK returns pointers extensively
- Always check if Azure SDK pointers are nil before accessing
- Use `safeString()` helper for string pointer conversion
- Extract resource names from Azure resource IDs properly
- Handle partial data gracefully (some fields may be nil)

Example:
```go
func safeString(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}
```

## Security Requirements

### Code Security
- Never hardcode credentials or API keys
- Validate all Azure API responses
- Handle nil pointers to prevent panics
- Prevent SQL/command injection vulnerabilities
- Sanitize all user inputs

### Security Analysis Rules
When implementing or modifying security analysis:
- **Critical**: Flag SSH (22), RDP (3389) exposed to internet (0.0.0.0/0)
- **Critical**: Flag database ports exposed (1433, 3306, 5432, 27017)
- **High**: Flag subnets without NSG protection
- **Medium**: Flag wide port ranges (*, 0-65535)
- **Low**: Flag security rules missing descriptions

## Testing Requirements

### Unit Tests
- Test all helper functions with edge cases
- Test nil inputs, empty collections, boundary values
- Test security analysis with known vulnerable configurations
- Mock Azure API responses for testing

### Test Pattern
Use table-driven tests consistently:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
    }{
        {"descriptive test case name", input1, expected1},
        {"another test case", input2, expected2},
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

### Adding a New Azure Resource Type
1. Add model struct in `pkg/models/topology.go`
2. Add field to `NetworkTopology` struct
3. Create extraction method in appropriate `pkg/azure/*.go` file
4. Add collection logic in `cmd/analyze.go`
5. Update analysis in `pkg/analyzer/`
6. Add to reporters in `pkg/reporter/`
7. Add to visualization in `pkg/visualization/graphviz.go`
8. Add tests for new functionality

### Adding a New Security Check
1. Add detection logic in `pkg/analyzer/security.go`
2. Use appropriate severity level (Critical/High/Medium/Low/Info)
3. Provide actionable recommendation
4. Add test case with vulnerable configuration
5. Update documentation

## Build & Development

### Common Commands
```bash
# Build
go build -o az-network-analyzer

# Test (run before commits)
go test ./... -v

# Format (run before commits)
go fmt ./...

# Lint
go vet ./...

# Dry run for testing without Azure
./az-network-analyzer analyze --dry-run -s test -g test
```

### CI/CD
- Binaries built for: Linux (amd64/arm64), Windows (amd64), macOS (amd64/arm64)
- Docker image uses multi-stage Alpine build
- CI runs tests on all pushes to main
- Use semantic versioning for releases

## Important Files to Know

- `main.go` - Entry point
- `cmd/analyze.go` - Main analysis orchestration
- `pkg/azure/client.go` - Azure client and helpers
- `pkg/analyzer/security.go` - Security risk detection
- `pkg/models/topology.go` - All data models
- `.github/workflows/ci.yml` - CI/CD pipeline

## Pull Request Guidelines

When creating or reviewing PRs:
1. Include tests for new functionality
2. Update documentation for user-facing changes
3. Follow existing code patterns and structure
4. Ensure all tests pass (`go test ./...`)
5. Check for security implications
6. Keep PRs focused and small
7. Use conventional commit messages

## Code Comments

- Add comments for complex logic that isn't self-explanatory
- Document all exported (public) functions
- Keep comments concise and up-to-date
- Use godoc format for package and function documentation

## Specific Implementation Notes

### Azure SDK Pagination
Many Azure SDK list operations return pagers - always iterate through all pages:
```go
pager := client.NewListPager(...)
for pager.More() {
    page, err := pager.NextPage(ctx)
    // handle error and process page
}
```

### Resource Dependencies
Some resources need parent resource info (e.g., subnets need VNet):
- Build dependency graph when needed
- Extract resource information from Azure resource IDs
- Handle cross-resource-group references

### Performance Considerations
- Use goroutines for parallel resource collection where safe
- Be mindful of Azure API rate limits
- Implement exponential backoff for retries
- Handle partial failures gracefully

## Documentation Standards

- Update README.md for user-facing changes
- Update CLAUDE.md for development patterns
- Add inline godoc comments for exported functions
- Include examples in documentation where helpful

## What NOT to Do

❌ Don't hardcode Azure credentials
❌ Don't panic - return errors instead
❌ Don't ignore nil pointer checks with Azure SDK
❌ Don't commit without running tests
❌ Don't remove working code without good reason
❌ Don't add dependencies without consideration
❌ Don't log or output sensitive data (secrets, keys, connection strings)

## Questions or Issues?

- Check existing issues and PRs first
- Review the README.md for usage examples
- Check CLAUDE.md for detailed implementation guidance
- Consult Azure SDK documentation for SDK-specific questions
