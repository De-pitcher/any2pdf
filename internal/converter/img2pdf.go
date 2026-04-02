package converter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/De-pitcher/any2pdf/internal/utils"
)

// ImageConverter handles conversion of images using img2pdf
type ImageConverter struct {
	timeout time.Duration
}

// NewImageConverter creates a new ImageConverter
func NewImageConverter() *ImageConverter {
	return &ImageConverter{
		timeout: 5 * time.Minute, // Default 5-minute timeout
	}
}

// Name returns the converter name
func (c *ImageConverter) Name() string {
	return "img2pdf"
}

// SetTimeout sets a custom timeout for conversions
func (c *ImageConverter) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// CheckAvailable verifies img2pdf is installed
func (c *ImageConverter) CheckAvailable() error {
	if !utils.CommandExists("img2pdf") {
		return &utils.ConverterNotFoundError{
			FileType:     "Image",
			MissingTool:  "img2pdf",
			InstallGuide: c.GetInstallInstructions(),
		}
	}
	return nil
}

// IsSupportedImageFormat checks if the file extension is a supported image format
func (c *ImageConverter) IsSupportedImageFormat(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	supported := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".tif":  true,
	}
	return supported[ext]
}

// Convert converts a single image or multiple images to PDF using img2pdf
func (c *ImageConverter) Convert(inputPath, outputPath string, opts Options) error {
	return c.ConvertWithContext(context.Background(), inputPath, outputPath, opts)
}

// ConvertWithContext converts image(s) to PDF with context support
func (c *ImageConverter) ConvertWithContext(ctx context.Context, inputPath, outputPath string, opts Options) error {
	// Check if img2pdf is available
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Validate input file exists
	if err := utils.ValidateInputFile(inputPath); err != nil {
		return err
	}

	// Validate it's a supported image format
	if !c.IsSupportedImageFormat(inputPath) {
		return &utils.UnsupportedFileTypeError{
			Extension: filepath.Ext(inputPath),
			FilePath:  inputPath,
		}
	}

	// Convert single image
	return c.convertSingle(ctx, inputPath, outputPath, opts)
}

// ConvertMultiple converts multiple images into a single PDF
func (c *ImageConverter) ConvertMultiple(inputPaths []string, outputPath string, opts Options) error {
	return c.ConvertMultipleWithContext(context.Background(), inputPaths, outputPath, opts)
}

// ConvertMultipleWithContext converts multiple images into a single PDF with context support
func (c *ImageConverter) ConvertMultipleWithContext(ctx context.Context, inputPaths []string, outputPath string, opts Options) error {
	// Check input first before checking tool availability
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input images provided")
	}

	// Check if img2pdf is available
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Validate all input files
	for _, path := range inputPaths {
		if err := utils.ValidateInputFile(path); err != nil {
			return fmt.Errorf("invalid input file %s: %w", path, err)
		}
		if !c.IsSupportedImageFormat(path) {
			return &utils.UnsupportedFileTypeError{
				Extension: filepath.Ext(path),
				FilePath:  path,
			}
		}
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build command arguments
	args := c.buildArgs(inputPaths, outputPath, opts)

	// Execute command
	cmd := exec.CommandContext(timeoutCtx, "img2pdf", args...)
	
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		exitCode := 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		return &utils.ConversionFailedError{
			InputFile:  strings.Join(inputPaths, ", "),
			OutputFile: outputPath,
			Converter:  c.Name(),
			ExitCode:   exitCode,
			Stderr:     string(output),
			Cause:      err,
		}
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("conversion completed but output file not found: %s", outputPath)
	}

	return nil
}

// convertSingle converts a single image to PDF
func (c *ImageConverter) convertSingle(ctx context.Context, inputPath, outputPath string, opts Options) error {
	// Use a temporary file if output path is not specified
	finalOutputPath := outputPath
	if finalOutputPath == "" {
		finalOutputPath = utils.GenerateOutputPath(inputPath, "")
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build command arguments for single image
	args := c.buildArgs([]string{inputPath}, finalOutputPath, opts)

	// Execute command
	cmd := exec.CommandContext(timeoutCtx, "img2pdf", args...)
	
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		exitCode := 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		// Clean up on failure
		os.Remove(finalOutputPath)

		return &utils.ConversionFailedError{
			InputFile:  inputPath,
			OutputFile: finalOutputPath,
			Converter:  c.Name(),
			ExitCode:   exitCode,
			Stderr:     string(output),
			Cause:      err,
		}
	}

	// Verify output file was created
	if _, err := os.Stat(finalOutputPath); os.IsNotExist(err) {
		return fmt.Errorf("conversion completed but output file not found: %s", finalOutputPath)
	}

	return nil
}

// buildArgs constructs the img2pdf command arguments
func (c *ImageConverter) buildArgs(inputPaths []string, outputPath string, opts Options) []string {
	args := []string{}

	// Add input files
	args = append(args, inputPaths...)

	// Add output flag
	args = append(args, "-o", outputPath)

	// Add DPI settings based on quality
	var dpi int
	switch opts.Quality {
	case Screen:
		dpi = 72 // Screen quality as specified
	case Printer:
		dpi = 300 // Printer quality
	default:
		dpi = 150 // Default quality
	}
	args = append(args, "--dpi", strconv.Itoa(dpi))

	return args
}

// GetInstallInstructions returns installation instructions for img2pdf
func (c *ImageConverter) GetInstallInstructions() string {
	return `img2pdf is not installed. Please install it:

Ubuntu/Debian:  sudo apt-get install img2pdf
                OR: pip3 install img2pdf

macOS:          pip3 install img2pdf
                OR: brew install img2pdf

Windows:        pip install img2pdf

For pip installation, ensure Python and pip are installed first.`
}
