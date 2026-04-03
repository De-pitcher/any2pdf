package converter

import "fmt"

// Quality represents the quality level for PDF conversion
type Quality int

const (
	Screen  Quality = iota // Lower DPI, smaller files (75 DPI)
	Default                // Balanced (150 DPI)
	Printer                // High quality (300 DPI)
)

// String returns the string representation of Quality
func (q Quality) String() string {
	switch q {
	case Screen:
		return "screen"
	case Printer:
		return "printer"
	default:
		return "default"
	}
}

// ParseQuality converts a string to Quality
func ParseQuality(s string) Quality {
	switch s {
	case "screen":
		return Screen
	case "printer":
		return Printer
	default:
		return Default
	}
}

// Options contains options for conversion
type Options struct {
	Quality Quality
	Quiet   bool
	Extra   map[string]interface{} // Extra converter-specific options (e.g., HTML page size, margins)
}

// Converter is the interface that all converters must implement
type Converter interface {
	// Convert converts the input file to PDF
	Convert(inputPath, outputPath string, opts Options) error

	// CheckAvailable verifies that the converter's dependencies are installed
	CheckAvailable() error

	// Name returns the name of the converter
	Name() string
}

// ConverterNotFoundError indicates that a required converter tool is not installed
type ConverterNotFoundError struct {
	Converter string
	Message   string
}

func (e *ConverterNotFoundError) Error() string {
	return fmt.Sprintf("converter not found: %s\n%s", e.Converter, e.Message)
}
