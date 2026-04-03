package converter

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHTMLConverter_Name(t *testing.T) {
	c := NewHTMLConverter()
	if c.Name() != "html" {
		t.Errorf("Expected name 'html', got '%s'", c.Name())
	}
}

func TestHTMLConverter_SetTimeout(t *testing.T) {
	c := NewHTMLConverter()
	newTimeout := 30 * time.Second
	c.SetTimeout(newTimeout)
	
	if c.timeout != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, c.timeout)
	}
}

func TestHTMLConverter_SetPreserveTemp(t *testing.T) {
	c := NewHTMLConverter()
	c.SetPreserveTemp(true)
	
	if !c.preserveTemp {
		t.Error("Expected preserveTemp to be true")
	}
}

func TestHTMLConverter_IsSupportedFormat(t *testing.T) {
	c := NewHTMLConverter()
	
	tests := []struct {
		filename string
		expected bool
	}{
		{"page.html", true},
		{"page.htm", true},
		{"page.HTML", true},
		{"page.xhtml", true},
		{"page.xml", true},
		{"page.mhtml", true},
		{"page.pdf", false},
		{"page.txt", false},
		{"page.doc", false},
		{"noextension", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := c.IsSupportedFormat(tt.filename)
			if result != tt.expected {
				t.Errorf("IsSupportedFormat(%s) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestHTMLConverter_FindBinaries(t *testing.T) {
	c := NewHTMLConverter()
	
	// At least one should be found on most systems (or none on minimal systems)
	// Just verify the detection ran without panicking
	if c.wkhtmltopdfPath != "" {
		t.Logf("Found wkhtmltopdf at: %s", c.wkhtmltopdfPath)
	}
	if c.chromePath != "" {
		t.Logf("Found Chrome at: %s", c.chromePath)
	}
	
	// Verify availability flags match paths
	if (c.wkhtmltopdfPath != "") != c.wkAvailable {
		t.Error("wkAvailable flag doesn't match wkhtmltopdfPath state")
	}
	if (c.chromePath != "") != c.chromeAvailable {
		t.Error("chromeAvailable flag doesn't match chromePath state")
	}
}

func TestHTMLConverter_BuildWkhtmltopdfArgs(t *testing.T) {
	c := NewHTMLConverter()
	
	tests := []struct {
		name     string
		input    string
		output   string
		opts     HTMLOptions
		quality  string
		contains []string
	}{
		{
			name:    "default options",
			input:   "input.html",
			output:  "output.pdf",
			opts:    HTMLOptions{},
			quality: "default",
			contains: []string{
				"--enable-local-file-access",
				"--print-media-type",
				"--page-size", "A4",
				"--margin-top", "10mm",
				"--dpi", "150",
				"input.html",
				"output.pdf",
			},
		},
		{
			name:   "custom page size and orientation",
			input:  "input.html",
			output: "output.pdf",
			opts: HTMLOptions{
				PageSize:    "Letter",
				Orientation: "landscape",
			},
			quality: "default",
			contains: []string{
				"--page-size", "Letter",
				"--orientation", "Landscape",
			},
		},
		{
			name:   "printer quality",
			input:  "input.html",
			output: "output.pdf",
			opts:   HTMLOptions{},
			quality: "printer",
			contains: []string{
				"--dpi", "300",
				"--disable-smart-shrinking",
			},
		},
		{
			name:   "screen quality",
			input:  "input.html",
			output: "output.pdf",
			opts:   HTMLOptions{},
			quality: "screen",
			contains: []string{
				"--dpi", "96",
				"--enable-smart-shrinking",
			},
		},
		{
			name:   "with JavaScript delay",
			input:  "input.html",
			output: "output.pdf",
			opts: HTMLOptions{
				JavaScriptDelay: 2000,
			},
			quality: "default",
			contains: []string{
				"--javascript-delay", "2000",
			},
		},
		{
			name:   "with zoom",
			input:  "input.html",
			output: "output.pdf",
			opts: HTMLOptions{
				Zoom: 1.5,
			},
			quality: "default",
			contains: []string{
				"--zoom", "1.50",
			},
		},
		{
			name:   "with TOC",
			input:  "input.html",
			output: "output.pdf",
			opts: HTMLOptions{
				EnableTOC: true,
			},
			quality: "default",
			contains: []string{
				"toc",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := c.buildWkhtmltopdfArgs(tt.input, tt.output, tt.opts, tt.quality)
			argsStr := strings.Join(args, " ")
			
			for _, expected := range tt.contains {
				if !strings.Contains(argsStr, expected) {
					t.Errorf("Expected args to contain '%s', got: %v", expected, argsStr)
				}
			}
		})
	}
}

func TestHTMLConverter_BuildChromeArgs(t *testing.T) {
	c := NewHTMLConverter()
	
	tests := []struct {
		name     string
		input    string
		output   string
		opts     HTMLOptions
		quality  string
		contains []string
	}{
		{
			name:    "default options",
			input:   "input.html",
			output:  "output.pdf",
			opts:    HTMLOptions{},
			quality: "default",
			contains: []string{
				"--headless",
				"--disable-gpu",
				"--print-to-pdf=output.pdf",
				"--print-to-pdf-no-header",
				"--print-to-pdf-paper-width=210",
				"--print-to-pdf-paper-height=297",
				"--virtual-time-budget=5000",
				"--print-to-pdf-scale=0.95",
			},
		},
		{
			name:   "Letter page size",
			input:  "input.html",
			output: "output.pdf",
			opts: HTMLOptions{
				PageSize: "Letter",
			},
			quality: "default",
			contains: []string{
				"--print-to-pdf-paper-width=216",
				"--print-to-pdf-paper-height=279",
			},
		},
		{
			name:   "Legal page size",
			input:  "input.html",
			output: "output.pdf",
			opts: HTMLOptions{
				PageSize: "Legal",
			},
			quality: "default",
			contains: []string{
				"--print-to-pdf-paper-width=216",
				"--print-to-pdf-paper-height=356",
			},
		},
		{
			name:   "printer quality",
			input:  "input.html",
			output: "output.pdf",
			opts:   HTMLOptions{},
			quality: "printer",
			contains: []string{
				"--print-to-pdf-scale=1.0",
			},
		},
		{
			name:   "screen quality",
			input:  "input.html",
			output: "output.pdf",
			opts:   HTMLOptions{},
			quality: "screen",
			contains: []string{
				"--print-to-pdf-scale=0.85",
			},
		},
		{
			name:   "with custom JS delay",
			input:  "input.html",
			output: "output.pdf",
			opts: HTMLOptions{
				JavaScriptDelay: 10000,
			},
			quality: "default",
			contains: []string{
				"--virtual-time-budget=10000",
			},
		},
		{
			name:    "URL input",
			input:   "https://example.com",
			output:  "output.pdf",
			opts:    HTMLOptions{},
			quality: "default",
			contains: []string{
				"https://example.com",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := c.buildChromeArgs(tt.input, tt.output, tt.opts, tt.quality)
			argsStr := strings.Join(args, " ")
			
			for _, expected := range tt.contains {
				if !strings.Contains(argsStr, expected) {
					t.Errorf("Expected args to contain '%s', got: %v", expected, argsStr)
				}
			}
		})
	}
}

func TestHTMLConverter_ExtractHTMLOptions(t *testing.T) {
	c := NewHTMLConverter()
	
	tests := []struct {
		name     string
		opts     Options
		expected HTMLOptions
	}{
		{
			name: "empty options",
			opts: Options{},
			expected: HTMLOptions{
				PageSize:    "A4",
				Orientation: "portrait",
			},
		},
		{
			name: "with page size and orientation",
			opts: Options{
				Extra: map[string]interface{}{
					"page_size":    "Letter",
					"orientation":  "landscape",
				},
			},
			expected: HTMLOptions{
				PageSize:    "Letter",
				Orientation: "landscape",
			},
		},
		{
			name: "with margins",
			opts: Options{
				Extra: map[string]interface{}{
					"margin_top":    "20mm",
					"margin_bottom": "20mm",
					"margin_left":   "15mm",
					"margin_right":  "15mm",
				},
			},
			expected: HTMLOptions{
				PageSize:     "A4",
				Orientation:  "portrait",
				MarginTop:    "20mm",
				MarginBottom: "20mm",
				MarginLeft:   "15mm",
				MarginRight:  "15mm",
			},
		},
		{
			name: "with JavaScript delay and zoom",
			opts: Options{
				Extra: map[string]interface{}{
					"javascript_delay": 3000,
					"zoom":            1.25,
				},
			},
			expected: HTMLOptions{
				PageSize:        "A4",
				Orientation:     "portrait",
				JavaScriptDelay: 3000,
				Zoom:            1.25,
			},
		},
		{
			name: "with TOC",
			opts: Options{
				Extra: map[string]interface{}{
					"enable_toc": true,
				},
			},
			expected: HTMLOptions{
				PageSize:    "A4",
				Orientation: "portrait",
				EnableTOC:   true,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.extractHTMLOptions(tt.opts)
			
			if result.PageSize != tt.expected.PageSize {
				t.Errorf("PageSize = %s, expected %s", result.PageSize, tt.expected.PageSize)
			}
			if result.Orientation != tt.expected.Orientation {
				t.Errorf("Orientation = %s, expected %s", result.Orientation, tt.expected.Orientation)
			}
			if result.MarginTop != tt.expected.MarginTop {
				t.Errorf("MarginTop = %s, expected %s", result.MarginTop, tt.expected.MarginTop)
			}
			if result.JavaScriptDelay != tt.expected.JavaScriptDelay {
				t.Errorf("JavaScriptDelay = %d, expected %d", result.JavaScriptDelay, tt.expected.JavaScriptDelay)
			}
			if result.Zoom != tt.expected.Zoom {
				t.Errorf("Zoom = %f, expected %f", result.Zoom, tt.expected.Zoom)
			}
			if result.EnableTOC != tt.expected.EnableTOC {
				t.Errorf("EnableTOC = %v, expected %v", result.EnableTOC, tt.expected.EnableTOC)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://example.com", true},
		{"http://example.com", true},
		{"file:///path/to/file.html", true},
		{"ftp://example.com", true},
		{"local-file.html", false},
		{"/absolute/path/file.html", false},
		{"relative/path/file.html", false},
		{"", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isURL(tt.input)
			if result != tt.expected {
				t.Errorf("isURL(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHTMLConverter_GetInstallInstructions(t *testing.T) {
	c := NewHTMLConverter()
	instructions := c.getInstallInstructions()
	
	if instructions == "" {
		t.Error("Expected non-empty install instructions")
	}
	
	// Should mention both wkhtmltopdf and Chrome
	if !strings.Contains(instructions, "wkhtmltopdf") {
		t.Error("Expected instructions to mention wkhtmltopdf")
	}
	if !strings.Contains(instructions, "Chrome") {
		t.Error("Expected instructions to mention Chrome")
	}
}

func TestHTMLConverter_ValidatePDFCreated(t *testing.T) {
	c := NewHTMLConverter()
	
	tests := []struct {
		name      string
		setup     func() string
		expectErr bool
		errMsg    string
	}{
		{
			name: "Non-existent file",
			setup: func() string {
				return "nonexistent.pdf"
			},
			expectErr: true,
			errMsg:    "output PDF not created",
		},
		{
			name: "Empty file",
			setup: func() string {
				f, _ := os.CreateTemp("", "empty-*.pdf")
				path := f.Name()
				f.Close()
				return path
			},
			expectErr: true,
			errMsg:    "output PDF is empty",
		},
		{
			name: "Invalid PDF header",
			setup: func() string {
				f, _ := os.CreateTemp("", "invalid-*.pdf")
				path := f.Name()
				f.WriteString("NOT A PDF")
				f.Close()
				return path
			},
			expectErr: true,
			errMsg:    "invalid PDF header",
		},
		{
			name: "Valid PDF",
			setup: func() string {
				f, _ := os.CreateTemp("", "valid-*.pdf")
				path := f.Name()
				f.WriteString("%PDF-1.4\n")
				f.WriteString("dummy content")
				f.Close()
				return path
			},
			expectErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			if path != "nonexistent.pdf" {
				defer os.Remove(path)
			}
			
			err := c.validatePDFCreated(path)
			
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestHTMLConverter_GetAvailableEngines(t *testing.T) {
	c := NewHTMLConverter()
	engines := c.GetAvailableEngines()
	
	// Should return a slice (possibly empty if no tools installed)
	if engines == nil {
		t.Error("Expected non-nil engines slice")
	}
	
	// If we found any engines, they should be valid
	for _, engine := range engines {
		if engine != "wkhtmltopdf" && engine != "chrome" {
			t.Errorf("Unexpected engine: %s", engine)
		}
	}
	
	t.Logf("Available engines: %v", engines)
}

func TestHTMLConverter_ConvertFromReader(t *testing.T) {
	c := NewHTMLConverter()
	
	// Skip if no converters available
	if !c.wkAvailable && !c.chromeAvailable {
		t.Skip("No HTML converters available")
	}
	
	// Create simple HTML content
	htmlContent := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body><h1>Test from reader</h1></body>
</html>`
	
	reader := strings.NewReader(htmlContent)
	
	// Create temp output file
	output := filepath.Join(os.TempDir(), "test-from-reader.pdf")
	defer os.Remove(output)
	
	ctx := context.Background()
	opts := Options{Quality: ParseQuality("default")}
	
	err := c.ConvertFromReader(ctx, reader, output, opts)
	if err != nil {
		t.Errorf("ConvertFromReader failed: %v", err)
	}
	
	// Verify PDF was created
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("Output PDF was not created")
	}
}
