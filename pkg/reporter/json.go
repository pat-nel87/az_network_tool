package reporter

import (
	"encoding/json"
	"time"

	"azure-network-analyzer/pkg/analyzer"
	"azure-network-analyzer/pkg/models"
)

// JSONReport is the complete JSON report structure
type JSONReport struct {
	Metadata  ReportMetadata            `json:"metadata"`
	Topology  *models.NetworkTopology   `json:"topology"`
	Analysis  *analyzer.AnalysisReport  `json:"analysis"`
}

// ReportMetadata contains report generation information
type ReportMetadata struct {
	GeneratedAt    time.Time `json:"generated_at"`
	ToolVersion    string    `json:"tool_version"`
	SubscriptionID string    `json:"subscription_id"`
	ResourceGroup  string    `json:"resource_group"`
}

// GenerateJSON creates a complete JSON report
func GenerateJSON(topology *models.NetworkTopology, analysis *analyzer.AnalysisReport) ([]byte, error) {
	report := JSONReport{
		Metadata: ReportMetadata{
			GeneratedAt:    time.Now(),
			ToolVersion:    "1.0.0",
			SubscriptionID: topology.SubscriptionID,
			ResourceGroup:  topology.ResourceGroup,
		},
		Topology: topology,
		Analysis: analysis,
	}

	return json.MarshalIndent(report, "", "  ")
}
