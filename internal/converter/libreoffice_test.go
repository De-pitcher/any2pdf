package converter

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestLibreOfficeConverter_Name(t *testing.T) {
	converter := NewLibreOfficeConverter()
	if converter.Name() != "libreoffice" {
		t.Errorf("expected name 'libreoffice', got '%s'", converter.Name())
	}
}

func TestLibreOfficeConverter_SetTimeout(t *testing.T) {
	converter := NewLibreOfficeConverter()
	
	// Default timeout should be 10 minutes
	if converter.timeout != 10*time.Minute {
		t.Errorf("expected default timeout of 10 minutes, got %v", converter.timeout)
	}

	// Set custom timeout
	customTimeout := 5 * time.Minute
	converter.SetTimeout(customTimeout)
	
	if converter.timeout != customTimeout {
		t.Errorf("expected timeout %v, got %v", customTimeout, converter.timeout)
	}
}

func TestLibreOfficeConverter_SetPreserveTemp(t *testing.T) {
	converter := NewLibreOfficeConverter()
	
	// Default should be false
	if converter.preserveTemp != false {
		t.Error("expected preserveTemp to be false by default")
	}

	converter.SetPreserveTemp(true)
	if converter.preserveTemp != true {
		t.Error("expected preserveTemp to be true after setting")
	}
}

func TestLibreOfficeConverter_IsSupportedFormat(t *testing.T) {
	converter := NewLibreOfficeConverter()

	tests := []struct {
		path     string
		expected bool
	}{
		// Microsoft Office formats
		{"document.docx", true},
		{"document.doc", true},
		{"document.DOCX", true}, // Case insensitive
		{"spreadsheet.xlsx", true},
		{"spreadsheet.xls", true},
		{"presentation.pptx", true},
		{"presentation.ppt", true},
		// OpenDocument formats
		{"document.odt", true},
		{"spreadsheet.ods", true},
		{"presentation.odp", true},
		// Rich Text Format
		{"document.rtf", true},
		// Unsupported
		{"document.pdf", false},
		{"document.txt", false},
		{"image.jpg", false},
		{"noextension", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := converter.IsSupportedFormat(tt.path)
			if result != tt.expected {
				t.Errorf("for %s: expected %v, got %v", tt.path, tt.expected, result)
			}
		})
	}
}

func TestLibreOfficeConverter_FindBinary(t *testing.T) {
	converter := NewLibreOfficeConverter()
	binary := converter.findBinary()
	
	// We can't guarantee LibreOffice is installed, but we can test the logic
	// If found, it should return a path
	if binary != "" {
		if _, err := os.Stat(binary); os.IsNotExist(err) {
			t.Errorf("findBinary returned path that doesn't exist: %s", binary)
		}
	}

	// Test caching
	if binary != "" {
		converter.binaryPath = binary
		secondCall := converter.findBinary()
		if secondCall != binary {
			t.Error("findBinary should return cached value")
		}
	}
}

func TestLibreOfficeConverter_GetFilter(t *testing.T) {
	converter := NewLibreOfficeConverter()

	tests := []struct {
		filename string
		expected string
	}{
		{"document.doc", "MS Word 97"},
		{"spreadsheet.xls", "MS Excel 97"},
		{"presentation.ppt", "MS PowerPoint 97"},
		{"document.rtf", "Rich Text Format"},
		{"document.docx", ""}, // Auto-detect
		{"document.odt", ""},  // Auto-detect
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := converter.getFilter(tt.filename)
			if result != tt.expected {
				t.Errorf("expected filter '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLibreOfficeConverter_BuildQualityArgs(t *testing.T) {
	converter := NewLibreOfficeConverter()

	// Test that quality args can be generated without errors
	qualities := []Quality{Screen, Default, Printer}
	for _, quality := range qualities {
		args := converter.buildQualityArgs(quality)
		// Args are currently empty but function should not panic
		if args == nil {
			t.Errorf("buildQualityArgs returned nil for quality %v", quality)
		}
	}
}

func TestLibreOfficeConverter_GetFormatSpecificArgs(t *testing.T) {
	converter := NewLibreOfficeConverter()

	extensions := []string{".docx", ".xlsx", ".pptx", ".odt", ".ods", ".odp"}
	for _, ext := range extensions {
		args := converter.getFormatSpecificArgs(ext)
		if args == nil {
			t.Errorf("getFormatSpecificArgs returned nil for extension %s", ext)
		}
	}
}

func TestLibreOfficeConverter_HandleConversionError(t *testing.T) {
	converter := NewLibreOfficeConverter()

	tests := []struct {
		name     string
		stderr   string
		contains string
	}{
		{
			name:     "Password protected",
			stderr:   "Error: document is password protected",
			contains: "password-protected",
		},
		{
			name:     "File locked",
			stderr:   "Error: file is locked by another process",
			contains: "locked or in use",
		},
		{
			name:     "Corrupted file",
			stderr:   "Error: document appears to be corrupted",
			contains: "corrupted",
		},
		{
			name:     "Normal error",
			stderr:   "Some other error",
			contains: "Some other error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.handleConversionError(nil, tt.stderr)
			if !contains(result, tt.contains) {
				t.Errorf("expected error message to contain '%s', got: %s", tt.contains, result)
			}
		})
	}
}

func TestLibreOfficeConverter_ValidatePDFCreated(t *testing.T) {
	converter := NewLibreOfficeConverter()

	t.Run("Non-existent file", func(t *testing.T) {
		err := converter.validatePDFCreated("nonexistent.pdf")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})

	t.Run("Empty file", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test-*.pdf")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.Close()

		err = converter.validatePDFCreated(tmpfile.Name())
		if err == nil {
			t.Error("expected error for empty file")
		}
	})

	t.Run("Invalid PDF", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test-*.pdf")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("Not a PDF")
		tmpfile.Close()

		err = converter.validatePDFCreated(tmpfile.Name())
		if err == nil {
			t.Error("expected error for invalid PDF")
		}
	})

	t.Run("Valid PDF", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test-*.pdf")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("%PDF-1.4\n")
		tmpfile.Close()

		err = converter.validatePDFCreated(tmpfile.Name())
		if err != nil {
			t.Errorf("unexpected error for valid PDF: %v", err)
		}
	})
}

