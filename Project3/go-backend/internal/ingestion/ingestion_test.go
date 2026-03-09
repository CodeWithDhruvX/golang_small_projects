package ingestion

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownParser_ParseMarkdown(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	parser := NewMarkdownParser(logger)

	tests := []struct {
		name     string
		content  string
		expected int // number of sections expected
	}{
		{
			name:     "simple heading",
			content:  "# Hello World",
			expected: 1,
		},
		{
			name:     "multiple headings",
			content:  "# Title\n## Section 1\n### Subsection",
			expected: 3,
		},
		{
			name:     "paragraphs and headings",
			content:  "# Title\nThis is a paragraph.\n## Section\nAnother paragraph.",
			expected: 4,
		},
		{
			name:     "code blocks",
			content:  "```go\nfunc main() {}\n```",
			expected: 1,
		},
		{
			name:     "list items",
			content:  "- Item 1\n- Item 2\n* Item 3",
			expected: 3,
		},
		{
			name:     "mixed content",
			content:  "# Title\nParagraph here.\n```javascript\nconsole.log('hello');\n```\n- List item\n## Section 2",
			expected: 5,
		},
		{
			name:     "empty content",
			content:  "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sections, err := parser.ParseMarkdown(ctx, []byte(tt.content))
			require.NoError(t, err)
			assert.Len(t, sections, tt.expected)

			// Verify section types
			for _, section := range sections {
				assert.NotEmpty(t, section.Type)
				assert.NotEmpty(t, section.Content)
			}
		})
	}
}

func TestMarkdownParser_ParseMarkdownWithContextCancellation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	content := strings.Repeat("# Heading\nParagraph\n", 100) // Large content

	sections, err := parser.ParseMarkdown(ctx, []byte(content))
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Empty(t, sections)
}

func TestMarkdownParser_ConvertToHTML(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple text",
			content:  "Hello World",
			expected: "Hello World<br>",
		},
		{
			name:     "multiline text",
			content:  "Line 1\nLine 2\nLine 3",
			expected: "Line 1<br>Line 2<br>Line 3<br>",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ConvertToHTML([]byte(tt.content))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarkdownParser_GetMarkdownMetadata(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	content := `# Main Title
This is a paragraph.

## Section 1
Another paragraph.

```go
func main() {
    fmt.Println("Hello")
}
```

- List item 1
- List item 2

### Subsection
Final paragraph.`

	metadata, err := parser.GetMarkdownMetadata(ctx, []byte(content))
	require.NoError(t, err)

	assert.Equal(t, 3, metadata["headings"])    # ##, ###
	assert.Equal(t, 2, metadata["paragraphs"])   // Two paragraphs
	assert.Equal(t, 1, metadata["code_blocks"])  // One code block
	assert.Equal(t, 2, metadata["list_items"])   // Two list items
	assert.Greater(t, metadata["content_length"], 0)
}

func TestMarkdownParser_ValidateMarkdown(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	tests := []struct {
		name      string
		content   string
		expectErr bool
	}{
		{
			name:      "valid content",
			content:   "# Hello World",
			expectErr: false,
		},
		{
			name:      "empty content",
			content:   "",
			expectErr: true,
		},
		{
			name:      "large valid content",
			content:   strings.Repeat("# Heading\nParagraph\n", 100),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.ValidateMarkdown([]byte(tt.content))
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSectionContent(t *testing.T) {
	// Test SectionContent struct
	section := SectionContent{
		Type:     "heading",
		Level:    1,
		Title:    "Test Title",
		Content:  "Test Content",
		Language: "",
	}

	assert.Equal(t, "heading", section.Type)
	assert.Equal(t, 1, section.Level)
	assert.Equal(t, "Test Title", section.Title)
	assert.Equal(t, "Test Content", section.Content)
}

func TestIngestionService_DocumentProcessing(t *testing.T) {
	// Mock ingestion service for testing
	service := &IngestionService{
		logger:        logrus.New(),
		markdownParser: NewMarkdownParser(logrus.New()),
		// Add other mock dependencies as needed
	}

	ctx := context.Background()

	// Test different file types
	tests := []struct {
		name        string
		filename    string
		contentType string
		content     []byte
		expectErr   bool
	}{
		{
			name:        "markdown file",
			filename:    "test.md",
			contentType: "text/markdown",
			content:     []byte("# Test\nThis is a test."),
			expectErr:   false,
		},
		{
			name:        "text file",
			filename:    "test.txt",
			contentType: "text/plain",
			content:     []byte("This is a plain text file."),
			expectErr:   false,
		},
		{
			name:        "empty file",
			filename:    "empty.txt",
			contentType: "text/plain",
			content:     []byte(""),
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would test the actual ingestion logic
			// For now, we'll test the markdown parsing part
			if tt.contentType == "text/markdown" {
				sections, err := service.markdownParser.ParseMarkdown(ctx, tt.content)
				if tt.expectErr && len(tt.content) == 0 {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotEmpty(t, sections)
				}
			}
		})
	}
}

func TestMarkdownParser_HeadingLevels(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	content := `# Level 1
## Level 2
### Level 3
#### Level 4
##### Level 5
###### Level 6`

	sections, err := parser.ParseMarkdown(ctx, []byte(content))
	require.NoError(t, err)
	assert.Len(t, sections, 6)

	// Verify heading levels
	expectedLevels := []int{1, 2, 3, 4, 5, 6}
	for i, section := range sections {
		assert.Equal(t, "heading", section.Type)
		assert.Equal(t, expectedLevels[i], section.Level)
	}
}

func TestMarkdownParser_CodeBlockDetection(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	content := "```go\nfunc main() {}\n```"

	sections, err := parser.ParseMarkdown(ctx, []byte(content))
	require.NoError(t, err)
	assert.Len(t, sections, 1)

	section := sections[0]
	assert.Equal(t, "code_block", section.Type)
	assert.Equal(t, "```go", section.Content)
}

func TestMarkdownParser_ListItemDetection(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "dash list items",
			content:  "- Item 1\n- Item 2\n- Item 3",
			expected: 3,
		},
		{
			name:     "asterisk list items",
			content:  "* Item 1\n* Item 2",
			expected: 2,
		},
		{
			name:     "mixed list items",
			content:  "- Item 1\n* Item 2\n- Item 3",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sections, err := parser.ParseMarkdown(ctx, []byte(tt.content))
			require.NoError(t, err)
			assert.Len(t, sections, tt.expected)

			for _, section := range sections {
				assert.Equal(t, "list_item", section.Type)
			}
		})
	}
}

