package visualization

import (
	"bytes"
	"context"
	"fmt"

	"azure-network-analyzer/pkg/models"

	"github.com/goccy/go-graphviz"
)

// TopologySizeWarning represents warnings about topology complexity
type TopologySizeWarning struct {
	TotalNodes   int
	TotalEdges   int
	IsLarge      bool
	IsVeryLarge  bool
	Message      string
}

// CheckTopologySize analyzes topology and returns warnings about size/complexity
func CheckTopologySize(topology *models.NetworkTopology) TopologySizeWarning {
	// Count total nodes
	nodes := 0
	nodes += len(topology.VirtualNetworks)
	for _, vnet := range topology.VirtualNetworks {
		nodes += len(vnet.Subnets)
	}
	nodes += len(topology.NSGs)
	nodes += len(topology.RouteTables)
	nodes += len(topology.LoadBalancers)
	nodes += len(topology.AppGateways)

	// Estimate edges (connections between nodes)
	edges := 0
	for _, vnet := range topology.VirtualNetworks {
		edges += len(vnet.Subnets)      // VNet -> Subnet edges
		edges += len(vnet.Peerings)     // Peering edges
	}
	edges += len(topology.NSGs) * 2       // NSG associations
	edges += len(topology.RouteTables) * 2 // Route table associations

	warning := TopologySizeWarning{
		TotalNodes: nodes,
		TotalEdges: edges,
	}

	// Thresholds based on stress testing results
	// - Under 100 nodes: Fast rendering (< 1 second)
	// - 100-500 nodes: Moderate (1-60 seconds)
	// - 500-1000 nodes: Slow (1-5 minutes)
	// - Over 1000 nodes: Very slow, potential memory issues

	if nodes > 1000 || edges > 2000 {
		warning.IsVeryLarge = true
		warning.IsLarge = true
		warning.Message = fmt.Sprintf(
			"WARNING: Topology is very large (%d nodes, %d edges). SVG rendering may take 5+ minutes and use significant memory. Consider using --viz-format=dot to generate DOT file only.",
			nodes, edges)
	} else if nodes > 500 || edges > 1000 {
		warning.IsLarge = true
		warning.Message = fmt.Sprintf(
			"NOTE: Topology is large (%d nodes, %d edges). SVG rendering may take 1-5 minutes.",
			nodes, edges)
	}

	return warning
}

// RenderSVG converts a DOT string to SVG format
func RenderSVG(dotContent string) ([]byte, error) {
	// Check DOT file size - warn if very large
	if len(dotContent) > 500000 { // 500KB
		fmt.Printf("WARNING: DOT file is large (%.2f MB). SVG rendering may take several minutes...\n",
			float64(len(dotContent))/1024/1024)
	}

	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create graphviz instance: %w", err)
	}
	defer g.Close()

	graph, err := graphviz.ParseBytes([]byte(dotContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DOT content: %w", err)
	}
	defer graph.Close()

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return nil, fmt.Errorf("failed to render SVG: %w", err)
	}

	return buf.Bytes(), nil
}

// RenderPNG converts a DOT string to PNG format
func RenderPNG(dotContent string) ([]byte, error) {
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create graphviz instance: %w", err)
	}
	defer g.Close()

	graph, err := graphviz.ParseBytes([]byte(dotContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DOT content: %w", err)
	}
	defer graph.Close()

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, graphviz.PNG, &buf); err != nil {
		return nil, fmt.Errorf("failed to render PNG: %w", err)
	}

	return buf.Bytes(), nil
}

// RenderPDF converts a DOT string to PDF format
// Note: PDF support depends on the underlying GraphViz installation
func RenderPDF(dotContent string) ([]byte, error) {
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create graphviz instance: %w", err)
	}
	defer g.Close()

	graph, err := graphviz.ParseBytes([]byte(dotContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DOT content: %w", err)
	}
	defer graph.Close()

	var buf bytes.Buffer
	// Use custom format string for PDF
	if err := g.Render(ctx, graph, "pdf", &buf); err != nil {
		return nil, fmt.Errorf("failed to render PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// RenderJPEG converts a DOT string to JPEG format
func RenderJPEG(dotContent string) ([]byte, error) {
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create graphviz instance: %w", err)
	}
	defer g.Close()

	graph, err := graphviz.ParseBytes([]byte(dotContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DOT content: %w", err)
	}
	defer graph.Close()

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, graphviz.JPG, &buf); err != nil {
		return nil, fmt.Errorf("failed to render JPEG: %w", err)
	}

	return buf.Bytes(), nil
}
