package detector

import (
	"testing"
)

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		expectedType FileType
		expectError  bool
	}{
		{"Text file", "document.txt", Text, false},
		{"Markdown file", "readme.md", Markdown, false},
		{"Word file", "report.docx", Word, false},
		{"Excel file", "spreadsheet.xlsx", Excel, false},
		{"PowerPoint file", "presentation.pptx", PowerPoint, false},
		{"JPEG image", "photo.jpg", Image, false},
		{"PNG image", "screenshot.png", Image, false},
		{"HTML file", "page.html", HTML, false},
		{"Unsupported file", "video.mp4", Unknown, true},
		{"No extension", "file", Unknown, true},
	}

	d := NewDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileType, err := d.DetectFileType(tt.filePath)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if fileType != tt.expectedType {
				t.Errorf("expected %v, got %v", tt.expectedType, fileType)
			}
		})
	}
}

func TestIsSupportedExtension(t *testing.T) {
	tests := []struct {
		extension string
		expected  bool
	}{
		{"txt", true},
		{"docx", true},
		{"jpg", true},
		{"mp4", false},
		{"exe", false},
	}

	d := NewDetector()

	for _, tt := range tests {
		t.Run(tt.extension, func(t *testing.T) {
			result := d.IsSupportedExtension(tt.extension)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetSupportedExtensions(t *testing.T) {
	d := NewDetector()
	extensions := d.GetSupportedExtensions()

	if len(extensions) == 0 {
		t.Error("expected non-empty list of extensions")
	}

	// Check for some required extensions
	required := []string{"txt", "md", "docx", "jpg", "html"}
	extensionMap := make(map[string]bool)
	for _, ext := range extensions {
		extensionMap[ext] = true
	}

	for _, req := range required {
		if !extensionMap[req] {
			t.Errorf("expected extension %s to be supported", req)
		}
	}
}
