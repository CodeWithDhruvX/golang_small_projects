package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/ai_recruiter?sslmode=disable"
	}

	// Connect to database
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Read migration file
	migrationSQL := `
DO $$
BEGIN
    -- Check if the column already exists
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name='emails' 
        AND column_name='gmail_id'
    ) THEN
        -- Add the gmail_id column
        ALTER TABLE emails ADD COLUMN gmail_id VARCHAR(255);
        
        -- Add unique constraint
        ALTER TABLE emails ADD CONSTRAINT emails_gmail_id_unique UNIQUE (gmail_id);
        
        -- Add index for better performance
        CREATE INDEX IF NOT EXISTS idx_emails_gmail_id ON emails(gmail_id);
        
        RAISE NOTICE 'gmail_id column added to emails table';
    ELSE
        RAISE NOTICE 'gmail_id column already exists in emails table';
    END IF;
END $$;
`

	// Execute migration
	_, err = conn.Exec(context.Background(), migrationSQL)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed successfully!")
}
