package converter

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestImageConverter_Name(t *testing.T) {
	converter := NewImageConverter()
	if converter.Name() != "img2pdf" {
		t.Errorf("expected name 'img2pdf', got '%s'", converter.Name())
	}
}

func TestImageConverter_SetTimeout(t *testing.T) {
	converter := NewImageConverter()
	
	// Default timeout should be 5 minutes
	if converter.timeout != 5*time.Minute {
		t.Errorf("expected default timeout of 5 minutes, got %v", converter.timeout)
	}

	// Set custom timeout
	customTimeout := 10 * time.Second
	converter.SetTimeout(customTimeout)
	
	if converter.timeout != customTimeout {
		t.Errorf("expected timeout %v, got %v", customTimeout, converter.timeout)
	}
}

func TestImageConverter_IsSupportedImageFormat(t *testing.T) {
	converter := NewImageConverter()

	tests := []struct {
		path     string
		expected bool
	}{
		{"image.jpg", true},
		{"image.jpeg", true},
		{"image.JPG", true},  // Case insensitive
		{"image.png", true},
		{"image.PNG", true},
		{"image.gif", true},
		{"image.bmp", true},
		{"image.tiff", true},
		{"image.tif", true},
		{"document.pdf", false},
		{"document.txt", false},
		{"video.mp4", false},
		{"noextension", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := converter.IsSupportedImageFormat(tt.path)
			if result != tt.expected {
				t.Errorf("for %s: expected %v, got %v", tt.path, tt.expected, result)
			}
		})
	}
}

func TestImageConverter_BuildArgs(t *testing.T) {
	converter := NewImageConverter()

	tests := []struct {
		name        string
		inputPaths  []string
		outputPath  string
		quality     Quality
		expectedDPI string
	}{
		{
			name:        "Single image, screen quality",
			inputPaths:  []string{"image.jpg"},
			outputPath:  "output.pdf",
			quality:     Screen,
			expectedDPI: "72",
		},
		{
			name:        "Single image, default quality",
			inputPaths:  []string{"image.png"},
			outputPath:  "output.pdf",
			quality:     Default,
			expectedDPI: "150",
		},
		{
			name:        "Single image, printer quality",
			inputPaths:  []string{"image.jpg"},
			outputPath:  "output.pdf",
			quality:     Printer,
			expectedDPI: "300",
		},
		{
			name:        "Multiple images",
			inputPaths:  []string{"img1.jpg", "img2.png", "img3.gif"},
			outputPath:  "combined.pdf",
			quality:     Default,
			expectedDPI: "150",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{Quality: tt.quality}
			args := converter.buildArgs(tt.inputPaths, tt.outputPath, opts)

			// Check that all input paths are in args
			for _, input := range tt.inputPaths {
				found := false
				for _, arg := range args {
					if arg == input {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("input path %s not found in args", input)
				}
			}

			// Check output path is in args
			foundOutput := false
			for i, arg := range args {
				if arg == "-o" && i+1 < len(args) && args[i+1] == tt.outputPath {
					foundOutput = true
					break
				}
			}
			if !foundOutput {
				t.Errorf("output path %s not found in args", tt.outputPath)
			}

			// Check DPI setting
			foundDPI := false
			for i, arg := range args {
				if arg == "--dpi" && i+1 < len(args) && args[i+1] == tt.expectedDPI {
					foundDPI = true
					break
				}
			}
			if !foundDPI {
				t.Errorf("expected DPI %s not found in args. Args: %v", tt.expectedDPI, args)
			}
		})
	}
}

func TestImageConverter_ConvertMultiple_EmptyInput(t *testing.T) {
	converter := NewImageConverter()
	
	err := converter.ConvertMultiple([]string{}, "output.pdf", Options{})
	
	if err == nil {
		t.Error("expected error for empty input paths, got nil")
	}
	
	if err.Error() != "no input images provided" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestImageConverter_ConvertWithContext_Timeout(t *testing.T) {
	// Skip if img2pdf is not installed
	converter := NewImageConverter()
	if err := converter.CheckAvailable(); err != nil {
		t.Skip("img2pdf not installed, skipping integration test")
	}

	// Create a context that times out immediately
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for context to timeout
	time.Sleep(2 * time.Millisecond)

	// Create a temporary test file
	tmpfile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Try to convert with expired context
	err = converter.ConvertWithContext(ctx, tmpfile.Name(), "output.pdf", Options{})
	
	if err == nil {
		t.Error("expected timeout error, got nil")
		os.Remove("output.pdf")
	}
}

func TestImageConverter_GetInstallInstructions(t *testing.T) {
	converter := NewImageConverter()
	instructions := converter.GetInstallInstructions()
	
	if instructions == "" {
		t.Error("install instructions should not be empty")
	}

	// Check that instructions mention common platforms
	platforms := []string{"Ubuntu", "macOS", "Windows", "pip"}
	for _, platform := range platforms {
		if !contains(instructions, platform) {
			t.Errorf("install instructions should mention %s", platform)
		}
	}
}

// Integration test - only runs if img2pdf is installed
func TestImageConverter_Convert_Integration(t *testing.T) {
	converter := NewImageConverter()
	
	// Skip if img2pdf is not installed
	if err := converter.CheckAvailable(); err != nil {
		t.Skip("img2pdf not installed, skipping integration test")
	}

	// Create a simple 1x1 PNG image for testing
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "test.png")
	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Create a minimal valid PNG file (1x1 white pixel)
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}

	if err := os.WriteFile(inputPath, pngData, 0644); err != nil {
		t.Fatal(err)
	}

	// Test conversion
	opts := Options{
		Quality: Default,
		Quiet:   true,
	}

	err := converter.Convert(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("output PDF file was not created")
	}

	// Verify it's a PDF (check magic bytes)
	pdfData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(pdfData) < 4 || string(pdfData[:4]) != "%PDF" {
		t.Error("output file does not appear to be a valid PDF")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
