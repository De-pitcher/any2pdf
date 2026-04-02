package converter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPandocConverter_Name(t *testing.T) {
	converter := NewPandocConverter()
	if converter.Name() != "pandoc" {
		t.Errorf("expected name 'pandoc', got '%s'", converter.Name())
	}
}

func TestPandocConverter_SetTimeout(t *testing.T) {
	converter := NewPandocConverter()
	
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

func TestPandocConverter_Metadata(t *testing.T) {
	converter := NewPandocConverter()
	
	// Initially should be empty
	if len(converter.metadata) != 0 {
		t.Error("metadata should be empty initially")
	}

	// Set metadata
	converter.SetMetadata("title", "Test Document")
	converter.SetMetadata("author", "Test Author")
	
	if converter.metadata["title"] != "Test Document" {
		t.Error("metadata title not set correctly")
	}
	if converter.metadata["author"] != "Test Author" {
		t.Error("metadata author not set correctly")
	}

	// Clear metadata
	converter.ClearMetadata()
	if len(converter.metadata) != 0 {
		t.Error("metadata should be empty after clear")
	}
}

func TestPandocConverter_IsSupportedFormat(t *testing.T) {
	converter := NewPandocConverter()

	tests := []struct {
		path     string
		expected bool
	}{
		{"document.txt", true},
		{"document.md", true},
		{"document.markdown", true},
		{"document.MD", true},  // Case insensitive
		{"document.rst", true},
		{"document.org", true},
		{"document.textile", true},
		{"document.mediawiki", true},
		{"document.pdf", false},
		{"document.docx", false},
		{"image.jpg", false},
		{"noextension", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := converter.IsSupportedFormat(tt.path)
			if result != tt.expected {
				t.Errorf("for %s: expected %v, got %v", tt.path, tt.expected, result)
			}
		})
	}
}

func TestPandocConverter_GetPreferredLatexEngine(t *testing.T) {
	converter := NewPandocConverter()
	engine := converter.getPreferredLatexEngine()
	
	// Should return either "xelatex" or "pdflatex"
	if engine != "xelatex" && engine != "pdflatex" {
		t.Errorf("unexpected LaTeX engine: %s", engine)
	}
}

