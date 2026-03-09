package ingestion

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// MarkdownParser handles Markdown document parsing
type MarkdownParser struct {
	logger *logrus.Logger
}

// NewMarkdownParser creates a new Markdown parser
func NewMarkdownParser(logger *logrus.Logger) *MarkdownParser {
	return &MarkdownParser{
		logger: logger,
	}
}

// ParseMarkdown extracts text content from a Markdown file (simplified version)
func (m *MarkdownParser) ParseMarkdown(ctx context.Context, content []byte) ([]SectionContent, error) {
	text := string(content)
	var sections []SectionContent
	
	// Simple regex-based parsing for now
	lines := strings.Split(text, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Detect headings
		if strings.HasPrefix(line, "#") {
			level := 0
			for _, char := range line {
				if char == '#' {
					level++
				} else {
					break
				}
			}
			title := strings.TrimSpace(line[level:])
			if title != "" {
				sections = append(sections, SectionContent{
					Type:     "heading",
					Level:    level,
					Title:    title,
					Content:  title,
				})
			}
		} else if strings.HasPrefix(line, "```") {
			// Code block
			sections = append(sections, SectionContent{
				Type:     "code_block",
				Content:  line,
			})
		} else if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			// List item
			content := strings.TrimSpace(line[1:])
			if content != "" {
				sections = append(sections, SectionContent{
					Type:    "list_item",
					Content: content,
				})
			}
		} else {
			// Regular paragraph
			sections = append(sections, SectionContent{
				Type:    "paragraph",
				Content: line,
			})
		}
		
		// Check context for cancellation
		select {
		case <-ctx.Done():
			return sections, ctx.Err()
		default:
		}
	}
	
	m.logger.Infof("Extracted %d sections from markdown", len(sections))
	return sections, nil
}

// ConvertToHTML converts markdown to HTML (simplified version)
func (m *MarkdownParser) ConvertToHTML(content []byte) string {
	// For now, just return basic HTML conversion
	text := string(content)
	html := strings.ReplaceAll(text, "\n", "<br>")
	return html
}

// GetMarkdownMetadata extracts metadata from markdown
func (m *MarkdownParser) GetMarkdownMetadata(ctx context.Context, content []byte) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})
	text := string(content)
	
	// Count various elements
	headings := len(regexp.MustCompile(`^#+\s`).FindAllString(text, -1))
	paragraphs := len(strings.Split(text, "\n\n"))
	codeBlocks := len(regexp.MustCompile("```").FindAllString(text, -1)) / 2
	listItems := len(regexp.MustCompile(`^[\-\*]\s`).FindAllString(text, -1))

	metadata["headings"] = headings
	metadata["paragraphs"] = paragraphs
	metadata["code_blocks"] = codeBlocks
	metadata["list_items"] = listItems
	metadata["content_length"] = len(content)

	return metadata, nil
}

// ValidateMarkdown checks if the content is valid markdown
func (m *MarkdownParser) ValidateMarkdown(content []byte) error {
	if len(content) == 0 {
		return fmt.Errorf("empty markdown content")
	}
	return nil
}

// SectionContent represents a section of markdown content
type SectionContent struct {
	Type     string `json:"type"`
	Level    int    `json:"level,omitempty"`
	Title    string `json:"title,omitempty"`
	Language string `json:"language,omitempty"`
	Content  string `json:"content"`
}
