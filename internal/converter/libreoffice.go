package converter

import (
	"github.com/De-pitcher/any2pdf/internal/utils"
	"fmt"
	"path/filepath"
	"time"
)

// LibreOfficeConverter handles conversion of Office documents using LibreOffice
type LibreOfficeConverter struct{}

// NewLibreOfficeConverter creates a new LibreOfficeConverter
func NewLibreOfficeConverter() *LibreOfficeConverter {
	return &LibreOfficeConverter{}
}

// Name returns the converter name
func (c *LibreOfficeConverter) Name() string {
	return "libreoffice"
}

// CheckAvailable verifies LibreOffice is installed
func (c *LibreOfficeConverter) CheckAvailable() error {
	// Try common LibreOffice command names
	commands := []string{"libreoffice", "soffice"}
	
	for _, cmd := range commands {
		if utils.CommandExists(cmd) {
			return nil
		}
	}
	
	return &utils.ConverterNotFoundError{
		FileType:     "Office Documents",
		MissingTool:  "libreoffice",
		InstallGuide: "https://www.libreoffice.org/download/",
	}
}

// Convert converts the file to PDF using LibreOffice
func (c *LibreOfficeConverter) Convert(inputPath, outputPath string, opts Options) error {
	// Check if LibreOffice is available
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Get the command name (libreoffice or soffice)
	cmdName := "libreoffice"
	if !utils.CommandExists(cmdName) {
		cmdName = "soffice"
	}

	// LibreOffice outputs to the same directory as input with .pdf extension
	// So we need to specify the output directory
	outputDir := filepath.Dir(outputPath)

	// Build LibreOffice command arguments
	args := []string{
		"--headless",                    // Run without GUI
		"--convert-to", "pdf",           // Convert to PDF
		"--outdir", outputDir,           // Output directory
		inputPath,                       // Input file
	}

	// Execute LibreOffice
	result := utils.ExecuteCommand(cmdName, args, utils.ExecOptions{
		Timeout: 120 * time.Second, // Office docs can take longer
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

	// LibreOffice creates output with input filename + .pdf
	// If user specified different output name, we need to rename
	expectedOutput := filepath.Join(outputDir, filepath.Base(inputPath)+".pdf")
	if expectedOutput != outputPath {
		// TODO: Implement file rename using os.Rename
	}

	return nil
}

// GetInstallInstructions returns installation instructions for LibreOffice
func (c *LibreOfficeConverter) GetInstallInstructions() string {
	return fmt.Sprintf(`LibreOffice is not installed. Please install it:

Ubuntu/Debian:  sudo apt-get install libreoffice
macOS:          brew install --cask libreoffice
Windows:        Download from https://www.libreoffice.org/download/`)
}
