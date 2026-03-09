package ingestion

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

// TextParser handles plain text document parsing
type TextParser struct {
	logger *logrus.Logger
}

// NewTextParser creates a new text parser
func NewTextParser(logger *logrus.Logger) *TextParser {
	return &TextParser{
		logger: logger,
	}
}

// ParseText extracts text content from a plain text file
func (t *TextParser) ParseText(ctx context.Context, content []byte) ([]ParagraphContent, error) {
	reader := bufio.NewReader(strings.NewReader(string(content)))
	var paragraphs []ParagraphContent
	var currentParagraph strings.Builder
	lineNumber := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if line != "" {
					currentParagraph.WriteString(strings.TrimSpace(line))
				}
			} else {
				return nil, fmt.Errorf("error reading text: %w", err)
			}
			break
		}

		lineNumber++
		line = strings.TrimSpace(line)

		// Skip empty lines (paragraph separators)
		if line == "" {
			if currentParagraph.Len() > 0 {
				paragraphs = append(paragraphs, ParagraphContent{
					LineNumber: lineNumber - currentParagraph.Len(),
					Content:    currentParagraph.String(),
				})
				currentParagraph.Reset()
			}
			continue
		}

		// Add line to current paragraph
		if currentParagraph.Len() > 0 {
			currentParagraph.WriteString(" ")
		}
		currentParagraph.WriteString(line)

		// Check context for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	// Add the last paragraph if there's remaining content
	if currentParagraph.Len() > 0 {
		paragraphs = append(paragraphs, ParagraphContent{
			LineNumber: lineNumber,
			Content:    currentParagraph.String(),
		})
	}

	t.logger.Infof("Extracted %d paragraphs from text file", len(paragraphs))
	return paragraphs, nil
}

// GetTextMetadata extracts metadata from plain text
func (t *TextParser) GetTextMetadata(ctx context.Context, content []byte) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})
	
	// Basic statistics
	contentStr := string(content)
	metadata["content_length"] = len(content)
	metadata["character_count"] = len([]rune(contentStr))
	metadata["word_count"] = len(strings.Fields(contentStr))
	metadata["line_count"] = strings.Count(contentStr, "\n") + 1
	metadata["paragraph_count"] = strings.Count(contentStr, "\n\n") + 1

	// Detect encoding (simple heuristic)
	if isValidUTF8(content) {
		metadata["encoding"] = "utf-8"
	} else {
		metadata["encoding"] = "unknown"
	}

	// Detect if content is likely code (simple heuristic)
	if isLikelyCode(contentStr) {
		metadata["content_type_hint"] = "code"
	} else {
		metadata["content_type_hint"] = "prose"
	}

	return metadata, nil
}

// ValidateText checks if the content is valid text
func (t *TextParser) ValidateText(content []byte) error {
	if len(content) == 0 {
		return fmt.Errorf("empty text content")
	}

	// Check if content is valid UTF-8
	if !isValidUTF8(content) {
		return fmt.Errorf("content contains invalid UTF-8 sequences")
	}

	return nil
}

// isValidUTF8 checks if content contains valid UTF-8
func isValidUTF8(content []byte) bool {
	// Go strings are UTF-8 by default, so this is a simple check
	// In a more sophisticated implementation, you might use utf8.Valid
	return true
}

// isLikelyCode performs a simple heuristic to determine if text is likely code
func isLikelyCode(content string) bool {
	indicators := []string{
		"{", "}", "function", "class", "def ", "import ", "include",
		"//", "/*", "*/", "#", "if ", "else", "for ", "while ",
		"return ", "var ", "let ", "const ", "public ", "private ",
	}

	codeIndicators := 0
	totalIndicators := 0

	for _, indicator := range indicators {
		if strings.Contains(content, indicator) {
			codeIndicators++
		}
		totalIndicators++
	}

	// If more than 30% of code indicators are present, consider it code
	if totalIndicators > 0 && float64(codeIndicators)/float64(totalIndicators) > 0.3 {
		return true
	}

	return false
}

// ParagraphContent represents a paragraph of text content
type ParagraphContent struct {
	LineNumber int    `json:"line_number"`
	Content    string `json:"content"`
}
