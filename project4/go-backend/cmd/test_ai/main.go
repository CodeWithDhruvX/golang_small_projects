package main

import (
	"context"
	"fmt"

	"ai-recruiter-assistant/internal/ai"
)

func main() {
	ollamaURL := "http://localhost:11434"
	osvc := ai.NewOllamaService(ollamaURL)

	ctx := context.Background()
	emailText := "Hi, I am looking for a Senior Go Developer. Are you interested?"
	candidateInfo := "Name: Dhruv Shah, Experience: 5 years in Go"

	fmt.Println("Testing GenerateReply with llama3.1:8b...")
	reply, err := osvc.GenerateReply(ctx, emailText, candidateInfo, "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Generated Reply:\n%s\n", reply)
}