func TestLibreOfficeConverter_CopyFile(t *testing.T) {
	converter := NewLibreOfficeConverter()

	// Create source file
	srcFile, err := os.CreateTemp("", "src-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())

	content := "Test content for copy"
	srcFile.WriteString(content)
	srcFile.Close()

	// Copy to destination
	dstFile := filepath.Join(t.TempDir(), "dst.txt")
	err = converter.copyFile(srcFile.Name(), dstFile)
	if err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}
	defer os.Remove(dstFile)

	// Verify content
	data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != content {
		t.Errorf("expected content '%s', got '%s'", content, string(data))
	}
}

func TestLibreOfficeConverter_ConvertBatch_EmptyInput(t *testing.T) {
	converter := NewLibreOfficeConverter()
	
	_, err := converter.ConvertBatch(context.Background(), []string{}, "/tmp", Options{})
	
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}

	if err.Error() != "no input files provided" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestLibreOfficeConverter_GetInstallInstructions(t *testing.T) {
	converter := NewLibreOfficeConverter()
	instructions := converter.GetInstallInstructions()
	
	if instructions == "" {
		t.Error("install instructions should not be empty")
	}

	// Check that instructions mention common platforms
	required := []string{"LibreOffice", "macOS", "Linux", "Windows", "brew", "apt-get", "choco"}
	for _, req := range required {
		if !contains(instructions, req) {
			t.Errorf("install instructions should mention %s", req)
		}
	}

	// Check for headless server instructions
	if !contains(instructions, "headless") {
		t.Error("install instructions should mention headless server setup")
	}
}

func TestLibreOfficeConverter_BuildCommand(t *testing.T) {
	converter := NewLibreOfficeConverter()
	
	// Set a mock binary path
	if runtime.GOOS == "windows" {
		converter.binaryPath = "C:\\Program Files\\LibreOffice\\program\\soffice.exe"
	} else {
		converter.binaryPath = "/usr/bin/libreoffice"
	}

	tests := []struct {
		name       string
		inputPath  string
		outputDir  string
		quality    Quality
	}{
		{
			name:      "Word document",
			inputPath: "document.docx",
			outputDir: "/tmp/output",
			quality:   Default,
		},
		{
			name:      "Excel spreadsheet",
			inputPath: "spreadsheet.xlsx",
			outputDir: "/tmp/output",
			quality:   Printer,
		},
		{
			name:      "PowerPoint presentation",
			inputPath: "presentation.pptx",
			outputDir: "/tmp/output",
			quality:   Screen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{Quality: tt.quality}
			cmd := converter.buildCommand(tt.inputPath, tt.outputDir, opts)

			if cmd == nil {
				t.Fatal("buildCommand returned nil")
			}

			if cmd.Path != converter.binaryPath {
				t.Errorf("expected path %s, got %s", converter.binaryPath, cmd.Path)
			}

			// Check that essential args are present
			argsStr := ""
			for _, arg := range cmd.Args {
				argsStr += arg + " "
			}

			if !contains(argsStr, "--headless") {
				t.Error("command should contain --headless flag")
			}

			if !contains(argsStr, "--convert-to") {
				t.Error("command should contain --convert-to flag")
			}

			if !contains(argsStr, tt.outputDir) {
				t.Error("command should contain output directory")
			}

			if !contains(argsStr, tt.inputPath) {
				t.Error("command should contain input file")
			}
		})
	}
}
