package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"ai-recruiter-assistant/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// CreateUser creates a new user
func (s *Storage) CreateUser(user *auth.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, experience, skills, current_salary, 
			expected_salary, notice_period, location, linkedin_url, github_url, resume_path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	
	_, err := s.db.Exec(context.Background(), query,
		user.ID, user.Email, user.PasswordHash, user.Name, user.Experience, user.Skills,
		user.CurrentSalary, user.ExpectedSalary, user.NoticePeriod, user.Location,
		user.LinkedInURL, user.GitHubURL, user.ResumePath, user.CreatedAt, user.UpdatedAt)
	
	return err
}

// GetUserByEmail retrieves a user by email
func (s *Storage) GetUserByEmail(email string) (*auth.User, error) {
	query := `
		SELECT id, email, password_hash, name, experience, skills, current_salary,
			expected_salary, notice_period, location, linkedin_url, github_url, resume_path, created_at, updated_at
		FROM users WHERE email = $1
	`
	
	var user auth.User
	var skills sql.NullString
	var experience sql.NullString
	var currentSalary sql.NullFloat64
	var expectedSalary sql.NullFloat64
	var noticePeriod sql.NullInt32
	var location sql.NullString
	var linkedinURL sql.NullString
	var githubURL sql.NullString
	var resumePath sql.NullString
	
	err := s.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &experience, &skills,
		&currentSalary, &expectedSalary, &noticePeriod, &location,
		&linkedinURL, &githubURL, &resumePath, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	user.Experience = experience.String
	user.CurrentSalary = currentSalary.Float64
	user.ExpectedSalary = expectedSalary.Float64
	user.NoticePeriod = int(noticePeriod.Int32)
	user.Location = location.String
	user.LinkedInURL = linkedinURL.String
	user.GitHubURL = githubURL.String
	user.ResumePath = resumePath.String

	// Parse skills array
	if skills.Valid {
		var skillArray []string
		if err := json.Unmarshal([]byte(skills.String), &skillArray); err == nil {
			user.Skills = skillArray
		}
	}
	
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *Storage) GetUserByID(id string) (*auth.User, error) {
	query := `
		SELECT id, email, password_hash, name, experience, skills, current_salary,
			expected_salary, notice_period, location, linkedin_url, github_url, resume_path, created_at, updated_at
		FROM users WHERE id = $1
	`
	
	var user auth.User
	var skills sql.NullString
	var experience sql.NullString
	var currentSalary sql.NullFloat64
	var expectedSalary sql.NullFloat64
	var noticePeriod sql.NullInt32
	var location sql.NullString
	var linkedinURL sql.NullString
	var githubURL sql.NullString
	var resumePath sql.NullString
	
	err := s.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &experience, &skills,
		&currentSalary, &expectedSalary, &noticePeriod, &location,
		&linkedinURL, &githubURL, &resumePath, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	user.Experience = experience.String
	user.CurrentSalary = currentSalary.Float64
	user.ExpectedSalary = expectedSalary.Float64
	user.NoticePeriod = int(noticePeriod.Int32)
	user.Location = location.String
	user.LinkedInURL = linkedinURL.String
	user.GitHubURL = githubURL.String
	user.ResumePath = resumePath.String

	// Parse skills array
	if skills.Valid {
		var skillArray []string
		if err := json.Unmarshal([]byte(skills.String), &skillArray); err == nil {
			user.Skills = skillArray
		}
	}
	
	return &user, nil
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(user *auth.User) error {
	query := `
		UPDATE users SET name = $2, experience = $3, skills = $4, current_salary = $5,
			expected_salary = $6, notice_period = $7, location = $8, linkedin_url = $9,
			github_url = $10, resume_path = $11, updated_at = $12
		WHERE id = $1
	`
	
	skillsJSON, _ := json.Marshal(user.Skills)
	
	_, err := s.db.Exec(context.Background(), query,
		user.ID, user.Name, user.Experience, skillsJSON, user.CurrentSalary,
		user.ExpectedSalary, user.NoticePeriod, user.Location, user.LinkedInURL,
		user.GitHubURL, user.ResumePath, user.UpdatedAt)
	
	return err
}

// CreateEmail creates a new email
func (s *Storage) CreateEmail(email *Email) (*Email, error) {
	if email.ID == "" {
		email.ID = uuid.New().String()
	}
	if email.CreatedAt.IsZero() {
		email.CreatedAt = time.Now()
	}
	if email.UpdatedAt.IsZero() {
		email.UpdatedAt = time.Now()
	}
	
	query := `
		INSERT INTO emails (id, user_id, subject, body, sender_email, sender_name, 
			is_recruiter, processed, gmail_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (gmail_id) DO NOTHING
		RETURNING id, user_id, subject, body, sender_email, sender_name, 
			is_recruiter, processed, gmail_id, created_at, updated_at
	`
	
	var createdEmail Email
	err := s.db.QueryRow(context.Background(), query,
		email.ID, email.UserID, email.Subject, email.Body, email.SenderEmail,
		email.SenderName, email.IsRecruiter, email.Processed, email.GmailID, email.CreatedAt, email.UpdatedAt).Scan(
		&createdEmail.ID, &createdEmail.UserID, &createdEmail.Subject, &createdEmail.Body,
		&createdEmail.SenderEmail, &createdEmail.SenderName, &createdEmail.IsRecruiter,
		&createdEmail.Processed, &createdEmail.GmailID, &createdEmail.CreatedAt, &createdEmail.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &createdEmail, nil
}

// GetEmailByID retrieves an email by ID
func (s *Storage) GetEmailByID(id, userID string) (*Email, error) {
	query := `
		SELECT id, user_id, subject, body, sender_email, sender_name, 
			is_recruiter, processed, COALESCE(gmail_id, '') as gmail_id, created_at, updated_at
		FROM emails WHERE id = $1 AND user_id = $2
	`
	
	var email Email
	err := s.db.QueryRow(context.Background(), query, id, userID).Scan(
		&email.ID, &email.UserID, &email.Subject, &email.Body,
		&email.SenderEmail, &email.SenderName, &email.IsRecruiter,
		&email.Processed, &email.GmailID, &email.CreatedAt, &email.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("email not found")
		}
		return nil, err
	}
	
	return &email, nil
}

// GetUserEmails retrieves emails for a user with pagination
func (s *Storage) GetUserEmails(userID string, page, limit int) ([]Email, error) {
	offset := (page - 1) * limit
	
	query := `
		SELECT id, user_id, subject, body, sender_email, sender_name, 
			is_recruiter, processed, COALESCE(gmail_id, '') as gmail_id, created_at, updated_at
		FROM emails WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := s.db.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var emails []Email
	for rows.Next() {
		var email Email
		err := rows.Scan(
			&email.ID, &email.UserID, &email.Subject, &email.Body,
			&email.SenderEmail, &email.SenderName, &email.IsRecruiter,
			&email.Processed, &email.GmailID, &email.CreatedAt, &email.UpdatedAt)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	
	return emails, nil
}

// GetUserEmailsByDateRange retrieves emails for a user with date range filtering and pagination
func (s *Storage) GetUserEmailsByDateRange(userID string, page, limit int, startDate, endDate time.Time) ([]Email, error) {
	offset := (page - 1) * limit
	
	query := `
		SELECT id, user_id, subject, body, sender_email, sender_name, 
			is_recruiter, processed, COALESCE(gmail_id, '') as gmail_id, created_at, updated_at
		FROM emails 
		WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`
	
	rows, err := s.db.Query(context.Background(), query, userID, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var emails []Email
	for rows.Next() {
		var email Email
		err := rows.Scan(
			&email.ID, &email.UserID, &email.Subject, &email.Body,
			&email.SenderEmail, &email.SenderName, &email.IsRecruiter,
			&email.Processed, &email.GmailID, &email.CreatedAt, &email.UpdatedAt)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	
	return emails, nil
}

// UpdateEmail updates an existing email
func (s *Storage) UpdateEmail(email *Email) error {
	query := `
		UPDATE emails SET subject = $2, body = $3, sender_email = $4, sender_name = $5,
			is_recruiter = $6, processed = $7, updated_at = $8
		WHERE id = $1 AND user_id = $9
	`
	
	_, err := s.db.Exec(context.Background(), query,
		email.ID, email.Subject, email.Body, email.SenderEmail, email.SenderName,
		email.IsRecruiter, email.Processed, email.UpdatedAt, email.UserID)
	
	return err
}

// DeleteEmail deletes an email
func (s *Storage) DeleteEmail(id, userID string) error {
	query := "DELETE FROM emails WHERE id = $1 AND user_id = $2"
	
	result, err := s.db.Exec(context.Background(), query, id, userID)
	if err != nil {
		return err
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("email not found")
	}
	
	return nil
}

// CreateApplication creates a new application
func (s *Storage) CreateApplication(application *Application) (*Application, error) {
	if application.ID == "" {
		application.ID = uuid.New().String()
	}
	if application.CreatedAt.IsZero() {
		application.CreatedAt = time.Now()
	}
	if application.UpdatedAt.IsZero() {
		application.UpdatedAt = time.Now()
	}
	
	query := `
		INSERT INTO applications (id, user_id, company, role, recruiter_email, recruiter_name, 
			status, email_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, user_id, company, role, recruiter_email, recruiter_name, 
			status, email_id, created_at, updated_at
	`
	
	var createdApplication Application
	err := s.db.QueryRow(context.Background(), query,
		application.ID, application.UserID, application.Company, application.Role,
		application.RecruiterEmail, application.RecruiterName, application.Status,
		application.EmailID, application.CreatedAt, application.UpdatedAt).Scan(
		&createdApplication.ID, &createdApplication.UserID, &createdApplication.Company,
		&createdApplication.Role, &createdApplication.RecruiterEmail, &createdApplication.RecruiterName,
		&createdApplication.Status, &createdApplication.EmailID, &createdApplication.CreatedAt, &createdApplication.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &createdApplication, nil
}

// GetApplicationByID retrieves an application by ID
func (s *Storage) GetApplicationByID(id, userID string) (*Application, error) {
	query := `
		SELECT id, user_id, company, role, recruiter_email, recruiter_name, 
			status, email_id, created_at, updated_at
		FROM applications WHERE id = $1 AND user_id = $2
	`
	
	var application Application
	err := s.db.QueryRow(context.Background(), query, id, userID).Scan(
		&application.ID, &application.UserID, &application.Company, &application.Role,
		&application.RecruiterEmail, &application.RecruiterName, &application.Status,
		&application.EmailID, &application.CreatedAt, &application.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("application not found")
		}
		return nil, err
	}
	
	return &application, nil
}

// GetUserApplications retrieves applications for a user with pagination and filtering
func (s *Storage) GetUserApplications(userID string, page, limit int, status string) ([]Application, error) {
	offset := (page - 1) * limit
	
	var query string
	var args []interface{}
	
	if status != "" {
		query = `
			SELECT id, user_id, company, role, recruiter_email, recruiter_name, 
				status, email_id, created_at, updated_at
			FROM applications WHERE user_id = $1 AND status = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{userID, status, limit, offset}
	} else {
		query = `
			SELECT id, user_id, company, role, recruiter_email, recruiter_name, 
				status, email_id, created_at, updated_at
			FROM applications WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{userID, limit, offset}
	}
	
	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var applications []Application
	for rows.Next() {
		var application Application
		err := rows.Scan(
			&application.ID, &application.UserID, &application.Company, &application.Role,
			&application.RecruiterEmail, &application.RecruiterName, &application.Status,
			&application.EmailID, &application.CreatedAt, &application.UpdatedAt)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}
	
	return applications, nil
}

// UpdateApplication updates an existing application
func (s *Storage) UpdateApplication(application *Application) (*Application, error) {
	query := `
		UPDATE applications SET company = $2, role = $3, recruiter_email = $4, recruiter_name = $5,
			status = $6, email_id = $7, updated_at = $8
		WHERE id = $1 AND user_id = $9
		RETURNING id, user_id, company, role, recruiter_email, recruiter_name, 
			status, email_id, created_at, updated_at
	`
	
	var updatedApplication Application
	err := s.db.QueryRow(context.Background(), query,
		application.ID, application.Company, application.Role, application.RecruiterEmail,
		application.RecruiterName, application.Status, application.EmailID,
		application.UpdatedAt, application.UserID).Scan(
		&updatedApplication.ID, &updatedApplication.UserID, &updatedApplication.Company,
		&updatedApplication.Role, &updatedApplication.RecruiterEmail, &updatedApplication.RecruiterName,
		&updatedApplication.Status, &updatedApplication.EmailID, &updatedApplication.CreatedAt, &updatedApplication.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("application not found")
		}
		return nil, err
	}
	
	return &updatedApplication, nil
}

// DeleteApplication deletes an application
func (s *Storage) DeleteApplication(id, userID string) error {
	query := "DELETE FROM applications WHERE id = $1 AND user_id = $2"
	
	result, err := s.db.Exec(context.Background(), query, id, userID)
	if err != nil {
		return err
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("application not found")
	}
	
	return nil
}

// CheckDuplicateApplication checks if a duplicate application exists
func (s *Storage) CheckDuplicateApplication(userID, company, recruiterEmail string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM applications 
		WHERE user_id = $1 AND company = $2 AND recruiter_email = $3)
	`
	
	var exists bool
	err := s.db.QueryRow(context.Background(), query, userID, company, recruiterEmail).Scan(&exists)
	
	return exists, err
}

// Placeholder implementations for remaining methods
func (s *Storage) CreateDocument(document *Document) error {
	if document.ID == "" {
		document.ID = uuid.New().String()
	}
	if document.CreatedAt.IsZero() {
		document.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO documents (id, user_id, content, source, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := s.db.Exec(context.Background(), query,
		document.ID, document.UserID, document.Content, document.Source, document.CreatedAt)
	
	return err
}

func (s *Storage) GetDocumentsByUserID(userID string, source string) ([]Document, error) {
	var query string
	var args []interface{}
	
	if source != "" {
		query = "SELECT id, user_id, content, source, created_at FROM documents WHERE user_id = $1 AND source = $2 ORDER BY created_at DESC"
		args = []interface{}{userID, source}
	} else {
		query = "SELECT id, user_id, content, source, created_at FROM documents WHERE user_id = $1 ORDER BY created_at DESC"
		args = []interface{}{userID}
	}
	
	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var documents []Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.ID, &doc.UserID, &doc.Content, &doc.Source, &doc.CreatedAt)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	
	return documents, nil
}

func (s *Storage) SearchDocuments(userID, query string, topK int) ([]Document, error) {
	// Simple text search fallback if no embeddings are ready
	sqlQuery := `
		SELECT id, user_id, content, source, created_at
		FROM documents
		WHERE user_id = $1 AND content ILIKE $2
		LIMIT $3
	`
	
	rows, err := s.db.Query(context.Background(), sqlQuery, userID, "%"+query+"%", topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var documents []Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.ID, &doc.UserID, &doc.Content, &doc.Source, &doc.CreatedAt)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	
	return documents, nil
}

func (s *Storage) DeleteDocument(id, userID string) error {
	query := "DELETE FROM documents WHERE id = $1 AND user_id = $2"
	_, err := s.db.Exec(context.Background(), query, id, userID)
	return err
}

func (s *Storage) CreateAIReply(reply *AIReply) error {
	if reply.ID == "" {
		reply.ID = uuid.New().String()
	}
	if reply.CreatedAt.IsZero() {
		reply.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO ai_replies (id, user_id, email_id, reply_content, model_used, response_time_ms, is_sent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := s.db.Exec(context.Background(), query,
		reply.ID, reply.UserID, reply.EmailID, reply.ReplyContent, reply.ModelUsed,
		reply.ResponseTime, reply.IsSent, reply.CreatedAt)
	
	return err
}

func (s *Storage) GetAIRepliesByUserID(userID string) ([]AIReply, error) {
	query := "SELECT id, user_id, email_id, reply_content, model_used, response_time_ms, is_sent, created_at FROM ai_replies WHERE user_id = $1 ORDER BY created_at DESC"
	
	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var replies []AIReply
	for rows.Next() {
		var r AIReply
		err := rows.Scan(&r.ID, &r.UserID, &r.EmailID, &r.ReplyContent, &r.ModelUsed, &r.ResponseTime, &r.IsSent, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		replies = append(replies, r)
	}
	
	return replies, nil
}

func (s *Storage) GetAIRepliesByEmailID(emailID, userID string) ([]AIReply, error) {
	query := "SELECT id, user_id, email_id, reply_content, model_used, response_time_ms, is_sent, created_at FROM ai_replies WHERE email_id = $1 AND user_id = $2 ORDER BY created_at DESC"
	
	rows, err := s.db.Query(context.Background(), query, emailID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var replies []AIReply
	for rows.Next() {
		var r AIReply
		err := rows.Scan(&r.ID, &r.UserID, &r.EmailID, &r.ReplyContent, &r.ModelUsed, &r.ResponseTime, &r.IsSent, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		replies = append(replies, r)
	}
	
	return replies, nil
}

func (s *Storage) UpdateAIReply(reply *AIReply) error {
	query := "UPDATE ai_replies SET is_sent = $1 WHERE id = $2 AND user_id = $3"
	_, err := s.db.Exec(context.Background(), query, reply.IsSent, reply.ID, reply.UserID)
	return err
}

func (s *Storage) StoreEmbedding(id string, embedding []float64, table string) error {
	// Simple validation to prevent SQL injection for table name
	switch table {
	case "emails", "documents":
		// OK
	default:
		return fmt.Errorf("invalid table name: %s", table)
	}

	query := fmt.Sprintf("UPDATE %s SET embedding = $1 WHERE id = $2", table)
	_, err := s.db.Exec(context.Background(), query, embedding, id)
	return err
}

func (s *Storage) SearchSimilar(embedding []float64, userID string, topK int, table string) ([]Document, error) {
	// Simple validation to prevent SQL injection for table name
	switch table {
	case "emails", "documents":
		// OK
	default:
		return nil, fmt.Errorf("invalid table name: %s", table)
	}

	// Double check if pgvector extension is available
	var exists bool
	err := s.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&exists)
	if err != nil || !exists {
		logrus.Warn("Vector extension not available, falling back to empty results")
		return []Document{}, nil
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, content, source, created_at
		FROM %s
		WHERE user_id = $1
		ORDER BY embedding <=> $2
		LIMIT $3
	`, table)

	// Note: If searching in emails table, we might need to map columns
	if table == "emails" {
		query = `
			SELECT id, user_id, body as content, 'email' as source, created_at
			FROM emails
			WHERE user_id = $1
			ORDER BY embedding <=> $2
			LIMIT $3
		`
	}

	rows, err := s.db.Query(context.Background(), query, userID, embedding, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.ID, &doc.UserID, &doc.Content, &doc.Source, &doc.CreatedAt)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

// CreateGmailIntegration creates a new Gmail integration
func (s *Storage) CreateGmailIntegration(integration *GmailIntegration) error {
	if integration.ID == "" {
		integration.ID = uuid.New().String()
	}
	query := `
		INSERT INTO gmail_integrations (id, user_id, access_token, refresh_token, token_expiry, 
			email, is_active, last_sync_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err := s.db.Exec(context.Background(), query,
		integration.ID, integration.UserID, integration.AccessToken, integration.RefreshToken,
		integration.TokenExpiry, integration.Email, integration.IsActive, integration.LastSyncAt,
		integration.CreatedAt, integration.UpdatedAt)
	
	return err
}

// GetGmailIntegration retrieves Gmail integration by user ID
func (s *Storage) GetGmailIntegration(userID string) (*GmailIntegration, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, token_expiry, email, 
			is_active, last_sync_at, created_at, updated_at
		FROM gmail_integrations WHERE user_id = $1
	`
	
	var integration GmailIntegration
	err := s.db.QueryRow(context.Background(), query, userID).Scan(
		&integration.ID, &integration.UserID, &integration.AccessToken, &integration.RefreshToken,
		&integration.TokenExpiry, &integration.Email, &integration.IsActive,
		&integration.LastSyncAt, &integration.CreatedAt, &integration.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("gmail integration not found")
		}
		return nil, err
	}
	
	return &integration, nil
}

// UpdateGmailIntegration updates an existing Gmail integration
func (s *Storage) UpdateGmailIntegration(integration *GmailIntegration) error {
	query := `
		UPDATE gmail_integrations 
		SET access_token = $2, refresh_token = $3, token_expiry = $4, email = $5,
			is_active = $6, last_sync_at = $7, updated_at = $8
		WHERE id = $1
	`
	
	_, err := s.db.Exec(context.Background(), query,
		integration.ID, integration.AccessToken, integration.RefreshToken,
		integration.TokenExpiry, integration.Email, integration.IsActive,
		integration.LastSyncAt, integration.UpdatedAt)
	
	return err
}

// DeleteGmailIntegration deletes a Gmail integration
func (s *Storage) DeleteGmailIntegration(userID string) error {
	query := `DELETE FROM gmail_integrations WHERE user_id = $1`
	_, err := s.db.Exec(context.Background(), query, userID)
	return err
}
