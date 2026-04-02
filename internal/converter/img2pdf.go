package converter

import (
	"any2pdf/internal/utils"
	"fmt"
	"strconv"
	"time"
)

// ImageConverter handles conversion of images using img2pdf
type ImageConverter struct{}

// NewImageConverter creates a new ImageConverter
func NewImageConverter() *ImageConverter {
	return &ImageConverter{}
}

// Name returns the converter name
func (c *ImageConverter) Name() string {
	return "img2pdf"
}

// CheckAvailable verifies img2pdf is installed
func (c *ImageConverter) CheckAvailable() error {
	if !utils.CommandExists("img2pdf") {
		return &utils.ConverterNotFoundError{
			FileType:     "Image",
			MissingTool:  "img2pdf",
			InstallGuide: "pip3 install img2pdf",
		}
	}
	return nil
}

// Convert converts the image to PDF using img2pdf
func (c *ImageConverter) Convert(inputPath, outputPath string, opts Options) error {
	// Check if img2pdf is available
	if err := c.CheckAvailable(); err != nil {
		return err
	}

	// Build img2pdf command arguments
	args := []string{
		inputPath,
		"-o", outputPath,
	}

	// Add DPI settings based on quality
	var dpi int
	switch opts.Quality {
	case Screen:
		dpi = 75
	case Printer:
		dpi = 300
	default:
		dpi = 150
	}
	
	args = append(args, "--dpi", strconv.Itoa(dpi))

	// Execute img2pdf
	result := utils.ExecuteCommand("img2pdf", args, utils.ExecOptions{
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

// GetInstallInstructions returns installation instructions for img2pdf
func (c *ImageConverter) GetInstallInstructions() string {
	return fmt.Sprintf(`img2pdf is not installed. Please install it:

Ubuntu/Debian:  sudo apt-get install img2pdf
macOS:          pip3 install img2pdf
Windows:        pip install img2pdf`)
}
