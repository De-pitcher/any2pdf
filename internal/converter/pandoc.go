package converter

import (
	"any2pdf/internal/utils"
	"fmt"
	"time"
)

// PandocConverter handles conversion of text and markdown files using pandoc
type PandocConverter struct{}

// NewPandocConverter creates a new PandocConverter
func NewPandocConverter() *PandocConverter {
	return &PandocConverter{}
}

// Name returns the converter name
func (c *PandocConverter) Name() string {
	return "pandoc"
}

// CheckAvailable verifies pandoc is installed
func (c *PandocConverter) CheckAvailable() error {
	if !utils.CommandExists("pandoc") {
		return &utils.ConverterNotFoundError{
			FileType:     "Text/Markdown",
			MissingTool:  "pandoc",
			InstallGuide: "https://pandoc.org/installing.html",
		}
	}
	return nil
}

// Convert converts the file to PDF using pandoc
func (c *PandocConverter) Convert(inputPath, outputPath string, opts Options) error {
	// Check if pandoc is available
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Build pandoc command arguments
	args := []string{
		inputPath,
		"-o", outputPath,
		"--pdf-engine=pdflatex", // or wkhtmltopdf, xelatex, etc.
	}

	// Add quality-specific options
	switch opts.Quality {
	case Screen:
		args = append(args, "-V", "papersize=a4", "-V", "fontsize=10pt")
	case Printer:
		args = append(args, "-V", "papersize=a4", "-V", "fontsize=12pt")
	default:
		args = append(args, "-V", "papersize=a4")
	}

	// Execute pandoc
	result := utils.ExecuteCommand("pandoc", args, utils.ExecOptions{
		Timeout: 60 * time.Second,
		Quiet:   opts.Quiet,
	})

	if result.Error != nil {
		return &utils.ConversionFailedError{
			InputFile:  inputPath,
			OutputFile: outputPath,
			Converter:  c.Name(),
			ExitCode:   result.ExitCode,
			Stderr:     result.Stderr,
			Cause:      result.Error,
		}
	}

	return nil
}

// GetInstallInstructions returns installation instructions for pandoc
func (c *PandocConverter) GetInstallInstructions() string {
	return fmt.Sprintf(`Pandoc is not installed. Please install it:

Ubuntu/Debian:  sudo apt-get install pandoc texlive-latex-base
macOS:          brew install pandoc
Windows:        Download from https://pandoc.org/installing.html

Note: For full PDF support, you may also need a LaTeX distribution.`)
}
