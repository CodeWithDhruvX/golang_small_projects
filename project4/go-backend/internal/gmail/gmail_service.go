package gmail

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/sirupsen/logrus"
	"ai-recruiter-assistant/internal/storage"
)

// GmailService handles Gmail API operations
type GmailService struct {
	config      *oauth2.Config
	storage     storage.StorageInterface
	redirectURL string
}

// NewGmailService creates a new Gmail service
func NewGmailService(clientID, clientSecret, redirectURL string, storage storage.StorageInterface) *GmailService {
	// Debug logging
	logrus.Infof("Initializing Gmail service with ClientID: %s, RedirectURL: %s", clientID, redirectURL)
	
	if clientID == "" || clientSecret == "" {
		logrus.Error("Gmail OAuth2 credentials are missing")
		return nil
	}
	
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			gmail.GmailReadonlyScope,
			gmail.GmailSendScope,
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GmailService{
		config:      config,
		storage:     storage,
		redirectURL: redirectURL,
	}
}

// GetConfig returns the OAuth2 config
func (gs *GmailService) GetConfig() *oauth2.Config {
	return gs.config
}

// GetAuthURL generates OAuth2 authorization URL
func (gs *GmailService) GetAuthURL(state string) string {
	return gs.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// ExchangeCodeForToken exchanges authorization code for tokens
func (gs *GmailService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := gs.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}

// StoreToken stores OAuth2 token in database
func (gs *GmailService) StoreToken(ctx context.Context, userID string, email string, token *oauth2.Token) error {
	tokenData, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Store or update Gmail integration
	gmailIntegration := &storage.GmailIntegration{
		UserID:          userID,
		AccessToken:     string(tokenData),
		RefreshToken:    token.RefreshToken,
		TokenExpiry:     token.Expiry,
		Email:           email,
		IsActive:        true,
		LastSyncAt:      time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Check if integration already exists
	existing, err := gs.storage.GetGmailIntegration(userID)
	if err != nil && err.Error() != "gmail integration not found" {
		return fmt.Errorf("failed to check existing integration: %w", err)
	}

	if existing != nil {
		gmailIntegration.ID = existing.ID
		return gs.storage.UpdateGmailIntegration(gmailIntegration)
	}

	return gs.storage.CreateGmailIntegration(gmailIntegration)
}

// GetToken retrieves OAuth2 token from database
func (gs *GmailService) GetToken(ctx context.Context, userID string) (*oauth2.Token, error) {
	integration, err := gs.storage.GetGmailIntegration(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gmail integration: %w", err)
	}

	if !integration.IsActive {
		return nil, fmt.Errorf("gmail integration is not active")
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(integration.AccessToken), &token)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	// Check if token needs refresh
	if token.Expiry.Before(time.Now().Add(5 * time.Minute)) {
		newToken, err := gs.RefreshToken(ctx, &token)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		return newToken, nil
	}

	return &token, nil
}

// RefreshToken refreshes an expired OAuth2 token
func (gs *GmailService) RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := gs.config.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Note: We'll need to update this to get the userID from context
	// For now, this is a placeholder that should be called with proper context
	return newToken, nil
}

// GetGmailClient creates a Gmail API client
func (gs *GmailService) GetGmailClient(ctx context.Context, userID string) (*gmail.Service, error) {
	token, err := gs.GetToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	client := gs.config.Client(ctx, token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create gmail service: %w", err)
	}

	return srv, nil
}

// FetchEmails fetches emails from Gmail
func (gs *GmailService) FetchEmails(ctx context.Context, userID string, maxResults int64) ([]*gmail.Message, error) {
	srv, err := gs.GetGmailClient(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gmail client: %w", err)
	}

	// Calculate date from one week ago
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	query := fmt.Sprintf("after:%s", oneWeekAgo.Format("2006/01/02"))

	logrus.Infof("Fetching emails from last week with query: %s", query)

	// Get messages from inbox from the last week
	msgs, err := srv.Users.Messages.List("me").Q(query).MaxResults(maxResults).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	logrus.Infof("Found %d messages from last week", len(msgs.Messages))

	var messages []*gmail.Message
	for _, msg := range msgs.Messages {
		fullMsg, err := srv.Users.Messages.Get("me", msg.Id).Format("full").Do()
		if err != nil {
			logrus.Errorf("Failed to get message %s: %v", msg.Id, err)
			continue
		}
		messages = append(messages, fullMsg)
	}

	return messages, nil
}

// SendEmail sends an email via Gmail API
func (gs *GmailService) SendEmail(ctx context.Context, userID string, to, subject, body string) error {
	srv, err := gs.GetGmailClient(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get gmail client: %w", err)
	}

	// Create email message
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	msg := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(message)),
	}

	_, err = srv.Users.Messages.Send("me", msg).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	logrus.Infof("Email sent successfully to %s", to)
	return nil
}

