package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/De-pitcher/any2pdf/internal/converter"
)

func TestHTMLIntegration_SimpleHTML(t *testing.T) {
	c := converter.NewHTMLConverter()
	
	// Skip if no HTML converters available
	engines := c.GetAvailableEngines()
	if len(engines) == 0 {
		t.Skip("No HTML converters available (wkhtmltopdf or Chrome required)")
	}
	
	t.Logf("Available engines: %v", engines)
	
	// Test file paths
	inputPath := filepath.Join("fixtures", "html", "simple.html")
	outputPath := filepath.Join(os.TempDir(), "test-simple-html.pdf")
	defer os.Remove(outputPath)
	
	// Verify input exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Fatalf("Test fixture not found: %s", inputPath)
	}
	
	// Convert
	ctx := context.Background()
	opts := converter.Options{
		Quality: converter.ParseQuality("default"),
		Quiet:   false,
	}
	
	err := c.ConvertWithContext(ctx, inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("HTML conversion failed: %v", err)
	}
	
	// Verify output
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
	
	if info.Size() == 0 {
		t.Fatal("Output PDF is empty")
	}
	
	t.Logf("Generated PDF: %d bytes", info.Size())
}

func TestHTMLIntegration_ComplexHTML(t *testing.T) {
	c := converter.NewHTMLConverter()
	
	// Skip if no HTML converters available
	engines := c.GetAvailableEngines()
	if len(engines) == 0 {
		t.Skip("No HTML converters available (wkhtmltopdf or Chrome required)")
	}
	
	// Test file paths
	inputPath := filepath.Join("fixtures", "html", "complex.html")
	outputPath := filepath.Join(os.TempDir(), "test-complex-html.pdf")
	defer os.Remove(outputPath)
	
	// Verify input exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Fatalf("Test fixture not found: %s", inputPath)
	}
	
	// Convert with JavaScript delay
	ctx := context.Background()
	opts := converter.Options{
		Quality: converter.ParseQuality("default"),
		Quiet:   false,
		Extra: map[string]interface{}{
			"javascript_delay": 2000, // 2 seconds for JS execution
		},
	}
	
	err := c.ConvertWithContext(ctx, inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("Complex HTML conversion failed: %v", err)
	}
	
	// Verify output
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
	
	if info.Size() == 0 {
		t.Fatal("Output PDF is empty")
	}
	
	t.Logf("Generated PDF with JavaScript: %d bytes", info.Size())
}

func TestHTMLIntegration_PrintMediaCSS(t *testing.T) {
	c := converter.NewHTMLConverter()
	
	// Skip if no HTML converters available
	engines := c.GetAvailableEngines()
	if len(engines) == 0 {
		t.Skip("No HTML converters available (wkhtmltopdf or Chrome required)")
	}
	
	// Test file paths
	inputPath := filepath.Join("fixtures", "html", "print-media.html")
	outputPath := filepath.Join(os.TempDir(), "test-print-media.pdf")
	defer os.Remove(outputPath)
	
	// Verify input exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Fatalf("Test fixture not found: %s", inputPath)
	}
	
	// Convert (print media should be automatically applied)
	ctx := context.Background()
	opts := converter.Options{
		Quality: converter.ParseQuality("default"),
		Quiet:   false,
	}
	
	err := c.ConvertWithContext(ctx, inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("Print media CSS conversion failed: %v", err)
	}
	
	// Verify output
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
	
	if info.Size() == 0 {
		t.Fatal("Output PDF is empty")
	}
	
	t.Logf("Generated PDF with print media: %d bytes", info.Size())
}

func TestHTMLIntegration_QualityLevels(t *testing.T) {
	c := converter.NewHTMLConverter()
	
	// Skip if no HTML converters available
	engines := c.GetAvailableEngines()
	if len(engines) == 0 {
		t.Skip("No HTML converters available (wkhtmltopdf or Chrome required)")
	}
	
	inputPath := filepath.Join("fixtures", "html", "simple.html")
	
	qualities := []string{"screen", "default", "printer"}
	
	for _, quality := range qualities {
		t.Run(quality, func(t *testing.T) {
			outputPath := filepath.Join(os.TempDir(), "test-quality-"+quality+".pdf")
			defer os.Remove(outputPath)
			
			ctx := context.Background()
			opts := converter.Options{
				Quality: converter.ParseQuality(quality),
				Quiet:   true,
			}
			
			err := c.ConvertWithContext(ctx, inputPath, outputPath, opts)
			if err != nil {
				t.Fatalf("Conversion failed for quality %s: %v", quality, err)
			}
			
			info, err := os.Stat(outputPath)
			if err != nil {
				t.Fatalf("Output file not created for quality %s: %v", quality, err)
			}
			
			t.Logf("Quality %s: %d bytes", quality, info.Size())
		})
	}
}

