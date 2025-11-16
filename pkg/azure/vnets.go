package azure

import (
	"context"
	"fmt"

	"azure-network-analyzer/pkg/models"
)

// GetVirtualNetworks retrieves all virtual networks in the specified resource group
func (c *AzureClient) GetVirtualNetworks(ctx context.Context, resourceGroup string) ([]models.VirtualNetwork, error) {
	client, err := c.getVNetsClient()
	if err != nil {
		return nil, err
	}

	var vnets []models.VirtualNetwork
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of VNets: %w", err)
		}

		for _, vnet := range page.Value {
			v := models.VirtualNetwork{
				ID:            safeString(vnet.ID),
				Name:          safeString(vnet.Name),
				ResourceGroup: resourceGroup,
				Location:      safeString(vnet.Location),
				AddressSpace:  []string{},
				Subnets:       []models.Subnet{},
				Peerings:      []models.VNetPeering{},
				DNSServers:    []string{},
				EnableDDoS:    false,
			}

			// Extract address space
			if vnet.Properties != nil && vnet.Properties.AddressSpace != nil {
				for _, addr := range vnet.Properties.AddressSpace.AddressPrefixes {
					if addr != nil {
						v.AddressSpace = append(v.AddressSpace, *addr)
					}
				}
			}

			// Extract DNS servers
			if vnet.Properties != nil && vnet.Properties.DhcpOptions != nil {
				for _, dns := range vnet.Properties.DhcpOptions.DNSServers {
					if dns != nil {
						v.DNSServers = append(v.DNSServers, *dns)
					}
				}
			}

			// Extract DDoS protection status
			if vnet.Properties != nil && vnet.Properties.EnableDdosProtection != nil {
				v.EnableDDoS = *vnet.Properties.EnableDdosProtection
			}

			// Extract subnets
			if vnet.Properties != nil && vnet.Properties.Subnets != nil {
				for _, subnet := range vnet.Properties.Subnets {
					s := c.extractSubnet(subnet)
					v.Subnets = append(v.Subnets, s)
				}
			}

			// Extract VNet peerings
			if vnet.Properties != nil && vnet.Properties.VirtualNetworkPeerings != nil {
				for _, peering := range vnet.Properties.VirtualNetworkPeerings {
					p := c.extractVNetPeering(peering)
					v.Peerings = append(v.Peerings, p)
				}
			}

			vnets = append(vnets, v)
		}
	}

	return vnets, nil
}

// GetSubnets retrieves all subnets for a specific VNet
func (c *AzureClient) GetSubnets(ctx context.Context, resourceGroup, vnetName string) ([]models.Subnet, error) {
	client, err := c.getSubnetsClient()
	if err != nil {
		return nil, err
	}

	var subnets []models.Subnet
	pager := client.NewListPager(resourceGroup, vnetName, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of subnets: %w", err)
		}

		for _, subnet := range page.Value {
			s := c.extractSubnet(subnet)
			subnets = append(subnets, s)
		}
	}

	return subnets, nil
}

// GetVNetPeerings retrieves all peerings for a specific VNet
func (c *AzureClient) GetVNetPeerings(ctx context.Context, resourceGroup, vnetName string) ([]models.VNetPeering, error) {
	client, err := c.getPeeringsClient()
	if err != nil {
		return nil, err
	}

	var peerings []models.VNetPeering
	pager := client.NewListPager(resourceGroup, vnetName, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of peerings: %w", err)
		}

		for _, peering := range page.Value {
			p := c.extractVNetPeering(peering)
			peerings = append(peerings, p)
		}
	}

	return peerings, nil
}
