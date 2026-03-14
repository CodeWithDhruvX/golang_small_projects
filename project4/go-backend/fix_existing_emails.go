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

	// Update existing emails with NULL gmail_id to have default values
	updateSQL := `
		UPDATE emails 
		SET gmail_id = 'imported_' || id::text 
		WHERE gmail_id IS NULL
	`

	result, err := conn.Exec(context.Background(), updateSQL)
	if err != nil {
		log.Fatalf("Failed to update existing emails: %v", err)
	}

	fmt.Printf("Updated %d emails with default gmail_id values\n", result.RowsAffected())

	// Verify the update
	var count int
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM emails WHERE gmail_id IS NULL").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to count NULL gmail_id emails: %v", err)
	}

	fmt.Printf("Remaining emails with NULL gmail_id: %d\n", count)

	if count == 0 {
		fmt.Println("All emails have been updated successfully!")
	} else {
		fmt.Println("Some emails still have NULL gmail_id values")
	}
}
