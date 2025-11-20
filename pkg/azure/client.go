package azure

import (
	"fmt"

	"azure-network-analyzer/pkg/models"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// AzureClient wraps Azure SDK clients for network resource operations
type AzureClient struct {
	cred           *azidentity.DefaultAzureCredential
	subscriptionID string

	// Cached clients - lazily initialized
	vnetsClient            *armnetwork.VirtualNetworksClient
	subnetsClient          *armnetwork.SubnetsClient
	peeringsClient         *armnetwork.VirtualNetworkPeeringsClient
	nsgsClient             *armnetwork.SecurityGroupsClient
	privateEndpointsClient *armnetwork.PrivateEndpointsClient
	routeTablesClient      *armnetwork.RouteTablesClient
	routesClient           *armnetwork.RoutesClient
	natGatewaysClient      *armnetwork.NatGatewaysClient
	vpnGatewaysClient      *armnetwork.VirtualNetworkGatewaysClient
	connectionsClient      *armnetwork.VirtualNetworkGatewayConnectionsClient
	erCircuitsClient       *armnetwork.ExpressRouteCircuitsClient
	erPeeringsClient       *armnetwork.ExpressRouteCircuitPeeringsClient
	erAuthorizationsClient *armnetwork.ExpressRouteCircuitAuthorizationsClient
	loadBalancersClient    *armnetwork.LoadBalancersClient
	appGatewaysClient      *armnetwork.ApplicationGatewaysClient
	azureFirewallsClient   *armnetwork.AzureFirewallsClient
}

// NewAzureClient creates a new Azure client with DefaultAzureCredential
func NewAzureClient(subscriptionID string) (*AzureClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	}

	return &AzureClient{
		cred:           cred,
		subscriptionID: subscriptionID,
	}, nil
}

// Client factory methods - lazily initialize clients as needed

func (c *AzureClient) getVNetsClient() (*armnetwork.VirtualNetworksClient, error) {
	if c.vnetsClient == nil {
		client, err := armnetwork.NewVirtualNetworksClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create VNets client: %w", err)
		}
		c.vnetsClient = client
	}
	return c.vnetsClient, nil
}

func (c *AzureClient) getSubnetsClient() (*armnetwork.SubnetsClient, error) {
	if c.subnetsClient == nil {
		client, err := armnetwork.NewSubnetsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Subnets client: %w", err)
		}
		c.subnetsClient = client
	}
	return c.subnetsClient, nil
}

func (c *AzureClient) getPeeringsClient() (*armnetwork.VirtualNetworkPeeringsClient, error) {
	if c.peeringsClient == nil {
		client, err := armnetwork.NewVirtualNetworkPeeringsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Peerings client: %w", err)
		}
		c.peeringsClient = client
	}
	return c.peeringsClient, nil
}

func (c *AzureClient) getNSGsClient() (*armnetwork.SecurityGroupsClient, error) {
	if c.nsgsClient == nil {
		client, err := armnetwork.NewSecurityGroupsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create NSGs client: %w", err)
		}
		c.nsgsClient = client
	}
	return c.nsgsClient, nil
}

func (c *AzureClient) getPrivateEndpointsClient() (*armnetwork.PrivateEndpointsClient, error) {
	if c.privateEndpointsClient == nil {
		client, err := armnetwork.NewPrivateEndpointsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Private Endpoints client: %w", err)
		}
		c.privateEndpointsClient = client
	}
	return c.privateEndpointsClient, nil
}

func (c *AzureClient) getRouteTablesClient() (*armnetwork.RouteTablesClient, error) {
	if c.routeTablesClient == nil {
		client, err := armnetwork.NewRouteTablesClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Route Tables client: %w", err)
		}
		c.routeTablesClient = client
	}
	return c.routeTablesClient, nil
}

func (c *AzureClient) getRoutesClient() (*armnetwork.RoutesClient, error) {
	if c.routesClient == nil {
		client, err := armnetwork.NewRoutesClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Routes client: %w", err)
		}
		c.routesClient = client
	}
	return c.routesClient, nil
}

