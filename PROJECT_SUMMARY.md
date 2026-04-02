# any2pdf Project Summary

## Overview

**any2pdf** is a CLI tool that converts various file formats to PDF using a modular, extensible architecture.

**Technology Stack:** Go 1.21+  
**License:** MIT  
**Architecture:** Modular with clean separation of concerns

---

## Project Structure

```
any2pdf/
├── cmd/
│   └── any2pdf/
│       └── main.go              # CLI entry point with urfave/cli
│
├── internal/
│   ├── detector/                # File type detection
│   │   ├── detector.go          # Detection logic
│   │   ├── detector_test.go     # Unit tests
│   │   └── types.go             # FileType enum & extension map
│   │
│   ├── router/                  # Converter routing
│   │   ├── router.go            # Main routing logic
│   │   └── registry.go          # Converter registry
│   │
│   ├── converter/               # Converter implementations
│   │   ├── converter.go         # Interface & Options
│   │   ├── converter_test.go    # Tests
│   │   ├── pandoc.go            # Text/Markdown → PDF
│   │   ├── libreoffice.go       # Office docs → PDF
│   │   ├── img2pdf.go           # Images → PDF
│   │   └── wkhtmltopdf.go       # HTML → PDF
│   │
│   └── utils/                   # Utilities
│       ├── errors.go            # Custom error types
│       ├── exec.go              # Command execution
│       ├── path.go              # Path utilities
│       └── path_test.go         # Tests
│
├── test/
│   └── fixtures/                # Test files
│       ├── text/sample.txt
│       ├── markdown/sample.md
│       ├── html/sample.html
│       └── README.md
│
├── docs/
│   └── QUICKSTART.md            # Quick start guide
│
├── README.md                    # Main documentation
├── ARCHITECTURE.md              # Architecture details
├── CONTRIBUTING.md              # Contribution guidelines
├── TODO.md                      # Development roadmap
├── LICENSE                      # MIT License
├── Makefile                     # Build automation
├── .gitignore                   # Git ignore rules
└── go.mod                       # Go module definition
```

---

## Core Components

### 1. **CLI (cmd/any2pdf/main.go)**
- Entry point using urfave/cli framework
- Parses arguments: `-o/--output`, `-q/--quiet`, `--quality`
- Commands: `convert` (default), `check` (dependency check)
- Wires together detector → router → converter

### 2. **Detector (internal/detector/)**
- Maps file extensions to FileType enum
- Supports: .txt, .md, .docx, .xlsx, .pptx, .jpg, .png, .html
- Returns `UnsupportedFileTypeError` for unknown types

### 3. **Router (internal/router/)**
- Registry pattern for converter lookup
- Routes FileType to appropriate Converter
- Validates converter availability before use

### 4. **Converters (internal/converter/)**
All implement the `Converter` interface:
```go
type Converter interface {
    Convert(inputPath, outputPath string, opts Options) error
    CheckAvailable() error
    Name() string
}
```

**Implementations:**
- `PandocConverter`: Text, Markdown (requires: pandoc)
- `LibreOfficeConverter`: Word, Excel, PowerPoint (requires: libreoffice)
- `ImageConverter`: JPG, PNG, etc. (requires: img2pdf)
- `HTMLConverter`: HTML files (requires: wkhtmltopdf)

### 5. **Utils (internal/utils/)**
- **errors.go**: Custom error types (UnsupportedFileType, ConverterNotFound, ConversionFailed, etc.)
- **exec.go**: Command execution with timeout support
- **path.go**: Path validation, extension extraction, output generation

---

## Quality Levels

- **screen**: 75 DPI, smaller files
- **default**: 150 DPI, balanced
- **printer**: 300 DPI, high quality

Mapped to converter-specific flags.

---

## Error Handling

Custom error types provide context:
- `UnsupportedFileTypeError` → file extension not supported
- `ConverterNotFoundError` → external tool missing (with install guide)
- `ConversionFailedError` → conversion failed (includes stderr, exit code)
- `FileNotFoundError` → input file doesn't exist
- `InvalidPathError` → invalid path

