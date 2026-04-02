# any2pdf

Convert any file to PDF with a single command.

## Features

- **Simple CLI**: `any2pdf input.docx -o output.pdf`
- **Auto-detection**: Automatically detects file type and uses the right converter
- **Multiple formats**: Supports .txt, .md, .docx, .xlsx, .pptx, .jpg, .png, .html
- **Extensible**: Easy to add support for new file types
- **Cross-platform**: Works on Linux, macOS, and Windows

## Installation

```bash
# Download pre-built binary from releases
# Or build from source:
go build -o any2pdf cmd/any2pdf/main.go
```

## Usage

```bash
# Basic conversion
any2pdf document.docx

# Specify output path
any2pdf document.docx -o /path/to/output.pdf

# Set quality
any2pdf image.jpg --quality printer

# Quiet mode
any2pdf file.xlsx -q
```

## Supported Formats

| Format | Extension | Converter Used |
|--------|-----------|----------------|
| Text | .txt | pandoc |
| Markdown | .md | pandoc |
| Word | .docx | libreoffice |
| Excel | .xlsx | libreoffice |
| PowerPoint | .pptx | libreoffice |
| Images | .jpg, .png | img2pdf |
| HTML | .html | wkhtmltopdf |

## Requirements

The following tools must be installed on your system:

- **pandoc**: For text/markdown conversion
- **libreoffice**: For Office documents (headless mode)
- **img2pdf**: For image conversion
- **wkhtmltopdf**: For HTML conversion

### Installation of dependencies

**Ubuntu/Debian:**
```bash
sudo apt-get install pandoc libreoffice img2pdf wkhtmltopdf
```

**macOS:**
```bash
brew install pandoc
brew install --cask libreoffice
pip3 install img2pdf
brew install wkhtmltopdf
```

**Windows:**
- Download tools from their official websites or use chocolatey/scoop

## Development

```bash
# Run tests
go test ./...

# Build
go build -o any2pdf cmd/any2pdf/main.go

# Install locally
go install ./cmd/any2pdf
```

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