func (c *AzureClient) getNATGatewaysClient() (*armnetwork.NatGatewaysClient, error) {
	if c.natGatewaysClient == nil {
		client, err := armnetwork.NewNatGatewaysClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create NAT Gateways client: %w", err)
		}
		c.natGatewaysClient = client
	}
	return c.natGatewaysClient, nil
}

func (c *AzureClient) getVPNGatewaysClient() (*armnetwork.VirtualNetworkGatewaysClient, error) {
	if c.vpnGatewaysClient == nil {
		client, err := armnetwork.NewVirtualNetworkGatewaysClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create VPN Gateways client: %w", err)
		}
		c.vpnGatewaysClient = client
	}
	return c.vpnGatewaysClient, nil
}

func (c *AzureClient) getConnectionsClient() (*armnetwork.VirtualNetworkGatewayConnectionsClient, error) {
	if c.connectionsClient == nil {
		client, err := armnetwork.NewVirtualNetworkGatewayConnectionsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create VPN Connections client: %w", err)
		}
		c.connectionsClient = client
	}
	return c.connectionsClient, nil
}

func (c *AzureClient) getERCircuitsClient() (*armnetwork.ExpressRouteCircuitsClient, error) {
	if c.erCircuitsClient == nil {
		client, err := armnetwork.NewExpressRouteCircuitsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create ExpressRoute Circuits client: %w", err)
		}
		c.erCircuitsClient = client
	}
	return c.erCircuitsClient, nil
}

func (c *AzureClient) getERPeeringsClient() (*armnetwork.ExpressRouteCircuitPeeringsClient, error) {
	if c.erPeeringsClient == nil {
		client, err := armnetwork.NewExpressRouteCircuitPeeringsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create ExpressRoute Peerings client: %w", err)
		}
		c.erPeeringsClient = client
	}
	return c.erPeeringsClient, nil
}

func (c *AzureClient) getERAuthorizationsClient() (*armnetwork.ExpressRouteCircuitAuthorizationsClient, error) {
	if c.erAuthorizationsClient == nil {
		client, err := armnetwork.NewExpressRouteCircuitAuthorizationsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create ExpressRoute Authorizations client: %w", err)
		}
		c.erAuthorizationsClient = client
	}
	return c.erAuthorizationsClient, nil
}

func (c *AzureClient) getLoadBalancersClient() (*armnetwork.LoadBalancersClient, error) {
	if c.loadBalancersClient == nil {
		client, err := armnetwork.NewLoadBalancersClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Load Balancers client: %w", err)
		}
		c.loadBalancersClient = client
	}
	return c.loadBalancersClient, nil
}

func (c *AzureClient) getAppGatewaysClient() (*armnetwork.ApplicationGatewaysClient, error) {
	if c.appGatewaysClient == nil {
		client, err := armnetwork.NewApplicationGatewaysClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Application Gateways client: %w", err)
		}
		c.appGatewaysClient = client
	}
	return c.appGatewaysClient, nil
}

func (c *AzureClient) getAzureFirewallsClient() (*armnetwork.AzureFirewallsClient, error) {
	if c.azureFirewallsClient == nil {
		client, err := armnetwork.NewAzureFirewallsClient(c.subscriptionID, c.cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure Firewalls client: %w", err)
		}
		c.azureFirewallsClient = client
	}
	return c.azureFirewallsClient, nil
}

// Helper functions for extracting data from Azure SDK types

// safeString safely dereferences a string pointer
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// extractResourceName extracts the resource name from an Azure resource ID
func extractResourceName(resourceID string) string {
	if resourceID == "" {
		return ""
	}

	// Handle trailing slash case - return empty string
	if resourceID[len(resourceID)-1] == '/' {
		return ""
	}

	lastSlash := len(resourceID) - 1
	for lastSlash >= 0 && resourceID[lastSlash] != '/' {
		lastSlash--
	}

	if lastSlash >= 0 && lastSlash < len(resourceID)-1 {
		return resourceID[lastSlash+1:]
	}

	return resourceID
}

