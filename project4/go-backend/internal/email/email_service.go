package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
	"ai-recruiter-assistant/internal/storage"
)

// EmailService handles email ingestion and processing
type EmailService struct {
	storage storage.StorageInterface
}

// NewEmailService creates a new email service
func NewEmailService(storage storage.StorageInterface) *EmailService {
	return &EmailService{
		storage: storage,
	}
}

// IMAPConfig holds IMAP server configuration
type IMAPConfig struct {
	Server   string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

// ImportEmailsFromIMAP imports emails from IMAP server
func (es *EmailService) ImportEmailsFromIMAP(ctx context.Context, userID string, config IMAPConfig) error {
	logrus.Infof("Starting IMAP email import for user: %s", userID)

	// Connect to IMAP server
	var c *client.Client
	var err error

	address := fmt.Sprintf("%s:%d", config.Server, config.Port)
	
	if config.UseTLS {
		c, err = client.DialTLS(address, &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		c, err = client.Dial(address)
	}
	
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(config.Username, config.Password); err != nil {
		return fmt.Errorf("failed to login to IMAP server: %w", err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("failed to select INBOX: %w", err)
	}

	logrus.Infof("Connected to mailbox with %d messages", mbox.Messages)

	// Get recent emails (last 30 days)
	since := time.Now().AddDate(0, 0, -30)
	criteria := imap.NewSearchCriteria()
	criteria.Since = since

	uids, err := c.Search(criteria)
	if err != nil {
		return fmt.Errorf("failed to search emails: %w", err)
	}

	if len(uids) == 0 {
		logrus.Info("No recent emails found")
		return nil
	}

	logrus.Infof("Found %d recent emails", len(uids))

	// Process emails in batches
 batchSize := 10
	for i := 0; i < len(uids); i += batchSize {
		end := i + batchSize
		if end > len(uids) {
			end = len(uids)
		}

		batchUIDs := uids[i:end]
		if err := es.processEmailBatch(ctx, c, userID, batchUIDs); err != nil {
			logrus.Errorf("Failed to process email batch %d-%d: %v", i, end, err)
			continue
		}
	}

	logrus.Infof("Successfully imported emails for user: %s", userID)
	return nil
}

// processEmailBatch processes a batch of emails
func (es *EmailService) processEmailBatch(ctx context.Context, c *client.Client, userID string, uids []uint32) error {
	// Set up sequence set
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	// Fetch message bodies
	messages := make(chan *imap.Message, len(uids))
	err := c.Fetch(seqSet, []imap.FetchItem{imap.FetchBody, imap.FetchEnvelope}, messages)
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	processedCount := 0
	for msg := range messages {
		if err := es.processEmailMessage(ctx, userID, msg); err != nil {
			logrus.Errorf("Failed to process email: %v", err)
			continue
		}
		processedCount++
	}

	logrus.Infof("Processed %d/%d emails in batch", processedCount, len(uids))
	return nil
}

// processEmailMessage processes a single email message
func (es *EmailService) processEmailMessage(ctx context.Context, userID string, msg *imap.Message) error {
	// Parse email content
	emailReader := msg.GetBody(&imap.BodySectionName{})
	if emailReader == nil {
		return fmt.Errorf("no email body found")
	}

	// Simple text extraction - for now we'll just read the raw content
	buf := new(strings.Builder)
	_, copyErr := io.Copy(buf, emailReader)
	if copyErr != nil {
		return fmt.Errorf("failed to read email content: %w", copyErr)
	}

	body := buf.String()
	if body == "" {
		return fmt.Errorf("no email body found")
	}

	// Extract email details
	subject := msg.Envelope.Subject
	from := msg.Envelope.From
	date := msg.Envelope.Date
	
	// Parse sender information
	var senderEmail, senderName string
	if from != nil && len(from) > 0 {
		senderEmail = from[0].Address()
		senderName = from[0].PersonalName
	}

	if senderEmail == "" {
		return fmt.Errorf("no sender email found")
	}

	// Create email record
	email := &storage.Email{
		UserID:      userID,
		Subject:     subject,
		Body:        body,
		SenderEmail: senderEmail,
		SenderName:  senderName,
		IsRecruiter: false, // Will be determined by AI classification
		Processed:   false,
	}

	// Store email
	createdEmail, err := es.storage.CreateEmail(email)
	if err != nil {
		return fmt.Errorf("failed to store email: %w", err)
	}

	logrus.Infof("Stored email from %s: %s", senderEmail, subject)
	
	// Parse and set creation date
	if !date.IsZero() {
		createdEmail.CreatedAt = date
		es.storage.UpdateEmail(createdEmail)
	}

	return nil
}

// ProcessEmailForAI processes an email for AI classification and analysis
func (es *EmailService) ProcessEmailForAI(ctx context.Context, emailID string) error {
	logrus.Infof("Processing email for AI analysis: %s", emailID)

	// Get email from storage
	_, err := es.storage.GetEmailByID(emailID, "")
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}

	// TODO: Call AI service for classification and requirement extraction
	// This will be implemented when we integrate the AI services
	
	logrus.Infof("Email processing queued for AI: %s", emailID)
	return nil
}
