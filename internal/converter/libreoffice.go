package converter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/De-pitcher/any2pdf/internal/utils"
)

// LibreOfficeConverter handles conversion of Office documents using LibreOffice
type LibreOfficeConverter struct {
	timeout      time.Duration
	binaryPath   string
	preserveTemp bool // For debugging
}

// NewLibreOfficeConverter creates a new LibreOfficeConverter
func NewLibreOfficeConverter() *LibreOfficeConverter {
	return &LibreOfficeConverter{
		timeout:      10 * time.Minute, // Office docs can be large
		preserveTemp: false,
	}
}

// Name returns the converter name
func (c *LibreOfficeConverter) Name() string {
	return "libreoffice"
}

// SetTimeout sets a custom timeout for conversions
func (c *LibreOfficeConverter) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetPreserveTemp sets whether to preserve temporary files (for debugging)
func (c *LibreOfficeConverter) SetPreserveTemp(preserve bool) {
	c.preserveTemp = preserve
}

// IsSupportedFormat checks if the file extension is supported
func (c *LibreOfficeConverter) IsSupportedFormat(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	supported := map[string]bool{
		// Microsoft Office formats
		".docx": true,
		".doc":  true,
		".xlsx": true,
		".xls":  true,
		".pptx": true,
		".ppt":  true,
		// OpenDocument formats
		".odt": true,
		".ods": true,
		".odp": true,
		// Rich Text Format
		".rtf": true,
	}
	return supported[ext]
}

// CheckAvailable verifies LibreOffice is installed
func (c *LibreOfficeConverter) CheckAvailable() error {
	binary := c.findBinary()
	if binary == "" {
		return &utils.ConverterNotFoundError{
			FileType:     "Office Documents",
			MissingTool:  "libreoffice",
			InstallGuide: c.GetInstallInstructions(),
		}
	}
	c.binaryPath = binary
	return nil
}

