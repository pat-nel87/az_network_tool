package azure

import (
	"context"
	"fmt"

	"azure-network-analyzer/pkg/models"
)

// GetLoadBalancers retrieves all load balancers in the specified resource group
func (c *AzureClient) GetLoadBalancers(ctx context.Context, resourceGroup string) ([]models.LoadBalancer, error) {
	client, err := c.getLoadBalancersClient()
	if err != nil {
		return nil, err
	}

	var loadBalancers []models.LoadBalancer
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of Load Balancers: %w", err)
		}

		for _, lb := range page.Value {
			balancer := models.LoadBalancer{
				ID:                  safeString(lb.ID),
				Name:                safeString(lb.Name),
				ResourceGroup:       resourceGroup,
				Location:            safeString(lb.Location),
				FrontendIPConfigs:   []models.FrontendIPConfig{},
				BackendAddressPools: []models.BackendAddressPool{},
				LoadBalancingRules:  []models.LoadBalancingRule{},
				Probes:              []models.Probe{},
				InboundNATRules:     []models.InboundNATRule{},
			}

			if lb.SKU != nil && lb.SKU.Name != nil {
				balancer.SKU = string(*lb.SKU.Name)
			}

			if lb.Properties != nil {
				// Determine if internal or public
				balancer.Type = "Internal"
				for _, feConfig := range lb.Properties.FrontendIPConfigurations {
					if feConfig.Properties != nil && feConfig.Properties.PublicIPAddress != nil {
						balancer.Type = "Public"
						break
					}
				}

				// Extract frontend IP configurations
				for _, feConfig := range lb.Properties.FrontendIPConfigurations {
					fe := c.extractFrontendIPConfig(feConfig)
					balancer.FrontendIPConfigs = append(balancer.FrontendIPConfigs, fe)
				}

				// Extract backend address pools
				for _, bePool := range lb.Properties.BackendAddressPools {
					be := c.extractBackendAddressPool(bePool)
					balancer.BackendAddressPools = append(balancer.BackendAddressPools, be)
				}

				// Extract load balancing rules
				for _, rule := range lb.Properties.LoadBalancingRules {
					r := c.extractLoadBalancingRule(rule)
					balancer.LoadBalancingRules = append(balancer.LoadBalancingRules, r)
				}

				// Extract probes
				for _, probe := range lb.Properties.Probes {
					p := c.extractProbe(probe)
					balancer.Probes = append(balancer.Probes, p)
				}

				// Extract inbound NAT rules
				for _, natRule := range lb.Properties.InboundNatRules {
					nat := c.extractInboundNATRule(natRule)
					balancer.InboundNATRules = append(balancer.InboundNATRules, nat)
				}
			}

			loadBalancers = append(loadBalancers, balancer)
		}
	}

	return loadBalancers, nil
}

// GetApplicationGateways retrieves all application gateways in the specified resource group
func (c *AzureClient) GetApplicationGateways(ctx context.Context, resourceGroup string) ([]models.ApplicationGateway, error) {
	client, err := c.getAppGatewaysClient()
	if err != nil {
		return nil, err
	}

	var appGateways []models.ApplicationGateway
	pager := client.NewListPager(resourceGroup, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of Application Gateways: %w", err)
		}

		for _, ag := range page.Value {
			gateway := models.ApplicationGateway{
				ID:                  safeString(ag.ID),
				Name:                safeString(ag.Name),
				ResourceGroup:       resourceGroup,
				Location:            safeString(ag.Location),
				FrontendIPConfigs:   []models.AppGWFrontendIPConfig{},
				FrontendPorts:       []models.AppGWFrontendPort{},
				BackendAddressPools: []models.AppGWBackendAddressPool{},
				BackendHTTPSettings: []models.AppGWBackendHTTPSettings{},
				HTTPListeners:       []models.AppGWHTTPListener{},
				RequestRoutingRules: []models.AppGWRequestRoutingRule{},
				Probes:              []models.AppGWProbe{},
			}

			if ag.Properties != nil {
				if ag.Properties.SKU != nil {
					if ag.Properties.SKU.Name != nil {
						gateway.SKU = string(*ag.Properties.SKU.Name)
					}
					if ag.Properties.SKU.Tier != nil {
						gateway.Tier = string(*ag.Properties.SKU.Tier)
					}
					if ag.Properties.SKU.Capacity != nil {
						gateway.Capacity = *ag.Properties.SKU.Capacity
					}
				}

				// Extract subnet from gateway IP configurations
				if len(ag.Properties.GatewayIPConfigurations) > 0 {
					ipConfig := ag.Properties.GatewayIPConfigurations[0]
					if ipConfig.Properties != nil && ipConfig.Properties.Subnet != nil && ipConfig.Properties.Subnet.ID != nil {
						gateway.SubnetID = *ipConfig.Properties.Subnet.ID
					}
				}

				// Check WAF configuration
				if ag.Properties.WebApplicationFirewallConfiguration != nil {
					if ag.Properties.WebApplicationFirewallConfiguration.Enabled != nil {
						gateway.WAFEnabled = *ag.Properties.WebApplicationFirewallConfiguration.Enabled
					}
					if ag.Properties.WebApplicationFirewallConfiguration.FirewallMode != nil {
						gateway.WAFMode = string(*ag.Properties.WebApplicationFirewallConfiguration.FirewallMode)
					}
				}

				// Extract frontend IP configurations
				for _, feConfig := range ag.Properties.FrontendIPConfigurations {
					fe := c.extractAppGWFrontendIPConfig(feConfig)
					gateway.FrontendIPConfigs = append(gateway.FrontendIPConfigs, fe)
				}

				// Extract frontend ports
				for _, port := range ag.Properties.FrontendPorts {
					fp := c.extractAppGWFrontendPort(port)
					gateway.FrontendPorts = append(gateway.FrontendPorts, fp)
				}

				// Extract backend address pools
				for _, bePool := range ag.Properties.BackendAddressPools {
					be := c.extractAppGWBackendAddressPool(bePool)
					gateway.BackendAddressPools = append(gateway.BackendAddressPools, be)
				}

				// Extract backend HTTP settings
				for _, settings := range ag.Properties.BackendHTTPSettingsCollection {
					s := c.extractAppGWBackendHTTPSettings(settings)
					gateway.BackendHTTPSettings = append(gateway.BackendHTTPSettings, s)
				}

				// Extract HTTP listeners
				for _, listener := range ag.Properties.HTTPListeners {
					l := c.extractAppGWHTTPListener(listener)
					gateway.HTTPListeners = append(gateway.HTTPListeners, l)
				}

				// Extract request routing rules
				for _, rule := range ag.Properties.RequestRoutingRules {
					r := c.extractAppGWRequestRoutingRule(rule)
					gateway.RequestRoutingRules = append(gateway.RequestRoutingRules, r)
				}

				// Extract probes
				for _, probe := range ag.Properties.Probes {
					p := c.extractAppGWProbe(probe)
					gateway.Probes = append(gateway.Probes, p)
				}
			}

			appGateways = append(appGateways, gateway)
		}
	}

	return appGateways, nil
}
