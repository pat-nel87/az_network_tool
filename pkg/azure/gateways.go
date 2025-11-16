package azure

import (
	"context"
	"fmt"

	"azure-network-analyzer/pkg/models"
)

// GetVPNGateways retrieves all VPN gateways in the specified resource group
func (c *AzureClient) GetVPNGateways(ctx context.Context, resourceGroup string) ([]models.VPNGateway, error) {
	client, err := c.getVPNGatewaysClient()
	if err != nil {
		return nil, err
	}

	var vpnGateways []models.VPNGateway
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of VPN Gateways: %w", err)
		}

		for _, gw := range page.Value {
			gateway := models.VPNGateway{
				ID:            safeString(gw.ID),
				Name:          safeString(gw.Name),
				ResourceGroup: resourceGroup,
				Location:      safeString(gw.Location),
				Connections:   []models.VPNConnection{},
			}

			if gw.Properties != nil {
				if gw.Properties.GatewayType != nil {
					gateway.GatewayType = string(*gw.Properties.GatewayType)
				}

				if gw.Properties.VPNType != nil {
					gateway.VpnType = string(*gw.Properties.VPNType)
				}

				if gw.Properties.SKU != nil && gw.Properties.SKU.Name != nil {
					gateway.SKU = string(*gw.Properties.SKU.Name)
				}

				// Extract VNet ID from IP configurations
				if len(gw.Properties.IPConfigurations) > 0 {
					ipConfig := gw.Properties.IPConfigurations[0]
					if ipConfig.Properties != nil && ipConfig.Properties.Subnet != nil && ipConfig.Properties.Subnet.ID != nil {
						gateway.VNetID = extractVNetIDFromSubnet(*ipConfig.Properties.Subnet.ID)
					}
				}

				// Extract BGP settings
				if gw.Properties.BgpSettings != nil {
					gateway.BGPSettings = &models.BGPSettings{}
					if gw.Properties.BgpSettings.Asn != nil {
						gateway.BGPSettings.ASN = *gw.Properties.BgpSettings.Asn
					}
					if gw.Properties.BgpSettings.BgpPeeringAddress != nil {
						gateway.BGPSettings.BGPPeeringAddress = *gw.Properties.BgpSettings.BgpPeeringAddress
					}
					if gw.Properties.BgpSettings.PeerWeight != nil {
						gateway.BGPSettings.PeerWeight = *gw.Properties.BgpSettings.PeerWeight
					}
				}
			}

			vpnGateways = append(vpnGateways, gateway)
		}
	}

	return vpnGateways, nil
}

// GetVPNConnections retrieves all connections for a specific VPN gateway
func (c *AzureClient) GetVPNConnections(ctx context.Context, resourceGroup, gatewayName string) ([]models.VPNConnection, error) {
	client, err := c.getConnectionsClient()
	if err != nil {
		return nil, err
	}

	var connections []models.VPNConnection
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of VPN connections: %w", err)
		}

		for _, conn := range page.Value {
			// Filter connections for this gateway
			if conn.Properties != nil && conn.Properties.VirtualNetworkGateway1 != nil {
				gwID := safeString(conn.Properties.VirtualNetworkGateway1.ID)
				if extractResourceName(gwID) == gatewayName {
					c := models.VPNConnection{
						ID:   safeString(conn.ID),
						Name: safeString(conn.Name),
					}

					if conn.Properties.ConnectionType != nil {
						c.ConnectionType = string(*conn.Properties.ConnectionType)
					}

					if conn.Properties.ConnectionStatus != nil {
						c.ConnectionStatus = string(*conn.Properties.ConnectionStatus)
					}

					c.SharedKey = conn.Properties.SharedKey != nil && *conn.Properties.SharedKey != ""

					if conn.Properties.EnableBgp != nil {
						c.EnableBGP = *conn.Properties.EnableBgp
					}

					// Get remote entity ID
					if conn.Properties.VirtualNetworkGateway2 != nil && conn.Properties.VirtualNetworkGateway2.ID != nil {
						c.RemoteEntityID = *conn.Properties.VirtualNetworkGateway2.ID
					} else if conn.Properties.LocalNetworkGateway2 != nil && conn.Properties.LocalNetworkGateway2.ID != nil {
						c.RemoteEntityID = *conn.Properties.LocalNetworkGateway2.ID
					} else if conn.Properties.Peer != nil && conn.Properties.Peer.ID != nil {
						c.RemoteEntityID = *conn.Properties.Peer.ID
					}

					connections = append(connections, c)
				}
			}
		}
	}

	return connections, nil
}

