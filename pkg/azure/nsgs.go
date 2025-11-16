package azure

import (
	"context"
	"fmt"

	"azure-network-analyzer/pkg/models"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// GetNetworkSecurityGroups retrieves all NSGs in the specified resource group
func (c *AzureClient) GetNetworkSecurityGroups(ctx context.Context, resourceGroup string) ([]models.NetworkSecurityGroup, error) {
	client, err := c.getNSGsClient()
	if err != nil {
		return nil, err
	}

	var nsgs []models.NetworkSecurityGroup
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of NSGs: %w", err)
		}

		for _, nsg := range page.Value {
			n := models.NetworkSecurityGroup{
				ID:            safeString(nsg.ID),
				Name:          safeString(nsg.Name),
				ResourceGroup: resourceGroup,
				Location:      safeString(nsg.Location),
				SecurityRules: []models.SecurityRule{},
				Associations: models.NSGAssociations{
					Subnets:           []string{},
					NetworkInterfaces: []string{},
				},
			}

			if nsg.Properties != nil {
				// Extract security rules (both default and custom)
				for _, rule := range nsg.Properties.SecurityRules {
					r := extractSecurityRule(rule)
					n.SecurityRules = append(n.SecurityRules, r)
				}

				// Also include default security rules
				for _, rule := range nsg.Properties.DefaultSecurityRules {
					r := extractSecurityRule(rule)
					n.SecurityRules = append(n.SecurityRules, r)
				}

				// Extract subnet associations
				for _, subnet := range nsg.Properties.Subnets {
					if subnet.ID != nil {
						n.Associations.Subnets = append(n.Associations.Subnets, *subnet.ID)
					}
				}

				// Extract NIC associations
				for _, nic := range nsg.Properties.NetworkInterfaces {
					if nic.ID != nil {
						n.Associations.NetworkInterfaces = append(n.Associations.NetworkInterfaces, *nic.ID)
					}
				}
			}

			nsgs = append(nsgs, n)
		}
	}

	return nsgs, nil
}

// GetNSGSecurityRules retrieves all security rules (custom + default) for a specific NSG
func (c *AzureClient) GetNSGSecurityRules(ctx context.Context, resourceGroup, nsgName string) ([]models.SecurityRule, error) {
	client, err := c.getNSGsClient()
	if err != nil {
		return nil, err
	}

	nsg, err := client.Get(ctx, resourceGroup, nsgName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get NSG %s: %w", nsgName, err)
	}

	var rules []models.SecurityRule

	if nsg.Properties != nil {
		// Custom rules
		for _, rule := range nsg.Properties.SecurityRules {
			r := extractSecurityRule(rule)
			rules = append(rules, r)
		}

		// Default rules
		for _, rule := range nsg.Properties.DefaultSecurityRules {
			r := extractSecurityRule(rule)
			rules = append(rules, r)
		}
	}

	return rules, nil
}

// GetNSGAssociations retrieves the associations (subnets and NICs) for a specific NSG
func (c *AzureClient) GetNSGAssociations(ctx context.Context, resourceGroup, nsgName string) (*models.NSGAssociations, error) {
	client, err := c.getNSGsClient()
	if err != nil {
		return nil, err
	}

	nsg, err := client.Get(ctx, resourceGroup, nsgName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get NSG %s: %w", nsgName, err)
	}

	associations := &models.NSGAssociations{
		Subnets:           []string{},
		NetworkInterfaces: []string{},
	}

	if nsg.Properties != nil {
		for _, subnet := range nsg.Properties.Subnets {
			if subnet.ID != nil {
				associations.Subnets = append(associations.Subnets, *subnet.ID)
			}
		}

		for _, nic := range nsg.Properties.NetworkInterfaces {
			if nic.ID != nil {
				associations.NetworkInterfaces = append(associations.NetworkInterfaces, *nic.ID)
			}
		}
	}

	return associations, nil
}

// extractSecurityRule extracts security rule details from an Azure security rule
func extractSecurityRule(rule *armnetwork.SecurityRule) models.SecurityRule {
	r := models.SecurityRule{
		Name:        safeString(rule.Name),
		Priority:    0,
		Direction:   "",
		Access:      "",
		Protocol:    "",
		Description: "",
	}

	if rule.Properties != nil {
		if rule.Properties.Priority != nil {
			r.Priority = *rule.Properties.Priority
		}

		if rule.Properties.Direction != nil {
			r.Direction = string(*rule.Properties.Direction)
		}

		if rule.Properties.Access != nil {
			r.Access = string(*rule.Properties.Access)
		}

		if rule.Properties.Protocol != nil {
			r.Protocol = string(*rule.Properties.Protocol)
		}

		if rule.Properties.SourceAddressPrefix != nil {
			r.SourceAddressPrefix = *rule.Properties.SourceAddressPrefix
		}

		if rule.Properties.SourcePortRange != nil {
			r.SourcePortRange = *rule.Properties.SourcePortRange
		}

		if rule.Properties.DestinationAddressPrefix != nil {
			r.DestinationAddressPrefix = *rule.Properties.DestinationAddressPrefix
		}

		if rule.Properties.DestinationPortRange != nil {
			r.DestinationPortRange = *rule.Properties.DestinationPortRange
		}

		if rule.Properties.Description != nil {
			r.Description = *rule.Properties.Description
		}
	}

	return r
}
