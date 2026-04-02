# Test Fixtures

This directory contains sample files for testing the conversion functionality.

## Directory Structure

```
fixtures/
├── text/
│   └── sample.txt
├── markdown/
│   └── sample.md
├── office/
│   ├── sample.docx
│   ├── sample.xlsx
│   └── sample.pptx
├── images/
│   ├── sample.jpg
│   └── sample.png
└── html/
    └── sample.html
```

## Adding Test Files

When adding test files:
1. Keep files small (< 100KB)
2. Use generic, non-copyrighted content
3. Include various edge cases
4. Document any special characteristics

## Usage

Integration tests use these fixtures to verify actual conversion works with real external tools.
