# Setup Guide for Windows

This guide will help you set up the development environment for any2pdf on Windows.

## Step 1: Install Go

### Option A: Using winget (Recommended)
```powershell
winget install GoLang.Go
```

### Option B: Using Chocolatey
```powershell
choco install golang
```

### Option C: Manual Installation
1. Download Go from https://go.dev/dl/
2. Run the installer (go1.21.x.windows-amd64.msi or later)
3. Follow the installation wizard
4. Restart PowerShell/Terminal

### Verify Installation
```powershell
go version
# Should output: go version go1.21.x windows/amd64
```

## Step 2: Initialize Go Module

The module has already been configured, just run:

```powershell
go mod download
```

This will download all dependencies.

## Step 3: Build the Project

```powershell
# Using Make (if you have Make installed)
make build

# Or directly with Go
go build -o build/any2pdf.exe cmd/any2pdf/main.go
```

## Step 4: Install External Dependencies

any2pdf requires external tools to perform conversions:

### Using Chocolatey (Easiest)
```powershell
# Install Chocolatey first if you don't have it
# Visit https://chocolatey.org/install

# Install dependencies
choco install pandoc
choco install libreoffice
choco install python3
pip3 install img2pdf
choco install wkhtmltopdf
```

### Manual Installation
- **Pandoc**: https://pandoc.org/installing.html
- **LibreOffice**: https://www.libreoffice.org/download/download/
- **img2pdf**: `pip3 install img2pdf` (requires Python)
- **wkhtmltopdf**: https://wkhtmltopdf.org/downloads.html

### Verify Dependencies
```powershell
./build/any2pdf.exe check
```

## Step 5: Run Tests

```powershell
go test ./...
```

## Step 6: Try It Out

```powershell
# Convert a sample file
./build/any2pdf.exe test/fixtures/text/sample.txt

# Should create test/fixtures/text/sample.pdf
```

## Troubleshooting

### "go: command not found"
- Restart your terminal after installing Go
- Check PATH with: `$env:Path`
- Go should be in `C:\Program Files\Go\bin`

### Module errors
```powershell
# Clean module cache and retry
go clean -modcache
go mod download
```

### Build errors
```powershell
# Update dependencies
go mod tidy

# Rebuild
go build -o build/any2pdf.exe cmd/any2pdf/main.go
```

## Optional: Install Make for Windows

To use the Makefile:

### Using Chocolatey
```powershell
choco install make
```

### Using Scoop
```powershell
scoop install make
```

Or use Go commands directly instead of make commands.

## Next Steps

Once everything is set up:
1. Run `go test ./...` to verify tests pass
2. Try converting different file types
3. Check [CONTRIBUTING.md](CONTRIBUTING.md) to start contributing
4. See [TODO.md](TODO.md) for features to implement

## Quick Reference

```powershell
# Build
go build -o build/any2pdf.exe cmd/any2pdf/main.go

# Test
go test ./...

# Run
./build/any2pdf.exe <file>

# Install globally
go install ./cmd/any2pdf
# Then use: any2pdf <file>
```
