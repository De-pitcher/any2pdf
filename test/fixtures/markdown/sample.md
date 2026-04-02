# Sample Markdown Document

This is a **comprehensive markdown** file for testing the pandoc converter.

## Introduction

This document demonstrates various markdown features that should be properly 
converted to PDF format. The pandoc converter should handle:

- Text formatting (bold, italic, code)
- Lists (ordered and unordered)
- Code blocks with syntax highlighting
- Links and references
- Headers and structure

## Text Formatting

**Bold text** is important.  
*Italic text* provides emphasis.  
***Bold and italic*** for maximum impact.  
`Inline code` for technical terms.

## Lists

### Unordered Lists

- First item
- Second item with a longer description that wraps to multiple lines
- Third item
  - Nested item A
  - Nested item B
    - Deeply nested item
- Fourth item

### Ordered Lists

1. First step
2. Second step
3. Third step with detailed explanation
4. Fourth step

## Code Examples

Here's a Python example:

```python
def fibonacci(n):
    """Calculate the nth Fibonacci number."""
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

print(fibonacci(10))
```

And a Go example:

```go
package main

import "fmt"

func main() {
    for i := 0; i < 5; i++ {
        fmt.Printf("Count: %d\n", i)
    }
}
```

## Links and References

Visit the [Pandoc website](https://pandoc.org) for more information.

Internal reference: See the [Introduction](#introduction) section above.

## Tables

| Feature      | Status      | Priority |
|--------------|-------------|----------|
| Text         | ✓ Complete  | High     |
| Images       | ✓ Complete  | High     |
| Tables       | ✓ Complete  | Medium   |
| Math         | In Progress | Low      |

## Blockquotes

> This is a blockquote that demonstrates how quoted text is rendered.
> 
> It can span multiple paragraphs and preserve formatting.
>
> — Anonymous

## Horizontal Rules

---

Text above the rule.

---

Text below the rule.

## Special Characters and Unicode

This tests Unicode support: こんにちは, 你好, Здравствуй, مرحبا

Special characters: © ® ™ € £ ¥ — – • 

Emoji support: 🎉 🚀 ✅ ❌ 🔍

## Conclusion

This comprehensive markdown file tests the pandoc converter's ability to handle various markup features and generate a well-formatted PDF document.

**End of sample document.**
