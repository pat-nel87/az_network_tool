# Azure Network Analyzer - Docker Quick Start

## Prerequisites
- Docker installed
- Azure CLI installed and logged in (`az login`)
- Subscription ID and Resource Group name

## Pull the Image
```bash
docker pull ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest
```

## Run Analysis

### Basic Usage (reports saved locally)
```bash
# Create output directory
mkdir -p ~/network-reports

# Run analyzer
docker run \
  -v ~/.azure:/root/.azure \
  -v ~/network-reports:/output \
  --workdir /output \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze \
  -s YOUR_SUBSCRIPTION_ID \
  -g YOUR_RESOURCE_GROUP \
  -f /output/report.md
```

### View Results
```bash
ls -lh ~/network-reports/
cat ~/network-reports/report.md
```

Files generated:
- `report.md` - Security findings and topology summary
- `network-topology-*.svg` - Visual network diagram

### Different Output Formats

**JSON Report:**
```bash
docker run \
  -v ~/.azure:/root/.azure \
  -v ~/network-reports:/output \
  --workdir /output \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze \
  -s YOUR_SUBSCRIPTION_ID \
  -g YOUR_RESOURCE_GROUP \
  --output-format json \
  -f /output/report.json
```

**HTML Report:**
```bash
docker run \
  -v ~/.azure:/root/.azure \
  -v ~/network-reports:/output \
  --workdir /output \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze \
  -s YOUR_SUBSCRIPTION_ID \
  -g YOUR_RESOURCE_GROUP \
  --output-format html \
  -f /output/report.html
```

### Skip Visualization (faster)
```bash
docker run \
  -v ~/.azure:/root/.azure \
  -v ~/network-reports:/output \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze \
  -s YOUR_SUBSCRIPTION_ID \
  -g YOUR_RESOURCE_GROUP \
  --visualize=false \
  -f /output/report.md
```

### DOT File Only (for large networks)
```bash
docker run \
  -v ~/.azure:/root/.azure \
  -v ~/network-reports:/output \
  --workdir /output \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze \
  -s YOUR_SUBSCRIPTION_ID \
  -g YOUR_RESOURCE_GROUP \
  --viz-format dot \
  -f /output/report.md
```

Then render locally with native Graphviz:
```bash
dot -Tsvg ~/network-reports/*.dot -o ~/network-reports/topology.svg
```

## Using Service Principal (CI/CD)
```bash
docker run \
  -e AZURE_CLIENT_ID=your-client-id \
  -e AZURE_CLIENT_SECRET=your-secret \
  -e AZURE_TENANT_ID=your-tenant-id \
  -v ~/network-reports:/output \
  --workdir /output \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze \
  -s YOUR_SUBSCRIPTION_ID \
  -g YOUR_RESOURCE_GROUP \
  -f /output/report.md
```

## Test with Mock Data
```bash
docker run \
  ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest \
  analyze \
  --dry-run \
  -s test-sub \
  -g test-rg
```

## Troubleshooting

**Permission denied on ~/.azure:**
```bash
chmod -R 755 ~/.azure
```

**Image not found:**
```bash
docker pull ghcr.io/YOUR_USERNAME/azure-network-analyzer:latest
```

**Check Azure login:**
```bash
az account show
```
