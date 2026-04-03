package converter

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// HTMLConverter converts HTML/web pages to PDF using wkhtmltopdf or Chrome/Chromium
type HTMLConverter struct {
	wkhtmltopdfPath string
	chromePath      string
	wkAvailable     bool
	chromeAvailable bool
	timeout         time.Duration
	preserveTemp    bool
}

// HTMLOptions contains HTML-specific conversion options
type HTMLOptions struct {
	PageSize        string // A4, Letter, Legal
	Orientation     string // portrait, landscape
	MarginTop       string // e.g., "10mm", "1in"
	MarginBottom    string
	MarginLeft      string
	MarginRight     string
	JavaScriptDelay int    // milliseconds
	Zoom            float64 // 0.75 to 2.0
	EnableTOC       bool
}

// NewHTMLConverter creates a new HTML converter with binary detection
func NewHTMLConverter() *HTMLConverter {
	c := &HTMLConverter{
		timeout: 10 * time.Second, // HTML pages may load external resources
	}
	c.detectBinaries()
	return c
}

// Name returns the converter name
func (c *HTMLConverter) Name() string {
	return "html"
}

// SetTimeout sets the conversion timeout
func (c *HTMLConverter) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetPreserveTemp sets whether to preserve temporary files
func (c *HTMLConverter) SetPreserveTemp(preserve bool) {
	c.preserveTemp = preserve
}

// IsSupportedFormat checks if the file format is supported
func (c *HTMLConverter) IsSupportedFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supported := []string{".html", ".htm", ".xhtml", ".xml", ".mhtml"}
	for _, s := range supported {
		if ext == s {
			return true
		}
	}
	return false
}

// detectBinaries finds available HTML-to-PDF conversion tools
func (c *HTMLConverter) detectBinaries() {
	c.wkhtmltopdfPath = c.findWkhtmltopdf()
	c.wkAvailable = c.wkhtmltopdfPath != ""
	
	c.chromePath = c.findChrome()
	c.chromeAvailable = c.chromePath != ""
}

// findWkhtmltopdf locates wkhtmltopdf binary
func (c *HTMLConverter) findWkhtmltopdf() string {
	// Try common installation paths
	candidates := []string{}
	
	switch runtime.GOOS {
	case "windows":
		candidates = []string{
			"wkhtmltopdf.exe",
			`C:\Program Files\wkhtmltopdf\bin\wkhtmltopdf.exe`,
			`C:\Program Files (x86)\wkhtmltopdf\bin\wkhtmltopdf.exe`,
		}
	case "darwin":
		candidates = []string{
			"wkhtmltopdf",
			"/usr/local/bin/wkhtmltopdf",
			"/opt/homebrew/bin/wkhtmltopdf",
		}
	case "linux":
		candidates = []string{
			"wkhtmltopdf",
			"/usr/bin/wkhtmltopdf",
			"/usr/local/bin/wkhtmltopdf",
		}
	}
	
	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			return path
		}
	}
	
	return ""
}

