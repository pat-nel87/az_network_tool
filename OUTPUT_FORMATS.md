# Output Formats Guide

Azure Network Topology Analyzer supports multiple output formats for both reports and visualizations.

## Report Formats

Use the `--output-format` or `-o` flag to specify the report format:

### Markdown (default)
```bash
az-network-analyzer analyze -s <subscription-id> -g <resource-group> -o markdown
```
- **Use Case**: Human-readable documentation, README files, wikis
- **Features**: Tables, bullet lists, headers, easy to read in plain text
- **File Extension**: `.md`

### JSON
```bash
az-network-analyzer analyze -s <subscription-id> -g <resource-group> -o json
```
- **Use Case**: Machine-readable output, automation, further processing
- **Features**: Complete structured data, easy to parse programmatically
- **File Extension**: `.json`

### HTML
```bash
az-network-analyzer analyze -s <subscription-id> -g <resource-group> -o html
```
- **Use Case**: Web viewing, sharing in browsers, presentations
- **Features**: Styled output, interactive viewing, can be opened directly in browser
- **File Extension**: `.html`

## Visualization Formats

Use the `--viz-format` flag to specify the network diagram format:

### SVG (default) - Scalable Vector Graphics
```bash
az-network-analyzer analyze -s <sub-id> -g <rg> --viz-format svg
```
- **Best For**: Web display, high-quality scalable graphics
- **Pros**:
  - Infinite zoom without quality loss
  - Small file size
  - Can be edited in vector graphics tools
  - Works in all modern browsers
- **Cons**: May not work in older applications
- **File Extension**: `.svg`

### PNG - Portable Network Graphics
```bash
az-network-analyzer analyze -s <sub-id> -g <rg> --viz-format png
```
- **Best For**: Screenshots, embedding in documents, presentations
- **Pros**:
  - Universal compatibility
  - Works in all image viewers
  - Good for PowerPoint, Word, etc.
- **Cons**: Fixed resolution, larger file size
- **File Extension**: `.png`

### JPEG/JPG - Joint Photographic Experts Group
```bash
az-network-analyzer analyze -s <sub-id> -g <rg> --viz-format jpg
```
- **Best For**: Email attachments, web sharing with smaller file sizes
- **Pros**:
  - Smaller file size than PNG
  - Universal compatibility
  - Good for sharing via email
- **Cons**: Lossy compression, not ideal for diagrams with text
- **File Extension**: `.jpg`

### PDF - Portable Document Format
```bash
az-network-analyzer analyze -s <sub-id> -g <rg> --viz-format pdf
```
- **Best For**: Printing, formal documentation, archiving
- **Pros**:
  - Professional format
  - Good for printing
  - Preserves layout
- **Cons**: **Limited support** - May not work in all environments due to WASM limitations
- **File Extension**: `.pdf`
- **Note**: PDF rendering is experimental and may fail in some environments

### DOT - GraphViz Source Format
```bash
az-network-analyzer analyze -s <sub-id> -g <rg> --viz-format dot
```
- **Best For**: Custom rendering, further editing, version control
- **Pros**:
  - Plain text format
  - Can be edited manually
  - Can be rendered with external GraphViz tools
  - Version control friendly
  - Can generate any GraphViz-supported format externally
- **Cons**: Requires external tools to view
- **File Extension**: `.dot`

## Using DOT Files with External GraphViz Tools

DOT files are the source format that can be rendered into any GraphViz-supported format:

### Install GraphViz
```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Windows
choco install graphviz
```

### Render DOT to Various Formats
```bash
# Generate the DOT file
az-network-analyzer analyze -s <sub-id> -g <rg> --viz-format dot

# Render to different formats using dot command
dot -Tsvg network-topology.dot -o output.svg
dot -Tpng network-topology.dot -o output.png
dot -Tpdf network-topology.dot -o output.pdf
dot -Tjpg network-topology.dot -o output.jpg
dot -Tgif network-topology.dot -o output.gif
dot -Tps network-topology.dot -o output.ps
```

### Additional GraphViz Formats

