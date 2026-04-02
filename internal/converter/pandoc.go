package converter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/De-pitcher/any2pdf/internal/utils"
)

// PandocConverter handles conversion of text and markdown files using pandoc
type PandocConverter struct {
	timeout  time.Duration
	metadata map[string]string
}

// NewPandocConverter creates a new PandocConverter
func NewPandocConverter() *PandocConverter {
	return &PandocConverter{
		timeout:  5 * time.Minute, // Default 5-minute timeout
		metadata: make(map[string]string),
	}
}

// Name returns the converter name
func (c *PandocConverter) Name() string {
	return "pandoc"
}

// SetTimeout sets a custom timeout for conversions
func (c *PandocConverter) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetMetadata sets metadata fields for the PDF
func (c *PandocConverter) SetMetadata(key, value string) {
	c.metadata[key] = value
}

// ClearMetadata clears all metadata fields
func (c *PandocConverter) ClearMetadata() {
	c.metadata = make(map[string]string)
}

// IsSupportedFormat checks if the file extension is supported by pandoc
func (c *PandocConverter) IsSupportedFormat(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	supported := map[string]bool{
		".txt":       true,
		".md":        true,
		".markdown":  true,
		".rst":       true,
		".org":       true,
		".textile":   true,
		".mediawiki": true,
	}
	return supported[ext]
}

// CheckAvailable verifies pandoc is installed
func (c *PandocConverter) CheckAvailable() error {
	if !utils.CommandExists("pandoc") {
		return &utils.ConverterNotFoundError{
			FileType:     "Text/Markdown",
			MissingTool:  "pandoc",
			InstallGuide: c.GetInstallInstructions(),
		}
	}
	return nil
}

// checkDependencies verifies both pandoc and LaTeX engine are available
func (c *PandocConverter) checkDependencies() error {
	// Check pandoc first
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Check for LaTeX engines (try xelatex first, then pdflatex)
	hasXeLaTeX := utils.CommandExists("xelatex")
	hasPDFLaTeX := utils.CommandExists("pdflatex")

	if !hasXeLaTeX && !hasPDFLaTeX {
		return fmt.Errorf(`❌ LaTeX engine not found

Pandoc requires a LaTeX distribution to generate PDFs.

Install LaTeX:
  macOS:   brew install basictex
           OR: brew install mactex
  Linux:   sudo apt-get install texlive-xetex
           OR: sudo apt-get install texlive-latex-base
  Windows: choco install miktex
           OR: download from https://miktex.org/

After installing, run: any2pdf check`)
	}

	return nil
}

// getPreferredLatexEngine returns the best available LaTeX engine
func (c *PandocConverter) getPreferredLatexEngine() string {
	// Prefer xelatex for Unicode support
	if utils.CommandExists("xelatex") {
		return "xelatex"
	}
	// Fallback to pdflatex
	if utils.CommandExists("pdflatex") {
		return "pdflatex"
	}
	return "xelatex" // Default, will error if not available
}

// Convert converts the file to PDF using pandoc
func (c *PandocConverter) Convert(inputPath, outputPath string, opts Options) error {
	return c.ConvertWithContext(context.Background(), inputPath, outputPath, opts)
}

// ConvertWithContext converts the file to PDF with context support
func (c *PandocConverter) ConvertWithContext(ctx context.Context, inputPath, outputPath string, opts Options) error {
	// Check dependencies
	if err := c.checkDependencies(); err != nil {
		return err
	}

	// Handle stdin input
	if inputPath != "-" {
		// Validate input file exists
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
	}

	// Generate output path if not specified
	finalOutputPath := outputPath
	if finalOutputPath == "" {
		if inputPath == "-" {
			finalOutputPath = "output.pdf"
		} else {
			finalOutputPath = utils.GenerateOutputPath(inputPath, "")
		}
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build command arguments
	args := c.buildArgs(inputPath, finalOutputPath, opts)

	// Execute pandoc
	cmd := exec.CommandContext(timeoutCtx, "pandoc", args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		exitCode := 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		// Clean up on failure
		os.Remove(finalOutputPath)

		// Enhance error message for common issues
		errorMsg := string(output)
		if strings.Contains(errorMsg, "! Package inputenc Error") || 
		   strings.Contains(errorMsg, "Unicode") {
			errorMsg += "\n\nℹ️  Tip: This may be a Unicode encoding issue. Ensure XeLaTeX is installed."
		}

		return &utils.ConversionFailedError{
			InputFile:  inputPath,
			OutputFile: finalOutputPath,
			Converter:  c.Name(),
			ExitCode:   exitCode,
			Stderr:     errorMsg,
			Cause:      err,
		}
	}

	// Verify output file was created
	if _, err := os.Stat(finalOutputPath); os.IsNotExist(err) {
		return fmt.Errorf("conversion completed but output file not found: %s", finalOutputPath)
	}

	return nil
}

// buildArgs constructs the pandoc command arguments
func (c *PandocConverter) buildArgs(inputPath, outputPath string, opts Options) []string {
	args := []string{}

	// Input file (or stdin)
	if inputPath != "-" {
		args = append(args, inputPath)
	}

	// Output file
	args = append(args, "-o", outputPath)

	// PDF engine (prefer xelatex for Unicode support)
	pdfEngine := c.getPreferredLatexEngine()
	args = append(args, "--pdf-engine="+pdfEngine)

	// Geometry (margins) based on quality
	margin := c.getMarginForQuality(opts.Quality)
	args = append(args, "-V", "geometry:margin="+margin)

	// Font size based on quality
	fontSize := c.getFontSizeForQuality(opts.Quality)
	args = append(args, "-V", "fontsize="+fontSize)

	// Font family (only for xelatex)
	if pdfEngine == "xelatex" {
		args = append(args, "-V", "mainfont=DejaVu Sans")
	}

	// Syntax highlighting
	args = append(args, "--highlight-style=tango")

	// Add metadata
	for key, value := range c.metadata {
		args = append(args, "--metadata", fmt.Sprintf("%s=%s", key, value))
	}

	return args
}

// getMarginForQuality returns the margin setting for the given quality level
func (c *PandocConverter) getMarginForQuality(quality Quality) string {
	switch quality {
	case Screen:
		return "0.75in" // Smaller margins for screen
	case Printer:
		return "1.2in" // Larger margins for printing
	default:
		return "1in" // Default margins
	}
}

// getFontSizeForQuality returns the font size for the given quality level
func (c *PandocConverter) getFontSizeForQuality(quality Quality) string {
	switch quality {
	case Screen:
		return "10pt" // Smaller font for screen
	case Printer:
		return "12pt" // Larger font for printing
	default:
		return "11pt" // Default font size
	}
}

// GetInstallInstructions returns installation instructions for pandoc
func (c *PandocConverter) GetInstallInstructions() string {
	return `❌ pandoc not found

Install pandoc:
  macOS:   brew install pandoc
  Linux:   sudo apt-get install pandoc
  Windows: choco install pandoc
           OR: https://pandoc.org/installing.html

Also need LaTeX engine:
  macOS:   brew install basictex
           OR: brew install mactex
  Linux:   sudo apt-get install texlive-xetex
           OR: sudo apt-get install texlive-latex-base
  Windows: choco install miktex
           OR: https://miktex.org/download

After installing, verify with: any2pdf check`
}
