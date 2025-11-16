package reporter

import (
	"fmt"
	"strings"
	"time"

	"azure-network-analyzer/pkg/analyzer"
	"azure-network-analyzer/pkg/models"
)

// GenerateHTML creates a rich HTML report with embedded CSS
func GenerateHTML(topology *models.NetworkTopology, analysis *analyzer.AnalysisReport) string {
	var html strings.Builder

	critical, high, medium, low := countBySeverity(analysis.SecurityFindings)

	html.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Azure Network Topology Report</title>
    <style>
        :root {
            --critical: #dc3545;
            --high: #fd7e14;
            --medium: #ffc107;
            --low: #17a2b8;
            --info: #6c757d;
            --success: #28a745;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }

        .report-container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 30px;
        }

        h1 {
            color: #0078d4;
            border-bottom: 3px solid #0078d4;
            padding-bottom: 10px;
        }

        h2 {
            color: #444;
            margin-top: 30px;
            border-bottom: 1px solid #ddd;
            padding-bottom: 5px;
        }

        h3 {
            color: #555;
            margin-top: 20px;
        }

        .metadata {
            background: #e9ecef;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
        }

        .metadata p {
            margin: 5px 0;
        }

        .summary-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin: 20px 0;
        }

        .summary-card {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            border-left: 4px solid #0078d4;
        }

        .summary-card h4 {
            margin: 0 0 5px 0;
            color: #666;
            font-size: 0.9em;
        }

        .summary-card .value {
            font-size: 1.8em;
            font-weight: bold;
            color: #0078d4;
        }

        .severity-badge {
            display: inline-block;
            padding: 3px 10px;
            border-radius: 12px;
            color: white;
            font-size: 0.85em;
            font-weight: bold;
        }

        .severity-critical { background: var(--critical); }
        .severity-high { background: var(--high); }
        .severity-medium { background: var(--medium); color: #333; }
        .severity-low { background: var(--low); }

        .finding {
            background: #fff;
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 15px;
            margin: 10px 0;
            border-left: 4px solid var(--info);
        }

        .finding.critical { border-left-color: var(--critical); }
        .finding.high { border-left-color: var(--high); }
        .finding.medium { border-left-color: var(--medium); }
        .finding.low { border-left-color: var(--low); }

        .finding h4 {
            margin: 0 0 10px 0;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .finding .details {
            font-size: 0.9em;
            color: #666;
        }

        .finding .recommendation {
            background: #e8f4f8;
            padding: 10px;
            border-radius: 3px;
            margin-top: 10px;
            font-size: 0.9em;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin: 15px 0;
            font-size: 0.9em;
        }

        th, td {
            border: 1px solid #ddd;
            padding: 10px;
            text-align: left;
        }

        th {
            background: #f0f0f0;
            font-weight: 600;
        }

        tr:nth-child(even) {
            background: #fafafa;
        }

        tr:hover {
            background: #f0f7ff;
        }

        .resource-section {
            margin: 20px 0;
            padding: 15px;
            background: #fafafa;
            border-radius: 5px;
        }

        .collapsible {
            cursor: pointer;
            user-select: none;
        }

        .collapsible:before {
            content: '▼ ';
            font-size: 0.8em;
        }

        .collapsible.collapsed:before {
            content: '▶ ';
        }

        .content {
            overflow: hidden;
            transition: max-height 0.3s ease;
        }

        .recommendations {
            background: #e8f4e8;
            padding: 15px;
            border-radius: 5px;
            border-left: 4px solid var(--success);
        }

        .recommendations ol {
            margin: 10px 0;
            padding-left: 20px;
        }

        .recommendations li {
            margin: 8px 0;
        }

        .orphaned {
            background: #fff3cd;
            padding: 15px;
            border-radius: 5px;
            margin: 10px 0;
        }

        footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #ddd;
            text-align: center;
            color: #888;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="report-container">
`)

	// Header
	html.WriteString(`        <h1>Azure Network Topology Report</h1>
        <div class="metadata">
`)
	html.WriteString(fmt.Sprintf(`            <p><strong>Subscription:</strong> %s</p>
`, topology.SubscriptionID))
	html.WriteString(fmt.Sprintf(`            <p><strong>Resource Group:</strong> %s</p>
`, topology.ResourceGroup))
	html.WriteString(fmt.Sprintf(`            <p><strong>Generated:</strong> %s</p>
`, time.Now().Format("2006-01-02 15:04:05 MST")))
	html.WriteString(`        </div>
`)

	// Executive Summary
	html.WriteString(`        <h2>Executive Summary</h2>
        <div class="summary-grid">
`)
	html.WriteString(fmt.Sprintf(`            <div class="summary-card">
                <h4>Virtual Networks</h4>
                <div class="value">%d</div>
            </div>
`, analysis.Summary.TotalVNets))
	html.WriteString(fmt.Sprintf(`            <div class="summary-card">
                <h4>Subnets</h4>
                <div class="value">%d</div>
            </div>
`, analysis.Summary.TotalSubnets))
	html.WriteString(fmt.Sprintf(`            <div class="summary-card">
                <h4>NSGs</h4>
                <div class="value">%d</div>
            </div>
`, analysis.Summary.TotalNSGs))
	html.WriteString(fmt.Sprintf(`            <div class="summary-card">
                <h4>Security Rules</h4>
                <div class="value">%d</div>
            </div>
`, analysis.Summary.TotalSecurityRules))
	html.WriteString(fmt.Sprintf(`            <div class="summary-card">
                <h4>Security Findings</h4>
                <div class="value">%d</div>
            </div>
`, len(analysis.SecurityFindings)))
	html.WriteString(`        </div>
`)

	// Security Findings
	if len(analysis.SecurityFindings) > 0 {
		html.WriteString(`        <h2>Security Findings</h2>
        <p>
`)
		html.WriteString(fmt.Sprintf(`            <span class="severity-badge severity-critical">Critical: %d</span>
            <span class="severity-badge severity-high">High: %d</span>
            <span class="severity-badge severity-medium">Medium: %d</span>
            <span class="severity-badge severity-low">Low: %d</span>
        </p>
`, critical, high, medium, low))

		// Critical findings
		if critical > 0 {
			html.WriteString(`        <h3>Critical Issues</h3>
`)
			for _, f := range analysis.SecurityFindings {
				if f.Severity == analyzer.SeverityCritical {
					html.WriteString(`        <div class="finding critical">
`)
					html.WriteString(fmt.Sprintf(`            <h4>
                <span class="severity-badge severity-critical">CRITICAL</span>
                %s
            </h4>
`, f.Description))
					html.WriteString(fmt.Sprintf(`            <div class="details">
                <strong>Resource:</strong> %s<br>
`, f.Resource))
					if f.Rule != "" {
						html.WriteString(fmt.Sprintf(`                <strong>Rule:</strong> %s<br>
`, f.Rule))
					}
					html.WriteString(fmt.Sprintf(`                <strong>Category:</strong> %s
            </div>
`, f.Category))
					html.WriteString(fmt.Sprintf(`            <div class="recommendation">
                <strong>Recommendation:</strong> %s
            </div>
`, f.Recommendation))
					html.WriteString(`        </div>
`)
				}
			}
		}

		// High findings
		if high > 0 {
			html.WriteString(`        <h3>High Severity Issues</h3>
`)
			for _, f := range analysis.SecurityFindings {
				if f.Severity == analyzer.SeverityHigh {
					html.WriteString(`        <div class="finding high">
`)
					html.WriteString(fmt.Sprintf(`            <h4>
                <span class="severity-badge severity-high">HIGH</span>
                %s
            </h4>
`, f.Description))
					html.WriteString(fmt.Sprintf(`            <div class="details">
                <strong>Resource:</strong> %s
            </div>
`, f.Resource))
					html.WriteString(fmt.Sprintf(`            <div class="recommendation">
                <strong>Recommendation:</strong> %s
            </div>
`, f.Recommendation))
					html.WriteString(`        </div>
`)
				}
			}
		}

		// Medium and Low findings
		if medium > 0 || low > 0 {
			html.WriteString(`        <h3>Other Issues</h3>
        <table>
            <tr>
                <th>Severity</th>
                <th>Description</th>
                <th>Resource</th>
            </tr>
`)
			for _, f := range analysis.SecurityFindings {
				if f.Severity == analyzer.SeverityMedium || f.Severity == analyzer.SeverityLow {
					badgeClass := "severity-medium"
					if f.Severity == analyzer.SeverityLow {
						badgeClass = "severity-low"
					}
					html.WriteString(fmt.Sprintf(`            <tr>
                <td><span class="severity-badge %s">%s</span></td>
                <td>%s</td>
                <td>%s</td>
            </tr>
`, badgeClass, f.Severity, f.Description, f.Resource))
				}
			}
			html.WriteString(`        </table>
`)
		}
	}

	// Recommendations
	if len(analysis.Recommendations) > 0 {
		html.WriteString(`        <h2>Recommendations</h2>
        <div class="recommendations">
            <ol>
`)
		for _, rec := range analysis.Recommendations {
			html.WriteString(fmt.Sprintf(`                <li>%s</li>
`, rec))
		}
		html.WriteString(`            </ol>
        </div>
`)
	}

	// Network Topology Details
	html.WriteString(`        <h2>Network Topology Details</h2>
`)

	// Virtual Networks
	if len(topology.VirtualNetworks) > 0 {
		html.WriteString(`        <h3>Virtual Networks</h3>
`)
		for _, vnet := range topology.VirtualNetworks {
			html.WriteString(`        <div class="resource-section">
`)
			html.WriteString(fmt.Sprintf(`            <h4>%s</h4>
            <p>
                <strong>Location:</strong> %s<br>
                <strong>Address Space:</strong> %s<br>
                <strong>DDoS Protection:</strong> %v
            </p>
`, vnet.Name, vnet.Location, strings.Join(vnet.AddressSpace, ", "), vnet.EnableDDoS))

			if len(vnet.Subnets) > 0 {
				html.WriteString(`            <strong>Subnets:</strong>
            <table>
                <tr>
                    <th>Name</th>
                    <th>Address Prefix</th>
                    <th>NSG</th>
                    <th>Route Table</th>
                </tr>
`)
				for _, subnet := range vnet.Subnets {
					nsg := "-"
					if subnet.NetworkSecurityGroup != nil {
						nsg = extractName(*subnet.NetworkSecurityGroup)
					}
					rt := "-"
					if subnet.RouteTable != nil {
						rt = extractName(*subnet.RouteTable)
					}
					html.WriteString(fmt.Sprintf(`                <tr>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>
`, subnet.Name, subnet.AddressPrefix, nsg, rt))
				}
				html.WriteString(`            </table>
`)
			}
			html.WriteString(`        </div>
`)
		}
	}

	// NSGs
	if len(topology.NSGs) > 0 {
		html.WriteString(`        <h3>Network Security Groups</h3>
`)
		for _, nsg := range topology.NSGs {
			html.WriteString(`        <div class="resource-section">
`)
			html.WriteString(fmt.Sprintf(`            <h4>%s</h4>
            <p><strong>Location:</strong> %s</p>
`, nsg.Name, nsg.Location))

			if len(nsg.SecurityRules) > 0 {
				html.WriteString(`            <table>
                <tr>
                    <th>Priority</th>
                    <th>Name</th>
                    <th>Direction</th>
                    <th>Access</th>
                    <th>Protocol</th>
                    <th>Source</th>
                    <th>Dest Port</th>
                </tr>
`)
				for _, rule := range nsg.SecurityRules {
					html.WriteString(fmt.Sprintf(`                <tr>
                    <td>%d</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>
`, rule.Priority, rule.Name, rule.Direction, rule.Access, rule.Protocol, rule.SourceAddressPrefix, rule.DestinationPortRange))
				}
				html.WriteString(`            </table>
`)
			}
			html.WriteString(`        </div>
`)
		}
	}

	// Footer
	html.WriteString(`        <footer>
            Generated by Azure Network Topology Analyzer v1.0.0
        </footer>
    </div>
</body>
</html>
`)

	return html.String()
}
