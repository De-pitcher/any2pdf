package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/De-pitcher/any2pdf/internal/converter"
)

// TestLibreOfficeConversion_Integration tests actual Office document conversion
func TestLibreOfficeConversion_Integration(t *testing.T) {
	libreOfficeConverter := converter.NewLibreOfficeConverter()

	// Skip if LibreOffice is not installed
	if err := libreOfficeConverter.CheckAvailable(); err != nil {
		t.Skip("LibreOffice not installed, skipping integration test")
	}

	// Create a test output directory
	outputDir := t.TempDir()

	tests := []struct {
		name       string
		fixture    string
		quality    converter.Quality
		shouldPass bool
	}{
		{
			name:       "Word document - default quality",
			fixture:    "fixtures/office/sample.docx",
			quality:    converter.Default,
			shouldPass: true,
		},
		{
			name:       "Excel spreadsheet - screen quality",
			fixture:    "fixtures/office/sample.xlsx",
			quality:    converter.Screen,
			shouldPass: true,
		},
		{
			name:       "PowerPoint presentation - printer quality",
			fixture:    "fixtures/office/sample.pptx",
			quality:    converter.Printer,
			shouldPass: true,
		},
		{
			name:       "OpenDocument text - default quality",
			fixture:    "fixtures/office/sample.odt",
			quality:    converter.Default,
			shouldPass: true,
		},
		{
			name:       "Rich Text Format - default quality",
			fixture:    "fixtures/office/sample.rtf",
			quality:    converter.Default,
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if fixture exists
			if _, err := os.Stat(tt.fixture); os.IsNotExist(err) {
				t.Skipf("fixture file %s not found", tt.fixture)
			}

			outputPath := filepath.Join(outputDir, tt.name+".pdf")

			opts := converter.Options{
				Quality: tt.quality,
				Quiet:   true,
			}

			err := libreOfficeConverter.Convert(tt.fixture, outputPath, opts)

			if tt.shouldPass && err != nil {
				t.Errorf("conversion failed: %v", err)
			}

			if tt.shouldPass {
				// Verify output exists
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Error("output PDF file was not created")
				}

				// Verify it's a PDF
				data, err := os.ReadFile(outputPath)
				if err != nil {
					t.Fatal(err)
				}
				if len(data) < 4 || string(data[:4]) != "%PDF" {
					t.Error("output file is not a valid PDF")
				}

				// Check file size is reasonable (should be more than 100 bytes)
				if len(data) < 100 {
					t.Errorf("PDF seems too small (%d bytes), conversion may have failed", len(data))
				}
			}
		})
	}
}

// TestBatchConversion_Integration tests batch conversion of multiple files
func TestBatchConversion_Integration(t *testing.T) {
	libreOfficeConverter := converter.NewLibreOfficeConverter()

	// Skip if LibreOffice is not installed
	if err := libreOfficeConverter.CheckAvailable(); err != nil {
		t.Skip("LibreOffice not installed, skipping integration test")
	}

	// Check if fixtures exist
	fixtures := []string{
		"fixtures/office/sample.docx",
		"fixtures/office/sample.xlsx",
	}

	for _, fixture := range fixtures {
		if _, err := os.Stat(fixture); os.IsNotExist(err) {
			t.Skipf("fixture file %s not found", fixture)
		}
	}

	outputDir := t.TempDir()

	opts := converter.Options{
		Quality: converter.Default,
		Quiet:   true,
	}

	outputs, err := libreOfficeConverter.ConvertBatch(nil, fixtures, outputDir, opts)
	if err != nil {
		t.Fatalf("batch conversion failed: %v", err)
	}

	// Verify we got the expected number of outputs
	if len(outputs) != len(fixtures) {
		t.Errorf("expected %d outputs, got %d", len(fixtures), len(outputs))
	}

	// Verify each output
	for _, outputPath := range outputs {
		// Verify output exists
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Errorf("output PDF file was not created: %s", outputPath)
			continue
		}

		// Verify it's a PDF
		data, err := os.ReadFile(outputPath)
		if err != nil {
			t.Errorf("failed to read output file: %v", err)
			continue
		}
		if len(data) < 4 || string(data[:4]) != "%PDF" {
			t.Errorf("output file is not a valid PDF: %s", outputPath)
		}
	}
}

// TestComplexDocumentFeatures tests handling of documents with complex features
func TestComplexDocumentFeatures_Integration(t *testing.T) {
	libreOfficeConverter := converter.NewLibreOfficeConverter()

	// Skip if LibreOffice is not installed
	if err := libreOfficeConverter.CheckAvailable(); err != nil {
		t.Skip("LibreOffice not installed, skipping integration test")
	}

	// Skip if fixture doesn't exist
	fixture := "fixtures/office/sample.docx"
	if _, err := os.Stat(fixture); os.IsNotExist(err) {
		t.Skipf("fixture file %s not found", fixture)
	}

	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "complex.pdf")

	opts := converter.Options{
		Quality: converter.Default,
		Quiet:   true,
	}

	err := libreOfficeConverter.Convert(fixture, outputPath, opts)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	// Verify PDF was created and is valid
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) < 4 || string(data[:4]) != "%PDF" {
		t.Fatal("output file is not a valid PDF")
	}

	// PDF should be reasonably sized (at least 1KB for a document with content)
	if len(data) < 1024 {
		t.Errorf("PDF seems suspiciously small (%d bytes)", len(data))
	}
}
