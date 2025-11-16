package azure

import (
	"context"
	"fmt"

	"azure-network-analyzer/pkg/models"
)

// GetRouteTables retrieves all route tables in the specified resource group
func (c *AzureClient) GetRouteTables(ctx context.Context, resourceGroup string) ([]models.RouteTable, error) {
	client, err := c.getRouteTablesClient()
	if err != nil {
		return nil, err
	}

	var routeTables []models.RouteTable
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of Route Tables: %w", err)
		}

		for _, rt := range page.Value {
			table := models.RouteTable{
				ID:                safeString(rt.ID),
				Name:              safeString(rt.Name),
				ResourceGroup:     resourceGroup,
				Location:          safeString(rt.Location),
				Routes:            []models.Route{},
				AssociatedSubnets: []string{},
			}

			if rt.Properties != nil {
				if rt.Properties.DisableBgpRoutePropagation != nil {
					table.DisableBGPRoutePropagation = *rt.Properties.DisableBgpRoutePropagation
				}

				// Extract routes
				for _, route := range rt.Properties.Routes {
					r := c.extractRoute(route)
					table.Routes = append(table.Routes, r)
				}

				// Extract associated subnets
				for _, subnet := range rt.Properties.Subnets {
					if subnet.ID != nil {
						table.AssociatedSubnets = append(table.AssociatedSubnets, *subnet.ID)
					}
				}
			}

			routeTables = append(routeTables, table)
		}
	}

	return routeTables, nil
}

// GetRoutes retrieves all routes for a specific route table
func (c *AzureClient) GetRoutes(ctx context.Context, resourceGroup, routeTableName string) ([]models.Route, error) {
	client, err := c.getRoutesClient()
	if err != nil {
		return nil, err
	}

	var routes []models.Route
	pager := client.NewListPager(resourceGroup, routeTableName, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of routes: %w", err)
		}

		for _, route := range page.Value {
			r := c.extractRoute(route)
			routes = append(routes, r)
		}
	}

	return routes, nil
}

// GetNATGateways retrieves all NAT gateways in the specified resource group
func (c *AzureClient) GetNATGateways(ctx context.Context, resourceGroup string) ([]models.NATGateway, error) {
	client, err := c.getNATGatewaysClient()
	if err != nil {
		return nil, err
	}

	var natGateways []models.NATGateway
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of NAT Gateways: %w", err)
		}

		for _, nat := range page.Value {
			gw := models.NATGateway{
				ID:                safeString(nat.ID),
				Name:              safeString(nat.Name),
				ResourceGroup:     resourceGroup,
				Location:          safeString(nat.Location),
				PublicIPAddresses: []string{},
				AssociatedSubnets: []string{},
			}

			if nat.Properties != nil {
				if nat.Properties.IdleTimeoutInMinutes != nil {
					gw.IdleTimeoutMinutes = *nat.Properties.IdleTimeoutInMinutes
				}

				// Extract public IPs
				for _, pip := range nat.Properties.PublicIPAddresses {
					if pip.ID != nil {
						gw.PublicIPAddresses = append(gw.PublicIPAddresses, *pip.ID)
					}
				}

				// Extract associated subnets
				for _, subnet := range nat.Properties.Subnets {
					if subnet.ID != nil {
						gw.AssociatedSubnets = append(gw.AssociatedSubnets, *subnet.ID)
					}
				}
			}

			natGateways = append(natGateways, gw)
		}
	}

	return natGateways, nil
}

// GetNATGatewayPublicIPs retrieves the public IP addresses associated with a NAT gateway
func (c *AzureClient) GetNATGatewayPublicIPs(ctx context.Context, resourceGroup, natGatewayName string) ([]string, error) {
	client, err := c.getNATGatewaysClient()
	if err != nil {
		return nil, err
	}

	nat, err := client.Get(ctx, resourceGroup, natGatewayName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get NAT gateway %s: %w", natGatewayName, err)
	}

	var publicIPs []string
	if nat.Properties != nil {
		for _, pip := range nat.Properties.PublicIPAddresses {
			if pip.ID != nil {
				publicIPs = append(publicIPs, *pip.ID)
			}
		}
	}

	return publicIPs, nil
}