func TestPandocConverter_GetMarginForQuality(t *testing.T) {
	converter := NewPandocConverter()

	tests := []struct {
		quality  Quality
		expected string
	}{
		{Screen, "0.75in"},
		{Default, "1in"},
		{Printer, "1.2in"},
	}

	for _, tt := range tests {
		t.Run(tt.quality.String(), func(t *testing.T) {
			result := converter.getMarginForQuality(tt.quality)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPandocConverter_GetFontSizeForQuality(t *testing.T) {
	converter := NewPandocConverter()

	tests := []struct {
		quality  Quality
		expected string
	}{
		{Screen, "10pt"},
		{Default, "11pt"},
		{Printer, "12pt"},
	}

	for _, tt := range tests {
		t.Run(tt.quality.String(), func(t *testing.T) {
			result := converter.getFontSizeForQuality(tt.quality)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPandocConverter_BuildArgs(t *testing.T) {
	converter := NewPandocConverter()

	tests := []struct {
		name           string
		inputPath      string
		outputPath     string
		quality        Quality
		metadata       map[string]string
		expectedMargin string
		expectedFont   string
	}{
		{
			name:           "Screen quality",
			inputPath:      "input.md",
			outputPath:     "output.pdf",
			quality:        Screen,
			expectedMargin: "0.75in",
			expectedFont:   "10pt",
		},
		{
			name:           "Default quality",
			inputPath:      "input.txt",
			outputPath:     "output.pdf",
			quality:        Default,
			expectedMargin: "1in",
			expectedFont:   "11pt",
		},
		{
			name:           "Printer quality",
			inputPath:      "input.md",
			outputPath:     "output.pdf",
			quality:        Printer,
			expectedMargin: "1.2in",
			expectedFont:   "12pt",
		},
		{
			name:           "With metadata",
			inputPath:      "input.md",
			outputPath:     "output.pdf",
			quality:        Default,
			metadata:       map[string]string{"title": "Test Doc"},
			expectedMargin: "1in",
			expectedFont:   "11pt",
		},
		{
			name:           "Stdin input",
			inputPath:      "-",
			outputPath:     "output.pdf",
			quality:        Default,
			expectedMargin: "1in",
			expectedFont:   "11pt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set metadata if provided
			if tt.metadata != nil {
				for k, v := range tt.metadata {
					converter.SetMetadata(k, v)
				}
			}

			opts := Options{Quality: tt.quality}
			args := converter.buildArgs(tt.inputPath, tt.outputPath, opts)

			// Check input path (unless stdin)
			if tt.inputPath != "-" {
				found := false
				for _, arg := range args {
					if arg == tt.inputPath {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("input path %s not found in args", tt.inputPath)
				}
			}

			// Check output path
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

			// Check margin setting
			marginFound := false
			for i := range args {
				if i+1 < len(args) && strings.Contains(args[i+1], "geometry:margin="+tt.expectedMargin) {
					marginFound = true
					break
				}
			}
			if !marginFound {
				t.Errorf("expected margin %s not found in args. Args: %v", tt.expectedMargin, args)
			}

			// Check font size
			fontFound := false
			for i := range args {
				if i+1 < len(args) && strings.Contains(args[i+1], "fontsize="+tt.expectedFont) {
					fontFound = true
					break
				}
			}
			if !fontFound {
				t.Errorf("expected font size %s not found in args. Args: %v", tt.expectedFont, args)
			}

			// Check PDF engine is present
			engineFound := false
			for _, arg := range args {
				if strings.HasPrefix(arg, "--pdf-engine=") {
					engineFound = true
					break
				}
			}
			if !engineFound {
				t.Error("PDF engine not specified in args")
			}

			// Check highlighting
			highlightFound := false
			for _, arg := range args {
				if arg == "--highlight-style=tango" {
					highlightFound = true
					break
				}
			}
			if !highlightFound {
				t.Error("syntax highlighting not found in args")
			}

			// Check metadata if provided
			if tt.metadata != nil {
				for key, value := range tt.metadata {
					metadataFound := false
					expectedMetadata := fmt.Sprintf("%s=%s", key, value)
					for i, arg := range args {
						if arg == "--metadata" && i+1 < len(args) && args[i+1] == expectedMetadata {
							metadataFound = true
							break
						}
					}
					if !metadataFound {
						t.Errorf("metadata %s not found in args", expectedMetadata)
					}
				}
			}

			// Clear metadata for next test
			converter.ClearMetadata()
		})
	}
}

func TestPandocConverter_GetInstallInstructions(t *testing.T) {
	converter := NewPandocConverter()
	instructions := converter.GetInstallInstructions()
	
	if instructions == "" {
		t.Error("install instructions should not be empty")
	}

	// Check that instructions mention common platforms and tools
	required := []string{"pandoc", "macOS", "Linux", "Windows", "LaTeX", "brew", "apt-get", "choco"}
	for _, req := range required {
		if !strings.Contains(instructions, req) {
			t.Errorf("install instructions should mention %s", req)
		}
	}
}

// Integration test - only runs if pandoc is installed
func TestPandocConverter_Convert_Integration(t *testing.T) {
	converter := NewPandocConverter()
	
	// Skip if pandoc or LaTeX is not installed
	if err := converter.checkDependencies(); err != nil {
		t.Skip("pandoc or LaTeX not installed, skipping integration test")
	}

	// Create a test markdown file
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "test.md")
	outputPath := filepath.Join(tmpDir, "output.pdf")

	testContent := `# Test Document

This is a **test** document with:

- List item 1
- List item 2

## Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```" + `

End of document.
`

	if err := os.WriteFile(inputPath, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test conversion with metadata
	converter.SetMetadata("title", "Test Document")
	converter.SetMetadata("author", "Test Author")

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

func TestPandocConverter_ConvertWithContext_Timeout(t *testing.T) {
	converter := NewPandocConverter()
	
	// Skip if pandoc is not installed
	if err := converter.CheckAvailable(); err != nil {
		t.Skip("pandoc not installed, skipping integration test")
	}

	// Create a context that times out immediately
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for context to timeout
	time.Sleep(2 * time.Millisecond)

	// Create a temporary test file
	tmpfile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("# Test\n")
	tmpfile.Close()

	// Try to convert with expired context
	err = converter.ConvertWithContext(ctx, tmpfile.Name(), "output.pdf", Options{})
	
	if err == nil {
		t.Error("expected timeout error, got nil")
		os.Remove("output.pdf")
	}
}
