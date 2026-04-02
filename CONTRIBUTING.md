# Contributing to any2pdf

Thank you for considering contributing to any2pdf! This document provides guidelines and instructions for contributing.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/yourusername/any2pdf/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - System information (OS, Go version, dependency versions)
   - Sample file that causes the issue (if applicable)

### Suggesting Enhancements

1. Check existing issues and pull requests
2. Create an issue describing:
   - The enhancement and its benefits
   - Proposed implementation approach
   - Any breaking changes

### Adding Support for New File Types

To add a new file format:

1. **Choose the right converter**
   - Research existing tools that convert the format to PDF
   - Prefer mature, widely-available tools
   - Document installation requirements

2. **Update the detector**
   - Add new FileType to `internal/detector/types.go`
   - Add extension mapping in `internal/detector/detector.go`

3. **Implement the converter**
   - Create new file in `internal/converter/`
   - Implement the `Converter` interface
   - Handle quality settings appropriately
   - Add comprehensive error handling

4. **Register the converter**
   - Add to registry in `internal/router/registry.go`

5. **Add tests**
   - Unit tests for the converter
   - Integration test with sample file
   - Add sample file to `test/fixtures/`

6. **Update documentation**
   - Add to README supported formats table
   - Document dependency installation
   - Add usage example

### Pull Request Process

1. **Fork and clone** the repository
   ```bash
   git clone https://github.com/yourusername/any2pdf.git
   cd any2pdf
   ```

2. **Create a branch** with a descriptive name
   ```bash
   git checkout -b feature/add-epub-support
   ```

3. **Make your changes**
   - Follow Go conventions and style
   - Add/update tests
   - Update documentation

4. **Test your changes**
   ```bash
   # Run all tests
   go test ./...
   
   # Check formatting
   go fmt ./...
   
   # Run linter
   go vet ./...
   
   # Build to ensure it compiles
   go build ./cmd/any2pdf
   ```

5. **Commit with clear messages**
   ```bash
   git commit -m "feat: add ePub to PDF conversion support"
   ```

6. **Push and create PR**
   ```bash
   git push origin feature/add-epub-support
   ```
   - Create PR on GitHub
   - Reference any related issues
   - Describe what changed and why

7. **Code review**
   - Address reviewer feedback
   - Keep the PR focused and small
   - Be responsive to comments

## Development Setup

### Prerequisites

- Go 1.21 or later
- External tools for testing:
  - pandoc
  - libreoffice
  - img2pdf
  - wkhtmltopdf

### Building

```bash
# Install dependencies
go mod download

# Build
go build -o any2pdf cmd/any2pdf/main.go

# Install locally
go install ./cmd/any2pdf

# Run tests
go test ./... -v

# Test coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Code Style

- Follow standard Go formatting: `go fmt`
- Use meaningful variable and function names
- Add comments for exported functions (godoc)
- Keep functions small and focused
- Handle errors explicitly, don't ignore them
- Use custom error types from `internal/utils/errors.go`

## Testing Guidelines

- Write unit tests for all new functionality
- Aim for 70%+ code coverage
- Use table-driven tests where appropriate
- Mock external dependencies in unit tests
- Include integration tests with real files
- Test error cases, not just happy paths

## Documentation

- Update README.md for user-facing changes
- Update ARCHITECTURE.md for structural changes
- Add godoc comments for exported identifiers
- Include usage examples for new features
- Keep documentation concise and practical

## Questions?

- Open an issue for discussion
- Check existing documentation
- Look at similar converters for examples

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Assume good intentions

Thank you for contributing! 🎉
