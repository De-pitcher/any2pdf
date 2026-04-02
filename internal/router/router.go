package router

import (
	"any2pdf/internal/converter"
	"any2pdf/internal/detector"
	"any2pdf/internal/utils"
)

// Router handles routing files to appropriate converters
type Router struct {
	detector *detector.Detector
	registry *ConverterRegistry
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{
		detector: detector.NewDetector(),
		registry: NewConverterRegistry(),
	}
}

// Route determines the appropriate converter for a file
func (r *Router) Route(filePath string) (converter.Converter, error) {
	// Detect file type
	fileType, err := r.detector.DetectFileType(filePath)
	if err != nil {
		return nil, err
	}

	// Get converter from registry
	conv, exists := r.registry.Get(fileType)
	if !exists {
		return nil, &utils.UnsupportedFileTypeError{
			Extension: utils.GetFileExtension(filePath),
			FilePath:  filePath,
		}
	}

	// Check if converter is available
	if err := conv.CheckAvailable(); err != nil {
		return nil, err
	}

	return conv, nil
}

// Convert performs the full conversion process
func (r *Router) Convert(inputPath, outputPath string, opts converter.Options) error {
	// Validate input file
	if err := utils.ValidateInputFile(inputPath); err != nil {
		return err
	}

	// Generate output path if not specified
	if outputPath == "" {
		outputPath = utils.GenerateOutputPath(inputPath, "")
	}

	// Route to converter
	conv, err := r.Route(inputPath)
	if err != nil {
		return err
	}

	// Perform conversion
	return conv.Convert(inputPath, outputPath, opts)
}

// GetRegistry returns the converter registry
func (r *Router) GetRegistry() *ConverterRegistry {
	return r.registry
}