// findBinary locates the LibreOffice binary on the system
func (c *LibreOfficeConverter) findBinary() string {
	// If already cached, return it
	if c.binaryPath != "" {
		return c.binaryPath
	}

	// Try common command names first
	commands := []string{"libreoffice", "soffice"}
	for _, cmd := range commands {
		if path, err := exec.LookPath(cmd); err == nil {
			return path
		}
	}

	// Try OS-specific paths
	switch runtime.GOOS {
	case "darwin": // macOS
		paths := []string{
			"/Applications/LibreOffice.app/Contents/MacOS/soffice",
			"/usr/local/bin/soffice",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}

	case "windows":
		paths := []string{
			"C:\\Program Files\\LibreOffice\\program\\soffice.exe",
			"C:\\Program Files (x86)\\LibreOffice\\program\\soffice.exe",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}

	case "linux":
		// Already tried which libreoffice/soffice above
		paths := []string{
			"/usr/bin/libreoffice",
			"/usr/bin/soffice",
			"/usr/local/bin/soffice",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// Convert converts a single Office document to PDF
func (c *LibreOfficeConverter) Convert(inputPath, outputPath string, opts Options) error {
	return c.ConvertWithContext(context.Background(), inputPath, outputPath, opts)
}

// ConvertWithContext converts a document to PDF with context support
func (c *LibreOfficeConverter) ConvertWithContext(ctx context.Context, inputPath, outputPath string, opts Options) error {
	// Check dependencies
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Validate input file
	if err := utils.ValidateInputFile(inputPath); err != nil {
		return err
	}

	// Validate format
	if !c.IsSupportedFormat(inputPath) {
		return &utils.UnsupportedFileTypeError{
			Extension: filepath.Ext(inputPath),
			FilePath:  inputPath,
		}
	}

	// Create temporary directory for conversion
	tempDir, err := os.MkdirTemp("", "any2pdf-libreoffice-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	if !c.preserveTemp {
		defer os.RemoveAll(tempDir)
	}

	// Convert using temp directory
	tempOutputName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) + ".pdf"
	tempOutputPath := filepath.Join(tempDir, tempOutputName)

	// Build command
	cmd := c.buildCommand(inputPath, tempDir, opts)
	cmd.Dir = tempDir

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd = exec.CommandContext(timeoutCtx, cmd.Path, cmd.Args[1:]...)
	cmd.Dir = tempDir

	output, err := cmd.CombinedOutput()

	if err != nil {
		exitCode := 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		// Check for specific error conditions
		errorMsg := string(output)
		errorMsg = c.handleConversionError(err, errorMsg)

		return &utils.ConversionFailedError{
			InputFile:  inputPath,
			OutputFile: outputPath,
			Converter:  c.Name(),
			ExitCode:   exitCode,
			Stderr:     errorMsg,
			Cause:      err,
		}
	}

	// Verify temporary PDF was created
	if err := c.validatePDFCreated(tempOutputPath); err != nil {
		return fmt.Errorf("LibreOffice completed but output PDF not found: %w", err)
	}

	// Move the PDF to the final output location
	finalOutputPath := outputPath
	if finalOutputPath == "" {
		finalOutputPath = utils.GenerateOutputPath(inputPath, "")
	}

	// Copy the file (use copy instead of rename for cross-device compatibility)
	if err := c.copyFile(tempOutputPath, finalOutputPath); err != nil {
		return fmt.Errorf("failed to move output PDF: %w", err)
	}

	return nil
}

// ConvertBatch converts multiple Office documents to PDF
func (c *LibreOfficeConverter) ConvertBatch(ctx context.Context, inputPaths []string, outputDir string, opts Options) ([]string, error) {
	if len(inputPaths) == 0 {
		return nil, fmt.Errorf("no input files provided")
	}

	// Check dependencies
	if err := c.CheckAvailable(); err != nil {
		return nil, err
	}

	// Validate all inputs
	for _, path := range inputPaths {
		if err := utils.ValidateInputFile(path); err != nil {
			return nil, fmt.Errorf("invalid input file %s: %w", path, err)
		}
		if !c.IsSupportedFormat(path) {
			return nil, &utils.UnsupportedFileTypeError{
				Extension: filepath.Ext(path),
				FilePath:  path,
			}
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Convert each file
	outputs := make([]string, 0, len(inputPaths))
	for _, inputPath := range inputPaths {
		baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
		outputPath := filepath.Join(outputDir, baseName+".pdf")

		if err := c.ConvertWithContext(ctx, inputPath, outputPath, opts); err != nil {
			return outputs, fmt.Errorf("failed to convert %s: %w", inputPath, err)
		}

		outputs = append(outputs, outputPath)
	}

	return outputs, nil
}

// buildCommand constructs the LibreOffice command
func (c *LibreOfficeConverter) buildCommand(inputPath, outputDir string, opts Options) *exec.Cmd {
	args := []string{
		c.binaryPath,
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outputDir,
	}

	// Add quality-specific arguments
	qualityArgs := c.buildQualityArgs(opts.Quality)
	args = append(args, qualityArgs...)

	// Add format-specific arguments
	ext := strings.ToLower(filepath.Ext(inputPath))
	formatArgs := c.getFormatSpecificArgs(ext)
	args = append(args, formatArgs...)

	// Add input file (must be last)
	args = append(args, inputPath)

	cmd := &exec.Cmd{
		Path: c.binaryPath,
		Args: args,
	}

	return cmd
}

// buildQualityArgs returns quality-specific PDF export options
func (c *LibreOfficeConverter) buildQualityArgs(quality Quality) []string {
	// Note: LibreOffice PDF options are passed differently than standard CLI args
	// These would need to be passed via --pdf-options or filter data
	// For now, we'll use simpler approach
	switch quality {
	case Screen:
		// Optimize for screen viewing - reduce image quality
		return []string{}
	case Printer:
		// Optimize for printing - preserve quality
		return []string{}
	default:
		return []string{}
	}
}

// getFormatSpecificArgs returns format-specific conversion arguments
func (c *LibreOfficeConverter) getFormatSpecificArgs(ext string) []string {
	args := []string{}

	switch ext {
	case ".xlsx", ".xls", ".ods":
		// Spreadsheet-specific options would go here
		// LibreOffice doesn't expose many CLI options for this
	case ".pptx", ".ppt", ".odp":
		// Presentation-specific options
	}

	return args
}

// getFilter returns the LibreOffice filter string for a given format (if needed)
func (c *LibreOfficeConverter) getFilter(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// Most formats work with auto-detection, but some need explicit filters
	filters := map[string]string{
		".doc":  "MS Word 97",
		".xls":  "MS Excel 97",
		".ppt":  "MS PowerPoint 97",
		".rtf":  "Rich Text Format",
	}

	if filter, exists := filters[ext]; exists {
		return filter
	}

	return "" // Auto-detect
}

// handleConversionError enhances error messages with common issues
func (c *LibreOfficeConverter) handleConversionError(err error, stderr string) string {
	errorMsg := stderr

	// Check for common errors
	if strings.Contains(stderr, "password") || strings.Contains(stderr, "encrypted") {
		errorMsg += "\n\nℹ️  This file appears to be password-protected. LibreOffice cannot convert encrypted documents."
	}

	if strings.Contains(stderr, "locked") {
		errorMsg += "\n\nℹ️  The file may be locked or in use by another program. Close it and try again."
	}

	if strings.Contains(stderr, "corrupted") || strings.Contains(stderr, "damaged") {
		errorMsg += "\n\nℹ️  The file appears to be corrupted. Try opening it in the original application first."
	}

	return errorMsg
}

// validatePDFCreated checks if the output PDF exists and is valid
func (c *LibreOfficeConverter) validatePDFCreated(outputPath string) error {
	info, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("output file not created: %s", outputPath)
	}
	if err != nil {
		return fmt.Errorf("cannot access output file: %w", err)
	}

	if info.Size() == 0 {
		return fmt.Errorf("output PDF is empty (0 bytes)")
	}

	// Verify it's a PDF file (check magic bytes)
	data := make([]byte, 4)
	f, err := os.Open(outputPath)
	if err != nil {
		return fmt.Errorf("cannot read output file: %w", err)
	}
	defer f.Close()

	if _, err := f.Read(data); err != nil {
		return fmt.Errorf("cannot verify PDF format: %w", err)
	}

	if string(data) != "%PDF" {
		return fmt.Errorf("output file is not a valid PDF")
	}

	return nil
}

// copyFile copies a file from src to dst
func (c *LibreOfficeConverter) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// GetInstallInstructions returns installation instructions for LibreOffice
func (c *LibreOfficeConverter) GetInstallInstructions() string {
	return `❌ LibreOffice not found

Install LibreOffice:
  macOS:   brew install --cask libreoffice
  Linux:   sudo apt-get install libreoffice-core
           OR: sudo apt-get install libreoffice
  Windows: choco install libreoffice
           OR: https://www.libreoffice.org/download/

For headless Linux servers:
  sudo apt-get install libreoffice-core libreoffice-common

After installing, verify with: any2pdf check`
}
