package converter

import (
	"github.com/De-pitcher/any2pdf/internal/utils"
	"fmt"
	"time"
)

// HTMLConverter handles conversion of HTML files using wkhtmltopdf
type HTMLConverter struct{}

// NewHTMLConverter creates a new HTMLConverter
func NewHTMLConverter() *HTMLConverter {
	return &HTMLConverter{}
}

// Name returns the converter name
func (c *HTMLConverter) Name() string {
	return "wkhtmltopdf"
}

// CheckAvailable verifies wkhtmltopdf is installed
func (c *HTMLConverter) CheckAvailable() error {
	if !utils.CommandExists("wkhtmltopdf") {
		return &utils.ConverterNotFoundError{
			FileType:     "HTML",
			MissingTool:  "wkhtmltopdf",
			InstallGuide: "https://wkhtmltopdf.org/downloads.html",
		}
	}
	return nil
}

// Convert converts the HTML file to PDF using wkhtmltopdf
func (c *HTMLConverter) Convert(inputPath, outputPath string, opts Options) error {
	// Check if wkhtmltopdf is available
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Build wkhtmltopdf command arguments
	args := []string{
		"--quiet", // Suppress output
	}

	// Add quality-specific options
	switch opts.Quality {
	case Screen:
		args = append(args, "--lowquality", "--dpi", "75")
	case Printer:
		args = append(args, "--dpi", "300")
	default:
		args = append(args, "--dpi", "150")
	}

	// Add input and output files
	args = append(args, inputPath, outputPath)

	// Execute wkhtmltopdf
	result := utils.ExecuteCommand("wkhtmltopdf", args, utils.ExecOptions{
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

// GetInstallInstructions returns installation instructions for wkhtmltopdf
func (c *HTMLConverter) GetInstallInstructions() string {
	return fmt.Sprintf(`wkhtmltopdf is not installed. Please install it:

Ubuntu/Debian:  sudo apt-get install wkhtmltopdf
macOS:          brew install wkhtmltopdf
Windows:        Download from https://wkhtmltopdf.org/downloads.html`)
}
