package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"document.txt", "txt"},
		{"report.DOCX", "docx"},
		{"image.JPG", "jpg"},
		{"/path/to/file.pdf", "pdf"},
		{"noextension", ""},
		{".hidden", "hidden"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetFileExtension(tt.path)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGenerateOutputPath(t *testing.T) {
	tests := []struct {
		input    string
		output   string
		expected string
	}{
		{"document.txt", "", "document.pdf"},
		{"report.docx", "", "report.pdf"},
		{"image.jpg", "custom.pdf", "custom.pdf"},
		{"/path/to/file.md", "", "/path/to/file.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GenerateOutputPath(tt.input, tt.output)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestValidateInputFile(t *testing.T) {
	// Create a temporary test file
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "testdir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorType   string
	}{
		{"Valid file", tmpfile.Name(), false, ""},
		{"Non-existent file", "nonexistent.txt", true, "FileNotFound"},
		{"Directory", tmpdir, true, "InvalidPath"},
		{"Empty path", "", true, "InvalidPath"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInputFile(tt.path)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectError && err != nil {
				if customErr, ok := err.(CustomError); ok {
					if customErr.Type() != tt.errorType {
						t.Errorf("expected error type %s, got %s", tt.errorType, customErr.Type())
					}
				}
			}
		})
	}
}

func TestCreateTempFile(t *testing.T) {
	path, err := CreateTempFile(".pdf")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer RemoveTempFile(path)

	// Check that file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("temp file was not created")
	}

	// Check extension
	if filepath.Ext(path) != ".pdf" {
		t.Errorf("expected .pdf extension, got %s", filepath.Ext(path))
	}
}

func TestEnsureAbsolutePath(t *testing.T) {
	absPath, err := EnsureAbsolutePath("relative/path.txt")
	if err != nil {
		t.Fatalf("failed to convert to absolute path: %v", err)
	}

	if !filepath.IsAbs(absPath) {
		t.Error("expected absolute path")
	}
}