func TestMarkdownParser_ParagraphDetection(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	content := "This is a paragraph.\nThis is another paragraph."

	sections, err := parser.ParseMarkdown(ctx, []byte(content))
	require.NoError(t, err)
	assert.Len(t, sections, 2)

	for _, section := range sections {
		assert.Equal(t, "paragraph", section.Type)
		assert.NotEmpty(t, section.Content)
	}
}

func TestMarkdownParser_ComplexDocument(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)

	content := `# Document Title

This is the introduction paragraph.

## Getting Started

Here's some code:

```python
def hello():
    print("Hello, World!")
```

### Installation Steps

1. First step
2. Second step
3. Third step

## Usage

- Use the command line
- Follow the documentation
- Enjoy the tool

### Advanced Features

For advanced users, there are additional features available.`

	sections, err := parser.ParseMarkdown(ctx, []byte(content))
	require.NoError(t, err)

	// Count different types
	var headings, paragraphs, codeBlocks, listItems int
	for _, section := range sections {
		switch section.Type {
		case "heading":
			headings++
		case "paragraph":
			paragraphs++
		case "code_block":
			codeBlocks++
		case "list_item":
			listItems++
		}
	}

	assert.Equal(t, 5, headings)   # # ## ### ### ##
	assert.Equal(t, 3, paragraphs)  // intro, usage, advanced
	assert.Equal(t, 1, codeBlocks)  // python code
	assert.Equal(t, 3, listItems)   // numbered list items

	// Test metadata extraction
	metadata, err := parser.GetMarkdownMetadata(ctx, []byte(content))
	require.NoError(t, err)

	assert.Equal(t, headings, metadata["headings"])
	assert.Equal(t, paragraphs, metadata["paragraphs"])
	assert.Equal(t, codeBlocks, metadata["code_blocks"])
	assert.Equal(t, listItems, metadata["list_items"])
}

func BenchmarkMarkdownParser_ParseMarkdown(b *testing.B) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)
	content := strings.Repeat("# Heading\nParagraph content here.\n", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseMarkdown(ctx, []byte(content))
	}
}

func BenchmarkMarkdownParser_GetMarkdownMetadata(b *testing.B) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	parser := NewMarkdownParser(logger)
	content := strings.Repeat("# Heading\nParagraph content here.\n```go\ncode\n```\n- item\n", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.GetMarkdownMetadata(ctx, []byte(content))
	}
}
