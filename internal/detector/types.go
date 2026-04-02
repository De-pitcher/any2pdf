package detector

// FileType represents the type of file to be converted
type FileType int

const (
	Unknown FileType = iota
	Text
	Markdown
	ReStructuredText
	OrgMode
	Textile
	Word
	Excel
	PowerPoint
	OpenDocument
	RichText
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
	case ReStructuredText:
		return "ReStructuredText"
	case OrgMode:
		return "OrgMode"
	case Textile:
		return "Textile"
	case Word:
		return "Word"
	case Excel:
		return "Excel"
	case PowerPoint:
		return "PowerPoint"
	case OpenDocument:
		return "OpenDocument"
	case RichText:
		return "RichText"
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
	
	// Other text markup formats
	"rst":       ReStructuredText,
	"org":       OrgMode,
	"textile":   Textile,
	"mediawiki": Markdown, // Treat as markdown variant for routing
	
	// Microsoft Word formats
	"doc":  Word,
	"docx": Word,
	
	// Microsoft Excel formats
	"xls":  Excel,
	"xlsx": Excel,
	
	// Microsoft PowerPoint formats
	"ppt":  PowerPoint,
	"pptx": PowerPoint,
	
	// OpenDocument formats
	"odt": OpenDocument, // Writer (word processing)
	"ods": OpenDocument, // Calc (spreadsheet)
	"odp": OpenDocument, // Impress (presentation)
	
	// Rich Text Format
	"rtf": RichText,
	
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
