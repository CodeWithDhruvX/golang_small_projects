package ingestion

import (
	"bytes"
	"context"
	"fmt"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/sirupsen/logrus"
)

// PDFParser handles PDF document parsing
type PDFParser struct {
	logger *logrus.Logger
}

// NewPDFParser creates a new PDF parser
func NewPDFParser(logger *logrus.Logger) *PDFParser {
	// Set unidoc license (community edition)
	// license.SetLicense("Community License") // Commented out for now
	
	return &PDFParser{
		logger: logger,
	}
}

// ParsePDF extracts text content from a PDF file
func (p *PDFParser) ParsePDF(ctx context.Context, content []byte) ([]PageContent, error) {
	// Create PDF reader
	pdfReader, err := model.NewPdfReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	// Get number of pages
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get number of pages: %w", err)
	}

	var pages []PageContent
	
	// Extract text from each page
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			p.logger.Warnf("Failed to get page %d: %v", pageNum, err)
			continue
		}

		// Extract text from page
		extractor, err := extractor.New(page)
		if err != nil {
			p.logger.Warnf("Failed to create extractor for page %d: %v", pageNum, err)
			continue
		}

		text, err := extractor.ExtractText()
		if err != nil {
			p.logger.Warnf("Failed to extract text from page %d: %v", pageNum, err)
			continue
		}

		if text != "" {
			pages = append(pages, PageContent{
				PageNumber: pageNum,
				Content:    text,
			})
		}

		// Check context for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	p.logger.Infof("Extracted text from %d pages", len(pages))
	return pages, nil
}

// GetPDFMetadata extracts metadata from PDF
func (p *PDFParser) GetPDFMetadata(ctx context.Context, content []byte) (map[string]interface{}, error) {
	pdfReader, err := model.NewPdfReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	// Get PDF info - temporarily commented out due to API change
	// pdfInfo, err := pdfReader.GetPdfInfo()
	// if err != nil {
	// 	p.logger.Warnf("Failed to get PDF info: %v", err)
	// 	return map[string]interface{}{}, nil
	// }

	metadata := make(map[string]interface{})

	// Extract basic metadata - temporarily commented out
	// if pdfInfo.Title != "" {
	// 	metadata["title"] = pdfInfo.Title
	// }
	// if pdfInfo.Author != "" {
	// 	metadata["author"] = pdfInfo.Author
	// }
	// if pdfInfo.Subject != "" {
	// 	metadata["subject"] = pdfInfo.Subject
	// }
	// if pdfInfo.Creator != "" {
	// 	metadata["creator"] = pdfInfo.Creator
	// }
	// if pdfInfo.Producer != "" {
	// 	metadata["producer"] = pdfInfo.Producer
	// }

	// Get page count
	numPages, err := pdfReader.GetNumPages()
	if err == nil {
		metadata["page_count"] = numPages
	}

	return metadata, nil
}

// ValidatePDF checks if the content is a valid PDF
func (p *PDFParser) ValidatePDF(content []byte) error {
	// Try to create PDF reader to validate
	_, err := model.NewPdfReader(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("invalid PDF format: %w", err)
	}
	return nil
}

// PageContent represents the content of a PDF page
type PageContent struct {
	PageNumber int    `json:"page_number"`
	Content    string `json:"content"`
}