---

## Extensibility

### Adding a New File Format

1. **Add to detector**: Update `ExtensionMap` in `internal/detector/types.go`
2. **Create converter**: Implement `Converter` interface in `internal/converter/`
3. **Register**: Add to registry in `internal/router/registry.go`
4. **Test**: Add unit tests and sample fixture file
5. **Document**: Update README supported formats table

**Example**: Adding ePub support
```go
// 1. Add FileType
const ePub FileType = ...

// 2. Update ExtensionMap
ExtensionMap["epub"] = ePub

// 3. Create converter
type ePubConverter struct{}
func (c *ePubConverter) Convert(...) error { ... }

// 4. Register
registry.Register(detector.ePub, NewePubConverter())
```

---

## Next Steps (MVP Development)

### Phase 1: Initialize & Build (Current Phase ✓)
- [x] Project structure created
- [x] Core modules implemented
- [x] Tests scaffolded
- [ ] Initialize Go module: `go mod init github.com/yourusername/any2pdf`
- [ ] Install dependencies: `go mod download`

### Phase 2: Build & Test
```bash
# Initialize module (update with your GitHub username)
go mod init github.com/yourusername/any2pdf
go mod tidy

# Run tests
go test ./...

# Build
make build

# Try it (after installing dependencies)
./build/any2pdf test/fixtures/text/sample.txt
```

### Phase 3: Install Dependencies
```bash
# Ubuntu/Debian
sudo apt-get install pandoc libreoffice img2pdf wkhtmltopdf

# macOS
brew install pandoc wkhtmltopdf
brew install --cask libreoffice
pip3 install img2pdf

# Verify
./build/any2pdf check
```

### Phase 4: Complete Implementation
- Fix LibreOffice output file renaming
- Add integration tests
- Test with real files
- Fix edge cases

### Phase 5: Release
- Set up GitHub Actions for CI/CD
- Create goreleaser config
- Build for multiple platforms
- Create first release (v0.1.0)

---

## Design Principles

✅ **Practical over perfect**: Uses mature external tools rather than reinventing  
✅ **Extensible**: Plugin-like converter system  
✅ **User-friendly errors**: Clear messages with installation guides  
✅ **Single binary**: Easy distribution, no runtime dependencies  
✅ **Testable**: Mockable interfaces, table-driven tests  
✅ **Cross-platform**: Works on Linux, macOS, Windows  

---

## Dependencies

### Build Dependencies
- Go 1.21+
- urfave/cli/v2 (CLI framework)

### Runtime Dependencies (External Tools)
- **pandoc**: Text/Markdown conversion
- **libreoffice**: Office documents
- **img2pdf**: Image conversion
- **wkhtmltopdf**: HTML conversion

Users install these separately. The tool checks availability at runtime.

---

## Usage Examples

```bash
# Basic conversion
any2pdf document.docx

# Custom output
any2pdf report.md -o final-report.pdf

# High quality image
any2pdf photo.jpg --quality printer

# Quiet mode
any2pdf data.xlsx -q

# Check dependencies
any2pdf check

# Help
any2pdf --help
```

---

## Current Status

✅ **Complete**: Project structure, architecture, skeleton code, tests, documentation  
⚠️ **Next**: Initialize Go module, run tests, build binary  
⏳ **Future**: Integration testing, CI/CD, first release  

---

## Commands Reference

```bash
# Development
make help           # Show all commands
make build          # Build binary
make test           # Run tests
make test-coverage  # Generate coverage report
make clean          # Clean build artifacts
make install        # Install to $GOPATH/bin
make check-deps     # Verify external tools installed

# First time setup
go mod init github.com/YOURNAME/any2pdf
go mod tidy
make build
```

---

**Ready to start coding!** 🚀

See [TODO.md](TODO.md) for detailed roadmap.
