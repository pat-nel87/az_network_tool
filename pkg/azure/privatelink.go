package azure

import (
	"context"
	"fmt"

	"azure-network-analyzer/pkg/models"
)

// GetPrivateEndpoints retrieves all private endpoints in the specified resource group
func (c *AzureClient) GetPrivateEndpoints(ctx context.Context, resourceGroup string) ([]models.PrivateEndpoint, error) {
	client, err := c.getPrivateEndpointsClient()
	if err != nil {
		return nil, err
	}

	var endpoints []models.PrivateEndpoint
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of Private Endpoints: %w", err)
		}

		for _, pe := range page.Value {
			endpoint := models.PrivateEndpoint{
				ID:            safeString(pe.ID),
				Name:          safeString(pe.Name),
				ResourceGroup: resourceGroup,
				Location:      safeString(pe.Location),
				SubnetID:      "",
				GroupIDs:      []string{},
			}

			if pe.Properties != nil {
				if pe.Properties.Subnet != nil && pe.Properties.Subnet.ID != nil {
					endpoint.SubnetID = *pe.Properties.Subnet.ID
				}

				// Get private IP address from network interfaces
				if len(pe.Properties.NetworkInterfaces) > 0 && pe.Properties.NetworkInterfaces[0].ID != nil {
					endpoint.PrivateIPAddress = "See NIC: " + extractResourceName(*pe.Properties.NetworkInterfaces[0].ID)
				}

				// Extract private link service connections
				for _, conn := range pe.Properties.PrivateLinkServiceConnections {
					if conn.Properties != nil {
						if conn.Properties.PrivateLinkServiceID != nil {
							endpoint.PrivateLinkServiceID = *conn.Properties.PrivateLinkServiceID
						}

						if conn.Properties.PrivateLinkServiceConnectionState != nil &&
							conn.Properties.PrivateLinkServiceConnectionState.Status != nil {
							endpoint.ConnectionState = *conn.Properties.PrivateLinkServiceConnectionState.Status
						}

						for _, groupID := range conn.Properties.GroupIDs {
							if groupID != nil {
								endpoint.GroupIDs = append(endpoint.GroupIDs, *groupID)
							}
						}
					}
				}

				// Also check manual connections
				for _, conn := range pe.Properties.ManualPrivateLinkServiceConnections {
					if conn.Properties != nil {
						if conn.Properties.PrivateLinkServiceID != nil && endpoint.PrivateLinkServiceID == "" {
							endpoint.PrivateLinkServiceID = *conn.Properties.PrivateLinkServiceID
						}

						if conn.Properties.PrivateLinkServiceConnectionState != nil &&
							conn.Properties.PrivateLinkServiceConnectionState.Status != nil &&
							endpoint.ConnectionState == "" {
							endpoint.ConnectionState = *conn.Properties.PrivateLinkServiceConnectionState.Status
						}

						for _, groupID := range conn.Properties.GroupIDs {
							if groupID != nil {
								endpoint.GroupIDs = append(endpoint.GroupIDs, *groupID)
							}
						}
					}
				}
			}

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

// GetPrivateDNSZones retrieves all private DNS zones in the specified resource group
// Note: This requires the privatedns SDK package
func (c *AzureClient) GetPrivateDNSZones(ctx context.Context, resourceGroup string) ([]models.PrivateDNSZone, error) {
	// TODO: Implement using github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns
	// For now, return empty slice as this requires an additional SDK package
	return []models.PrivateDNSZone{}, nil
}

// GetPrivateDNSZoneVNetLinks retrieves VNet links for a specific private DNS zone
func (c *AzureClient) GetPrivateDNSZoneVNetLinks(ctx context.Context, resourceGroup, zoneName string) ([]models.VNetLink, error) {
	// TODO: Implement using github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns
	// For now, return empty slice
	return []models.VNetLink{}, nil
}
