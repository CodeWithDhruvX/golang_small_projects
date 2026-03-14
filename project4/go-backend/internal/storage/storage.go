package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Storage handles all database operations
type Storage struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

// NewStorage creates a new storage instance
func NewStorage(db *pgxpool.Pool, redis *redis.Client) *Storage {
	return &Storage{
		db:    db,
		redis: redis,
	}
}

// NewDatabase creates a new database connection
func NewDatabase(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Successfully connected to database")
	return pool, nil
}

// NewRedisClient creates a new Redis client
func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		// Fallback to default Redis URL if parsing fails
		opt = &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		logrus.Warnf("Failed to connect to Redis: %v. Continuing without cache.", err)
		return nil, err
	}

	logrus.Info("Successfully connected to Redis")
	return client, nil
}

// RunMigrations runs database migrations
func (s *Storage) RunMigrations() error {
	logrus.Info("Running database migrations...")

	// Check if vector extension exists
	var exists bool
	err := s.db.QueryRow(context.Background(), 
		"SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check vector extension: %w", err)
	}

	if !exists {
		logrus.Warn("Vector extension not found. Please ensure PostgreSQL has PGVector extension installed.")
	}

	// Create tables if they don't exist
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			experience TEXT,
			skills TEXT[],
			current_salary DECIMAL(10,2),
			expected_salary DECIMAL(10,2),
			notice_period INTEGER,
			location VARCHAR(255),
			linkedin_url VARCHAR(500),
			github_url VARCHAR(500),
			resume_path VARCHAR(500),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS emails (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			subject VARCHAR(1000) NOT NULL,
			body TEXT NOT NULL,
			sender_email VARCHAR(255) NOT NULL,
			sender_name VARCHAR(255),
			is_recruiter BOOLEAN DEFAULT FALSE,
			processed BOOLEAN DEFAULT FALSE,
			gmail_id VARCHAR(255) UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS applications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			company VARCHAR(255) NOT NULL,
			role VARCHAR(255) NOT NULL,
			recruiter_email VARCHAR(255) NOT NULL,
			recruiter_name VARCHAR(255),
			status VARCHAR(50) DEFAULT 'Applied',
			email_id UUID REFERENCES emails(id) ON DELETE SET NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(user_id, company, recruiter_email)
		)`,
		`CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			source VARCHAR(100) NOT NULL,
			metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS ai_replies (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			email_id UUID NOT NULL REFERENCES emails(id) ON DELETE CASCADE,
			reply_content TEXT NOT NULL,
			model_used VARCHAR(100),
			tokens_used INTEGER,
			response_time_ms INTEGER,
			is_sent BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS email_processing_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email_id UUID REFERENCES emails(id) ON DELETE CASCADE,
			processing_step VARCHAR(100) NOT NULL,
			status VARCHAR(50) NOT NULL,
			message TEXT,
			metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS gmail_integrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			access_token TEXT NOT NULL,
			refresh_token TEXT,
			token_expiry TIMESTAMP WITH TIME ZONE,
			email VARCHAR(255),
			is_active BOOLEAN DEFAULT TRUE,
			last_sync_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(user_id)
		)`,
	}

	for _, query := range queries {
		_, err := s.db.Exec(context.Background(), query)
		if err != nil {
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_emails_user_id ON emails(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_emails_sender_email ON emails(sender_email)",
		"CREATE INDEX IF NOT EXISTS idx_emails_gmail_id ON emails(gmail_id)",
		"CREATE INDEX IF NOT EXISTS idx_emails_created_at ON emails(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_applications_user_id ON applications(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status)",
		"CREATE INDEX IF NOT EXISTS idx_documents_user_id ON documents(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_documents_source ON documents(source)",
		"CREATE INDEX IF NOT EXISTS idx_ai_replies_user_id ON ai_replies(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_ai_replies_email_id ON ai_replies(email_id)",
		"CREATE INDEX IF NOT EXISTS idx_gmail_integrations_user_id ON gmail_integrations(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_gmail_integrations_email ON gmail_integrations(email)",
	}

	for _, index := range indexes {
		_, err := s.db.Exec(context.Background(), index)
		if err != nil {
			logrus.Warnf("Failed to create index: %v", err)
		}
	}

	// Create vector indexes if vector extension is available
	if exists {
		// Try to add embedding columns if vector extension is available
		embeddingColumns := []string{
			"ALTER TABLE emails ADD COLUMN IF NOT EXISTS embedding vector(768)",
			"ALTER TABLE documents ADD COLUMN IF NOT EXISTS embedding vector(768)",
		}

		for _, column := range embeddingColumns {
			_, err := s.db.Exec(context.Background(), column)
			if err != nil {
				logrus.Warnf("Failed to add embedding column: %v", err)
			}
		}

		vectorIndexes := []string{
			"CREATE INDEX IF NOT EXISTS idx_emails_embedding ON emails USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100)",
			"CREATE INDEX IF NOT EXISTS idx_documents_embedding ON documents USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100)",
		}

		for _, index := range vectorIndexes {
			_, err := s.db.Exec(context.Background(), index)
			if err != nil {
				logrus.Warnf("Failed to create vector index: %v", err)
			}
		}
	}

	// Fix existing emails with NULL gmail_id
	logrus.Info("Fixing existing emails with NULL gmail_id...")
	result, err := s.db.Exec(context.Background(), 
		"UPDATE emails SET gmail_id = 'imported_' || id::text WHERE gmail_id IS NULL")
	if err != nil {
		logrus.Warnf("Failed to fix NULL gmail_id values: %v", err)
	} else {
		logrus.Infof("Fixed %d emails with NULL gmail_id", result.RowsAffected())
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}

// Close closes the storage connections
func (s *Storage) Close() {
	if s.redis != nil {
		s.redis.Close()
	}
	if s.db != nil {
		s.db.Close()
	}
	logrus.Info("Storage connections closed")
}

// GetDB returns the database connection pool
func (s *Storage) GetDB() *pgxpool.Pool {
	return s.db
}

// GetRedis returns the Redis client
func (s *Storage) GetRedis() *redis.Client {
	return s.redis
}
