package router

import (
	"github.com/De-pitcher/any2pdf/internal/converter"
	"github.com/De-pitcher/any2pdf/internal/detector"
)

//ConverterRegistry maps file types to converters
type ConverterRegistry struct {
	converters map[detector.FileType]converter.Converter
}

// NewConverterRegistry creates a new converter registry
func NewConverterRegistry() *ConverterRegistry {
	registry := &ConverterRegistry{
		converters: make(map[detector.FileType]converter.Converter),
	}
	
	// Register default converters
	registry.registerDefaultConverters()
	
	return registry
}

// registerDefaultConverters registers all built-in converters
func (r *ConverterRegistry) registerDefaultConverters() {
	// Register Pandoc for text and markdown formats
	pandocConverter := converter.NewPandocConverter()
	r.Register(detector.Text, pandocConverter)
	r.Register(detector.Markdown, pandocConverter)
	r.Register(detector.ReStructuredText, pandocConverter)
	r.Register(detector.OrgMode, pandocConverter)
	r.Register(detector.Textile, pandocConverter)
	
	// Register LibreOffice for Office documents
	libreOfficeConverter := converter.NewLibreOfficeConverter()
	r.Register(detector.Word, libreOfficeConverter)
	r.Register(detector.Excel, libreOfficeConverter)
	r.Register(detector.PowerPoint, libreOfficeConverter)
	r.Register(detector.OpenDocument, libreOfficeConverter)
	r.Register(detector.RichText, libreOfficeConverter)
	
	// Register img2pdf for images
	imageConverter := converter.NewImageConverter()
	r.Register(detector.Image, imageConverter)
	
	// Register wkhtmltopdf for HTML
	htmlConverter := converter.NewHTMLConverter()
	r.Register(detector.HTML, htmlConverter)
}

// Register registers a converter for a file type
func (r *ConverterRegistry) Register(fileType detector.FileType, conv converter.Converter) {
	r.converters[fileType] = conv
}

// Get returns the converter for the given file type
func (r *ConverterRegistry) Get(fileType detector.FileType) (converter.Converter, bool) {
	conv, exists := r.converters[fileType]
	return conv, exists
}

// GetAll returns all registered converters
func (r *ConverterRegistry) GetAll() map[detector.FileType]converter.Converter {
	return r.converters
}