// extractVNetIDFromSubnet extracts the VNet ID from a subnet ID
func extractVNetIDFromSubnet(subnetID string) string {
	if subnetID == "" {
		return ""
	}

	subnetsIndex := -1
	for i := 0; i <= len(subnetID)-9; i++ {
		if subnetID[i:i+9] == "/subnets/" {
			subnetsIndex = i
			break
		}
	}

	if subnetsIndex > 0 {
		return subnetID[:subnetsIndex]
	}

	// If /subnets/ is at position 0, return empty string
	if subnetsIndex == 0 {
		return ""
	}

	return subnetID
}

// Extractor methods for complex Azure SDK types

func (c *AzureClient) extractSubnet(subnet *armnetwork.Subnet) models.Subnet {
	s := models.Subnet{
		ID:               safeString(subnet.ID),
		Name:             safeString(subnet.Name),
		AddressPrefix:    "",
		PrivateEndpoints: []string{},
		ServiceEndpoints: []string{},
		Delegations:      []string{},
	}

	if subnet.Properties != nil {
		if subnet.Properties.AddressPrefix != nil {
			s.AddressPrefix = *subnet.Properties.AddressPrefix
		}

		// NSG association
		if subnet.Properties.NetworkSecurityGroup != nil && subnet.Properties.NetworkSecurityGroup.ID != nil {
			s.NetworkSecurityGroup = subnet.Properties.NetworkSecurityGroup.ID
		}

		// Route table association
		if subnet.Properties.RouteTable != nil && subnet.Properties.RouteTable.ID != nil {
			s.RouteTable = subnet.Properties.RouteTable.ID
		}

		// NAT gateway association
		if subnet.Properties.NatGateway != nil && subnet.Properties.NatGateway.ID != nil {
			s.NATGateway = subnet.Properties.NatGateway.ID
		}

		// Private endpoints
		for _, pe := range subnet.Properties.PrivateEndpoints {
			if pe.ID != nil {
				s.PrivateEndpoints = append(s.PrivateEndpoints, *pe.ID)
			}
		}

		// Service endpoints
		for _, se := range subnet.Properties.ServiceEndpoints {
			if se.Service != nil {
				s.ServiceEndpoints = append(s.ServiceEndpoints, *se.Service)
			}
		}

		// Delegations
		for _, del := range subnet.Properties.Delegations {
			if del.Properties != nil && del.Properties.ServiceName != nil {
				s.Delegations = append(s.Delegations, *del.Properties.ServiceName)
			}
		}
	}

	return s
}

func (c *AzureClient) extractVNetPeering(peering *armnetwork.VirtualNetworkPeering) models.VNetPeering {
	p := models.VNetPeering{
		ID:             safeString(peering.ID),
		Name:           safeString(peering.Name),
		RemoteVNetID:   "",
		RemoteVNetName: "",
		PeeringState:   "",
	}

	if peering.Properties != nil {
		if peering.Properties.RemoteVirtualNetwork != nil && peering.Properties.RemoteVirtualNetwork.ID != nil {
			p.RemoteVNetID = *peering.Properties.RemoteVirtualNetwork.ID
			p.RemoteVNetName = extractResourceName(p.RemoteVNetID)
		}

		if peering.Properties.PeeringState != nil {
			p.PeeringState = string(*peering.Properties.PeeringState)
		}

		if peering.Properties.AllowVirtualNetworkAccess != nil {
			p.AllowVNetAccess = *peering.Properties.AllowVirtualNetworkAccess
		}

		if peering.Properties.AllowForwardedTraffic != nil {
			p.AllowForwardedTraffic = *peering.Properties.AllowForwardedTraffic
		}

		if peering.Properties.AllowGatewayTransit != nil {
			p.AllowGatewayTransit = *peering.Properties.AllowGatewayTransit
		}

		if peering.Properties.UseRemoteGateways != nil {
			p.UseRemoteGateways = *peering.Properties.UseRemoteGateways
		}
	}

	return p
}

