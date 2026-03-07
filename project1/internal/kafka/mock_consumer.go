package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type MockConsumer struct{}

func NewMockConsumer() (*MockConsumer, error) {
	return &MockConsumer{}, nil
}

func (c *MockConsumer) StartMockConsumer(ctx context.Context, topic string, processMessage func(*sarama.ConsumerMessage) error) error {
	fmt.Println("[MOCK] Consumer started - simulating message consumption")
	
	// Simulate periodic message processing
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[MOCK] Consumer stopped")
			return nil
		case <-ticker.C:
			// Simulate receiving a message
			mockOrder := map[string]interface{}{
				"id":           "mock-order-" + fmt.Sprintf("%d", time.Now().Unix()),
				"customer_id":  "mock-customer",
				"items":        []map[string]interface{}{{"product_id": "mock-product", "quantity": 1, "price": 99.99}},
				"total_amount": 99.99,
				"status":       "PENDING",
				"created_at":   time.Now().UTC().Format(time.RFC3339),
			}
			
			messageBytes, _ := json.Marshal(mockOrder)
			mockMsg := &sarama.ConsumerMessage{
				Topic: topic,
				Value: messageBytes,
			}
			
			if err := processMessage(mockMsg); err != nil {
				fmt.Printf("[MOCK] Error processing message: %v\n", err)
			}
		}
	}
}
