# Office Document Test Fixtures

This directory contains sample Office documents for testing the LibreOffice converter.

## Available Fixtures

- `sample.rtf` - Rich Text Format document (✅ included)

## Missing Fixtures (Add Your Own)

To run the full integration tests, add these files:

### Word Document: `sample.docx`
Create a Word document with:
- Multiple paragraphs
- Bold, italic, and underlined text
- A simple table (2x2)
- A bulleted list
- An image (optional)

Save as `.docx` format.

### Excel Spreadsheet: `sample.xlsx`
Create an Excel file with:
- Multiple sheets (Sheet1, Sheet2)
- Some data and formulas (e.g., `=SUM(A1:A10)`)
- Cell formatting (bold headers, colors)
- A simple chart (optional)

Save as `.xlsx` format.

### PowerPoint Presentation: `sample.pptx`
Create a PowerPoint with:
- 3-4 slides
- Title slide with text
- Content slide with bullet points
- Slide with an image (optional)

Save as `.pptx` format.

### OpenDocument Text: `sample.odt`
Create in LibreOffice Writer or save a Word doc as `.odt` format.

## How to Add Fixtures

1. Create the documents in their respective applications
2. Save them in this directory with the exact filenames above
3. Keep files small (< 100KB each)
4. Use generic, non-confidential content

## Testing Without Fixtures

The integration tests will skip gracefully if fixture files are missing.
Only the RTF test will run by default.

To run all tests once fixtures are added:
```bash
go test ./test -v -run LibreOffice
```

## Creating Minimal Test Documents

### Minimal DOCX (using Python + python-docx):
```python
from docx import Document
doc = Document()
doc.add_heading('Sample Document', 0)
doc.add_paragraph('This is a test document.')
doc.save('sample.docx')
```

### Minimal XLSX (using Python + openpyxl):
```python
from openpyxl import Workbook
wb = Workbook()
ws = wb.active
ws['A1'] = 'Column 1'
ws['B1'] = 'Column 2'
ws['A2'] = 'Data 1'
ws['B2'] = 'Data 2'
wb.save('sample.xlsx')
```

### Minimal PPTX (using Python + python-pptx):
```python
from pptx import Presentation
prs = Presentation()
slide = prs.slides.add_slide(prs.slide_layouts[0])
title = slide.shapes.title
title.text = 'Sample Presentation'
prs.save('sample.pptx')
```

## Note

These fixtures are for testing only. Real-world documents may have more complex
features (macros, embedded objects, etc.) that should be tested separately.
