package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/De-pitcher/any2pdf/internal/converter"
)

// TestImageConversion_Integration tests actual image to PDF conversion
func TestImageConversion_Integration(t *testing.T) {
	imageConverter := converter.NewImageConverter()

	// Skip if img2pdf is not installed
	if err := imageConverter.CheckAvailable(); err != nil {
		t.Skip("img2pdf not installed, skipping integration test")
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
			name:       "PNG image - default quality",
			fixture:    "fixtures/images/sample.png",
			quality:    converter.Default,
			shouldPass: true,
		},
		{
			name:       "JPG image - screen quality",
			fixture:    "fixtures/images/sample.jpg",
			quality:    converter.Screen,
			shouldPass: true,
		},
		{
			name:       "JPG image - printer quality",
			fixture:    "fixtures/images/sample.jpg",
			quality:    converter.Printer,
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

			err := imageConverter.Convert(tt.fixture, outputPath, opts)

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
			}
		})
	}
}

// TestMultipleImagesConversion tests combining multiple images into one PDF
func TestMultipleImagesConversion_Integration(t *testing.T) {
	imageConverter := converter.NewImageConverter()

	// Skip if img2pdf is not installed
	if err := imageConverter.CheckAvailable(); err != nil {
		t.Skip("img2pdf not installed, skipping integration test")
	}

	// Check if fixtures exist
	fixtures := []string{
		"fixtures/images/sample.png",
		"fixtures/images/sample.jpg",
	}

	for _, fixture := range fixtures {
		if _, err := os.Stat(fixture); os.IsNotExist(err) {
			t.Skipf("fixture file %s not found", fixture)
		}
	}

	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "combined.pdf")

	opts := converter.Options{
		Quality: converter.Default,
		Quiet:   true,
	}

	err := imageConverter.ConvertMultiple(fixtures, outputPath, opts)
	if err != nil {
		t.Fatalf("multi-image conversion failed: %v", err)
	}

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
}
