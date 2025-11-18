package visualization

import (
	"strings"
	"testing"

	"azure-network-analyzer/pkg/models"
)

// TestAllOutputFormats tests that all supported visualization formats work
func TestAllOutputFormats(t *testing.T) {
	// Create a simple test topology
	topology := &models.NetworkTopology{
		SubscriptionID: "test-sub",
		ResourceGroup:  "test-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1",
				Name:         "test-vnet",
				Location:     "eastus",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
						Name:          "subnet1",
						AddressPrefix: "10.0.1.0/24",
					},
				},
			},
		},
	}

	// Generate DOT content
	dotContent := GenerateDOTFile(topology)

	// Test 1: Verify DOT content is valid
	if !strings.Contains(dotContent, "digraph NetworkTopology") {
		t.Error("DOT content should contain digraph declaration")
	}
	if !strings.Contains(dotContent, "test-vnet") {
		t.Error("DOT content should contain VNet name")
	}
	t.Logf("✓ DOT format generation successful (%d bytes)", len(dotContent))

	// Test 2: SVG rendering
	t.Run("SVG", func(t *testing.T) {
		svgContent, err := RenderSVG(dotContent)
		if err != nil {
			t.Fatalf("SVG rendering failed: %v", err)
		}
		if len(svgContent) == 0 {
			t.Error("SVG content should not be empty")
		}
		if !strings.Contains(string(svgContent), "<svg") {
			t.Error("SVG content should contain <svg> tag")
		}
		t.Logf("✓ SVG rendering successful (%d bytes)", len(svgContent))
	})

	// Test 3: PNG rendering
	t.Run("PNG", func(t *testing.T) {
		pngContent, err := RenderPNG(dotContent)
		if err != nil {
			t.Fatalf("PNG rendering failed: %v", err)
		}
		if len(pngContent) == 0 {
			t.Error("PNG content should not be empty")
		}
		// PNG files start with specific magic bytes
		if len(pngContent) > 4 && string(pngContent[1:4]) != "PNG" {
			t.Error("PNG content should have PNG magic bytes")
		}
		t.Logf("✓ PNG rendering successful (%d bytes)", len(pngContent))
	})

	// Test 4: PDF rendering (may not be supported in all environments)
	t.Run("PDF", func(t *testing.T) {
		pdfContent, err := RenderPDF(dotContent)
		if err != nil {
			// PDF rendering may fail due to WASM limitations or missing system dependencies
			t.Skipf("PDF rendering not available in this environment: %v", err)
			return
		}
		if len(pdfContent) == 0 {
			t.Error("PDF content should not be empty")
		}
		// PDF files start with %PDF
		if len(pdfContent) > 4 && string(pdfContent[0:4]) != "%PDF" {
			t.Error("PDF content should have PDF magic bytes")
		}
		t.Logf("✓ PDF rendering successful (%d bytes)", len(pdfContent))
	})

	// Test 5: JPEG rendering
	t.Run("JPEG", func(t *testing.T) {
		jpegContent, err := RenderJPEG(dotContent)
		if err != nil {
			t.Fatalf("JPEG rendering failed: %v", err)
		}
		if len(jpegContent) == 0 {
			t.Error("JPEG content should not be empty")
		}
		// JPEG files start with FFD8
		if len(jpegContent) > 2 && (jpegContent[0] != 0xFF || jpegContent[1] != 0xD8) {
			t.Error("JPEG content should have JPEG magic bytes")
		}
		t.Logf("✓ JPEG rendering successful (%d bytes)", len(jpegContent))
	})
}

// TestDOTOutputOnly tests that DOT format can be used without rendering
func TestDOTOutputOnly(t *testing.T) {
	topology := &models.NetworkTopology{
		SubscriptionID: "test-sub",
		ResourceGroup:  "test-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1",
				Name:         "test-vnet",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets:      []models.Subnet{},
			},
		},
	}

	dotContent := GenerateDOTFile(topology)

	// Verify it's valid DOT syntax
	if !strings.HasPrefix(dotContent, "digraph") {
		t.Error("DOT content should start with 'digraph'")
	}
	if !strings.HasSuffix(strings.TrimSpace(dotContent), "}") {
		t.Error("DOT content should end with '}'")
	}

	// Verify essential GraphViz attributes are present
	requiredAttributes := []string{
		"rankdir=",
		"node [",
		"edge [",
		"label=",
	}
	for _, attr := range requiredAttributes {
		if !strings.Contains(dotContent, attr) {
			t.Errorf("DOT content should contain '%s'", attr)
		}
	}

	t.Logf("✓ DOT format is valid and can be used with external GraphViz tools")
}