// findChrome locates Chrome/Chromium binary
func (c *HTMLConverter) findChrome() string {
	candidates := []string{}
	
	switch runtime.GOOS {
	case "windows":
		candidates = []string{
			"chrome.exe",
			"chromium.exe",
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files\Chromium\Application\chrome.exe`,
		}
	case "darwin":
		candidates = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
	case "linux":
		candidates = []string{
			"google-chrome",
			"google-chrome-stable",
			"chromium",
			"chromium-browser",
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
		}
	}
	
	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			return path
		}
	}
	
	return ""
}

// Convert converts HTML to PDF (interface implementation)
func (c *HTMLConverter) Convert(inputPath, outputPath string, opts Options) error {
	ctx := context.Background()
	return c.ConvertWithContext(ctx, inputPath, outputPath, opts)
}

// CheckAvailable verifies that at least one HTML converter is available
func (c *HTMLConverter) CheckAvailable() error {
	if !c.wkAvailable && !c.chromeAvailable {
		return &ConverterNotFoundError{
			Converter: "html",
			Message:   c.getInstallInstructions(),
		}
	}
	return nil
}

// ConvertWithContext converts HTML to PDF with smart fallback
func (c *HTMLConverter) ConvertWithContext(ctx context.Context, input, output string, opts Options) error {
	// Create context with timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	
	// Validate input
	if !c.IsSupportedFormat(input) && !isURL(input) {
		return fmt.Errorf("unsupported format: %s", input)
	}
	
	// Extract HTML-specific options
	htmlOpts := c.extractHTMLOptions(opts)
	
	// Try wkhtmltopdf first (faster for simple pages)
	if c.wkAvailable {
		err := c.convertWithWkhtmltopdf(ctxTimeout, input, output, htmlOpts, opts.Quality.String())
		if err == nil {
			return c.validatePDFCreated(output)
		}
		// Log error but continue to Chrome fallback
		if !opts.Quiet {
			fmt.Fprintf(os.Stderr, "wkhtmltopdf failed: %v, trying Chrome...\n", err)
		}
	}
	
	// Fallback to Chrome/Chromium
	if c.chromeAvailable {
		err := c.convertWithChrome(ctxTimeout, input, output, htmlOpts, opts.Quality.String())
		if err == nil {
			return c.validatePDFCreated(output)
		}
		return fmt.Errorf("chrome conversion failed: %w", err)
	}
	
	// No converter available
	return &ConverterNotFoundError{
		Converter: "html",
		Message:   c.getInstallInstructions(),
	}
}

// convertWithWkhtmltopdf converts using wkhtmltopdf
func (c *HTMLConverter) convertWithWkhtmltopdf(ctx context.Context, input, output string, opts HTMLOptions, quality string) error {
	args := c.buildWkhtmltopdfArgs(input, output, opts, quality)
	
	cmd := exec.CommandContext(ctx, c.wkhtmltopdfPath, args...)
	
	// Capture stderr for error reporting
	var stderr strings.Builder
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("wkhtmltopdf error: %s", stderr.String())
		}
		return fmt.Errorf("wkhtmltopdf failed: %w", err)
	}
	
	return nil
}

// convertWithChrome converts using Chrome/Chromium headless
func (c *HTMLConverter) convertWithChrome(ctx context.Context, input, output string, opts HTMLOptions, quality string) error {
	args := c.buildChromeArgs(input, output, opts, quality)
	
	cmd := exec.CommandContext(ctx, c.chromePath, args...)
	
	// Capture stderr for error reporting
	var stderr strings.Builder
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("chrome error: %s", stderr.String())
		}
		return fmt.Errorf("chrome failed: %w", err)
	}
	
	return nil
}

// buildWkhtmltopdfArgs constructs command arguments for wkhtmltopdf
func (c *HTMLConverter) buildWkhtmltopdfArgs(input, output string, opts HTMLOptions, quality string) []string {
	args := []string{
		"--enable-local-file-access",
		"--print-media-type",
	}
	
	// Page size
	pageSize := opts.PageSize
	if pageSize == "" {
		pageSize = "A4"
	}
	args = append(args, "--page-size", pageSize)
	
	// Orientation
	orientation := opts.Orientation
	if orientation == "" {
		orientation = "portrait"
	}
	if orientation == "landscape" {
		args = append(args, "--orientation", "Landscape")
	}
	
	// Margins
	marginTop := opts.MarginTop
	if marginTop == "" {
		marginTop = "10mm"
	}
	marginBottom := opts.MarginBottom
	if marginBottom == "" {
		marginBottom = "10mm"
	}
	marginLeft := opts.MarginLeft
	if marginLeft == "" {
		marginLeft = "10mm"
	}
	marginRight := opts.MarginRight
	if marginRight == "" {
		marginRight = "10mm"
	}
	
	args = append(args,
		"--margin-top", marginTop,
		"--margin-bottom", marginBottom,
		"--margin-left", marginLeft,
		"--margin-right", marginRight,
	)
	
	// Quality-specific options
	switch quality {
	case "screen":
		args = append(args, "--dpi", "96")
		args = append(args, "--enable-smart-shrinking")
	case "printer":
		args = append(args, "--dpi", "300")
		args = append(args, "--disable-smart-shrinking")
	default: // "default"
		args = append(args, "--dpi", "150")
		args = append(args, "--disable-smart-shrinking")
	}
	
	// JavaScript delay
	if opts.JavaScriptDelay > 0 {
		args = append(args, "--javascript-delay", fmt.Sprintf("%d", opts.JavaScriptDelay))
	}
	
	// Zoom
	if opts.Zoom > 0 {
		args = append(args, "--zoom", fmt.Sprintf("%.2f", opts.Zoom))
	}
	
	// Table of contents
	if opts.EnableTOC {
		args = append(args, "toc")
	}
	
	// Input and output
	args = append(args, input, output)
	
	return args
}

// buildChromeArgs constructs command arguments for Chrome headless
func (c *HTMLConverter) buildChromeArgs(input, output string, opts HTMLOptions, quality string) []string {
	args := []string{
		"--headless",
		"--disable-gpu",
		"--no-sandbox", // Required for some environments
		"--disable-dev-shm-usage",
		"--print-to-pdf=" + output,
		"--print-to-pdf-no-header",
	}
	
	// Page size - Chrome uses different format
	paperSize := opts.PageSize
	if paperSize == "" {
		paperSize = "A4"
	}
	
	// Convert to Chrome's paper size format
	switch strings.ToUpper(paperSize) {
	case "A4":
		// Chrome uses 210x297mm for A4
		args = append(args, "--print-to-pdf-paper-width=210", "--print-to-pdf-paper-height=297")
	case "LETTER":
		// 8.5x11 inches = 216x279mm
		args = append(args, "--print-to-pdf-paper-width=216", "--print-to-pdf-paper-height=279")
	case "LEGAL":
		// 8.5x14 inches = 216x356mm
		args = append(args, "--print-to-pdf-paper-width=216", "--print-to-pdf-paper-height=356")
	}
	
	// Orientation via page dimensions (swap width/height for landscape)
	if opts.Orientation == "landscape" {
		// Need to swap width and height in existing args
		// This is complex, so we'll handle it separately if needed
	}
	
	// Virtual time budget for JavaScript (milliseconds)
	jsDelay := opts.JavaScriptDelay
	if jsDelay == 0 {
		jsDelay = 5000 // Default 5 seconds
	}
	args = append(args, fmt.Sprintf("--virtual-time-budget=%d", jsDelay))
	
	// Quality-specific options
	switch quality {
	case "screen":
		args = append(args, "--print-to-pdf-scale=0.85")
	case "printer":
		args = append(args, "--print-to-pdf-scale=1.0")
	default: // "default"
		args = append(args, "--print-to-pdf-scale=0.95")
	}
	
	// Handle URL vs file path
	if isURL(input) {
		args = append(args, input)
	} else {
		// Convert to file:// URL
		absPath, err := filepath.Abs(input)
		if err == nil {
			fileURL := "file://" + filepath.ToSlash(absPath)
			args = append(args, fileURL)
		} else {
			args = append(args, input)
		}
	}
	
	return args
}

// extractHTMLOptions extracts HTML-specific options from generic Options
func (c *HTMLConverter) extractHTMLOptions(opts Options) HTMLOptions {
	htmlOpts := HTMLOptions{
		PageSize:    "A4",
		Orientation: "portrait",
	}
	
	// Extract from opts.Extra map if present
	if opts.Extra != nil {
		if pageSize, ok := opts.Extra["page_size"].(string); ok {
			htmlOpts.PageSize = pageSize
		}
		if orientation, ok := opts.Extra["orientation"].(string); ok {
			htmlOpts.Orientation = orientation
		}
		if marginTop, ok := opts.Extra["margin_top"].(string); ok {
			htmlOpts.MarginTop = marginTop
		}
		if marginBottom, ok := opts.Extra["margin_bottom"].(string); ok {
			htmlOpts.MarginBottom = marginBottom
		}
		if marginLeft, ok := opts.Extra["margin_left"].(string); ok {
			htmlOpts.MarginLeft = marginLeft
		}
		if marginRight, ok := opts.Extra["margin_right"].(string); ok {
			htmlOpts.MarginRight = marginRight
		}
		if jsDelay, ok := opts.Extra["javascript_delay"].(int); ok {
			htmlOpts.JavaScriptDelay = jsDelay
		}
		if zoom, ok := opts.Extra["zoom"].(float64); ok {
			htmlOpts.Zoom = zoom
		}
		if enableTOC, ok := opts.Extra["enable_toc"].(bool); ok {
			htmlOpts.EnableTOC = enableTOC
		}
	}
	
	return htmlOpts
}

// validatePDFCreated checks if the PDF was successfully created
func (c *HTMLConverter) validatePDFCreated(output string) error {
	info, err := os.Stat(output)
	if err != nil {
		return fmt.Errorf("output PDF not created: %w", err)
	}
	
	if info.Size() == 0 {
		return fmt.Errorf("output PDF is empty")
	}
	
	// Check PDF header
	file, err := os.Open(output)
	if err != nil {
		return fmt.Errorf("cannot open output PDF: %w", err)
	}
	defer file.Close()
	
	header := make([]byte, 4)
	_, err = io.ReadFull(file, header)
	if err != nil {
		return fmt.Errorf("cannot read PDF header: %w", err)
	}
	
	if string(header) != "%PDF" {
		return fmt.Errorf("invalid PDF header")
	}
	
	return nil
}

// getInstallInstructions returns platform-specific installation instructions
func (c *HTMLConverter) getInstallInstructions() string {
	msg := "No HTML converter available.\n\n"
	
	switch runtime.GOOS {
	case "darwin":
		msg += "Install one of:\n"
		msg += "  wkhtmltopdf (fast):  brew install wkhtmltopdf\n"
		msg += "  Chrome (best):       brew install --cask google-chrome\n"
	case "linux":
		msg += "Install one of:\n"
		msg += "  wkhtmltopdf (fast):  sudo apt-get install wkhtmltopdf\n"
		msg += "                       OR sudo dnf install wkhtmltopdf\n"
		msg += "  Chrome (best):       sudo apt-get install chromium-browser\n"
		msg += "                       OR sudo dnf install chromium\n"
	case "windows":
		msg += "Install one of:\n"
		msg += "  wkhtmltopdf (fast):  Download from https://wkhtmltopdf.org/\n"
		msg += "  Chrome (best):       winget install Google.Chrome\n"
	}
	
	return msg
}

// isURL checks if the input is a URL (http://, https://, file://)
func isURL(input string) bool {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "file://") {
		return true
	}
	
	// Try parsing as URL
	u, err := url.Parse(input)
	if err != nil {
		return false
	}
	
	return u.Scheme != ""
}

// ConvertFromReader converts HTML from stdin/reader to PDF
func (c *HTMLConverter) ConvertFromReader(ctx context.Context, reader io.Reader, output string, opts Options) error {
	// Create temporary HTML file
	tmpFile, err := os.CreateTemp("", "any2pdf-html-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	
	// Cleanup temp file unless preserve flag is set
	if !c.preserveTemp {
		defer os.Remove(tmpPath)
	}
	
	// Copy reader content to temp file
	_, err = io.Copy(tmpFile, reader)
	tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	
	// Convert the temp file
	return c.ConvertWithContext(ctx, tmpPath, output, opts)
}

// GetAvailableEngines returns which conversion engines are available
func (c *HTMLConverter) GetAvailableEngines() []string {
	engines := []string{}
	if c.wkAvailable {
		engines = append(engines, "wkhtmltopdf")
	}
	if c.chromeAvailable {
		engines = append(engines, "chrome")
	}
	return engines
}