The DOT format can be rendered to many other formats using the `dot` command:
- **eps** - Encapsulated PostScript
- **gif** - Graphics Interchange Format
- **bmp** - Windows Bitmap
- **ps** - PostScript
- **ps2** - PostScript for PDF
- **cmapx** - Client-side imagemap
- **imap** - Server-side imagemap
- **wbmp** - Wireless Bitmap

## Complete Examples

### Generate HTML report with PNG diagram
```bash
az-network-analyzer analyze \
  -s my-subscription \
  -g my-resource-group \
  --output-format html \
  --viz-format png \
  --output ./reports/network-report.html
```

### Generate JSON report with SVG diagram for CI/CD
```bash
az-network-analyzer analyze \
  -s $SUBSCRIPTION_ID \
  -g $RESOURCE_GROUP \
  --output-format json \
  --viz-format svg \
  --output ./artifacts/network-topology.json
```

### Generate DOT file for custom rendering
```bash
# Generate DOT file
az-network-analyzer analyze \
  -s my-subscription \
  -g my-resource-group \
  --viz-format dot \
  --output ./topology/report.md

# Render with custom options using GraphViz
dot -Tsvg -Gdpi=300 network-topology.dot -o high-res.svg
```

### Docker with multiple output formats
```bash
# Generate markdown report with SVG
docker run --rm \
  -v $(pwd)/output:/output \
  -e AZURE_TENANT_ID=$AZURE_TENANT_ID \
  -e AZURE_CLIENT_ID=$AZURE_CLIENT_ID \
  -e AZURE_CLIENT_SECRET=$AZURE_CLIENT_SECRET \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze -s $SUB_ID -g $RG -o markdown --viz-format svg -f /output/report.md

# Generate JSON report with PNG
docker run --rm \
  -v $(pwd)/output:/output \
  -e AZURE_TENANT_ID=$AZURE_TENANT_ID \
  -e AZURE_CLIENT_ID=$AZURE_CLIENT_ID \
  -e AZURE_CLIENT_SECRET=$AZURE_CLIENT_SECRET \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze -s $SUB_ID -g $RG -o json --viz-format png -f /output/topology.json
```

## Testing Formats

To test all output formats work correctly:

```bash
# Run format tests
go test ./pkg/visualization -run TestAllOutputFormats -v

# Test report generation
az-network-analyzer analyze -s test -g test --dry-run -o json
az-network-analyzer analyze -s test -g test --dry-run -o markdown
az-network-analyzer analyze -s test -g test --dry-run -o html
```

## Format Recommendations by Use Case

| Use Case | Report Format | Viz Format | Why |
|----------|--------------|------------|-----|
| Documentation | Markdown | SVG | Best for wikis, README files, scalable diagrams |
| Presentations | HTML | PNG | Professional look, works in PowerPoint |
| Automation/CI | JSON | DOT | Machine-readable, can be processed further |
| Email Sharing | Markdown | JPG | Small file sizes, universal compatibility |
| Formal Reports | HTML | PDF | Professional formatting, good for printing |
| Version Control | Markdown | DOT | Text-based, easy to diff and track changes |
| Web Dashboard | JSON | SVG | Data + scalable graphics for web display |

## Troubleshooting

### PDF Not Working
If PDF rendering fails, use DOT format and render externally:
```bash
az-network-analyzer analyze -s <sub> -g <rg> --viz-format dot
dot -Tpdf network-topology.dot -o output.pdf
```

### Large Topologies
For very large networks (100+ resources):
1. Use DOT format to avoid memory issues
2. Render with external GraphViz tools
3. Consider filtering to specific resource types
4. Use PNG/JPG instead of SVG for better performance

### File Locations
By default, files are saved to the current directory with auto-generated names:
- Report: `network-report-<resource-group>-<timestamp>.<ext>`
- Visualization: `network-topology-<resource-group>-<timestamp>.<ext>`

Specify custom paths with `-f` or `--output`:
```bash
az-network-analyzer analyze -s <sub> -g <rg> -f /path/to/report.md
# This will also save the visualization to /path/to/network-topology-<rg>-<timestamp>.svg
```