func (c *AzureClient) extractRoute(route *armnetwork.Route) models.Route {
	r := models.Route{
		Name:          safeString(route.Name),
		AddressPrefix: "",
		NextHopType:   "",
	}

	if route.Properties != nil {
		if route.Properties.AddressPrefix != nil {
			r.AddressPrefix = *route.Properties.AddressPrefix
		}

		if route.Properties.NextHopType != nil {
			r.NextHopType = string(*route.Properties.NextHopType)
		}

		if route.Properties.NextHopIPAddress != nil {
			r.NextHopIPAddress = *route.Properties.NextHopIPAddress
		}
	}

	return r
}

func (c *AzureClient) extractERPeering(peering *armnetwork.ExpressRouteCircuitPeering) models.ERPeering {
	p := models.ERPeering{
		Name: safeString(peering.Name),
	}

	if peering.Properties != nil {
		if peering.Properties.PeeringType != nil {
			p.PeeringType = string(*peering.Properties.PeeringType)
		}
		if peering.Properties.State != nil {
			p.State = string(*peering.Properties.State)
		}
		if peering.Properties.AzureASN != nil {
			p.AzureASN = *peering.Properties.AzureASN
		}
		if peering.Properties.PeerASN != nil {
			p.PeerASN = *peering.Properties.PeerASN
		}
		if peering.Properties.PrimaryPeerAddressPrefix != nil {
			p.PrimaryPeerAddressPrefix = *peering.Properties.PrimaryPeerAddressPrefix
		}
		if peering.Properties.SecondaryPeerAddressPrefix != nil {
			p.SecondaryPeerAddressPrefix = *peering.Properties.SecondaryPeerAddressPrefix
		}
		if peering.Properties.VlanID != nil {
			p.VlanID = *peering.Properties.VlanID
		}
	}

	return p
}

func (c *AzureClient) extractERAuthorization(auth *armnetwork.ExpressRouteCircuitAuthorization) models.ERAuthorization {
	a := models.ERAuthorization{
		Name: safeString(auth.Name),
	}

	if auth.Properties != nil {
		a.AuthorizationKey = auth.Properties.AuthorizationKey != nil && *auth.Properties.AuthorizationKey != ""
		if auth.Properties.AuthorizationUseStatus != nil {
			a.AuthorizationUseStatus = string(*auth.Properties.AuthorizationUseStatus)
		}
	}

	return a
}

func (c *AzureClient) extractFrontendIPConfig(feConfig *armnetwork.FrontendIPConfiguration) models.FrontendIPConfig {
	fe := models.FrontendIPConfig{
		Name: safeString(feConfig.Name),
	}

	if feConfig.Properties != nil {
		if feConfig.Properties.PrivateIPAddress != nil {
			fe.PrivateIPAddress = *feConfig.Properties.PrivateIPAddress
		}
		if feConfig.Properties.PublicIPAddress != nil && feConfig.Properties.PublicIPAddress.ID != nil {
			fe.PublicIPAddressID = *feConfig.Properties.PublicIPAddress.ID
		}
		if feConfig.Properties.Subnet != nil && feConfig.Properties.Subnet.ID != nil {
			fe.SubnetID = *feConfig.Properties.Subnet.ID
		}
	}

	return fe
}

func (c *AzureClient) extractBackendAddressPool(bePool *armnetwork.BackendAddressPool) models.BackendAddressPool {
	be := models.BackendAddressPool{
		Name:             safeString(bePool.Name),
		BackendIPConfigs: []string{},
	}

	if bePool.Properties != nil {
		for _, ipConfig := range bePool.Properties.BackendIPConfigurations {
			if ipConfig.ID != nil {
				be.BackendIPConfigs = append(be.BackendIPConfigs, *ipConfig.ID)
			}
		}
	}

	return be
}

