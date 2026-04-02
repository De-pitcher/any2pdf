package detector

// FileType represents the type of file to be converted
type FileType int

const (
	Unknown FileType = iota
	Text
	Markdown
	Word
	Excel
	PowerPoint
	Image
	HTML
)

// String returns the string representation of FileType
func (ft FileType) String() string {
	switch ft {
	case Text:
		return "Text"
	case Markdown:
		return "Markdown"
	case Word:
		return "Word"
	case Excel:
		return "Excel"
	case PowerPoint:
		return "PowerPoint"
	case Image:
		return "Image"
	case HTML:
		return "HTML"
	default:
		return "Unknown"
	}
}

// ExtensionMap maps file extensions to FileTypes
var ExtensionMap = map[string]FileType{
	// Text formats
	"txt": Text,
	
	// Markdown formats
	"md":       Markdown,
	"markdown": Markdown,
	
	// Microsoft Word formats
	"doc":  Word,
	"docx": Word,
	
	// Microsoft Excel formats
	"xls":  Excel,
	"xlsx": Excel,
	
	// Microsoft PowerPoint formats
	"ppt":  PowerPoint,
	"pptx": PowerPoint,
	
	// Image formats
	"jpg":  Image,
	"jpeg": Image,
	"png":  Image,
	"gif":  Image,
	"bmp":  Image,
	"tiff": Image,
	"tif":  Image,
	
	// HTML formats
	"html": HTML,
	"htm":  HTML,
}
