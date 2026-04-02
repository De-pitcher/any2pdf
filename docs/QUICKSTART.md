# Quick Start Guide

Welcome to any2pdf! This guide will help you get started quickly.

## Installation

```bash
# Download the latest release from GitHub
# Or build from source:
git clone https://github.com/yourusername/any2pdf.git
cd any2pdf
make install
```

## Install Dependencies

Before using any2pdf, you need to install the required external tools:

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install pandoc libreoffice img2pdf wkhtmltopdf texlive-latex-base
```

**macOS:**
```bash
brew install pandoc wkhtmltopdf
brew install --cask libreoffice
pip3 install img2pdf
```

**Check Installation:**
```bash
any2pdf check
```

## Basic Usage

### Convert a file (output name auto-generated)
```bash
any2pdf document.docx
# Creates: document.pdf
```

### Specify output path
```bash
any2pdf report.docx -o /path/to/output.pdf
```

### Set quality level
```bash
any2pdf image.jpg --quality printer
```

### Quiet mode
```bash
any2pdf document.md -q
```

## Examples

```bash
# Text to PDF
any2pdf notes.txt

# Markdown to PDF
any2pdf README.md -o documentation.pdf

# Word document to PDF
any2pdf report.docx

# Excel spreadsheet to PDF
any2pdf data.xlsx

# PowerPoint to PDF
any2pdf slides.pptx

# Image to PDF
any2pdf photo.jpg --quality printer

# HTML to PDF
any2pdf webpage.html
```

## Next Steps

- Read [ARCHITECTURE.md](ARCHITECTURE.md) to understand how it works
- Check [CONTRIBUTING.md](CONTRIBUTING.md) to contribute
- See [TODO.md](TODO.md) for planned features

## Getting Help

```bash
# Show help
any2pdf --help

# Check version
any2pdf --version

# Verify dependencies
any2pdf check
```

## Troubleshooting

**Error: converter not found**
- Run `any2pdf check` to see which dependencies are missing
- Install the required tool from the README

**Error: unsupported file type**
- Check the supported formats list in README
- The file extension must match a supported format

**Conversion fails**
- Verify the input file is valid and not corrupted
- Try opening it in the native application first
- Check the error message for specific details

For more help, open an issue on GitHub.
