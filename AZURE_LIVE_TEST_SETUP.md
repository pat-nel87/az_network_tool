# Azure Live Testing Setup Instructions

## For Claude Code: Execute these steps to set up Azure test resources

### Prerequisites
- Azure CLI installed (`brew install azure-cli` on macOS)
- Logged into Azure (`az login`)

---

## Step 1: Verify Azure Login
```bash
az account show --query '{name:name, id:id}' --output json
```
Save the subscription ID from the output.

---

## Step 2: Create Test Resource Group
```bash
az group create -n test-network-rg -l eastus
```

---

## Step 3: Create Virtual Network with Subnets
```bash
# Create VNet with first subnet
az network vnet create -g test-network-rg -n test-vnet \
  --address-prefix 10.0.0.0/16 \
  --subnet-name web-subnet --subnet-prefix 10.0.1.0/24

# Add second subnet
az network vnet subnet create -g test-network-rg \
  --vnet-name test-vnet -n db-subnet --address-prefix 10.0.2.0/24
```

---

## Step 4: Create Network Security Group with Rule
```bash
# Create NSG
az network nsg create -g test-network-rg -n test-nsg

# Add SSH allow rule (for testing - this is intentionally risky)
az network nsg rule create -g test-network-rg --nsg-name test-nsg \
  -n AllowSSH --priority 100 --access Allow --protocol Tcp \
  --direction Inbound --source-address-prefixes '*' \
  --destination-port-ranges 22

# Associate NSG to web-subnet
az network vnet subnet update -g test-network-rg \
  --vnet-name test-vnet -n web-subnet --network-security-group test-nsg
```

---

## Step 5: Create Route Table with Route
```bash
# Create route table
az network route-table create -g test-network-rg -n test-rt

# Add route to internet
az network route-table route create -g test-network-rg \
  --route-table-name test-rt -n to-internet \
  --address-prefix 0.0.0.0/0 --next-hop-type Internet

# Associate route table to db-subnet
az network vnet subnet update -g test-network-rg \
  --vnet-name test-vnet -n db-subnet --route-table test-rt
```

---

## Step 6: Verify Resources Created
```bash
az network vnet show -g test-network-rg -n test-vnet --query '{name:name, addressSpace:addressSpace.addressPrefixes}' -o json
az network nsg show -g test-network-rg -n test-nsg --query '{name:name, rulesCount:length(securityRules)}' -o json
az network route-table show -g test-network-rg -n test-rt --query '{name:name, routesCount:length(routes)}' -o json
```

---

## Step 7: Test the Azure Network Analyzer
```bash
# Navigate to project directory
cd /path/to/az_network_tool

# Build the analyzer
go build -o az-network-analyzer

# Get subscription ID
SUBSCRIPTION_ID=$(az account show --query id --output tsv)

# Run the analyzer
./az-network-analyzer analyze \
  --subscription "$SUBSCRIPTION_ID" \
  --resource-group "test-network-rg"
```

---

## Expected Results
The analyzer should discover:
- **1 Virtual Network** (test-vnet) with address space 10.0.0.0/16
  - **2 Subnets**: web-subnet (10.0.1.0/24) and db-subnet (10.0.2.0/24)
- **1 Network Security Group** (test-nsg) with 1 custom rule
  - Associated to web-subnet
  - Contains risky SSH rule allowing all sources
- **1 Route Table** (test-rt) with 1 route
  - Associated to db-subnet
  - Route to 0.0.0.0/0 via Internet

---

## Cleanup (Important!)
Delete all test resources when done to avoid any potential charges:
```bash
az group delete -n test-network-rg --yes --no-wait
```

---

## Troubleshooting

### Authentication Issues
If `az login` fails, try:
```bash
az login --use-device-code
```

### Permission Issues
Ensure your Azure account has at least "Reader" role on the subscription.

### Module Issues
If Go compilation fails:
```bash
go mod tidy
go mod download
```

---

## Cost Information
All resources created are **FREE**:
- Virtual Networks: No charge
- Subnets: No charge
- Network Security Groups: No charge
- Route Tables: No charge

Only resources that cost money (NOT created here):
- VPN Gateways (~$140/month)
- Application Gateways (~$25/month)
- Load Balancers (~$18/month)
- Private Endpoints (~$7/month)
