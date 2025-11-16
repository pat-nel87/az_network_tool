package visualization

import (
	"bytes"
	"context"
	"fmt"

	"github.com/goccy/go-graphviz"
)

// RenderSVG converts a DOT string to SVG format
func RenderSVG(dotContent string) ([]byte, error) {
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