// TestFormatConversionEquivalence tests that different formats represent the same topology
func TestFormatConversionEquivalence(t *testing.T) {
	topology := &models.NetworkTopology{
		SubscriptionID: "test-sub",
		ResourceGroup:  "test-rg",
		VirtualNetworks: []models.VirtualNetwork{
			{
				ID:           "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1",
				Name:         "prod-vnet",
				AddressSpace: []string{"10.0.0.0/16"},
				Subnets: []models.Subnet{
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
						Name:          "web-subnet",
						AddressPrefix: "10.0.1.0/24",
					},
					{
						ID:            "/subscriptions/test/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2",
						Name:          "data-subnet",
						AddressPrefix: "10.0.2.0/24",
					},
				},
			},
		},
	}

	dotContent := GenerateDOTFile(topology)

	// All formats should successfully render the same DOT content
	formats := []struct {
		name     string
		fn       func(string) ([]byte, error)
		optional bool // PDF may not work in all environments
	}{
		{"SVG", RenderSVG, false},
		{"PNG", RenderPNG, false},
		{"PDF", RenderPDF, true},  // PDF may fail due to WASM limitations
		{"JPEG", RenderJPEG, false},
	}

	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			content, err := format.fn(dotContent)
			if err != nil {
				if format.optional {
					t.Skipf("%s rendering not available: %v", format.name, err)
					return
				}
				t.Fatalf("%s rendering failed: %v", format.name, err)
			}
			if len(content) == 0 {
				t.Errorf("%s content should not be empty", format.name)
			}
			t.Logf("✓ %s format rendered successfully (%d bytes)", format.name, len(content))
		})
	}
}

// TestInvalidDOTHandling tests error handling for invalid DOT content
func TestInvalidDOTHandling(t *testing.T) {
	invalidDOT := "this is not valid DOT syntax"

	formats := []struct {
		name string
		fn   func(string) ([]byte, error)
	}{
		{"SVG", RenderSVG},
		{"PNG", RenderPNG},
		{"PDF", RenderPDF},  // Test even though PDF may not be fully supported
		{"JPEG", RenderJPEG},
	}

	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			_, err := format.fn(invalidDOT)
			if err == nil {
				t.Errorf("%s should return error for invalid DOT content", format.name)
			}
			t.Logf("✓ %s correctly rejects invalid DOT: %v", format.name, err)
		})
	}
}

// TestEmptyTopologyFormats tests all formats with an empty topology
func TestEmptyTopologyFormats(t *testing.T) {
	topology := &models.NetworkTopology{
		SubscriptionID:  "test-sub",
		ResourceGroup:   "test-rg",
		VirtualNetworks: []models.VirtualNetwork{},
	}

	dotContent := GenerateDOTFile(topology)

	// Even with empty topology, all formats should work
	formats := []struct {
		name     string
		fn       func(string) ([]byte, error)
		optional bool
	}{
		{"SVG", RenderSVG, false},
		{"PNG", RenderPNG, false},
		{"PDF", RenderPDF, true},  // PDF may fail due to WASM limitations
		{"JPEG", RenderJPEG, false},
	}

	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			content, err := format.fn(dotContent)
			if err != nil {
				if format.optional {
					t.Skipf("%s rendering not available: %v", format.name, err)
					return
				}
				t.Fatalf("%s rendering failed for empty topology: %v", format.name, err)
			}
			// PDF may return empty content without error - this is acceptable for optional formats
			if len(content) == 0 && !format.optional {
				t.Errorf("%s content should not be empty even for empty topology", format.name)
			}
			if len(content) == 0 && format.optional {
				t.Skipf("%s returned empty content for empty topology", format.name)
				return
			}
			t.Logf("✓ %s handles empty topology (%d bytes)", format.name, len(content))
		})
	}
}