func TestHTMLIntegration_PageSizeAndMargins(t *testing.T) {
	c := converter.NewHTMLConverter()
	
	// Skip if no HTML converters available
	engines := c.GetAvailableEngines()
	if len(engines) == 0 {
		t.Skip("No HTML converters available (wkhtmltopdf or Chrome required)")
	}
	
	inputPath := filepath.Join("fixtures", "html", "simple.html")
	
	tests := []struct {
		name       string
		pageSize   string
		marginTop  string
		marginLeft string
	}{
		{"A4 standard", "A4", "10mm", "10mm"},
		{"Letter large margins", "Letter", "20mm", "25mm"},
		{"Legal minimal margins", "Legal", "5mm", "5mm"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(os.TempDir(), "test-pagesize-"+tt.pageSize+".pdf")
			defer os.Remove(outputPath)
			
			ctx := context.Background()
			opts := converter.Options{
				Quality: converter.ParseQuality("default"),
				Quiet:   true,
				Extra: map[string]interface{}{
					"page_size":   tt.pageSize,
					"margin_top":  tt.marginTop,
					"margin_left": tt.marginLeft,
				},
			}
			
			err := c.ConvertWithContext(ctx, inputPath, outputPath, opts)
			if err != nil {
				t.Fatalf("Conversion failed for %s: %v", tt.name, err)
			}
			
			info, err := os.Stat(outputPath)
			if err != nil {
				t.Fatalf("Output file not created for %s: %v", tt.name, err)
			}
			
			t.Logf("%s: %d bytes", tt.name, info.Size())
		})
	}
}

func TestHTMLIntegration_Orientation(t *testing.T) {
	c := converter.NewHTMLConverter()
	
	// Skip if no HTML converters available
	engines := c.GetAvailableEngines()
	if len(engines) == 0 {
		t.Skip("No HTML converters available (wkhtmltopdf or Chrome required)")
	}
	
	inputPath := filepath.Join("fixtures", "html", "simple.html")
	
	orientations := []string{"portrait", "landscape"}
	
	for _, orientation := range orientations {
		t.Run(orientation, func(t *testing.T) {
			outputPath := filepath.Join(os.TempDir(), "test-orientation-"+orientation+".pdf")
			defer os.Remove(outputPath)
			
			ctx := context.Background()
			opts := converter.Options{
				Quality: converter.ParseQuality("default"),
				Quiet:   true,
				Extra: map[string]interface{}{
					"orientation": orientation,
				},
			}
			
			err := c.ConvertWithContext(ctx, inputPath, outputPath, opts)
			if err != nil {
				t.Fatalf("Conversion failed for %s: %v", orientation, err)
			}
			
			info, err := os.Stat(outputPath)
			if err != nil {
				t.Fatalf("Output file not created for %s: %v", orientation, err)
			}
			
			t.Logf("Orientation %s: %d bytes", orientation, info.Size())
		})
	}
}

func TestHTMLIntegration_FallbackLogic(t *testing.T) {
	c := converter.NewHTMLConverter()
	
	// Skip if no HTML converters available
	engines := c.GetAvailableEngines()
	if len(engines) == 0 {
		t.Skip("No HTML converters available (wkhtmltopdf or Chrome required)")
	}
	
	t.Logf("Testing fallback with available engines: %v", engines)
	
	// The converter should automatically try wkhtmltopdf first, then Chrome
	// We just verify it works with at least one
	inputPath := filepath.Join("fixtures", "html", "simple.html")
	outputPath := filepath.Join(os.TempDir(), "test-fallback.pdf")
	defer os.Remove(outputPath)
	
	ctx := context.Background()
	opts := converter.Options{
		Quality: converter.ParseQuality("default"),
		Quiet:   false,
	}
	
	err := c.ConvertWithContext(ctx, inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("Fallback logic failed: %v", err)
	}
	
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
	
	if info.Size() == 0 {
		t.Fatal("Output PDF is empty")
	}
	
	t.Logf("Fallback conversion successful: %d bytes", info.Size())
}