// GetUserProfile gets user's Gmail profile
func (gs *GmailService) GetUserProfile(ctx context.Context, userID string) (*gmail.Profile, error) {
	srv, err := gs.GetGmailClient(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gmail client: %w", err)
	}

	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return profile, nil
}

// GenerateState generates a random state string for OAuth2
func GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ParseMessageContent extracts email content from Gmail message
func ParseMessageContent(msg *gmail.Message) (subject, from, to, body string, err error) {
	// Extract headers
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "Subject":
			subject = header.Value
		case "From":
			from = header.Value
		case "To":
			to = header.Value
		}
	}

	// Extract body
	if msg.Payload.Body != nil && msg.Payload.Body.Data != "" {
		bodyData, err := decodeBase64URL(msg.Payload.Body.Data)
		if err != nil {
			logrus.Errorf("Failed to decode simple message body: %v", err)
			return subject, from, to, "", nil
		}
		body = string(bodyData)
	} else {
		// Look for body in parts
		body = extractBodyFromParts(msg.Payload.Parts)
	}

	return subject, from, to, body, nil
}

// extractBodyFromParts extracts body from message parts
func extractBodyFromParts(parts []*gmail.MessagePart) string {
	var htmlBody string
	
	for _, part := range parts {
		if part.MimeType == "text/plain" && part.Body != nil && part.Body.Data != "" {
			bodyData, err := decodeBase64URL(part.Body.Data)
			if err != nil {
				logrus.Errorf("Failed to decode text part body: %v", err)
				continue
			}
			return string(bodyData)
		}
		
		if part.MimeType == "text/html" && part.Body != nil && part.Body.Data != "" {
			bodyData, err := decodeBase64URL(part.Body.Data)
			if err != nil {
				logrus.Errorf("Failed to decode html part body: %v", err)
				continue
			}
			htmlBody = string(bodyData)
		}
		
		// Recursively check nested parts
		if len(part.Parts) > 0 {
			if body := extractBodyFromParts(part.Parts); body != "" {
				return body
			}
		}
	}
	
	// If no plain text body found but html body exists, return stripped html
	if htmlBody != "" {
		return stripHTML(htmlBody)
	}
	
	return ""
}

// decodeBase64URL decodes base64url encoded string (handles both padded and unpadded)
func decodeBase64URL(s string) ([]byte, error) {
	// Add padding if missing
	if l := len(s) % 4; l > 0 {
		s += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(s)
}

// stripHTML removes basic HTML tags from a string
func stripHTML(s string) string {
	// Replace <br> and <div>/ <p> closures with newlines for basic readability
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "</div>", "\n")
	s = strings.ReplaceAll(s, "</p>", "\n\n")

	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	
	// Clean up double newlines and spaces
	content := result.String()
	content = strings.TrimSpace(content)
	
	return content
}

// MarkAsRead marks an email as read
func (gs *GmailService) MarkAsRead(ctx context.Context, userID string, messageID string) error {
	srv, err := gs.GetGmailClient(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get gmail client: %w", err)
	}

	_, err = srv.Users.Messages.Modify("me", messageID, &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()

	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}
