package utils

import "fmt"

// CustomError is the base interface for all custom errors
type CustomError interface {
	error
	Type() string
}

// UnsupportedFileTypeError indicates the file type is not supported
type UnsupportedFileTypeError struct {
	Extension string
	FilePath  string
}

func (e *UnsupportedFileTypeError) Error() string {
	return fmt.Sprintf("unsupported file type: %s (file: %s)", e.Extension, e.FilePath)
}

func (e *UnsupportedFileTypeError) Type() string {
	return "UnsupportedFileType"
}

// ConverterNotFoundError indicates no converter is available for the file type
type ConverterNotFoundError struct {
	FileType     string
	MissingTool  string
	InstallGuide string
}

func (e *ConverterNotFoundError) Error() string {
	msg := fmt.Sprintf("converter not found for file type: %s", e.FileType)
	if e.MissingTool != "" {
		msg += fmt.Sprintf("\nRequired tool: %s", e.MissingTool)
	}
	if e.InstallGuide != "" {
		msg += fmt.Sprintf("\nInstallation: %s", e.InstallGuide)
	}
	return msg
}

func (e *ConverterNotFoundError) Type() string {
	return "ConverterNotFound"
}

// ConversionFailedError indicates the conversion process failed
type ConversionFailedError struct {
	InputFile  string
	OutputFile string
	Converter  string
	ExitCode   int
	Stderr     string
	Cause      error
}

func (e *ConversionFailedError) Error() string {
	msg := fmt.Sprintf("conversion failed: %s -> %s (using %s)", e.InputFile, e.OutputFile, e.Converter)
	if e.ExitCode > 0 {
		msg += fmt.Sprintf("\nExit code: %d", e.ExitCode)
	}
	if e.Stderr != "" {
		msg += fmt.Sprintf("\nError output: %s", e.Stderr)
	}
	if e.Cause != nil {
		msg += fmt.Sprintf("\nCause: %v", e.Cause)
	}
	return msg
}

func (e *ConversionFailedError) Type() string {
	return "ConversionFailed"
}

func (e *ConversionFailedError) Unwrap() error {
	return e.Cause
}

// FileNotFoundError indicates the input file doesn't exist
type FileNotFoundError struct {
	FilePath string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found: %s", e.FilePath)
}

func (e *FileNotFoundError) Type() string {
	return "FileNotFound"
}

// InvalidPathError indicates an invalid file path
type InvalidPathError struct {
	Path   string
	Reason string
}

func (e *InvalidPathError) Error() string {
	return fmt.Sprintf("invalid path: %s (%s)", e.Path, e.Reason)
}

func (e *InvalidPathError) Type() string {
	return "InvalidPath"
}
