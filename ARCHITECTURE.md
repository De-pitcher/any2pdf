# Architecture

## Overview

any2pdf follows a simple, modular architecture with clear separation of concerns.

```
┌─────────────┐
│     CLI     │  Parse args, handle I/O
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Detector   │  Identify file type
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Router    │  Select appropriate converter
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Converters  │  Execute conversion (pandoc, libreoffice, etc.)
└─────────────┘
```

## Core Modules

### 1. CLI (`cmd/any2pdf/main.go`)
- Entry point
- Argument parsing with `urfave/cli` or `cobra`
- User output and error display
- Delegates to internal packages

### 2. Detector (`internal/detector/`)
- Detects file type from extension
- Falls back to MIME type detection if needed
- Returns `FileType` enum

### 3. Router (`internal/router/`)
- Maps `FileType` to appropriate `Converter`
- Registry pattern for converter lookup
- Validates converter availability

### 4. Converters (`internal/converter/`)
- Interface: `type Converter interface { Convert(input, output string, opts Options) error }`
- Implementations:
  - `pandoc.go` - Text/Markdown
  - `libreoffice.go` - Office documents
  - `img2pdf.go` - Images
  - `wkhtmltopdf.go` - HTML
- Each converter:
  - Checks if external tool is installed
  - Builds command with appropriate flags
  - Executes and streams output
  - Returns structured errors

### 5. Utils (`internal/utils/`)
- Path validation
- Executable checking
- Temp file management
- Error types

## Design Principles

### Extensibility
New file types can be added by:
1. Adding entry to `detector/types.go`
2. Implementing `Converter` interface
3. Registering in `router/registry.go`

No changes to CLI or core logic needed.

### Error Handling
- Custom error types for:
  - `UnsupportedFileType`
  - `ConverterNotFound`
  - `ConversionFailed`
- Errors include context (file, converter used, exit code)
- User-friendly messages vs debug info (controlled by quiet flag)

### Dependency Management
- External tools (pandoc, etc.) are NOT bundled
- Check availability at runtime with helpful error messages
- Guide users to installation instructions
- Future: Optional Docker mode to bundle everything

### Quality Levels
Map to converter-specific settings:
- `screen`: Lower DPI, smaller files (75 DPI)
- `default`: Balanced (150 DPI)
- `printer`: High quality (300 DPI)

## File Structure

```
any2pdf/
├── cmd/
│   └── any2pdf/
│       └── main.go              # CLI entry point
├── internal/
│   ├── detector/
│   │   ├── detector.go          # File type detection
│   │   └── types.go             # FileType enum
│   ├── router/
│   │   ├── router.go            # Converter selection
│   │   └── registry.go          # Converter registry
│   ├── converter/
│   │   ├── converter.go         # Interface & common code
│   │   ├── pandoc.go            # Pandoc converter
│   │   ├── libreoffice.go       # LibreOffice converter
│   │   ├── img2pdf.go           # Image converter
│   │   └── wkhtmltopdf.go       # HTML converter
│   └── utils/
│       ├── exec.go              # Execute external commands
│       ├── path.go              # Path utilities
│       └── errors.go            # Custom error types
├── pkg/
│   └── any2pdf/
│       └── any2pdf.go           # Public API (if library use is desired)
├── test/
│   ├── fixtures/                # Test files
│   └── integration_test.go      # Integration tests
├── go.mod
├── go.sum
├── README.md
├── ARCHITECTURE.md
├── LICENSE
└── .gitignore
```

## Technology Choices

### Language: Go
- Single binary compilation
- Excellent cross-platform support
- Simple concurrency (if we add batch processing)
- Strong standard library for CLI tools
- Fast startup time

### CLI Framework: urfave/cli v2
- Simple, idiomatic API
- Good flag handling
- Built-in help generation
- Lightweight (alternative: cobra, but heavier)

### External Tools Strategy
- Shell out to existing mature tools
- Don't reinvent the wheel
- Each tool is battle-tested
- Smaller binary size
- Trade-off: Requires installation

## Future Enhancements (Post-MVP)

1. **Batch processing**: `any2pdf *.docx`
2. **Config file**: `~/.any2pdf.yaml` for defaults
3. **Docker mode**: Bundled with all dependencies
4. **Watch mode**: Auto-convert on file change
5. **Plugin system**: Load custom converters
6. **Web service**: HTTP API for conversion
7. **More formats**: ePub, SVG, CAD files, etc.
