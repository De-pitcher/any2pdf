# TODO - any2pdf MVP

## Phase 1: Project Setup
- [x] Create project structure
- [x] Write architecture documentation
- [ ] Initialize Go module: `go mod init github.com/De-pitcher/any2pdf`
- [ ] Create `.gitignore` for Go projects
- [ ] Choose and add LICENSE file (MIT recommended for max adoption)
- [ ] Set up basic CI/CD (GitHub Actions for testing/building)

## Phase 2: Core Infrastructure
- [ ] Implement `internal/utils/errors.go`
  - Custom error types
  - Error wrapping helpers
- [ ] Implement `internal/utils/exec.go`
  - Check if command exists
  - Execute command with timeout
  - Capture stdout/stderr
- [ ] Implement `internal/utils/path.go`
  - Validate input file exists
  - Generate output filename if not specified
  - Temp file creation/cleanup

## Phase 3: File Detection
- [ ] Define `FileType` enum in `internal/detector/types.go`
  - Text, Markdown, Word, Excel, PowerPoint, Image, HTML, Unknown
- [ ] Implement extension-based detection in `internal/detector/detector.go`
  - Map common extensions to FileType
  - Handle case-insensitive extensions
- [ ] Add tests for detector

## Phase 4: Converter Interface & Registry
- [ ] Define `Converter` interface in `internal/converter/converter.go`
  - `Convert(input, output string, opts Options) error`
  - `CheckAvailable() error`
  - `Name() string`
- [ ] Implement converter registry in `internal/router/registry.go`
  - Register converters for file types
  - Lookup converter by FileType
- [ ] Implement router in `internal/router/router.go`
  - Detect file type
  - Find converter
  - Validate converter available

## Phase 5: Converter Implementations
- [ ] **Pandoc converter** (`internal/converter/pandoc.go`)
  - Support .txt, .md
  - Quality mapping
  - Test with sample files
- [ ] **LibreOffice converter** (`internal/converter/libreoffice.go`)
  - Support .docx, .xlsx, .pptx
  - Headless mode
  - Test with sample files
- [ ] **img2pdf converter** (`internal/converter/img2pdf.go`)
  - Support .jpg, .png
  - Quality/DPI settings
  - Test with sample files
- [ ] **wkhtmltopdf converter** (`internal/converter/wkhtmltopdf.go`)
  - Support .html
  - Quality settings
  - Test with sample files

## Phase 6: CLI Implementation
- [ ] Set up CLI framework in `cmd/any2pdf/main.go`
  - Use `urfave/cli/v2`
  - Define flags: `-o`, `-q`, `--quality`
  - Add version flag
- [ ] Implement main conversion flow
  - Parse arguments
  - Call router
  - Handle errors with user-friendly messages
  - Respect quiet flag
- [ ] Add `--check` command to verify dependencies
- [ ] Add help text and examples

## Phase 7: Testing
- [ ] Unit tests for detector (90%+ coverage)
- [ ] Unit tests for utils (90%+ coverage)
- [ ] Unit tests for router
- [ ] Mock tests for converters (test command building)
- [ ] Integration tests with real tools (if available in CI)
- [ ] Create test fixtures directory with sample files
- [ ] Error case testing

## Phase 8: Documentation
- [ ] Complete README with installation instructions
- [ ] Add usage examples
- [ ] Document each supported format
- [ ] Create CONTRIBUTING.md
- [ ] Add code comments/godoc
- [ ] Create dependency installation guide per OS

## Phase 9: Build & Release
- [ ] Create Makefile for common tasks
- [ ] Set up goreleaser config
- [ ] Create GitHub Actions workflow for releases
- [ ] Build for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- [ ] Create install script (curl | sh)
- [ ] Tag v0.1.0 and create first release

## Phase 10: Distribution
- [ ] Publish to GitHub Releases
- [ ] Create Homebrew formula (for macOS)
- [ ] Create installation instructions for package managers
- [ ] Announce on relevant forums/communities

## Future Enhancements (Post-MVP)
- [ ] Batch processing support
- [ ] Configuration file support (~/.any2pdf.yaml)
- [ ] Progress bars for long conversions
- [ ] Parallel processing for multiple files
- [ ] Docker image with all dependencies
- [ ] Web UI/API mode
- [ ] More file formats (ePub, SVG, etc.)
- [ ] Custom converter plugins

## Quick Start Command Sequence

```bash
# 1. Initialize
cd any2pdf
go mod init github.com/De-pitcher/any2pdf
go get github.com/urfave/cli/v2

# 2. First working version
# Implement in order: utils -> detector -> converter interface -> one converter -> CLI

# 3. Test
go test ./...

# 4. Build
go build -o any2pdf cmd/any2pdf/main.go

# 5. Try it
./any2pdf test.txt
```

## MVP Success Criteria

- ✅ Converts at least 3 different format types (.txt, .docx, .jpg)
- ✅ Auto-detects file types
- ✅ Handles errors gracefully
- ✅ Single binary works on Linux/macOS/Windows
- ✅ Clear user documentation
- ✅ 70%+ test coverage
- ✅ Can be installed without touching source code
