package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"ai-recruiter-assistant/internal/gmail"
	"ai-recruiter-assistant/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := "postgres://postgres:postgres@localhost:5432/ai_recruiter?sslmode=disable"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize storage (nil for redis as it's not needed for this script)
	store := storage.NewStorage(pool, nil)
	gmailSvc := gmail.NewGmailService(
		os.Getenv("GMAIL_CLIENT_ID"),
		os.Getenv("GMAIL_CLIENT_SECRET"),
		os.Getenv("GMAIL_REDIRECT_URL"),
		store,
	)

	// Fetch emails with empty bodies
	rows, err := pool.Query(ctx, "SELECT id, user_id, gmail_id FROM emails WHERE (body = '' OR body IS NULL) AND gmail_id NOT LIKE 'mock_%'")
	if err != nil {
		fmt.Printf("Error querying: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	type EmailRef struct {
		ID      string
		UserID  string
		GmailID string
	}
	var toRepair []EmailRef
	for rows.Next() {
		var e EmailRef
		if err := rows.Scan(&e.ID, &e.UserID, &e.GmailID); err == nil {
			toRepair = append(toRepair, e)
		}
	}

	fmt.Printf("Found %d emails to repair\n", len(toRepair))

	for _, e := range toRepair {
		fmt.Printf("Repairing email %s (GmailID: %s)...\n", e.ID, e.GmailID)
		
		srv, err := gmailSvc.GetGmailClient(ctx, e.UserID)
		if err != nil {
			fmt.Printf("  Error getting client for user %s: %v\n", e.UserID, err)
			continue
		}

		fullMsg, err := srv.Users.Messages.Get("me", e.GmailID).Format("full").Do()
		if err != nil {
			fmt.Printf("  Error fetching message from Gmail: %v\n", err)
			continue
		}

		_, _, _, body, err := gmail.ParseMessageContent(fullMsg)
		if err != nil {
			fmt.Printf("  Error parsing content: %v\n", err)
			continue
		}

		if body != "" {
			_, err = pool.Exec(ctx, "UPDATE emails SET body = $1 WHERE id = $2", body, e.ID)
			if err != nil {
				fmt.Printf("  Error updating database: %v\n", err)
			} else {
				fmt.Printf("  Successfully repaired! Body length: %d\n", len(body))
			}
		} else {
			fmt.Printf("  Parsed body is still empty.\n")
		}
		
		// Rate limiting sleep
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("Repair process finished.")
}