func (c *AzureClient) extractLoadBalancingRule(rule *armnetwork.LoadBalancingRule) models.LoadBalancingRule {
	r := models.LoadBalancingRule{
		Name: safeString(rule.Name),
	}

	if rule.Properties != nil {
		if rule.Properties.Protocol != nil {
			r.Protocol = string(*rule.Properties.Protocol)
		}
		if rule.Properties.FrontendPort != nil {
			r.FrontendPort = *rule.Properties.FrontendPort
		}
		if rule.Properties.BackendPort != nil {
			r.BackendPort = *rule.Properties.BackendPort
		}
		if rule.Properties.EnableFloatingIP != nil {
			r.EnableFloatingIP = *rule.Properties.EnableFloatingIP
		}
		if rule.Properties.IdleTimeoutInMinutes != nil {
			r.IdleTimeoutMinutes = *rule.Properties.IdleTimeoutInMinutes
		}
		if rule.Properties.LoadDistribution != nil {
			r.LoadDistribution = string(*rule.Properties.LoadDistribution)
		}
	}

	return r
}

func (c *AzureClient) extractProbe(probe *armnetwork.Probe) models.Probe {
	p := models.Probe{
		Name: safeString(probe.Name),
	}

	if probe.Properties != nil {
		if probe.Properties.Protocol != nil {
			p.Protocol = string(*probe.Properties.Protocol)
		}
		if probe.Properties.Port != nil {
			p.Port = *probe.Properties.Port
		}
		if probe.Properties.IntervalInSeconds != nil {
			p.IntervalInSeconds = *probe.Properties.IntervalInSeconds
		}
		if probe.Properties.NumberOfProbes != nil {
			p.NumberOfProbes = *probe.Properties.NumberOfProbes
		}
		if probe.Properties.RequestPath != nil {
			p.RequestPath = *probe.Properties.RequestPath
		}
	}

	return p
}

func (c *AzureClient) extractInboundNATRule(natRule *armnetwork.InboundNatRule) models.InboundNATRule {
	nat := models.InboundNATRule{
		Name: safeString(natRule.Name),
	}

	if natRule.Properties != nil {
		if natRule.Properties.Protocol != nil {
			nat.Protocol = string(*natRule.Properties.Protocol)
		}
		if natRule.Properties.FrontendPort != nil {
			nat.FrontendPort = *natRule.Properties.FrontendPort
		}
		if natRule.Properties.BackendPort != nil {
			nat.BackendPort = *natRule.Properties.BackendPort
		}
		if natRule.Properties.EnableFloatingIP != nil {
			nat.EnableFloatingIP = *natRule.Properties.EnableFloatingIP
		}
	}

	return nat
}

func (c *AzureClient) extractAppGWFrontendIPConfig(feConfig *armnetwork.ApplicationGatewayFrontendIPConfiguration) models.AppGWFrontendIPConfig {
	fe := models.AppGWFrontendIPConfig{
		Name: safeString(feConfig.Name),
	}

	if feConfig.Properties != nil {
		if feConfig.Properties.PrivateIPAddress != nil {
			fe.PrivateIPAddress = *feConfig.Properties.PrivateIPAddress
		}
		if feConfig.Properties.PublicIPAddress != nil && feConfig.Properties.PublicIPAddress.ID != nil {
			fe.PublicIPAddressID = *feConfig.Properties.PublicIPAddress.ID
		}
	}

	return fe
}

func (c *AzureClient) extractAppGWFrontendPort(port *armnetwork.ApplicationGatewayFrontendPort) models.AppGWFrontendPort {
	fp := models.AppGWFrontendPort{
		Name: safeString(port.Name),
	}

	if port.Properties != nil && port.Properties.Port != nil {
		fp.Port = *port.Properties.Port
	}

	return fp
}

func (c *AzureClient) extractAppGWBackendAddressPool(bePool *armnetwork.ApplicationGatewayBackendAddressPool) models.AppGWBackendAddressPool {
	be := models.AppGWBackendAddressPool{
		Name:             safeString(bePool.Name),
		BackendAddresses: []string{},
	}

	if bePool.Properties != nil {
		for _, addr := range bePool.Properties.BackendAddresses {
			if addr.IPAddress != nil {
				be.BackendAddresses = append(be.BackendAddresses, *addr.IPAddress)
			} else if addr.Fqdn != nil {
				be.BackendAddresses = append(be.BackendAddresses, *addr.Fqdn)
			}
		}
	}

	return be
}