// GetExpressRouteCircuits retrieves all ExpressRoute circuits in the specified resource group
func (c *AzureClient) GetExpressRouteCircuits(ctx context.Context, resourceGroup string) ([]models.ExpressRouteCircuit, error) {
	client, err := c.getERCircuitsClient()
	if err != nil {
		return nil, err
	}

	var circuits []models.ExpressRouteCircuit
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of ExpressRoute Circuits: %w", err)
		}

		for _, circuit := range page.Value {
			er := models.ExpressRouteCircuit{
				ID:             safeString(circuit.ID),
				Name:           safeString(circuit.Name),
				ResourceGroup:  resourceGroup,
				Location:       safeString(circuit.Location),
				Peerings:       []models.ERPeering{},
				Authorizations: []models.ERAuthorization{},
			}

			if circuit.SKU != nil {
				if circuit.SKU.Tier != nil {
					er.SKUTier = string(*circuit.SKU.Tier)
				}
				if circuit.SKU.Family != nil {
					er.SKUFamily = string(*circuit.SKU.Family)
				}
			}

			if circuit.Properties != nil {
				if circuit.Properties.ServiceProviderProperties != nil {
					if circuit.Properties.ServiceProviderProperties.ServiceProviderName != nil {
						er.ServiceProviderName = *circuit.Properties.ServiceProviderProperties.ServiceProviderName
					}
					if circuit.Properties.ServiceProviderProperties.PeeringLocation != nil {
						er.PeeringLocation = *circuit.Properties.ServiceProviderProperties.PeeringLocation
					}
					if circuit.Properties.ServiceProviderProperties.BandwidthInMbps != nil {
						er.BandwidthInMbps = *circuit.Properties.ServiceProviderProperties.BandwidthInMbps
					}
				}

				if circuit.Properties.CircuitProvisioningState != nil {
					er.CircuitProvisioningState = *circuit.Properties.CircuitProvisioningState
				}

				// Extract peerings
				for _, peering := range circuit.Properties.Peerings {
					p := c.extractERPeering(peering)
					er.Peerings = append(er.Peerings, p)
				}

				// Extract authorizations
				for _, auth := range circuit.Properties.Authorizations {
					a := c.extractERAuthorization(auth)
					er.Authorizations = append(er.Authorizations, a)
				}
			}

			circuits = append(circuits, er)
		}
	}

	return circuits, nil
}

// GetERPeerings retrieves all peerings for a specific ExpressRoute circuit
func (c *AzureClient) GetERPeerings(ctx context.Context, resourceGroup, circuitName string) ([]models.ERPeering, error) {
	client, err := c.getERPeeringsClient()
	if err != nil {
		return nil, err
	}

	var peerings []models.ERPeering
	pager := client.NewListPager(resourceGroup, circuitName, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of ER peerings: %w", err)
		}

		for _, peering := range page.Value {
			p := c.extractERPeering(peering)
			peerings = append(peerings, p)
		}
	}

	return peerings, nil
}

// GetERAuthorizations retrieves all authorizations for a specific ExpressRoute circuit
func (c *AzureClient) GetERAuthorizations(ctx context.Context, resourceGroup, circuitName string) ([]models.ERAuthorization, error) {
	client, err := c.getERAuthorizationsClient()
	if err != nil {
		return nil, err
	}

	var authorizations []models.ERAuthorization
	pager := client.NewListPager(resourceGroup, circuitName, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of ER authorizations: %w", err)
		}

		for _, auth := range page.Value {
			a := c.extractERAuthorization(auth)
			authorizations = append(authorizations, a)
		}
	}

	return authorizations, nil
}
