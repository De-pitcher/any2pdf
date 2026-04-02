package detector

import (
	"github.com/De-pitcher/any2pdf/internal/utils"
)

// Detector handles file type detection
type Detector struct{}

// NewDetector creates a new Detector instance
func NewDetector() *Detector {
	return &Detector{}
}

// DetectFileType determines the file type based on extension
func (d *Detector) DetectFileType(filePath string) (FileType, error) {
	// Get file extension
	ext := utils.GetFileExtension(filePath)
	if ext == "" {
		return Unknown, &utils.UnsupportedFileTypeError{
			Extension: "none",
			FilePath:  filePath,
		}
	}

	// Look up file type
	fileType, exists := ExtensionMap[ext]
	if !exists {
		return Unknown, &utils.UnsupportedFileTypeError{
			Extension: ext,
			FilePath:  filePath,
		}
	}

	return fileType, nil
}

// IsSupportedExtension checks if an extension is supported
func (d *Detector) IsSupportedExtension(ext string) bool {
	_, exists := ExtensionMap[ext]
	return exists
}

// GetSupportedExtensions returns a list of all supported extensions
func (d *Detector) GetSupportedExtensions() []string {
	extensions := make([]string, 0, len(ExtensionMap))
	for ext := range ExtensionMap {
		extensions = append(extensions, ext)
	}
	return extensions
}
