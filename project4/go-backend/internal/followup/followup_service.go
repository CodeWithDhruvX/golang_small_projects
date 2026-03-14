package followup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-recruiter-assistant/internal/auth"
	"github.com/sirupsen/logrus"
	"ai-recruiter-assistant/internal/ai"
	"ai-recruiter-assistant/internal/storage"
)

// FollowUpService handles follow-up email generation
type FollowUpService struct {
	storage storage.StorageInterface
	ollama  *ai.OllamaService
}

// NewFollowUpService creates a new follow-up service
func NewFollowUpService(storage storage.StorageInterface, ollama *ai.OllamaService) *FollowUpService {
	return &FollowUpService{
		storage: storage,
		ollama:  ollama,
	}
}

// FollowUpRequest represents a follow-up email request
type FollowUpRequest struct {
	ApplicationID string `json:"application_id"`
	DaysSinceLast int    `json:"days_since_last"`
	Type          string `json:"type"` // "gentle_reminder", "status_inquiry", "final_followup"
}

// FollowUpEmail represents a generated follow-up email
type FollowUpEmail struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Type    string `json:"type"`
}

// GenerateFollowUpEmail generates a follow-up email for a pending application
func (fs *FollowUpService) GenerateFollowUpEmail(ctx context.Context, userID string, req FollowUpRequest) (*FollowUpEmail, error) {
	logrus.Infof("Generating follow-up email for user: %s, application: %s", userID, req.ApplicationID)

	// Get application details
	application, err := fs.storage.GetApplicationByID(req.ApplicationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Get user profile for context
	user, err := fs.storage.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Build context for AI generation
	context := fs.buildFollowUpContext(user, application, req)

	// Generate follow-up email using AI
	prompt := fs.buildFollowUpPrompt(context, req.Type)

	response, err := fs.ollama.GenerateText(ctx, "llama3.1:8b", prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate follow-up email: %w", err)
	}

	// Parse the AI response to extract subject and body
	subject, body := fs.parseAIResponse(response.Response)

	followUpEmail := &FollowUpEmail{
		Subject: subject,
		Body:    body,
		Type:    req.Type,
	}

	logrus.Infof("Generated follow-up email: %s", subject)
	return followUpEmail, nil
}

// buildFollowUpContext builds context for follow-up generation
func (fs *FollowUpService) buildFollowUpContext(user *auth.User, application *storage.Application, req FollowUpRequest) string {
	return fmt.Sprintf(`
Candidate: %s
Position Applied: %s at %s
Recruiter: %s (%s)
Days Since Last Contact: %d
Application Status: %s

Candidate Details:
- Experience: %s
- Expected Salary: %.2f
- Notice Period: %d days
- Location: %s
`,
		user.Name,
		application.Role,
		application.Company,
		application.RecruiterName,
		application.RecruiterEmail,
		req.DaysSinceLast,
		application.Status,
		user.Experience,
		user.ExpectedSalary,
		user.NoticePeriod,
		user.Location,
	)
}

// buildFollowUpPrompt builds the AI prompt for follow-up generation
func (fs *FollowUpService) buildFollowUpPrompt(context, followUpType string) string {
	basePrompt := fmt.Sprintf(`You are helping a job candidate generate a professional follow-up email. Use the following context:

%s

Generate a follow-up email that is:`, context)

	switch followUpType {
	case "gentle_reminder":
		return basePrompt + `
- Polite and professional
- Brief (2-3 sentences max)
- Shows continued interest
- Not pushy or demanding
- Appropriate for 3-7 days since last contact

Format your response as:
Subject: [brief subject line]
[Email body]`

	case "status_inquiry":
		return basePrompt + `
- Professional and respectful
- Asks for update on application status
- Shows continued enthusiasm
- Appropriate for 1-2 weeks since last contact

Format your response as:
Subject: [brief subject line]
[Email body]`

	case "final_followup":
		return basePrompt + `
- Professional and courteous
- Indicates this is a final follow-up
- Leaves door open for future opportunities
- Appropriate for 2+ weeks since last contact

Format your response as:
Subject: [brief subject line]
[Email body]`

	default:
		return basePrompt + `
- Professional and appropriate
- Based on the context provided

Format your response as:
Subject: [brief subject line]
[Email body]`
	}
}

// parseAIResponse parses the AI response to extract subject and body
func (fs *FollowUpService) parseAIResponse(response string) (string, string) {
	lines := strings.Split(response, "\n")
	
	var subject string
	var bodyLines []string
	inBody := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		if !inBody {
			if len(line) > 8 && line[:8] == "Subject:" {
				subject = line[8:]
				continue
			}
			if subject != "" {
				inBody = true
				bodyLines = append(bodyLines, line)
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}
	
	// If no subject found, use default
	if subject == "" {
		subject = "Following up on " + (bodyLines[0])[:50] + "..."
	}
	
	// Join body lines
	body := ""
	for i, line := range bodyLines {
		if i > 0 {
			body += "\n"
		}
		body += line
	}
	
	return subject, body
}

// GetPendingFollowUps returns applications that need follow-up
func (fs *FollowUpService) GetPendingFollowUps(ctx context.Context, userID string) ([]storage.Application, error) {
	// Get all applications for the user
	applications, err := fs.storage.GetUserApplications(userID, 1, 100, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	var pendingFollowUps []storage.Application
	now := time.Now()

	for _, app := range applications {
		// Only consider applications that are "Applied" or "Interview Scheduled"
		if app.Status != "Applied" && app.Status != "Interview Scheduled" {
			continue
		}

		// Check days since last update
		daysSinceLast := int(now.Sub(app.UpdatedAt).Hours() / 24)

		// Different follow-up intervals based on status
		needsFollowUp := false
		switch app.Status {
		case "Applied":
			needsFollowUp = daysSinceLast >= 7 // Follow up after 7 days
		case "Interview Scheduled":
			needsFollowUp = daysSinceLast >= 14 // Follow up after 2 weeks
		}

		if needsFollowUp {
			pendingFollowUps = append(pendingFollowUps, app)
		}
	}

	return pendingFollowUps, nil
}

// ScheduleFollowUp creates a scheduled follow-up reminder
func (fs *FollowUpService) ScheduleFollowUp(ctx context.Context, userID, applicationID string, followUpDate time.Time) error {
	// TODO: Implement follow-up scheduling
	// This could be stored in a separate follow_ups table with reminder dates
	logrus.Infof("Scheduled follow-up for user: %s, application: %s, date: %v", userID, applicationID, followUpDate)
	return nil
}