func (c *AzureClient) extractAppGWBackendHTTPSettings(settings *armnetwork.ApplicationGatewayBackendHTTPSettings) models.AppGWBackendHTTPSettings {
	s := models.AppGWBackendHTTPSettings{
		Name: safeString(settings.Name),
	}

	if settings.Properties != nil {
		if settings.Properties.Port != nil {
			s.Port = *settings.Properties.Port
		}
		if settings.Properties.Protocol != nil {
			s.Protocol = string(*settings.Properties.Protocol)
		}
		if settings.Properties.CookieBasedAffinity != nil {
			s.CookieBasedAffinity = string(*settings.Properties.CookieBasedAffinity)
		}
		if settings.Properties.RequestTimeout != nil {
			s.RequestTimeout = *settings.Properties.RequestTimeout
		}
		if settings.Properties.Probe != nil && settings.Properties.Probe.ID != nil {
			s.ProbeName = extractResourceName(*settings.Properties.Probe.ID)
		}
	}

	return s
}

func (c *AzureClient) extractAppGWHTTPListener(listener *armnetwork.ApplicationGatewayHTTPListener) models.AppGWHTTPListener {
	l := models.AppGWHTTPListener{
		Name: safeString(listener.Name),
	}

	if listener.Properties != nil {
		if listener.Properties.FrontendIPConfiguration != nil && listener.Properties.FrontendIPConfiguration.ID != nil {
			l.FrontendIPConfig = extractResourceName(*listener.Properties.FrontendIPConfiguration.ID)
		}
		if listener.Properties.FrontendPort != nil && listener.Properties.FrontendPort.ID != nil {
			l.FrontendPort = extractResourceName(*listener.Properties.FrontendPort.ID)
		}
		if listener.Properties.Protocol != nil {
			l.Protocol = string(*listener.Properties.Protocol)
		}
		if listener.Properties.HostName != nil {
			l.HostName = *listener.Properties.HostName
		}
	}

	return l
}

func (c *AzureClient) extractAppGWRequestRoutingRule(rule *armnetwork.ApplicationGatewayRequestRoutingRule) models.AppGWRequestRoutingRule {
	r := models.AppGWRequestRoutingRule{
		Name: safeString(rule.Name),
	}

	if rule.Properties != nil {
		if rule.Properties.RuleType != nil {
			r.RuleType = string(*rule.Properties.RuleType)
		}
		if rule.Properties.HTTPListener != nil && rule.Properties.HTTPListener.ID != nil {
			r.HTTPListener = extractResourceName(*rule.Properties.HTTPListener.ID)
		}
		if rule.Properties.BackendAddressPool != nil && rule.Properties.BackendAddressPool.ID != nil {
			r.BackendAddressPool = extractResourceName(*rule.Properties.BackendAddressPool.ID)
		}
		if rule.Properties.BackendHTTPSettings != nil && rule.Properties.BackendHTTPSettings.ID != nil {
			r.BackendHTTPSettings = extractResourceName(*rule.Properties.BackendHTTPSettings.ID)
		}
		if rule.Properties.Priority != nil {
			r.Priority = *rule.Properties.Priority
		}
	}

	return r
}

func (c *AzureClient) extractAppGWProbe(probe *armnetwork.ApplicationGatewayProbe) models.AppGWProbe {
	p := models.AppGWProbe{
		Name: safeString(probe.Name),
	}

	if probe.Properties != nil {
		if probe.Properties.Protocol != nil {
			p.Protocol = string(*probe.Properties.Protocol)
		}
		if probe.Properties.Host != nil {
			p.Host = *probe.Properties.Host
		}
		if probe.Properties.Path != nil {
			p.Path = *probe.Properties.Path
		}
		if probe.Properties.Interval != nil {
			p.Interval = *probe.Properties.Interval
		}
		if probe.Properties.Timeout != nil {
			p.Timeout = *probe.Properties.Timeout
		}
		if probe.Properties.UnhealthyThreshold != nil {
			p.UnhealthyThreshold = *probe.Properties.UnhealthyThreshold
		}
	}

	return p
}
