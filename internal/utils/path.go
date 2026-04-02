package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ValidateInputFile checks if the input file exists and is readable
func ValidateInputFile(path string) error {
	if path == "" {
		return &InvalidPathError{Path: path, Reason: "empty path"}
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return &FileNotFoundError{FilePath: path}
	}
	if err != nil {
		return &InvalidPathError{Path: path, Reason: err.Error()}
	}

	if info.IsDir() {
		return &InvalidPathError{Path: path, Reason: "path is a directory"}
	}

	return nil
}

// GenerateOutputPath creates an output path if one is not specified
func GenerateOutputPath(inputPath, outputPath string) string {
	if outputPath != "" {
		return outputPath
	}

	// Replace extension with .pdf
	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(inputPath, ext)
	return baseName + ".pdf"
}

// GetFileExtension returns the file extension in lowercase without the dot
func GetFileExtension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	// Remove the leading dot and convert to lowercase
	return strings.ToLower(ext[1:])
}

// EnsureAbsolutePath converts a path to absolute if it's relative
func EnsureAbsolutePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", &InvalidPathError{Path: path, Reason: "cannot convert to absolute path"}
	}
	
	return absPath, nil
}

// CreateTempFile creates a temporary file with the given extension
func CreateTempFile(extension string) (string, error) {
	tmpfile, err := os.CreateTemp("", "any2pdf-*"+extension)
	if err != nil {
		return "", err
	}
	
	path := tmpfile.Name()
	tmpfile.Close()
	
	return path, nil
}

// RemoveTempFile removes a temporary file, ignoring errors
func RemoveTempFile(path string) {
	os.Remove(path)
}
