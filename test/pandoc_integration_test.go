package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/De-pitcher/any2pdf/internal/converter"
)

// TestPandocConversion_Integration tests actual text/markdown to PDF conversion
func TestPandocConversion_Integration(t *testing.T) {
	pandocConverter := converter.NewPandocConverter()

	// Skip if pandoc or LaTeX is not installed
	if err := pandocConverter.CheckAvailable(); err != nil {
		t.Skip("pandoc not installed, skipping integration test")
	}

	// Create a test output directory
	outputDir := t.TempDir()

	tests := []struct {
		name       string
		fixture    string
		quality    converter.Quality
		metadata   map[string]string
		shouldPass bool
	}{
		{
			name:       "Markdown file - default quality",
			fixture:    "fixtures/markdown/sample.md",
			quality:    converter.Default,
			shouldPass: true,
		},
		{
			name:       "Text file - screen quality",
			fixture:    "fixtures/text/sample.txt",
			quality:    converter.Screen,
			shouldPass: true,
		},
		{
			name:       "Markdown file - printer quality with metadata",
			fixture:    "fixtures/markdown/sample.md",
			quality:    converter.Printer,
			metadata:   map[string]string{"title": "Sample Document", "author": "Test Author"},
			shouldPass: true,
		},
		{
			name:       "Text file - default quality",
			fixture:    "fixtures/text/sample.txt",
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

			// Set metadata if provided
			if tt.metadata != nil {
				for k, v := range tt.metadata {
					pandocConverter.SetMetadata(k, v)
				}
			}

			opts := converter.Options{
				Quality: tt.quality,
				Quiet:   true,
			}

			err := pandocConverter.Convert(tt.fixture, outputPath, opts)

			// Clear metadata for next test
			pandocConverter.ClearMetadata()

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

// TestMarkdownFeatures tests various markdown features
func TestMarkdownFeatures_Integration(t *testing.T) {
	pandocConverter := converter.NewPandocConverter()

	// Skip if pandoc or LaTeX is not installed
	if err := pandocConverter.CheckAvailable(); err != nil {
		t.Skip("pandoc not installed, skipping integration test")
	}

	// Create a test markdown file with various features
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "features.md")
	outputPath := filepath.Join(tmpDir, "features.pdf")

	testContent := `# Markdown Feature Test

## Text Formatting

This document tests various **markdown features** including *italic text*, 
***bold italic***, and ` + "`inline code`" + `.

## Lists

### Unordered List
- First item
- Second item
  - Nested item
  - Another nested item
- Third item

### Ordered List
1. First step
2. Second step
3. Third step

## Code Blocks

` + "```python\ndef hello():\n    print(\"Hello, World!\")\n```" + `

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, Go!\")\n}\n```" + `

## Links and Images

[Visit Example](https://example.com)

## Tables

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |

## Blockquotes

> This is a blockquote
> with multiple lines

## Horizontal Rule

---

## Conclusion

This tests the pandoc converter's ability to handle various markdown features.
`

	if err := os.WriteFile(inputPath, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Set metadata
	pandocConverter.SetMetadata("title", "Markdown Feature Test")
	pandocConverter.SetMetadata("author", "any2pdf Test Suite")

	opts := converter.Options{
		Quality: converter.Default,
		Quiet:   true,
	}

	err := pandocConverter.Convert(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	// Verify output file exists and is valid PDF
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) < 4 || string(data[:4]) != "%PDF" {
		t.Error("output file is not a valid PDF")
	}

	// Check file size is reasonable (should be more than 1KB for all that content)
	if len(data) < 1024 {
		t.Errorf("PDF seems too small (%d bytes), conversion may have failed", len(data))
	}
}
